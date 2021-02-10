package gwhandlers

import (
	"github.com/rosbit/go-wget"
)

func JsonCall(url string, method string, postData interface{}) (res map[string]interface{}, err error) {
	_, err = wget.JsonCallJ(url, method, postData, nil, &res)
	return
}
