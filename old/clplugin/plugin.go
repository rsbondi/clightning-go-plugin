// helper package for creating clightning plugins
package clplugin

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Each method used by or plugin has the actual code to create a response.
// Also a description must be provided in the manifest, and is stored here.
// Required methods(init, getmanifest) are stored in this format but the
// description is not used.
type PluginMethod struct {
	Method      Rpcfun
	Description string
}

// This is the base struct for plugin creation.
// Methods and options used by the plugin are stored here and used to create
// the manifest.
// Additional members RpcFilename and LightningDir are saved during call to
// _init
// Several helper mehtods are attached for initializing and running a plugin.
type Plugin struct {
	Methods      map[string]PluginMethod
	Options      map[string]RpcInitOptions
	RpcFilename  string
	LightningDir string
	Init         InitCb
}

func NewPlugin() *Plugin {
	p := &Plugin{
		Methods: make(map[string]PluginMethod),
		Options: make(map[string]RpcInitOptions),
	}
	return p
}

// Add a method to the map of methods that will be called by your plugin.
// The methods are called by name when requested by the clightning daemon.
// The description is used in creating the response to `getmanifest`.
func (p *Plugin) AddMethod(name string, description string, method Rpcfun) {
	if _, exists := p.Methods[name]; exists {
		panic(fmt.Sprintf("attempted to add method %s but it already exists", name))
	}

	p.Methods[name] = PluginMethod{method, description}
}

// Add a option to the map of options that may be used by your plugin
// as additional command line arguments to the lightningd daemon.
// The map of options is used in creating the response to `getmanifest`.
func (p *Plugin) AddOption(name string, defaultVal string, description string) {
	if _, exists := p.Options[name]; exists {
		panic(fmt.Sprintf("attempted to add option %s but it already exists", name))
	}

	p.Options[name] = RpcInitOptions{
		Name:        name,
		Default:     defaultVal,
		Description: description,
		Type:        "string",
	}
}

// optionally used if additional initialization is needed
type InitCb func(json.RawMessage)

// Add additional initialization if needed.  Internally a call to `init` will
// set the lightning directory and rpc file, but here you can do more if needed
func (p *Plugin) AddInit(cb InitCb) {
	p.Init = cb
}

// This grabs common information that may optionally be used by your plugin.
// LightningDir and RpcFilename are saved separately, and later can be used to
// combined to make rpc calls to the daemon through a call to RpcFile().  If
// additional initialization is needed, a call back is provided that is set in
// AddInit(cb)
func (p *Plugin) _init(msg json.RawMessage) interface{} {
	var params RpcInitParams
	if err := json.Unmarshal(msg, &params); err != nil {
		panic(fmt.Sprintf("Plugin init failed: %s", err.Error()))
	}

	p.LightningDir = params.Configuration.LightningDir
	p.RpcFilename = params.Configuration.RpcFile

	if p.Init != nil {
		p.Init(msg)
	}
	return "ok"

}

// Returns a string of where to connect to the rpc interface of the daemon.
// This can be useful for getting information from your node via rpc.
func (p *Plugin) RpcFile() string {
	return fmt.Sprintf("%s/%s", p.LightningDir, p.RpcFilename)
}

// Helper type for detecting presense of string in slice
type strarr []string

// Detect presense of string in slice.  Returns true if needle found in
// haystack and false if not.
func (haystack strarr) Contains(needle string) bool {
	for _, n := range haystack {
		if needle == n {
			return true
		}
	}
	return false
}

// Automatically generate a response to `getmanifest`
func (p *Plugin) _getManifest(json.RawMessage) interface{} {
	var methods []RpcMethods
	skip := strarr{"init", "getmanifest"}

	for k, m := range p.Methods {
		if skip.Contains(k) {
			continue
		}
		met := RpcMethods{
			Name:        k,
			Description: m.Description,
		}

		methods = append(methods, met)
	}

	var options []RpcInitOptions
	for _, o := range p.Options {
		options = append(options, o)
	}

	return RpcInit{
		Rpcmethods: methods,
		Options:    options,
	}
}

// Function to pass log messages through to the daemon
// Possible levels, info, warn, debug
func (p *Plugin) Log(level string, msg string) {
	writer := bufio.NewWriter(os.Stdout)
	log := RpcLog{
		Level:   level,
		Message: msg,
	}
	rpc := LogCommand{
		Method:  "log",
		Jsonrpc: "2.0",
		Params:  log,
	}

	json.NewEncoder(writer).Encode(rpc)
	writer.Flush()
	writer.Reset(os.Stdout)

}

// Called from plugin instance to launch plugin. Additionally addes required
// `getmanifest` and `init` methods.  It will loop continually, monitoring
// stdin for commands and responding to requests.  The appropriate JSONRPC
// formatted response will be sent to stdout filling in the result with the
// response from the call to the mapped method.
func (p *Plugin) Run() {
	p.AddMethod("getmanifest", "", p._getManifest)
	p.AddMethod("init", "", p._init)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	decoder := json.NewDecoder(reader)
	encoder := json.NewEncoder(writer)
	for {
		var msg json.RawMessage
		cmd := RpcCommand{
			Params: &msg,
		}
		err := decoder.Decode(&cmd)

		if err != nil {
		}
		m, ok := p.Methods[cmd.Method]
		if ok {
			method := m.Method
			rpcResponse := RpcResult{
				Id:      cmd.Id,
				Jsonrpc: "2.0",
				Result:  method(msg),
			}

			encoder.Encode(rpcResponse)
			writer.Flush()
			writer.Reset(os.Stdout)
			reader.Reset(os.Stdin)
		}
		time.Sleep(50 * time.Millisecond) // TODO: Is this ok?  Maybe make parameter or let user handle in instance?
	}

}
