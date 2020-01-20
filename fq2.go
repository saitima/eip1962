package eip

import (
	"fmt"
	"math/big"
)

type fe2 [2]fieldElement

type fq2 struct {
	f               *field
	nonResidue      fieldElement
	t               []fieldElement
	frobeniusCoeffs *fe2
}

func newFq2(f *field, nonResidue []byte) (*fq2, error) {
	nonResidue_ := f.newFieldElement()
	if nonResidue != nil {
		var err error
		nonResidue_, err = f.newFieldElementFromBytes(nonResidue)
		if err != nil {
			return nil, err
		}
	}
	t := make([]fieldElement, 4)
	for i := 0; i < 4; i++ {
		t[i] = f.newFieldElement()
		f.copy(t[i], f.zero)
	}
	return &fq2{f, nonResidue_, t, nil}, nil
}

func (fq *fq2) newElement() *fe2 {
	fe := &fe2{fq.f.newFieldElement(), fq.f.newFieldElement()}
	fq.f.copy(fe[0], fq.f.zero)
	fq.f.copy(fe[1], fq.f.zero)
	return fe
}

func (fq *fq2) fromBytes(in []byte) (*fe2, error) {
	byteLen := fq.f.limbSize * 8
	if len(in) < len(&fe2{})*byteLen {
		return nil, fmt.Errorf("input string should be larger than 64 bytes")
	}
	c := fq.newElement()
	var err error
	for i := 0; i < len(&fe2{}); i++ {
		c[i], err = fq.f.newFieldElementFromBytes(in[i*byteLen : (i+1)*byteLen])
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (fq *fq2) toBytes(a *fe2) []byte {
	byteLen := fq.f.limbSize * 8
	out := make([]byte, len(a)*byteLen)
	for i := 0; i < len(a); i++ {
		copy(out[i*byteLen:(i+1)*byteLen], fq.f.toBytes(a[i]))
	}
	return out
}

func (fq *fq2) toString(a *fe2) string {
	return fmt.Sprintf(
		"c0: %s c1: %s\n",
		fq.f.toString(a[0]),
		fq.f.toString(a[1]),
	)
}

func (fq *fq2) toStringNoTransform(a *fe2) string {
	return fmt.Sprintf(
		"c0: %s c1: %s\n",
		fq.f.toStringNoTransform(a[0]),
		fq.f.toStringNoTransform(a[1]),
	)
}

func (fq *fq2) zero() *fe2 {
	return fq.newElement()
}

func (fq *fq2) one() *fe2 {
	a := fq.newElement()
	fq.f.copy(a[0], fq.f.one)
	return a
}

func (fq *fq2) isZero(a *fe2) bool {
	return fq.f.isZero(a[0]) && fq.f.isZero(a[1])
}

func (fq *fq2) equal(a, b *fe2) bool {
	return fq.f.equal(a[0], b[0]) && fq.f.equal(a[1], b[1])
}

func (fq *fq2) copy(c, a *fe2) *fe2 {
	fq.f.copy(c[0], a[0])
	fq.f.copy(c[1], a[1])
	return c
}

func (fq *fq2) add(c, a, b *fe2) *fe2 {
	fq.f.add(c[0], a[0], b[0])
	fq.f.add(c[1], a[1], b[1])
	return c
}

func (fq *fq2) double(c, a *fe2) *fe2 {
	fq.f.double(c[0], a[0])
	fq.f.double(c[1], a[1])
	return c
}

func (fq *fq2) sub(c, a, b *fe2) *fe2 {
	fq.f.sub(c[0], a[0], b[0])
	fq.f.sub(c[1], a[1], b[1])
	return c
}

func (fq *fq2) neg(c, a *fe2) *fe2 {
	fq.f.neg(c[0], a[0])
	fq.f.neg(c[1], a[1])
	return c
}

func (fq *fq2) conjugate(c, a *fe2) *fe2 {
	fq.copy(c, a)
	fq.f.neg(c[1], a[1])
	return c
}

func (fq *fq2) mulByNonResidue(c, a fieldElement) {
	fq.f.mul(c, a, fq.nonResidue)
}

func (fq *fq2) mulByNonResidue12(c, a *fe2) {
	t := fq.t
	fq.f.sub(t[0], a[0], a[1])
	fq.f.add(c[1], a[0], a[1])
	fq.f.copy(c[0], t[0])
}

func (fq *fq2) mul(c, a, b *fe2) {
	t := fq.t
	// c0 = (a0 * b0) + β * (a1 * b1)
	// c1 = (a0 + a1) * (b0 + b1) - (a0 * b0 + a1 * b1)
	fq.f.mul(t[1], a[0], b[0])     // v0 = a0 * b0
	fq.f.mul(t[2], a[1], b[1])     // v1 = a1 * b1
	fq.f.add(t[0], t[1], t[2])     // v0 + v1
	fq.mulByNonResidue(t[2], t[2]) // β * v1
	fq.f.add(t[3], t[1], t[2])     // β * v1 + v0
	fq.f.add(t[1], a[0], a[1])     // a0 + a1
	fq.f.add(t[2], b[0], b[1])     // b0 + b1
	fq.f.mul(t[1], t[1], t[2])     // (a0 + a1)(b0 + b1)
	fq.f.copy(c[0], t[3])          // c0 = β * v1 + v0
	fq.f.sub(c[1], t[1], t[0])     // c1 = (a0 + a1)(b0 + b1) - (v0+v1)
}

func (fq *fq2) square(c, a *fe2) {
	t := fq.t
	// c0 = (a0 - a1) * (a0 - β * a1) + a0 * a1 + β * a0 * a1
	// c1 = 2 * a0 * a1
	fq.f.sub(t[0], a[0], a[1])     // v0 = a0 - a1
	fq.mulByNonResidue(t[1], a[1]) // β * a1
	fq.f.sub(t[2], a[0], t[1])     // v3 = a0 -  β * a1
	fq.f.mul(t[1], a[0], a[1])     // v2 = a0 * a1
	fq.f.mul(t[0], t[0], t[2])     // v0 * v3
	fq.f.add(t[0], t[1], t[0])     // v0 = v0 * v3 + v2
	fq.mulByNonResidue(t[2], t[1]) // β * v2
	fq.f.add(c[0], t[0], t[2])     // c0 = v0 + β * v2
	fq.f.double(c[1], t[1])        // c1 = 2*v2
}

func (fq *fq2) inverse(c, a *fe2) {
	t := fq.t
	// c0 = a0 * (a0^2 - β * a1^2)^-1
	// c1 = a1 * (a0^2 - β * a1^2)^-1
	fq.f.square(t[0], a[0])        // v0 = a0^2
	fq.f.square(t[1], a[1])        // v1 = a1^2
	fq.mulByNonResidue(t[1], t[1]) // β * v1
	fq.f.sub(t[1], t[0], t[1])     // v0 = v0 - β * v1
	fq.f.inverse(t[0], t[1])       // v1 = v0^-1
	fq.f.mul(c[0], a[0], t[0])     // c0 = a0 * v1
	fq.f.mul(t[0], a[1], t[0])     // a1 * v1
	fq.f.neg(c[1], t[0])           // c1 = -a1 * v1
}

func (fq *fq2) exp(c, a *fe2, e *big.Int) {
	z := fq.one()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.square(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.copy(c, z)
}

func (fq *fq2) mulByFq(c, a *fe2, b fieldElement) {
	fq.f.mul(c[0], a[0], b)
	fq.f.mul(c[1], a[1], b)
}

func (fq *fq2) frobeniusMap(c, a *fe2, power uint) {
	fq.f.copy(c[0], a[0])
	fq.f.mul(c[1], a[1], fq.frobeniusCoeffs[power%2])
}

func (fq *fq2) calculateFrobeniusCoeffs() bool {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = fq.newElement()
	}
	power, rem := new(big.Int), new(big.Int)
	power.Sub(fq.f.pbig, big.NewInt(1))
	power.DivMod(power, big.NewInt(2), rem)
	if rem.Uint64() != 0 {
		return false
	}
	fq.f.exp(fq.frobeniusCoeffs[1], fq.nonResidue, power)
	fq.f.copy(fq.frobeniusCoeffs[0], fq.f.one)
	return true
}

func (fq *fq2) calculateFrobeniusCoeffsWithPrecomputation(f1 fieldElement) {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = fq.newElement()
	}
	fq.f.copy(fq.frobeniusCoeffs[0], fq.f.one)
	fq.f.square(fq.frobeniusCoeffs[1], f1)
}
