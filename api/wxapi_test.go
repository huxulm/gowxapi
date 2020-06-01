package api

import (
	"encoding/json"
	"testing"

	"github.com/jackdon/gowxapi/config"
)

func TestGetAccessToken(t *testing.T) {
	appID := config.C.AppInfo.AppID
	secret := config.C.AppInfo.Secret
	accessToken, _ := GetAccessToken(appID, secret)
	aB, _ := json.Marshal(accessToken)
	println(string(aB))
}
