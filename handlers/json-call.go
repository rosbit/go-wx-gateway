package gwhandlers

import (
	"github.com/rosbit/gnet"
)

func JsonCall(url string, method string, postData interface{}) (res map[string]interface{}, err error) {
	_, err = gnet.JSONCallJ(url, &res, gnet.M(method), gnet.Params(postData))
	return
}
