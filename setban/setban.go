package main

import (
	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
	"log"
	"os"
	"time"
)

type SetBan struct {
	Id      string `json:"id"`
	Command string `json:"command"`
	Bantime int    `json:"bantime,omitempty"`
}

func (b *SetBan) New() interface{} {
	return &SetBan{}
}

func (b *SetBan) Name() string {
	return "setban"
}

func (b *SetBan) Call() (jrpc2.Result, error) {
	log.Printf("set ban called %s %s", b.Command, b.Id)

	if b.Command == "add" {
		log.Printf("adding ban for %s", b.Id)
		now := time.Now()
		ban := &Banned{
			Id:          b.Id,
			BanCreated:  now.Unix(),
			BannedUntil: now.Unix() + DEFAULT_BAN_TIME,
		}

		err := lightning.Disconnect(b.Id, true)
		if err != nil {
			log.Printf("disconnect error: %s", err.Error())
		}

		banned[ban.Id] = ban
	} else if b.Command == "remove" {
		delete(banned, b.Id)
	}

	return listbanned(), nil
}

type Banned struct {
	Id          string `json:"id"`
	BannedUntil int64  `json:"banned_until"`
	BanCreated  int64  `json:"ban_created"`
}

type ListBanned struct{}

func (b *ListBanned) New() interface{} {
	return &ListBanned{}
}

func (b *ListBanned) Name() string {
	return "listbanned"
}

func (b *ListBanned) Call() (jrpc2.Result, error) {
	return listbanned(), nil
}

func listbanned() []*Banned {
	bans := make([]*Banned, 0)
	for _, v := range banned {
		bans = append(bans, v)
	}
	return bans
}

var lightning *glightning.Lightning
var plugin *glightning.Plugin
var banned map[string]*Banned

const DEFAULT_BAN_TIME = 60 * 60 * 24

func main() {
	plugin = glightning.NewPlugin(onInit)
	lightning = glightning.NewLightning()
	banned = make(map[string]*Banned, 0)

	registerMethods(plugin)
	registerSubscriptions(plugin)
	err := plugin.Start(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func onInit(plugin *glightning.Plugin, options map[string]string, config *glightning.Config) {
	log.Printf("successfully init'd! %s\n", config.RpcFile)
	lightning.StartUp(config.RpcFile, config.LightningDir)
}

func registerMethods(p *glightning.Plugin) {
	rpcBan := glightning.NewRpcMethod(&SetBan{}, "Ban peers from connecting")
	rpcBan.LongDesc = `add or remove a ban for a peer by {id} {command(add|remove)} with optional {bantime}`
	rpcBan.Category = "peers"
	p.RegisterMethod(rpcBan)

	rpcListBanned := glightning.NewRpcMethod(&ListBanned{}, "List of banned peers")
	rpcListBanned.LongDesc = `shows list`
	rpcListBanned.Category = "peers"
	p.RegisterMethod(rpcListBanned)

}

func OnConnect(c *glightning.ConnectEvent) {
	log.Printf("connect called: id %s at %s:%d", c.PeerId, c.Address.Addr, c.Address.Port)
	if _, ok := banned[c.PeerId]; ok {
		err := lightning.Disconnect(c.PeerId, true)
		if err != nil {
			log.Printf("disconnect error: %s", err.Error())
		}
	}
}

func registerSubscriptions(p *glightning.Plugin) {
	p.SubscribeConnect(OnConnect)
}
