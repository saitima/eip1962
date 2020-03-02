package eip

import (
	"fmt"
	"io"
	"math/big"
)

type fe6C [3]fe2

type fq6C struct {
	f               *fq2
	nonResidue      *fe2
	t               []*fe2
	frobeniusCoeffs *[2][6]*fe2
}

func newFq6Cubic(fq2 *fq2, nonResidueBuf []byte) (*fq6C, error) {
	nonResidue := fq2.new()
	if nonResidueBuf != nil {
		var err error
		nonResidue, err = fq2.fromBytes(nonResidueBuf)
		if err != nil {
			return nil, err
		}
	}
	t := make([]*fe2, 6)
	for i := 0; i < 6; i++ {
		t[i] = fq2.new()
	}
	return &fq6C{fq2, nonResidue, t, nil}, nil
}

func (fq6 *fq6C) byteSize() int {
	fq2 := fq6.fq2()
	return fq2.byteSize() * 3
}

func (fq6 *fq6C) new() *fe6C {
	return fq6.zero()
}

func (fq6 *fq6C) modulus() *big.Int {
	fq2 := fq6.fq2()
	return fq2.modulus()
}

func (fq6 *fq6C) fromBytes(in []byte) (*fe6C, error) {
	fq2 := fq6.fq2()
	byteSize := fq2.byteSize()
	if len(in) != byteSize*3 {
		return nil, fmt.Errorf("input string should be larger than %d bytes", byteSize*3)
	}
	var err error
	c0, err := fq2.fromBytes(in[:byteSize])
	if err != nil {
		return nil, err
	}
	c1, err := fq2.fromBytes(in[byteSize : byteSize*2])
	if err != nil {
		return nil, err
	}
	c2, err := fq2.fromBytes(in[2*byteSize:])
	if err != nil {
		return nil, err
	}
	return &fe6C{*c0, *c1, *c2}, nil
}

func (fq6 *fq6C) toBytes(a *fe6C) []byte {
	fq2 := fq6.fq2()
	byteSize := fq2.byteSize()
	out := make([]byte, 3*byteSize)
	copy(out[:byteSize], fq2.toBytes(&a[0]))
	copy(out[byteSize:2*byteSize], fq2.toBytes(&a[1]))
	copy(out[2*byteSize:], fq2.toBytes(&a[2]))
	return out
}

func (fq6 *fq6C) toString(a *fe6C) string {
	fq2 := fq6.fq2()
	return fmt.Sprintf("%s\n%s\n%s", fq2.toString(&a[0]), fq2.toString(&a[1]), fq2.toString(&a[2]))
}

func (fq6 *fq6C) toStringNoTransform(a *fe6C) string {
	fq2 := fq6.fq2()
	return fmt.Sprintf("%s\n%s\n%s", fq2.toStringNoTransform(&a[0]), fq2.toStringNoTransform(&a[1]), fq2.toStringNoTransform(&a[2]))
}

func (fq6 *fq6C) zero() *fe6C {
	fq2 := fq6.fq2()
	return &fe6C{*fq2.zero(), *fq2.zero(), *fq2.zero()}
}

func (fq6 *fq6C) one() *fe6C {
	fq2 := fq6.fq2()
	a := fq6.zero()
	fq2.copy(&a[0], fq2.one())
	return a
}

func (fq6 *fq6C) rand(r io.Reader) *fe6C {
	fq2 := fq6.fq2()
	return &fe6C{*fq2.rand(r), *fq2.rand(r), *fq2.rand(r)}
}

func (fq6 *fq6C) isZero(a *fe6C) bool {
	fq2 := fq6.fq2()
	return fq2.isZero(&a[0]) && fq2.isZero(&a[1]) && fq2.isZero(&a[2])
}

func (fq6 *fq6C) isOne(a *fe6C) bool {
	fq2 := fq6.fq2()
	return fq2.isOne(&a[0]) && fq2.isZero(&a[1]) && fq2.isZero(&a[2])
}

