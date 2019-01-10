package plugin

import (
	"encoding/json"
	"fmt"
)

type PluginMethod struct {
	Method      rpcfun
	Description string
}

type Plugin struct {
	Methods map[string]PluginMethod
	Options map[string]RpcInitOptions
}

func (p *Plugin) AddMethod(name string, description string, method rpcfun) {
	if _, exists := p.Methods[name]; exists {
		panic(fmt.Sprintf("attempted to add method %s but it already exists"))
	}

	p.Methods[name] = PluginMethod{method, description}
}

func (p *Plugin) AddOption(name string, defaultVal string, description string) {
	if _, exists := p.Options[name]; exists {
		panic(fmt.Sprintf("attempted to add option %s but it already exists"))
	}

	p.Options[name] = RpcInitOptions{
		Name:        name,
		Default:     defaultVal,
		Description: description,
		Type:        "string",
	}
}

func (p *Plugin) _init() {

}

type strarr []string

func (haystack strarr) Contains(needle string) bool {
	for _, n := range haystack {
		if needle == n {
			return true
		}
	}
	return false
}

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

	return RpcInit{
		Rpcmethods: methods,
	}
}

func (p *Plugin) Run() {
	p.AddMethod("getmanifest", "", p._getManifest)
}
