package clplugin

import (
	"encoding/json"
)

// The options for generating the manifest are stored here
type RpcInitOptions struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

// The methods for generating the manifest are dervied from this struct
type RpcMethods struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// This is the complete init message
type RpcInit struct {
	Options    []RpcInitOptions `json:"options"`
	Rpcmethods []RpcMethods     `json:"rpcmethods"`
}

// When a response to an rpc call is generated, it will use this format.  The
// `Result` is a generic interace and can represent any object/format
type RpcResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
}

// When recieving command via rpc from the daemon, it will be in this format.
// It may also optionally be used by your plugin if it needs to talk to the
// daemon via rpc.
type RpcCommand struct {
	Id      int
	Method  string
	Params  interface{}
	Jsonrpc string
}

// All methods created by this package for your plugin as well as any method
// created by your plugin should use this
type Rpcfun func(json.RawMessage) interface{}

// This is the container for the values sent from the daemon which may be used
// by your plugin. These are automatically stored by the `Plugin._init` method.
// Their use is not required but handy if talking to the daemon via rpc.
type RpcInitConfig struct {
	LightningDir string `json:"lightning-dir"`
	RpcFile      string `json:"rpc-file"`
}

// The configuration send during `init` is stored here.
type RpcInitParams struct {
	Configuration RpcInitConfig `json:"configuration"`
}

// Plugin log parameters
type RpcLog struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

// This is how to pass log info to the daemon
type LogCommand struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  RpcLog `json:"params"`
}
