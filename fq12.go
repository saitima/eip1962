package eip

import (
	"fmt"
	"io"
	"math/big"
)

type fe12 [2]fe6C

type fq12 struct {
	f               *fq6C
	nonResidue      *fe6C
	t               []*fe6C
	t2              []*fe2
	frobeniusCoeffs *[12]*fe2
}

func newFq12(fq6 *fq6C, nonResidueBuf []byte) (*fq12, error) {
	nonResidue := fq6.new()
	if nonResidueBuf != nil {
		var err error
		nonResidue, err = fq6.fromBytes(nonResidueBuf)
		if err != nil {
			return nil, err
		}
	}
	t := make([]*fe6C, 4)
	for i := 0; i < 4; i++ {
		t[i] = fq6.new()
	}
	fq2 := fq6.f
	t2 := make([]*fe2, 10)
	for i := 0; i < 10; i++ {
		t2[i] = fq2.new()
	}
	return &fq12{fq6, nonResidue, t, t2, nil}, nil
}

func (fq12 *fq12) new() *fe12 {
	return fq12.zero()
}

func (fq12 *fq12) modulus() *big.Int {
	fq6 := fq12.fq6()
	return fq6.modulus()
}

func (fq12 *fq12) byteSize() int {
	fq6 := fq12.fq6()
	return fq6.byteSize() * 2
}

func (fq12 *fq12) rand(r io.Reader) *fe12 {
	fq6 := fq12.fq6()
	return &fe12{*fq6.rand(r), *fq6.rand(r)}
}

func (fq12 *fq12) fromBytes(in []byte) (*fe12, error) {
	fq6 := fq12.fq6()
	byteSize := fq6.byteSize()
	if len(in) != byteSize*2 {
		return nil, fmt.Errorf("input string should be larger than %d bytes", byteSize*2)
	}
	var err error
	c0, err := fq6.fromBytes(in[:byteSize])
	if err != nil {
		return nil, err
	}
	c1, err := fq6.fromBytes(in[byteSize:])
	if err != nil {
		return nil, err
	}
	return &fe12{*c0, *c1}, nil
}

func (fq12 *fq12) toBytes(a *fe12) []byte {
	fq6 := fq12.fq6()
	byteSize := fq6.byteSize()
	out := make([]byte, 2*byteSize)
	copy(out[:byteSize], fq6.toBytes(&a[0]))
	copy(out[byteSize:], fq6.toBytes(&a[1]))
	return out
}

func (fq12 *fq12) toString(a *fe12) string {
	fq6 := fq12.fq6()
	return fmt.Sprintf("%s\n%s", fq6.toString(&a[0]), fq6.toString(&a[1]))
}

func (fq12 *fq12) toStringNoTransform(a *fe12) string {
	fq6 := fq12.fq6()
	return fmt.Sprintf("%s\n%s", fq6.toStringNoTransform(&a[0]), fq6.toStringNoTransform(&a[1]))
}

func (fq12 *fq12) zero() *fe12 {
	fq6 := fq12.fq6()
	return &fe12{*fq6.zero(), *fq6.zero()}
}

func (fq12 *fq12) one() *fe12 {
	fq6 := fq12.fq6()
	a := fq12.zero()
	fq6.copy(&a[0], fq6.one())
	return a
}

func (fq12 *fq12) isZero(a *fe12) bool {
	fq6 := fq12.fq6()
	return fq6.isZero(&a[0]) && fq6.isZero(&a[1])
}

func (fq12 *fq12) isOne(a *fe12) bool {
	fq6 := fq12.fq6()
	return fq6.isOne(&a[0]) && fq6.isZero(&a[1])
}

func (fq12 *fq12) equal(a, b *fe12) bool {
	fq6 := fq12.fq6()
	return fq6.equal(&a[0], &b[0]) && fq6.equal(&a[1], &b[1])
}

func (fq12 *fq12) copy(c, a *fe12) *fe12 {
	fq6 := fq12.fq6()
	fq6.copy(&c[0], &a[0])
	fq6.copy(&c[1], &a[1])
	return c
}