func (fq6 *fq6C) equal(a, b *fe6C) bool {
	fq2 := fq6.fq2()
	return fq2.equal(&a[0], &b[0]) && fq2.equal(&a[1], &b[1]) && fq2.equal(&a[2], &b[2])
}

func (fq6 *fq6C) copy(c, a *fe6C) *fe6C {
	fq2 := fq6.fq2()
	fq2.copy(&c[0], &a[0])
	fq2.copy(&c[1], &a[1])
	fq2.copy(&c[2], &a[2])
	return c
}

func (fq6 *fq6C) add(c, a, b *fe6C) *fe6C {
	fq2 := fq6.fq2()
	fq2.add(&c[0], &a[0], &b[0])
	fq2.add(&c[1], &a[1], &b[1])
	fq2.add(&c[2], &a[2], &b[2])
	return c
}

func (fq6 *fq6C) double(c, a *fe6C) *fe6C {
	fq2 := fq6.fq2()
	fq2.double(&c[0], &a[0])
	fq2.double(&c[1], &a[1])
	fq2.double(&c[2], &a[2])
	return c
}

func (fq6 *fq6C) sub(c, a, b *fe6C) *fe6C {
	fq2 := fq6.fq2()
	fq2.sub(&c[0], &a[0], &b[0])
	fq2.sub(&c[1], &a[1], &b[1])
	fq2.sub(&c[2], &a[2], &b[2])
	return c
}

func (fq6 *fq6C) neg(c, a *fe6C) *fe6C {
	fq2 := fq6.fq2()
	fq2.neg(&c[0], &a[0])
	fq2.neg(&c[1], &a[1])
	fq2.neg(&c[2], &a[2])
	return c
}

func (fq6 *fq6C) conjugate(c, a *fe6C) *fe6C {
	fq2 := fq6.fq2()
	fq6.copy(c, a)
	fq2.neg(&c[1], &a[1])
	return c
}

func (fq6 *fq6C) mulByNonResidue(c, a *fe2) {
	fq2 := fq6.fq2()
	fq2.mul(c, a, fq6.nonResidue)
}

func (fq6 *fq6C) mul(c, a, b *fe6C) {
	fq2, t := fq6.fq2(), fq6.t
	fq2.mul(t[0], &a[0], &b[0])
	fq2.mul(t[1], &a[1], &b[1])
	fq2.mul(t[2], &a[2], &b[2])
	fq2.add(t[3], &a[1], &a[2])
	fq2.add(t[4], &b[1], &b[2])
	fq2.mul(t[3], t[3], t[4])
	fq2.add(t[4], t[1], t[2])
	fq2.sub(t[3], t[3], t[4])
	fq6.mulByNonResidue(t[3], t[3])
	fq2.add(t[5], t[0], t[3])
	fq2.add(t[3], &a[0], &a[1])
	fq2.add(t[4], &b[0], &b[1])
	fq2.mul(t[3], t[3], t[4])
	fq2.add(t[4], t[0], t[1])
	fq2.sub(t[3], t[3], t[4])
	fq6.mulByNonResidue(t[4], t[2])
	fq2.add(&c[1], t[3], t[4])
	fq2.add(t[3], &a[0], &a[2])
	fq2.add(t[4], &b[0], &b[2])
	fq2.mul(t[3], t[3], t[4])
	fq2.add(t[4], t[0], t[2])
	fq2.sub(t[3], t[3], t[4])
	fq2.add(&c[2], t[1], t[3])
	fq2.copy(&c[0], t[5])
}

func (fq6 *fq6C) square(c, a *fe6C) {
	fq2, t := fq6.fq2(), fq6.t
	fq2.square(t[0], &a[0])
	fq2.mul(t[1], &a[0], &a[1])
	fq2.add(t[1], t[1], t[1])
	fq2.sub(t[2], &a[0], &a[1])
	fq2.add(t[2], t[2], &a[2])
	fq2.square(t[2], t[2])
	fq2.mul(t[3], &a[1], &a[2])
	fq2.add(t[3], t[3], t[3])
	fq2.square(t[4], &a[2])
	fq6.mulByNonResidue(t[5], t[3])
	fq2.add(&c[0], t[0], t[5])
	fq6.mulByNonResidue(t[5], t[4])
	fq2.add(&c[1], t[1], t[5])
	fq2.add(t[1], t[1], t[2])
	fq2.add(t[1], t[1], t[3])
	fq2.add(t[0], t[0], t[4])
	fq2.sub(&c[2], t[1], t[0])
}

