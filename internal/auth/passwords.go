package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing failed.")
		return "", err
	}
	
	return string(hashed_password), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}