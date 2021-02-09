package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
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

	url := r.FormValue("u")
	if url == "" {
		writeError(w, http.StatusBadRequest, "u(rl) parameter expected")
		return
	}

	nonce, timestamp, signature, err := wxtools.SignJSAPI(service, url)
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

