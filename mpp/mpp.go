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
	rpcPayment := glightning.NewRpcMethod(&PaymentMpp{}, "Group multi part payments!")
	rpcPayment.LongDesc = `This is an alternative to listsendpays to accomodate multi part payments.
Instead of showing all parts, they are grouped together by payment_hash and status.
It will default to showing only completed payments with the omptional boolean parameter {includeall} 
which will also include pending and failed`
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
