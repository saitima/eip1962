package eip

import (
	"fmt"
	"io"
	"math/big"
)

type fe2 [2]fe

type fq2 struct {
	f               *fq
	nonResidue      fe
	t               []fe
	frobeniusCoeffs *fe2
}

func newFq2(fq *fq, nonResidueBuf []byte) (*fq2, error) {
	nonResidue := fq.new()
	if nonResidueBuf != nil {
		var err error
		nonResidue, err = fq.fromBytes(nonResidueBuf)
		if err != nil {
			return nil, err
		}
	}
	t := make([]fe, 4)
	for i := 0; i < 4; i++ {
		t[i] = fq.new()
	}
	return &fq2{fq, nonResidue, t, nil}, nil
}

func (f *fq2) byteSize() int {
	fq := f.fq()
	return fq.byteSize() * 2
}

func (f *fq2) new() *fe2 {
	fq := f.fq()
	return &fe2{fq.new(), fq.new()}
}

func (f *fq2) modulus() *big.Int {
	fq := f.fq()
	return fq.modulus()
}

func (f *fq2) rand(r io.Reader) *fe2 {
	fq := f.fq()
	return &fe2{fq.rand(r), fq.rand(r)}
}

func (f *fq2) fromBytes(in []byte) (*fe2, error) {
	fq := f.fq()
	byteSize := fq.byteSize()
	if len(in) != byteSize*2 {
		return nil, fmt.Errorf("input string should be larger than %d bytes", byteSize*2)
	}
	var err error
	c0, err := fq.fromBytes(in[:byteSize])
	if err != nil {
		return nil, err
	}
	c1, err := fq.fromBytes(in[byteSize:])
	if err != nil {
		return nil, err
	}
	return &fe2{c0, c1}, nil
}

func (f *fq2) toBytes(a *fe2) []byte {
	fq := f.fq()
	byteSize := fq.byteSize()
	out := make([]byte, 2*byteSize)
	copy(out[:byteSize], fq.toBytes(a[0]))
	copy(out[byteSize:], fq.toBytes(a[1]))
	return out
}

func (f *fq2) toBytesDense(a *fe2) []byte {
	fq := f.fq()
	byteSize := fq.modulusByteLen
	out := make([]byte, 2*byteSize)
	copy(out[:byteSize], fq.toBytesDense(a[0]))
	copy(out[byteSize:], fq.toBytesDense(a[1]))
	return out
}

func (f *fq2) toString(a *fe2) string {
	fq := f.fq()
	return fmt.Sprintf("%s\n%s", fq.toString(a[0]), fq.toString(a[1]))
}

func (f *fq2) toStringNoTransform(a *fe2) string {
	fq := f.fq()
	return fmt.Sprintf("%s\n%s", fq.toStringNoTransform(a[0]), fq.toStringNoTransform(a[1]))
}

func (f *fq2) zero() *fe2 {
	return f.new()
}

func (f *fq2) one() *fe2 {
	fq := f.fq()
	a := f.new()
	fq.copy(a[0], fq.one)
	return a
}

func (f *fq2) twistOne() *fe2 {
	fq := f.fq()
	a := f.new()
	fq.copy(a[1], fq.one)
	return a
}

func (f *fq2) isZero(a *fe2) bool {
	fq := f.fq()
	return fq.isZero(a[0]) && fq.isZero(a[1])
}

func (f *fq2) isOne(a *fe2) bool {
	fq := f.fq()
	return fq.isOne(a[0]) && fq.isZero(a[1])
}

func (f *fq2) equal(a, b *fe2) bool {
	fq := f.fq()
	return fq.equal(a[0], b[0]) && fq.equal(a[1], b[1])
}

func (f *fq2) copy(c, a *fe2) *fe2 {
	fq := f.fq()
	fq.copy(c[0], a[0])
	fq.copy(c[1], a[1])
	return c
}

func (f *fq2) isNonResidue(a *fe2, degree int) bool {
	zero := big.NewInt(0)
	result := f.new()
	p := f.modulus()
	exp := new(big.Int).Sub(p, big.NewInt(1))
	exp, rem := new(big.Int).DivMod(exp, big.NewInt(int64(degree)), zero)
	if rem.Cmp(zero) != 0 {
		return false
	}
	f.exp(result, a, exp)
	if f.equal(result, f.one()) {
		return false
	}
	return true
}

func (f *fq2) add(c, a, b *fe2) *fe2 {
	fq := f.fq()
	fq.add(c[0], a[0], b[0])
	fq.add(c[1], a[1], b[1])
	return c
}

func (f *fq2) double(c, a *fe2) *fe2 {
	fq := f.fq()
	fq.double(c[0], a[0])
	fq.double(c[1], a[1])
	return c
}

func (f *fq2) sub(c, a, b *fe2) *fe2 {
	fq := f.fq()
	fq.sub(c[0], a[0], b[0])
	fq.sub(c[1], a[1], b[1])
	return c
}

func (f *fq2) neg(c, a *fe2) *fe2 {
	fq := f.fq()
	fq.neg(c[0], a[0])
	fq.neg(c[1], a[1])
	return c
}

