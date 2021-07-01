module view_fluent_wifi

go 1.16

require (
	github.com/golang/protobuf v1.5.2
	github.com/synerex/proto_fluentd v0.1.1
	github.com/synerex/provider_channel_monitor/proto_view v0.0.0
	github.com/synerex/synerex_api v0.4.3
	github.com/synerex/synerex_sxutil v0.6.7
	google.golang.org/genproto v0.0.0-20210630183607-d20f26d13c79 // indirect
	google.golang.org/grpc v1.39.0 // indirect
)

replace github.com/synerex/provider_channel_monitor/proto_view => ../../proto_view
