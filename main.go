package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-provider-time/internal/provider"
	"log"

	"github.com/hashicorp/terraform-provider-time/internal/tftime"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-mux/tf6to5server"
)

func main() {
	ctx := context.Background()
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	downgradedProvider, err := tf6to5server.DowngradeServer(
		ctx,
		func() tfprotov6.ProviderServer {
			return providerserver.NewProtocol6(provider.New())()
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov5.ProviderServer{
		func() tfprotov5.ProviderServer {
			return downgradedProvider
		},
		tftime.Provider().GRPCProvider,
	}

	muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf5server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/hashicorp/time",
		muxServer.ProviderServer,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
