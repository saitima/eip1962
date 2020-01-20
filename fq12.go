package eip

import (
	"fmt"
	"math/big"
)

type fe12 [2]fe6

type fq12 struct {
	f               *fq6
	nonResidue      *fe6
	t               []*fe6
	t2              []*fe2
	frobeniusCoeffs *[12]*fe2
}

func newFq12(f *fq6, nonResidue []byte) (*fq12, error) {
	nonResidue_ := f.newElement()
	if nonResidue != nil {
		var err error
		nonResidue_, err = f.fromBytes(nonResidue)
		if err != nil {
			return nil, err
		}
	}
	t := make([]*fe6, 4)
	for i := 0; i < 4; i++ {
		t[i] = f.zero()
	}
	t2 := make([]*fe2, 10)
	for i := 0; i < 10; i++ {
		t2[i] = f.f.zero()
	}

	return &fq12{f, nonResidue_, t, t2, nil}, nil
}

func (fq *fq12) newElement() *fe12 {
	return fq.zero()
}

func (fq *fq12) fromBytes(in []byte) (*fe12, error) {
	byteLen := fq.f.f.f.limbSize * 8 * 2 * 3
	if len(in) < len(&fe12{})*byteLen {
		return nil, fmt.Errorf("input string should be larger than 64 bytes")
	}
	c := fq.newElement()
	for i := 0; i < len(&fe12{}); i++ {
		elem, err := fq.f.fromBytes(in[i*byteLen : (i+1)*byteLen])
		if err != nil {
			return nil, err
		}
		c[i] = *elem
	}
	return c, nil
}

func (fq *fq12) toBytes(a *fe12) []byte {
	byteLen := fq.f.f.f.limbSize * 8 * 2 * 3
	out := make([]byte, len(a)*byteLen)
	for i := 0; i < len(a); i++ {
		copy(out[i*byteLen:(i+1)*byteLen], fq.f.toBytes(&a[i]))
	}
	return out
}

func (fq *fq12) toString(a *fe12) string {
	return fmt.Sprintf(
		"c0: %s c1: %s\n",
		fq.f.toString(&a[0]),
		fq.f.toString(&a[1]),
	)
}

func (fq *fq12) toStringNoTransform(a *fe12) string {
	return fmt.Sprintf(
		"c0: %s c1: %s\n",
		fq.f.toStringNoTransform(&a[0]),
		fq.f.toStringNoTransform(&a[1]),
	)
}

func (fq *fq12) zero() *fe12 {
	return &fe12{
		*fq.f.zero(),
		*fq.f.zero(),
	}
}

func (fq *fq12) one() *fe12 {
	a := fq.zero()
	fq.f.copy(&a[0], fq.f.one())
	return a
}

func (fq *fq12) isZero(a *fe12) bool {
	return fq.f.isZero(&a[0]) && fq.f.isZero(&a[1])
}

func (fq *fq12) equal(a, b *fe12) bool {
	return fq.f.equal(&a[0], &b[0]) && fq.f.equal(&a[1], &b[1])
}

func (fq *fq12) copy(c, a *fe12) *fe12 {
	fq.f.copy(&c[0], &a[0])
	fq.f.copy(&c[1], &a[1])
	return c
}

func (fq *fq12) add(c, a, b *fe12) *fe12 {
	fq.f.add(&c[0], &a[0], &b[0])
	fq.f.add(&c[1], &a[1], &b[1])
	return c
}

func (fq *fq12) double(c, a *fe12) *fe12 {
	fq.f.double(&c[0], &a[0])
	fq.f.double(&c[1], &a[1])
	return c
}

func (fq *fq12) sub(c, a, b *fe12) *fe12 {
	fq.f.sub(&c[0], &a[0], &b[0])
	fq.f.sub(&c[1], &a[1], &b[1])
	return c
}

func (fq *fq12) neg(c, a *fe12) *fe12 {
	fq.f.neg(&c[0], &a[0])
	fq.f.neg(&c[1], &a[1])
	return c
}

func (fq *fq12) conjugate(c, a *fe12) *fe12 {
	fq.copy(c, a)
	fq.f.neg(&c[1], &a[1])
	return c
}

func (fq *fq12) mulByNonResidue(c, a *fe6) {
	o := fq.f.f.newElement()
	fq.f.mulByNonResidue(o, &a[2])
	fq.f.f.copy(&c[2], &a[1])
	fq.f.f.copy(&c[1], &a[0])
	fq.f.f.copy(&c[0], o)
}

