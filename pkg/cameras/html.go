package cameras

import "regexp"

// videoNameRegex is the matching regex for the videos generated by the
// IP Camera OneWay OWIPCAM45. This filename structure is "PYYMMDD_hhmmss_hhmmss.ext", where:
// "YY" are the two digits of the year
// "MM" is the month
// "DD" is the day of the month
// The first "hhmmss" corresponds at the time when the video started
// The second "hhmmss" corresponds at the time when the video stopped
// "ext" is the extension of the file
var videoNameRegex = regexp.MustCompile("^P[0-9]{6}_[0-9]{6}_[0-9]{6}\\..+$")

// GetInfoFromFilenameOWIPCAM45 gets the info contained in the filename of the videos recorded by a OWIPCAM45.
// See videoNameRegex to know the structure of the name.
func GetInfoFromFilenameOWIPCAM45(s string) (year, month, day, rest string) {
	return s[1:3], s[3:5], s[5:7], s[8:]
}

// IsValidVideoName checks if the string provided matches with videoNameRegex.
func IsValidVideoName(s string) bool {
	if len(s) < 21 {
		return false
	}
	return videoNameRegex.MatchString(s) && s[15:21] != "999999"
}

// IsValidFolderName checks if the string provided is a valid name of a folder generated by the
// IP Camera OneWay OWIPCAM45. This folder name is like "YYYYMMDD/" where:
// "YYYY" is the year
// "MM" is the month
// "DD" is the day of the month
func IsValidFolderName(s string) bool {
	if len(s) != 9 {
		return false
	}

	return isMadeUpOfNumbers(s[:len(s)-1])
}

// isMadeUpOfNumbers checks if the string provided is made up of numeral digits
func isMadeUpOfNumbers(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
