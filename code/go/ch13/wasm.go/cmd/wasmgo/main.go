package main

import (
	"flag"
	"fmt"
	"os"

	"wasm.go/binary"
	"wasm.go/instance"
	"wasm.go/interpreter"
	"wasm.go/validator"
)

func main() {
	dumpFlag := flag.Bool("d", false, "dump Wasm file")
	checkFlag := flag.Bool("c", false, "check Wasm file")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf(`Usage: 
	wasmgo    filename
	wasmgo -d filename
	wasmgo -c filename
`)
		os.Exit(1)
	}

	module, err := binary.DecodeFile(flag.Args()[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *dumpFlag {
		dump(module)
	} else if *checkFlag {
		check(module)
	} else {
		instantiateAndExecMainFunc(module)
	}
}

func check(module binary.Module) {
	if err := validator.Validate(module); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("OK!")
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
