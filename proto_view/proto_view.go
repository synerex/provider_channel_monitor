package proto_view

import {
	sxutil "github.com/synerex/synerex_sxutil"

}

var channelSubscribers := make(map[int]func(client *sxutil.SXServiceClient))


func addSubscriber(channel int, client *sxutil.SXServiceClient){
	fmt.Printf("View Sub: %d",channel)
	channelSubscribers[channel] = client
}

func subscribeAll(client *sxutil.SXSynerexClient  ){

	for ch, subscFunc := range channelSubscribers {
		chStr := fmt.Sprintf("Clt:GridMon:%d",ch)
		chClient := sxutil.NewSXServiceClient(client, ch, chStr)
		go subscFunc(chClient)
	}

}