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
	"github.com/rosbit/go-wx-api/auth"
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

		// init wx API
		wxService := wxapi.InitWxAPIWithParams(wxParams, service.WorkerNum, os.Stdout)
		endpoints := service.Endpoints

		// add uri signature checker
		signatureChecker := wxapi.NewWxSignatureChecker(paramConf.Token, service.Timeout, []string{endpoints.ServicePath})
		api.Use(negroni.HandlerFunc(signatureChecker))

		// set router
		router.Get(endpoints.ServicePath,  wxService.Echo)
		router.Post(endpoints.ServicePath, wxService.Request)
		router.Get(endpoints.RedirectPath, wxService.Redirect)

		// set msg handlers and menu redirector
		if service.MsgProxyPass != "" {
			msgHandler := gwhandlers.NewMsgHandler(service.MsgProxyPass)
			wxService.RegisterWxMsghandler(msgHandler)
		} else {
			wxService.RegisterWxMsghandler(wxmsg.MsgHandler)
		}

		var menuRedirect wxauth.RedirectHandler
		if service.MenuHandler != "" {
			menuRedirect = gwhandlers.CreateMenuRedirector(service.MenuHandler)
		} else {
			menuRedirect = wxauth.ToAppIdRedirectHandler(wxauth.HandleRedirect)
		}
		wxService.RegisterRedictHandler(menuRedirect)
	}
	api.UseHandler(router)

	listenParam := fmt.Sprintf("%s:%d", serviceConf.ListenHost, serviceConf.ListenPort)
	fmt.Printf("%v\n", http.ListenAndServe(listenParam, api))
	return nil
}

