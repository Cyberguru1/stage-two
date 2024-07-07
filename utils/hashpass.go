package utils

import "golang.org/x/crypto/bcrypt"


func HashPassword(password string) (string, error) {

	if b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost);  err != nil {
		return "", err
	}else {
		return string(b), nil
	}

	return "", nil
}

func ComparePassword(attempt, password string) error {
	bytePassword, byteHashedPassword := []byte(attempt), []byte(password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}