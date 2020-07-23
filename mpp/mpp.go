package main

import (
	"log"
	"os"

	"github.com/niftynei/glightning/glightning"
)

var plugin *glightning.Plugin
var lightning *glightning.Lightning
var myid string

func main() {
	plugin = glightning.NewPlugin(onInit)
	lightning = glightning.NewLightning()

	registerMethods(plugin)
	err := plugin.Start(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func registerMethods(p *glightning.Plugin) {
	rpcPayment := glightning.NewRpcMethod(&PaymentMpp{}, "A bunch of stuff about payments!")
	rpcPayment.LongDesc = `Various metrics about payments `
	p.RegisterMethod(rpcPayment)

}

func onInit(plugin *glightning.Plugin, options map[string]glightning.Option, config *glightning.Config) {
	log.Printf("successfully init'd! %s\n", config.RpcFile)
	lightning.StartUp(config.RpcFile, config.LightningDir)
	info, err := lightning.GetInfo()

	if err != nil {
		log.Printf("forward: %s\n", err.Error())
	}
	myid = info.Id
}
