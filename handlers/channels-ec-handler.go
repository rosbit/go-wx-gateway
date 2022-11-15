package gwhandlers

import (
	"github.com/rosbit/go-wx-api/v2/msg"
	"fmt"
)

type ChannelsEcEventhandler struct {
	service string
	proxyPass string
}

func NewChannelsEcHandler(service, proxyPass string) *ChannelsEcEventhandler {
	return &ChannelsEcEventhandler{service:service, proxyPass:proxyPass}
}

func (h *ChannelsEcEventhandler) jsonCall(fromUser, toUser, url string, receivedMsg wxmsg.ReceivedJSONEvent) []byte {
	var res Res
	if err := JsonCall(url, "POST", receivedMsg, &res); err != nil {
		fmt.Printf("failed to JsonCall(%s): %v\n", url, err)
		return nil
	}
	msg := res.Msg
	return []byte(msg)
}

func (h *ChannelsEcEventhandler) jsonCallEvent(fromUser, toUser, eventType string, receivedMsg wxmsg.ReceivedJSONEvent) []byte {
	url := fmt.Sprintf("%s/event/%s", h.proxyPass, eventType)
	return h.jsonCall(fromUser, toUser, url, receivedMsg)
}

func (h *ChannelsEcEventhandler) HandleOrderCancelEvent(event *wxmsg.OrderCancelEvent) []byte {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *ChannelsEcEventhandler) HandleOrderPayEvent(event *wxmsg.OrderPayEvent) []byte {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *ChannelsEcEventhandler) HandleOrderConfirmEvent(event *wxmsg.OrderConfirmEvent) []byte {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *ChannelsEcEventhandler) HandleOrderSettleEvent(event *wxmsg.OrderSettleEvent) []byte {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}

func (h *ChannelsEcEventhandler) HandleAftersaleUpdateEvent (event *wxmsg.AftersaleUpdateEvent) []byte {
	return h.jsonCallEvent(event.FromUserName, event.ToUserName, event.Event, event)
}
