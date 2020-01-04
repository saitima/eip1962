package x86

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func generateDiv2(size int, single bool) {
	funcName := "div_two"
	if !single {
		funcName = fmt.Sprintf("%s_%d", funcName, size)
	}
	TEXT(funcName, NOSPLIT, fmt.Sprintf("func(a *[%d]uint64)", size))
	tape := newTape(nil)
	A := tape.newReprAtParam(size, "a", RDI)
	XORQ(RAX, RAX)
	A.previous()
	for i := 0; i < size; i++ {
		RCRQ(Imm(1), A.previous().s)
	}
	RET()
}

func generateMul2(size int, single bool) {
	funcName := "mul_two"
	if !single {
		funcName = fmt.Sprintf("%s_%d", funcName, size)
	}
	TEXT(funcName, NOSPLIT, fmt.Sprintf("func(a *[%d]uint64)", size))
	tape := newTape(nil)
	A := tape.newReprAtParam(size, "a", RDI)
	XORQ(RAX, RAX)
	for i := 0; i < size; i++ {
		RCLQ(Imm(1), A.next(_ITER).s)
	}
	RET()
}

func generateEq(size int, single bool) {
	funcName := "eq"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
	}
	TEXT(funcName, NOSPLIT, fmt.Sprintf("func(a, b *[%d]uint64) bool", size))
	tape := newTape(nil)
	A := tape.newReprAtParam(size, "a", RDI)
	B := tape.newReprAtParam(size, "b", RSI)
	r := NewParamAddr("ret", 16)
	t := R8
	MOVB(U8(0), r)
	for i := 0; i < size; i++ {
		A.next(_ITER).moveTo(t, _NO_ASSIGN)
		B.next(_ITER).cmp(t)
		JNE(LabelRef("ret"))
	}
	MOVB(U8(1), r)
	Label("ret")
	RET()
}

func generateCopy(size int, single bool) {
	funcName := "cpy"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
	}
	TEXT(funcName, NOSPLIT, fmt.Sprintf("func(dst, src *[%d]uint64)", size))
	tape := newTape(nil)
	A := tape.newReprAtParam(size, "dst", RDI)
	B := tape.newReprAtParam(size, "src", RSI)
	t := R8
	for i := 0; i < size; i++ {
		B.next(_ITER).moveTo(t, _NO_ASSIGN)
		A.next(_ITER).load(t, nil)
	}
	RET()
}

// func cmp4(a *[4]uint64, b *[4]uint64) int8
// TEXT ·cmpx(SB), NOSPLIT, $0-17
// 	MOVQ a+0(FP), DI
// 	MOVQ b+8(FP), SI
// 	MOVQ 0(DI), R8
// 	SUBQ 0(SI), R8
// 	JB  lt
// 	JA  gt
// 	MOVB $0x00, ret+16(FP)
// 	JMP  ret

// gt:
// 	MOVB $0x01, ret+16(FP)
// 	JMP  ret

// lt:
// 	MOVB $0xff, ret+16(FP)

// ret:
// 	RET

func generateCmp(size int, single bool) {
	funcName := "cmp"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
	}
	TEXT(funcName, NOSPLIT, fmt.Sprintf("func(a, b *[%d]uint64) int8", size))
	tape := newTape(nil)
	A := tape.newReprAtParam(size, "a", RDI)
	B := tape.newReprAtParam(size, "b", RSI)
	r := NewParamAddr("ret", 16)
	A.previous()
	B.previous()
	t := R8
	for i := 0; i < size; i++ {
		A.previous().moveTo(t, _NO_ASSIGN)
		B.previous().cmp(t)
		JB(LabelRef("gt"))
		JA(LabelRef("lt"))
	}
	MOVB(U8(0), r)
	JMP(LabelRef("ret"))
	Label("gt")
	MOVB(U8(1), r)
	JMP(LabelRef("ret"))
	Label("lt")
	MOVB(U8(0xff), r)
	Label("ret")
	RET()
}

