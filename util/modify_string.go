package util

import "strings"

func GenerateSlug(str string) string {
	slug := strings.ToLower(str)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

func GenerateShortName(str string) string {
	words := strings.Fields(str)
	var shortsName string
	for _, word := range words {
		shortsName += string(word[0])
	}

	return strings.ToUpper(shortsName)
}
