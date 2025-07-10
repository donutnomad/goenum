package main

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

type EnumInfo struct {
	Name          string
	Type          string
	Values        []EnumValue
	Comment       string
	Options       EnumOptions
	PackageName   string
	FileName      string
	BaseType      string
	ContainerName string
	AllTags       []string // All unique tags across all values
	ConstBlock    string   // Raw text of the entire const block
}

type EnumValue struct {
	Name            string
	Value           string
	Comment         string
	OriginalComment string   // Raw comment text for display in generated code
	PrecedingLines  []string // All lines (comments, separators, empty lines) before this enum
	Names           []string
	IsInvalid       bool
	Tags            []string
	Transitions     []string
	IsFinal         bool
}

type EnumOptions struct {
	SQL          bool
	JSON         bool
	YAML         bool
	Text         bool
	Binary       bool
	SerdeFormat  string // "name" or "value"
	GenName      bool
	StateMachine bool
}

// FileTemplateData is the data structure passed to the template for generating the output file.
type FileTemplateData struct {
	PackageName string
	Enums       []EnumInfo
	HasSQL      bool
	HasYAML     bool
}

func main() {
	if len(os.Args) < 2 {
		println("Usage: goenum <file.go>")
		os.Exit(1)
	}

	filename := os.Args[1]
	enums, err := parseFile(filename)
	if err != nil {
		println("Error parsing file:", err.Error())
		os.Exit(1)
	}

	if len(enums) > 0 {
		err := generateEnumsFile(enums, filename)
		if err != nil {
			println("Error generating enums file:", err.Error())
			os.Exit(1)
		}
	}
}

func parseFile(filename string) ([]EnumInfo, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var enums []EnumInfo
	enumsMap := make(map[string]*EnumInfo) // Map from type name to EnumInfo

	// First pass: find all goenums type declarations
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if genDecl.Doc != nil {
				for _, comment := range genDecl.Doc.List {
					if strings.Contains(comment.Text, "goenums:") {
						enum := parseEnumFromTypeSpec(typeSpec, comment.Text, node.Name.Name, filename)
						enums = append(enums, enum)
						enumsMap[enum.Type] = &enums[len(enums)-1]
						break
					}
				}
			}
		}
	}

	// Second pass: find const declarations and associate values
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		constBlockText := string(fileContent[fset.Position(genDecl.Pos()).Offset:fset.Position(genDecl.End()).Offset])

		for enumType, enumInfo := range enumsMap {
			values := parseConstValuesWithContext(constBlockText, genDecl, enumType)
			if len(values) > 0 {
				enumInfo.Values = append(enumInfo.Values, values...)
				if enumInfo.ConstBlock == "" {
					enumInfo.ConstBlock = constBlockText
				}
			}
		}
	}

	// Third pass: post-process (e.g., collect tags)
	for i := range enums {
		enum := &enums[i]
		tagSet := make(map[string]bool)
		for _, value := range enum.Values {
			for _, tag := range value.Tags {
				tagSet[tag] = true
			}
		}
		enum.AllTags = nil
		for tag := range tagSet {
			enum.AllTags = append(enum.AllTags, tag)
		}
	}

	return enums, nil
}

func FirstUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func parseEnumFromTypeSpec(ts *ast.TypeSpec, comment string, packageName string, filename string) EnumInfo {
	enum := EnumInfo{
		Name:        FirstUpper(ts.Name.Name),
		Type:        ts.Name.Name,
		PackageName: packageName,
		FileName:    filename,
		Options:     parseOptions(comment),
	}

	// Determine base type
	if ident, ok := ts.Type.(*ast.Ident); ok {
		enum.BaseType = ident.Name
	}

	// Generate container name using pluralization
	enum.ContainerName = FirstUpper(GetContainerName(enum.Name))

	return enum
}

