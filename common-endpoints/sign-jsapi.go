package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
	"github.com/rosbit/mgin"
	"net/http"
)

// POST ${commonEndpoints.SignJSAPI}
// s=<service-name-in-conf>&u=<url-calling-jsapi-in-urlencoding>
func SignJSAPI(c *mgin.Context) {
	var params struct {
		Service string `form:"s"`
		Url     string `form:"u"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	nonce, timestamp, signature, err := wxtools.SignJSAPI(params.Service, params.Url)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"params": map[string]interface{}{
			"nonce": nonce,
			"timestamp": timestamp,
			"signature": signature,
		},
	})
}

