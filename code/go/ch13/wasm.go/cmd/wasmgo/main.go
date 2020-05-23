package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"wasm.go/aot"
	"wasm.go/binary"
	"wasm.go/instance"
	"wasm.go/interpreter"
	"wasm.go/validator"
)

func main() {
	dumpFlag := flag.Bool("d", false, "dump Wasm file")
	checkFlag := flag.Bool("c", false, "check Wasm file")
	aotFlag := flag.Bool("a", false, "compile Wasm file to Go plugin")

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf(`Usage: 
	wasmgo    filename
	wasmgo -d filename
	wasmgo -c filename
	wasmgo -a filename
`)
		os.Exit(1)
	}

	filename := flag.Args()[0]
	if *dumpFlag {
		dump(decode(filename))
	} else if *checkFlag {
		check(decode(filename))
	} else if *aotFlag {
		aot.Compile(decode(filename))
	} else if strings.HasSuffix(filename, ".so") {
		execSO(filename)
	} else {
		instantiateAndExecMainFunc(decode(filename))
	}
}

func decode(filename string) binary.Module {
	module, err := binary.DecodeFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return module
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

func execSO(filename string) {
	mm := map[string]instance.Module{
		"env": newEnv(),
	}
	i, err := aot.Load(filename, mm)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	_, err = i.InvokeFunc("main")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
