/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/urfave/negroni"
	"github.com/gernest/alien"
	"net/http"
	"fmt"
	"os"
	"github.com/rosbit/go-wx-api"
	"github.com/rosbit/go-wx-api/conf"
	"github.com/rosbit/go-wx-api/msg"
	"wx-gateway/common-endpoints"
	"wx-gateway/conf"
	"wx-gateway/handlers"
)

func StartWxGateway() error {
	api := negroni.New()
	api.Use(negroni.NewRecovery())
	api.Use(negroni.NewLogger())

	router := alien.New()
	serviceConf := gwconf.ServiceConf
	for _, service := range serviceConf.Services {
		paramConf := service.WxParams
		wxParams, err := wxconf.NewWxParams(paramConf.Token, paramConf.AppId, paramConf.AppSecret, paramConf.AesKey)
		if err != nil {
			return fmt.Errorf("failed to init servie %s: %v", service.Name, err)
		}
		ce.CacheWxParams(service.Name, wxParams)
		// wxParamsCache[service.Name] = wxParams

		// init wx API
		wxService := wxapi.InitWxAPIWithParams(wxParams, service.WorkerNum, os.Stdout)
		endpoints := service.Endpoints

		// add uri signature checker
		signatureChecker := wxapi.NewWxSignatureChecker(paramConf.Token, service.Timeout, []string{endpoints.ServicePath})
		api.Use(negroni.HandlerFunc(signatureChecker))

		// set router
		router.Get(endpoints.ServicePath,  wxService.Echo)
		router.Post(endpoints.ServicePath, wxService.Request)
		if len(endpoints.RedirectPath) > 0 {
			router.Get(endpoints.RedirectPath, wxService.Redirect)
		}

		// set msg handlers and menu redirector
		if len(service.MsgProxyPass) > 0 {
			msgHandler := gwhandlers.NewMsgHandler(service.MsgProxyPass, wxParams, serviceConf.DontAppendUserInfo)
			wxService.RegisterWxMsghandler(msgHandler)
		} else {
			wxService.RegisterWxMsghandler(wxmsg.MsgHandler)
		}

		if len(service.RedirectURL) > 0 {
			if len(endpoints.RedirectPath) == 0 {
				return fmt.Errorf("listen-endpoints/redirect-path in servie %s must be specfied if you want to use redirect-url", service.Name)
			}
			wxService.RegisterRedirectUrl(service.RedirectURL, service.RedirectUserInfoFlag)
		}
	}

	commonEndpoints := &serviceConf.CommonEndpoints
	if len(commonEndpoints.HealthCheck) > 0 {
		router.Get(commonEndpoints.HealthCheck, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "OK\n")
		})
	}
	if len(commonEndpoints.WxQr) > 0 {
		router.Get(commonEndpoints.WxQr, ce.CreateWxQr)
	}
	if len(commonEndpoints.WxUser) > 0 {
		router.Get(commonEndpoints.WxUser, ce.GetWxUserInfo)
	}
	if len(commonEndpoints.SnsAPI) > 0 {
		router.Get(commonEndpoints.SnsAPI, ce.SnsAPI)
	}
	if len(commonEndpoints.ShortUrl) > 0 {
		router.Post(commonEndpoints.ShortUrl, ce.CreateShorturl)
	}
	if len(commonEndpoints.TmplMsg) > 0 {
		router.Post(commonEndpoints.TmplMsg, ce.SendTmplMsg)
	}
	if len(commonEndpoints.SignJSAPI) > 0  {
		router.Post(commonEndpoints.SignJSAPI, ce.SignJSAPI)
	}
	api.UseHandler(router)

	if len(serviceConf.TokenCacheDir) > 0  {
		wxconf.TokenStorePath = serviceConf.TokenCacheDir
	}

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}

