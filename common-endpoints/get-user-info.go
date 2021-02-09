package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
	"net/http"
)

// GET ${commonEndpoints.WxUser}?s=<service-name-in-conf>&o=<openId>
func GetWxUserInfo(w http.ResponseWriter, r *http.Request) {
	service := r.FormValue("s")
	if service == "" {
		writeError(w, http.StatusBadRequest, "s(ervice) parameter expected")
		return
	}

	openId := r.FormValue("o")
	if openId == "" {
		writeError(w, http.StatusBadRequest, "o(penId) parameter expected")
		return
	}

	userInfo, err := wxtools.GetUserInfo(service, openId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"userInfo": userInfo,
	})
}

