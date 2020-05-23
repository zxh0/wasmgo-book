package main

import (
	"fmt"

	"wasm.go/binary"
)

type dumper struct {
	module              binary.Module
	importedFuncCount   int
	importedTableCount  int
	importedMemCount    int
	importedGlobalCount int
}

func dump(module binary.Module) {
	d := &dumper{module: module}

	fmt.Printf("Version: 0x%02x\n", d.module.Version)
	d.dumpTypeSec()
	d.dumpImportSec()
	d.dumpFuncSec()
	d.dumpTableSec()
	d.dumpMemSec()
	d.dumpGlobalSec()
	d.dumpExportSec()
	d.dumpStartSec()
	d.dumpElemSec()
	d.dumpCodeSec()
	d.dumpDataSec()
	d.dumpCustomSec()
}

func (d *dumper) dumpTypeSec() {
	fmt.Printf("Type[%d]:\n", len(d.module.TypeSec))
	for i, ft := range d.module.TypeSec {
		fmt.Printf("  type[%d]: %s\n", i, ft)
	}
}

func (d *dumper) dumpImportSec() {
	fmt.Printf("Import[%d]:\n", len(d.module.ImportSec))
	for _, imp := range d.module.ImportSec {
		switch imp.Desc.Tag {
		case binary.ImportTagFunc:
			fmt.Printf("  func[%d]: %s.%s, sig=%d\n",
				d.importedFuncCount, imp.Module, imp.Name, imp.Desc.FuncType)
			d.importedFuncCount++
		case binary.ImportTagTable:
			fmt.Printf("  table[%d]: %s.%s, %s\n",
				d.importedTableCount, imp.Module, imp.Name, imp.Desc.Table.Limits)
			d.importedTableCount++
		case binary.ImportTagMem:
			fmt.Printf("  memory[%d]: %s.%s, %s\n",
				d.importedMemCount, imp.Module, imp.Name, imp.Desc.Mem)
			d.importedMemCount++
		case binary.ImportTagGlobal:
			fmt.Printf("  global[%d]: %s.%s, %s\n",
				d.importedGlobalCount, imp.Module, imp.Name, imp.Desc.Global)
			d.importedGlobalCount++
		}
	}
	return
}

func (d *dumper) dumpFuncSec() {
	fmt.Printf("Function[%d]:\n", len(d.module.FuncSec))
	for i, sig := range d.module.FuncSec {
		fmt.Printf("  func[%d]: sig=%d\n",
			d.importedFuncCount+i, sig)
	}
}

func (d *dumper) dumpTableSec() {
	fmt.Printf("Table[%d]:\n", len(d.module.TableSec))
	for i, t := range d.module.TableSec {
		fmt.Printf("  table[%d]: %s\n",
			d.importedTableCount+i, t.Limits)
	}
}

func (d *dumper) dumpMemSec() {
	fmt.Printf("Memory[%d]:\n", len(d.module.MemSec))
	for i, limits := range d.module.MemSec {
		fmt.Printf("  memory[%d]: %s\n",
			d.importedMemCount+i, limits)
	}
}

func (d *dumper) dumpGlobalSec() {
	fmt.Printf("Global[%d]:\n", len(d.module.GlobalSec))
	for i, g := range d.module.GlobalSec {
		fmt.Printf("  global[%d]: %s\n",
			d.importedGlobalCount+i, g.Type)
	}
}

func (d *dumper) dumpExportSec() {
	fmt.Printf("Export[%d]:\n", len(d.module.ExportSec))
	for _, exp := range d.module.ExportSec {
		switch exp.Desc.Tag {
		case binary.ExportTagFunc:
			fmt.Printf("  func[%d]: name=%s\n", int(exp.Desc.Idx), exp.Name)
		case binary.ExportTagTable:
			fmt.Printf("  table[%d]: name=%s\n", int(exp.Desc.Idx), exp.Name)
		case binary.ExportTagMem:
			fmt.Printf("  memory[%d]: name=%s\n", int(exp.Desc.Idx), exp.Name)
		case binary.ExportTagGlobal:
			fmt.Printf("  global[%d]: name=%s\n", int(exp.Desc.Idx), exp.Name)
		}
	}
}

func (d *dumper) dumpStartSec() {
	fmt.Printf("Start:\n")
	if d.module.StartSec != nil {
		fmt.Printf("  func=%d\n", *d.module.StartSec)
	}
}

func (d *dumper) dumpElemSec() {
	fmt.Printf("Element[%d]:\n", len(d.module.ElemSec))
	for i, elem := range d.module.ElemSec {
		fmt.Printf("  elem[%d]: table=%d\n", i, elem.Table) // TODO
	}
}

func (d *dumper) dumpCodeSec() {
	fmt.Printf("Code[%d]:\n", len(d.module.CodeSec))
	for i, code := range d.module.CodeSec {
		fmt.Printf("  func[%d]: locals=[", d.importedFuncCount+i) // TODO
		if len(code.Locals) > 0 {
			for i, locals := range code.Locals {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s x %d",
					binary.ValTypeToStr(locals.Type), locals.N)
			}
		}
		fmt.Println("]")
	}
}

func (d *dumper) dumpDataSec() {
	fmt.Printf("Data[%d]:\n", len(d.module.DataSec))
	for i, data := range d.module.DataSec {
		fmt.Printf("  data[%d]: mem=%d\n", i, data.Mem) // TODO
	}
}

func (d *dumper) dumpCustomSec() {
	fmt.Printf("Custom[%d]:\n", len(d.module.CustomSecs))
	for i, cs := range d.module.CustomSecs {
		fmt.Printf("  custom[%d]: name=%s\n", i, cs.Name) // TODO
	}
}
