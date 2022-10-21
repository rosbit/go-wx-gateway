package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
	"github.com/rosbit/mgin"
	"net/http"
)

// POST ${commonEndpoints.TmplMsg}
// {
//    "s": "service-name-in-conf",
//    "to": "to-user-id",
//    "tid": "template-id",
//    "url": "optional url to jump",
//    "mp": {
//        "appid": "mini program appid",
//        "pagepath": "pagepath",
//    },
//    "data": {
//       "k1": "v1",
//       "k2": "v2",
//       "....": "..."
//    }
// }
func SendTmplMsg(c *mgin.Context) {
	var params struct {
		Service  string `json:"s"`
		ToUserId string `json:"to"`
		TmplId   string `json:"tid"`
		Url      string `json:"url"`
		MiniProg struct {
			AppId    string `json:"appid"`
			PagePath string `json:"pagepath"`
		} `json:"mp"`
		Data map[string]interface{} `json:"data"`
	}
	if status, err := c.ReadJSON(&params); err != nil {
		c.Error(status, err.Error())
		return
	}
	if params.Service == "" {
		c.Error(http.StatusBadRequest, "s(ervice) parameter expected")
		return
	}

	if params.ToUserId == "" {
		c.Error(http.StatusBadRequest, "to(user id) parameter expected")
		return
	}
	if params.TmplId == "" {
		c.Error(http.StatusBadRequest, "tid(template id) parameter expected")
		return
	}
	if len(params.Data) == 0 {
		c.Error(http.StatusBadRequest, "data parameter as a map expected")
		return
	}

	res, err := wxtools.SendTemplateMessage(params.Service, params.ToUserId, params.TmplId, params.Data, params.Url, params.MiniProg.AppId, params.MiniProg.PagePath)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

