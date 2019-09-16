package html

import "regexp"

// A represents a simple <a> tag in HTML containing an HREF and the text itself
type A struct {
	HREF, Text string
}

// aRegex is the regex for parsing simple HTML <a> tags
var aRegex = regexp.MustCompile("<a href=\"([0-9A-Za-z_./]+)\">([0-9A-Za-z_./]+)</a>")

// GetAList gets a list of all the <a> tags in the data provided
func GetAList(data []byte) []A {
	matches := aRegex.FindAllSubmatch(data, -1)
	result := make([]A, 0, len(matches))
	for _, match := range matches {
		if match[1] == nil || match[2] == nil {
			continue
		}
		result = append(result, A{
			HREF: string(match[1]),
			Text: string(match[2]),
		})
	}
	return result
}