func parseOptions(comment string) EnumOptions {
	options := EnumOptions{}

	// Parse goenums: comment
	re := regexp.MustCompile(`goenums:\s*(.+)`)
	matches := re.FindStringSubmatch(comment)
	if len(matches) < 2 {
		return options
	}

	optionsStr := matches[1]
	parts := strings.Fields(optionsStr)

	for _, part := range parts {
		switch {
		case part == "-sql":
			options.SQL = true
		case part == "-json":
			options.JSON = true
		case part == "-yaml":
			options.YAML = true
		case part == "-text":
			options.Text = true
		case part == "-binary":
			options.Binary = true
		case part == "-serde/name":
			options.SerdeFormat = "name"
		case part == "-serde/value":
			options.SerdeFormat = "value"
		case part == "-genName":
			options.GenName = true
		case part == "-statemachine":
			options.StateMachine = true
		}
	}

	return options
}

func parseConstValuesWithContext(constBlock string, decl *ast.GenDecl, enumType string) []EnumValue {
	// First parse using AST to get the structured data
	astValues := parseConstValues(decl, enumType)

	// Then parse the raw text to associate comments
	lines := strings.Split(constBlock, "\n")
	var values []EnumValue
	var currentPrecedingLines []string

	// Create a map of enum names from AST parsing
	astMap := make(map[string]EnumValue)
	for _, v := range astValues {
		astMap[v.Name] = v
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip const ( and )
		if trimmed == "const (" || trimmed == ")" {
			continue
		}

		// Check if this line contains an enum definition
		isEnumDefinition := false
		var enumName string

		if strings.Contains(line, "=") && !strings.HasPrefix(trimmed, "//") {
			// Explicit value definition
			parts := strings.Fields(trimmed)
			if len(parts) >= 1 {
				enumName = parts[0]
				isEnumDefinition = true
			}
		} else if !strings.HasPrefix(trimmed, "//") && trimmed != "" &&
			!strings.Contains(trimmed, "////") && !strings.Contains(trimmed, "//////") {
			// This might be an enum without explicit value (continuing iota)
			parts := strings.Fields(trimmed)
			if len(parts) >= 1 {
				candidateName := parts[0]
				// Check if this name exists in our AST map
				if _, exists := astMap[candidateName]; exists {
					enumName = candidateName
					isEnumDefinition = true
				}
			}
		}

		if isEnumDefinition && enumName != "" {
			if astValue, exists := astMap[enumName]; exists {
				// Copy AST data and add preceding lines
				value := astValue
				value.PrecedingLines = make([]string, len(currentPrecedingLines))
				copy(value.PrecedingLines, currentPrecedingLines)
				values = append(values, value)
				currentPrecedingLines = nil // Reset for next enum
			}
		} else {
			// This is a comment, separator, or empty line
			currentPrecedingLines = append(currentPrecedingLines, line)
		}
	}

	return values
}

func parseConstValues(decl *ast.GenDecl, enumType string) []EnumValue {
	var values []EnumValue
	iotaValue := -1
	var lastValueSpec *ast.ValueSpec

	for _, spec := range decl.Specs {
		if vs, ok := spec.(*ast.ValueSpec); ok {
			// If type is not specified, inherit from the previous spec
			if vs.Type == nil && lastValueSpec != nil {
				vs.Type = lastValueSpec.Type
			}

			// Check if this spec belongs to our enum type
			if typeIdent, ok := vs.Type.(*ast.Ident); !ok || typeIdent.Name != enumType {
				continue
			}

			// If values are not specified, inherit from the previous spec
			if len(vs.Values) == 0 && lastValueSpec != nil {
				vs.Values = lastValueSpec.Values
			}

			// Increment iota for each new spec
			iotaValue++

			for i, name := range vs.Names {
				value := EnumValue{
					Name: name.Name,
				}

				// Parse value - handle iota and explicit values
				if vs.Values != nil && i < len(vs.Values) {
					valExpr := vs.Values[i]
					if bl, ok := valExpr.(*ast.BasicLit); ok {
						value.Value = bl.Value
						if bl.Kind == token.INT {
							if _, err := strconv.Atoi(bl.Value); err == nil {
								// This is an explicit integer value, but iota continues based on position
							}
						}
					} else if ident, ok := valExpr.(*ast.Ident); ok && ident.Name == "iota" {
						value.Value = strconv.Itoa(iotaValue)
					} else {
						// This handles cases where the value is an expression,
						// or where iota is used implicitly.
						value.Value = strconv.Itoa(iotaValue)
					}
				} else {
					// No explicit value, use iotaValue
					value.Value = strconv.Itoa(iotaValue)
				}

				// Parse comment
				var commentText string
				if vs.Comment != nil {
					commentText = vs.Comment.Text()
				}
				if vs.Doc != nil {
					if commentText != "" {
						commentText = vs.Doc.Text() + "\n" + commentText
					} else {
						commentText = vs.Doc.Text()
					}
				}

				if commentText != "" {
					value.Comment = commentText
					value.OriginalComment = commentText
					parseValueComment(&value, commentText)
				}

				values = append(values, value)
			}
			lastValueSpec = vs
		}
	}

	return values
}

