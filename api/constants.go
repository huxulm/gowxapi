package api

import "os"

var (
	// APPID is wx appid
	APPID = os.Getenv("appid")
	// SECRET is wx secret
	SECRET = os.Getenv("secret")
)

// BaseURL is api.weixin.qq.com base url which any of other apis are based on
const BaseURL = "https://api.weixin.qq.com/"

// Code2SessionURL is used with wx.login in backend need query params: `appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code`
const Code2SessionURL = BaseURL + "sns/jscode2session"

// AccessTokenURL is ...
const AccessTokenURL = BaseURL + "cgi-bin/token"
