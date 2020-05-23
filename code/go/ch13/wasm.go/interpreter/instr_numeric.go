package interpreter

import (
	"math"
	"math/bits"
)

// const
func i32Const(vm *vm, args interface{}) {
	vm.pushS32(args.(int32))
}
func i64Const(vm *vm, args interface{}) {
	vm.pushS64(args.(int64))
}
func f32Const(vm *vm, args interface{}) {
	vm.pushF32(args.(float32))
}
func f64Const(vm *vm, args interface{}) {
	vm.pushF64(args.(float64))
}

// i32 test & rel
func i32Eqz(vm *vm, _ interface{}) {
	vm.pushBool(vm.popU32() == 0)
}
func i32Eq(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushBool(v1 == v2)
}
func i32Ne(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushBool(v1 != v2)
}
func i32LtS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS32(), vm.popS32()
	vm.pushBool(v1 < v2)
}
func i32LtU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushBool(v1 < v2)
}
func i32GtS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS32(), vm.popS32()
	vm.pushBool(v1 > v2)
}
func i32GtU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushBool(v1 > v2)
}
func i32LeS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS32(), vm.popS32()
	vm.pushBool(v1 <= v2)
}
func i32LeU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushBool(v1 <= v2)
}
func i32GeS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS32(), vm.popS32()
	vm.pushBool(v1 >= v2)
}
func i32GeU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushBool(v1 >= v2)
}

// i64 test & rel
func i64Eqz(vm *vm, _ interface{}) {
	vm.pushBool(vm.popU64() == 0)
}
func i64Eq(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushBool(v1 == v2)
}
func i64Ne(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushBool(v1 != v2)
}
func i64LtS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS64(), vm.popS64()
	vm.pushBool(v1 < v2)
}
func i64LtU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushBool(v1 < v2)
}
func i64GtS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS64(), vm.popS64()
	vm.pushBool(v1 > v2)
}
func i64GtU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushBool(v1 > v2)
}
func i64LeS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS64(), vm.popS64()
	vm.pushBool(v1 <= v2)
}
func i64LeU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushBool(v1 <= v2)
}
func i64GeS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS64(), vm.popS64()
	vm.pushBool(v1 >= v2)
}
func i64GeU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushBool(v1 >= v2)
}

// f32 rel
func f32Eq(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushBool(v1 == v2)
}
func f32Ne(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushBool(v1 != v2)
}
func f32Lt(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushBool(v1 < v2)
}
func f32Gt(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushBool(v1 > v2)
}
func f32Le(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushBool(v1 <= v2)
}
func f32Ge(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushBool(v1 >= v2)
}

// f64 rel
func f64Eq(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushBool(v1 == v2)
}
func f64Ne(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushBool(v1 != v2)
}
func f64Lt(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushBool(v1 < v2)
}
func f64Gt(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushBool(v1 > v2)
}
func f64Le(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushBool(v1 <= v2)
}
func f64Ge(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushBool(v1 >= v2)
}

// i32 arithmetic & bitwise
func i32Clz(vm *vm, _ interface{}) {
	vm.pushU32(uint32(bits.LeadingZeros32(vm.popU32())))
}
func i32Ctz(vm *vm, _ interface{}) {
	vm.pushU32(uint32(bits.TrailingZeros32(vm.popU32())))
}
func i32PopCnt(vm *vm, _ interface{}) {
	vm.pushU32(uint32(bits.OnesCount32(vm.popU32())))
}
func i32Add(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 + v2)
}
func i32Sub(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 - v2)
}
func i32Mul(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 * v2)
}
func i32DivS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS32(), vm.popS32()
	if v1 == math.MinInt32 && v2 == -1 {
		panic(errIntOverflow)
	}
	vm.pushS32(v1 / v2)
}
func i32DivU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 / v2)
}
func i32RemS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS32(), vm.popS32()
	vm.pushS32(v1 % v2)
}
func i32RemU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 % v2)
}
func i32And(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 & v2)
}
func i32Or(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 | v2)
}
func i32Xor(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 ^ v2)
}
func i32Shl(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 << (v2 % 32))
}
func i32ShrS(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popS32()
	vm.pushS32(v1 >> (v2 % 32))
}
func i32ShrU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(v1 >> (v2 % 32))
}
func i32Rotl(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(bits.RotateLeft32(v1, int(v2)))
}
func i32Rotr(vm *vm, _ interface{}) {
	v2, v1 := vm.popU32(), vm.popU32()
	vm.pushU32(bits.RotateLeft32(v1, -int(v2)))
}

