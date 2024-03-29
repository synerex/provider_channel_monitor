package view_pcounter

import (
	//	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	view "github.com/synerex/provider_channel_monitor/proto_view"

	pcounter "github.com/synerex/proto_pcounter"

	api "github.com/synerex/synerex_api"

	sxutil "github.com/synerex/synerex_sxutil"
)

func init() {
	fmt.Printf("Initial view People Counter\n")
	view.AddSubscriber(11, subscribePcounterSupply)
}

func supplyPcounterCallback(clt *sxutil.SXServiceClient, sp *api.Supply) {

	pc := &pcounter.PCounter{}
	//	log.Printf("Unmarshal %v\nCdata:%v",sp, sp.Cdata)
	if sp.Cdata != nil {
		err := proto.Unmarshal(sp.Cdata.Entity, pc)
		if err == nil { // get GridEye Data
			ts0 := ptypes.TimestampString(pc.Ts)
			log.Printf("PCounter %s,%v", ts0, pc)
			return
		} 
	}
	log.Printf("Unmarshal error on View_Pcoutner %s", sp.SupplyName)
}

func subscribePcounterSupply(client *sxutil.SXServiceClient) {
	//
	log.Printf("SubscribePcounterSupply")
	//	ctx := context.Background() //
	sxutil.SimpleSubscribeSupply(client, supplyPcounterCallback) // error prone..	
	//	client.SubscribeSupply(ctx, supplyPcounterCallback)
	//	log.Printf("Error on subscribe with Pcounter")
}
