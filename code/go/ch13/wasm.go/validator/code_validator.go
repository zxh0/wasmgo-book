package validator

import (
	"fmt"

	"wasm.go/binary"
)

// https://webassembly.github.io/spec/core/appendix/algorithm.html

/*
type val_type = I32 | I64 | F32 | F64
type opd_stack = stack(val_type | Unknown)
type ctrl_stack = stack(ctrl_frame)
*/

const (
	Unknown = 0

	I32 = binary.ValTypeI32
	I64 = binary.ValTypeI64
	F32 = binary.ValTypeF32
	F64 = binary.ValTypeF64
)

type valType = byte
type opdStack = []valType
type ctrlStack = []ctrlFrame

/*
type ctrl_frame = {
  opcode      : opcode
  start_types : list(val_type)
  end_types   : list(val_type)
  height      : nat
  unreachable : bool
}
*/
type ctrlFrame struct {
	opcode      byte
	startTypes  []valType
	endTypes    []valType
	height      int
	unreachable bool
}

/*
var opds : opd_stack
var ctrls : ctrl_stack
*/
type codeValidator struct {
	opds       opdStack
	ctrls      ctrlStack
	mv         *moduleValidator
	codeIdx    int
	localCount int
	instrPath  map[int]string // depth -> opname
}

func validateCode(mv *moduleValidator, idx int,
	code binary.Code, ft binary.FuncType) {

	cv := &codeValidator{
		mv:        mv,
		codeIdx:   idx,
		instrPath: map[int]string{},
	}
	cv.validateCode(code, ft)
}

func (cv *codeValidator) error(msg string) {
	panic(fmt.Errorf("code[%d], %s: %s",
		cv.codeIdx, cv.getInstrPath(), msg))
}
func (cv *codeValidator) errorf(format string, a ...interface{}) {
	cv.error(fmt.Sprintf(format, a...))
}

func (cv *codeValidator) getInstrPath() string {
	path := ""
	for i := 0; i < len(cv.instrPath); i++ {
		if i > 0 {
			path += "/"
		}
		path += cv.instrPath[i]
	}
	return path
}

/* ctrl_frame */

/*
func label_types(frame : ctrl_frame) : list(val_types) =
  return (
    if frame.opcode == loop then
      frame.start_types
    else
      frame.end_types
  )
*/
func (frame ctrlFrame) labelTypes() []valType {
	if frame.opcode == binary.Loop {
		return frame.startTypes
	}
	return frame.endTypes
}

/* operand stack */

/*
func push_opd(type : val_type | Unknown) =
  opds.push(type)
*/
func (cv *codeValidator) pushOpd(vt valType) {
	cv.opds = append(cv.opds, vt)
}

/*
func pop_opd() : val_type | Unknown =
  if (opds.size() = ctrls[0].height && ctrls[0].unreachable) return Unknown
  error_if(opds.size() = ctrls[0].height)
  return opds.pop()
*/
func (cv *codeValidator) popOpd() valType {
	if ctrl0 := cv.getCtrl(0); len(cv.opds) == ctrl0.height {
		if ctrl0.unreachable {
			return Unknown
		}
		cv.error("type mismatch") // TODO
	}
	r := cv.opds[len(cv.opds)-1]
	cv.opds = cv.opds[:len(cv.opds)-1]
	return r
}

/*
func pop_opd(expect : val_type | Unknown) : val_type | Unknown =
  let actual = pop_opd()
  if (actual = Unknown) return expect
  if (expect = Unknown) return actual
  error_if(actual =/= expect)
  return actual
*/
func (cv *codeValidator) popOpdOf(expect valType) valType {
	actual := cv.popOpd()
	if actual == Unknown {
		return expect
	}
	if expect == Unknown {
		return actual
	}
	if actual != expect {
		cv.error("type mismatch") // TODO
	}
	return actual
}

