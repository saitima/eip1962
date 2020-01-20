package eip

import (
	"fmt"
	"math/big"
)

type fe6 [3]fe2

type fq6 struct {
	f               *fq2
	nonResidue      *fe2
	t               []*fe2
	frobeniusCoeffs *[2][6]*fe2
}

func newFq6(f *fq2, nonResidue []byte) (*fq6, error) {
	nonResidue_ := f.newElement()
	if nonResidue != nil {
		var err error
		nonResidue_, err = f.fromBytes(nonResidue)
		if err != nil {
			return nil, err
		}
	}
	t := make([]*fe2, 6)
	for i := 0; i < 6; i++ {
		t[i] = f.newElement()
		f.copy(t[i], f.zero())
	}
	return &fq6{f, nonResidue_, t, nil}, nil
}

func (fq *fq6) newElement() *fe6 {
	return fq.zero()
}

func (fq *fq6) fromBytes(in []byte) (*fe6, error) {
	byteLen := fq.f.f.limbSize * 8 * 2
	totalLen := len(&fe6{}) * byteLen
	if len(in) < totalLen {
		return nil, fmt.Errorf("input string should be larger than %d bytes", totalLen)
	}
	c := fq.newElement()
	for i := 0; i < len(&fe6{}); i++ {
		elem, err := fq.f.fromBytes(in[i*byteLen : (i+1)*byteLen])
		if err != nil {
			return nil, err
		}
		c[i] = *elem
	}
	return c, nil
}

func (fq *fq6) toBytes(a *fe6) []byte {
	byteLen := fq.f.f.limbSize * 8 * 2
	out := make([]byte, len(a)*byteLen)
	for i := 0; i < len(a); i++ {
		copy(out[i*byteLen:(i+1)*byteLen], fq.f.toBytes(&a[i]))
	}
	return out
}

func (fq *fq6) toString(a *fe6) string {
	return fmt.Sprintf(
		"c0: %s c1: %s\n",
		fq.f.toString(&a[0]),
		fq.f.toString(&a[1]),
	)
}

func (fq *fq6) toStringNoTransform(a *fe6) string {
	return fmt.Sprintf(
		"c0: %s c1: %s\n",
		fq.f.toStringNoTransform(&a[0]),
		fq.f.toStringNoTransform(&a[1]),
	)
}

func (fq *fq6) zero() *fe6 {
	return &fe6{
		*fq.f.zero(),
		*fq.f.zero(),
		*fq.f.zero(),
	}
}

func (fq *fq6) one() *fe6 {
	a := fq.zero()
	fq.f.copy(&a[0], fq.f.one())
	return a
}

func (fq *fq6) isZero(a *fe6) bool {
	return fq.f.isZero(&a[0]) && fq.f.isZero(&a[1]) && fq.f.isZero(&a[2])
}

func (fq *fq6) equal(a, b *fe6) bool {
	return fq.f.equal(&a[0], &b[0]) && fq.f.equal(&a[1], &b[1]) && fq.f.equal(&a[2], &b[2])
}

func (fq *fq6) copy(c, a *fe6) *fe6 {
	fq.f.copy(&c[0], &a[0])
	fq.f.copy(&c[1], &a[1])
	fq.f.copy(&c[2], &a[2])
	return c
}

func (fq *fq6) add(c, a, b *fe6) *fe6 {
	fq.f.add(&c[0], &a[0], &b[0])
	fq.f.add(&c[1], &a[1], &b[1])
	fq.f.add(&c[2], &a[2], &b[2])
	return c
}

func (fq *fq6) double(c, a *fe6) *fe6 {
	fq.f.double(&c[0], &a[0])
	fq.f.double(&c[1], &a[1])
	fq.f.double(&c[2], &a[2])
	return c
}

func (fq *fq6) sub(c, a, b *fe6) *fe6 {
	fq.f.sub(&c[0], &a[0], &b[0])
	fq.f.sub(&c[1], &a[1], &b[1])
	fq.f.sub(&c[2], &a[2], &b[2])
	return c
}

func (fq *fq6) neg(c, a *fe6) *fe6 {
	fq.f.neg(&c[0], &a[0])
	fq.f.neg(&c[1], &a[1])
	fq.f.neg(&c[2], &a[2])
	return c
}

func (fq *fq6) conjugate(c, a *fe6) *fe6 {
	fq.copy(c, a)
	fq.f.neg(&c[1], &a[1])
	return c
}

func (fq *fq6) mulByNonResidue(c, a *fe2) {
	fq.f.mul(c, a, fq.nonResidue)
}

