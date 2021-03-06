package proto_view

import (
	"fmt"

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
		go subscFunc(chClient)
	}

}

// SubscribeChannel subscribe selected channel
func SubscribeChannels(client *sxutil.SXSynerexClient, channelTypes []uint32) {

	for _, ch := range channelTypes {
		subscFunc, ok := channelSubscribers[int(ch)]
		if ok {
			chStr := fmt.Sprintf("Clt:ChMon:%d", ch)
			chClient := sxutil.NewSXServiceClient(client, uint32(ch), chStr)
			go subscFunc(chClient)
		}
	}

}
