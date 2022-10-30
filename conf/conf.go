/**
 * global conf
 * ENV:
 *   CONF_FILE      --- 配置文件名
 *   TZ             --- 时区名称"Asia/Shanghai"
 *
 * JSON:
 *  {
      "listen-host": "",
      "listen-port": 7080,
      "services": [
         {
			 "is-channels-ec": false,
             "name": "echo_server",
             "workerNum": 5,
             "timeout": 0,
             "wx-params": {
                 "token": "hello_rosbit",
                 "app-id": "",
                 "app-secret": "",
                 "aes-key":  null
             },
             "listen-endpoints": {
                 "service-path": "/wx",
                 "redirect-path": "/redirect"
             },
             "msg-proxy-pass": "http://yourhost.or.ip.here",
             "redirect-url": "http://yourhost.or.ip/path/to/redirect",
             "redirect-userinfo-flag": "login, register or any-strings else if you want use snsapi_userinfo",
         }
      ],
      "token-cache-dir": "/root/dir/to/save/token",
      "common-endpoints": {
          "health-check": "/health",
          "wx-qr": "/qr",
          "wx-user": "/userinfo",
          "sns-auth2": "/sns-auth2",
          "short-url": "/short-url",
          "tmpl-msg": "/tmpl-msg",
          "sign-jsapi": "/sign-jsapi",
		  "channels-ec-order-detail": "/channles-ec-order-detail",
		  "channels-ec-refund-detail": "/channels-ec-refund-detail"
      },
      "dont-append-userinfo": true
   }
 *
 * Rosbit Xu
 */
package gwconf

import (
	"fmt"
	"os"
	"time"
	"encoding/json"
)

type WxParamsConf struct {
	Token     string `json:"token"`
	AppId     string `json:"app-id"`
	AppSecret string `json:"app-secret"`
	AesKey    string `json:"aes-key"`
}

type WxServiceConf struct {
	ListenHost     string `json:"listen-host"`
	ListenPort     int    `json:"listen-port"`
	Services       []struct {
		IsChannelsEc bool   `json:"is-channels-ec"`
		Name         string `json:"name"`
		WorkerNum    int    `json:"workerNum"`
		Timeout      int    `json:"timeout"`
		WxParams     WxParamsConf `json:"wx-params"`
		Endpoints    struct {
			ServicePath  string `json:"service-path"`
			RedirectPath string `json:"redirect-path"`
		}  `json:"listen-endpoints"`
		MsgProxyPass string `json:"msg-proxy-pass"`
		RedirectURL string  `json:"redirect-url"`
		RedirectUserInfoFlag string `json:"redirect-userinfo-flag"`
	} `json:"services"`
	TokenCacheDir string `json:"token-cache-dir"`
	CommonEndpoints struct {
		HealthCheck string `json:"health-check"`
		WxQr        string `json:"wx-qr"`
		WxUser      string `json:"wx-user"`
		SnsAPI      string `json:"sns-auth2"`
		ShortUrl    string `json:"short-url"`
		TmplMsg     string `json:"tmpl-msg"`
		SignJSAPI   string `json:"sign-jsapi"`
		ChannelsEcOrderDetail string `json:"channels-ec-order-detail"`
		ChannelsEcRefundDetail string `json:"channels-ec-refund-detail"`
	} `json:"common-endpoints"`
	DontAppendUserInfo bool `json:"dont-append-userinfo"`
}

var (
	ServiceConf WxServiceConf
	Loc = time.FixedZone("UTC+8", 8*60*60)
)

func getEnv(name string, result *string, must bool) error {
	s := os.Getenv(name)
	if s == "" {
		if must {
			return fmt.Errorf("env \"%s\" not set", name)
		}
	}
	*result = s
	return nil
}

func CheckGlobalConf() error {
	var p string
	getEnv("TZ", &p, false)
	if p != "" {
		if loc, err := time.LoadLocation(p); err == nil {
			Loc = loc
		}
	}

	var confFile string
	if err := getEnv("CONF_FILE", &confFile, true); err != nil {
		return err
	}

	fp, err := os.Open(confFile)
	if err != nil {
		return err
	}
	defer fp.Close()

	dec := json.NewDecoder(fp)
	if err := dec.Decode(&ServiceConf); err != nil {
		return err
	}

	return nil
}

func DumpConf() {
	fmt.Printf("conf: %v\n", ServiceConf)
	fmt.Printf("TZ time location: %v\n", Loc)
}
