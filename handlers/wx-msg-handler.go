package gwhandlers

import (
	"github.com/rosbit/go-wx-api/v2/msg"
	"github.com/rosbit/go-wx-api/v2/auth"
	"encoding/json"
	"io/ioutil"
	"bytes"
	"fmt"
)

type WxMsgHandler struct {
	service string
	proxyPass string
	dontAppendUserInfo bool
}

func NewMsgHandler(service, proxyPass string, dontAppendUserInfo bool) *WxMsgHandler {
	return &WxMsgHandler{service, proxyPass, dontAppendUserInfo}
}

func (h *WxMsgHandler) jsonCall(fromUser, toUser, url string, receivedMsg wxmsg.ReceivedMsg) wxmsg.ReplyMsg {
	var res map[string]interface{}
	var err error
	if h.dontAppendUserInfo {
		res, err = JsonCall(url, "POST", receivedMsg)
	} else {
		res, err = wxauth.GetUserInfo(h.service, fromUser)
		b := &bytes.Buffer{}               // b: <empty>
		je := json.NewEncoder(b)
		je.Encode(receivedMsg)             // b: {JSON}\n
		b.Truncate(b.Len()-2)              // b: {JSON      //  "}\n" removed
		b.WriteString(`,"userInfo":`)      // b: {JSON,"userInfo":
		je.Encode(res)                     // b: {JSON,"userInfo":JSON\n
		b.WriteString(`,"userInfoError":`) // b: {JSON,"userInfo":JSON\n,"userInfoError":
		je.Encode(func()string{
			if err == nil {
				return ""
			}
			return err.Error()
		}())                   // b: {JSON,"userInfo":JSON\n,"userInfoError":"xxx"\n
		b.WriteByte('}')       // b: {JSON,"userInfo":JSON\n,"userInfoError":"xxx"\n}
		res, err = JsonCall(url, "POST", ioutil.NopCloser(b))
	}
	if err != nil {
		fmt.Printf("failed to JsonCall(%s): %v\n", url, err)
		return nil
	}
	typ, ok := res["type"]
	if ! ok {
		typ = "text"
	}
	msg, ok := res["msg"];
	if !ok {
		fmt.Printf("no \"msg\" item in %v\n", res)
		return nil
	}
	switch typ {
	case "text":
		return wxmsg.NewReplyTextMsg(fromUser, toUser, msg.(string))
	case "image":
		return wxmsg.NewReplyImageMsg(fromUser, toUser, msg.(string))
	case "voice":
		return wxmsg.NewReplyVoiceMsg(fromUser, toUser, msg.(string))
	case "video":
		title, ok := res["title"]
		if !ok {
			title = "[no title]"
		}
		desc, ok := res["desc"]
		if !ok {
			desc = "[no desc]"
		}
		return wxmsg.NewReplyVideoMsg(fromUser, toUser, msg.(string), title.(string), desc.(string))
	case "success":
		return wxmsg.NewSuccessMsg()
	default:
		fmt.Printf("unknwon type %s", typ)
		return nil
	}
}

func (h *WxMsgHandler) jsonCallMsg(fromUser, toUser, msgType string, receivedMsg wxmsg.ReceivedMsg) wxmsg.ReplyMsg {
	url := fmt.Sprintf("%s/msg/%s", h.proxyPass, msgType)
	return h.jsonCall(fromUser, toUser, url, receivedMsg)
}

func (h *WxMsgHandler) jsonCallEvent(fromUser, toUser, eventType string, receivedMsg wxmsg.ReceivedMsg) wxmsg.ReplyMsg {
	url := fmt.Sprintf("%s/event/%s", h.proxyPass, eventType)
	return h.jsonCall(fromUser, toUser, url, receivedMsg)
}

func (h *WxMsgHandler) HandleTextMsg(msg *wxmsg.TextMsg) wxmsg.ReplyMsg {
	return h.jsonCallMsg(msg.FromUserName, msg.ToUserName, msg.MsgType, msg)
}

func (h *WxMsgHandler) HandleImageMsg(msg *wxmsg.ImageMsg) wxmsg.ReplyMsg {
	return h.jsonCallMsg(msg.FromUserName, msg.ToUserName, msg.MsgType, msg)
}

func (h *WxMsgHandler) HandleVoiceMsg(msg *wxmsg.VoiceMsg) wxmsg.ReplyMsg {
	return h.jsonCallMsg(msg.FromUserName, msg.ToUserName, msg.MsgType, msg)
}

func (h *WxMsgHandler) HandleVideoMsg(msg *wxmsg.VideoMsg) wxmsg.ReplyMsg {
	return h.jsonCallMsg(msg.FromUserName, msg.ToUserName, msg.MsgType, msg)
}

func (h *WxMsgHandler) HandleLocationMsg(msg *wxmsg.LocationMsg) wxmsg.ReplyMsg {
	return h.jsonCallMsg(msg.FromUserName, msg.ToUserName, msg.MsgType, msg)
}

func (h *WxMsgHandler) HandleLinkMsg(msg *wxmsg.LinkMsg) wxmsg.ReplyMsg {
	return h.jsonCallMsg(msg.FromUserName, msg.ToUserName, msg.MsgType, msg)
}

func (h *WxMsgHandler) HandleClickEvent(event *wxmsg.ClickEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.MsgType, event)
}

func (h *WxMsgHandler) HandleViewEvent(event *wxmsg.ViewEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleScancodePushEvent(event *wxmsg.ScancodeEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleScancodeWaitEvent(event *wxmsg.ScancodeEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleSubscribeEvent(event *wxmsg.SubscribeEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleUnsubscribeEvent(event *wxmsg.SubscribeEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleScanEvent(event *wxmsg.SubscribeEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleWhereEvent(event *wxmsg.WhereEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandlePhotoEvent(event *wxmsg.PhotoEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleLocationEvent(event *wxmsg.LocationEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleMassSentEvent(event *wxmsg.MassSentEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *WxMsgHandler) HandleTemplateSentEvent(event *wxmsg.TemplateSentEvent) wxmsg.ReplyMsg {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

