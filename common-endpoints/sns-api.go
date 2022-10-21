package ce

import (
	"github.com/rosbit/go-wx-api/v2/oauth2"
	"github.com/rosbit/go-wx-api/v2/auth"
	"github.com/rosbit/mgin"
	"net/http"
)

// GET ${commonEndpoints.SnsAPI}?s=<service-name-in-conf>&code=<code-from-wx-server>&scope={userinfo|base}
func SnsAPI(c *mgin.Context) {
	var params struct {
		Service string `query:"s"`
		Code    string `query:"code"`
		Scope   string `query:"scope" optional:"true"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	switch params.Scope {
	case "userinfo","base":
	case "", "snsapi_base":
		params.Scope = "base"
	case "snsapi_userinfo":
		params.Scope = "userinfo"
	default:
		c.Error(http.StatusBadRequest, `scope must be "useinfo", "base", "sns_userinfo" or "sns_base"`)
		return
	}

	wxUser := wxoauth2.NewWxUser(params.Service)
	openId, err := wxUser.GetOpenId(params.Code)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	var userInfo *wxauth.WxUserInfo
	if params.Scope == "base" {
		userInfo, err = wxUser.GetInfoByAccessToken()
	} else {
		if err = wxUser.GetInfo(); err == nil {
			userInfo = wxUser.UserInfo
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"openId": openId,
		"userInfo": userInfo,
		"error": func()string{if err == nil {return ""}; return err.Error()}(),
	})
}