func (fq *fq12) mul(c, a, b *fe12) {
	t := fq.t
	// c0 = (a0 * b0) + β * (a1 * b1)
	// c1 = (a0 + a1) * (b0 + b1) - (a0 * b0 + a1 * b1)
	fq.f.mul(t[1], &a[0], &b[0])   // v0 = a0 * b0
	fq.f.mul(t[2], &a[1], &b[1])   // v1 = a1 * b1
	fq.f.add(t[0], t[1], t[2])     // v0 + v1
	fq.mulByNonResidue(t[2], t[2]) // β * v1
	fq.f.add(t[3], t[1], t[2])     // β * v1 + v0
	fq.f.add(t[1], &a[0], &a[1])   // a0 + a1
	fq.f.add(t[2], &b[0], &b[1])   // b0 + b1
	fq.f.mul(t[1], t[1], t[2])     // (a0 + a1)(b0 + b1)
	fq.f.copy(&c[0], t[3])         // c0 = β * v1 + v0
	fq.f.sub(&c[1], t[1], t[0])    // c1 = (a0 + a1)(b0 + b1) - (v0+v1)
}

func (fq *fq12) square(c, a *fe12) {
	t := fq.t
	// c0 = (a0 - a1) * (a0 - β * a1) + a0 * a1 + β * a0 * a1
	// c1 = 2 * a0 * a1
	fq.f.sub(t[0], &a[0], &a[1])    // v0 = a0 - a1
	fq.mulByNonResidue(t[1], &a[1]) // β * a1
	fq.f.sub(t[2], &a[0], t[1])     // v3 = a0 -  β * a1
	fq.f.mul(t[1], &a[0], &a[1])    // v2 = a0 * a1
	fq.f.mul(t[0], t[0], t[2])      // v0 * v3
	fq.f.add(t[0], t[1], t[0])      // v0 = v0 * v3 + v2
	fq.mulByNonResidue(t[2], t[1])  // β * v2
	fq.f.add(&c[0], t[0], t[2])     // c0 = v0 + β * v2
	fq.f.double(&c[1], t[1])        // c1 = 2*v2
}

func (fq *fq12) inverse(c, a *fe12) {
	t := fq.t
	// c0 = a0 * (a0^2 - β * a1^2)^-1
	// c1 = a1 * (a0^2 - β * a1^2)^-1
	fq.f.square(t[0], &a[0])       // v0 = a0^2
	fq.f.square(t[1], &a[1])       // v1 = a1^2
	fq.mulByNonResidue(t[1], t[1]) // β * v1
	fq.f.sub(t[1], t[0], t[1])     // v0 = v0 - β * v1
	fq.f.inverse(t[0], t[1])       // v1 = v0^-1
	fq.f.mul(&c[0], &a[0], t[0])   // c0 = a0 * v1
	fq.f.mul(t[0], &a[1], t[0])    // a1 * v1
	fq.f.neg(&c[1], t[0])          // c1 = -a1 * v1
}

