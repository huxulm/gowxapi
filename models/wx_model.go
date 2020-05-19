package models

// WxResp is ...
type WxResp struct {
	ErrCode *int    `json:"errcode"` // 错误码
	ErrMsg  *string `json:"errmsg"`  // 错误信息
}

// WxRespCode2Session is ...
type WxRespCode2Session struct {
	WxResp
	OpenID     *string `json:"openid"`      // 用户唯一标识
	SessionKey *string `json:"session_key"` // 会话密钥
	UnionID    *string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回，详见 UnionID 机制说明。
}

// WxRespAccessToken is ...
type WxRespAccessToken struct {
	WxResp
	AccessToken string `json:"access_token"` //	获取到的凭证
	ExpiresIn   int    `json:"expires_in"`   //	凭证有效时间，单位：秒。目前是7200秒之内的值。
}
