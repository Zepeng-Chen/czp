package internal

import "golang.org/x/crypto/bcrypt"

func HashPasswd(passwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 12) // the cost should be at least 10
	return string(bytes), err
}

func CheckPasswdHash(passwd, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd))
	return err == nil
}
