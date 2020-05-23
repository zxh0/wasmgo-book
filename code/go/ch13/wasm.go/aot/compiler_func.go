package aot

type funcCompiler struct {
	printer
}

func newFuncCompiler() funcCompiler {
	return funcCompiler{newPrinter()}
}

func (c *funcCompiler) genParams(paramCount int) {
	for i := 0; i < paramCount; i++ {
		c.printf("a%d", i)
		c.printIf(i < paramCount-1, ", ", " uint64")
	}
}

func (c *funcCompiler) genResults(resultCount int) {
	if resultCount > 0 {
		c.printIf(resultCount > 1, " (", " ")
		for i := 0; i < resultCount; i++ {
			c.printIf(i > 0, ", ", "")
			c.printf("uint64")
		}
		c.printIf(resultCount > 1, ")", "")
	}
}
