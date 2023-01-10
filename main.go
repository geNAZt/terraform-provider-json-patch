package main

import (
	"github.com/geNAZt/terraform-provider-json-patch/pkg/provider"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	tf5server "github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"google.golang.org/grpc"
)

func main() {
	opts := &plugin.ServeOpts{}
	grpcProviderFunc := func() tfprotov5.ProviderServer {
		return schema.NewGRPCProviderServer(provider.Provider())
	}

	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: plugin.Handshake,
		VersionedPlugins: map[int]goplugin.PluginSet{
			5: {
				plugin.ProviderPluginName: &tf5server.GRPCProviderPlugin{
					GRPCProvider: func() tfprotov5.ProviderServer {
						return grpcProviderFunc()
					},
				},
			},
		},
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			return grpc.NewServer(append(opts,
				grpc.MaxSendMsgSize(64<<20 /* 64MB */),
				grpc.MaxRecvMsgSize(64<<20 /* 64MB */))...)
		},
		Logger: opts.Logger,
		Test:   opts.TestConfig,
	})
}
