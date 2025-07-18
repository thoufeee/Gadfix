package utlis

import "golang.org/x/crypto/bcrypt"

// create hashpassword
func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// check hashpassword
func CheckHash(password, newpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(newpassword))
	return err == nil
}
