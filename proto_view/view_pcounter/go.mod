module view_pcounter

go 1.16

require (
	github.com/golang/protobuf v1.4.3
	github.com/synerex/synerex_api v0.4.2
	github.com/synerex/synerex_sxutil v0.6.3
)

require (
	github.com/synerex/proto_pcounter v0.0.6
	github.com/synerex/provider_channel_monitor/proto_view v0.0.0
)

replace github.com/synerex/provider_channel_monitor/proto_view => ../../proto_view
