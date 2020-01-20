package eip

import (
	"fmt"
	"math/big"
)

type fe6q [2]*fe3

type fq6q struct {
	f               *fq3
	nonResidue      *fe3
	t               []*fe3
	frobeniusCoeffs *[6]fieldElement
}

func newFq6Quadratic(f *fq3, nonResidue []byte) (*fq6q, error) {
	nonResidue_ := f.newElement()
	if nonResidue != nil {
		var err error
		nonResidue_, err = f.fromBytes(nonResidue)
		if err != nil {
			return nil, err
		}
	}
	t := make([]*fe3, 4)
	for i := 0; i < 4; i++ {
		t[i] = f.newElement()
	}
	return &fq6q{f, nonResidue_, t, nil}, nil
}

func (fq *fq6q) newElement() *fe6q {
	fe := &fe6q{fq.f.newElement(), fq.f.newElement()}
	fq.f.copy(fe[0], fq.f.zero())
	fq.f.copy(fe[1], fq.f.zero())
	return fe
}

func (fq *fq6q) fromBytes(in []byte) (*fe6q, error) {
	byteLen := fq.f.f.limbSize * 8 * 3
	if len(in) < len(&fe6q{})*byteLen {
		return nil, fmt.Errorf("input string should be larger than %d bytes given %d", byteLen, len(in))
	}
	c := fq.newElement()
	var err error
	for i := 0; i < len(&fe6q{}); i++ {
		c[i], err = fq.f.fromBytes(in[i*byteLen : (i+1)*byteLen])
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (fq *fq6q) toBytes(a *fe6q) []byte {
	byteLen := fq.f.f.limbSize * 8 * 3
	out := make([]byte, len(a)*byteLen)
	for i := 0; i < len(a); i++ {
		copy(out[i*byteLen:(i+1)*byteLen], fq.f.toBytes(a[i]))
	}
	return out
}

func (fq *fq6q) toString(a *fe6q) string {
	return fmt.Sprintf(
		"c0: %s c1: %s\n",
		fq.f.toString(a[0]),
		fq.f.toString(a[1]),
	)
}

func (fq *fq6q) toStringNoTransform(a *fe6q) string {
	return fmt.Sprintf(
		"c0: %s c1: %s\n",
		fq.f.toStringNoTransform(a[0]),
		fq.f.toStringNoTransform(a[1]),
	)
}

func (fq *fq6q) zero() *fe6q {
	return fq.newElement()
}

func (fq *fq6q) one() *fe6q {
	a := fq.newElement()
	fq.f.copy(a[0], fq.f.one())
	return a
}

func (fq *fq6q) isZero(a *fe6q) bool {
	return fq.f.isZero(a[0]) && fq.f.isZero(a[1])
}

func (fq *fq6q) equal(a, b *fe6q) bool {
	return fq.f.equal(a[0], b[0]) && fq.f.equal(a[1], b[1])
}

func (fq *fq6q) copy(c, a *fe6q) *fe6q {
	fq.f.copy(c[0], a[0])
	fq.f.copy(c[1], a[1])
	return c
}

func (fq *fq6q) add(c, a, b *fe6q) *fe6q {
	fq.f.add(c[0], a[0], b[0])
	fq.f.add(c[1], a[1], b[1])
	return c
}

func (fq *fq6q) double(c, a *fe6q) *fe6q {
	fq.f.double(c[0], a[0])
	fq.f.double(c[1], a[1])
	return c
}

func (fq *fq6q) sub(c, a, b *fe6q) *fe6q {
	fq.f.sub(c[0], a[0], b[0])
	fq.f.sub(c[1], a[1], b[1])
	return c
}

func (fq *fq6q) neg(c, a *fe6q) *fe6q {
	fq.f.neg(c[0], a[0])
	fq.f.neg(c[1], a[1])
	return c
}

func (fq *fq6q) conjugate(c, a *fe6q) *fe6q {
	fq.copy(c, a)
	fq.f.neg(c[1], a[1])
	return c
}

func (fq *fq6q) mulByNonResidue(c, a *fe3) {
	o := fq.f.newElement()
	fq.f.f.copy(o[2], a[1])
	fq.f.f.copy(o[1], a[0])
	fq.f.mulByNonResidue(o[0], a[2])
	fq.f.copy(c, o)
}

func (fq *fq6q) mul(c, a, b *fe6q) {
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

func (fq *fq6q) square(c, a *fe6q) {
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

func (fq *fq6q) inverse(c, a *fe6q) {
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

func (fq *fq6q) exp(c, a *fe6q, e *big.Int) {
	z := fq.one()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.square(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.copy(c, z)
}

func (fq *fq6q) frobeniusMap(c, a *fe6q, power uint) {
	fq.copy(c, a)
	fq.f.frobeniusMap(c[0], a[0], power)
	fq.f.frobeniusMap(c[1], a[1], power)
	fq.f.mulByFq(c[1], c[1], fq.frobeniusCoeffs[power%6])
}

func (fq *fq6q) calculateFrobeniusCoeffs() bool {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = new([6]fieldElement)
		for i := 0; i < len(fq.frobeniusCoeffs); i++ {
			fq.frobeniusCoeffs[i] = fq.f.f.newFieldElement()
		}
	}
	modulus := fq.f.f.pbig
	f0 := fq.f.f.one
	fq.f.f.copy(fq.frobeniusCoeffs[0], f0)
	fq.f.f.copy(fq.frobeniusCoeffs[2], fq.f.f.zero)
	fq.f.f.copy(fq.frobeniusCoeffs[4], fq.f.f.zero)
	fq.f.f.copy(fq.frobeniusCoeffs[5], fq.f.f.zero)
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	for i := 1; i <= 3; i++ {
		if i == 2 {
			qPower.Mul(qPower, modulus)
			continue
		}
		power.Sub(qPower, big.NewInt(1))
		power.DivMod(power, big.NewInt(6), rem)
		if rem.Uint64() != 0 {
			return false
		}
		fq.f.f.exp(fq.frobeniusCoeffs[i], fq.f.nonResidue, power)
		qPower.Mul(qPower, modulus)
	}
	return true
}

func (fq *fq6q) calculateFrobeniusCoeffsWithPrecomputation(f1 fieldElement) {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = new([6]fieldElement)
		for i := 0; i < len(fq.frobeniusCoeffs); i++ {
			fq.frobeniusCoeffs[i] = fq.f.f.newFieldElement()
			fq.f.f.copy(fq.frobeniusCoeffs[i], fq.f.f.zero)
		}
	}
	fq.f.f.copy(fq.frobeniusCoeffs[0], fq.f.f.one)
	fq.f.f.copy(fq.frobeniusCoeffs[1], f1)
	fq.f.f.exp(fq.frobeniusCoeffs[3], fq.frobeniusCoeffs[1], big.NewInt(3))
}
