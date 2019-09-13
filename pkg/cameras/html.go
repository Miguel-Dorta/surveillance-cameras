package cameras

func ValidFolderName(s string) bool {
	if len(s) == 0 {
		return false
	}

	s = s[:len(s)-1]
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
