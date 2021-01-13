package main

import (
	"log"
	"os"

	"github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src"
	log2 "github.com/ConsenSys/orchestrate-hashicorp-vault-plugin/src/log"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/plugin"
)

func main() {
	client := &api.PluginAPIClientMeta{}
	err := client.FlagSet().Parse(os.Args[1:])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log2.InitLogger()
	
	err = plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: src.NewVaultBackend,
		TLSProviderFunc:    api.VaultPluginTLSProvider(client.GetTLSConfig()),
		Logger:             log2.Default(),
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
