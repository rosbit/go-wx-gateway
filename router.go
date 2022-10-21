/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/rosbit/go-wx-api/v2/msg"
	"github.com/rosbit/go-wx-api/v2/log"
	"github.com/rosbit/go-wx-api/v2"
	"github.com/rosbit/mgin"
	"wx-gateway/common-endpoints"
	"wx-gateway/handlers"
	"wx-gateway/conf"
	"fmt"
	"os"
	"net/http"
)

func StartWxGateway() error {
	handlers := []mgin.Handler{mgin.WithLogger("wx-gateway")}

	serviceConf := gwconf.ServiceConf
	if len(serviceConf.TokenCacheDir) > 0  {
		wxapi.InitWx(serviceConf.TokenCacheDir)
	}
	wxlog.SetLogger(os.Stderr)

	for _, service := range serviceConf.Services {
		paramConf := &service.WxParams
		if err := wxapi.SetWxParams(service.Name, paramConf.Token, paramConf.AppId, paramConf.AppSecret, paramConf.AesKey); err != nil {
			return err
		}

		endpoints := &service.Endpoints
		// add uri signature checker as Handler
		signatureChecker := wxapi.NewWxSignatureChecker(paramConf.Token, service.Timeout, []string{endpoints.ServicePath})
		handlers = append(handlers, mgin.WrapMiddleFunc(signatureChecker))
	}
	api := mgin.NewMgin(handlers...)

	for _, service := range serviceConf.Services {
		paramConf := &service.WxParams
		endpoints := &service.Endpoints

		// set echo handler
		api.Get(endpoints.ServicePath,  wxapi.CreateEcho(paramConf.Token))

		// set msg handlers
		var msgHandler wxmsg.WxMsgHandler
		if len(service.MsgProxyPass) > 0 {
			msgHandler = gwhandlers.NewMsgHandler(service.Name, service.MsgProxyPass, serviceConf.DontAppendUserInfo)
		} else {
			msgHandler = wxmsg.MsgHandler
		}
		api.Post(endpoints.ServicePath, wxapi.CreateMsgHandler(service.Name, service.WorkerNum, msgHandler))

		// set oauth2 rediretor
		if len(service.RedirectURL) > 0 {
			if len(endpoints.RedirectPath) == 0 {
				return fmt.Errorf("listen-endpoints/redirect-path in servie %s must be specfied if you want to use redirect-url", service.Name)
			}
			api.Get(endpoints.RedirectPath, wxapi.CreateOAuth2Redirector(service.Name, service.WorkerNum, service.RedirectURL, service.RedirectUserInfoFlag))
		}
	}

	commonEndpoints := &serviceConf.CommonEndpoints
	if len(commonEndpoints.HealthCheck) > 0 {
		api.GET(commonEndpoints.HealthCheck, func(c *mgin.Context) {
			c.String(http.StatusOK, "OK\n")
		})
	}
	if len(commonEndpoints.WxQr) > 0 {
		api.GET(commonEndpoints.WxQr, ce.CreateWxQr)
	}
	if len(commonEndpoints.WxUser) > 0 {
		api.GET(commonEndpoints.WxUser, ce.GetWxUserInfo)
	}
	if len(commonEndpoints.SnsAPI) > 0 {
		api.GET(commonEndpoints.SnsAPI, ce.SnsAPI)
	}
	if len(commonEndpoints.ShortUrl) > 0 {
		api.POST(commonEndpoints.ShortUrl, ce.CreateShorturl)
	}
	if len(commonEndpoints.TmplMsg) > 0 {
		api.POST(commonEndpoints.TmplMsg, ce.SendTmplMsg)
	}
	if len(commonEndpoints.SignJSAPI) > 0  {
		api.POST(commonEndpoints.SignJSAPI, ce.SignJSAPI)
	}

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}

