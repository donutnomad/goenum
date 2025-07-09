package main

import (
	"strings"
)

// PluralizeEnglish converts English words to their plural forms
func PluralizeEnglish(word string) string {
	word = strings.ToLower(word)

	// Special irregular plurals
	irregulars := map[string]string{
		"status":     "statuses",
		"child":      "children",
		"person":     "people",
		"man":        "men",
		"woman":      "women",
		"tooth":      "teeth",
		"foot":       "feet",
		"mouse":      "mice",
		"goose":      "geese",
		"ox":         "oxen",
		"datum":      "data",
		"medium":     "media",
		"crisis":     "crises",
		"analysis":   "analyses",
		"basis":      "bases",
		"oasis":      "oases",
		"thesis":     "theses",
		"axis":       "axes",
		"matrix":     "matrices",
		"index":      "indices",
		"vertex":     "vertices",
		"appendix":   "appendices",
		"criterion":  "criteria",
		"phenomenon": "phenomena",
		"focus":      "foci",
		"radius":     "radii",
		"nucleus":    "nuclei",
		"cactus":     "cacti",
		"fungus":     "fungi",
		"stimulus":   "stimuli",
		"alumnus":    "alumni",
	}

	if plural, exists := irregulars[word]; exists {
		return plural
	}

	// Words ending in 's', 'ss', 'sh', 'ch', 'x', 'z' -> add 'es'
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "ss") ||
		strings.HasSuffix(word, "sh") || strings.HasSuffix(word, "ch") ||
		strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z") {
		return word + "es"
	}

	// Words ending in consonant + 'y' -> change 'y' to 'ies'
	if strings.HasSuffix(word, "y") && len(word) > 1 {
		secondLast := word[len(word)-2]
		if !isVowel(secondLast) {
			return word[:len(word)-1] + "ies"
		}
	}

	// Words ending in 'f' or 'fe' -> change to 'ves'
	if strings.HasSuffix(word, "f") {
		return word[:len(word)-1] + "ves"
	}
	if strings.HasSuffix(word, "fe") {
		return word[:len(word)-2] + "ves"
	}

	// Words ending in consonant + 'o' -> add 'es'
	if strings.HasSuffix(word, "o") && len(word) > 1 {
		secondLast := word[len(word)-2]
		if !isVowel(secondLast) {
			return word + "es"
		}
	}

	// Default: just add 's'
	return word + "s"
}

func isVowel(c byte) bool {
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u'
}

// GetContainerName generates the container name for an enum type
func GetContainerName(typeName string) string {
	// Remove common enum suffixes and get the base name
	name := typeName
	suffixes := []string{"Status", "Type", "Kind", "State", "Mode", "Level", "Priority", "Category"}

	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) {
			baseLen := len(name) - len(suffix)
			if baseLen > 0 {
				name = name[:baseLen]
				break
			}
		}
	}

	// If no suffix was removed or name is empty, use the original name
	if name == "" || name == typeName {
		name = typeName
	}

	return PluralizeEnglish(name)
}