func (fq12 *fq12) add(c, a, b *fe12) *fe12 {
	fq6 := fq12.fq6()
	fq6.add(&c[0], &a[0], &b[0])
	fq6.add(&c[1], &a[1], &b[1])
	return c
}

func (fq12 *fq12) double(c, a *fe12) *fe12 {
	fq6 := fq12.fq6()
	fq6.double(&c[0], &a[0])
	fq6.double(&c[1], &a[1])
	return c
}

func (fq12 *fq12) sub(c, a, b *fe12) *fe12 {
	fq6 := fq12.fq6()
	fq6.sub(&c[0], &a[0], &b[0])
	fq6.sub(&c[1], &a[1], &b[1])
	return c
}

func (fq12 *fq12) neg(c, a *fe12) *fe12 {
	fq6 := fq12.fq6()
	fq6.neg(&c[0], &a[0])
	fq6.neg(&c[1], &a[1])
	return c
}

func (fq12 *fq12) conjugate(c, a *fe12) *fe12 {
	fq6 := fq12.fq6()
	fq12.copy(c, a)
	fq6.neg(&c[1], &a[1])
	return c
}

func (fq12 *fq12) mulByNonResidue(c, a *fe6C) {
	fq6, fq2 := fq12.fq6(), fq12.fq2()
	o := fq2.new()
	fq6.mulByNonResidue(o, &a[2])
	fq2.copy(&c[2], &a[1])
	fq2.copy(&c[1], &a[0])
	fq2.copy(&c[0], o)
}

func (fq12 *fq12) mul(c, a, b *fe12) {
	fq6, t := fq12.fq6(), fq12.t
	// c0 = (a0 * b0) + β * (a1 * b1)
	// c1 = (a0 + a1) * (b0 + b1) - (a0 * b0 + a1 * b1)
	fq6.mul(t[1], &a[0], &b[0])      // v0 = a0 * b0
	fq6.mul(t[2], &a[1], &b[1])      // v1 = a1 * b1
	fq6.add(t[0], t[1], t[2])        // v0 + v1
	fq12.mulByNonResidue(t[2], t[2]) // β * v1
	fq6.add(t[3], t[1], t[2])        // β * v1 + v0
	fq6.add(t[1], &a[0], &a[1])      // a0 + a1
	fq6.add(t[2], &b[0], &b[1])      // b0 + b1
	fq6.mul(t[1], t[1], t[2])        // (a0 + a1)(b0 + b1)
	fq6.copy(&c[0], t[3])            // c0 = β * v1 + v0
	fq6.sub(&c[1], t[1], t[0])       // c1 = (a0 + a1)(b0 + b1) - (v0+v1)
}

func (fq12 *fq12) square(c, a *fe12) {
	fq6, t := fq12.fq6(), fq12.t
	// c0 = (a0 - a1) * (a0 - β * a1) + a0 * a1 + β * a0 * a1
	// c1 = 2 * a0 * a1
	fq6.sub(t[0], &a[0], &a[1])       // v0 = a0 - a1
	fq12.mulByNonResidue(t[1], &a[1]) // β * a1
	fq6.sub(t[2], &a[0], t[1])        // v3 = a0 -  β * a1
	fq6.mul(t[1], &a[0], &a[1])       // v2 = a0 * a1
	fq6.mul(t[0], t[0], t[2])         // v0 * v3
	fq6.add(t[0], t[1], t[0])         // v0 = v0 * v3 + v2
	fq12.mulByNonResidue(t[2], t[1])  // β * v2
	fq6.add(&c[0], t[0], t[2])        // c0 = v0 + β * v2
	fq6.double(&c[1], t[1])           // c1 = 2*v2
}

