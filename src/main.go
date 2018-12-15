package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func jsonInit(msg json.RawMessage) interface{} {
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
		//log.Println(err)
	}
	var name string
	if len(s) > 0 {
		name = s[0]
	} else {
		name = "unkonwn user"
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
