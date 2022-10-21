package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
	"github.com/rosbit/mgin"
	"net/http"
)

// GET ${commonEndpoints.WxQr}?s=<service-name-in-conf>&t=<type-name,temp|forever>[&sceneid=xx][&e=<expire-secs-for-type-temp>]
func CreateWxQr(c *mgin.Context) {
	var params struct {
		Service string `query:"s"`
		QrType  string `query:"t"`
		SceneId string `query:"sceneid" optional:"true"`
		ExpireSecs int `query:"e" optional:"true"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	switch params.QrType {
	case "temp", "forever":
	default:
		c.Error(http.StatusBadRequest, `t(ype) value must be "temp" or "forever"`)
		return
	}

	if params.SceneId == "" {
		params.SceneId = "0"
	}

	var ticketURL2ShowQrCode, urlIncluedInQrcode string
	var err error
	switch params.QrType {
	case "temp":
		expireSecs := 30
		if params.ExpireSecs > 0 {
			expireSecs = params.ExpireSecs
		}
		ticketURL2ShowQrCode, urlIncluedInQrcode, err = wxtools.CreateTempQrStrScene(params.Service, params.SceneId, expireSecs)
	case "forever":
		ticketURL2ShowQrCode, urlIncluedInQrcode, err = wxtools.CreateQrStrScene(params.Service, params.SceneId)
	}

	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"result": map[string]string {
			"ticketURL2ShowQrCode": ticketURL2ShowQrCode,
			"urlIncluedInQrcode": urlIncluedInQrcode,
		},
	})
}
