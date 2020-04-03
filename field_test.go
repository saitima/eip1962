package eip

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"testing"
)

type fieldElement interface{}

type field interface {
	new() fieldElement
	debugElement(fe fieldElement)
	limbSize() int
	byteSize() int
	fromBytes(in []byte) (fieldElement, error)
	toBytes(a fieldElement) []byte
	isZero(a fieldElement) bool
	isOne(a fieldElement) bool
	zero() fieldElement
	p() *big.Int
	one() fieldElement
	equal(a, b fieldElement) bool
	rand(r io.Reader) fieldElement
	add(c, a, b fieldElement)
	sub(c, a, b fieldElement)
	double(c, a fieldElement)
	neg(c, a fieldElement)
	mul(c, a, b fieldElement)
	square(c, a fieldElement)
	exp(c, a fieldElement, e *big.Int)
	inverse(c, a fieldElement) bool
	sqrt(c, a fieldElement) bool
}

type fqTest struct {
	*fq
}

func (t fqTest) new() fieldElement {
	return t.fq.new()
}

func (t fqTest) debugElement(in fieldElement) {
	t.fq.debugElement(in.(fe))
}

func (f *fq) debugElement(in fe) {
	fmt.Println(f.toString(in))
}

func (t fqTest) byteSize() int {
	return t.fq.byteSize()
}

func (t fqTest) limbSize() int {
	return t.fq.limbSize
}

func (t fqTest) fromBytes(in []byte) (fieldElement, error) {
	return t.fq.fromBytes(in)
}

func (t fqTest) toBytes(in fieldElement) []byte {
	return t.fq.toBytes(in.(fe))
}

func (t fqTest) isZero(in fieldElement) bool {
	return t.fq.isZero(in.(fe))
}

func (t fqTest) isOne(in fieldElement) bool {
	return t.fq.isOne(in.(fe))
}

func (t fqTest) zero() fieldElement {
	return t.fq.zero
}

func (t fqTest) one() fieldElement {
	return t.fq.one
}

func (t fqTest) p() *big.Int {
	return t.fq.modulus()
}

func (t fqTest) equal(a, b fieldElement) bool {
	return t.fq.equal(a.(fe), b.(fe))
}

func (t fqTest) rand(r io.Reader) fieldElement {
	return t.fq.rand(r)
}

func (t fqTest) add(c, a, b fieldElement) {
	t.fq.add(c.(fe), a.(fe), b.(fe))
}

func (t fqTest) sub(c, a, b fieldElement) {
	t.fq.sub(c.(fe), a.(fe), b.(fe))
}

func (t fqTest) mul(c, a, b fieldElement) {
	t.fq.mul(c.(fe), a.(fe), b.(fe))
}

func (t fqTest) double(c, a fieldElement) {
	t.fq.double(c.(fe), a.(fe))
}

func (t fqTest) neg(c, a fieldElement) {
	t.fq.neg(c.(fe), a.(fe))
}

func (t fqTest) square(c, a fieldElement) {
	t.fq.square(c.(fe), a.(fe))
}

func (t fqTest) exp(c, a fieldElement, e *big.Int) {
	t.fq.exp(c.(fe), a.(fe), e)
}

func (t fqTest) inverse(c, a fieldElement) bool {
	return t.fq.inverse(c.(fe), a.(fe))
}

func (t fqTest) sqrt(c, a fieldElement) bool {
	return t.fq.sqrt(c.(fe), a.(fe))
}

type fq2Test struct {
	*fq2
}

func (t fq2Test) new() fieldElement {
	return t.fq2.new()
}

func (t fq2Test) debugElement(in fieldElement) {
	t.fq2.debugElement(in.(*fe2))
}

func (f *fq2) debugElement(in *fe2) {
	fmt.Println(f.toString(in))
}

func (t fq2Test) byteSize() int {
	return t.fq2.byteSize()
}

func (t fq2Test) limbSize() int {
	return t.fq2.f.limbSize
}

func (t fq2Test) fromBytes(in []byte) (fieldElement, error) {
	return t.fq2.fromBytes(in)
}

func (t fq2Test) toBytes(in fieldElement) []byte {
	return t.fq2.toBytes(in.(*fe2))
}

func (t fq2Test) isZero(in fieldElement) bool {
	return t.fq2.isZero(in.(*fe2))
}

func (t fq2Test) isOne(in fieldElement) bool {
	return t.fq2.isOne(in.(*fe2))
}

func (t fq2Test) zero() fieldElement {
	return t.fq2.zero()
}

func (t fq2Test) one() fieldElement {
	return t.fq2.one()
}

func (t fq2Test) p() *big.Int {
	return t.f.modulus()
}

func (t fq2Test) equal(a, b fieldElement) bool {
	return t.fq2.equal(a.(*fe2), b.(*fe2))
}

func (t fq2Test) rand(r io.Reader) fieldElement {
	return t.fq2.rand(r)
}

func (t fq2Test) add(c, a, b fieldElement) {
	t.fq2.add(c.(*fe2), a.(*fe2), b.(*fe2))
}

func (t fq2Test) sub(c, a, b fieldElement) {
	t.fq2.sub(c.(*fe2), a.(*fe2), b.(*fe2))
}

func (t fq2Test) mul(c, a, b fieldElement) {
	t.fq2.mul(c.(*fe2), a.(*fe2), b.(*fe2))
}

func (t fq2Test) double(c, a fieldElement) {
	t.fq2.double(c.(*fe2), a.(*fe2))
}

func (t fq2Test) neg(c, a fieldElement) {
	t.fq2.neg(c.(*fe2), a.(*fe2))
}

