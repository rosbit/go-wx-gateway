package gwhandlers

import (
	"github.com/rosbit/go-wget"
	"encoding/json"
	"fmt"
)

func JsonCall(url string, method string, postData interface{}) (map[string]interface{}, error) {
	status, content, _, err := wget.PostJson(url, method, postData, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("status %d", status)
	}

	var res map[string]interface{}
	if err = json.Unmarshal(content, &res); err != nil {
		return nil, err
	}
	return res, nil
}