func (fq6 *fq6C) inverse(c, a *fe6C) bool {
	fq2, t := fq6.fq2(), fq6.t
	fq2.square(t[0], &a[0])                 // v0 = a0^2
	fq2.mul(t[1], &a[1], &a[2])             // v5 = a1 * a2
	fq6.mulByNonResidue(t[1], t[1])         // α * v5
	fq2.sub(t[0], t[0], t[1])               // A = v0 - α * v5
	fq2.square(t[1], &a[1])                 // v1 = a1^2
	fq2.mul(t[2], &a[0], &a[2])             // v4 = a0 * a2
	fq2.sub(t[1], t[1], t[2])               // C = v1 - v4
	fq2.square(t[2], &a[2])                 // v2 = a2^2
	fq6.mulByNonResidue(t[2], t[2])         // α * v2
	fq2.mul(t[3], &a[0], &a[1])             // v3 = a0 * a1
	fq2.sub(t[2], t[2], t[3])               // B = α * v2 - v3
	fq2.mul(t[3], &a[2], t[2])              // B * a2
	fq2.mul(t[4], &a[1], t[1])              // C * a1
	fq2.add(t[3], t[3], t[4])               // C * a1 + B * a2
	fq6.mulByNonResidue(t[3], t[3])         // α * (C * a1 + B * a2)
	fq2.mul(t[4], &a[0], t[0])              // A * a0
	fq2.add(t[3], t[3], t[4])               // v6 = A * a0 * α * (C * a1 + B * a2)
	if ok := fq2.inverse(t[3], t[3]); !ok { // F = v6^-1
		fq6.copy(c, fq6.zero())
		return false
	}
	fq2.mul(&c[0], t[0], t[3]) // c0 = A * F
	fq2.mul(&c[1], t[2], t[3]) // c1 = B * F
	fq2.mul(&c[2], t[1], t[3]) // c2 = C * F
	return true
}

func (fq6 *fq6C) exp(c, a *fe6C, e *big.Int) {
	z := fq6.one()
	found := false
	for i := e.BitLen() - 1; i >= 0; i-- {
		if found {
			fq6.square(z, z)
		} else {
			found = e.Bit(i) == 1
		}
		if e.Bit(i) == 1 {
			fq6.mul(z, z, a)
		}
	}
	fq6.copy(c, z)
}

func (fq6 *fq6C) mulBy01(e *fe6C, c0, c1 *fe2) {
	fq2, t := fq6.fq2(), fq6.t
	fq2.mul(t[0], &e[0], c0)
	fq2.mul(t[1], &e[1], c1)
	fq2.add(t[5], &e[1], &e[2])
	fq2.mul(t[2], c1, t[5])
	fq2.sub(t[2], t[2], t[1])
	fq6.mulByNonResidue(t[2], t[2])
	fq2.add(t[2], t[2], t[0])
	fq2.add(t[5], &e[0], &e[2])
	fq2.mul(t[3], c0, t[5])
	fq2.sub(t[3], t[3], t[0])
	fq2.add(t[3], t[3], t[1])
	fq2.add(t[4], c0, c1)
	fq2.add(t[5], &e[0], &e[1])
	fq2.mul(t[4], t[4], t[5])
	fq2.sub(t[4], t[4], t[0])
	fq2.sub(t[4], t[4], t[1])
	fq2.copy(&e[0], t[2])
	fq2.copy(&e[1], t[4])
	fq2.copy(&e[2], t[3])
}

