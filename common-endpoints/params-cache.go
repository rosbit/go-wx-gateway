package ce

import (
	"github.com/rosbit/go-wx-api/conf"
)

var wxParamsCache = map[string]*wxconf.WxParamsT{}

func CacheWxParams(service string, wxParams *wxconf.WxParamsT) {
	wxParamsCache[service] = wxParams
}
