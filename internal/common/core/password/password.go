package password

import (
	"strings"
	"unicode/utf8"
)

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func Validate(password string, minLen int, maxLen int, incSym bool, incNum bool) (bool, bool, bool) {
	lenPass := false
	symPass := false
	numPass := false

	if utf8.RuneCountInString(password) >= minLen && utf8.RuneCountInString(password) <= maxLen {
		lenPass = true
	}

	if strings.ContainsAny(password, "!*-,.$\"&()+=") {
		symPass = true
	}

	if strings.ContainsAny(password, "0123456789") {
		numPass = true
	}

	return lenPass, symPass, numPass
}
