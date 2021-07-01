package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
	"sync"

	view "github.com/synerex/provider_channel_monitor/proto_view"
	_ "github.com/synerex/provider_channel_monitor/proto_view/view_fluent_wifi"
	_ "github.com/synerex/provider_channel_monitor/proto_view/view_grideye"
	_ "github.com/synerex/provider_channel_monitor/proto_view/view_pcounter"
	sxutil "github.com/synerex/synerex_sxutil"
)

var (
	nodesrv  = flag.String("nodesrv", "127.0.0.1:9990", "Node ID Server")
	local    = flag.String("local", "", "Local Synerex Server")
	channels = flag.String("channels", "19", "Monitor Channels")
	all      = flag.Bool("all", false, "Subscribe all views")
	mu       sync.Mutex
)

const dateFmt = "2006-01-02T15:04:05.999Z"

//dataServer(geClient)

func main() {
	log.Printf("ChannelMonitor(%s) built %s sha1 %s", sxutil.GitVer, sxutil.BuildTime, sxutil.Sha1Ver)
	flag.Parse()
	go sxutil.HandleSigInt()
	sxutil.RegisterDeferFunction(sxutil.UnRegisterNode)

	chst := strings.Split(*channels, ",")

	channelTypes := make([]uint32, len(chst))

	for i, v := range chst {
		vv, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("strconv err in channels")
		} else {
			channelTypes[i] = uint32(vv)
		}
	}

	srv, rerr := sxutil.RegisterNode(*nodesrv, "ChM["+*channels+"]", channelTypes, nil)

	if rerr != nil {
		log.Fatal("Can't register node:", rerr)
	}
	if *local != "" { // quick hack for AWS local network
		srv = *local
	}
	log.Printf("Connecting SynerexServer at [%s]", srv)

	//	wg := sync.WaitGroup{} // for syncing other goroutines

	client := sxutil.GrpcConnectServer(srv)

	if client == nil {
		log.Fatal("Can't connect Synerex Server")
	} else {
		log.Print("Connecting SynerexServer")
	}
	wg := sync.WaitGroup{}
	wg.Add(1)

	// receive all
	// we should filter views
	if *all {
		view.SubscribeAll(client)
	} else {
		view.SubscribeChannels(client, channelTypes)
	}

	wg.Wait()

}
