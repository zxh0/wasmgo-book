package aot

import (
	"fmt"
	"wasm.go/binary"
)

func Compile(module binary.Module) {
	c := &moduleCompiler{
		printer:    newPrinter(),
		moduleInfo: newModuleInfo(module),
	}
	c.compile()
	fmt.Println(c.sb.String())
}
