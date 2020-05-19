package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackdon/gowxapi/models"
	"github.com/julienschmidt/httprouter"
)

// Login is used to accept params after client invoking wx.login
// and request auth.code2Session in server side
func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jsCode := GetParam(r, "js_code")
	wxResp, err := Code2Session(jsCode)
	if err != nil {
		defer RespError(&w)
	}
	if err != nil {
		defer RespError(&w)
	} else {
		var resp models.Resp
		code, msg, data := 0, "", map[string]string{}
		if *wxResp.ErrCode == 0 {
			code, msg = 0, "成功"
			data = map[string]string{
				"session_key": *wxResp.SessionKey, "open_id": *wxResp.OpenID,
			}
		} else {
			code, msg, data = 1, *wxResp.ErrMsg, nil
		}
		resp = models.Resp{&code, &msg, &data}
		respJSON, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, string(respJSON))
	}
}
