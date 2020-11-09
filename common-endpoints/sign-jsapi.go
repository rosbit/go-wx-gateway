package ce

import (
	"github.com/rosbit/go-wx-api/auth"
	"github.com/rosbit/go-wx-api/tools"
	"fmt"
	"net/http"
)

// POST ${commonEndpoints.SignJSAPI}
// s=<service-name-in-conf>&u=<url-calling-jsapi-in-urlencoding>
func SignJSAPI(w http.ResponseWriter, r *http.Request) {
	service := r.FormValue("s")
	if service == "" {
		writeError(w, http.StatusBadRequest, "s(ervice) parameter expected")
		return
	}

	wxParams, ok := wxParamsCache[service]
	if !ok {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("unknown service name %s", service))
		return
	}

	url := r.FormValue("u")
	if url == "" {
		writeError(w, http.StatusBadRequest, "u(rl) parameter expected")
		return
	}

	accessToken, err := wxauth.NewAccessTokenWithParams(wxParams).Get()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	nonce, timestamp, signature, err := wxtools.SignJSAPI(accessToken, url)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"params": map[string]interface{}{
			"nonce": nonce,
			"timestamp": timestamp,
			"signature": signature,
		},
	})
}

