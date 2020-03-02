package eip

import (
	"fmt"
	"io"
	"math/big"
)

type fe4 [2]fe2

type fq4 struct {
	f               *fq2
	nonResidue      *fe2
	t               []*fe2
	frobeniusCoeffs *[4]fe
}

func newFq4(fq2 *fq2, nonResidueBuf []byte) (*fq4, error) {
	nonResidue := fq2.new()
	if nonResidueBuf != nil {
		var err error
		nonResidue, err = fq2.fromBytes(nonResidueBuf)
		if err != nil {
			return nil, err
		}
	}
	t := make([]*fe2, 4)
	for i := 0; i < 4; i++ {
		t[i] = fq2.new()
	}
	return &fq4{fq2, nonResidue, t, nil}, nil
}

func (fq4 *fq4) byteSize() int {
	fq2 := fq4.fq2()
	return fq2.byteSize() * 4
}

func (fq4 *fq4) new() *fe4 {
	return fq4.zero()
}

func (fq4 *fq4) modulus() *big.Int {
	fq2 := fq4.fq2()
	return fq2.modulus()
}

func (fq4 *fq4) rand(r io.Reader) *fe4 {
	fq2 := fq4.fq2()
	return &fe4{*fq2.rand(r), *fq2.rand(r)}
}

func (fq4 *fq4) fromBytes(in []byte) (*fe4, error) {
	fq2 := fq4.fq2()
	byteSize := fq2.byteSize()
	if len(in) != byteSize*2 {
		return nil, fmt.Errorf("input string should be larger than %d bytes", byteSize*2)
	}
	var err error
	u0, err := fq2.fromBytes(in[:byteSize])
	if err != nil {
		return nil, err
	}
	u1, err := fq2.fromBytes(in[byteSize:])
	if err != nil {
		return nil, err
	}
	return &fe4{*u0, *u1}, nil
}

func (fq4 *fq4) toBytes(a *fe4) []byte {
	fq2 := fq4.fq2()
	byteSize := fq2.byteSize()
	out := make([]byte, 2*byteSize)
	copy(out[:byteSize], fq2.toBytes(&a[0]))
	copy(out[byteSize:], fq2.toBytes(&a[1]))
	return out
}

func (fq4 *fq4) toString(a *fe4) string {
	fq2 := fq4.fq2()
	return fmt.Sprintf("%s\n%s", fq2.toString(&a[0]), fq2.toString(&a[1]))
}

func (fq4 *fq4) toStringNoTransform(a *fe4) string {
	fq2 := fq4.fq2()
	return fmt.Sprintf("%s\n%s", fq2.toStringNoTransform(&a[0]), fq2.toStringNoTransform(&a[1]))
}

func (fq4 *fq4) zero() *fe4 {
	fq2 := fq4.fq2()
	return &fe4{*fq2.zero(), *fq2.zero()}
}

func (fq4 *fq4) one() *fe4 {
	fq2 := fq4.fq2()
	a := fq4.new()
	fq2.copy(&a[0], fq2.one())
	return a
}

func (fq4 *fq4) isZero(a *fe4) bool {
	fq2 := fq4.fq2()
	return fq2.isZero(&a[0]) && fq2.isZero(&a[1])
}

func (fq4 *fq4) isOne(a *fe4) bool {
	fq2 := fq4.fq2()
	return fq2.isOne(&a[0]) && fq2.isZero(&a[1])
}

func (fq4 *fq4) equal(a, b *fe4) bool {
	fq2 := fq4.fq2()
	return fq2.equal(&a[0], &b[0]) && fq2.equal(&a[1], &b[1])
}

func (fq4 *fq4) copy(c, a *fe4) *fe4 {
	fq2 := fq4.fq2()
	fq2.copy(&c[0], &a[0])
	fq2.copy(&c[1], &a[1])
	return c
}

func (fq4 *fq4) add(c, a, b *fe4) *fe4 {
	fq2 := fq4.fq2()
	fq2.add(&c[0], &a[0], &b[0])
	fq2.add(&c[1], &a[1], &b[1])
	return c
}

func (fq4 *fq4) double(c, a *fe4) *fe4 {
	fq2 := fq4.fq2()
	fq2.double(&c[0], &a[0])
	fq2.double(&c[1], &a[1])
	return c
}

func (fq4 *fq4) sub(c, a, b *fe4) *fe4 {
	fq2 := fq4.fq2()
	fq2.sub(&c[0], &a[0], &b[0])
	fq2.sub(&c[1], &a[1], &b[1])
	return c
}

func (fq4 *fq4) neg(c, a *fe4) *fe4 {
	fq2 := fq4.fq2()
	fq2.neg(&c[0], &a[0])
	fq2.neg(&c[1], &a[1])
	return c
}

func (fq4 *fq4) conjugate(c, a *fe4) *fe4 {
	fq2 := fq4.fq2()
	fq4.copy(c, a)
	fq2.neg(&c[1], &a[1])
	return c
}

func (fq4 *fq4) mulByNonResidue(c, a *fe2) {
	fq, fq2 := fq4.fq(), fq4.fq2()
	o := fq2.new()
	fq.copy(o[1], a[0])
	fq2.mulByNonResidue(o[0], a[1])
	fq2.copy(c, o)
}

