package auth

import (
	"errors"
	"net/http"
	"strings"
)

// extract API Key from header of HTTP request
// let say we want the header to look like
// EX:
// Authorization: APIKey {value here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authorization info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malform auth header")
	}
	// THis make us has to specify ApiKey and space the key in the header
	if vals[0] != "ApiKey" {
		return "", errors.New("malform first part of header")
	}
	return vals[1], nil
}
