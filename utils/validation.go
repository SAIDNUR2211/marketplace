package utils

import (
	"strings"
)

func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return true
}
