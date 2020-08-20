package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/niftynei/glightning/glightning"
)

var plugin *glightning.Plugin
var lightning *glightning.Lightning
var myid string
var dbPath string

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
	p.RegisterMethod(rpcForwards)

	rpcForwardView := glightning.NewRpcMethod(&ForwardView{}, "View of stuff about forwarding!")
	rpcForwardView.LongDesc = `View of various metrics about forwarding `
	p.RegisterMethod(rpcForwardView)

	rpcPayment := glightning.NewRpcMethod(&Payment{}, "A bunch of stuff about payments!")
	rpcPayment.LongDesc = `Various metrics about payments `
	p.RegisterMethod(rpcPayment)

	rpcPaymentView := glightning.NewRpcMethod(&PaymentView{}, "A view of stuff about payments!")
	rpcPaymentView.LongDesc = `View of various metrics about payments `
	p.RegisterMethod(rpcPaymentView)

	rpcActivityView := glightning.NewRpcMethod(&ChannelActivity{}, "A view of stuff about channels!")
	rpcActivityView.LongDesc = `How sats flow through active channels`
	p.RegisterMethod(rpcActivityView)
}

func onInit(plugin *glightning.Plugin, options map[string]glightning.Option, config *glightning.Config) {
	log.Printf("successfully init'd! %s\n", config.RpcFile)
	lightning.StartUp(config.RpcFile, config.LightningDir)
	dbPath = filepath.Join(config.LightningDir, "lightningd.sqlite3")
	info, err := lightning.GetInfo()

	if err != nil {
		log.Printf("forward: %s\n", err.Error())
	}
	myid = info.Id
}
