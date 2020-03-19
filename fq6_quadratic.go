package eip

import (
	"fmt"
	"io"
	"math/big"
)

type fe6Q [2]fe3

type fq6Q struct {
	f               *fq3
	nonResidue      *fe3
	t               []*fe3
	frobeniusCoeffs *[6]fe
}

func newFq6Quadratic(fq3 *fq3, nonResidueBuf []byte) (*fq6Q, error) {
	nonResidue := fq3.new()
	if nonResidueBuf != nil {
		var err error
		nonResidue, err = fq3.fromBytes(nonResidueBuf)
		if err != nil {
			return nil, err
		}
	}
	t := make([]*fe3, 4)
	for i := 0; i < 4; i++ {
		t[i] = fq3.new()
	}
	return &fq6Q{fq3, nonResidue, t, nil}, nil
}

func (fq6 *fq6Q) byteSize() int {
	fq3 := fq6.fq3()
	return fq3.byteSize() * 2
}

func (fq6 *fq6Q) new() *fe6Q {
	return fq6.zero()
}

func (fq6 *fq6Q) modulus() *big.Int {
	fq3 := fq6.fq3()
	return fq3.modulus()
}

func (fq6 *fq6Q) rand(r io.Reader) *fe6Q {
	fq3 := fq6.fq3()
	return &fe6Q{*fq3.rand(r), *fq3.rand(r)}
}

func (fq6 *fq6Q) fromBytes(in []byte) (*fe6Q, error) {
	fq3 := fq6.fq3()
	byteSize := fq3.byteSize()
	if len(in) != byteSize*2 {
		return nil, fmt.Errorf("input string should be larger than %d bytes", byteSize*2)
	}
	var err error
	c0, err := fq3.fromBytes(in[:byteSize])
	if err != nil {
		return nil, err
	}
	c1, err := fq3.fromBytes(in[byteSize:])
	if err != nil {
		return nil, err
	}
	return &fe6Q{*c0, *c1}, nil
}

func (fq6 *fq6Q) toBytes(a *fe6Q) []byte {
	fq3 := fq6.fq3()
	byteSize := fq3.byteSize()
	out := make([]byte, 2*byteSize)
	copy(out[:byteSize], fq3.toBytes(&a[0]))
	copy(out[byteSize:], fq3.toBytes(&a[1]))
	return out
}

func (fq6 *fq6Q) toString(a *fe6Q) string {
	fq3 := fq6.fq3()
	return fmt.Sprintf("%s\n%s", fq3.toString(&a[0]), fq3.toString(&a[1]))
}

func (fq6 *fq6Q) toStringNoTransform(a *fe6Q) string {
	fq3 := fq6.fq3()
	return fmt.Sprintf("%s\n%s", fq3.toStringNoTransform(&a[0]), fq3.toStringNoTransform(&a[1]))
}

func (fq6 *fq6Q) zero() *fe6Q {
	fq3 := fq6.fq3()
	return &fe6Q{*fq3.zero(), *fq3.zero()}
}

func (fq6 *fq6Q) one() *fe6Q {
	fq3, a := fq6.fq3(), fq6.new()
	fq3.copy(&a[0], fq3.one())
	return a
}

func (fq6 *fq6Q) isZero(a *fe6Q) bool {
	fq3 := fq6.fq3()
	return fq3.isZero(&a[0]) && fq3.isZero(&a[1])
}

func (fq6 *fq6Q) isOne(a *fe6Q) bool {
	fq3 := fq6.fq3()
	return fq3.isOne(&a[0]) && fq3.isZero(&a[1])
}

func (fq6 *fq6Q) equal(a, b *fe6Q) bool {
	fq3 := fq6.fq3()
	return fq3.equal(&a[0], &b[0]) && fq3.equal(&a[1], &b[1])
}

func (fq6 *fq6Q) copy(c, a *fe6Q) *fe6Q {
	fq3 := fq6.fq3()
	fq3.copy(&c[0], &a[0])
	fq3.copy(&c[1], &a[1])
	return c
}

func (fq6 *fq6Q) add(c, a, b *fe6Q) *fe6Q {
	fq := fq6.fq3()
	fq.add(&c[0], &a[0], &b[0])
	fq.add(&c[1], &a[1], &b[1])
	return c
}

func (fq6 *fq6Q) double(c, a *fe6Q) *fe6Q {
	fq3 := fq6.fq3()
	fq3.double(&c[0], &a[0])
	fq3.double(&c[1], &a[1])
	return c
}

func (fq6 *fq6Q) sub(c, a, b *fe6Q) *fe6Q {
	fq3 := fq6.fq3()
	fq3.sub(&c[0], &a[0], &b[0])
	fq3.sub(&c[1], &a[1], &b[1])
	return c
}

func (fq6 *fq6Q) neg(c, a *fe6Q) *fe6Q {
	fq3 := fq6.fq3()
	fq3.neg(&c[0], &a[0])
	fq3.neg(&c[1], &a[1])
	return c
}

func (fq6 *fq6Q) conjugate(c, a *fe6Q) *fe6Q {
	fq3 := fq6.fq3()
	fq6.copy(c, a)
	fq3.neg(&c[1], &a[1])
	return c
}

func (fq6 *fq6Q) mulByNonResidue(c, a *fe3) {
	fq, fq3 := fq6.fq(), fq6.fq3()
	o := fq3.new()
	fq.copy(o[2], a[1])
	fq.copy(o[1], a[0])
	fq3.mulByNonResidue(o[0], a[2])
	fq3.copy(c, o)
}

