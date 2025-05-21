package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	prefix := "ApiKey "
	authHeader := headers.Get("Authorization")
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header missing ApiKey prefix")
	}
	apiKeyToken := strings.TrimPrefix(authHeader, prefix)
	return apiKeyToken, nil
	
}