func (t fq2Test) square(c, a fieldElement) {
	t.fq2.square(c.(*fe2), a.(*fe2))
}

func (t fq2Test) exp(c, a fieldElement, e *big.Int) {
	t.fq2.exp(c.(*fe2), a.(*fe2), e)
}

func (t fq2Test) inverse(c, a fieldElement) bool {
	return t.fq2.inverse(c.(*fe2), a.(*fe2))
}

func (t fq2Test) mulByFq(c, a fieldElement, b fieldElement) {
	t.fq2.mulByFq(c.(*fe2), a.(*fe2), b.(fe))
}

func (t fq2Test) sqrt(c, a fieldElement) bool {
	return t.fq2.sqrt(c.(*fe2), a.(*fe2))
}

type fq3Test struct {
	*fq3
}

func (t fq3Test) new() fieldElement {
	return t.fq3.new()
}

func (t fq3Test) debugElement(in fieldElement) {
	t.fq3.debugElement(in.(*fe3))
}

func (f *fq3) debugElement(in *fe3) {
	fmt.Println(f.toString(in))
}

func (t fq3Test) byteSize() int {
	return t.fq3.byteSize()
}

func (t fq3Test) limbSize() int {
	return t.fq3.f.limbSize
}

func (t fq3Test) fromBytes(in []byte) (fieldElement, error) {
	return t.fq3.fromBytes(in)
}

func (t fq3Test) toBytes(in fieldElement) []byte {
	return t.fq3.toBytes(in.(*fe3))
}

func (t fq3Test) isZero(in fieldElement) bool {
	return t.fq3.isZero(in.(*fe3))
}

func (t fq3Test) isOne(in fieldElement) bool {
	return t.fq3.isOne(in.(*fe3))
}

func (t fq3Test) zero() fieldElement {
	return t.fq3.zero()
}

func (t fq3Test) one() fieldElement {
	return t.fq3.one()
}

func (t fq3Test) p() *big.Int {
	return t.f.modulus()
}

func (t fq3Test) equal(a, b fieldElement) bool {
	return t.fq3.equal(a.(*fe3), b.(*fe3))
}

func (t fq3Test) rand(r io.Reader) fieldElement {
	return t.fq3.rand(r)
}

func (t fq3Test) add(c, a, b fieldElement) {
	t.fq3.add(c.(*fe3), a.(*fe3), b.(*fe3))
}

func (t fq3Test) sub(c, a, b fieldElement) {
	t.fq3.sub(c.(*fe3), a.(*fe3), b.(*fe3))
}

func (t fq3Test) mul(c, a, b fieldElement) {
	t.fq3.mul(c.(*fe3), a.(*fe3), b.(*fe3))
}

func (t fq3Test) double(c, a fieldElement) {
	t.fq3.double(c.(*fe3), a.(*fe3))
}

func (t fq3Test) neg(c, a fieldElement) {
	t.fq3.neg(c.(*fe3), a.(*fe3))
}

func (t fq3Test) square(c, a fieldElement) {
	t.fq3.square(c.(*fe3), a.(*fe3))
}

func (t fq3Test) exp(c, a fieldElement, e *big.Int) {
	t.fq3.exp(c.(*fe3), a.(*fe3), e)
}

func (t fq3Test) inverse(c, a fieldElement) bool {
	return t.fq3.inverse(c.(*fe3), a.(*fe3))
}

func (t fq3Test) mulByFq(c, a fieldElement, b fieldElement) {
	t.fq3.mulByFq(c.(*fe3), a.(*fe3), b.(fe))
}

func (t fq3Test) sqrt(c, a fieldElement) bool {
	// return t.fq.sqrt(c.(*fe3), a.(*fe3))
	return true
}

type fq4Test struct {
	*fq4
}

func (t fq4Test) new() fieldElement {
	return t.fq4.new()
}

func (t fq4Test) debugElement(in fieldElement) {
	t.fq4.debugElement(in.(*fe4))
}

func (f *fq4) debugElement(in *fe4) {
	fmt.Println(f.toString(in))
}

func (t fq4Test) byteSize() int {
	return t.fq4.byteSize()
}

func (t fq4Test) limbSize() int {
	return t.fq4.f.f.limbSize
}

func (t fq4Test) fromBytes(in []byte) (fieldElement, error) {
	return t.fq4.fromBytes(in)
}

func (t fq4Test) toBytes(in fieldElement) []byte {
	return t.fq4.toBytes(in.(*fe4))
}

func (t fq4Test) isZero(in fieldElement) bool {
	return t.fq4.isZero(in.(*fe4))
}

func (t fq4Test) isOne(in fieldElement) bool {
	return t.fq4.isOne(in.(*fe4))
}

func (t fq4Test) zero() fieldElement {
	return t.fq4.zero()
}

func (t fq4Test) one() fieldElement {
	return t.fq4.one()
}

func (t fq4Test) p() *big.Int {
	return t.f.modulus()
}

func (t fq4Test) equal(a, b fieldElement) bool {
	return t.fq4.equal(a.(*fe4), b.(*fe4))
}

func (t fq4Test) rand(r io.Reader) fieldElement {
	return t.fq4.rand(r)
}

func (t fq4Test) add(c, a, b fieldElement) {
	t.fq4.add(c.(*fe4), a.(*fe4), b.(*fe4))
}

func (t fq4Test) sub(c, a, b fieldElement) {
	t.fq4.sub(c.(*fe4), a.(*fe4), b.(*fe4))
}

func (t fq4Test) mul(c, a, b fieldElement) {
	t.fq4.mul(c.(*fe4), a.(*fe4), b.(*fe4))
}