func (f *fq2) conjugate(c, a *fe2) *fe2 {
	fq := f.fq()
	f.copy(c, a)
	fq.neg(c[1], a[1])
	return c
}

func (f *fq2) mulByNonResidue(c, a fe) {
	fq := f.fq()
	fq.mul(c, a, f.nonResidue)
}

func (f *fq2) mul(c, a, b *fe2) {
	fq, t := f.fq(), f.t
	// c0 = (a0 * b0) + β * (a1 * b1)
	// c1 = (a0 + a1) * (b0 + b1) - (a0 * b0 + a1 * b1)
	fq.mul(t[1], a[0], b[0])      // v0 = a0 * b0
	fq.mul(t[2], a[1], b[1])      // v1 = a1 * b1
	fq.add(t[0], t[1], t[2])      // v0 + v1
	f.mulByNonResidue(t[2], t[2]) // β * v1
	fq.add(t[3], t[1], t[2])      // β * v1 + v0
	fq.add(t[1], a[0], a[1])      // a0 + a1
	fq.add(t[2], b[0], b[1])      // b0 + b1
	fq.mul(t[1], t[1], t[2])      // (a0 + a1)(b0 + b1)
	fq.copy(c[0], t[3])           // c0 = β * v1 + v0
	fq.sub(c[1], t[1], t[0])      // c1 = (a0 + a1)(b0 + b1) - (v0+v1)
}

func (f *fq2) square(c, a *fe2) {
	fq, t := f.fq(), f.t
	// c0 = (a0 - a1) * (a0 - β * a1) + a0 * a1 + β * a0 * a1
	// c1 = 2 * a0 * a1
	fq.sub(t[0], a[0], a[1])      // v0 = a0 - a1
	f.mulByNonResidue(t[1], a[1]) // β * a1
	fq.sub(t[2], a[0], t[1])      // v3 = a0 -  β * a1
	fq.mul(t[1], a[0], a[1])      // v2 = a0 * a1
	fq.mul(t[0], t[0], t[2])      // v0 * v3
	fq.add(t[0], t[1], t[0])      // v0 = v0 * v3 + v2
	f.mulByNonResidue(t[2], t[1]) // β * v2
	fq.add(c[0], t[0], t[2])      // c0 = v0 + β * v2
	fq.double(c[1], t[1])         // c1 = 2*v2
}

func (f *fq2) inverse(c, a *fe2) bool {
	fq, t := f.fq(), f.t
	// c0 = a0 * (a0^2 - β * a1^2)^-1
	// c1 = a1 * (a0^2 - β * a1^2)^-1
	fq.square(t[0], a[0])         // v0 = a0^2
	fq.square(t[1], a[1])         // v1 = a1^2
	f.mulByNonResidue(t[1], t[1]) // β * v1
	fq.sub(t[1], t[0], t[1])      // v0 = v0 - β * v1
	if ok := fq.inverse(t[0], t[1]); !ok {
		f.copy(c, f.zero())
		return false
	} // v1 = v0^-1;
	fq.mul(c[0], a[0], t[0]) // c0 = a0 * v1
	fq.mul(t[0], a[1], t[0]) // a1 * v1
	fq.neg(c[1], t[0])       // c1 = -a1 * v1
	return true
}

func (f *fq2) exp(c, a *fe2, e *big.Int) {
	z := f.one()
	found := false
	for i := e.BitLen() - 1; i >= 0; i-- {
		if found {
			f.square(z, z)
		} else {
			found = e.Bit(i) == 1
		}
		if e.Bit(i) == 1 {
			f.mul(z, z, a)
		}
	}
	f.copy(c, z)
}

func (f *fq2) mulByFq(c, a *fe2, b fe) {
	fq := f.fq()
	fq.mul(c[0], a[0], b)
	fq.mul(c[1], a[1], b)
}

func (f *fq2) frobeniusMap(c, a *fe2, power uint) {
	fq := f.fq()
	fq.copy(c[0], a[0])
	fq.mul(c[1], a[1], f.frobeniusCoeffs[power%2])
}

func (fq2 *fq2) calculateFrobeniusCoeffs() bool {
	fq := fq2.fq()
	zero, one, two := big.NewInt(0), big.NewInt(1), big.NewInt(2)
	if fq2.frobeniusCoeffs == nil {
		fq2.frobeniusCoeffs = fq2.new()
	}
	power, rem := new(big.Int), new(big.Int)
	power.Sub(fq2.modulus(), one)
	power.DivMod(power, two, rem)
	if rem.Cmp(zero) != 0 {
		return false
	}
	fq.exp(fq2.frobeniusCoeffs[1], fq2.nonResidue, power)
	fq.copy(fq2.frobeniusCoeffs[0], fq.one)
	return true
}

func (fq2 *fq2) calculateFrobeniusCoeffsWithPrecomputation(f1 fe) {
	fq := fq2.fq()
	if fq2.frobeniusCoeffs == nil {
		fq2.frobeniusCoeffs = fq2.new()
	}
	fq.copy(fq2.frobeniusCoeffs[0], fq.one)
	fq.square(fq2.frobeniusCoeffs[1], f1)
}

func (f *fq2) fq() *fq {
	return f.f
}
