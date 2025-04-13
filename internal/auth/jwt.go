package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	})
	
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	
	claims := jwt.RegisteredClaims{}
//TODO: Read more in-depth on keyfunc argument below. Had to copy pasta from boot.dev's solution file.
	keyFunc := func(token *jwt.Token) (interface{}, error) { 
		return []byte(tokenSecret), nil }

	token, err := jwt.ParseWithClaims(
		tokenString, 
		&claims, 
		keyFunc)
	if err != nil {
		log.Println("TOKEN NOT PARSED CORRECTLY")
		log.Print(token)
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