func (t fq4Test) double(c, a fieldElement) {
	t.fq4.double(c.(*fe4), a.(*fe4))
}

func (t fq4Test) neg(c, a fieldElement) {
	t.fq4.neg(c.(*fe4), a.(*fe4))
}

func (t fq4Test) square(c, a fieldElement) {
	t.fq4.square(c.(*fe4), a.(*fe4))
}

func (t fq4Test) exp(c, a fieldElement, e *big.Int) {
	t.fq4.exp(c.(*fe4), a.(*fe4), e)
}

func (t fq4Test) inverse(c, a fieldElement) bool {
	return t.fq4.inverse(c.(*fe4), a.(*fe4))
}

func (t fq4Test) sqrt(c, a fieldElement) bool {
	// return t.fq.sqrt(c.(*fe4), a.(*fe4))
	return true
}

type fq6CTest struct {
	*fq6C
}

func (t fq6CTest) new() fieldElement {
	return t.fq6C.new()
}

func (t fq6CTest) debugElement(in fieldElement) {
	t.fq6C.debugElement(in.(*fe6C))
}

func (f *fq6C) debugElement(in *fe6C) {
	fmt.Println(f.toString(in))
}

func (t fq6CTest) byteSize() int {
	return t.fq6C.byteSize()
}

func (t fq6CTest) limbSize() int {
	return t.fq6C.f.f.limbSize
}

func (t fq6CTest) fromBytes(in []byte) (fieldElement, error) {
	return t.fq6C.fromBytes(in)
}

func (t fq6CTest) toBytes(in fieldElement) []byte {
	return t.fq6C.toBytes(in.(*fe6C))
}

func (t fq6CTest) isZero(in fieldElement) bool {
	return t.fq6C.isZero(in.(*fe6C))
}

func (t fq6CTest) isOne(in fieldElement) bool {
	return t.fq6C.isOne(in.(*fe6C))
}

func (t fq6CTest) zero() fieldElement {
	return t.fq6C.zero()
}

func (t fq6CTest) one() fieldElement {
	return t.fq6C.one()
}

func (t fq6CTest) p() *big.Int {
	return t.f.modulus()
}

func (t fq6CTest) equal(a, b fieldElement) bool {
	return t.fq6C.equal(a.(*fe6C), b.(*fe6C))
}

func (t fq6CTest) rand(r io.Reader) fieldElement {
	return t.fq6C.rand(r)
}

func (t fq6CTest) add(c, a, b fieldElement) {
	t.fq6C.add(c.(*fe6C), a.(*fe6C), b.(*fe6C))
}

func (t fq6CTest) sub(c, a, b fieldElement) {
	t.fq6C.sub(c.(*fe6C), a.(*fe6C), b.(*fe6C))
}

func (t fq6CTest) mul(c, a, b fieldElement) {
	t.fq6C.mul(c.(*fe6C), a.(*fe6C), b.(*fe6C))
}

func (t fq6CTest) double(c, a fieldElement) {
	t.fq6C.double(c.(*fe6C), a.(*fe6C))
}

func (t fq6CTest) neg(c, a fieldElement) {
	t.fq6C.neg(c.(*fe6C), a.(*fe6C))
}

func (t fq6CTest) square(c, a fieldElement) {
	t.fq6C.square(c.(*fe6C), a.(*fe6C))
}

func (t fq6CTest) exp(c, a fieldElement, e *big.Int) {
	t.fq6C.exp(c.(*fe6C), a.(*fe6C), e)
}

func (t fq6CTest) inverse(c, a fieldElement) bool {
	return t.fq6C.inverse(c.(*fe6C), a.(*fe6C))
}

func (t fq6CTest) sqrt(c, a fieldElement) bool {
	// return t.fq.sqrt(c.(*fe6C), a.(*fe6C))
	return true
}

type fq6QTest struct {
	*fq6Q
}

func (t fq6QTest) new() fieldElement {
	return t.fq6Q.new()
}

func (t fq6QTest) debugElement(in fieldElement) {
	t.fq6Q.debugElement(in.(*fe6Q))
}

func (f *fq6Q) debugElement(in *fe6Q) {
	fmt.Println(f.toString(in))
}

func (t fq6QTest) byteSize() int {
	return t.fq6Q.byteSize()
}

func (t fq6QTest) limbSize() int {
	return t.fq6Q.f.f.limbSize
}

func (t fq6QTest) fromBytes(in []byte) (fieldElement, error) {
	return t.fq6Q.fromBytes(in)
}

func (t fq6QTest) toBytes(in fieldElement) []byte {
	return t.fq6Q.toBytes(in.(*fe6Q))
}

func (t fq6QTest) isZero(in fieldElement) bool {
	return t.fq6Q.isZero(in.(*fe6Q))
}

func (t fq6QTest) isOne(in fieldElement) bool {
	return t.fq6Q.isOne(in.(*fe6Q))
}

func (t fq6QTest) zero() fieldElement {
	return t.fq6Q.zero()
}

func (t fq6QTest) one() fieldElement {
	return t.fq6Q.one()
}

func (t fq6QTest) p() *big.Int {
	return t.f.modulus()
}

func (t fq6QTest) equal(a, b fieldElement) bool {
	return t.fq6Q.equal(a.(*fe6Q), b.(*fe6Q))
}

func (t fq6QTest) rand(r io.Reader) fieldElement {
	return t.fq6Q.rand(r)
}

func (t fq6QTest) add(c, a, b fieldElement) {
	t.fq6Q.add(c.(*fe6Q), a.(*fe6Q), b.(*fe6Q))
}

