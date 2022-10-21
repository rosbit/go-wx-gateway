package ce

import (
	"github.com/rosbit/go-wx-api/v2/tools"
	"github.com/rosbit/mgin"
	"net/http"
)

// POST ${commonEndpoints.ShortUrl}
// s=<service-name-in-conf>&u=<long-url>
func CreateShorturl(c *mgin.Context) {
	var params struct {
		Service string `form:"s"`
		LongUrl string `form:"u"`
	}
	if code, err := c.ReadParams(&params); err != nil {
		c.Error(code, err.Error())
		return
	}

	shortUrl, err := wxtools.MakeShorturl(params.Service, params.LongUrl)
	if err != nil {
		c.Error(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"msg": "OK",
		"short-url": shortUrl,
	})
}

