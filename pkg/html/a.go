package html

import "regexp"

type A struct {
	HREF, Text string
}

var aRegex = regexp.MustCompile("<a href=\"([0-9A-Za-z_./]+)\">([0-9A-Za-z_./]+)</a>")

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
