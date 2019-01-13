package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

var alias string // this example uses node's alias if no name provided

func jsonInit(msg json.RawMessage) interface{} {
	var params RpcInitParams
	if err := json.Unmarshal(msg, &params); err != nil {

	}

	// example to make rpc calls to lightningd, connect socket from the following info provided in init call
	// not the best code but illustrative
	c, err := net.Dial("unix", fmt.Sprintf("%s/%s", params.Configuration.LightningDir, params.Configuration.RpcFile))
	if err != nil {
	}

	c.Write([]byte(`{"jsonrpc":"2.0","id":"plugininit","method":"getinfo","params":[]}`))
	buf := make([]byte, 1024)
	n, err := c.Read(buf[:])
	if err != nil {
	}

	var info RpcInfoResult
	rpcinf := RpcInfo{
		Result: info,
	}
	if err := json.Unmarshal(buf[0:n], &rpcinf); err != nil {
	}

	alias = rpcinf.Result.Alias // or to use lightningd --greeting params.Options.Greeting

	return "ok"
}

func jsonGetManifest(msg json.RawMessage) interface{} {
	response := RpcInit{
		Options: []RpcInitOptions{
			{
				Name:        "greeting",
				Type:        "string",
				Default:     "World",
				Description: "What name should I call you?",
			},
		},
		Rpcmethods: []RpcMethods{
			{
				Name:        "hello",
				Description: "Returns a personalized greeting for {name}",
			},
		},
	}
	return response
}

func jsonHello(msg json.RawMessage) interface{} {
	var s []string
	if err := json.Unmarshal(msg, &s); err != nil {

	}

	var name string
	if len(s) > 0 {
		name = s[0]
	} else {
		name = alias
	}

	return fmt.Sprintf("Greetings from plugin %s", name)
}

func main() {
	commands := map[string]rpcfun{
		"init":        jsonInit,
		"hello":       jsonHello,
		"getmanifest": jsonGetManifest,
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	for {
		var msg json.RawMessage
		cmd := RpcCommand{
			Params: &msg,
		}
		err := json.NewDecoder(reader).Decode(&cmd)

		if err != nil {
		}
		method, ok := commands[cmd.Method]
		if ok {
			rpcResponse := RpcResult{
				Id:      cmd.Id,
				Jsonrpc: "2.0",
				Result:  method(msg),
			}

			json.NewEncoder(writer).Encode(rpcResponse)
			writer.Flush()
			writer.Reset(os.Stdout)
			reader.Reset(os.Stdin)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