func parseValueComment(value *EnumValue, comment string) {
	lines := strings.Split(comment, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimPrefix(line, "//"))
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if line == "invalid" {
			value.IsInvalid = true
			continue
		}

		if strings.HasPrefix(line, "state:") {
			stateInfo := strings.TrimSpace(strings.TrimPrefix(line, "state:"))
			if stateInfo == "[final]" {
				value.IsFinal = true
			} else if strings.HasPrefix(stateInfo, "->") {
				transitionsStr := strings.TrimSpace(strings.TrimPrefix(stateInfo, "->"))
				if transitionsStr != "" {
					value.Transitions = strings.Split(transitionsStr, ",")
					for i, t := range value.Transitions {
						value.Transitions[i] = strings.TrimSpace(t)
					}
				}
			}
			continue
		}

		if strings.HasPrefix(line, "tag:") {
			tagsStr := strings.TrimSpace(strings.TrimPrefix(line, "tag:"))
			if tagsStr != "" {
				value.Tags = strings.Split(tagsStr, ",")
				for i, t := range value.Tags {
					value.Tags[i] = strings.TrimSpace(t)
				}
			}
			continue
		}

		// Parse multiple names (first line without : and not starting with //)
		if !strings.Contains(line, ":") && !strings.HasPrefix(line, "//") && strings.TrimSpace(line) != "" && len(value.Names) == 0 {
			names := strings.Split(line, ",")
			for _, name := range names {
				name = strings.TrimSpace(name)
				if name != "" && name != "invalid" {
					value.Names = append(value.Names, name)
				}
			}
		}
	}

	// If no names specified, use the constant name
	if len(value.Names) == 0 {
		value.Names = []string{value.Name}
	}
}

func convertStateNames(line string, nameMap map[string]string) string {
	// Handle state: -> target1, target2 format
	if strings.Contains(line, "state: ->") {
		parts := strings.Split(line, "->")
		if len(parts) == 2 {
			prefix := parts[0] // "// state: "
			targets := strings.TrimSpace(parts[1])

			// Split targets by comma
			targetList := strings.Split(targets, ",")
			var convertedTargets []string

			for _, target := range targetList {
				target = strings.TrimSpace(target)
				if convertedName, exists := nameMap[target]; exists {
					convertedTargets = append(convertedTargets, convertedName)
				} else {
					convertedTargets = append(convertedTargets, target)
				}
			}

			return prefix + "-> " + strings.Join(convertedTargets, ", ")
		}
	}

	return line
}