func (fq *fq12) exp(c, a *fe12, e *big.Int) {
	z := fq.one()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.square(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.copy(c, z)
}

func (fq *fq12) fp4Square(c0, c1, a0, a1 *fe2) {
	t := make([]*fe2, 3)
	fq2 := fq.f.f
	t[0] = fq2.newElement()
	t[1] = fq2.newElement()
	t[2] = fq2.newElement()
	fq2.square(t[0], a0)
	fq2.square(t[1], a1)
	fq2.mulByNonResidue12(t[2], t[1])
	// fp2.mulByFq(t[2], t[1], nonResidue1)
	fq2.add(c0, t[2], t[0])
	fq2.add(t[2], a0, a1)
	fq2.square(t[2], t[2])
	fq2.sub(t[2], t[2], t[0])
	fq2.sub(c1, t[2], t[1])
}

func (fq *fq12) cyclotomicSquare(c, a *fe12) {
	t := fq.t2
	fp2 := fq.f.f
	fq.fp4Square(t[3], t[4], &a[0][0], &a[1][1])
	fp2.sub(t[2], t[3], &a[0][0])
	fp2.double(t[2], t[2])
	fp2.add(&c[0][0], t[2], t[3])
	fp2.add(t[2], t[4], &a[1][1])
	fp2.double(t[2], t[2])
	fp2.add(&c[1][1], t[2], t[4])
	fq.fp4Square(t[3], t[4], &a[1][0], &a[0][2])
	fq.fp4Square(t[5], t[6], &a[0][1], &a[1][2])
	fp2.sub(t[2], t[3], &a[0][1])
	fp2.double(t[2], t[2])
	fp2.add(&c[0][1], t[2], t[3])
	fp2.add(t[2], t[4], &a[1][2])
	fp2.double(t[2], t[2])
	fp2.add(&c[1][2], t[2], t[4])
	fp2.mulByNonResidue12(t[3], t[6])
	fp2.add(t[2], t[3], &a[1][0])
	fp2.double(t[2], t[2])
	fp2.add(&c[1][0], t[2], t[3])
	fp2.sub(t[2], t[5], &a[0][2])
	fp2.double(t[2], t[2])
	fp2.add(&c[0][2], t[2], t[5])
}

func (fq *fq12) cyclotomicSquare2(c, a *fe12) {
	t := new([9]*fe2)
	for i := 0; i < 9; i++ {
		t[i] = fq.f.f.newElement()
	}
	fq.f.f.mul(t[0], &a[0][0], &a[1][1])
	fq.f.f.add(t[1], &a[0][0], &a[1][1])
	fq.f.mulByNonResidue(t[2], &a[1][1])
	fq.f.f.add(t[2], t[2], &a[0][0])
	fq.f.mulByNonResidue(t[3], t[0])
	fq.f.f.mul(t[4], t[1], t[2])
	fq.f.f.sub(t[4], t[4], t[0])
	fq.f.f.sub(t[4], t[4], t[3])
	fq.f.f.double(t[5], t[0])
	fq.f.f.mul(t[0], &a[1][0], &a[0][2])
	fq.f.f.add(t[1], &a[1][0], &a[0][2])
	fq.f.mulByNonResidue(t[2], &a[0][2])
	fq.f.f.add(t[2], t[2], &a[1][0])
	fq.f.mulByNonResidue(t[3], t[0])
	fq.f.f.mul(t[6], t[1], t[2])
	fq.f.f.sub(t[6], t[6], t[0])
	fq.f.f.sub(t[6], t[6], t[3])
	fq.f.f.double(t[7], t[0])
	fq.f.f.mul(t[0], &a[0][1], &a[1][2])
	fq.f.f.add(t[1], &a[0][1], &a[1][2])
	fq.f.mulByNonResidue(t[2], &a[1][2])
	fq.f.f.add(t[2], t[2], &a[0][1])
	fq.f.mulByNonResidue(t[3], t[0])
	fq.f.f.mul(t[8], t[1], t[2])
	fq.f.f.sub(t[8], t[8], t[0])
	fq.f.f.sub(t[8], t[8], t[3])
	fq.f.f.double(t[0], t[0])
	fq.f.mulByNonResidue(t[0], t[0])
	fq.f.f.sub(t[1], t[4], &a[0][0])
	fq.f.f.double(t[1], t[1])
	fq.f.f.add(t[1], t[1], t[4])
	fq.f.f.copy(&c[0][0], t[1])
	fq.f.f.add(t[1], t[5], &a[1][1])
	fq.f.f.double(t[1], t[1])
	fq.f.f.add(t[1], t[1], t[5])
	fq.f.f.copy(&c[1][1], t[1])
	fq.f.f.add(t[1], t[0], &a[1][0])
	fq.f.f.double(t[1], t[1])
	fq.f.f.add(t[1], t[1], t[0])
	fq.f.f.copy(&c[1][0], t[1])
	fq.f.f.sub(t[1], t[8], &a[0][2])
	fq.f.f.double(t[1], t[1])
	fq.f.f.add(t[1], t[1], t[8])
	fq.f.f.copy(&c[0][2], t[1])
	fq.f.f.sub(t[1], t[6], &a[0][1])
	fq.f.f.double(t[1], t[1])
	fq.f.f.add(t[1], t[1], t[6])
	fq.f.f.copy(&c[0][1], t[1])
	fq.f.f.add(t[1], t[7], &a[1][2])
	fq.f.f.double(t[1], t[1])
	fq.f.f.add(t[1], t[1], t[7])
	fq.f.f.copy(&c[1][2], t[1])
}

func (fq *fq12) cyclotomicExp(c, a *fe12, e *big.Int) {
	z := fq.one()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.cyclotomicSquare(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.copy(c, z)
}

func (fq *fq12) cyclotomicExp2(c, a *fe12, e *big.Int) {
	z := fq.one()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.cyclotomicSquare2(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.copy(c, z)
}

func (fq *fq12) mulBy034(f *fe12, c0, c3, c4 *fe2) {
	o := fq.f.f.newElement() // base filed temp mem could be used :/
	t := fq.t
	fq.f.f.mul(&t[0][0], &f[0][0], c0)
	fq.f.f.mul(&t[0][1], &f[0][1], c0)
	fq.f.f.mul(&t[0][2], &f[0][2], c0)
	fq.f.copy(t[1], &f[1])
	fq.f.mulBy01(t[1], c3, c4)
	fq.f.f.add(o, c0, c3)
	fq.f.add(t[2], &f[1], &f[0])
	fq.f.mulBy01(t[2], o, c4)
	fq.f.sub(t[2], t[2], t[0])
	fq.f.sub(&f[1], t[2], t[1])
	fq.mulByNonResidue(t[1], t[1])
	fq.f.add(&f[0], t[0], t[1])
}

func (fq *fq12) mulBy014(e *fe12, c0, c1, c4 *fe2) {
	o := fq.f.f.newElement() // base field temp mem could be used :/
	t := fq.t
	fq.f.copy(t[0], &e[0])
	fq.f.mulBy01(t[0], c0, c1)
	fq.f.copy(t[1], &e[1])
	fq.f.mulBy1(t[1], c4)
	fq.f.f.add(o, c1, c4)
	fq.f.add(&e[1], &e[1], &e[0])
	fq.f.mulBy01(&e[1], c0, o)
	fq.f.sub(&e[1], &e[1], t[0])
	fq.f.sub(&e[1], &e[1], t[1])
	fq.mulByNonResidue(t[1], t[1])
	fq.f.add(&e[0], t[1], t[0])
}

func (fq *fq12) frobeniusMap(c, a *fe12, power uint) {
	fq.f.frobeniusMap(&c[0], &a[0], power)
	fq.f.frobeniusMap(&c[1], &a[1], power)
	fq.f.f.mul(&c[1][0], &c[1][0], fq.frobeniusCoeffs[power%12])
	fq.f.f.mul(&c[1][1], &c[1][1], fq.frobeniusCoeffs[power%12])
	fq.f.f.mul(&c[1][2], &c[1][2], fq.frobeniusCoeffs[power%12])
}

func (fq *fq12) calculateFrobeniusCoeffs() bool {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = new([12]*fe2)
		for i := 0; i < 12; i++ {
			fq.frobeniusCoeffs[i] = fq.f.f.newElement()
		}
	}
	modulus := fq.f.f.f.pbig
	f0 := fq.f.f.one()
	fq.f.f.copy(fq.frobeniusCoeffs[0], f0)
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	for i := 1; i <= 6; i++ {
		if i > 3 && i < 6 {
			continue
		}
		power := power.Sub(qPower, big.NewInt(1))
		power, rem = power.DivMod(power, big.NewInt(6), rem)
		if rem.Uint64() != 0 {
			return false
		}
		fq.f.f.exp(fq.frobeniusCoeffs[i], fq.f.nonResidue, power)
		if i == 3 {
			qPower.Mul(qPower, qPower)
		} else {
			qPower.Mul(qPower, modulus)
		}
	}
	return true
}

func (fq *fq12) calculateFrobeniusCoeffsWithPrecomputation(f1, f2 *fe2) bool {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = new([12]*fe2)
		for i := 0; i < 12; i++ {
			fq.frobeniusCoeffs[i] = fq.f.f.newElement()
		}
	}
	fq.f.f.copy(fq.frobeniusCoeffs[0], fq.f.f.one())
	fq.f.f.copy(fq.frobeniusCoeffs[1], f1)
	fq.f.f.copy(fq.frobeniusCoeffs[2], f2)
	fq.f.f.frobeniusMap(fq.frobeniusCoeffs[3], f2, 1)
	fq.f.f.mul(fq.frobeniusCoeffs[3], fq.frobeniusCoeffs[3], fq.frobeniusCoeffs[1])
	fq.f.f.exp(fq.frobeniusCoeffs[6], fq.frobeniusCoeffs[2], big.NewInt(3))
	return true
}