func (fq4 *fq4) mul(c, a, b *fe4) {
	fq2, t := fq4.fq2(), fq4.t
	// c0 = (a0 * b0) + β * (a1 * b1)
	// c1 = (a0 + a1) * (b0 + b1) - (a0 * b0 + a1 * b1)
	fq2.mul(t[1], &a[0], &b[0])     // v0 = a0 * b0
	fq2.mul(t[2], &a[1], &b[1])     // v1 = a1 * b1
	fq2.add(t[0], t[1], t[2])       // v0 + v1
	fq4.mulByNonResidue(t[2], t[2]) // β * v1
	fq2.add(t[3], t[1], t[2])       // β * v1 + v0
	fq2.add(t[1], &a[0], &a[1])     // a0 + a1
	fq2.add(t[2], &b[0], &b[1])     // b0 + b1
	fq2.mul(t[1], t[1], t[2])       // (a0 + a1)(b0 + b1)
	fq2.copy(&c[0], t[3])           // c0 = β * v1 + v0
	fq2.sub(&c[1], t[1], t[0])      // c1 = (a0 + a1)(b0 + b1) - (v0+v1)
}

func (fq4 *fq4) square(c, a *fe4) {
	fq2, t := fq4.fq2(), fq4.t
	// c0 = (a0 - a1) * (a0 - β * a1) + a0 * a1 + β * a0 * a1
	// c1 = 2 * a0 * a1
	fq2.sub(t[0], &a[0], &a[1])      // v0 = a0 - a1
	fq4.mulByNonResidue(t[1], &a[1]) // β * a1
	fq2.sub(t[2], &a[0], t[1])       // v3 = a0 -  β * a1
	fq2.mul(t[1], &a[0], &a[1])      // v2 = a0 * a1
	fq2.mul(t[0], t[0], t[2])        // v0 * v3
	fq2.add(t[0], t[1], t[0])        // v0 = v0 * v3 + v2
	fq4.mulByNonResidue(t[2], t[1])  // β * v2
	fq2.add(&c[0], t[0], t[2])       // c0 = v0 + β * v2
	fq2.double(&c[1], t[1])          // c1 = 2*v2
}

func (fq4 *fq4) inverse(c, a *fe4) bool {
	fq2, t := fq4.fq2(), fq4.t
	// c0 = a0 * (a0^2 - β * a1^2)^-1
	// c1 = a1 * (a0^2 - β * a1^2)^-1
	fq2.square(t[0], &a[0])         // v0 = a0^2
	fq2.square(t[1], &a[1])         // v1 = a1^2
	fq4.mulByNonResidue(t[1], t[1]) // β * v1
	fq2.sub(t[1], t[0], t[1])       // v0 = v0 - β * v1
	if ok := fq2.inverse(t[0], t[1]); !ok {
		fq4.copy(c, fq4.zero())
		return false
	} // v1 = v0^-1
	fq2.mul(&c[0], &a[0], t[0]) // c0 = a0 * v1
	fq2.mul(t[0], &a[1], t[0])  // a1 * v1
	fq2.neg(&c[1], t[0])        // c1 = -a1 * v1
	return true
}

func (fq4 *fq4) exp(c, a *fe4, e *big.Int) {
	z := fq4.one()
	found := false
	for i := e.BitLen() - 1; i >= 0; i-- {
		if found {
			fq4.square(z, z)
		} else {
			found = e.Bit(i) == 1
		}
		if e.Bit(i) == 1 {
			fq4.mul(z, z, a)
		}
	}
	fq4.copy(c, z)
}

func (fq4 *fq4) frobeniusMap(c, a *fe4, power uint) {
	fq2 := fq4.fq2()
	fq2.copy(&c[0], &a[0])
	fq2.frobeniusMap(&c[0], &a[0], power)
	fq2.frobeniusMap(&c[1], &a[1], power)
	fq2.mulByFq(&c[1], &c[1], fq4.frobeniusCoeffs[power%4])
}

func (fq4 *fq4) calculateFrobeniusCoeffs() bool {
	fq := fq4.fq()
	fq2 := fq4.fq2()
	modulus := fq.modulus()
	zero, one, four := big.NewInt(0), big.NewInt(1), big.NewInt(4)
	if fq4.frobeniusCoeffs == nil {
		fq4.frobeniusCoeffs = new([4]fe)
		for i := 0; i < 4; i++ {
			fq4.frobeniusCoeffs[i] = fq.new()
		}
	}
	f0 := fq.one
	fq.copy(fq4.frobeniusCoeffs[0], f0)
	fq.copy(fq4.frobeniusCoeffs[3], fq.zero)
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	for i := 1; i <= 2; i++ {
		power.Sub(qPower, one)
		power.DivMod(power, four, rem)
		if rem.Cmp(zero) != 0 {
			return false
		}
		fq.exp(fq4.frobeniusCoeffs[i], fq2.nonResidue, power)
		qPower.Mul(qPower, modulus)
	}
	return true
}

func (fq4 *fq4) calculateFrobeniusCoeffsWithPrecomputation(f1 fe) {
	fq := fq4.fq()
	if fq4.frobeniusCoeffs == nil {
		fq4.frobeniusCoeffs = new([4]fe)
		for i := 0; i < 4; i++ {
			fq4.frobeniusCoeffs[i] = fq.new()
		}
	}
	fq.copy(fq4.frobeniusCoeffs[0], fq.one)
	fq.copy(fq4.frobeniusCoeffs[1], f1)
	fq.square(fq4.frobeniusCoeffs[2], fq4.frobeniusCoeffs[1])
}

func (fq4 *fq4) fq() *fq {
	return fq4.f.f
}

func (fq4 *fq4) fq2() *fq2 {
	return fq4.f
}
