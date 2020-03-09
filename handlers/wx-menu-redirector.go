/**
 * 微信网页授权处理，用于处理服务号菜单点击事件
 * Rosbit Xu
 */
package gwhandlers

import (
	"github.com/rosbit/go-wx-api/auth"
)

/**
 * @deprecated
 * 为了更充分发挥菜单处理的能力，请在配置文件中加上"menu-redirect-url"，该配置将完全忽略下面的实现。
 *
 * 根据服务号菜单state做跳转
 * @param appId   公众号的appId
 * @param openId  订阅用户的openId
 * @param state   微信网页授权中的参数，用来标识某个菜单
 * @return
 *   c    需要显示服务号对话框中的内容
 *   h    需要在微信内嵌浏览器中设置的header信息，包括Cookie
 *   r    需要通过302跳转的URL。如果r不是空串，c的内容被忽略
 *   err  如果没有错误返回nil，非nil表示错误
 */
func CreateMenuRedirector(menuHandler string) wxauth.RedirectHandler {
	return func(appId, openId, state string) (c string, h map[string]string, r string, err error) {
		res, e := JsonCall(menuHandler, "POST", map[string]string{"appId": appId, "openId": openId, "state": state})
		if e != nil {
			err = e
			return
		}

		if cc, ok := res["c"]; ok {
			c = cc.(string)
		}
		if hh, ok := res["h"]; ok {
			h1 := hh.(map[string]interface{})
			h = make(map[string]string, len(h1))
			for k, v := range h1 {
				h[k] = v.(string)
			}
		}
		if rr, ok := res["r"]; ok {
			r = rr.(string)
		}

		return
	}
}
