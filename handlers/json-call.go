package gwhandlers

import (
	"github.com/rosbit/gnet"
	"os"
)

type Res struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

func JsonCall(url string, method string, postData interface{}, res *Res) (err error) {
	_, err = gnet.JSONCallJ(url, &res, gnet.M(method), gnet.Params(postData), gnet.BodyLogger(os.Stderr))
	return
}
