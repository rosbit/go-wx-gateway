package ce

import (
	"github.com/rosbit/go-wx-api/auth"
	"github.com/rosbit/go-wx-api/tools"
	"strconv"
	"fmt"
	"net/http"
)

// GET ${commonEndpoints.WxQr}?s=<service-name-in-conf>&t=<type-name,temp|forever>[&sceneid=xx][&e=<expire-secs-for-type-temp>]
func CreateWxQr(w http.ResponseWriter, r *http.Request) {
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

	qrType := r.FormValue("t")
	if qrType == "" {
		writeError(w, http.StatusBadRequest, "t(type) parameter expected")
		return
	}
	switch qrType {
	case "temp", "forever":
	default:
		writeError(w, http.StatusBadRequest, `t(ype) value must be "temp" or "forever"`)
		return
	}

	sceneid := r.FormValue("sceneid")
	if sceneid == "" {
		sceneid = "0"
	}

	accessToken, err := wxauth.NewAccessTokenWithParams(wxParams).Get()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var ticketURL2ShowQrCode, urlIncluedInQrcode string
	switch qrType {
	case "temp":
		expireSecs := 30
		e := r.FormValue("e")
		if e == "" {
			expireSecs, _ := strconv.Atoi(e)
			if expireSecs <= 0 {
				expireSecs = 30
			}
		}
		ticketURL2ShowQrCode, urlIncluedInQrcode, err = wxtools.CreateTempQrStrScene(accessToken, sceneid, expireSecs)
	case "forever":
		ticketURL2ShowQrCode, urlIncluedInQrcode, err = wxtools.CreateQrStrScene(accessToken, sceneid)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJson(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]string {
			"ticketURL2ShowQrCode": ticketURL2ShowQrCode,
			"urlIncluedInQrcode": urlIncluedInQrcode,
		},
	})
}
