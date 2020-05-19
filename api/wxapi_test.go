package api

import (
	"encoding/json"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	accessToken, _ := GetAccessToken()
	aB, _ := json.Marshal(accessToken)
	println(string(aB))
}