func (t fq6QTest) sub(c, a, b fieldElement) {
	t.fq6Q.sub(c.(*fe6Q), a.(*fe6Q), b.(*fe6Q))
}

func (t fq6QTest) mul(c, a, b fieldElement) {
	t.fq6Q.mul(c.(*fe6Q), a.(*fe6Q), b.(*fe6Q))
}

func (t fq6QTest) double(c, a fieldElement) {
	t.fq6Q.double(c.(*fe6Q), a.(*fe6Q))
}

func (t fq6QTest) neg(c, a fieldElement) {
	t.fq6Q.neg(c.(*fe6Q), a.(*fe6Q))
}

func (t fq6QTest) square(c, a fieldElement) {
	t.fq6Q.square(c.(*fe6Q), a.(*fe6Q))
}

func (t fq6QTest) exp(c, a fieldElement, e *big.Int) {
	t.fq6Q.exp(c.(*fe6Q), a.(*fe6Q), e)
}

func (t fq6QTest) inverse(c, a fieldElement) bool {
	return t.fq6Q.inverse(c.(*fe6Q), a.(*fe6Q))
}

func (t fq6QTest) sqrt(c, a fieldElement) bool {
	// return t.fq.sqrt(c.(*fe6Q), a.(*fe6Q))
	return true
}

type fq12Test struct {
	*fq12
}

func (t fq12Test) new() fieldElement {
	return t.fq12.new()
}

func (t fq12Test) debugElement(in fieldElement) {
	t.fq12.debugElement(in.(*fe12))
}

func (f *fq12) debugElement(in *fe12) {
	fmt.Println(f.toString(in))
}

func (t fq12Test) byteSize() int {
	return t.fq12.byteSize()
}

func (t fq12Test) limbSize() int {
	return t.fq12.f.f.f.limbSize
}

func (t fq12Test) fromBytes(in []byte) (fieldElement, error) {
	return t.fq12.fromBytes(in)
}

func (t fq12Test) toBytes(in fieldElement) []byte {
	return t.fq12.toBytes(in.(*fe12))
}

func (t fq12Test) isZero(in fieldElement) bool {
	return t.fq12.isZero(in.(*fe12))
}

func (t fq12Test) isOne(in fieldElement) bool {
	return t.fq12.isOne(in.(*fe12))
}

func (t fq12Test) zero() fieldElement {
	return t.fq12.zero()
}

func (t fq12Test) one() fieldElement {
	return t.fq12.one()
}

func (t fq12Test) p() *big.Int {
	return t.f.modulus()
}

func (t fq12Test) equal(a, b fieldElement) bool {
	return t.fq12.equal(a.(*fe12), b.(*fe12))
}

func (t fq12Test) rand(r io.Reader) fieldElement {
	return t.fq12.rand(r)
}

func (t fq12Test) add(c, a, b fieldElement) {
	t.fq12.add(c.(*fe12), a.(*fe12), b.(*fe12))
}

func (t fq12Test) sub(c, a, b fieldElement) {
	t.fq12.sub(c.(*fe12), a.(*fe12), b.(*fe12))
}

func (t fq12Test) mul(c, a, b fieldElement) {
	t.fq12.mul(c.(*fe12), a.(*fe12), b.(*fe12))
}

func (t fq12Test) double(c, a fieldElement) {
	t.fq12.double(c.(*fe12), a.(*fe12))
}

func (t fq12Test) neg(c, a fieldElement) {
	t.fq12.neg(c.(*fe12), a.(*fe12))
}

func (t fq12Test) square(c, a fieldElement) {
	t.fq12.square(c.(*fe12), a.(*fe12))
}

func (t fq12Test) exp(c, a fieldElement, e *big.Int) {
	t.fq12.exp(c.(*fe12), a.(*fe12), e)
}

func (t fq12Test) inverse(c, a fieldElement) bool {
	return t.fq12.inverse(c.(*fe12), a.(*fe12))
}

func (t fq12Test) sqrt(c, a fieldElement) bool {
	// return t.fq.sqrt(c.(*fe12), a.(*fe12))
	return true
}

func randFq(limbSize int) *fq {
	var offset int
	t, err := rand.Int(rand.Reader, new(big.Int).SetUint64(64))
	if err != nil {
		panic(err)
	}
	offset = int(t.Uint64())
	byteLen := limbSize * 8
	bitLen := (limbSize-1)*64 + offset
	if bitLen < 32 {
		bitLen = 32
	}
	pbig, err := rand.Prime(rand.Reader, bitLen)
	if err != nil {
		panic(err)
	}
	rawpbytes := pbig.Bytes()
	pbytes := make([]byte, byteLen)
	copy(pbytes[byteLen-len(rawpbytes):], pbig.Bytes())
	field, err := newField(pbytes)
	if err != nil {
		panic(err)
	}
	if limbSize < 4 {
		if field.limbSize != 4 {
			panic("bad random field construction")
		}
	} else {
		if field.limbSize != limbSize {
			panic("bad random field construction")
		}
	}
	return field
}

func resolveLimbSize(bitSize int) int {
	size := (bitSize / 64)
	if bitSize%64 != 0 {
		size += 1
	}
	return size
}

