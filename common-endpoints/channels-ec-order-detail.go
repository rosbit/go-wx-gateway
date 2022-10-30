package ce

import (
	"github.com/rosbit/go-wx-api/v2/channels-ec-order"
	"github.com/rosbit/mgin"
	"net/http"
)

// GET ${commonEndpoints.ChannelsEcOrderDetail}?s=<service-name-in-conf>&o=<order-id>
func ChannelsEcOrderDetail(c *mgin.Context) {
	var params struct {
		Service string `query:"s"`
		OrderId string `query:"o"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	ord, err := ceord.GetOrderDetail(params.Service, params.OrderId)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"order": ord,
	})
}

