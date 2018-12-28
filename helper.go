package vagrant_go

func contains(array []string, needle string) bool {
	for _, element := range array {
		if element == needle {
			return true
		}
	}

	return false
}