func (fq6 *fq6Q) mul(c, a, b *fe6Q) {
	fq3, t := fq6.fq3(), fq6.t
	// c0 = (a0 * b0) + β * (a1 * b1)
	// c1 = (a0 + a1) * (b0 + b1) - (a0 * b0 + a1 * b1)
	fq3.mul(t[1], &a[0], &b[0])     // v0 = a0 * b0
	fq3.mul(t[2], &a[1], &b[1])     // v1 = a1 * b1
	fq3.add(t[0], t[1], t[2])       // v0 + v1
	fq6.mulByNonResidue(t[2], t[2]) // β * v1
	fq3.add(t[3], t[1], t[2])       // β * v1 + v0
	fq3.add(t[1], &a[0], &a[1])     // a0 + a1
	fq3.add(t[2], &b[0], &b[1])     // b0 + b1
	fq3.mul(t[1], t[1], t[2])       // (a0 + a1)(b0 + b1)
	fq3.copy(&c[0], t[3])           // c0 = β * v1 + v0
	fq3.sub(&c[1], t[1], t[0])      // c1 = (a0 + a1)(b0 + b1) - (v0+v1)
}

func (fq6 *fq6Q) square(c, a *fe6Q) {
	fq3, t := fq6.fq3(), fq6.t
	// c0 = (a0 - a1) * (a0 - β * a1) + a0 * a1 + β * a0 * a1
	// c1 = 2 * a0 * a1
	fq3.sub(t[0], &a[0], &a[1])      // v0 = a0 - a1
	fq6.mulByNonResidue(t[1], &a[1]) // β * a1
	fq3.sub(t[2], &a[0], t[1])       // v3 = a0 -  β * a1
	fq3.mul(t[1], &a[0], &a[1])      // v2 = a0 * a1
	fq3.mul(t[0], t[0], t[2])        // v0 * v3
	fq3.add(t[0], t[1], t[0])        // v0 = v0 * v3 + v2
	fq6.mulByNonResidue(t[2], t[1])  // β * v2
	fq3.add(&c[0], t[0], t[2])       // c0 = v0 + β * v2
	fq3.double(&c[1], t[1])          // c1 = 2*v2
}

func (fq6 *fq6Q) inverse(c, a *fe6Q) bool {
	fq3, t := fq6.fq3(), fq6.t
	// c0 = a0 * (a0^2 - β * a1^2)^-1
	// c1 = a1 * (a0^2 - β * a1^2)^-1
	fq3.square(t[0], &a[0])                 // v0 = a0^2
	fq3.square(t[1], &a[1])                 // v1 = a1^2
	fq6.mulByNonResidue(t[1], t[1])         // β * v1
	fq3.sub(t[1], t[0], t[1])               // v0 = v0 - β * v1
	if ok := fq3.inverse(t[0], t[1]); !ok { // v1 = v0^-1
		fq6.copy(c, fq6.zero())
		return false
	}
	fq3.mul(&c[0], &a[0], t[0]) // c0 = a0 * v1
	fq3.mul(t[0], &a[1], t[0])  // a1 * v1
	fq3.neg(&c[1], t[0])        // c1 = -a1 * v1
	return true
}

func (fq6 *fq6Q) exp(c, a *fe6Q, e *big.Int) {
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

func (fq6 *fq6Q) frobeniusMap(c, a *fe6Q, power uint) {
	fq3 := fq6.fq3()
	fq6.copy(c, a)
	fq3.frobeniusMap(&c[0], &a[0], power)
	fq3.frobeniusMap(&c[1], &a[1], power)
	fq3.mulByFq(&c[1], &c[1], fq6.frobeniusCoeffs[power%6])
}

func (fq6 *fq6Q) calculateFrobeniusCoeffs() bool {
	fq, fq3 := fq6.fq(), fq6.fq3()
	if fq6.frobeniusCoeffs == nil {
		fq6.frobeniusCoeffs = new([6]fe)
		for i := 0; i < len(fq6.frobeniusCoeffs); i++ {
			fq6.frobeniusCoeffs[i] = fq.new()
		}
	}
	modulus := fq.modulus()
	fq.copy(fq6.frobeniusCoeffs[0], fq.one)
	qPower, rem, power := new(big.Int).Set(modulus), new(big.Int), new(big.Int)
	zero := new(big.Int)
	for i := 1; i <= 3; i++ {
		if i == 2 {
			qPower.Mul(qPower, modulus)
			continue
		}
		power.Sub(qPower, big.NewInt(1))
		power.DivMod(power, big.NewInt(6), rem)
		if rem.Cmp(zero) != 0 {
			return false
		}
		fq.exp(fq6.frobeniusCoeffs[i], fq3.nonResidue, power)
		qPower.Mul(qPower, modulus)
	}
	return true
}

func (fq6 *fq6Q) calculateFrobeniusCoeffsWithPrecomputation(f1 fe) {
	fq := fq6.fq()
	if fq6.frobeniusCoeffs == nil {
		fq6.frobeniusCoeffs = new([6]fe)
		for i := 0; i < len(fq6.frobeniusCoeffs); i++ {
			fq6.frobeniusCoeffs[i] = fq.new()
		}
	}
	fq.copy(fq6.frobeniusCoeffs[0], fq.one)
	fq.copy(fq6.frobeniusCoeffs[1], f1)
	fq.exp(fq6.frobeniusCoeffs[3], fq6.frobeniusCoeffs[1], big.NewInt(3))
}

func (fq6 *fq6Q) fq() *fq {
	return fq6.f.f
}

func (fq6 *fq6Q) fq3() *fq3 {
	return fq6.f
}