func (fq *fq6) mul(c, a, b *fe6) {
	t := fq.t
	fq.f.mul(t[0], &a[0], &b[0])
	fq.f.mul(t[1], &a[1], &b[1])
	fq.f.mul(t[2], &a[2], &b[2])
	fq.f.add(t[3], &a[1], &a[2])
	fq.f.add(t[4], &b[1], &b[2])
	fq.f.mul(t[3], t[3], t[4])
	fq.f.add(t[4], t[1], t[2])
	fq.f.sub(t[3], t[3], t[4])
	fq.mulByNonResidue(t[3], t[3])
	fq.f.add(t[5], t[0], t[3])
	fq.f.add(t[3], &a[0], &a[1])
	fq.f.add(t[4], &b[0], &b[1])
	fq.f.mul(t[3], t[3], t[4])
	fq.f.add(t[4], t[0], t[1])
	fq.f.sub(t[3], t[3], t[4])
	fq.mulByNonResidue(t[4], t[2])
	fq.f.add(&c[1], t[3], t[4])
	fq.f.add(t[3], &a[0], &a[2])
	fq.f.add(t[4], &b[0], &b[2])
	fq.f.mul(t[3], t[3], t[4])
	fq.f.add(t[4], t[0], t[2])
	fq.f.sub(t[3], t[3], t[4])
	fq.f.add(&c[2], t[1], t[3])
	fq.f.copy(&c[0], t[5])
}

func (fq *fq6) square(c, a *fe6) {
	t := fq.t
	fq.f.square(t[0], &a[0])
	fq.f.mul(t[1], &a[0], &a[1])
	fq.f.add(t[1], t[1], t[1])
	fq.f.sub(t[2], &a[0], &a[1])
	fq.f.add(t[2], t[2], &a[2])
	fq.f.square(t[2], t[2])
	fq.f.mul(t[3], &a[1], &a[2])
	fq.f.add(t[3], t[3], t[3])
	fq.f.square(t[4], &a[2])
	fq.mulByNonResidue(t[5], t[3])
	fq.f.add(&c[0], t[0], t[5])
	fq.mulByNonResidue(t[5], t[4])
	fq.f.add(&c[1], t[1], t[5])
	fq.f.add(t[1], t[1], t[2])
	fq.f.add(t[1], t[1], t[3])
	fq.f.add(t[0], t[0], t[4])
	fq.f.sub(&c[2], t[1], t[0])
}

func (fq *fq6) inverse(c, a *fe6) {
	t := fq.t
	fq.f.square(t[0], &a[0])       // v0 = a0^2
	fq.f.mul(t[1], &a[1], &a[2])   // v5 = a1 * a2
	fq.mulByNonResidue(t[1], t[1]) // α * v5
	fq.f.sub(t[0], t[0], t[1])     // A = v0 - α * v5
	fq.f.square(t[1], &a[1])       // v1 = a1^2
	fq.f.mul(t[2], &a[0], &a[2])   // v4 = a0 * a2
	fq.f.sub(t[1], t[1], t[2])     // C = v1 - v4
	fq.f.square(t[2], &a[2])       // v2 = a2^2
	fq.mulByNonResidue(t[2], t[2]) // α * v2
	fq.f.mul(t[3], &a[0], &a[1])   // v3 = a0 * a1
	fq.f.sub(t[2], t[2], t[3])     // B = α * v2 - v3
	fq.f.mul(t[3], &a[2], t[2])    // B * a2
	fq.f.mul(t[4], &a[1], t[1])    // C * a1
	fq.f.add(t[3], t[3], t[4])     // C * a1 + B * a2
	fq.mulByNonResidue(t[3], t[3]) // α * (C * a1 + B * a2)
	fq.f.mul(t[4], &a[0], t[0])    // A * a0
	fq.f.add(t[3], t[3], t[4])     // v6 = A * a0 * α * (C * a1 + B * a2)
	fq.f.inverse(t[3], t[3])       // F = v6^-1
	fq.f.mul(&c[0], t[0], t[3])    // c0 = A * F
	fq.f.mul(&c[1], t[2], t[3])    // c1 = B * F
	fq.f.mul(&c[2], t[1], t[3])    // c2 = C * F
}

