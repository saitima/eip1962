#include "textflag.h"

// func mul_two_1(a *[1]uint64)
TEXT ·mul_two_1(SB), NOSPLIT, $0-8
	MOVQ a+0(FP), DI
	XORQ AX, AX
	RCLQ $0x01, (DI)
	RET

// func div_two_1(a *[1]uint64)
TEXT ·div_two_1(SB), NOSPLIT, $0-8
	MOVQ a+0(FP), DI
	XORQ AX, AX
	RCRQ $0x01, (DI)
	RET

// func cpy(dst *[1]uint64, src *[1]uint64)
TEXT ·cpy1(SB), NOSPLIT, $0-16
	MOVQ dst+0(FP), DI
	MOVQ src+8(FP), SI
	MOVQ (SI), R8
	MOVQ R8, (DI)
	RET

// func eq(a *[1]uint64, b *[1]uint64) bool
TEXT ·eq1(SB), NOSPLIT, $0-17
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	MOVB $0x00, ret+16(FP)
	MOVQ (DI), R8
	CMPQ (SI), R8
	JNE  ret
	MOVB $0x01, ret+16(FP)

ret:
	RET

// func cmp(a *[1]uint64, b *[1]uint64) int8
TEXT ·cmp1(SB), NOSPLIT, $0-17
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	MOVQ (DI), R8
	CMPQ (SI), R8
	JB   gt
	JA   lt
	MOVB $0x00, ret+16(FP)
	JMP  ret

gt:
	MOVB $0x01, ret+16(FP)
	JMP  ret

lt:
	MOVB $0xff, ret+16(FP)

ret:
	RET

// func add(c *[1]uint64, a *[1]uint64, b *[1]uint64, p *[1]uint64)
TEXT ·add1(SB), NOSPLIT, $0-32
	// |
	MOVQ a+8(FP), DI
	MOVQ b+16(FP), SI
	XORQ AX, AX

	// |
	MOVQ (DI), CX
	ADDQ (SI), CX
	ADCQ $0x00, AX

	// |
	MOVQ p+24(FP), SI
	MOVQ CX, DX
	SUBQ (SI), DX
	SBBQ $0x00, AX

	// |
	MOVQ    c+0(FP), DI
	CMOVQCC DX, CX
	MOVQ    CX, (DI)
	RET

// func addn(a *[1]uint64, b *[1]uint64) uint64
TEXT ·addn1(SB), NOSPLIT, $0-24
	// |
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI

	// |
	MOVQ (DI), CX
	ADDQ (SI), CX
	ADCQ $0x00, AX

	// |
	MOVQ CX, (DI)
	MOVQ AX, ret+16(FP)
	RET

// func double(c *[1]uint64, a *[1]uint64, p *[1]uint64)
TEXT ·double1(SB), NOSPLIT, $0-24
	// |
	MOVQ a+8(FP), DI
	XORQ AX, AX
	MOVQ (DI), CX
	ADDQ CX, CX
	ADCQ $0x00, AX

	// |
	MOVQ p+16(FP), SI
	MOVQ CX, DX
	SUBQ (SI), DX
	SBBQ $0x00, AX

	// |
	MOVQ    c+0(FP), DI
	CMOVQCC DX, CX
	MOVQ    CX, (DI)
	RET

// func sub(c *[1]uint64, a *[1]uint64, b *[1]uint64, p *[1]uint64)
TEXT ·sub1(SB), NOSPLIT, $0-32
	// |
	MOVQ a+8(FP), DI
	MOVQ b+16(FP), SI
	XORQ AX, AX
	MOVQ (DI), CX
	SUBQ (SI), CX

	// |
	MOVQ    p+24(FP), SI
	MOVQ    (SI), DX
	CMOVQCC AX, DX

	// |
	MOVQ c+0(FP), DI
	ADDQ DX, CX
	MOVQ CX, (DI)
	RET

// func subn(a *[1]uint64, b *[1]uint64) uint64
TEXT ·subn1(SB), NOSPLIT, $0-24
	// |
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	XORQ AX, AX

	// |
	MOVQ (DI), CX
	SUBQ (SI), CX
	ADCQ $0x00, AX

	// |
	MOVQ CX, (DI)
	MOVQ AX, ret+16(FP)
	RET

// func _neg(c *[1]uint64, a *[1]uint64, p *[1]uint64)
TEXT ·_neg1(SB), NOSPLIT, $0-24
	// |
	MOVQ a+8(FP), DI

	// |
	MOVQ p+16(FP), SI
	MOVQ (SI), CX
	SUBQ (DI), CX

	// |
	MOVQ c+0(FP), DI
	MOVQ CX, (DI)
	RET

// func mul(c *[2]uint64, a *[1]uint64, b *[1]uint64, p *[1]uint64, inp uint64)
TEXT ·mul1(SB), NOSPLIT, $0-40
	// | 

/* inputs 				*/

	MOVQ a+8(FP), DI
	MOVQ b+16(FP), SI

	// | 
	MOVQ (SI), DX
	MULXQ (DI), R8, R9

/* swap 				*/

	MOVQ p+24(FP), R15

	// | 
  MOVQ  R8, DX
	MULXQ inp+32(FP), DX, DI

  MULXQ (R15), AX, DI
  ADDQ AX, R8
	ADCQ DI, R9
	ADCQ $0x00, R8

/* reduction 				*/

	MOVQ R9, AX
	SUBQ (R15), AX
	SBBQ $0x00, R8

	// |
	MOVQ    c+0(FP), DI
	CMOVQCC AX, R9
	MOVQ    R9, (DI)
	RET

	// | 

/* end 				*/