/*
func push_opds(types : list(val_type)) = foreach (t in types) push_opd(t)
func pop_opds(types : list(val_type)) = foreach (t in reverse(types)) pop_opd(t)
*/
func (cv *codeValidator) pushOpds(types []valType) {
	for _, t := range types {
		cv.pushOpd(t)
	}
}
func (cv *codeValidator) popOpds(types []valType) {
	for i := len(types) - 1; i >= 0; i-- {
		cv.popOpdOf(types[i])
	}
}

func (cv *codeValidator) pushI32() { cv.pushOpd(I32) }
func (cv *codeValidator) pushI64() { cv.pushOpd(I64) }
func (cv *codeValidator) pushF32() { cv.pushOpd(F32) }
func (cv *codeValidator) pushF64() { cv.pushOpd(F64) }

func (cv *codeValidator) popI32() { cv.popOpdOf(I32) }
func (cv *codeValidator) popI64() { cv.popOpdOf(I64) }
func (cv *codeValidator) popF32() { cv.popOpdOf(F32) }
func (cv *codeValidator) popF64() { cv.popOpdOf(F64) }

/* control stack */

func (cv *codeValidator) getCtrl(n int) ctrlFrame {
	if n >= len(cv.ctrls) {
		cv.error("")
	}
	return cv.ctrls[len(cv.ctrls)-1-n]
}

/*
func push_ctrl(opcode : opcode, in : list(val_type), out : list(val_type)) =
  let frame = ctrl_frame(opcode, in, out, opds.size(), false)
  ctrls.push(frame)
  push_opds(in)
*/
func (cv *codeValidator) pushCtrl(opcode byte, in, out []valType) {
	frame := ctrlFrame{opcode, in, out, len(cv.opds), false}
	cv.ctrls = append(cv.ctrls, frame)
	cv.pushOpds(in)
}

/*
func pop_ctrl() : ctrl_frame =
  error_if(ctrls.is_empty())
  let frame = ctrls[0]
  pop_opds(frame.end_types)
  error_if(opds.size() =/= frame.height)
  ctrls.pop()
  return frame
*/
func (cv *codeValidator) popCtrl() ctrlFrame {
	if len(cv.ctrls) == 0 {
		cv.error("")
	}
	frame := cv.getCtrl(0)
	cv.popOpds(frame.endTypes)
	if len(cv.opds) != frame.height {
		cv.error("type mismatch") // TODO
	}
	cv.ctrls = cv.ctrls[:len(cv.ctrls)-1]
	return frame
}

/*
func unreachable() =
  opds.resize(ctrls[0].height)
  ctrls[0].unreachable := true
*/
func (cv *codeValidator) unreachable() {
	cv.opds = cv.opds[:cv.getCtrl(0).height]
	//cv.getCtrl(0).unreachable = true
	cv.ctrls[len(cv.ctrls)-1].unreachable = true
}

/* code validation */

func (cv *codeValidator) validateCode(
	code binary.Code, ft binary.FuncType) {

	cv.pushOpds(ft.ParamTypes)
	cv.localCount = len(ft.ParamTypes)
	for _, local := range code.Locals {
		for i := 0; i < int(local.N); i++ {
			cv.pushOpd(local.Type)
			cv.localCount++
		}
	}
	cv.pushCtrl(binary.Block, nil, ft.ResultTypes)
	cv.validateExpr(code.Expr)
	cv.pushOpds(cv.popCtrl().endTypes)
}

func (cv *codeValidator) validateExpr(expr []binary.Instruction) {
	depth := len(cv.instrPath)
	for _, instr := range expr {
		cv.instrPath[depth] = instr.GetOpname()
		cv.validateInstr(instr)
	}
	delete(cv.instrPath, depth)
}