// i64 arithmetic & bitwise
func i64Clz(vm *vm, _ interface{}) {
	vm.pushU64(uint64(bits.LeadingZeros64(vm.popU64())))
}
func i64Ctz(vm *vm, _ interface{}) {
	vm.pushU64(uint64(bits.TrailingZeros64(vm.popU64())))
}
func i64PopCnt(vm *vm, _ interface{}) {
	vm.pushU64(uint64(bits.OnesCount64(vm.popU64())))
}
func i64Add(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 + v2)
}
func i64Sub(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 - v2)
}
func i64Mul(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 * v2)
}
func i64DivS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS64(), vm.popS64()
	if v1 == math.MinInt64 && v2 == -1 {
		panic(errIntOverflow)
	}
	vm.pushS64(v1 / v2)
}
func i64DivU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 / v2)
}
func i64RemS(vm *vm, _ interface{}) {
	v2, v1 := vm.popS64(), vm.popS64()
	vm.pushS64(v1 % v2)
}
func i64RemU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 % v2)
}
func i64And(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 & v2)
}
func i64Or(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 | v2)
}
func i64Xor(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 ^ v2)
}
func i64Shl(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 << (v2 % 64))
}
func i64ShrS(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popS64()
	vm.pushS64(v1 >> (v2 % 64))
}
func i64ShrU(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(v1 >> (v2 % 64))
}
func i64Rotl(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(bits.RotateLeft64(v1, int(v2)))
}
func i64Rotr(vm *vm, _ interface{}) {
	v2, v1 := vm.popU64(), vm.popU64()
	vm.pushU64(bits.RotateLeft64(v1, -int(v2)))
}

// f32 arithmetic
func f32Abs(vm *vm, _ interface{}) {
	vm.pushF32(float32(math.Abs(float64(vm.popF32()))))
}
func f32Neg(vm *vm, _ interface{}) {
	vm.pushF32(-vm.popF32())
}
func f32Ceil(vm *vm, _ interface{}) {
	vm.pushF32(float32(math.Ceil(float64(vm.popF32()))))
}
func f32Floor(vm *vm, _ interface{}) {
	vm.pushF32(float32(math.Floor(float64(vm.popF32()))))
}
func f32Trunc(vm *vm, _ interface{}) {
	vm.pushF32(float32(math.Trunc(float64(vm.popF32()))))
}
func f32Nearest(vm *vm, _ interface{}) {
	vm.pushF32(float32(math.RoundToEven(float64(vm.popF32()))))
}
func f32Sqrt(vm *vm, _ interface{}) {
	vm.pushF32(float32(math.Sqrt(float64(vm.popF32()))))
}
func f32Add(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushF32(v1 + v2)
}
func f32Sub(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushF32(v1 - v2)
}
func f32Mul(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushF32(v1 * v2)
}
func f32Div(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushF32(v1 / v2)
}
func f32Min(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	v1NaN := math.IsNaN(float64(v1))
	v2NaN := math.IsNaN(float64(v2))
	if v1NaN && !v2NaN {
		vm.pushF32(v1)
		return
	} else if v2NaN && !v1NaN {
		vm.pushF32(v2)
		return
	}
	vm.pushF32(float32(math.Min(float64(v1), float64(v2))))
}
func f32Max(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	v1NaN := math.IsNaN(float64(v1))
	v2NaN := math.IsNaN(float64(v2))
	if v1NaN && !v2NaN {
		vm.pushF32(v1)
		return
	} else if v2NaN && !v1NaN {
		vm.pushF32(v2)
		return
	}
	vm.pushF32(float32(math.Max(float64(v1), float64(v2))))
}
func f32CopySign(vm *vm, _ interface{}) {
	v2, v1 := vm.popF32(), vm.popF32()
	vm.pushF32(float32(math.Copysign(float64(v1), float64(v2))))
}

