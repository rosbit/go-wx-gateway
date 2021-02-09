package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
	"strconv"
	"net/http"
)

// GET ${commonEndpoints.WxQr}?s=<service-name-in-conf>&t=<type-name,temp|forever>[&sceneid=xx][&e=<expire-secs-for-type-temp>]
func CreateWxQr(w http.ResponseWriter, r *http.Request) {
	service := r.FormValue("s")
	if service == "" {
		writeError(w, http.StatusBadRequest, "s(ervice) parameter expected")
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

	var ticketURL2ShowQrCode, urlIncluedInQrcode string
	var err error
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
		ticketURL2ShowQrCode, urlIncluedInQrcode, err = wxtools.CreateTempQrStrScene(service, sceneid, expireSecs)
	case "forever":
		ticketURL2ShowQrCode, urlIncluedInQrcode, err = wxtools.CreateQrStrScene(service, sceneid)
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
