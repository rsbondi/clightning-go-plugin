package main

import (
	"github.com/niftynei/glightning/glightning"
	"log"
	"os"
)

/*
   "command": "getlog [level]",
   "command": "listchannels [short_channel_id]
   "command": "listconfigs [config]",
   "command": "listinvoices [label]",
   "command": "listnodes [id]",
   "command": "listpeers [id] [level]",
   "command": "listsendpays [bolt11] [payment_hash]",
   "command": "paystatus [bolt11]",
*/

var lightning *glightning.Lightning
var plugin *glightning.Plugin
var lightningdir string

func main() {
	plugin = glightning.NewPlugin(onInit)
	lightning = glightning.NewLightning()

	registerMethods(plugin)
	err := plugin.Start(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func onInit(plugin *glightning.Plugin, options map[string]string, config *glightning.Config) {
	log.Printf("successfully init'd! %s\n", config.RpcFile)
	lightning.StartUp(config.RpcFile, config.LightningDir)
	lightningdir = config.LightningDir
	channels, _ := lightning.ListChannels()
	log.Printf("You know about %d channels", len(channels))
	log.Printf("Is this initial node startup? %v\n", config.Startup)
}

func registerMethods(p *glightning.Plugin) {
	rpcLog := glightning.NewRpcMethod(&GetLogExt{}, "Show logs, filtered by {levels} (info|unusual|debug|io)")
	rpcLog.LongDesc = `Show logs, with optional log [{levels}] with any combination of these (info|unusual|debug|io) `
	rpcLog.Category = "utility"
	p.RegisterMethod(rpcLog)

	rpcChannels := glightning.NewRpcMethod(&ListChannelsExt{}, "Show channels {short_channel_ids} or {sources}") // TODO: what is source?
	rpcChannels.LongDesc = `TBD`
	rpcChannels.Category = "channels"
	p.RegisterMethod(rpcChannels)

	rpcConfig := glightning.NewRpcMethod(&ListConfigsExt{}, "List filtered configuration options filtered by [configs].")
	rpcConfig.LongDesc = `listconfigsext [configs]\nOutputs an object, with each field a config options\n(Option names which start with # are comments)\nFiltered with [configs], object only has those fields`
	rpcConfig.Category = "utility"
	p.RegisterMethod(rpcConfig)

	rpcInvoice := glightning.NewRpcMethod(&ListInvoicesExt{}, "Show invoice {label}")
	rpcInvoice.LongDesc = `TBD`
	rpcInvoice.Category = "payment"
	p.RegisterMethod(rpcInvoice)

	rpcNode := glightning.NewRpcMethod(&ListNodesExt{}, "Show node {ids} in our local network view")
	rpcNode.LongDesc = `TBD`
	rpcNode.Category = "network"
	p.RegisterMethod(rpcNode)

	rpcPeers := glightning.NewRpcMethod(&ListPeersExt{}, "Show current peers, if {level} is set, include logs for {ids}") // TODO: how are log levels represented?
	rpcPeers.LongDesc = `TBD`
	rpcPeers.Category = "network"
	p.RegisterMethod(rpcPeers)

	rpcPays := glightning.NewRpcMethod(&ListSendpaysExt{}, "Show sendpay, old and current, filtering by {payment_hash}.")
	rpcPays.LongDesc = `TBD`
	rpcPays.Category = "payment"
	p.RegisterMethod(rpcPays)

}