func generateAdd(size int, fixedmod bool, single bool) {
	funcName := "add"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
	}
	if fixedmod {
		TEXT(funcName, NOSPLIT, fmt.Sprintf("func(c, a, b *[%d]uint64)", size))
	} else {
		TEXT(funcName, NOSPLIT, fmt.Sprintf("func(c, a, b, p *[%d]uint64)", size))
	}
	Commentf("|")
	tape := newTape(RBX, RAX)
	A := tape.newReprAtParam(size, "a", RDI)
	B := tape.newReprAtParam(size, "b", RSI)
	C_sum := tape.newReprAlloc(size)
	XORQ(RAX, RAX)
	Commentf("|")
	for i := 0; i < size; i++ {
		C_sum.next(_ITER).loadAdd(
			*A.next(_ITER),
			*B.next(_ITER), i != 0)
	}
	reduceAdded(tape, C_sum, fixedmod, single)
	tape.ret()
	RET()
}

func generateAddNoCar(size int, single bool) {
	funcName := "addn"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
	}
	TEXT(funcName, NOSPLIT, fmt.Sprintf("func(a, b *[%d]uint64) uint64", size))
	Commentf("|")
	tape := newTape(RBX, RAX)
	A := tape.newReprAtParam(size, "a", RDI)
	B := tape.newReprAtParam(size, "b", RSI)
	C_sum := tape.newReprAlloc(size)
	MOVQ(RAX, RAX)
	Commentf("|")
	for i := 0; i < size; i++ {
		C_sum.next(_ITER).loadAdd(
			*A.next(_ITER),
			*B.next(_ITER), i != 0)
	}
	ADCQ(Imm(0), RAX)
	Commentf("|")
	for i := 0; i < size; i++ {
		C_sum.next(_ITER).moveTo(A.next(_ITER), _NO_ASSIGN)
	}
	Store(RAX, ReturnIndex(0))
	tape.ret()
	RET()
}

func generateDouble(size int, fixedmod bool, single bool) {
	funcName := "double"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
	}
	if fixedmod {
		TEXT(funcName, NOSPLIT, fmt.Sprintf("func(c, a *[%d]uint64)", size))
	} else {
		TEXT(funcName, NOSPLIT, fmt.Sprintf("func(c, a, p *[%d]uint64)", size))
	}
	Commentf("|")
	tape := newTape(RBX, RAX)
	if !fixedmod {
		tape.reserveGp(RSI)
	}
	A := tape.newReprAtParam(size, "a", RDI)
	C_sum := tape.newReprAlloc(size)
	XORQ(RAX, RAX)
	for i := 0; i < size; i++ {
		C_sum.next(_ITER).loadDouble(*A.next(_ITER), i != 0)
	}
	reduceAdded(tape, C_sum, fixedmod, single)
	tape.ret()
	RET()
}

func reduceAdded(tape *tape, C_sum *repr, fixedmod bool, single bool) {
	size := C_sum.size
	modulusName := "·modulus"
	if !single {
		modulusName = fmt.Sprintf("%s%d", modulusName, size)
	}
	ADCQ(Imm(0), RAX)
	Commentf("|")
	var modulus *repr
	if fixedmod {
		modulus = tape.newReprAtMemory(size, NewDataAddr(Symbol{Name: modulusName}, 0))
	} else {
		modulus = tape.newReprAtParam(size, "p", RSI)
	}
	C_red := tape.newReprAlloc(size)
	for i := 0; i < size; i++ {
		C_red.next(_ITER).loadSubSafe(*C_sum.next(_ITER), *modulus.next(_ITER), i != 0)
	}
	SBBQ(Imm(0), RAX)
	Commentf("|")
	C := tape.newReprAtParam(size, "c", RDI)
	for i := 0; i < size; i++ {
		C_red.next(_ITER).moveIfNotCFAux(*C_sum.next(_ITER), *C.next(_ITER))
	}
}

