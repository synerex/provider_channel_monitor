package proto_view

import (
	"fmt"
	"log"

	api "github.com/synerex/synerex_api"
	sxutil "github.com/synerex/synerex_sxutil"
)

var channelSubscribers = make(map[int]func(*sxutil.SXServiceClient))

// AddSubscriber set data
func AddSubscriber(channel int, clFunc func(*sxutil.SXServiceClient)) {
	fmt.Printf("View Sub: %d \n", channel)
	channelSubscribers[channel] = clFunc
}

// SubscribeAll subscribe Functions
func SubscribeAll(client *sxutil.SXSynerexClient) {

	for ch, subscFunc := range channelSubscribers {
		chStr := fmt.Sprintf("Clt:GridMon:%d", ch)
		chClient := sxutil.NewSXServiceClient(client, uint32(ch), chStr)
		subscFunc(chClient)
	}

}

func simpleSupplyCallback(clt *sxutil.SXServiceClient, sp *api.Supply) {

	clen := 0
	if sp.Cdata != nil {
		clen = len(sp.Cdata.Entity)
	}
	log.Printf("SupplyName [%s] arg[%s] len:%d id:%d", sp.SupplyName, sp.ArgJson, clen, sp.SenderId)
}

func simpleDemandCallback(clt *sxutil.SXServiceClient, dm *api.Demand) {

	clen := 0
	if dm.Cdata != nil {
		clen = len(dm.Cdata.Entity)
	}
	log.Printf("DemandName [%s] arg[%s] len:%d id:%d", dm.DemandName, dm.ArgJson, clen, dm.SenderId)
}

// SubscribeChannel subscribe selected channel
func SubscribeChannels(client *sxutil.SXSynerexClient, channelTypes []uint32) {

	for _, ch := range channelTypes {
		subscFunc, ok := channelSubscribers[int(ch)]
		chStr := fmt.Sprintf("Clt:ChMon:%d", ch)
		if ok {
			chClient := sxutil.NewSXServiceClient(client, uint32(ch), chStr)
			subscFunc(chClient)
		} else { // there is no specified viewer
			chClient := sxutil.NewSXServiceClient(client, uint32(ch), chStr)
			sxutil.SimpleSubscribeSupply(chClient, simpleSupplyCallback)
			sxutil.SimpleSubscribeDemand(chClient, simpleDemandCallback)
		}
	}

}