func (fq6 *fq6C) mulBy1(e *fe6C, c1 *fe2) {
	fq2, t := fq6.fq2(), fq6.t
	fq2.mul(t[0], &e[1], c1)
	fq2.add(t[1], &e[1], &e[2])
	fq2.mul(t[1], t[1], c1)
	fq2.sub(t[1], t[1], t[0])
	fq6.mulByNonResidue(t[1], t[1])
	fq2.add(t[2], &e[0], &e[1])
	fq2.mul(t[2], t[2], c1)
	fq2.sub(&e[1], t[2], t[0])
	fq2.copy(&e[0], t[1])
	fq2.copy(&e[2], t[0])
}

func (fq6 *fq6C) frobeniusMap(c, a *fe6C, power uint) {
	fq2 := fq6.fq2()
	fq6.copy(c, a)
	fq2.frobeniusMap(&c[0], &a[0], power)
	fq2.frobeniusMap(&c[1], &a[1], power)
	fq2.frobeniusMap(&c[2], &a[2], power)
	fq2.mul(&c[1], &c[1], fq6.frobeniusCoeffs[0][power%6])
	fq2.mul(&c[2], &c[2], fq6.frobeniusCoeffs[1][power%6])
}

func (fq6 *fq6C) calculateFrobeniusCoeffs() bool {
	fq2 := fq6.fq2()
	modulus := fq2.modulus()
	zero, one, three := big.NewInt(0), big.NewInt(1), big.NewInt(3)
	if fq6.frobeniusCoeffs == nil {
		fq6.frobeniusCoeffs = new([2][6]*fe2)
		for i := 0; i < 6; i++ {
			fq6.frobeniusCoeffs[0][i] = fq2.new()
			fq6.frobeniusCoeffs[1][i] = fq2.new()
		}
	}
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	fq2.copy(fq6.frobeniusCoeffs[0][0], fq2.one())
	fq2.copy(fq6.frobeniusCoeffs[1][0], fq2.one())
	for i := 1; i <= 3; i++ {
		power.Sub(qPower, one)
		power.DivMod(power, three, rem)
		if rem.Cmp(zero) != 0 {
			return false
		}
		fq2.exp(fq6.frobeniusCoeffs[0][i], fq6.nonResidue, power)
		fq2.square(fq6.frobeniusCoeffs[1][i], fq6.frobeniusCoeffs[0][i])
		qPower.Mul(qPower, modulus)
	}
	return true
}

func (fq6 *fq6C) calculateFrobeniusCoeffsWithPrecomputation(f1, f2 *fe2) bool {
	fq2 := fq6.fq2()
	if fq6.frobeniusCoeffs == nil {
		fq6.frobeniusCoeffs = new([2][6]*fe2)
		for i := 0; i < 6; i++ {
			fq6.frobeniusCoeffs[0][i] = fq2.new()
			fq6.frobeniusCoeffs[1][i] = fq2.new()
		}
	}
	fq2.copy(fq6.frobeniusCoeffs[0][0], fq2.one())
	fq2.copy(fq6.frobeniusCoeffs[1][0], fq2.one())
	fq2.square(fq6.frobeniusCoeffs[0][1], f1)
	fq2.square(fq6.frobeniusCoeffs[0][2], f2)
	fq2.frobeniusMap(fq6.frobeniusCoeffs[0][3], fq6.frobeniusCoeffs[0][2], 1)
	fq2.mul(fq6.frobeniusCoeffs[0][3], fq6.frobeniusCoeffs[0][3], fq6.frobeniusCoeffs[0][1])
	fq2.square(fq6.frobeniusCoeffs[1][1], fq6.frobeniusCoeffs[0][1])
	fq2.square(fq6.frobeniusCoeffs[1][2], fq6.frobeniusCoeffs[0][2])
	fq2.square(fq6.frobeniusCoeffs[1][3], fq6.frobeniusCoeffs[0][3])
	return true
}

func (fq6 *fq6C) fq() *fq {
	return fq6.f.f
}

func (fq6 *fq6C) fq2() *fq2 {
	return fq6.f
}
