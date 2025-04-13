package auth

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secretToken := "ThIsIsAsEcReTtOkEn9876"		
	expiresIn := time.Duration(rand.Intn(3600)) * time.Second // testing 0s to 1 hr		
	jwt, err := MakeJWT(userID, secretToken, expiresIn)
	if err != nil {
		t.Errorf("Error creating jwt: %s\n", err)
		return
	}
	log.Printf("JWT ==> %s\n", jwt)
}

func TestValidateValidJWT(t *testing.T) {
	// parameters for jwt
	userID := uuid.New()
	secretToken := "ThIsIsAsEcReTtOkEn9876"		
	expiresIn := time.Duration(rand.Intn(750)) * time.Second 
	fmt.Println("-------------INPUT---------------")
	fmt.Printf("userID ==> %s\n", userID)
	fmt.Printf("secretToken ==> %s\n", secretToken)
	fmt.Printf("expiresIn ==> %s\n", expiresIn)
	fmt.Println("---------------------------------")
	
	//make jwt
	jwt, err := MakeJWT(userID, secretToken, expiresIn)
	if err != nil {
		t.Errorf("Error creating jwt: %s\n", err)
		return
	}

	fmt.Println("---------------------------------")
	fmt.Printf("JWT ==> %s\n", jwt)
	fmt.Println("---------------------------------")

	// attempt to validate
	validatedUserID, err := ValidateJWT(jwt, secretToken)
	fmt.Println("---------------------------------")
	fmt.Printf("validatedUserID ==> %s\n", validatedUserID)
	fmt.Println("---------------------------------")

	if validatedUserID != userID || err != nil {
		t.Errorf("testing testing: %s\n", err)
	}
}

func TestValidateExpiredJWT(t *testing.T) {
	// parameters for jwt
	userID := uuid.New()
	secretToken := "ThIsIsAsEcReTtOkEn9876"		
	expiresIn := time.Duration(rand.Intn(750)) * time.Millisecond 
	fmt.Println("-------------INPUT---------------")
	fmt.Printf("userID ==> %s\n", userID)
	fmt.Printf("secretToken ==> %s\n", secretToken)
	fmt.Printf("expiresIn ==> %s\n", expiresIn)
	fmt.Println("---------------------------------")
	
	//make jwt
	jwt, err := MakeJWT(userID, secretToken, expiresIn)
	if err != nil {
		t.Errorf("Error creating jwt: %s\n", err)
		return
	}

	fmt.Println("---------------------------------")
	fmt.Printf("JWT ==> %s\n", jwt)
	fmt.Println("---------------------------------")


	// let jwt expire
	time.Sleep(time.Second)

	// attempt to validate
	validatedUserID, err := ValidateJWT(jwt, secretToken)
	fmt.Println("---------------------------------")
	fmt.Printf("validatedUserID ==> %s\n", validatedUserID)
	fmt.Println("---------------------------------")

	if validatedUserID != userID || err != nil {
		t.Errorf("testing testing: %s\n", err)
	}
}