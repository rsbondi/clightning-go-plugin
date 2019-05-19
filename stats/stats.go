package main

import (
	"log"
	"os"

	"github.com/niftynei/glightning/glightning"
)

var plugin *glightning.Plugin

func main() {
	plugin = glightning.NewPlugin(onInit)
	lightning = glightning.NewLightning()

	registerOptions(plugin)
	registerMethods(plugin)
	err := plugin.Start(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func registerOptions(p *glightning.Plugin) {
	//	p.RegisterOption(glightning.NewOption("name", "How you'd like to be called", "Mary"))
}

func registerMethods(p *glightning.Plugin) {
	rpcForwards := glightning.NewRpcMethod(&Forwards{}, "A bunch of stuff about forwarding!")
	rpcForwards.LongDesc = `Various metrics about forwarding `
	rpcForwards.Usage = " "
	p.RegisterMethod(rpcForwards)

	rpcForwardView := glightning.NewRpcMethod(&ForwardView{}, "View of stuff about forwarding!")
	rpcForwardView.LongDesc = `View of various metrics about forwarding `
	rpcForwardView.Usage = " "
	p.RegisterMethod(rpcForwardView)

	rpcPayment := glightning.NewRpcMethod(&Payment{}, "A bunch of stuff about payments!")
	rpcPayment.LongDesc = `Various metrics about payments `
	rpcPayment.Usage = " "
	p.RegisterMethod(rpcPayment)

}

func onInit(plugin *glightning.Plugin, options map[string]string, config *glightning.Config) {
	log.Printf("successfully init'd! %s\n", config.RpcFile)
	lightning.StartUp(config.RpcFile, config.LightningDir)
	info, err := lightning.GetInfo()

	if err != nil {
		log.Printf("forward: %s\n", err.Error())
	}
	myid = info.Id
}