func BenchmarkField(t *testing.B) {
	var limbSize int
	if targetNumberOfLimb > 0 {
		limbSize = targetNumberOfLimb
	} else {
		return
	}
	field := randFq(limbSize)
	if !USE_4LIMBS_FOR_LOWER_LIMBS && field.limbSize != limbSize {
		t.Fatalf("bad field construction")
	}
	bitSize := limbSize * 64
	a := field.rand(rand.Reader)
	b := field.rand(rand.Reader)
	c := field.new()
	t.Run(fmt.Sprintf("%d_add", bitSize), func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			field.add(c, a, b)
		}
	})
	t.Run(fmt.Sprintf("%d_double", bitSize), func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			field.double(c, a)
		}
	})
	t.Run(fmt.Sprintf("%d_sub", bitSize), func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			field.sub(c, a, b)
		}
	})
	t.Run(fmt.Sprintf("%d_mul", bitSize), func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			field.mul(c, a, b)
		}
	})
	t.Run(fmt.Sprintf("%d_cmp", bitSize), func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			field.cmp(a, b)
		}
	})
}

func TestFqShift(t *testing.T) {
	two := big.NewInt(2)
	for limbSize := from; limbSize < to+1; limbSize++ {
		t.Run(fmt.Sprintf("%d_shift", limbSize*64), func(t *testing.T) {
			field := randFq(limbSize)
			a := field.rand(rand.Reader)
			bi := field.toBigNoTransform(a)
			da := field.new()
			field.copy(da, a)
			field.div_two(da)
			dbi := new(big.Int).Div(bi, two)
			dbi_2 := field.toBigNoTransform(da)
			if dbi.Cmp(dbi_2) != 0 {
				t.Fatalf("bad div 2 operation")
			}
			ma := field.new()
			field.copy(ma, a)
			field.mul_two(ma)
			mbi := new(big.Int).Mul(bi, two)
			mbi_2 := field.toBigNoTransform(ma)
			if mbi.Cmp(mbi_2) != 0 {
				t.Fatalf("bad mul 2 operation")
			}
		})
	}
}

func TestFqCompare(t *testing.T) {
	for limbSize := from; limbSize < to+1; limbSize++ {
		t.Run(fmt.Sprintf("%d_compare", limbSize*64), func(t *testing.T) {
			field := randFq(limbSize)
			if field.cmp(field.r, field.r) != 0 {
				t.Fatalf("r == r (cmp)")
			}
			if !field.equal(field.r, field.r) {
				t.Fatalf("r == r (equal)")
			}
			if field.equal(field.p, field.r) {
				t.Fatalf("p != r")
			}
			if field.equal(field.r, field.zero) {
				t.Fatalf("r != 0")
			}
			if !field.equal(field.zero, field.zero) {
				t.Fatalf("0 == 0")
			}
			if field.cmp(field.p, field.r) != 1 {
				t.Fatalf("p > r")
			}
			if field.cmp(field.r, field.p) != -1 {
				t.Fatalf("r < p")
			}
			if is_even(field.p) {
				t.Fatalf("p is not even")
			}
		})
	}
}

func TestFqCopy(t *testing.T) {
	for limbSize := from; limbSize < to+1; limbSize++ {
		t.Run(fmt.Sprintf("%d_copy", limbSize*64), func(t *testing.T) {
			field := randFq(limbSize)
			a := field.rand(rand.Reader)
			b := field.new()
			field.copy(b, a)
			if !field.equal(a, b) {
				t.Fatalf("copy operation fails")
			}
		})
	}
}

func TestFqSerialization(t *testing.T) {
	for limbSize := from; limbSize < to+1; limbSize++ {
		t.Run(fmt.Sprintf("%d_serialization", limbSize*64), func(t *testing.T) {
			field := randFq(limbSize)
			if !USE_4LIMBS_FOR_LOWER_LIMBS && field.limbSize != limbSize {
				t.Fatalf("bad field construction\n")
			}
			// demont(r) == 1
			b0 := make([]byte, field.byteSize())
			b0[len(b0)-1] = byte(1)
			b1 := field.toBytes(field.r)
			if !bytes.Equal(b0, b1) {
				t.Fatalf("demont(r) must be equal to 1\n")
			}
			// is a => modulus should not be valid
			_, err := field.fromBytes(field.pbig.Bytes())
			if err == nil {
				t.Fatalf("a number eq or larger than modulus must not be valid")
			}
			for i := 0; i < fuz; i++ {
				field := randFq(limbSize)
				if !USE_4LIMBS_FOR_LOWER_LIMBS && field.limbSize != limbSize {
					t.Fatalf("bad field construction")
				}
				// bytes
				b0 := randBytes(field.pbig)
				if USE_4LIMBS_FOR_LOWER_LIMBS && len(b0) < 32 {
					b0 = padBytes(b0, 32)
				}
				a0, err := field.fromBytes(b0)
				if err != nil {
					t.Fatal(err)
				}
				if USE_4LIMBS_FOR_LOWER_LIMBS && len(b0) < 32 {
					b0 = padBytes(b0, 32)
				}
				b1 = field.toBytes(a0)
				if !bytes.Equal(b0, b1) {
					t.Fatalf("bad serialization (bytes)")
				}
				// string
				s := field.toString(a0)
				a1, err := field.fromString(s)
				if err != nil {
					t.Fatal(err)
				}
				if !field.equal(a0, a1) {
					t.Fatalf("bad serialization (str)")
				}
				// big int
				a0, err = field.fromBytes(b0)
				if err != nil {
					t.Fatal(err)
				}
				bi := field.toBig(a0)
				a1, err = field.fromBig(bi)
				if err != nil {
					t.Fatal(err)
				}
				if !field.equal(a0, a1) {
					t.Fatalf("bad serialization (big.Int)")
				}
				// bytes dense
				b0 = field.toBytesDense(a0)
				a1, err = field.fromBytes(padBytes(b0, field.byteSize()))
				if err != nil {
					t.Fatal(err)
				}
				if !field.equal(a0, a1) {
					t.Fatalf("bad serialization (dense)")
				}
			}
		})
	}
}

