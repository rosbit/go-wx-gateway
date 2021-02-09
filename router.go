/**
 * REST API router
 * Rosbit Xu
 */
package main

import (
	"github.com/rosbit/go-wx-api/v2/msg"
	"github.com/rosbit/go-wx-api/v2/log"
	"github.com/rosbit/go-wx-api/v2"
	"github.com/urfave/negroni"
	"github.com/go-zoo/bone"
	"wx-gateway/common-endpoints"
	"wx-gateway/handlers"
	"wx-gateway/conf"
	"fmt"
	"os"
	"net/http"
)

func StartWxGateway() error {
	api := negroni.New()
	api.Use(negroni.NewRecovery())
	api.Use(negroni.NewLogger())

	router := bone.New()
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
		// add uri signature checker
		signatureChecker := wxapi.NewWxSignatureChecker(paramConf.Token, service.Timeout, []string{endpoints.ServicePath})
		api.Use(negroni.HandlerFunc(signatureChecker))

		// set echo handler
		router.Get(endpoints.ServicePath,  wxapi.CreateEcho(paramConf.Token))

		// set msg handlers
		var msgHandler wxmsg.WxMsgHandler
		if len(service.MsgProxyPass) > 0 {
			msgHandler = gwhandlers.NewMsgHandler(service.Name, service.MsgProxyPass, serviceConf.DontAppendUserInfo)
		} else {
			msgHandler = wxmsg.MsgHandler
		}
		router.Post(endpoints.ServicePath, wxapi.CreateMsgHandler(service.Name, service.WorkerNum, msgHandler))

		// set oauth2 rediretor
		if len(service.RedirectURL) > 0 {
			if len(endpoints.RedirectPath) == 0 {
				return fmt.Errorf("listen-endpoints/redirect-path in servie %s must be specfied if you want to use redirect-url", service.Name)
			}
			router.Get(endpoints.RedirectPath, wxapi.CreateOAuth2Redirector(service.Name, service.WorkerNum, service.RedirectURL, service.RedirectUserInfoFlag))
		}
	}

	commonEndpoints := &serviceConf.CommonEndpoints
	if len(commonEndpoints.HealthCheck) > 0 {
		router.Get(commonEndpoints.HealthCheck, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintf(w, "OK\n")
		}))
	}
	if len(commonEndpoints.WxQr) > 0 {
		router.Get(commonEndpoints.WxQr, http.HandlerFunc(ce.CreateWxQr))
	}
	if len(commonEndpoints.WxUser) > 0 {
		router.Get(commonEndpoints.WxUser, http.HandlerFunc(ce.GetWxUserInfo))
	}
	if len(commonEndpoints.SnsAPI) > 0 {
		router.Get(commonEndpoints.SnsAPI, http.HandlerFunc(ce.SnsAPI))
	}
	if len(commonEndpoints.ShortUrl) > 0 {
		router.Post(commonEndpoints.ShortUrl, http.HandlerFunc(ce.CreateShorturl))
	}
	if len(commonEndpoints.TmplMsg) > 0 {
		router.Post(commonEndpoints.TmplMsg, http.HandlerFunc(ce.SendTmplMsg))
	}
	if len(commonEndpoints.SignJSAPI) > 0  {
		router.Post(commonEndpoints.SignJSAPI, http.HandlerFunc(ce.SignJSAPI))
	}
	api.UseHandler(router)

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}

