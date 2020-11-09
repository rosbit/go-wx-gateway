package ce

import (
	"github.com/rosbit/go-wx-api/auth"
	"github.com/rosbit/go-wx-api/tools"
	"fmt"
	"net/http"
)

// POST ${commonEndpoints.TmplMsg}
// {
//    "s": "service-name-in-conf",
//    "to": "to-user-id",
//    "tid": "template-id",
//    "url": "optional url to jump",
//    "mp": {
//        "appid": "mini program appid",
//        "pagepath": "pagepath",
//    },
//    "data": {
//       "k1": "v1",
//       "k2": "v2",
//       "....": "..."
//    }
// }
func SendTmplMsg(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Service  string `json:"s"`
		ToUserId string `json:"to"`
		TmplId   string `json:"tid"`
		Url      string `json:"url"`
		MiniProg struct {
			AppId    string `json:"appid"`
			PagePath string `json:"pagepath"`
		} `json:"mp"`
		Data map[string]interface{} `json:"data"`
	}
	if status, err := readJson(r, &params); err != nil {
		writeError(w, status, err.Error())
		return
	}
	if params.Service == "" {
		writeError(w, http.StatusBadRequest, "s(ervice) parameter expected")
		return
	}

	wxParams, ok := wxParamsCache[params.Service]
	if !ok {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("unknown service name %s", params.Service))
		return
	}
	if params.ToUserId == "" {
		writeError(w, http.StatusBadRequest, "to(user id) parameter expected")
		return
	}
	if params.TmplId == "" {
		writeError(w, http.StatusBadRequest, "tid(template id) parameter expected")
		return
	}
	if len(params.Data) == 0 {
		writeError(w, http.StatusBadRequest, "data parameter as a map expected")
		return
	}

	accessToken, err := wxauth.NewAccessTokenWithParams(wxParams).Get()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res, err := wxtools.SendTemplateMessage(accessToken, params.ToUserId, params.TmplId, params.Data, params.Url, params.MiniProg.AppId, params.MiniProg.PagePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, res)
}

