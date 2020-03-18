package eip

import (
	"fmt"
	"io"
	"math/big"
)

type fe3 [3]fe

type fq3 struct {
	f               *fq
	nonResidue      fe
	t               []fe
	frobeniusCoeffs *[2]*fe3
}

func newFq3(fq *fq, nonResidueBuf []byte) (*fq3, error) {
	nonResidue := fq.new()
	if nonResidueBuf != nil {
		var err error
		nonResidue, err = fq.fromBytes(nonResidueBuf)
		if err != nil {
			return nil, err
		}
	}
	t := make([]fe, 6)
	for i := 0; i < 6; i++ {
		t[i] = fq.new()
	}
	return &fq3{fq, nonResidue, t, nil}, nil
}

func (fq3 *fq3) byteSize() int {
	fq := fq3.fq()
	return fq.byteSize() * 3
}

func (fq3 *fq3) new() *fe3 {
	fq := fq3.fq()
	return &fe3{fq.new(), fq.new(), fq.new()}
}

func (fq3 *fq3) modulus() *big.Int {
	fq := fq3.fq()
	return fq.modulus()
}

func (fq3 *fq3) rand(r io.Reader) *fe3 {
	fq := fq3.fq()
	return &fe3{fq.rand(r), fq.rand(r), fq.rand(r)}
}

func (fq3 *fq3) fromBytes(in []byte) (*fe3, error) {
	fq := fq3.fq()
	byteSize := fq.byteSize()
	if len(in) != byteSize*3 {
		return nil, fmt.Errorf("input string should be larger than %d bytes", byteSize)
	}
	var err error
	c0, err := fq.fromBytes(in[:byteSize])
	if err != nil {
		return nil, err
	}
	c1, err := fq.fromBytes(in[byteSize : byteSize*2])
	if err != nil {
		return nil, err
	}
	c2, err := fq.fromBytes(in[byteSize*2:])
	if err != nil {
		return nil, err
	}
	return &fe3{c0, c1, c2}, nil
}

func (fq3 *fq3) toBytes(a *fe3) []byte {
	fq := fq3.fq()
	byteSize := fq.byteSize()
	out := make([]byte, byteSize*3)
	copy(out[:byteSize], fq.toBytes(a[0]))
	copy(out[byteSize:byteSize*2], fq.toBytes(a[1]))
	copy(out[byteSize*2:byteSize*3], fq.toBytes(a[2]))
	return out
}

func (fq3 *fq3) toBytesDense(a *fe3) []byte {
	fq := fq3.fq()
	byteSize := fq.modulusByteLen
	out := make([]byte, 3*byteSize)
	copy(out[:byteSize], fq.toBytesDense(a[0]))
	copy(out[byteSize:byteSize*2], fq.toBytesDense(a[1]))
	copy(out[byteSize*2:byteSize*3], fq.toBytesDense(a[2]))
	return out
}

func (fq3 *fq3) toString(a *fe3) string {
	fq := fq3.fq()
	return fmt.Sprintf("%s\n%s\n%s", fq.toString(a[0]), fq.toString(a[1]), fq.toString(a[2]))
}

func (fq3 *fq3) toStringNoTransform(a *fe3) string {
	fq := fq3.fq()
	return fmt.Sprintf("%s\n%s\n%s", fq.toStringNoTransform(a[0]), fq.toStringNoTransform(a[1]), fq.toStringNoTransform(a[2]))
}

func (fq3 *fq3) zero() *fe3 {
	return fq3.new()
}

func (fq3 *fq3) one() *fe3 {
	fq := fq3.fq()
	a := fq3.new()
	fq.copy(a[0], fq.one)
	return a
}

func (fq3 *fq3) twistOne() *fe3 {
	fq := fq3.fq()
	a := fq3.new()
	fq.copy(a[1], fq.one)
	return a
}

func (fq3 *fq3) isZero(a *fe3) bool {
	fq := fq3.fq()
	return fq.isZero(a[0]) && fq.isZero(a[1]) && fq.isZero(a[2])
}

func (fq3 *fq3) isOne(a *fe3) bool {
	fq := fq3.fq()
	return fq.isOne(a[0]) && fq.isZero(a[1]) && fq.isZero(a[2])
}

func (fq3 *fq3) equal(a, b *fe3) bool {
	fq := fq3.fq()
	return fq.equal(a[0], b[0]) && fq.equal(a[1], b[1]) && fq.equal(a[2], b[2])
}

func (fq3 *fq3) copy(c, a *fe3) *fe3 {
	fq := fq3.fq()
	fq.copy(c[0], a[0])
	fq.copy(c[1], a[1])
	fq.copy(c[2], a[2])
	return c
}

