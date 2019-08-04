package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
	"golang.org/x/crypto/hkdf"
)

const VERSION = "0.0.1-WIP"

var plugin *glightning.Plugin
var lightning *glightning.Lightning
var lightningdir string
var bitcoinNet *chaincfg.Params

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

type ListConfigsRequest struct{}

func (c *ListConfigsRequest) Name() string {
	return "listconfigs"
}

type ListConfigsResponse struct {
	Network string `json:"network"`
}

func ListConfigs() (*ListConfigsResponse, error) {
	result := &ListConfigsResponse{}
	req := &ListConfigsRequest{}
	err := lightning.Request(req, result)
	return result, err
}
func onInit(plugin *glightning.Plugin, options map[string]string, config *glightning.Config) {
	log.Printf("versiion: %s initialized", VERSION)

	lightningdir = config.LightningDir
	options["rpc-file"] = fmt.Sprintf("%s/%s", config.LightningDir, config.RpcFile)

	lightning.StartUp(config.RpcFile, config.LightningDir)

	cfg, err := ListConfigs()
	if err != nil {
		log.Fatal(err)
	}

	switch cfg.Network {
	case "bitcoin":
		bitcoinNet = &chaincfg.MainNetParams
	case "regtest":
		bitcoinNet = &chaincfg.RegressionNetParams
	case "signet":
		panic("unsupported network")
	default:
		bitcoinNet = &chaincfg.TestNet3Params
	}
}

func registerOptions(p *glightning.Plugin) {

}

type DumpKeys struct {
	IncludePriv bool
}

type DumpKeysResult struct {
	Xpriv string `json:"xpriv,omitempty"`
	Xpub  string `json:"xpub"`
}

func (m *DumpKeys) Call() (jrpc2.Result, error) {
	f, err := os.Open(lightningdir + "/hsm_secret")
	if err != nil {
		return nil, err
	}
	hsm_secret := make([]byte, 32)
	_, err = f.Read(hsm_secret)
	if err != nil {
		return nil, err
	}

	salt := []byte{0x0}
	bip32_seed := hkdf.New(sha256.New, hsm_secret, salt, []byte("bip32 seed"))
	b := make([]byte, 32)
	bip32_seed.Read(b)

	key, err := hdkeychain.NewMaster(b, bitcoinNet)
	if err != nil {
		return nil, err
	}
	base1, _ := key.Child(0)
	master, _ := base1.Child(0)

	pub, _ := master.Neuter()
	dump := DumpKeysResult{
		Xpub: pub.String(),
	}
	if m.IncludePriv {
		dump.Xpriv = master.String()
	}
	return dump, nil
}

func (f *DumpKeys) Name() string {
	return "dump_keys"
}

func (f *DumpKeys) New() interface{} {
	return &DumpKeys{}
}

func registerMethods(p *glightning.Plugin) {
	dump := glightning.NewRpcMethod(&DumpKeys{}, `Dump extended keys`)
	dump.LongDesc = `optional parameter {include_priv} set to true if you want to include the private key, returns json object {"xpriv": "xpriv..."","xpub":"xpub..."}`
	dump.Usage = "include_priv"
	p.RegisterMethod(dump)

}
