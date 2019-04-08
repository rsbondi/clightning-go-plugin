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
	TotalForward uint64      `json:"totalforward"`
	PercentGain  float64     `json:"total_percent_gain"`
}

type ForwardChan struct {
	MsatFees      uint64  `json:"fee_msat"`
	MsatForward   uint64  `json:"forward_msat"`
	FailedFees    uint64  `json:"fee_fail"`
	FailedForward uint64  `json:"forward_fail"`
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
	var totalfees uint64
	var totalforwards uint64
	for _, f := range forwards {
		chins[f.InChannel] = append(chins[f.InChannel], f)
		chouts[f.OutChannel] = append(chouts[f.OutChannel], f)
		if f.Status == "settled" {
			totalfees += f.Fee
			totalforwards += f.MilliSatoshiOut
		}
	}

	chinsfinal := make(map[string]*ForwardChan, 0)
	choutsfinal := make(map[string]*ForwardChan, 0)

	processChannels(chins, chinsfinal, funds, totalfees)
	processChannels(chouts, choutsfinal, funds, totalfees)

	for k, f := range chinsfinal { // we have total fees now, so calc pie
		chinsfinal[k].PercentPie = float64(f.MsatFees) / float64(totalfees)
	}

	c := ForwardSplit{
		Chins:        chinsfinal,
		Chouts:       choutsfinal,
		TotalFunding: totalfunding,
		TotalFees:    totalfees,
		TotalForward: totalforwards,
		PercentGain:  float64(totalfees) / float64(totalfunding),
	}

	return c, nil
}

func processChannels(src map[string][]glightning.Forwarding,
	dest map[string]*ForwardChan,
	funds map[string]uint64,
	totalfees uint64) {
	for k, _ := range src {
		fees := uint64(0)
		forwarded := uint64(0)
		feesfail := uint64(0)
		forwardedfail := uint64(0)
		for _, f := range src[k] {
			if f.Status == "settled" {
				fees += f.Fee
				forwarded += f.MilliSatoshiOut
			} else {
				feesfail += f.Fee
				forwardedfail += f.MilliSatoshiOut
			}
		}
		var gain float64
		if funds[k] > 0 {
			gain = float64(fees) / float64(funds[k])

		}

		dest[k] = &ForwardChan{
			MsatFees:      fees,
			MsatForward:   forwarded,
			FailedFees:    feesfail,
			FailedForward: forwardedfail,
			Funding:       funds[k],
			PercentGain:   gain,
			PercentPie:    float64(fees) / float64(totalfees),
		}
	}
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
