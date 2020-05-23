package aot

import "wasm.go/binary"

type externalFuncCompiler struct {
	funcCompiler
}

func newExternalFuncCompiler() *externalFuncCompiler {
	return &externalFuncCompiler{newFuncCompiler()}
}

func (c *externalFuncCompiler) compile(idx int, ft binary.FuncType) string {
	c.printf("func (m *aotModule) f%d(", idx)
	c.genParams(len(ft.ParamTypes))
	c.print(")")
	c.genResults(len(ft.ResultTypes))
	c.print(" {\n")
	c.genFuncBody(idx, ft)
	c.println("}")
	return c.sb.String()
}

func (c *externalFuncCompiler) genFuncBody(idx int, ft binary.FuncType) {
	c.printIf(len(ft.ResultTypes) > 0,
		"\tresults, err := ",
		"\t_, err := ")
	c.printf("m.importedFuncs[%d].Call(", idx)
	for i, vt := range ft.ParamTypes {
		c.printIf(i > 0, ", ", "")
		switch vt {
		case binary.ValTypeI32:
			c.printf("int32(a%d)", i)
		case binary.ValTypeI64:
			c.printf("int64(a%d)", i)
		case binary.ValTypeF32:
			c.printf("_f32(a%d)", i)
		case binary.ValTypeF64:
			c.printf("_f64(a%d)", i)
		}
	}
	c.println(")")
	c.println("\tif err != nil {} // TODO")
	if len(ft.ResultTypes) > 0 {
		c.print("\treturn ")
		for i, vt := range ft.ResultTypes {
			c.printIf(i > 0, ", ", "")
			switch vt {
			case binary.ValTypeI32:
				c.printf("uint64(results[%d].(int32))", i)
			case binary.ValTypeI64:
				c.printf("uint64(results[%d].(int64))", i)
			case binary.ValTypeF32:
				c.printf("_u32(results[%d].(float32))", i)
			case binary.ValTypeF64:
				c.printf("_u64(results[%d].(float64))", i)
			}
		}
	}
}
