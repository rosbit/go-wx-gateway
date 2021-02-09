package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
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

	longUrl := r.FormValue("u")
	if longUrl == "" {
		writeError(w, http.StatusBadRequest, "u(rl) parameter expected")
		return
	}

	shortUrl, err := wxtools.MakeShorturl(service, longUrl)
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