/*
func validate(opcode) = switch (opcode)
  case (i32.add)
    pop_opd(I32)
    pop_opd(I32)
    push_opd(I32)
  case (drop)
    pop_opd()
  case (select)
    pop_opd(I32)
    let t1 = pop_opd()
    let t2 = pop_opd(t1)
    push_opd(t2)
  case (unreachable)
    unreachable()
  case (block t1*->t2*)
    pop_opds([t1*])
    push_ctrl(block, [t1*], [t2*])
  case (loop t1*->t2*)
    pop_opds([t1*])
    push_ctrl(loop, [t1*], [t2*])
  case (if t1*->t2*)
    pop_opd(I32)
    pop_opds([t1*])
    push_ctrl(if, [t1*], [t2*])
  case (end)
    let frame = pop_ctrl()
    push_opds(frame.end_types)
  case (else)
    let frame = pop_ctrl()
    error_if(frame.opcode =/= if)
    push_ctrl(else, frame.start_types, frame.end_types)
  case (br n)
    error_if(ctrls.size() < n)
    pop_opds(label_types(ctrls[n]))
    unreachable()
  case (br_if n)
    error_if(ctrls.size() < n)
    pop_opd(I32)
    pop_opds(label_types(ctrls[n]))
    push_opds(label_types(ctrls[n]))
  case (br_table n* m)
    error_if(ctrls.size() < m)
    foreach (n in n*)
      error_if(ctrls.size() < n || label_types(ctrls[n]) =/= label_types(ctrls[m]))
    pop_opd(I32)
    pop_opds(label_types(ctrls[m]))
    unreachable()
*/
func (cv *codeValidator) validateInstr(instr binary.Instruction) {
	switch instr.Opcode {
	case binary.Unreachable:
		cv.unreachable()
	case binary.Nop:
	case binary.Block, binary.Loop:
		blockArgs := instr.Args.(binary.BlockArgs)
		bt := cv.mv.module.GetBlockType(blockArgs.BT)
		cv.popOpds(bt.ParamTypes)
		cv.pushCtrl(instr.Opcode, bt.ParamTypes, bt.ResultTypes)
		cv.validateExpr(blockArgs.Instrs)
		cv.pushOpds(cv.popCtrl().endTypes)
	case binary.If:
		ifArgs := instr.Args.(binary.IfArgs)
		bt := cv.mv.module.GetBlockType(ifArgs.BT)
		cv.popI32()
		cv.popOpds(bt.ParamTypes)
		cv.pushCtrl(binary.If, bt.ParamTypes, bt.ResultTypes)
		cv.validateExpr(ifArgs.Instrs1)
		// else
		frame := cv.popCtrl()
		if frame.opcode != binary.If {
			cv.error("TODO")
		}
		cv.pushCtrl(binary.Else_, frame.startTypes, frame.endTypes)
		cv.validateExpr(ifArgs.Instrs2)
		// end
		cv.pushOpds(cv.popCtrl().endTypes)
	case binary.Br:
		n := int(instr.Args.(uint32))
		if len(cv.ctrls) < n {
			cv.error("unknown label")
		}
		cv.popOpds(cv.getCtrl(n).labelTypes())
		cv.unreachable()
	case binary.BrIf:
		n := int(instr.Args.(uint32))
		if len(cv.ctrls) < n {
			cv.error("unknown label")
		}
		cv.popI32()
		cv.popOpds(cv.getCtrl(n).labelTypes())
		cv.pushOpds(cv.getCtrl(n).labelTypes())
	case binary.BrTable:
		brTableArgs := instr.Args.(binary.BrTableArgs)
		m := int(brTableArgs.Default)
		if len(cv.ctrls) < m {
			cv.error("unknown label")
		}
		for _, n := range brTableArgs.Labels {
			if len(cv.ctrls) < int(n) {
				cv.error("unknown label")
			}
			t1 := cv.getCtrl(int(n)).labelTypes()
			t2 := cv.getCtrl(m).labelTypes()
			if !isValTypesEq(t1, t2) {
				cv.error("type mismatch")
			}
		}
		cv.popI32()
		cv.popOpds(cv.getCtrl(m).labelTypes())
		cv.unreachable()
	case binary.Return:
		n := len(cv.ctrls) - 1
		cv.popOpds(cv.getCtrl(n).labelTypes())
		cv.unreachable()
	case binary.Call:
		fIdx := instr.Args.(uint32)
		ft, ok := cv.mv.getFuncType(int(fIdx))
		if !ok {
			cv.error("unknown function")
		}
		cv.popOpds(ft.ParamTypes)
		cv.pushOpds(ft.ResultTypes)
	case binary.CallIndirect:
		if cv.mv.getTableCount() == 0 {
			cv.error("unknown table")
		}
		ftIdx := instr.Args.(uint32)
		if int(ftIdx) >= cv.mv.getTypeCount() {
			cv.error("unknown type")
		}
		ft := cv.mv.module.TypeSec[ftIdx]
		cv.popI32()
		cv.popOpds(ft.ParamTypes)
		cv.pushOpds(ft.ResultTypes)
	case binary.Drop:
		cv.popOpd()
	case binary.Select:
		cv.popI32()
		t1 := cv.popOpd()
		t2 := cv.popOpdOf(t1)
		cv.pushOpd(t2)
	case binary.LocalGet:
		n := int(instr.Args.(uint32))
		if n >= cv.localCount {
			cv.errorf("unknown local: %d", n)
		}
		cv.pushOpd(cv.opds[n])
	case binary.LocalSet:
		n := int(instr.Args.(uint32))
		if n >= cv.localCount {
			cv.errorf("unknown local: %d", n)
		}
		cv.popOpdOf(cv.opds[n])
	case binary.LocalTee:
		n := int(instr.Args.(uint32))
		if n >= cv.localCount {
			cv.errorf("unknown local: %d", n)
		}
		cv.popOpdOf(cv.opds[n])
		cv.pushOpd(cv.opds[n])
	case binary.GlobalGet:
		n := int(instr.Args.(uint32))
		if n >= len(cv.mv.globalTypes) {
			cv.errorf("unknown global: %d", n)
		}
		cv.pushOpd(cv.mv.globalTypes[n].ValType)
	case binary.GlobalSet:
		n := int(instr.Args.(uint32))
		if n >= len(cv.mv.globalTypes) {
			cv.errorf("unknown global: %d", n)
		}
		gt := cv.mv.globalTypes[n]
		if gt.Mut != 1 {
			cv.errorf(" global is immutable: %d", n)
		}
		cv.popOpdOf(gt.ValType)
	case binary.I32Load:
		cv.i32Load(instr.Args, 32)
	case binary.F32Load:
		cv.f32Load(instr.Args, 32)
	case binary.I64Load:
		cv.i64Load(instr.Args, 64)
	case binary.F64Load:
		cv.f64Load(instr.Args, 64)
	case binary.I32Load8S, binary.I32Load8U:
		cv.i32Load(instr.Args, 8)
	case binary.I32Load16S, binary.I32Load16U:
		cv.i32Load(instr.Args, 16)
	case binary.I64Load8S, binary.I64Load8U:
		cv.i64Load(instr.Args, 8)
	case binary.I64Load16S, binary.I64Load16U:
		cv.i64Load(instr.Args, 16)
	case binary.I64Load32S, binary.I64Load32U:
		cv.i64Load(instr.Args, 32)
	case binary.I32Store:
		cv.i32Store(instr.Args, 32)
	case binary.I64Store:
		cv.i64Store(instr.Args, 64)
	case binary.F32Store:
		cv.f32Store(instr.Args, 32)
	case binary.F64Store:
		cv.f64Store(instr.Args, 64)
	case binary.I32Store8:
		cv.i32Store(instr.Args, 8)
	case binary.I32Store16:
		cv.i32Store(instr.Args, 16)
	case binary.I64Store8:
		cv.i64Store(instr.Args, 8)
	case binary.I64Store16:
		cv.i64Store(instr.Args, 16)
	case binary.I64Store32:
		cv.i64Store(instr.Args, 32)
	case binary.MemorySize:
		cv.checkMem()
		cv.pushI32()
	case binary.MemoryGrow:
		cv.checkMem()
		cv.popI32()
		cv.pushI32()
	case binary.I32Const:
		cv.pushI32()
	case binary.I64Const:
		cv.pushI64()
	case binary.F32Const:
		cv.pushF32()
	case binary.F64Const:
		cv.pushF64()
	case binary.I32Eqz:
		cv.popI32()
		cv.pushI32()
	case binary.I32Eq, binary.I32Ne,
		binary.I32LtS, binary.I32LtU,
		binary.I32GtS, binary.I32GtU,
		binary.I32LeS, binary.I32LeU,
		binary.I32GeS, binary.I32GeU:
		cv.popI32()
		cv.popI32()
		cv.pushI32()
	case binary.I64Eqz:
		cv.popI64()
		cv.pushI32()
	case binary.I64Eq, binary.I64Ne,
		binary.I64LtS, binary.I64LtU,
		binary.I64GtS, binary.I64GtU,
		binary.I64LeS, binary.I64LeU,
		binary.I64GeS, binary.I64GeU:
		cv.popI64()
		cv.popI64()
		cv.pushI32()
	case binary.F32Eq, binary.F32Ne,
		binary.F32Lt, binary.F32Gt,
		binary.F32Le, binary.F32Ge:
		cv.popF32()
		cv.popF32()
		cv.pushI32()
	case binary.F64Eq, binary.F64Ne,
		binary.F64Lt, binary.F64Gt,
		binary.F64Le, binary.F64Ge:
		cv.popF64()
		cv.popF64()
		cv.pushI32()
	case binary.I32Clz, binary.I32Ctz, binary.I32PopCnt:
		cv.popI32()
		cv.pushI32()
	case binary.I32Add, binary.I32Sub, binary.I32Mul,
		binary.I32DivS, binary.I32DivU,
		binary.I32RemS, binary.I32RemU,
		binary.I32And, binary.I32Or, binary.I32Xor,
		binary.I32Shl, binary.I32ShrS, binary.I32ShrU,
		binary.I32Rotl, binary.I32Rotr:
		cv.popOpdOf(I32)
		cv.popOpdOf(I32)
		cv.pushOpd(I32)
	case binary.I64Clz, binary.I64Ctz, binary.I64PopCnt:
		cv.popI64()
		cv.pushI64()
	case binary.I64Add, binary.I64Sub, binary.I64Mul,
		binary.I64DivS, binary.I64DivU,
		binary.I64RemS, binary.I64RemU,
		binary.I64And, binary.I64Or, binary.I64Xor,
		binary.I64Shl, binary.I64ShrS, binary.I64ShrU,
		binary.I64Rotl, binary.I64Rotr:
		cv.popI64()
		cv.popI64()
		cv.pushI64()
	case binary.F32Abs, binary.F32Neg,
		binary.F32Ceil, binary.F32Floor,
		binary.F32Trunc, binary.F32Nearest,
		binary.F32Sqrt:
		cv.popF32()
		cv.pushF32()
	case binary.F32Add, binary.F32Sub,
		binary.F32Mul, binary.F32Div,
		binary.F32Min, binary.F32Max,
		binary.F32CopySign:
		cv.popF32()
		cv.popF32()
		cv.pushF32()
	case binary.F64Abs, binary.F64Neg,
		binary.F64Ceil, binary.F64Floor,
		binary.F64Trunc, binary.F64Nearest,
		binary.F64Sqrt:
		cv.popF64()
		cv.pushF64()
	case binary.F64Add, binary.F64Sub,
		binary.F64Mul, binary.F64Div,
		binary.F64Min, binary.F64Max,
		binary.F64CopySign:
		cv.popF64()
		cv.popF64()
		cv.pushF64()
	case binary.I32WrapI64:
		cv.popI64()
		cv.pushI32()
	case binary.I32TruncF32S, binary.I32TruncF32U:
		cv.popF32()
		cv.pushI32()
	case binary.I32TruncF64S, binary.I32TruncF64U:
		cv.popF64()
		cv.pushI32()
	case binary.I64ExtendI32S, binary.I64ExtendI32U:
		cv.popI32()
		cv.pushI64()
	case binary.I64TruncF32S, binary.I64TruncF32U:
		cv.popF32()
		cv.pushI64()
	case binary.I64TruncF64S, binary.I64TruncF64U:
		cv.popF64()
		cv.pushI64()
	case binary.F32ConvertI32S, binary.F32ConvertI32U:
		cv.popI32()
		cv.pushF32()
	case binary.F32ConvertI64S, binary.F32ConvertI64U:
		cv.popI64()
		cv.pushF32()
	case binary.F32DemoteF64:
		cv.popF64()
		cv.pushF32()
	case binary.F64ConvertI32S, binary.F64ConvertI32U:
		cv.popI32()
		cv.pushF64()
	case binary.F64ConvertI64S, binary.F64ConvertI64U:
		cv.popI64()
		cv.pushF64()
	case binary.F64PromoteF32:
		cv.popF32()
		cv.pushF64()
	case binary.I32ReinterpretF32:
		cv.popF32()
		cv.pushI32()
	case binary.I64ReinterpretF64:
		cv.popF64()
		cv.pushI64()
	case binary.F32ReinterpretI32:
		cv.popI32()
		cv.pushF32()
	case binary.F64ReinterpretI64:
		cv.popI64()
		cv.pushF64()
	case binary.I32Extend8S, binary.I32Extend16S:
		cv.popI32()
		cv.pushI32()
	case binary.I64Extend8S, binary.I64Extend16S, binary.I64Extend32S:
		cv.popI64()
		cv.pushI64()
	case binary.TruncSat:
		switch instr.Args.(byte) {
		case 0, 1:
			cv.popF32()
			cv.pushI32()
		case 2, 3:
			cv.popF64()
			cv.pushI32()
		case 4, 5:
			cv.popF32()
			cv.pushI64()
		case 6, 7:
			cv.popF64()
			cv.pushI64()
		}
	default:
		cv.errorf("unknown opcode: 0x%x", instr.Opcode)
	}
}

