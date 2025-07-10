package main

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
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

	for _, enum := range enums {
		err := generateEnumWithTemplate(enum)
		if err != nil {
			println("Error generating enum", enum.Name+":", err.Error())
			os.Exit(1)
		}
	}
}

func parseFile(filename string) ([]EnumInfo, error) {
	// Read the raw file content first
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
	var currentEnum *EnumInfo

	for _, decl := range node.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok {
						if d.Doc != nil {
							for _, comment := range d.Doc.List {
								if strings.Contains(comment.Text, "goenums:") {
									enum := parseEnumFromTypeSpec(ts, comment.Text, node.Name.Name, filename)
									enums = append(enums, enum)
									currentEnum = &enums[len(enums)-1]
									break
								}
							}
						}
					}
				}
			} else if d.Tok == token.CONST && currentEnum != nil {
				// Extract raw const block text
				start := fset.Position(d.Pos()).Offset
				end := fset.Position(d.End()).Offset
				currentEnum.ConstBlock = string(fileContent[start:end])
				
				values := parseConstValuesWithContext(currentEnum.ConstBlock, d, currentEnum.Name)
				currentEnum.Values = values

				// Collect all unique tags
				tagSet := make(map[string]bool)
				for _, value := range values {
					for _, tag := range value.Tags {
						tagSet[tag] = true
					}
				}
				for tag := range tagSet {
					currentEnum.AllTags = append(currentEnum.AllTags, tag)
				}
			}
		}
	}

	return enums, nil
}

func parseEnumFromTypeSpec(ts *ast.TypeSpec, comment string, packageName string, filename string) EnumInfo {
	enum := EnumInfo{
		Name:        strings.Title(ts.Name.Name),
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
	enum.ContainerName = GetContainerName(enum.Name)

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
		if strings.Contains(line, "=") && !strings.HasPrefix(trimmed, "//") {
			// Extract enum name
			parts := strings.Fields(trimmed)
			if len(parts) >= 1 {
				enumName := parts[0]
				if astValue, exists := astMap[enumName]; exists {
					// Copy AST data and add preceding lines
					value := astValue
					value.PrecedingLines = make([]string, len(currentPrecedingLines))
					copy(value.PrecedingLines, currentPrecedingLines)
					values = append(values, value)
					currentPrecedingLines = nil // Reset for next enum
				}
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

	for _, spec := range decl.Specs {
		if vs, ok := spec.(*ast.ValueSpec); ok {
			for i, name := range vs.Names {
				value := EnumValue{
					Name: name.Name,
				}

				// Parse value
				if vs.Values != nil && i < len(vs.Values) {
					if bl, ok := vs.Values[i].(*ast.BasicLit); ok {
						value.Value = bl.Value
					}
				}

				// Parse comment - check both inline and doc comments
				var commentText string
				if vs.Comment != nil {
					commentText = vs.Comment.Text()
				}
				// Also check for doc comments on the previous line
				if vs.Doc != nil {
					if commentText != "" {
						commentText = vs.Doc.Text() + "\n" + commentText
					} else {
						commentText = vs.Doc.Text()
					}
				}

				if commentText != "" {
					value.Comment = commentText
					value.OriginalComment = commentText // Save raw comment for display
					parseValueComment(&value, commentText)
				}

				values = append(values, value)
			}
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

func generateEnumWithTemplate(enum EnumInfo) error {
	outputFile := strings.TrimSuffix(enum.FileName, ".go") + "_enums.go"

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	// Execute template
	tmpl := template.Must(template.New("enum").Funcs(template.FuncMap{
		"Title":     strings.Title,
		"ToLower":   strings.ToLower,
		"ToUpper":   strings.ToUpper,
		"HasPrefix": strings.HasPrefix,
		"Join":      strings.Join,
		"FirstUpper": func(s string) string {
			if len(s) == 0 {
				return s
			}
			return strings.ToUpper(s[:1]) + s[1:]
		},
		"FormatPrecedingLines": func(lines []string, enumValues []EnumValue) string {
			if len(lines) == 0 {
				return ""
			}
			
			// Create a map of original enum names to title case names
			nameMap := make(map[string]string)
			for _, v := range enumValues {
				nameMap[v.Name] = strings.Title(v.Name)
			}
			
			var result []string
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" {
					result = append(result, "")
				} else {
					// Process state transitions to convert names
					processedLine := trimmed
					if strings.Contains(trimmed, "state:") {
						processedLine = convertStateNames(trimmed, nameMap)
					}
					result = append(result, "\t"+processedLine)
				}
			}
			
			return strings.Join(result, "\n")
		},
	}).Parse(enumTemplate))

	return tmpl.Execute(w, enum)
}

const enumTemplate = `package {{.PackageName}}

import (
	"database/sql/driver"
	"fmt"
	"github.com/donutnomad/goenum/enums"
	"iter"
	{{- if .Options.SQL}}
	{{- end}}
	{{- if .Options.YAML}}
	"gopkg.in/yaml.v3"
	{{- end}}
)

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
{{FormatPrecedingLines .PrecedingLines $.Values}}
	{{- end}}
	{{Title .Name}} {{$.Name}}
	{{- end}}
}

// {{Title .ContainerName}} is a main entry point using the {{.Name}} type.
// It is a container for all enum values and provides a convenient way to access all enum values and perform
// operations, with convenience methods for common use cases.
var {{Title .ContainerName}} = {{.Type}}Container{
	{{- range .Values}}
	{{Title .Name}}: {{$.Name}}{
		{{$.Type}}: {{.Name}},
	},
	{{- end}}
}

// {{ToLower .Name}}NamesMap maps enum values to their names array
var {{ToLower .Name}}NamesMap = map[{{.Name}}][]string{
	{{- range .Values}}
	{{Title $.ContainerName}}.{{Title .Name}}: {
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
	{{Title $.ContainerName}}.{{Title .Name}}: {
		{{- range .Tags}}
		"{{.}}",
		{{- end}}
	},
	{{- end}}
	{{- end}}
}
{{- end}}

// {{Title .ContainerName}}Raw is a type alias for the underlying enum type {{.Type}}.
// It provides direct access to the raw enum values for cases where you need
// to work with the underlying type directly.
type {{Title .ContainerName}}Raw = {{.Type}}

// allSlice returns a slice of all enum values.
func (t {{.Type}}Container) allSlice() []{{.Name}} {
	return []{{.Name}}{
		{{- range .Values}}
		{{Title $.ContainerName}}.{{Title .Name}},
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
		for _, v := range {{Title .ContainerName}}.allSlice() {
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
	if t == {{Title $.ContainerName}}.{{Title .Name}} {
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
	for v := range {{Title .ContainerName}}.All() {
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
func (t {{$.Type}}Container) {{FirstUpper .}}Slice() []{{$.Name}} {
	var result []{{$.Name}}
	for _, v := range t.allSlice() {
		if v.Is{{FirstUpper .}}() {
			result = append(result, v)
		}
	}
	return result
}

// Is{{FirstUpper .}} returns true if this enum value has the "{{.}}" tag.
func (t {{$.Name}}) Is{{FirstUpper .}}() bool {
	if tags, ok := {{ToLower $.Name}}TagsMap[t]; ok {
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
	if t == {{Title $.ContainerName}}.{{Title .Name}} {
		return []{{$.Name}}{
			{{- range .Transitions}}
			{{Title $.ContainerName}}.{{Title .}},
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
	if t == {{Title $.ContainerName}}.{{Title .Name}} {
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
		{{Title $.ContainerName}}.{{Title .Name}},
		{{- end}}
		{{- end}}
	}
}
{{- end}}
`
