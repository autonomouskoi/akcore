package main

import (
	"context"
	"fmt"
	"log"
	"os"

	extism "github.com/extism/go-sdk"
)

func main() {
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			/*
				extism.WasmUrl{
					Url: "https://github.com/extism/plugins/releases/latest/download/count_vowels.wasm",
				},
			*/
			extism.WasmFile{
				Path: "bonk.wasm",
			},
		},
	}

	ping := extism.NewHostFunctionWithStack("ping", func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
		pv, err := p.ReadString(stack[0])
		if err != nil {
			log.Print("error doing the do: ", err)
		}
		log.Print("Ping value: ", pv)
	},
		[]extism.ValueType{extism.ValueTypePTR},
		[]extism.ValueType{extism.ValueTypePTR},
	)

	ctx := context.Background()
	config := extism.PluginConfig{
		EnableWasi: true,
	}
	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{ping})

	if err != nil {
		fmt.Printf("Failed to initialize plugin: %v\n", err)
		os.Exit(1)
	}

	/*
		data := []byte("Hello, World!")
		exit, out, err := plugin.Call("count_vowels", data)
	*/
	exit, out, err := plugin.Call("greet", []byte("hork"))
	if err != nil {
		fmt.Println(err)
		os.Exit(int(exit))
	}

	response := string(out)
	fmt.Println(response)
}