func (fq3 *fq3) isNonResidue(a *fe3, degree int) bool {
	zero := big.NewInt(0)
	result := fq3.new()
	p := fq3.modulus()
	exp := new(big.Int).Sub(p, big.NewInt(1))
	exp, rem := new(big.Int).DivMod(exp, big.NewInt(int64(degree)), zero)
	if rem.Cmp(zero) != 0 {
		return false
	}
	fq3.exp(result, a, exp)
	if fq3.equal(result, fq3.one()) {
		return false
	}
	return true
}

func (fq3 *fq3) add(c, a, b *fe3) *fe3 {
	fq := fq3.fq()
	fq.add(c[0], a[0], b[0])
	fq.add(c[1], a[1], b[1])
	fq.add(c[2], a[2], b[2])
	return c
}

func (fq3 *fq3) double(c, a *fe3) *fe3 {
	fq := fq3.fq()
	fq.double(c[0], a[0])
	fq.double(c[1], a[1])
	fq.double(c[2], a[2])
	return c
}

func (fq3 *fq3) sub(c, a, b *fe3) *fe3 {
	fq := fq3.fq()
	fq.sub(c[0], a[0], b[0])
	fq.sub(c[1], a[1], b[1])
	fq.sub(c[2], a[2], b[2])
	return c
}

func (fq3 *fq3) neg(c, a *fe3) *fe3 {
	fq := fq3.fq()
	fq.neg(c[0], a[0])
	fq.neg(c[1], a[1])
	fq.neg(c[2], a[2])
	return c
}

func (fq3 *fq3) conjugate(c, a *fe3) *fe3 {
	fq := fq3.fq()
	fq3.copy(c, a)
	fq.neg(c[1], a[1])
	return c
}

func (fq3 *fq3) mulByNonResidue(c, a fe) {
	fq := fq3.fq()
	fq.mul(c, a, fq3.nonResidue)
}

func (fq3 *fq3) mul(c, a, b *fe3) {
	fq, t := fq3.fq(), fq3.t
	fq.mul(t[0], a[0], b[0])
	fq.mul(t[1], a[1], b[1])
	fq.mul(t[2], a[2], b[2])
	fq.add(t[3], a[1], a[2])
	fq.add(t[4], b[1], b[2])
	fq.mul(t[3], t[3], t[4])
	fq.add(t[4], t[1], t[2])
	fq.sub(t[3], t[3], t[4])
	fq3.mulByNonResidue(t[3], t[3])
	fq.add(t[5], t[0], t[3])
	fq.add(t[3], a[0], a[1])
	fq.add(t[4], b[0], b[1])
	fq.mul(t[3], t[3], t[4])
	fq.add(t[4], t[0], t[1])
	fq.sub(t[3], t[3], t[4])
	fq3.mulByNonResidue(t[4], t[2])
	fq.add(c[1], t[3], t[4])
	fq.add(t[3], a[0], a[2])
	fq.add(t[4], b[0], b[2])
	fq.mul(t[3], t[3], t[4])
	fq.add(t[4], t[0], t[2])
	fq.sub(t[3], t[3], t[4])
	fq.add(c[2], t[1], t[3])
	fq.copy(c[0], t[5])
}

func (fq3 *fq3) square(c, a *fe3) {
	fq, t := fq3.fq(), fq3.t
	fq.square(t[0], a[0])
	fq.mul(t[1], a[0], a[1])
	fq.add(t[1], t[1], t[1])
	fq.sub(t[2], a[0], a[1])
	fq.add(t[2], t[2], a[2])
	fq.square(t[2], t[2])
	fq.mul(t[3], a[1], a[2])
	fq.add(t[3], t[3], t[3])
	fq.square(t[4], a[2])
	fq3.mulByNonResidue(t[5], t[3])
	fq.add(c[0], t[0], t[5])
	fq3.mulByNonResidue(t[5], t[4])
	fq.add(c[1], t[1], t[5])
	fq.add(t[1], t[1], t[2])
	fq.add(t[1], t[1], t[3])
	fq.add(t[0], t[0], t[4])
	fq.sub(c[2], t[1], t[0])
}