func (fq12 *fq12) inverse(c, a *fe12) bool {
	fq6, t := fq12.fq6(), fq12.t
	// c0 = a0 * (a0^2 - β * a1^2)^-1
	// c1 = a1 * (a0^2 - β * a1^2)^-1
	fq6.square(t[0], &a[0])                 // v0 = a0^2
	fq6.square(t[1], &a[1])                 // v1 = a1^2
	fq12.mulByNonResidue(t[1], t[1])        // β * v1
	fq6.sub(t[1], t[0], t[1])               // v0 = v0 - β * v1
	if ok := fq6.inverse(t[0], t[1]); !ok { // v1 = v0^-1
		fq12.copy(c, fq12.zero())
		return false
	}
	fq6.mul(&c[0], &a[0], t[0]) // c0 = a0 * v1
	fq6.mul(t[0], &a[1], t[0])  // a1 * v1
	fq6.neg(&c[1], t[0])        // c1 = -a1 * v1
	return true
}

func (fq12 *fq12) exp(c, a *fe12, e *big.Int) {
	z := fq12.one()
	found := false
	for i := e.BitLen() - 1; i >= 0; i-- {
		if found {
			fq12.square(z, z)
		} else {
			found = e.Bit(i) == 1
		}
		if e.Bit(i) == 1 {
			fq12.mul(z, z, a)
		}
	}
	fq12.copy(c, z)
}

func (fq12 *fq12) fp4Square(c0, c1, a0, a1 *fe2) {
	t := make([]*fe2, 3) // fix
	fq2, fq6 := fq12.fq2(), fq12.fq6()
	t[0] = fq2.new()
	t[1] = fq2.new()
	t[2] = fq2.new()
	fq2.square(t[0], a0)
	fq2.square(t[1], a1)
	fq6.mulByNonResidue(t[2], t[1])
	fq2.add(c0, t[2], t[0])
	fq2.add(t[2], a0, a1)
	fq2.square(t[2], t[2])
	fq2.sub(t[2], t[2], t[0])
	fq2.sub(c1, t[2], t[1])
}

func (fq12 *fq12) cyclotomicSquare(c, a *fe12) {
	fq2, fq6, t := fq12.fq2(), fq12.fq6(), fq12.t2
	fq12.fp4Square(t[3], t[4], &a[0][0], &a[1][1])
	fq2.sub(t[2], t[3], &a[0][0])
	fq2.double(t[2], t[2])
	fq2.add(&c[0][0], t[2], t[3])
	fq2.add(t[2], t[4], &a[1][1])
	fq2.double(t[2], t[2])
	fq2.add(&c[1][1], t[2], t[4])
	fq12.fp4Square(t[3], t[4], &a[1][0], &a[0][2])
	fq12.fp4Square(t[5], t[6], &a[0][1], &a[1][2])
	fq2.sub(t[2], t[3], &a[0][1])
	fq2.double(t[2], t[2])
	fq2.add(&c[0][1], t[2], t[3])
	fq2.add(t[2], t[4], &a[1][2])
	fq2.double(t[2], t[2])
	fq2.add(&c[1][2], t[2], t[4])
	fq6.mulByNonResidue(t[3], t[6])
	fq2.add(t[2], t[3], &a[1][0])
	fq2.double(t[2], t[2])
	fq2.add(&c[1][0], t[2], t[3])
	fq2.sub(t[2], t[5], &a[0][2])
	fq2.double(t[2], t[2])
	fq2.add(&c[0][2], t[2], t[5])
}

func (fq12 *fq12) cyclotomicExp(c, a *fe12, e *big.Int) {
	z := fq12.one()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq12.cyclotomicSquare(z, z)
		if e.Bit(i) == 1 {
			fq12.mul(z, z, a)
		}
	}
	fq12.copy(c, z)
}

func (fq12 *fq12) mulBy034(a *fe12, c0, c3, c4 *fe2) {
	fq2, fq6, t := fq12.fq2(), fq12.fq6(), fq12.t
	o := fq2.new()
	fq2.mul(&t[0][0], &a[0][0], c0)
	fq2.mul(&t[0][1], &a[0][1], c0)
	fq2.mul(&t[0][2], &a[0][2], c0)
	fq6.copy(t[1], &a[1])
	fq6.mulBy01(t[1], c3, c4)
	fq2.add(o, c0, c3)
	fq6.add(t[2], &a[1], &a[0])
	fq6.mulBy01(t[2], o, c4)
	fq6.sub(t[2], t[2], t[0])
	fq6.sub(&a[1], t[2], t[1])
	fq12.mulByNonResidue(t[1], t[1])
	fq6.add(&a[0], t[0], t[1])
}