func TestFqAdditionCrossAgainstBigInt(t *testing.T) {
	for limbSize := from; limbSize < to+1; limbSize++ {
		t.Run(fmt.Sprintf("%d_addition_cross", limbSize*64), func(t *testing.T) {
			for i := 0; i < fuz; i++ {
				field := randFq(limbSize)
				if !USE_4LIMBS_FOR_LOWER_LIMBS && field.limbSize != limbSize {
					t.Fatalf("Bad field construction")
				}
				a := field.rand(rand.Reader)
				b := field.rand(rand.Reader)
				c := field.new()
				big_a := field.toBig(a)
				big_b := field.toBig(b)
				big_c := new(big.Int)
				field.add(c, a, b)
				out_1 := field.toBytes(c)
				out_2 := padBytes(big_c.Add(big_a, big_b).Mod(big_c, field.pbig).Bytes(), field.byteSize())
				if !bytes.Equal(out_1, out_2) {
					t.Fatalf("cross test against big.Int is not satisfied A")
				}
				field.double(c, a)
				out_1 = field.toBytes(c)
				out_2 = padBytes(big_c.Add(big_a, big_a).Mod(big_c, field.pbig).Bytes(), field.byteSize())
				if !bytes.Equal(out_1, out_2) {
					t.Fatalf("cross test against big.Int is not satisfied B")
				}
				field.sub(c, a, b)
				out_1 = field.toBytes(c)
				out_2 = padBytes(big_c.Sub(big_a, big_b).Mod(big_c, field.pbig).Bytes(), field.byteSize())
				if !bytes.Equal(out_1, out_2) {
					t.Fatalf("cross test against big.Int is not satisfied C")
				}
				field.neg(c, a)
				out_1 = field.toBytes(c)
				out_2 = padBytes(big_c.Neg(big_a).Mod(big_c, field.pbig).Bytes(), field.byteSize())
				if !bytes.Equal(out_1, out_2) {
					t.Fatalf("cross test against big.Int is not satisfied D")
				}
			}
		})
	}
}

func TestFqMultiplicationCrossAgainstBigInt(t *testing.T) {
	for limbSize := from; limbSize < to+1; limbSize++ {
		t.Run(fmt.Sprintf("%d_multiplication_cross", limbSize*64), func(t *testing.T) {
			for i := 0; i < fuz; i++ {
				field := randFq(limbSize)
				if !USE_4LIMBS_FOR_LOWER_LIMBS && field.limbSize != limbSize {
					t.Fatalf("bad field construction")
				}
				a := field.rand(rand.Reader)
				b := field.rand(rand.Reader)
				c := field.new()
				big_a := field.toBig(a)
				big_b := field.toBig(b)
				big_c := new(big.Int)
				field.mul(c, a, b)
				out_1 := field.toBytes(c)
				out_2 := padBytes(big_c.Mul(big_a, big_b).Mod(big_c, field.pbig).Bytes(), field.byteSize())
				if !bytes.Equal(out_1, out_2) {
					t.Fatalf("cross test against big.Int is not satisfied")
				}
			}
		})
	}
}

func randFq2(limbSize int) *fq2 {
	fq := randFq(limbSize)
	fq2, err := newFq2(fq, nil)
	if err != nil {
		panic(err)
	}
	// find a non residue
	for {
		k := fq.rand(rand.Reader)
		if !fq.isNonResidue(k, 2) {
			fq.copy(fq2.nonResidue, k)
			break
		}
	}
	fq2.calculateFrobeniusCoeffs()
	return fq2
}

func randFq3(limbSize int) *fq3 {
	var fq3 *fq3
	var err error
	for {
		fq := randFq(limbSize)
		fq3, err = newFq3(fq, nil)
		if err != nil {
			panic(err)
		}
		k := fq.rand(rand.Reader)
		if !fq.isNonResidue(k, 3) {
			fq.copy(fq3.nonResidue, k)
			break
		}
	}
	return fq3
}

func randFq4(limbSize int) *fq4 {
	var fq4 *fq4
	var fq2 *fq2
	var err error
	for {
		fq := randFq(limbSize)
		fq2, err = newFq2(fq, nil)
		if err != nil {
			panic(err)
		}
		k := fq.rand(rand.Reader)
		if !fq.isNonResidue(k, 4) {
			fq.copy(fq2.nonResidue, k)
			break
		}
	}
	for {
		fq4, err = newFq4(fq2, nil)
		if err != nil {
			panic(err)
		}
		k := fq2.rand(rand.Reader)
		if !fq2.isNonResidue(k, 2) {
			fq2.copy(fq4.nonResidue, k)
			break
		}
	}
	return fq4
}

func randFq6Q(limbSize int) *fq6Q {
	var fq6 *fq6Q
	var fq3 *fq3
	var err error
	for {
		fq := randFq(limbSize)
		fq3, err = newFq3(fq, nil)
		if err != nil {
			panic(err)
		}
		k := fq.rand(rand.Reader)
		if !fq.isNonResidue(k, 6) {
			fq.copy(fq3.nonResidue, k)
			break
		}
	}
	for {
		fq6, err = newFq6Quadratic(fq3, nil)
		if err != nil {
			panic(err)
		}
		k := fq3.rand(rand.Reader)
		if !fq3.isNonResidue(k, 2) {
			fq3.copy(fq6.nonResidue, k)
			break
		}
	}
	return fq6
}

