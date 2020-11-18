package helper

// Contains checks if value exists in the list
func Contains(key string, utList []string) bool {
	for _, ut := range utList {
		if key == ut {
			return true
		}
	}
	return false
}