func (fq12 *fq12) mulBy014(e *fe12, c0, c1, c4 *fe2) {
	fq2, fq6, t := fq12.fq2(), fq12.fq6(), fq12.t
	o := fq2.new()
	fq6.copy(t[0], &e[0])
	fq6.mulBy01(t[0], c0, c1)
	fq6.copy(t[1], &e[1])
	fq6.mulBy1(t[1], c4)
	fq2.add(o, c1, c4)
	fq6.add(&e[1], &e[1], &e[0])
	fq6.mulBy01(&e[1], c0, o)
	fq6.sub(&e[1], &e[1], t[0])
	fq6.sub(&e[1], &e[1], t[1])
	fq12.mulByNonResidue(t[1], t[1])
	fq6.add(&e[0], t[1], t[0])
}

func (fq12 *fq12) frobeniusMap(c, a *fe12, power uint) {
	fq2, fq6 := fq12.fq2(), fq12.fq6()
	fq6.frobeniusMap(&c[0], &a[0], power)
	fq6.frobeniusMap(&c[1], &a[1], power)
	fq2.mul(&c[1][0], &c[1][0], fq12.frobeniusCoeffs[power%12])
	fq2.mul(&c[1][1], &c[1][1], fq12.frobeniusCoeffs[power%12])
	fq2.mul(&c[1][2], &c[1][2], fq12.frobeniusCoeffs[power%12])
}

func (fq12 *fq12) calculateFrobeniusCoeffs() bool {
	fq6, fq2 := fq12.fq6(), fq12.fq2()
	if fq12.frobeniusCoeffs == nil {
		fq12.frobeniusCoeffs = new([12]*fe2)
		for i := 0; i < 12; i++ {
			fq12.frobeniusCoeffs[i] = fq2.new()
		}
	}
	modulus := fq2.modulus()
	fq2.copy(fq12.frobeniusCoeffs[0], fq2.one())
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	zero := new(big.Int)
	for i := 1; i <= 6; i++ {
		if i > 3 && i < 6 {
			continue
		}
		power := power.Sub(qPower, big.NewInt(1))
		power, rem = power.DivMod(power, big.NewInt(6), rem)
		if rem.Cmp(zero) != 0 {

			return false
		}
		fq2.exp(fq12.frobeniusCoeffs[i], fq6.nonResidue, power)
		if i == 3 {
			qPower.Mul(qPower, qPower)
		} else {
			qPower.Mul(qPower, modulus)
		}
	}
	return true
}

func (fq12 *fq12) calculateFrobeniusCoeffsWithPrecomputation(f1, f2 *fe2) bool {
	fq2 := fq12.fq2()
	if fq12.frobeniusCoeffs == nil {
		fq12.frobeniusCoeffs = new([12]*fe2)
		for i := 0; i < 12; i++ {
			fq12.frobeniusCoeffs[i] = fq2.new()
		}
	}
	fq2.copy(fq12.frobeniusCoeffs[0], fq2.one())
	fq2.copy(fq12.frobeniusCoeffs[1], f1)
	fq2.copy(fq12.frobeniusCoeffs[2], f2)
	fq2.frobeniusMap(fq12.frobeniusCoeffs[3], f2, 1)
	fq2.mul(fq12.frobeniusCoeffs[3], fq12.frobeniusCoeffs[3], fq12.frobeniusCoeffs[1])
	fq2.exp(fq12.frobeniusCoeffs[6], fq12.frobeniusCoeffs[2], big.NewInt(3))
	return true
}

func (fq12 *fq12) fq() *fq {
	return fq12.f.f.f
}

func (fq12 *fq12) fq2() *fq2 {
	return fq12.f.f
}

func (fq12 *fq12) fq6() *fq6C {
	return fq12.f
}
