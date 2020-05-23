package interpreter

import "math"

type operandStack struct {
	slots []uint64
}

func (s *operandStack) stackSize() int {
	return len(s.slots)
}

func (s *operandStack) getOperand(idx uint32) uint64 {
	return s.slots[idx]
}
func (s *operandStack) setOperand(idx uint32, val uint64) {
	s.slots[idx] = val
}

func (s *operandStack) pushU64s(vals []uint64) {
	s.slots = append(s.slots, vals...)
}
func (s *operandStack) popU64s(n int) []uint64 {
	vals := s.slots[len(s.slots)-n:]
	s.slots = s.slots[:len(s.slots)-n]
	return vals
}

func (s *operandStack) pushU64(val uint64) {
	s.slots = append(s.slots, val)
}
func (s *operandStack) popU64() uint64 {
	val := s.slots[len(s.slots)-1]
	s.slots = s.slots[:len(s.slots)-1]
	return val
}

func (s *operandStack) pushS64(val int64) {
	s.pushU64(uint64(val))
}
func (s *operandStack) popS64() int64 {
	return int64(s.popU64())
}

func (s *operandStack) pushU32(val uint32) {
	s.pushU64(uint64(val))
}
func (s *operandStack) popU32() uint32 {
	return uint32(s.popU64())
}

func (s *operandStack) pushS32(val int32) {
	s.pushU32(uint32(val))
}
func (s *operandStack) popS32() int32 {
	return int32(s.popU32())
}

func (s *operandStack) pushF32(val float32) {
	s.pushU32(math.Float32bits(val))
}
func (s *operandStack) popF32() float32 {
	return math.Float32frombits(s.popU32())
}

func (s *operandStack) pushF64(val float64) {
	s.pushU64(math.Float64bits(val))
}
func (s *operandStack) popF64() float64 {
	return math.Float64frombits(s.popU64())
}

func (s *operandStack) pushBool(val bool) {
	if val {
		s.pushU64(1)
	} else {
		s.pushU64(0)
	}
}
func (s *operandStack) popBool() bool {
	return s.popU64() != 0
}
