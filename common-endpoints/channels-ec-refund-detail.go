package ce

import (
	"github.com/rosbit/go-wx-api/v2/channels-ec-order"
	"github.com/rosbit/mgin"
	"net/http"
)

// GET ${commonEndpoints.ChannelsEcRefundDetail}?s=<service-name-in-conf>&o=<aftersale-order-id>
func ChannelsEcRefundDetail(c *mgin.Context) {
	var params struct {
		Service string `query:"s"`
		AftersaleOrderId string `query:"o"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	ord, err := ceord.GetRefundOrderDetail(params.Service, params.AftersaleOrderId)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"aftersaleOrder": ord,
	})
}

