package view_grideye

import (
	//	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	view "github.com/synerex/provider_channel_monitor/proto_view"

	grideye "github.com/synerex/proto_grideye"

	api "github.com/synerex/synerex_api"

	sxutil "github.com/synerex/synerex_sxutil"
)

func init() {
	fmt.Printf("Initial view GridEye\n")
	view.AddSubscriber(19, subscribeGridEyeSupply)
}

func supplyGridEyeCallback(clt *sxutil.SXServiceClient, sp *api.Supply) {

	ge := &grideye.GridEye{}
	if sp.Cdata != nil {
		err := proto.Unmarshal(sp.Cdata.Entity, ge)
		if err == nil { // get GridEye Data
			ts0 := ptypes.TimestampString(ge.Ts)
			ld := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%d", ts0, ge.DeviceId, ge.Hostname, ge.Location, ge.Mac, ge.Ip, ge.Seq)
			for _, ev := range ge.Data {
				ts := ptypes.TimestampString(ev.Ts)
				line := fmt.Sprintf("%s,%s,%s,%s,%d,%v", ld, ts, ev.Typ, ev.Id, ev.Seq, ev.Temps)
				log.Printf("GridEye:%s", line)
			}
			return
		} 
	}
	log.Printf("Unmarshal error on View_Pcoutner %s", sp.SupplyName)
}

func subscribeGridEyeSupply(client *sxutil.SXServiceClient) {
	//
	log.Printf("SubscribeGridEyeSupply\n")
	sxutil.SimpleSubscribeSupply(client, supplyGridEyeCallback) // error prone..
	//	ctx := context.Background() //
	//	client.SubscribeSupply(ctx, supplyGridEyeCallback)
	//	log.Printf("Error on subscribe with GridEye\n")
}
