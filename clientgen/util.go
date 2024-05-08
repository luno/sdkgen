package clientgen

import "regexp"

// SanitiseVariableName removes special characters from strings intended to be used as variables
func SanitiseVariableName(input string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9]+")
	return regex.ReplaceAllString(input, "")
}
