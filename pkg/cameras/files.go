package cameras

import (
	"errors"
)

// Gets CamID, Year, Month and Day from a string like CMIDyyyyMMdd
func GetInfoFromFilename(filename string) (camID, year, month, day string, err error) {
	if len(filename) < 12 {
		return "", "", "", "", errors.New("incorrect format: too short")
	}

	for _, b := range filename[4:12] {
		if b < '0' || b > '9' {
			return "", "", "", "", errors.New("incorrect format: cannot parse date")
		}
	}

	return filename[:4], filename[4:8], filename[8:10], filename[10:12], nil
}
