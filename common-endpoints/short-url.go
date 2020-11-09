package ce

import (
	"github.com/rosbit/go-wx-api/auth"
	"github.com/rosbit/go-wx-api/tools"
	"fmt"
	"net/http"
)

// POST ${commonEndpoints.ShortUrl}
// s=<service-name-in-conf>&u=<long-url>
func CreateShorturl(w http.ResponseWriter, r *http.Request) {
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

	longUrl := r.FormValue("u")
	if longUrl == "" {
		writeError(w, http.StatusBadRequest, "u(rl) parameter expected")
		return
	}

	accessToken, err := wxauth.NewAccessTokenWithParams(wxParams).Get()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	shortUrl, err := wxtools.MakeShorturl(accessToken, longUrl)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJson(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"short-url": shortUrl,
	})
}

