package aot

import "wasm.go/binary"

type exportedFuncCompiler struct {
	printer
	importedFuncCount int
}

func newExportedFuncCompiler(importedFuncCount int) *exportedFuncCompiler {
	return &exportedFuncCompiler{
		printer:           newPrinter(),
		importedFuncCount: importedFuncCount,
	}
}

func (c *exportedFuncCompiler) compile(expIdx, fIdx int, ft binary.FuncType) string {
	c.printf("func (m *aotModule) exported%d(args []interface{}) ([]interface{}, error) {\n", expIdx)
	if fIdx < c.importedFuncCount {
		c.printf("\treturn m.f%d(args...)\n", fIdx)
	} else {
		c.print("\t")
		c.genResults(len(ft.ResultTypes))
		c.printf("m.f%d(", fIdx)
		c.genParams(ft)
		c.println(")")
		c.genReturn(ft)
	}
	c.println("}")
	return c.sb.String()
}

// r0, r1, ... := f()
func (c *exportedFuncCompiler) genResults(resultCount int) {
	if resultCount > 0 {
		for i := 0; i < resultCount; i++ {
			c.printIf(i > 0, ", ", "")
			c.printf("r%d", i)
		}
		c.print(" := ")
	}
}

func (c *exportedFuncCompiler) genParams(ft binary.FuncType) {
	for i, vt := range ft.ParamTypes {
		c.printIf(i > 0, ", ", "")
		switch vt {
		case binary.ValTypeI32:
			c.printf("uint64(args[%d].(int32))", i)
		case binary.ValTypeI64:
			c.printf("uint64(args[%d].(int64))", i)
		case binary.ValTypeF32:
			c.printf("_u32(args[%d].(float32))", i)
		case binary.ValTypeF64:
			c.printf("_u64(args[%d].(float64))", i)
		}
	}
}

func (c *exportedFuncCompiler) genReturn(ft binary.FuncType) {
	if len(ft.ResultTypes) == 0 {
		c.println("\treturn nil, nil")
	} else {
		c.print("\treturn []interface{}{")
		for i, vt := range ft.ResultTypes {
			c.printIf(i > 0, ", ", "")
			switch vt {
			case binary.ValTypeI32:
				c.printf("int32(r%d)", i)
			case binary.ValTypeI64:
				c.printf("int64(r%d)", i)
			case binary.ValTypeF32:
				c.printf("_f32(r%d)", i)
			case binary.ValTypeF64:
				c.printf("_f64(r%d)", i)
			}
		}
		c.println("}, nil")
	}
}