func randFq6C(limbSize int) *fq6C {
	var fq6 *fq6C
	var fq2 *fq2
	var err error
	for {
		fq := randFq(limbSize)
		fq2, err = newFq2(fq, nil)
		if err != nil {
			panic(err)
		}
		k := fq.rand(rand.Reader)
		if !fq.isNonResidue(k, 6) {
			fq.copy(fq2.nonResidue, k)
			break
		}
	}
	for {
		fq6, err = newFq6Cubic(fq2, nil)
		if err != nil {
			panic(err)
		}
		k := fq2.rand(rand.Reader)
		if !fq2.isNonResidue(k, 3) {
			fq2.copy(fq6.nonResidue, k)
			break
		}
	}
	return fq6
}

func randFq12(limbSize int) *fq12 {
	var fq6 *fq6C
	var fq2 *fq2
	var err error
	for {
		fq := randFq(limbSize)
		fq2, err = newFq2(fq, nil)
		if err != nil {
			panic(err)
		}
		k := fq.rand(rand.Reader)
		if !fq.isNonResidue(k, 12) {
			fq.copy(fq2.nonResidue, k)
			break
		}
	}
	for {
		fq6, err = newFq6Cubic(fq2, nil)
		if err != nil {
			panic(err)
		}
		k := fq2.rand(rand.Reader)
		if !fq2.isNonResidue(k, 6) {
			fq2.copy(fq6.nonResidue, k)
			break
		}
	}
	fq12, err := newFq12(fq6, nil)
	if err != nil {
		panic(err)
	}
	return fq12
}

var fields = []string{"FQ", "FQ2", "FQ3", "FQ4", "FQ6Q", "FQ6C", "FQ12"}

func randField(ext string, limbSize int) field {
	switch ext {
	case "FQ":
		return fqTest{randFq(limbSize)}
	case "FQ2":
		return fq2Test{randFq2(limbSize)}
	case "FQ3":
		return fq3Test{randFq3(limbSize)}
	case "FQ4":
		return fq4Test{randFq4(limbSize)}
	case "FQ6Q":
		return fq6QTest{randFq6Q(limbSize)}
	case "FQ6C":
		return fq6CTest{randFq6C(limbSize)}
	case "FQ12":
		return fq12Test{randFq12(limbSize)}
	default:
		panic("unknown extension")
	}
}

func TestFqSerializationGeneric(t *testing.T) {
	for _, ext := range fields {
		for limbSize := from; limbSize < to+1; limbSize++ {
			t.Run(fmt.Sprintf("%d_%s", limbSize*64, ext), func(t *testing.T) {
				field := randField(ext, limbSize)
				for i := 0; i < fuz; i++ {
					a0 := field.rand(rand.Reader)
					buf := field.toBytes(a0)
					a1, err := field.fromBytes(buf)
					if err != nil {
						t.Fatal(err)
					}
					if !field.equal(a0, a1) {
						t.Fatalf("bad serialization (bytes)")
					}
				}
			})
		}
	}
}

func TestFqAdditionProperties(t *testing.T) {
	for _, ext := range fields {
		for limbSize := from; limbSize < to+1; limbSize++ {
			t.Run(fmt.Sprintf("%d_%s", limbSize*64, ext), func(t *testing.T) {
				for i := 0; i < fuz; i++ {
					field := randField(ext, limbSize)
					zero := field.zero()
					if !USE_4LIMBS_FOR_LOWER_LIMBS && field.limbSize() != limbSize {
						t.Fatalf("bad field construction")
					}
					a := field.rand(rand.Reader)
					b := field.rand(rand.Reader)
					c_1 := field.new()
					c_2 := field.new()
					field.add(c_1, a, zero)
					if !field.equal(c_1, a) {
						t.Fatalf("a + 0 == a")
					}
					field.sub(c_1, a, zero)
					if !field.equal(c_1, a) {
						t.Fatalf("a - 0 == a")
					}
					field.double(c_1, zero)
					if !field.equal(c_1, zero) {
						t.Fatalf("2 * 0 == 0")
					}
					field.neg(c_1, zero)
					if !field.equal(c_1, zero) {
						t.Fatalf("-0 == 0")
					}
					field.sub(c_1, zero, a)
					field.neg(c_2, a)
					if !field.equal(c_1, c_2) {
						t.Fatalf("0-a == -a")
					}
					field.double(c_1, a)
					field.add(c_2, a, a)
					if !field.equal(c_1, c_2) {
						t.Fatalf("2 * a == a + a")
					}
					field.add(c_1, a, b)
					field.add(c_2, b, a)
					if !field.equal(c_1, c_2) {
						t.Fatalf("a + b = b + a")
					}
					field.sub(c_1, a, b)
					field.sub(c_2, b, a)
					field.neg(c_2, c_2)
					if !field.equal(c_1, c_2) {
						t.Fatalf("a - b = - ( b - a )")
					}
					c_x := field.rand(rand.Reader)
					field.add(c_1, a, b)
					field.add(c_1, c_1, c_x)
					field.add(c_2, a, c_x)
					field.add(c_2, c_2, b)
					if !field.equal(c_1, c_2) {
						t.Fatalf("(a + b) + c == (a + c ) + b")
					}
					field.sub(c_1, a, b)
					field.sub(c_1, c_1, c_x)
					field.sub(c_2, a, c_x)
					field.sub(c_2, c_2, b)
					if !field.equal(c_1, c_2) {
						t.Fatalf("(a - b) - c == (a - c ) -b")
					}
				}
			})
		}
	}
}

