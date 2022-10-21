package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
	"github.com/rosbit/mgin"
	"net/http"
)

// GET ${commonEndpoints.WxUser}?s=<service-name-in-conf>&o=<openId>
func GetWxUserInfo(c *mgin.Context) {
	var params struct {
		Service string `query:"s"`
		OpenId  string `query:"o"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	userInfo, err := wxtools.GetUserInfo(params.Service, params.OpenId)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"userInfo": userInfo,
	})
}