func (fq *fq6) exp(c, a *fe6, e *big.Int) {
	z := fq.one()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.square(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.copy(c, z)
}

func (fq *fq6) mulBy01(e *fe6, c0, c1 *fe2) {
	t := fq.t
	fq.f.mul(t[0], &e[0], c0)
	fq.f.mul(t[1], &e[1], c1)
	fq.f.add(t[5], &e[1], &e[2])
	fq.f.mul(t[2], c1, t[5])
	fq.f.sub(t[2], t[2], t[1])
	fq.mulByNonResidue(t[2], t[2])
	fq.f.add(t[2], t[2], t[0])
	fq.f.add(t[5], &e[0], &e[2])
	fq.f.mul(t[3], c0, t[5])
	fq.f.sub(t[3], t[3], t[0])
	fq.f.add(t[3], t[3], t[1])
	fq.f.add(t[4], c0, c1)
	fq.f.add(t[5], &e[0], &e[1])
	fq.f.mul(t[4], t[4], t[5])
	fq.f.sub(t[4], t[4], t[0])
	fq.f.sub(t[4], t[4], t[1])
	fq.f.copy(&e[0], t[2])
	fq.f.copy(&e[1], t[4])
	fq.f.copy(&e[2], t[3])
}

func (fq *fq6) mulBy1(e *fe6, c1 *fe2) {
	t := fq.t
	fq.f.mul(t[0], &e[1], c1)
	fq.f.add(t[1], &e[1], &e[2])
	fq.f.mul(t[1], t[1], c1)
	fq.f.sub(t[1], t[1], t[0])
	fq.mulByNonResidue(t[1], t[1])
	fq.f.add(t[2], &e[0], &e[1])
	fq.f.mul(t[2], t[2], c1)
	fq.f.sub(&e[1], t[2], t[0])
	fq.f.copy(&e[0], t[1])
	fq.f.copy(&e[2], t[0])
}

func (fq *fq6) frobeniusMap(c, a *fe6, power uint) {
	fq.copy(c, a)
	fq.f.frobeniusMap(&c[0], &a[0], power)
	fq.f.frobeniusMap(&c[1], &a[1], power)
	fq.f.frobeniusMap(&c[2], &a[2], power)
	fq.f.mul(&c[1], &c[1], fq.frobeniusCoeffs[0][power%6])
	fq.f.mul(&c[2], &c[2], fq.frobeniusCoeffs[1][power%6])
}

func (fq *fq6) calculateFrobeniusCoeffs() bool {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = new([2][6]*fe2)
		for i := 0; i < 6; i++ {
			fq.frobeniusCoeffs[0][i] = fq.f.newElement()
			fq.frobeniusCoeffs[1][i] = fq.f.newElement()
		}
	}
	modulus := fq.f.f.pbig
	bigOne, bigThree := big.NewInt(1), big.NewInt(3)
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	fq.f.copy(fq.frobeniusCoeffs[0][0], fq.f.one())
	fq.f.copy(fq.frobeniusCoeffs[1][0], fq.f.one())
	for i := 1; i <= 3; i++ {
		power.Sub(qPower, bigOne)
		power.DivMod(power, bigThree, rem)
		if rem.Uint64() != 0 {
			return false
		}
		fq.f.exp(fq.frobeniusCoeffs[0][i], fq.nonResidue, power)
		fq.f.square(fq.frobeniusCoeffs[1][i], fq.frobeniusCoeffs[0][i])
		qPower.Mul(qPower, modulus)
	}
	return true
}

func (fq *fq6) calculateFrobeniusCoeffsWithPrecomputation(f1, f2 *fe2) bool {
	if fq.frobeniusCoeffs == nil {
		fq.frobeniusCoeffs = new([2][6]*fe2)
		for i := 0; i < 6; i++ {
			fq.frobeniusCoeffs[0][i] = fq.f.newElement()
			fq.frobeniusCoeffs[1][i] = fq.f.newElement()
		}
	}
	fq.f.copy(fq.frobeniusCoeffs[0][0], fq.f.one())
	fq.f.copy(fq.frobeniusCoeffs[1][0], fq.f.one())
	fq.f.square(fq.frobeniusCoeffs[0][1], f1)
	fq.f.square(fq.frobeniusCoeffs[0][2], f2)
	fq.f.frobeniusMap(fq.frobeniusCoeffs[0][3], fq.frobeniusCoeffs[0][2], 1)
	fq.f.mul(fq.frobeniusCoeffs[0][3], fq.frobeniusCoeffs[0][3], fq.frobeniusCoeffs[0][1])
	for i := 1; i <= 3; i++ {
		fq.f.square(fq.frobeniusCoeffs[1][i], fq.frobeniusCoeffs[0][i])
	}
	fq.f.copy(fq.frobeniusCoeffs[0][4], fq.f.zero())
	fq.f.copy(fq.frobeniusCoeffs[1][4], fq.f.zero())
	fq.f.copy(fq.frobeniusCoeffs[0][5], fq.f.zero())
	fq.f.copy(fq.frobeniusCoeffs[1][5], fq.f.zero())
	return true
}
