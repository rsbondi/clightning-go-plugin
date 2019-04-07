package main

import (
	"fmt"
	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
	"log"
	"os"
)

type Forwards struct{}

type ForwardSplit struct {
	Chins        interface{} `json:"channels_in"`
	Chouts       interface{} `json:"channels_out"`
	TotalFunding uint64      `json:"totalfunding"`
	TotalFees    uint64      `json:"totalfees"`
	PercentGain  float64     `json:"total_percent_gain"`
}

type ForwardChan struct {
	MsatForwarded uint64  `json:"msat"`
	Funding       uint64  `json:"funding"`
	PercentGain   float64 `json:"percent_gain"`
	PercentPie    float64 `json:"percent_pie"`
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

	chinsfinal := make(map[string]*ForwardChan, 0)
	choutsfinal := make(map[string]*ForwardChan, 0)

	var totalfees uint64
	for k, _ := range chins {
		fees := uint64(0)
		for _, f := range chins[k] {
			fees += f.Fee
		}
		totalfees += fees
		chinsfinal[k] = &ForwardChan{
			MsatForwarded: fees,
			Funding:       funds[k],
			PercentGain:   0,
			PercentPie:    0,
		}
	}

	for k, f := range chinsfinal { // we have total fees now, so calc pie
		chinsfinal[k].PercentPie = float64(f.MsatForwarded) / float64(totalfees)
	}

	for k, _ := range chouts {
		fees := uint64(0)
		for _, f := range chouts[k] {
			fees += f.Fee
		}
		choutsfinal[k] = &ForwardChan{
			MsatForwarded: fees,
			Funding:       funds[k],
			PercentGain:   float64(fees) / float64(funds[k]),
			PercentPie:    float64(fees) / float64(totalfees),
		}
	}

	c := ForwardSplit{
		Chins:        chinsfinal,
		Chouts:       choutsfinal,
		TotalFunding: totalfunding,
		TotalFees:    totalfees,
		PercentGain:  float64(totalfees) / float64(totalfunding),
	}

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
