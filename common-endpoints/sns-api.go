package ce

import (
	"github.com/rosbit/go-wx-api/v2/oauth2"
	"fmt"
	"net/http"
)

// GET ${commonEndpoints.SnsAPI}?s=<service-name-in-conf>&code=<code-from-wx-server>&scope={userinfo|base}
func SnsAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("SnsAPI called\n")
	service := r.FormValue("s")
	if service == "" {
		writeError(w, http.StatusBadRequest, "s(ervice) parameter expected")
		return
	}

	scope := r.FormValue("scope")
	switch scope {
	case "userinfo","base":
	case "", "snsapi_base":
		scope = "base"
	case "snsapi_userinfo":
		scope = "userinfo"
	default:
		writeError(w, http.StatusBadRequest, `scope must be "useinfo", "base", "sns_userinfo" or "sns_base"`)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "code parameter expected")
		return
	}

	wxUser := wxoauth2.NewWxUser(service)
	openId, err := wxUser.GetOpenId(code)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var userInfo map[string]interface{}
	if scope == "base" {
		userInfo, err = wxUser.GetInfoByAccessToken()
	} else {
		if err = wxUser.GetInfo(); err == nil {
			userInfo = wxUser.UserInfo
		}
	}

	writeJson(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"openId": openId,
		"userInfo": userInfo,
		"error": func()string{if err == nil {return ""}; return err.Error()}(),
	})
}

