package plugin

import (
	"encoding/json"
)

type RpcInitOptions struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

type RpcMethods struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RpcInit struct {
	Options    []RpcInitOptions `json:"options"`
	Rpcmethods []RpcMethods     `json:"rpcmethods"`
}

type RpcResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
}

type RpcCommand struct {
	Id      int
	Method  string
	Params  interface{}
	Jsonrpc string
}

type rpcfun func(json.RawMessage) interface{}

type RpcInitConfig struct {
	LightningDir string `json:"lightning-dir"`
	RpcFile      string `json:"rpc-file"`
}
