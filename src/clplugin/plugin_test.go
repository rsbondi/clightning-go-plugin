package clplugin

import (
	"encoding/json"
	"testing"
)

var initmsg = `{ "configuration": {"lightning-dir": "/home/richard/.lightning","rpc-file":"lightning-rpc"}}`

func TestPluginCreate(t *testing.T) {
	p := NewPlugin()
	if p == nil {
		t.Errorf("Plugin creation failed")
	}

}

func TestPluginInit(t *testing.T) {
	p := NewPlugin()
	var msg json.RawMessage = json.RawMessage(initmsg)
	p._init(msg)
	want := "/home/richard/.lightning"
	if p.LightningDir != want {
		t.Errorf("Plugin initialization failed, expected \"%s\", got \"%s\"", want, p.LightningDir)
	}
	want = "lightning-rpc"
	if p.RpcFilename != want {
		t.Errorf("Plugin initialization failed, expected \"%s\", got \"%s\"", want, p.LightningDir)
	}
}

func TestRpcFile(t *testing.T) {
	p := NewPlugin()
	var msg json.RawMessage = json.RawMessage(initmsg)
	p._init(msg)
	want := "/home/richard/.lightning/lightning-rpc"
	have := p.RpcFile()
	if have != want {
		t.Errorf("Plugin initialization failed, expected \"%s\", got \"%s\"", want, have)
	}
}

func TestManifest(t *testing.T) {
	p := NewPlugin()
	p.AddMethod("fundprice", "show fund summary with price", func(msg json.RawMessage) interface{} { return "" })

	p.AddOption("fiat", "USD", "Ticker symbol for fiat currency.")
	p.AddOption("crypto", "BTC", "Ticker symbol for crypto currency.")
	var msg json.RawMessage = json.RawMessage("{}")

	init := p._getManifest(msg).(RpcInit)

	want := "show fund summary with price"
	have := init.Rpcmethods[0].Description
	if have != want {
		t.Errorf("Plugin manifest failed, expected \"%s\", got \"%s\"", want, have)
	}

	var fiat RpcInitOptions
	var crypto RpcInitOptions

	for _, opt := range init.Options {
		if opt.Name == "fiat" {
			fiat = opt
		}
		if opt.Name == "crypto" {
			crypto = opt
		}
	}

	want = "USD"
	have = fiat.Default
	if have != want {
		t.Errorf("Plugin manifest failed, expected \"%s\", got \"%s\"", want, have)
	}

	want = "BTC"
	have = crypto.Default
	if have != want {
		t.Errorf("Plugin manifest failed, expected \"%s\", got \"%s\"", want, have)
	}

}