func generateEnumsFile(enums []EnumInfo, sourceFilename string) error {
	outputFile := strings.TrimSuffix(sourceFilename, ".go") + "_enums.go"

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	// Determine required imports
	hasSQL := false
	hasYAML := false
	for _, e := range enums {
		if e.Options.SQL {
			hasSQL = true
		}
		if e.Options.YAML {
			hasYAML = true
		}
	}

	data := FileTemplateData{
		PackageName: enums[0].PackageName,
		Enums:       enums,
		HasSQL:      hasSQL,
		HasYAML:     hasYAML,
	}

	// Execute template
	tmpl := template.Must(template.New("enumsFile").Funcs(template.FuncMap{
		"FirstUpper": FirstUpper,
		"ToLower":    strings.ToLower,
		"FormatPrecedingLines": func(lines []string, enumValues []EnumValue, currentValue EnumValue) string {
			if len(lines) == 0 {
				return ""
			}

			// Create a map of original enum names to handle state name conversions
			nameMap := make(map[string]string)
			for _, v := range enumValues {
				nameMap[v.Name] = v.Name
			}

			var result []string
			isFirstCommentLine := true

			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" {
					result = append(result, "")
				} else {
					processedLine := trimmed
					if strings.Contains(trimmed, "state:") {
						// No longer need to convert state names here, it's handled in the template
					}

					if isFirstCommentLine && strings.HasPrefix(trimmed, "//") && !strings.Contains(trimmed, "////") {
						processedLine = processedLine + " (" + currentValue.Value + ")"
						isFirstCommentLine = false
					}

					result = append(result, "\t"+processedLine)
				}
			}

			return strings.Join(result, "\n")
		},
	}).Parse(fileTemplate))

	return tmpl.Execute(w, data)
}

