package main

import (
	"fmt"
	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
	"log"
	"os"
)

type Forwards struct {
	Channel string
}

type ForwardSplit struct {
	Chins        interface{} `json:"chins"`
	Chouts       interface{} `json:"chouts"`
	TotalFunding uint64      `json:"totalfunding"`
	TotalFees    uint64      `json:"totalfees"`
	Percent      float64     `json:"totalpercent"`
}

type ForwardChan struct {
	MsatForwarded uint64
	Funding       uint64
	Percent       float64
}

var lightning *glightning.Lightning
var myid string

func (f *Forwards) New() interface{} {
	return &Forwards{}
}

func (f *Forwards) Name() string {
	return "forwardstats"
}

func (z *Forwards) Call() (jrpc2.Result, error) {

	forwards, err := lightning.ListForwards()

	if err != nil {
		return fmt.Sprintf("forward: %s\n", err.Error()), nil
	}

	if len(forwards) == 0 {
		return "no forwarding information available", nil
	}

	peers, err := lightning.ListPeers()
	if err != nil {
		return fmt.Sprintf("forward: %s\n", err.Error()), nil
	}

	funds := make(map[string]uint64, 0)
	var totalfunding uint64
	for _, p := range peers {
		for _, c := range p.Channels {
			funds[c.ShortChannelId] = c.FundingAllocations[myid]
			totalfunding += c.FundingAllocations[myid]
		}
	}

	chins := make(map[string][]glightning.Forwarding, 0)
	chouts := make(map[string][]glightning.Forwarding, 0)
	for _, f := range forwards {
		log.Printf("forward: %s\n", f)
		chins[f.InChannel] = append(chins[f.InChannel], f)
		chouts[f.OutChannel] = append(chouts[f.OutChannel], f)
	}

	chinsout := make(map[string]ForwardChan, 0)
	choutsout := make(map[string]ForwardChan, 0)

	for k, _ := range chins {
		fees := uint64(0)
		for _, f := range chins[k] {
			fees += f.Fee
		}
		chinsout[k] = ForwardChan{fees, funds[k], 0}
	}

	var totalfees uint64
	for k, _ := range chouts {
		fees := uint64(0)
		for _, f := range chouts[k] {
			fees += f.Fee
		}
		totalfees += fees
		choutsout[k] = ForwardChan{fees, funds[k], float64(funds[k]+fees) / float64(funds[k])}
	}

	c := ForwardSplit{chinsout, choutsout, totalfunding, totalfees, float64(totalfunding+totalfees) / float64(totalfunding)}

	return c, nil
}

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
	rpcForwards := glightning.NewRpcMethod(&Forwards{}, "A bunch of stuff about {channel}!")
	rpcForwards.LongDesc = `Various metrics about routing `
	rpcForwards.Usage = "[channel]"
	p.RegisterMethod(rpcForwards)
}

func OnConnect(c *glightning.ConnectEvent) {
	log.Printf("connect called: id %s at %s:%d", c.PeerId, c.Address.Addr, c.Address.Port)
}

func OnDisconnect(d *glightning.DisconnectEvent) {
	log.Printf("disconnect called for %s\n", d.PeerId)
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
