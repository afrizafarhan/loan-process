package common

import (
	"strings"
)

var Tenor = []uint{3, 6, 9, 12, 24}

var StatusLoanRequest = []string{"accepted", "rejected"}

func InArray(array []uint, search interface{}) bool {
	for _, val := range array {
		if val == search {
			return true
		}
	}
	return false
}

func InAllowedImageExtension(ext string) bool {
	var allowed = []string{"png", "jpg", "jpeg"}
	for _, val := range allowed {
		if val == ext {
			return true
		}
	}
	return false
}

func GetImageExtension(filename string) string {
	nameExt := strings.Split(filename, ".")
	return nameExt[len(nameExt)-1]
}
