package auth

import (
	"testing"

	"github.com/bdjekel/chirpy/internal/auth"
)

func TestMakeJWT(t *testing.T) {
	auth.MakeJWT()
}

idCount := 15
userIDs := []string{}

for i:=0; i < idCount; i++ {
	newID := gen_random_uuid()
	userIDs = append(userIDs, newID)
}

secretToken := "ThIsIsAsEcReTtOkEn9876"
expiresIn := rand.Intn(3600) * time.Second // 