/* memory */

func (cv *codeValidator) i32Load(args interface{}, bitWidth int) {
	cv.load(binary.ValTypeI32, bitWidth, args)
}
func (cv *codeValidator) i64Load(args interface{}, bitWidth int) {
	cv.load(binary.ValTypeI64, bitWidth, args)
}
func (cv *codeValidator) f32Load(args interface{}, bitWidth int) {
	cv.load(binary.ValTypeF32, bitWidth, args)
}
func (cv *codeValidator) f64Load(args interface{}, bitWidth int) {
	cv.load(binary.ValTypeF64, bitWidth, args)
}

func (cv *codeValidator) i32Store(args interface{}, bitWidth int) {
	cv.store(binary.ValTypeI32, bitWidth, args)
}
func (cv *codeValidator) i64Store(args interface{}, bitWidth int) {
	cv.store(binary.ValTypeI64, bitWidth, args)
}
func (cv *codeValidator) f32Store(args interface{}, bitWidth int) {
	cv.store(binary.ValTypeF32, bitWidth, args)
}
func (cv *codeValidator) f64Store(args interface{}, bitWidth int) {
	cv.store(binary.ValTypeF64, bitWidth, args)
}

func (cv *codeValidator) load(vt binary.ValType, bitWidth int, args interface{}) {
	cv.checkMem()
	cv.checkAlign(bitWidth, args)
	cv.popI32()
	cv.pushOpd(vt)
}
func (cv *codeValidator) store(vt binary.ValType, bitWidth int, args interface{}) {
	cv.checkMem()
	cv.checkAlign(bitWidth, args)
	cv.popOpdOf(vt)
	cv.popI32()
}
func (cv *codeValidator) checkMem() {
	if cv.mv.getMemCount() == 0 {
		cv.error("unknown memory")
	}
}
func (cv *codeValidator) checkAlign(bitWidth int, args interface{}) {
	align := args.(binary.MemArg).Align
	if 1<<align > bitWidth/8 {
		cv.errorf("alignment must not be larger than natural alignment (%d)",
			bitWidth/8)
	}
}

/* helper */

func isValTypesEq(a, b []valType) bool {
	if len(a) != len(b) {
		return false
	}
	for i, vt := range a {
		if vt != b[i] {
			return false
		}
	}
	return true
}
