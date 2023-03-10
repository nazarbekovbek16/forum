package service

func ValidMailAddress(address string) bool {
	for i, v := range address {
		if v == '@' {
			address = address[i:]
		}
	}

	for i, v := range address {
		if v == '.' && i != len(address)-1 {
			return true
		}
	}
	return false
}
