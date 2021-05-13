package app

import (
	"errors"
	"net/mail"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

//validUserID validate user id
//must contain only letters & numbers
func validUserID(userID string) bool {
	var alphaNum = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	if !alphaNum(userID) {
		return false
	}
	return true
}

//validPwd validates password
func validPwd(pwd string) bool {
	// Password rules :
	// at least : 1 upper case, 1 lower case, 1 number, 1 special char, no whitespace
	// min 8 chars

	var (
		upper, lower, number, special bool
		count                         int
	)

	for _, char := range pwd {
		switch {
		case unicode.IsUpper(char):
			upper = true
			count++
		case unicode.IsLower(char):
			lower = true
			count++
		case unicode.IsNumber(char):
			number = true
			count++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special = true
			count++
		default:
			return false
		}
	}

	if !upper || !lower || !number || !special || count < 8 {
		return false
	}

	_, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return false
	}

	return true
}

//validEmail validates email
func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}

//validateRegisterInput validates email, userID and passwords
func validateRegisterInput(email, userID, password1, password2 string) error {
	if email == "" || userID == "" || password1 == "" || password2 == "" {
		return errors.New("please enter all the fields")
	}

	if !validEmail(email) {
		return errors.New("please enter a valid email address")
	}

	if !validUserID(userID) {
		return errors.New("user id contain invalid characters")
	}

	if !validPwd(password1) || !validPwd(password2) || (password1 != password2) {
		return errors.New("please enter a valid password")
	}

	return nil
}