func generateSub(size int, fixedmod bool, single bool) {
	funcName := "sub"
	modulusName := "·modulus"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
		modulusName = fmt.Sprintf("%s%d", modulusName, size)
	}
	if fixedmod {
		TEXT(funcName, NOSPLIT, fmt.Sprintf("func(c, a, b *[%d]uint64)", size))
	} else {
		TEXT(funcName, NOSPLIT, fmt.Sprintf("func(c, a, b, p *[%d]uint64)", size))
	}
	Commentf("|")
	tape := newTape(RBX, RAX)
	A := tape.newReprAtParam(size, "a", RDI)
	B := tape.newReprAtParam(size, "b", RSI)
	C_sub := tape.newReprAlloc(size)
	zero := tape.newReprNoAlloc(size)
	for i := 0; i < size; i++ {
		zero.next(_ITER).set(RAX)
	}
	XORQ(RAX, RAX)
	for i := 0; i < size; i++ {
		C_sub.next(_ITER).loadSub(*A.next(_ITER), *B.next(_ITER), i != 0)
	}
	Commentf("|")
	var modulus *repr
	if fixedmod {
		tape.free(B.base)
		modulus = tape.newReprAtMemory(size, NewDataAddr(Symbol{Name: modulusName}, 0))
	} else {
		modulus = tape.newReprAtParam(size, "p", B.base)
	}
	C_mod := tape.newReprAlloc(size)
	for i := 0; i < size; i++ {
		zero.next(_ITER).moveIfNotCFAux(*modulus.next(_ITER), *C_mod.next(_ITER))
	}
	Commentf("|")
	C := tape.newReprAtParam(size, "c", RDI)
	for i := 0; i < size; i++ {
		C.next(_ITER).loadAdd(*C_sub.next(_ITER), *C_mod.next(_ITER), i != 0)
	}
	tape.ret()
	RET()
}

func generateSubNoCar(size int, single bool) {
	funcName := "subn"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
	}
	TEXT(funcName, NOSPLIT, fmt.Sprintf("func(a, b *[%d]uint64) uint64", size))
	Commentf("|")
	tape := newTape(RBX, RAX)
	A := tape.newReprAtParam(size, "a", RDI)
	B := tape.newReprAtParam(size, "b", RSI)
	C_sum := tape.newReprAlloc(size)
	XORQ(RAX, RAX)
	Commentf("|")
	for i := 0; i < size; i++ {
		C_sum.next(_ITER).loadSub(*A.next(_ITER), *B.next(_ITER), i != 0)
	}
	ADCQ(Imm(0), RAX)
	Commentf("|")
	for i := 0; i < size; i++ {
		C_sum.next(_ITER).moveTo(A.next(_ITER), _NO_ASSIGN)
	}
	Store(RAX, ReturnIndex(0))
	tape.ret()
	RET()
}

func generateNeg(size int, fixedmod bool, single bool) {
	funcName := "_neg"
	modulusName := "·modulus"
	if !single {
		funcName = fmt.Sprintf("%s%d", funcName, size)
		modulusName = fmt.Sprintf("%s%d", modulusName, size)
	}
	TEXT(funcName, NOSPLIT, fmt.Sprintf("func(c, a, p *[%d]uint64)", size))
	Commentf("|")
	tape := newTape(RBX, RAX)
	A := tape.newReprAtParam(size, "a", RDI)
	if !fixedmod {
		tape.reserveGp(RSI)
	}
	C_sub := tape.newReprAlloc(size)
	Commentf("|")
	var modulus *repr
	if fixedmod {
		modulus = tape.newReprAtMemory(size, NewDataAddr(Symbol{Name: modulusName}, 0))
	} else {
		modulus = tape.newReprAtParam(size, "p", RSI)
	}
	for i := 0; i < size; i++ {
		C_sub.next(_ITER).loadSub(*modulus.next(_ITER), *A.next(_ITER), i != 0)
	}
	Commentf("|")
	C := tape.newReprAtParam(size, "c", RDI)
	for i := 0; i < size; i++ {
		C_sub.next(_ITER).moveTo(C.next(_ITER), _NO_ASSIGN)
	}
	tape.ret()
	RET()
}