// f64 arithmetic
func f64Abs(vm *vm, _ interface{}) {
	vm.pushF64(math.Abs(vm.popF64()))
}
func f64Neg(vm *vm, _ interface{}) {
	vm.pushF64(-vm.popF64())
}
func f64Ceil(vm *vm, _ interface{}) {
	vm.pushF64(math.Ceil(vm.popF64()))
}
func f64Floor(vm *vm, _ interface{}) {
	vm.pushF64(math.Floor(vm.popF64()))
}
func f64Trunc(vm *vm, _ interface{}) {
	vm.pushF64(math.Trunc(vm.popF64()))
}
func f64Nearest(vm *vm, _ interface{}) {
	vm.pushF64(math.RoundToEven(vm.popF64()))
}
func f64Sqrt(vm *vm, _ interface{}) {
	vm.pushF64(math.Sqrt(vm.popF64()))
}
func f64Add(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushF64(v1 + v2)
}
func f64Sub(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushF64(v1 - v2)
}
func f64Mul(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushF64(v1 * v2)
}
func f64Div(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushF64(v1 / v2)
}
func f64Min(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	v1NaN := math.IsNaN(v1)
	v2NaN := math.IsNaN(v2)
	if v1NaN && !v2NaN {
		vm.pushF64(v1)
		return
	} else if v2NaN && !v1NaN {
		vm.pushF64(v2)
		return
	}
	vm.pushF64(math.Min(v1, v2))
}
func f64Max(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	v1NaN := math.IsNaN(v1)
	v2NaN := math.IsNaN(v2)
	if v1NaN && !v2NaN {
		vm.pushF64(v1)
		return
	} else if v2NaN && !v1NaN {
		vm.pushF64(v2)
		return
	}
	vm.pushF64(math.Max(v1, v2))
}
func f64CopySign(vm *vm, _ interface{}) {
	v2, v1 := vm.popF64(), vm.popF64()
	vm.pushF64(math.Copysign(v1, v2))
}

// conversions
func i32WrapI64(vm *vm, _ interface{}) {
	vm.pushU32(uint32(vm.popU64()))
}
func i32TruncF32S(vm *vm, _ interface{}) {
	f := math.Trunc(float64(vm.popF32()))
	if f > math.MaxInt32 || f < math.MinInt32 {
		panic(errIntOverflow)
	}
	if math.IsNaN(f) {
		panic(errConvertToInt)
	}
	vm.pushS32(int32(f))
}
func i32TruncF32U(vm *vm, _ interface{}) {
	f := math.Trunc(float64(vm.popF32()))
	if f > math.MaxUint32 || f < 0 {
		panic(errIntOverflow)
	}
	if math.IsNaN(f) {
		panic(errConvertToInt)
	}
	vm.pushU32(uint32(f))
}
func i32TruncF64S(vm *vm, _ interface{}) {
	f := math.Trunc(vm.popF64())
	if f > math.MaxInt32 || f < math.MinInt32 {
		panic(errIntOverflow)
	}
	if math.IsNaN(f) {
		panic(errConvertToInt)
	}
	vm.pushS32(int32(f))
}
func i32TruncF64U(vm *vm, _ interface{}) {
	f := math.Trunc(vm.popF64())
	if f > math.MaxUint32 || f < 0 {
		panic(errIntOverflow)
	}
	if math.IsNaN(f) {
		panic(errConvertToInt)
	}
	vm.pushU32(uint32(f))
}
func i64ExtendI32S(vm *vm, _ interface{}) {
	vm.pushS64(int64(vm.popS32()))
}
func i64ExtendI32U(vm *vm, _ interface{}) {
	vm.pushU64(uint64(vm.popU32()))
}
func i64TruncF32S(vm *vm, _ interface{}) {
	f := math.Trunc(float64(vm.popF32()))
	if f >= math.MaxInt64 || f < math.MinInt64 {
		panic(errIntOverflow)
	}
	if math.IsNaN(f) {
		panic(errConvertToInt)
	}
	vm.pushS64(int64(f))
}
func i64TruncF32U(vm *vm, _ interface{}) {
	f := math.Trunc(float64(vm.popF32()))
	if f >= math.MaxUint64 || f < 0 {
		panic(errIntOverflow)
	}
	if math.IsNaN(f) {
		panic(errConvertToInt)
	}
	vm.pushU64(uint64(f))
}
func i64TruncF64S(vm *vm, _ interface{}) {
	f := math.Trunc(vm.popF64())
	if f >= math.MaxInt64 || f < math.MinInt64 {
		panic(errIntOverflow)
	}
	if math.IsNaN(f) {
		panic(errConvertToInt)
	}
	vm.pushS64(int64(f))
}
func i64TruncF64U(vm *vm, _ interface{}) {
	f := math.Trunc(vm.popF64())
	if f >= math.MaxUint64 || f < 0 {
		panic(errIntOverflow)
	}
	if math.IsNaN(f) {
		panic(errConvertToInt)
	}
	vm.pushU64(uint64(f))
}
func f32ConvertI32S(vm *vm, _ interface{}) {
	vm.pushF32(float32(vm.popS32()))
}
func f32ConvertI32U(vm *vm, _ interface{}) {
	vm.pushF32(float32(vm.popU32()))
}
func f32ConvertI64S(vm *vm, _ interface{}) {
	vm.pushF32(float32(vm.popS64()))
}
func f32ConvertI64U(vm *vm, _ interface{}) {
	vm.pushF32(float32(vm.popU64()))
}
func f32DemoteF64(vm *vm, _ interface{}) {
	vm.pushF32(float32(vm.popF64()))
}
func f64ConvertI32S(vm *vm, _ interface{}) {
	vm.pushF64(float64(vm.popS32()))
}
func f64ConvertI32U(vm *vm, _ interface{}) {
	vm.pushF64(float64(vm.popU32()))
}
func f64ConvertI64S(vm *vm, _ interface{}) {
	vm.pushF64(float64(vm.popS64()))
}
func f64ConvertI64U(vm *vm, _ interface{}) {
	vm.pushF64(float64(vm.popU64()))
}
func f64PromoteF32(vm *vm, _ interface{}) {
	vm.pushF64(float64(vm.popF32()))
}
func i32ReinterpretF32(vm *vm, _ interface{}) {
	//vm.pushU32(math.Float32bits(vm.popF32()))
}
func i64ReinterpretF64(vm *vm, _ interface{}) {
	//vm.pushU64(math.Float64bits(vm.popF64()))
}
func f32ReinterpretI32(vm *vm, _ interface{}) {
	//vm.pushF32(math.Float32frombits(vm.popU32()))
}
func f64ReinterpretI64(vm *vm, _ interface{}) {
	//vm.pushF64(math.Float64frombits(vm.popU64()))
}

