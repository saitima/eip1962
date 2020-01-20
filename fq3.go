package eip

import (
	"fmt"
	"math/big"
)

type fe3 [3]fieldElement

type fq3 struct {
	f               *field
	nonResidue      fieldElement
	t               []fieldElement
	frobeniusCoeffs *[2]*fe3
}

func newFq3(f *field, nonResidue []byte) (*fq3, error) {
	nonResidue_ := f.newFieldElement()
	if nonResidue != nil {
		var err error
		nonResidue_, err = f.newFieldElementFromBytes(nonResidue)
		if err != nil {
			return nil, err
		}
	}
	t := make([]fieldElement, 6)
	for i := 0; i < 6; i++ {
		t[i] = f.newFieldElement()
		f.copy(t[i], f.zero)
	}
	return &fq3{f, nonResidue_, t, nil}, nil
}

func (fq *fq3) newElement() *fe3 {
	fe := &fe3{fq.f.newFieldElement(), fq.f.newFieldElement(), fq.f.newFieldElement()}
	fq.f.copy(fe[0], fq.f.zero)
	fq.f.copy(fe[1], fq.f.zero)
	fq.f.copy(fe[2], fq.f.zero)
	return fe
}

func (fq *fq3) fromBytes(in []byte) (*fe3, error) {
	byteLen := fq.f.limbSize * 8
	totalLen := len(&fe3{}) * byteLen
	if len(in) < totalLen {
		return nil, fmt.Errorf("input string should be larger than %d bytes", totalLen)
	}
	c := fq.newElement()
	var err error
	for i := 0; i < len(&fe3{}); i++ {
		c[i], err = fq.f.newFieldElementFromBytes(in[i*byteLen : (i+1)*byteLen])
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (fq *fq3) toBytes(a *fe3) []byte {
	byteLen := fq.f.limbSize * 8
	out := make([]byte, len(a)*byteLen)
	for i := 0; i < len(a); i++ {
		copy(out[i*byteLen:(i+1)*byteLen], fq.f.toBytes(a[i]))
	}
	return out
}

func (fq *fq3) toString(a *fe3) string {
	return fmt.Sprintf(
		"c0: %s c1: %s c2: %s \n",
		fq.f.toString(a[0]),
		fq.f.toString(a[1]),
		fq.f.toString(a[2]),
	)
}

func (fq *fq3) toStringNoTransform(a *fe3) string {
	return fmt.Sprintf(
		"c0: %s c1: %s c2: %s\n",
		fq.f.toStringNoTransform(a[0]),
		fq.f.toStringNoTransform(a[1]),
		fq.f.toStringNoTransform(a[2]),
	)
}

func (fq *fq3) zero() *fe3 {
	return fq.newElement()
}

func (fq *fq3) one() *fe3 {
	a := fq.newElement()
	fq.f.copy(a[0], fq.f.one)
	return a
}

func (fq *fq3) isZero(a *fe3) bool {
	return fq.f.isZero(a[0]) && fq.f.isZero(a[1]) && fq.f.isZero(a[2])
}

func (fq *fq3) equal(a, b *fe3) bool {
	return fq.f.equal(a[0], b[0]) && fq.f.equal(a[1], b[1]) && fq.f.equal(a[2], b[2])
}

func (fq *fq3) copy(c, a *fe3) *fe3 {
	fq.f.copy(c[0], a[0])
	fq.f.copy(c[1], a[1])
	fq.f.copy(c[2], a[2])
	return c
}

func (fq *fq3) add(c, a, b *fe3) *fe3 {
	fq.f.add(c[0], a[0], b[0])
	fq.f.add(c[1], a[1], b[1])
	fq.f.add(c[2], a[2], b[2])
	return c
}

func (fq *fq3) double(c, a *fe3) *fe3 {
	fq.f.double(c[0], a[0])
	fq.f.double(c[1], a[1])
	fq.f.double(c[2], a[2])
	return c
}

func (fq *fq3) sub(c, a, b *fe3) *fe3 {
	fq.f.sub(c[0], a[0], b[0])
	fq.f.sub(c[1], a[1], b[1])
	fq.f.sub(c[2], a[2], b[2])
	return c
}

func (fq *fq3) neg(c, a *fe3) *fe3 {
	fq.f.neg(c[0], a[0])
	fq.f.neg(c[1], a[1])
	fq.f.neg(c[2], a[2])
	return c
}

func (fq *fq3) conjugate(c, a *fe3) *fe3 {
	fq.copy(c, a)
	fq.f.neg(c[1], a[1])
	return c
}

func (fq *fq3) mulByNonResidue(c, a fieldElement) {
	fq.f.mul(c, a, fq.nonResidue)
}

func (fq *fq3) mul(c, a, b *fe3) {
	t := fq.t
	fq.f.mul(t[0], a[0], b[0])
	fq.f.mul(t[1], a[1], b[1])
	fq.f.mul(t[2], a[2], b[2])
	fq.f.add(t[3], a[1], a[2])
	fq.f.add(t[4], b[1], b[2])
	fq.f.mul(t[3], t[3], t[4])
	fq.f.add(t[4], t[1], t[2])
	fq.f.sub(t[3], t[3], t[4])
	fq.mulByNonResidue(t[3], t[3])
	fq.f.add(t[5], t[0], t[3])
	fq.f.add(t[3], a[0], a[1])
	fq.f.add(t[4], b[0], b[1])
	fq.f.mul(t[3], t[3], t[4])
	fq.f.add(t[4], t[0], t[1])
	fq.f.sub(t[3], t[3], t[4])
	fq.mulByNonResidue(t[4], t[2])
	fq.f.add(c[1], t[3], t[4])
	fq.f.add(t[3], a[0], a[2])
	fq.f.add(t[4], b[0], b[2])
	fq.f.mul(t[3], t[3], t[4])
	fq.f.add(t[4], t[0], t[2])
	fq.f.sub(t[3], t[3], t[4])
	fq.f.add(c[2], t[1], t[3])
	fq.f.copy(c[0], t[5])
}

func (fq *fq3) square(c, a *fe3) {
	t := fq.t
	fq.f.square(t[0], a[0])
	fq.f.mul(t[1], a[0], a[1])
	fq.f.add(t[1], t[1], t[1])
	fq.f.sub(t[2], a[0], a[1])
	fq.f.add(t[2], t[2], a[2])
	fq.f.square(t[2], t[2])
	fq.f.mul(t[3], a[1], a[2])
	fq.f.add(t[3], t[3], t[3])
	fq.f.square(t[4], a[2])
	fq.mulByNonResidue(t[5], t[3])
	fq.f.add(c[0], t[0], t[5])
	fq.mulByNonResidue(t[5], t[4])
	fq.f.add(c[1], t[1], t[5])
	fq.f.add(t[1], t[1], t[2])
	fq.f.add(t[1], t[1], t[3])
	fq.f.add(t[0], t[0], t[4])
	fq.f.sub(c[2], t[1], t[0])
}

func (fq *fq3) inverse(c, a *fe3) {
	t := fq.t
	fq.f.square(t[0], a[0])        // v0 = a0^2
	fq.f.mul(t[1], a[1], a[2])     // v5 = a1 * a2
	fq.mulByNonResidue(t[1], t[1]) // α * v5
	fq.f.sub(t[0], t[0], t[1])     // A = v0 - α * v5
	fq.f.square(t[1], a[1])        // v1 = a1^2
	fq.f.mul(t[2], a[0], a[2])     // v4 = a0 * a2
	fq.f.sub(t[1], t[1], t[2])     // C = v1 - v4
	fq.f.square(t[2], a[2])        // v2 = a2^2
	fq.mulByNonResidue(t[2], t[2]) // α * v2
	fq.f.mul(t[3], a[0], a[1])     // v3 = a0 * a1
	fq.f.sub(t[2], t[2], t[3])     // B = α * v2 - v3
	fq.f.mul(t[3], a[2], t[2])     // B * a2
	fq.f.mul(t[4], a[1], t[1])     // C * a1
	fq.f.add(t[3], t[3], t[4])     // C * a1 + B * a2
	fq.mulByNonResidue(t[3], t[3]) // α * (C * a1 + B * a2)
	fq.f.mul(t[4], a[0], t[0])     // A * a0
	fq.f.add(t[3], t[3], t[4])     // v6 = A * a0 * α * (C * a1 + B * a2)
	fq.f.inverse(t[3], t[3])       // F = v6^-1
	fq.f.mul(c[0], t[0], t[3])     // c0 = A * F
	fq.f.mul(c[1], t[2], t[3])     // c1 = B * F
	fq.f.mul(c[2], t[1], t[3])     // c2 = C * F
}

func (fq *fq3) exp(c, a *fe3, e *big.Int) {
	z := fq.one()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.square(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.copy(c, z)
}

func (fq *fq3) mulByFq(c, a *fe3, b fieldElement) {
	fq.f.mul(c[0], a[0], b)
	fq.f.mul(c[1], a[1], b)
	fq.f.mul(c[2], a[2], b)
}

func (fq *fq3) frobeniusMap(c, a *fe3, power uint) {
	fq.copy(c, a)
	fq.f.mul(c[1], a[1], fq.frobeniusCoeffs[0][power%3])
	fq.f.mul(c[2], a[2], fq.frobeniusCoeffs[1][power%3])
}

func (fq *fq3) calculateFrobeniusCoeffs() bool {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = new([2]*fe3)
		fq.frobeniusCoeffs[0] = fq.newElement()
		fq.frobeniusCoeffs[1] = fq.newElement()
	}
	modulus := fq.f.pbig
	bigOne, bigTwo, bigThree := big.NewInt(1), big.NewInt(2), big.NewInt(3)
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	for i := 1; i <= 2; i++ {
		power.Sub(qPower, bigOne)
		power.DivMod(power, bigThree, rem)
		if rem.Uint64() != 0 {
			return false
		}
		fq.f.exp(fq.frobeniusCoeffs[0][i], fq.nonResidue, power)
		fq.f.exp(fq.frobeniusCoeffs[1][i], fq.frobeniusCoeffs[0][i], bigTwo)
		qPower.Mul(qPower, modulus)
	}
	fq.f.copy(fq.frobeniusCoeffs[0][0], fq.f.one)
	fq.f.copy(fq.frobeniusCoeffs[1][0], fq.f.one)
	return true
}

func (fq *fq3) calculateFrobeniusCoeffsWithPrecomputation(f1 fieldElement) {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = new([2]*fe3)
		fq.frobeniusCoeffs[0] = fq.newElement()
		fq.frobeniusCoeffs[1] = fq.newElement()
	}

	fq.f.copy(fq.frobeniusCoeffs[0][0], fq.f.one)
	fq.f.copy(fq.frobeniusCoeffs[1][0], fq.f.one)
	fq.f.square(fq.frobeniusCoeffs[0][1], f1)
	fq.f.square(fq.frobeniusCoeffs[0][2], fq.frobeniusCoeffs[0][1])
	for i := 1; i <= 2; i++ {
		fq.f.square(fq.frobeniusCoeffs[1][i], fq.frobeniusCoeffs[0][i])
		fq.f.square(fq.frobeniusCoeffs[1][i], fq.frobeniusCoeffs[0][i])
	}
}
