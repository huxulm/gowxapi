package api

import (
	"net/http"

	resty "github.com/go-resty/resty/v2"
	m "github.com/jackdon/gowxapi/models"
)

// GetParam is a util func to get query param value by name
func GetParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// Code2Session 登录凭证校验。通过 wx.login 接口获得临时登录凭证 code
// 后传到开发者服务器调用此接口完成登录流程。更多使用方法详见
// 小程序登录: https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/login.html
func Code2Session(jsCode string) (*m.WxRespCode2Session, error) {
	var respData *m.WxRespCode2Session
	// Create a Resty Client
	client := resty.New()
	qp := make(map[string]string)
	qp["js_code"] = jsCode
	qp["appid"] = APPID
	qp["secret"] = SECRET
	qp["grant_type"] = "authorization_code"
	resp, err := client.R().
		EnableTrace().
		SetQueryParams(qp).
		Get(Code2SessionURL)
	if err != nil {
		return respData, err
	}
	respData = new(m.WxRespCode2Session)
	client.JSONUnmarshal(resp.Body(), &respData)
	return respData, nil
}

// GetAccessToken 获取小程序全局唯一后台接口调用凭据（access_token）。
// 调用绝大多数后台接口时都需使用 access_token，开发者需要进行妥善保存。
func GetAccessToken() (*m.WxRespAccessToken, error) {
	var respData *m.WxRespAccessToken
	// Create a Resty Client
	client := resty.New()
	qp := make(map[string]string)
	qp["appid"] = APPID
	qp["secret"] = SECRET
	qp["grant_type"] = "client_credential"
	resp, err := client.R().
		EnableTrace().
		SetQueryParams(qp).
		Get(AccessTokenURL)
	if err != nil {
		return respData, err
	}
	respData = new(m.WxRespAccessToken)
	client.JSONUnmarshal(resp.Body(), &respData)
	return respData, nil
}