func i32Extend8S(vm *vm, _ interface{}) {
	vm.pushS32(int32(int8(vm.popS32())))
}
func i32Extend16S(vm *vm, _ interface{}) {
	vm.pushS32(int32(int16(vm.popS32())))
}
func i64Extend8S(vm *vm, _ interface{}) {
	vm.pushS64(int64(int8(vm.popS64())))
}
func i64Extend16S(vm *vm, _ interface{}) {
	vm.pushS64(int64(int16(vm.popS64())))
}
func i64Extend32S(vm *vm, _ interface{}) {
	vm.pushS64(int64(int32(vm.popS64())))
}

func truncSat(vm *vm, args interface{}) {
	switch args.(byte) {
	case 0: // i32.trunc_sat_f32_s
		v := truncSatS(float64(vm.popF32()), 32)
		vm.pushS32(int32(v))
	case 1: // i32.trunc_sat_f32_u
		v := truncSatU(float64(vm.popF32()), 32)
		vm.pushU32(uint32(v))
	case 2: // i32.trunc_sat_f64_s
		v := truncSatS(vm.popF64(), 32)
		vm.pushS32(int32(v))
	case 3: // i32.trunc_sat_f64_u
		v := truncSatU(vm.popF64(), 32)
		vm.pushU32(uint32(v))
	case 4: // i64.trunc_sat_f32_s
		v := truncSatS(float64(vm.popF32()), 64)
		vm.pushS64(v)
	case 5: // i64.trunc_sat_f32_u
		v := truncSatU(float64(vm.popF32()), 64)
		vm.pushU64(v)
	case 6: // i64.trunc_sat_f64_s
		v := truncSatS(vm.popF64(), 64)
		vm.pushS64(v)
	case 7: // i64.trunc_sat_f64_u
		v := truncSatU(vm.popF64(), 64)
		vm.pushU64(v)
	default:
		panic("unreachable")
	}
}

func truncSatU(z float64, n int) uint64 {
	if math.IsNaN(z) {
		return 0
	}
	if math.IsInf(z, -1) {
		return 0
	}
	max := (uint64(1) << n) - 1
	if math.IsInf(z, 1) {
		return max
	}
	if x := math.Trunc(z); x < 0 {
		return 0
	} else if x >= float64(max) {
		return max
	} else {
		return uint64(x)
	}
}
func truncSatS(z float64, n int) int64 {
	if math.IsNaN(z) {
		return 0
	}
	min := -(int64(1) << (n - 1))
	max := (int64(1) << (n - 1)) - 1
	if math.IsInf(z, -1) {
		return min
	}
	if math.IsInf(z, 1) {
		return max
	}
	if x := math.Trunc(z); x < float64(min) {
		return min
	} else if x >= float64(max) {
		return max
	} else {
		return int64(x)
	}
}