const fileTemplate = `// Code generated by goenum. DO NOT EDIT.

package {{.PackageName}}

import (
	"database/sql/driver"
	"fmt"
	"github.com/donutnomad/goenum/enums"
	"iter"
	{{- if .HasYAML}}
	"gopkg.in/yaml.v3"
	{{- end}}
)
{{range $enum := .Enums}}
// =================================================================================================
// {{.Name}}
// =================================================================================================

// {{.Name}} is a type that represents a single enum value.
// It combines the core information about the enum constant and its defined fields.
type {{.Name}} struct {
	{{.Type}}
}

// Verify that {{.Name}} implements the Enum interface
var _ enums.Enum[{{.BaseType}}, {{.Name}}] = {{.Name}}{}

// {{.Type}}Container is the container for all enum values.
// It is private and should not be used directly use the public methods on the {{.Name}} type.
type {{.Type}}Container struct {
	{{- range .Values}}
	{{- if .PrecedingLines}}
{{FormatPrecedingLines .PrecedingLines $enum.Values .}}
	{{- end}}
	{{FirstUpper .Name}} {{$enum.Name}}
	{{- end}}
}

// {{.ContainerName}} is a main entry point using the {{.Name}} type.
// It is a container for all enum values and provides a convenient way to access all enum values and perform
// operations, with convenience methods for common use cases.
var {{.ContainerName}} = {{.Type}}Container{
	{{- range .Values}}
	{{FirstUpper .Name}}: {{$enum.Name}}{ {{.Name}} },
	{{- end}}
}

// {{ToLower .Name}}NamesMap maps enum values to their names array
var {{ToLower .Name}}NamesMap = map[{{.Name}}][]string{
	{{- range .Values}}
	{{$enum.ContainerName}}.{{FirstUpper .Name}}: {
		{{- range .Names}}
		"{{.}}",
		{{- end}}
	},
	{{- end}}
}

{{- if .AllTags}}
// {{ToLower .Name}}TagsMap maps enum values to their tags array
var {{ToLower .Name}}TagsMap = map[{{.Name}}][]string{
	{{- range .Values}}
	{{- if .Tags}}
	{{$enum.ContainerName}}.{{FirstUpper .Name}}: {
		{{- range .Tags}}
		"{{.}}",
		{{- end}}
	},
	{{- end}}
	{{- end}}
}
{{- end}}

// {{.Name}}Raw is a type alias for the underlying enum type {{.Type}}.
// It provides direct access to the raw enum values for cases where you need
// to work with the underlying type directly.
type {{.Name}}Raw = {{.Type}}

// allSlice returns a slice of all enum values.
func (t {{.Type}}Container) allSlice() []{{.Name}} {
	return []{{.Name}}{
		{{- range .Values}}
		{{$enum.ContainerName}}.{{FirstUpper .Name}},
		{{- end}}
	}
}

// Val implements the Enum interface.
func (t {{.Name}}) Val() {{.BaseType}} {
	return {{.BaseType}}(t.{{.Type}})
}

// All implements the Enum interface.
func (t {{.Name}}) All() iter.Seq[{{.Name}}] {
	return func(yield func({{.Name}}) bool) {
		for _, v := range {{.ContainerName}}.allSlice() {
			if !v.IsValid() {
				continue
			}
			if !yield(v) {
				return
			}
		}
	}
}

// IsValid implements the Enum interface.
func (t {{.Name}}) IsValid() bool {
	{{- range .Values}}
	{{- if .IsInvalid}}
	if t == {{$enum.ContainerName}}.{{FirstUpper .Name}} {
		return false
	}
	{{- end}}
	{{- end}}
	return true
}

// Name implements the Enum interface.
// Returns the first name of the enum value.
func (t {{.Name}}) Name() string {
	if names, ok := {{ToLower .Name}}NamesMap[t]; ok && len(names) > 0 {
		return names[0]
	}
	return ""
}

// NameWith returns the name at the specified index.
// If the index is out of bounds, returns the last name.
func (t {{.Name}}) NameWith(idx int) string {
	names, ok := {{ToLower .Name}}NamesMap[t]
	if !ok || len(names) == 0 {
		return ""
	}
	if idx < 0 || idx >= len(names) {
		return names[len(names)-1]
	}
	return names[idx]
}

// Names returns all names of the enum value.
func (t {{.Name}}) Names() []string {
	if names, ok := {{ToLower .Name}}NamesMap[t]; ok {
		return names
	}
	return []string{}
}

// String implements the Stringer interface.
func (t {{.Name}}) String() string {
	if names, ok := {{ToLower .Name}}NamesMap[t]; ok && len(names) > 0 {
		return names[0]
	}
	return fmt.Sprintf("{{.Type}}(%v)", t.{{.Type}})
}

// SerdeFormat implements the Enum interface.
func (t {{.Name}}) SerdeFormat() enums.Format {
	{{- if eq .Options.SerdeFormat "name"}}
	return enums.FormatName
	{{- else}}
	return enums.FormatValue
	{{- end}}
}

// FromName implements the Enum interface.
func (t {{.Name}}) FromName(name string) ({{.Name}}, bool) {
	for enumValue, names := range {{ToLower .Name}}NamesMap {
		for _, n := range names {
			if n == name {
				return enumValue, enumValue.IsValid()
			}
		}
	}
	var zero {{.Name}}
	return zero, false
}

// FromValue implements the Enum interface.
func (t {{.Name}}) FromValue(value {{.BaseType}}) ({{.Name}}, bool) {
	for v := range {{.ContainerName}}.All() {
		if v.Val() == value {
			return v, true
		}
	}
	var zero {{.Name}}
	return zero, false
}

{{- if .AllTags}}
{{- range .AllTags}}

// {{FirstUpper .}}Slice returns all enum values that have the "{{.}}" tag.
func (t {{$enum.Type}}Container) {{FirstUpper .}}Slice() []{{$enum.Name}} {
	var result []{{$enum.Name}}
	for _, v := range t.allSlice() {
		if v.Is{{FirstUpper .}}() {
			result = append(result, v)
		}
	}
	return result
}

// Is{{FirstUpper .}} returns true if this enum value has the "{{.}}" tag.
func (t {{$enum.Name}}) Is{{FirstUpper .}}() bool {
	if tags, ok := {{ToLower $enum.Name}}TagsMap[t]; ok {
		for _, tag := range tags {
			if tag == "{{.}}" {
				return true
			}
		}
	}
	return false
}
{{- end}}
{{- end}}

// All container methods for convenience
func (t {{.Type}}Container) All() iter.Seq[{{.Name}}] {
	return {{.Name}}{}.All()
}

func (t {{.Type}}Container) FromName(name string) ({{.Name}}, bool) {
	return {{.Name}}{}.FromName(name)
}

func (t {{.Type}}Container) FromValue(value {{.BaseType}}) ({{.Name}}, bool) {
	return {{.Name}}{}.FromValue(value)
}

{{- if .Options.SQL}}

// Scan implements the database/sql.Scanner interface for {{.Name}}.
func (t *{{.Name}}) Scan(value any) error {
	result, err := enums.SQLScan(*t, value)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}

// Value implements the database/sql/driver.Valuer interface for {{.Name}}.
func (t {{.Name}}) Value() (driver.Value, error) {
	return enums.SQLValue(t)
}
{{- end}}

{{- if .Options.JSON}}

// MarshalJSON implements the json.Marshaler interface for {{.Name}}.
func (t {{.Name}}) MarshalJSON() ([]byte, error) {
	return enums.MarshalJSON(t, t.{{.Type}})
}

// UnmarshalJSON implements the json.Unmarshaler interface for {{.Name}}.
func (t *{{.Name}}) UnmarshalJSON(data []byte) error {
	result, err := enums.UnmarshalJSON(*t, data)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}
{{- end}}

{{- if .Options.YAML}}

// MarshalYAML implements the yaml.Marshaler interface for {{.Name}}.
func (t {{.Name}}) MarshalYAML() (any, error) {
	return enums.MarshalYAML(t, t.{{.Type}})
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for {{.Name}}.
func (t *{{.Name}}) UnmarshalYAML(node *yaml.Node) error {
	result, err := enums.UnmarshalYAML(*t, node)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}
{{- end}}

{{- if .Options.Text}}

// MarshalText implements the encoding.TextMarshaler interface for {{.Name}}.
func (t {{.Name}}) MarshalText() ([]byte, error) {
	return enums.MarshalText(t, t.{{.Type}})
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for {{.Name}}.
func (t *{{.Name}}) UnmarshalText(data []byte) error {
	result, err := enums.UnmarshalText(*t, data)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}
{{- end}}

{{- if .Options.Binary}}

// MarshalBinary implements the encoding.BinaryMarshaler interface for {{.Name}}.
func (t {{.Name}}) MarshalBinary() ([]byte, error) {
	return enums.MarshalBinary(t, t.{{.Type}})
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for {{.Name}}.
func (t *{{.Name}}) UnmarshalBinary(data []byte) error {
	result, err := enums.UnmarshalBinary(*t, data)
	if err != nil {
		return err
	}
	*t = *result
	return nil
}
{{- end}}

{{- if .Options.StateMachine}}

// CanTransitionTo checks if the current state can transition to the target state.
func (t {{.Name}}) CanTransitionTo(target {{.Name}}) bool {
	transitions := t.ValidTransitions()
	for _, validTarget := range transitions {
		if validTarget == target {
			return true
		}
	}
	return false
}

// ValidTransitions returns all valid target states that this state can transition to.
func (t {{.Name}}) ValidTransitions() []{{.Name}} {
	{{- range .Values}}
	{{- if .Transitions}}
	if t == {{$enum.ContainerName}}.{{FirstUpper .Name}} {
		return []{{$enum.Name}}{
			{{- range .Transitions}}
			{{$enum.ContainerName}}.{{FirstUpper .}},
			{{- end}}
		}
	}
	{{- end}}
	{{- end}}
	return []{{.Name}}{}
}

// IsTerminalState returns true if this state is a terminal (final) state.
func (t {{.Name}}) IsTerminalState() bool {
	{{- range .Values}}
	{{- if .IsFinal}}
	if t == {{$enum.ContainerName}}.{{FirstUpper .Name}} {
		return true
	}
	{{- end}}
	{{- end}}
	return false
}

// TerminalStateSlice returns a slice of all terminal states.
func (t {{.Name}}) TerminalStateSlice() []{{.Name}} {
	return []{{.Name}}{
		{{- range .Values}}
		{{- if .IsFinal}}
		{{$enum.ContainerName}}.{{FirstUpper .Name}},
		{{- end}}
		{{- end}}
	}
}
{{- end}}
{{end}}
`
