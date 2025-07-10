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
	// Directly pluralize the full type name instead of removing suffixes
	// But first preserve the original case pattern
	if typeName == "" {
		return ""
	}
	
	// For compound words like "SwapStatus", we need to handle the pluralization correctly
	// We need to find the last word and pluralize it while preserving the case
	
	// First, try to find the last "word" in a camelCase/PascalCase string
	lastWordStart := 0
	for i := len(typeName) - 1; i > 0; i-- {
		if typeName[i] >= 'A' && typeName[i] <= 'Z' {
			lastWordStart = i
			break
		}
	}
	
	if lastWordStart > 0 {
		// We found a compound word like "SwapStatus" 
		prefix := typeName[:lastWordStart]
		lastWord := typeName[lastWordStart:]
		pluralLastWord := PluralizeEnglish(strings.ToLower(lastWord))
		// Capitalize the first letter of the pluralized last word
		pluralLastWord = strings.ToUpper(pluralLastWord[:1]) + pluralLastWord[1:]
		return prefix + pluralLastWord
	}
	
	// For simple words, just pluralize and preserve case
	plural := PluralizeEnglish(strings.ToLower(typeName))
	
	// Preserve original case pattern
	if len(typeName) > 0 && typeName[0] >= 'A' && typeName[0] <= 'Z' {
		// PascalCase - make first letter uppercase
		plural = strings.ToUpper(plural[:1]) + plural[1:]
	}
	
	return plural
}