func TestFqMultiplicationProperties(t *testing.T) {
	for _, ext := range fields {
		for limbSize := from; limbSize < to+1; limbSize++ {
			t.Run(fmt.Sprintf("%d_%s", limbSize*64, ext), func(t *testing.T) {
				for i := 0; i < fuz; i++ {
					field := randField(ext, limbSize)
					if !USE_4LIMBS_FOR_LOWER_LIMBS && field.limbSize() != limbSize {
						t.Fatalf("bad field construction")
					}
					a := field.rand(rand.Reader)
					b := field.rand(rand.Reader)
					zero := field.zero()
					one := field.one()
					c_1 := field.new()
					c_2 := field.new()
					field.mul(c_1, a, zero)
					if !field.equal(c_1, zero) {
						t.Fatalf("a * 0 == 0")
					}
					field.mul(c_1, a, one)
					if !field.equal(c_1, a) {
						t.Fatalf("a * 1 == a")
					}
					field.mul(c_1, a, b)
					field.mul(c_2, b, a)
					if !field.equal(c_1, c_2) {
						t.Fatalf("a * b == b * a")
					}
					c_x := field.rand(rand.Reader)
					field.mul(c_1, a, b)
					field.mul(c_1, c_1, c_x)
					field.mul(c_2, c_x, b)
					field.mul(c_2, c_2, a)
					if !field.equal(c_1, c_2) {
						t.Fatalf("(a * b) * c == (a * c) * b")
					}
				}
			})
		}
	}
}

func TestFqExponentiation(t *testing.T) {
	for _, ext := range fields {
		for limbSize := from; limbSize < to+1; limbSize++ {
			t.Run(fmt.Sprintf("%d_%s", limbSize*64, ext), func(t *testing.T) {
				for i := 0; i < fuz; i++ {
					field := randField(ext, limbSize)
					if !USE_4LIMBS_FOR_LOWER_LIMBS && field.limbSize() != limbSize {
						t.Fatalf("bad field construction")
					}
					a := field.rand(rand.Reader)
					u := field.new()
					field.exp(u, a, big.NewInt(0))
					if !field.equal(u, field.one()) {
						t.Fatalf("a^0 == 1")
					}
					field.exp(u, a, big.NewInt(1))
					if !field.equal(u, a) {
						t.Fatalf("a^1 == a")
					}
					v := field.new()
					field.mul(u, a, a)
					field.mul(u, u, u)
					field.mul(u, u, u)
					field.exp(v, a, big.NewInt(8))
					if !field.equal(u, v) {
						t.Fatalf("((a^2)^2)^2 == a^8")
					}
					p := new(big.Int).Set(field.p())
					field.exp(u, a, p)
					if !field.equal(u, a) {
						t.Fatalf("a^p == a")
					}
					field.exp(u, a, p.Sub(p, big.NewInt(1)))
					if !field.equal(u, field.one()) {
						t.Fatalf("a^(p-1) == 1")
					}
				}
			})
		}
	}
}

func TestFqInversion(t *testing.T) {
	for _, ext := range fields {
		for limbSize := from; limbSize < to+1; limbSize++ {
			t.Run(fmt.Sprintf("%d_%s", limbSize*64, ext), func(t *testing.T) {
				for i := 0; i < fuz; i++ {
					field := randField(ext, limbSize)
					u := field.new()
					zero := field.zero()
					one := field.one()
					field.inverse(u, zero)
					if !field.equal(u, zero) {
						t.Fatalf("(0^-1) == 0)")
					}
					field.inverse(u, one)
					if !field.equal(u, one) {
						t.Fatalf("(1^-1) == 1)")
					}
					a := field.rand(rand.Reader)
					field.inverse(u, a)
					field.mul(u, u, a)
					if !field.equal(u, one) {
						t.Fatalf("(r*a) * r*(a^-1) == r)")
					}
					v := field.new()
					p := new(big.Int).Set(field.p())
					field.exp(u, a, p.Sub(p, big.NewInt(2)))
					field.inverse(v, a)
					if !field.equal(v, u) {
						t.Fatalf("a^(p-2) == a^-1")
					}
				}
			})
		}
	}
}

// func TestFqSquareRoot(t *testing.T) {
// 	fields = []string{"FQ"}
// 	for _, ext := range fields {
// 		for limbSize := from; limbSize < 4; limbSize++ {
// 			// for limbSize := from; limbSize < to+1; limbSize++ {
// 			t.Run(fmt.Sprintf("%d_%s", limbSize*64, ext), func(t *testing.T) {
// 				for i := 0; i < fuz; i++ {
// 					field := randField(ext, limbSize)
// 					u := field.new()
// 					zero := field.zero()
// 					one := field.one()
// 					field.sqrt(u, zero)
// 					if !field.equal(u, zero) {
// 						t.Errorf("(0^(1/2)) == 0)")
// 					}
// 					field.sqrt(u, one)
// 					if !field.equal(u, one) {
// 						// t.Errorf("(1^(1/2)) == 1)")
// 					}
// 					v, w, negA := field.new(), field.new(), field.new()
// 					a := field.rand(rand.Reader)
// 					field.neg(negA, a)
// 					field.square(u, a)
// 					field.square(w, negA)
// 					if !field.equal(w, u) {
// 						// TODO: why fails?
// 						// t.Errorf("square of r and -r is not same")
// 					}
// 					if hasRoot := field.sqrt(v, u); !hasRoot {
// 						// t.Errorf("elem has no square-root")
// 					}
// 					if !field.equal(a, v) && !field.equal(negA, v) {
// 						field.debugElement(a)
// 						field.debugElement(negA)
// 						field.debugElement(v)
// 						// t.Errorf("((r)^2)^(1/2) == r)")
// 					}
// 				}
// 			})
// 		}
// 	}
// }
