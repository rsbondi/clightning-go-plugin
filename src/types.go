package main

import "encoding/json"

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

func MakeRpcError(id string, code int, message string) string {
	errResponse := RpcError{
		Id:      id,
		Jsonrpc: "2.0",
		Error: RpcErrorObj{
			Code:    code,
			Message: message,
		},
	}
	err, _ := json.Marshal(errResponse)
	return string(err[:])

}

type RpcErrorObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RpcError struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Error   RpcErrorObj `json:"error"`
}

type rpcfun func(json.RawMessage) interface{}
