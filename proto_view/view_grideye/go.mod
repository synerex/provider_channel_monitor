module view_grideye

go 1.16

require (
	github.com/golang/protobuf v1.4.1
	github.com/synerex/proto_grideye v0.0.2
	github.com/synerex/synerex_api v0.3.1
	github.com/synerex/synerex_proto v0.1.10
	github.com/synerex/synerex_sxutil v0.4.9
)


require github.com/synerex/provider_channel_monitor/proto_view v0.0.0


replace github.com/synerex/provider_channel_monitor/proto_view  => ../../proto_view