func (fq3 *fq3) inverse(c, a *fe3) bool {
	fq, t := fq3.fq(), fq3.t
	fq.square(t[0], a[0])           // v0 = a0^2
	fq.mul(t[1], a[1], a[2])        // v5 = a1 * a2
	fq3.mulByNonResidue(t[1], t[1]) // α * v5
	fq.sub(t[0], t[0], t[1])        // A = v0 - α * v5
	fq.square(t[1], a[1])           // v1 = a1^2
	fq.mul(t[2], a[0], a[2])        // v4 = a0 * a2
	fq.sub(t[1], t[1], t[2])        // C = v1 - v4
	fq.square(t[2], a[2])           // v2 = a2^2
	fq3.mulByNonResidue(t[2], t[2]) // α * v2
	fq.mul(t[3], a[0], a[1])        // v3 = a0 * a1
	fq.sub(t[2], t[2], t[3])        // B = α * v2 - v3
	fq.mul(t[3], a[2], t[2])        // B * a2
	fq.mul(t[4], a[1], t[1])        // C * a1
	fq.add(t[3], t[3], t[4])        // C * a1 + B * a2
	fq3.mulByNonResidue(t[3], t[3]) // α * (C * a1 + B * a2)
	fq.mul(t[4], a[0], t[0])        // A * a0
	fq.add(t[3], t[3], t[4])        // v6 = A * a0 * α * (C * a1 + B * a2)
	if ok := fq.inverse(t[3], t[3]); !ok {
		fq3.copy(c, fq3.zero())
		return false
	} // F = v6^-1
	fq.mul(c[0], t[0], t[3]) // c0 = A * F
	fq.mul(c[1], t[2], t[3]) // c1 = B * F
	fq.mul(c[2], t[1], t[3]) // c2 = C * F
	return true
}

func (fq3 *fq3) exp(c, a *fe3, e *big.Int) {
	z := fq3.one()
	found := false
	for i := e.BitLen() - 1; i >= 0; i-- {
		if found {
			fq3.square(z, z)
		} else {
			found = e.Bit(i) == 1
		}
		if e.Bit(i) == 1 {
			fq3.mul(z, z, a)
		}
	}
	fq3.copy(c, z)
}

func (fq3 *fq3) mulByFq(c, a *fe3, b fe) {
	fq := fq3.fq()
	fq.mul(c[0], a[0], b)
	fq.mul(c[1], a[1], b)
	fq.mul(c[2], a[2], b)
}

func (fq3 *fq3) frobeniusMap(c, a *fe3, power uint) {
	fq := fq3.fq()
	fq3.copy(c, a)
	fq.mul(c[1], a[1], fq3.frobeniusCoeffs[0][power%3])
	fq.mul(c[2], a[2], fq3.frobeniusCoeffs[1][power%3])
}

func (fq3 *fq3) calculateFrobeniusCoeffs() bool {
	fq := fq3.fq()
	if fq3.frobeniusCoeffs == nil {
		fq3.frobeniusCoeffs = new([2]*fe3)
		fq3.frobeniusCoeffs[0] = fq3.new()
		fq3.frobeniusCoeffs[1] = fq3.new()
	}
	modulus := fq.modulus()
	zero, one, two, three := big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(3)
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	for i := 1; i <= 2; i++ {
		power.Sub(qPower, one)
		power.DivMod(power, three, rem)
		if rem.Cmp(zero) != 0 {
			return false
		}
		fq.exp(fq3.frobeniusCoeffs[0][i], fq3.nonResidue, power)
		fq.exp(fq3.frobeniusCoeffs[1][i], fq3.frobeniusCoeffs[0][i], two)
		qPower.Mul(qPower, modulus)
	}
	fq.copy(fq3.frobeniusCoeffs[0][0], fq.one)
	fq.copy(fq3.frobeniusCoeffs[1][0], fq.one)
	return true
}

func (fq3 *fq3) calculateFrobeniusCoeffsWithPrecomputation(f1 fe) {
	fq := fq3.fq()
	if fq3.frobeniusCoeffs == nil {
		fq3.frobeniusCoeffs = new([2]*fe3)
		fq3.frobeniusCoeffs[0] = fq3.new()
		fq3.frobeniusCoeffs[1] = fq3.new()
	}
	fq.copy(fq3.frobeniusCoeffs[0][0], fq.one)
	fq.copy(fq3.frobeniusCoeffs[1][0], fq.one)
	fq.square(fq3.frobeniusCoeffs[0][1], f1)
	fq.square(fq3.frobeniusCoeffs[0][2], fq3.frobeniusCoeffs[0][1])
	fq.square(fq3.frobeniusCoeffs[1][1], fq3.frobeniusCoeffs[0][1])
	fq.square(fq3.frobeniusCoeffs[1][2], fq3.frobeniusCoeffs[0][2])
}

func (fq3 *fq3) fq() *fq {
	return fq3.f
}
