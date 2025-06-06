package auth

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
	// fmt.Println("-------------INPUT---------------")
	// fmt.Printf("userID ==> %s\n", userID)
	// fmt.Printf("secretToken ==> %s\n", secretToken)
	// fmt.Printf("expiresIn ==> %s\n", expiresIn)
	// fmt.Println("---------------------------------")
	
	//make jwt
	jwt, err := MakeJWT(userID, secretToken, expiresIn)
	if err != nil {
		t.Errorf("Error creating jwt: %s\n", err)
		return
	}



	// let jwt expire
	time.Sleep(time.Second)

	// attempt to validate
	expectedErrMsg := "token has invalid claims: token is expired"
	validatedUserID, err := ValidateJWT(jwt, secretToken)
	if err != nil {
		errMsg := fmt.Sprint(err)
		// fmt.Println("------------------")
		// fmt.Printf(">>%s<<", errMsg)
		// fmt.Println()
		// fmt.Printf(">>%s<<", expectedErrMsg)
		// fmt.Println("------------------")
		if errMsg == expectedErrMsg {
			fmt.Println("expired token test passed")
			return
		}
		t.Error(err)
	}
	fmt.Printf("FAIL. Test did not invalidate JWT. See valid token:\n%s", validatedUserID)
}

func TestValidateWrongSecretJWT(t *testing.T) {
	// parameters for jwt
	userID := uuid.New()
	secretToken := "ThIsIsAsEcReTtOkEn9876"		
	expiresIn := time.Duration(rand.Intn(750)) * time.Millisecond 
	// fmt.Println("-------------INPUT---------------")
	// fmt.Printf("userID ==> %s\n", userID)
	// fmt.Printf("secretToken ==> %s\n", secretToken)
	// fmt.Printf("expiresIn ==> %s\n", expiresIn)
	// fmt.Println("---------------------------------")
	
	//make jwt
	jwt, err := MakeJWT(userID, secretToken, expiresIn)
	if err != nil {
		t.Errorf("Error creating jwt: %s\n", err)
		return
	}

	// attempt to validate
	wrongSecretToken := "incorrectamundo"
	expectedErrMsg := "token signature is invalid: signature is invalid"
	_, err = ValidateJWT(jwt, wrongSecretToken)
	if err != nil {
		errMsg := fmt.Sprint(err)
		// fmt.Println("------------------")
		// fmt.Printf(">>%s<<", errMsg)
		// fmt.Println()
		// fmt.Printf(">>%s<<", expectedErrMsg)
		// fmt.Println("------------------")
		if errMsg == expectedErrMsg {
			fmt.Println("incorrect secret token test passed")
			return
		}
		t.Error(err)
	}
	fmt.Printf("FAIL. Test did not invalidate JWT. See valid secret token:\n%s", secretToken)
}

func TestValidateCorruptedJWT(t *testing.T) {
	// parameters for jwt
	userID := uuid.New()
	secretToken := "ThIsIsAsEcReTtOkEn9876"		
	expiresIn := time.Duration(rand.Intn(7500)) * time.Millisecond + (10 * time.Second)
	// fmt.Println("-------------INPUT---------------")
	// fmt.Printf("userID ==> %s\n", userID)
	// fmt.Printf("secretToken ==> %s\n", secretToken)
	// fmt.Printf("expiresIn ==> %s\n", expiresIn)
	// fmt.Println("---------------------------------")
	
	//make jwt
	jwt, err := MakeJWT(userID, secretToken, expiresIn)
	if err != nil {
		t.Errorf("Error creating jwt: %s\n", err)
		return
	}

	// make wrong jwt
	wrongUserID := uuid.New()
	wrongSecretToken := "incorrectamundo"
	wrongJwt, err := MakeJWT(wrongUserID, wrongSecretToken, expiresIn)
	if err != nil {
		t.Errorf("Error creating jwt: %s\n", err)
		return
	}

	// attempt to validate
	expectedErrMsg := "token signature is invalid: signature is invalid"
	_, err = ValidateJWT(wrongJwt, secretToken)
	if err != nil {
		errMsg := fmt.Sprint(err)
		fmt.Println("------------------")
		fmt.Printf(">>%s<<", errMsg)
		fmt.Println()
		fmt.Printf(">>%s<<", expectedErrMsg)
		fmt.Println("------------------")
		if errMsg == expectedErrMsg {
			fmt.Println("incorrect jwt test passed")
			return
		}
		t.Error(err)
	}
	fmt.Printf("FAIL. Test did not invalidate JWT. See valid jwt:\n%s", jwt)
	fmt.Printf("\nand invalid jwt:\n%s", wrongJwt)
}


func TestGetBearerTokenValid(t *testing.T) {
	// define test headers
	headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
	}

	//call GetBearerToken on said headers
	bearerToken, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("Bearer token text failed: %s", err)
	}
	fmt.Println("test passed. See bearerToken below:")
	fmt.Println(bearerToken)
}

func TestGetBearerTokenNoBearerPrefix(t *testing.T) {
	// define test headers
	headers := http.Header{
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
	}

	//call GetBearerToken on said headers
	bearerToken, err := GetBearerToken(headers)
	if err != nil {
		errString := fmt.Sprint(err)
		if errString == "authorization header missing Bearer prefix" {
			fmt.Println("missing bearer prefix test passed.")
			return
		}
		t.Error("incorrect error message")
	}
	t.Errorf("error not thrown. should be missing bearer prefix. See bearerToken below:\n>>%s<<\n", bearerToken)
}