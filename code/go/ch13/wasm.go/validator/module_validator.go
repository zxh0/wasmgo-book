package validator

import (
	"fmt"

	"wasm.go/binary"
)

type moduleValidator struct {
	module           binary.Module
	importedFuncs    []binary.Import
	importedTables   []binary.Import
	importedMemories []binary.Import
	importedGlobals  []binary.Import
	globalTypes      []binary.GlobalType
}

func Validate(module binary.Module) (err error) {
	defer func() {
		if _err := recover(); _err != nil {
			switch x := _err.(type) {
			case error:
				err = x
			default:
				panic(_err)
			}
		}
	}()
	v := &moduleValidator{module: module}
	v.validate()
	return
}

func (v *moduleValidator) validate() {
	v.validateImportSec()
	v.validateFuncSec()
	v.validateTableSec()
	v.validateMemSec()
	v.validateGlobalSec()
	v.validateExportSec()
	v.validateStartSec()
	v.validateElemSec()
	v.validateCodeSec()
	v.validateDataSec()
}

func (v *moduleValidator) validateImportSec() {
	for i, imp := range v.module.ImportSec {
		switch imp.Desc.Tag {
		case binary.ImportTagFunc:
			v.importedFuncs = append(v.importedFuncs, imp)
			if int(imp.Desc.FuncType) >= v.getTypeCount() {
				panic(fmt.Errorf("import[%d]: unknown type: %d",
					i, imp.Desc.FuncType))
			}
		case binary.ImportTagTable:
			if len(v.importedTables) > 0 {
				panic(fmt.Errorf("multiple tables"))
			}
			v.importedTables = append(v.importedTables, imp)
			if err := validateTableType(imp.Desc.Table.Limits); err != "" {
				panic(fmt.Errorf("import[%d]: %s", i, err))
			}
		case binary.ImportTagMem:
			if len(v.importedMemories) > 0 {
				panic(fmt.Errorf("multiple memories"))
			}
			v.importedMemories = append(v.importedMemories, imp)
			if err := validateMemoryType(imp.Desc.Mem); err != "" {
				panic(fmt.Errorf("import[%d]: %s", i, err))
			}
		case binary.ImportTagGlobal:
			v.importedGlobals = append(v.importedGlobals, imp)
			v.globalTypes = append(v.globalTypes, imp.Desc.Global)
		}
	}
}
func (v *moduleValidator) validateFuncSec() {
	for i, ftIdx := range v.module.FuncSec {
		if int(ftIdx) >= v.getTypeCount() {
			panic(fmt.Errorf("func[%d]: unknown type: %d", i, ftIdx))
		}
	}
}
func (v *moduleValidator) validateTableSec() {
	for i, table := range v.module.TableSec {
		if i+v.getImportedTableCount() > 0 {
			panic(fmt.Errorf("multiple tables"))
		}
		if err := validateTableType(table.Limits); err != "" {
			panic(fmt.Errorf("table[%d]: %s", i, err))
		}
	}
}
func (v *moduleValidator) validateMemSec() {
	for i, mem := range v.module.MemSec {
		if i+v.getImportedMemCount() > 0 {
			panic(fmt.Errorf("multiple memories"))
		}
		if err := validateMemoryType(mem); err != "" {
			panic(fmt.Errorf("mem[%d]: %s", i, err))
		}
	}
}
func (v *moduleValidator) validateGlobalSec() {
	for i, g := range v.module.GlobalSec {
		if err := v.validateConstExpr(g.Init, g.Type.ValType); err != "" {
			panic(fmt.Errorf("global[%d]: %s",
				i+v.getImportedGlobalCount(), err))
		}
		v.globalTypes = append(v.globalTypes, g.Type)
	}
}
func (v *moduleValidator) validateExportSec() {
	exportedNames := map[string]bool{}
	for i, exp := range v.module.ExportSec {
		if exportedNames[exp.Name] {
			panic(fmt.Errorf("duplicate export name: %s", exp.Name))
		} else {
			exportedNames[exp.Name] = true
		}

		switch exp.Desc.Tag {
		case binary.ExportTagFunc:
			if int(exp.Desc.Idx) >= v.getFuncCount() {
				panic(fmt.Errorf("export[%d]: unknown function: %d",
					i, exp.Desc.Idx))
			}
		case binary.ExportTagTable:
			if int(exp.Desc.Idx) >= v.getTableCount() {
				panic(fmt.Errorf("export[%d]: unknown table: %d",
					i, exp.Desc.Idx))
			}
		case binary.ExportTagMem:
			if int(exp.Desc.Idx) >= v.getMemCount() {
				panic(fmt.Errorf("export[%d]: unknown memory: %d",
					i, exp.Desc.Idx))
			}
		case binary.ExportTagGlobal:
			if int(exp.Desc.Idx) >= v.getGlobalCount() {
				panic(fmt.Errorf("export[%d]: unknown global: %d",
					i, exp.Desc.Idx))
			}
		}
	}
}
func (v *moduleValidator) validateStartSec() {
	if v.module.StartSec != nil {
		idx := int(*v.module.StartSec)
		ft, ok := v.getFuncType(idx)
		if !ok {
			panic(fmt.Errorf("start function: unknown function: %d", idx))
		}
		if len(ft.ParamTypes) > 0 || len(ft.ResultTypes) > 0 {
			panic(fmt.Errorf("start function: invalid type: %d", idx))
		}
	}
}
func (v *moduleValidator) validateElemSec() {
	for i, elem := range v.module.ElemSec {
		if int(elem.Table) >= v.getTableCount() {
			panic(fmt.Errorf("elem[%d]: unknown table: %d", i, elem.Table))
		}
		if err := v.validateConstExpr(elem.Offset, binary.ValTypeI32); err != "" {
			panic(fmt.Errorf("elem[%d]: %s", i, err))
		}
		for j, funcIdx := range elem.Init {
			if int(funcIdx) >= v.getFuncCount() {
				panic(fmt.Errorf("elem[%d][%d]: unknown function: %d", i, j, funcIdx))
			}
		}
	}
}
func (v *moduleValidator) validateCodeSec() {
	if len(v.module.CodeSec) != len(v.module.FuncSec) {
		panic(fmt.Errorf("invalid code count"))
	}
	for i, code := range v.module.CodeSec {
		ftIdx := v.module.FuncSec[i]
		ft := v.module.TypeSec[ftIdx]
		validateCode(v, i, code, ft)
	}
}
func (v *moduleValidator) validateDataSec() {
	for i, data := range v.module.DataSec {
		if int(data.Mem) >= v.getMemCount() {
			panic(fmt.Errorf("data[%d]: unknown memory: %d", i, data.Mem))
		}
		if err := v.validateConstExpr(data.Offset, binary.ValTypeI32); err != "" {
			panic(fmt.Errorf("data[%d]: %s", i, err))
		}
	}
}

