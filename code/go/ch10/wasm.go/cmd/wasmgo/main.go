package main

import (
	"flag"
	"fmt"
	"os"

	"wasm.go/binary"
	"wasm.go/instance"
	"wasm.go/interpreter"
)

func main() {
	dumpFlag := flag.Bool("d", false, "dump")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: wasmgo [-d] filename")
		os.Exit(1)
	}

	module, err := binary.DecodeFile(flag.Args()[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *dumpFlag {
		dump(module)
	} else {
		instantiateAndExecMainFunc(module)
	}
}

func instantiateAndExecMainFunc(module binary.Module) {
	mm := map[string]instance.Module{"env": newEnv()}
	m, err := interpreter.New(module, mm)
	if err == nil {
		_, err = m.InvokeFunc("main")
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
