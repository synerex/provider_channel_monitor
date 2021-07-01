module github.com/synerex/provider_channel_monitor

go 1.16


require (
	github.com/synerex/synerex_sxutil v0.6.7
	github.com/synerex/synerex_proto v0.1.12
	github.com/synerex/provider_channel_monitor/proto_view v0.0.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/synerex/proto_fluentd v0.1.1 // indirect
	github.com/synerex/proto_grideye v0.0.2 // indirect
	github.com/synerex/provider_channel_monitor/proto_view/view_grideye v0.0.0
	github.com/synerex/provider_channel_monitor/proto_view/view_pcounter v0.0.0
	github.com/synerex/provider_channel_monitor/proto_view/view_fluent_wifi v0.0.0
)

replace github.com/synerex/provider_channel_monitor/proto_view => ./proto_view

replace github.com/synerex/provider_channel_monitor/proto_view/view_grideye => ./proto_view/view_grideye

replace github.com/synerex/provider_channel_monitor/proto_view/view_pcounter => ./proto_view/view_pcounter

replace github.com/synerex/provider_channel_monitor/proto_view/view_fluent_wifi => ./proto_view/view_fluent_wifi