// TODO
func (v *moduleValidator) validateConstExpr(expr []binary.Instruction,
	expectedType binary.ValType) (errMsg string) {

	if len(expr) > 1 {
		for _, instr := range expr {
			switch instr.Opcode {
			case binary.I32Const, binary.I64Const,
				binary.F32Const, binary.F64Const,
				binary.GlobalGet:
			default:
				return "constant expression required"
			}
		}
		return "type mismatch" // TODO
	}

	var actualType byte = 0
	if len(expr) > 0 {
		switch expr[0].Opcode {
		case binary.I32Const:
			actualType = binary.ValTypeI32
		case binary.I64Const:
			actualType = binary.ValTypeI64
		case binary.F32Const:
			actualType = binary.ValTypeF32
		case binary.F64Const:
			actualType = binary.ValTypeF64
		case binary.GlobalGet:
			gIdx := expr[0].Args.(uint32)
			if int(gIdx) >= len(v.globalTypes) {
				return fmt.Sprintf("unknown global: %d", gIdx)
			}
			actualType = v.globalTypes[gIdx].ValType
		default:
			return "constant expression required"
		}
	}
	if actualType != expectedType {
		return "type mismatch" // TODO
	}

	return ""
}

func (v *moduleValidator) getImportedFuncCount() int {
	return len(v.importedFuncs)
}
func (v *moduleValidator) getImportedTableCount() int {
	return len(v.importedTables)
}
func (v *moduleValidator) getImportedMemCount() int {
	return len(v.importedMemories)
}
func (v *moduleValidator) getImportedGlobalCount() int {
	return len(v.importedGlobals)
}

func (v *moduleValidator) getInternalFuncCount() int {
	return len(v.module.FuncSec)
}
func (v *moduleValidator) getInternalTableCount() int {
	return len(v.module.TableSec)
}
func (v *moduleValidator) getInternalMemCount() int {
	return len(v.module.MemSec)
}
func (v *moduleValidator) getInternalGlobalCount() int {
	return len(v.module.GlobalSec)
}

func (v *moduleValidator) getTypeCount() int {
	return len(v.module.TypeSec)
}
func (v *moduleValidator) getFuncCount() int {
	return v.getImportedFuncCount() + v.getInternalFuncCount()
}
func (v *moduleValidator) getTableCount() int {
	return v.getImportedTableCount() + v.getInternalTableCount()
}
func (v *moduleValidator) getMemCount() int {
	return v.getImportedMemCount() + v.getInternalMemCount()
}
func (v *moduleValidator) getGlobalCount() int {
	return v.getImportedGlobalCount() + v.getInternalGlobalCount()
}

func (v *moduleValidator) getFuncType(fIdx int) (binary.FuncType, bool) {
	if fIdx < v.getImportedFuncCount() {
		ftIdx := v.importedFuncs[fIdx].Desc.FuncType
		return v.module.TypeSec[ftIdx], true
	}
	if fIdx < v.getFuncCount() {
		ftIdx := v.module.FuncSec[fIdx-v.getImportedFuncCount()]
		return v.module.TypeSec[ftIdx], true
	}
	return binary.FuncType{}, false
}

func validateTableType(limits binary.Limits) string {
	return validateLimits(limits, (1<<32)-1, "table")
}
func validateMemoryType(limits binary.Limits) string {
	return validateLimits(limits, 1<<16, "mem")
}
func validateLimits(limits binary.Limits, k uint32, kind string) (errMsg string) {
	if limits.Min > k {
		if kind == "mem" {
			return "memory size must be at most 65536 pages (4GiB)"
		} else {
			// TODO
		}
	}
	if limits.Tag == 1 {
		if limits.Max > k {
			if kind == "mem" {
				return "memory size must be at most 65536 pages (4GiB)"
			} else {
				// TODO
			}
		}
		if limits.Max < limits.Min {
			return "size minimum must not be greater than maximum"
		}
	}
	return ""
}
