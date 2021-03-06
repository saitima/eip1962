package eip

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"unsafe"

	"golang.org/x/sys/cpu"
)

// fe is a pointer that addresses a field element given limb size
type fe = unsafe.Pointer

// for x86 we decide the instruction set in runtime
var nonADXBMI2 = !(cpu.X86.HasADX && cpu.X86.HasBMI2) || forceNonADXBMI2()

type fq struct {
	limbSize       int
	modulusBitLen  int
	modulusByteLen int
	p              fe
	inp            uint64
	one            fe
	_one           fe
	zero           fe
	r              fe
	r2             fe
	pbig           *big.Int
	rbig           *big.Int
	equal          func(a, b fe) bool
	cmp            func(a, b fe) int8
	copy           func(dst, stc fe)
	_mul           func(c, a, b, p fe, inp uint64)
	_add           func(c, a, b, p fe)
	_double        func(c, a, p fe)
	_sub           func(c, a, b, p fe)
	_neg           func(c, a, p fe)
	addn           func(a, b fe) uint64
	subn           func(a, b fe) uint64
	div_two        func(a fe)
	mul_two        func(a fe)
}

func newField(p []byte) (*fq, error) {
	var err error
	f := new(fq)
	pbig := new(big.Int).SetBytes(p)
	f.modulusByteLen = len(pbig.Bytes())
	f.modulusBitLen = pbig.BitLen()
	f.pbig = pbig
	f.p, f.limbSize, err = newFieldElementFromBytes(p)
	if err != nil {
		return nil, err
	}
	R := new(big.Int)
	R.SetBit(R, f.byteSize()*8, 1).Mod(R, f.pbig)
	R2 := new(big.Int)
	R2.Mul(R, R).Mod(R2, f.pbig)
	inpT := new(big.Int).ModInverse(new(big.Int).Neg(f.pbig), new(big.Int).SetBit(new(big.Int), 64, 1))
	f.r = newFieldElementFromBigUnchecked(f.limbSize, R)
	f.rbig = R
	f.one = newFieldElementFromBigUnchecked(f.limbSize, R)
	f.r2 = newFieldElementFromBigUnchecked(f.limbSize, R2)
	f._one = newFieldElementFromBigUnchecked(f.limbSize, big.NewInt(1))
	f.zero = newFieldElementFromBigUnchecked(f.limbSize, new(big.Int))
	if inpT == nil {
		return nil, fmt.Errorf("field is not applicable\n%s", hex.EncodeToString(p))
	}
	f.inp = inpT.Uint64()
	switch f.limbSize {
	case 1:
		f.equal = eq1
		f.copy = cpy1
		f.cmp = cmp1
		f.addn = addn1
		f.subn = subn1
		f._add = add1
		f._sub = sub1
		f._double = double1
		f._neg = _neg1
		f.div_two = div_two_1
		f.mul_two = mul_two_1
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_1
		} else {
			f._mul = mul1
		}
	case 2:
		f.equal = eq2
		f.copy = cpy2
		f.cmp = cmp2
		f.addn = addn2
		f.subn = subn2
		f._add = add2
		f._sub = sub2
		f._double = double2
		f._neg = _neg2
		f.div_two = div_two_2
		f.mul_two = mul_two_2
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_2
		} else {
			f._mul = mul2
		}
	case 3:
		f.equal = eq3
		f.copy = cpy3
		f.cmp = cmp3
		f.addn = addn3
		f.subn = subn3
		f._add = add3
		f._sub = sub3
		f._double = double3
		f._neg = _neg3
		f.div_two = div_two_3
		f.mul_two = mul_two_3
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_3
		} else {
			f._mul = mul3
		}
	case 4:
		f.equal = eq4
		f.copy = cpy4
		f.cmp = cmp4
		f.addn = addn4
		f.subn = subn4
		f._add = add4
		f._sub = sub4
		f._double = double4
		f._neg = _neg4
		f.div_two = div_two_4
		f.mul_two = mul_two_4
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_4
		} else {
			f._mul = mul4
		}
	case 5:
		f.equal = eq5
		f.copy = cpy5
		f.cmp = cmp5
		f.addn = addn5
		f.subn = subn5
		f._add = add5
		f._sub = sub5
		f._double = double5
		f._neg = _neg5
		f.div_two = div_two_5
		f.mul_two = mul_two_5
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_5
		} else {
			f._mul = mul5
		}
	case 6:
		f.equal = eq6
		f.copy = cpy6
		f.cmp = cmp6
		f.addn = addn6
		f.subn = subn6
		f._add = add6
		f._sub = sub6
		f._double = double6
		f._neg = _neg6
		f.div_two = div_two_6
		f.mul_two = mul_two_6
		f._mul = mul6
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_6
		} else {
			f._mul = mul6
		}
	case 7:
		f.equal = eq7
		f.copy = cpy7
		f.cmp = cmp7
		f.addn = addn7
		f.subn = subn7
		f._add = add7
		f._sub = sub7
		f._double = double7
		f._neg = _neg7
		f.div_two = div_two_7
		f.mul_two = mul_two_7
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_7
		} else {
			f._mul = mul7
		}
	case 8:
		f.equal = eq8
		f.copy = cpy8
		f.cmp = cmp8
		f.addn = addn8
		f.subn = subn8
		f._add = add8
		f._sub = sub8
		f._double = double8
		f._neg = _neg8
		f.div_two = div_two_8
		f.mul_two = mul_two_8
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_8
		} else {
			f._mul = mul8
		}
	case 9:
		f.equal = eq9
		f.copy = cpy9
		f.cmp = cmp9
		f.addn = addn9
		f.subn = subn9
		f._add = add9
		f._sub = sub9
		f._double = double9
		f._neg = _neg9
		f.div_two = div_two_9
		f.mul_two = mul_two_9
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_9
		} else {
			f._mul = mul9
		}
	case 10:
		f.equal = eq10
		f.copy = cpy10
		f.cmp = cmp10
		f.addn = addn10
		f.subn = subn10
		f._add = add10
		f._sub = sub10
		f._double = double10
		f._neg = _neg10
		f.div_two = div_two_10
		f.mul_two = mul_two_10
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_10
		} else {
			f._mul = mul10
		}
	case 11:
		f.equal = eq11
		f.copy = cpy11
		f.cmp = cmp11
		f.addn = addn11
		f.subn = subn11
		f._add = add11
		f._sub = sub11
		f._double = double11
		f._neg = _neg11
		f.div_two = div_two_11
		f.mul_two = mul_two_11
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_11
		} else {
			f._mul = mul11
		}
	case 12:
		f.equal = eq12
		f.copy = cpy12
		f.cmp = cmp12
		f.addn = addn12
		f.subn = subn12
		f._add = add12
		f._sub = sub12
		f._double = double12
		f._neg = _neg12
		f.div_two = div_two_12
		f.mul_two = mul_two_12
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_12
		} else {
			f._mul = mul12
		}
	case 13:
		f.equal = eq13
		f.copy = cpy13
		f.cmp = cmp13
		f.addn = addn13
		f.subn = subn13
		f._add = add13
		f._sub = sub13
		f._double = double13
		f._neg = _neg13
		f.div_two = div_two_13
		f.mul_two = mul_two_13
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_13
		} else {
			f._mul = mul13
		}
	case 14:
		f.equal = eq14
		f.copy = cpy14
		f.cmp = cmp14
		f.addn = addn14
		f.subn = subn14
		f._add = add14
		f._sub = sub14
		f._double = double14
		f._neg = _neg14
		f.div_two = div_two_14
		f.mul_two = mul_two_14
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_14
		} else {
			f._mul = mul14
		}
	case 15:
		f.equal = eq15
		f.copy = cpy15
		f.cmp = cmp15
		f.addn = addn15
		f.subn = subn15
		f._add = add15
		f._sub = sub15
		f._double = double15
		f._neg = _neg15
		f.div_two = div_two_15
		f.mul_two = mul_two_15
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_15
		} else {
			f._mul = mul15
		}
	case 16:
		f.equal = eq16
		f.copy = cpy16
		f.cmp = cmp16
		f.addn = addn16
		f.subn = subn16
		f._add = add16
		f._sub = sub16
		f._double = double16
		f._neg = _neg16
		f.div_two = div_two_16
		f.mul_two = mul_two_16
		if nonADXBMI2 {
			f._mul = mul_no_adx_bmi2_16
		} else {
			f._mul = mul16
		}
	default:
		return nil, fmt.Errorf("limb size %d is not implemented", f.limbSize)
	}
	return f, nil
}

func (f *fq) toMont(c, a fe) {
	f._mul(c, a, f.r2, f.p, f.inp)
}

func (f *fq) fromMont(c, a fe) {
	f._mul(c, a, f._one, f.p, f.inp)
}

func (f *fq) add(c, a, b fe) {
	f._add(c, a, b, f.p)
}

func (f *fq) double(c, a fe) {
	f._double(c, a, f.p)
}

func (f *fq) sub(c, a, b fe) {
	f._sub(c, a, b, f.p)
}

func (f *fq) neg(c, a fe) {
	if f.equal(a, f.zero) {
		f.copy(c, f.zero)
		return
	}
	f._neg(c, a, f.p)
}

func (f *fq) mul(c, a, b fe) {
	f._mul(c, a, b, f.p, f.inp)
}

func (f *fq) square(c, a fe) {
	f._mul(c, a, a, f.p, f.inp)
}

func (f *fq) exp(c, a fe, e *big.Int) {
	z := f.new()
	f.copy(z, f.r)
	for i := e.BitLen(); i >= 0; i-- {
		f.mul(z, z, z)
		if e.Bit(i) == 1 {
			f.mul(z, z, a)
		}
	}
	f.copy(c, z)
}

func (f *fq) sqrt(c, a fe) bool {
	negOne, tmp, b := f.new(), f.new(), f.new()
	f.neg(negOne, f.one)

	// power = (p-3)/4
	power := new(big.Int).Rsh(f.pbig, 2)
	// b = a^((p-3)/4)
	f.exp(b, a, power)
	// tmp = b^2 = a^(p-3/2) = a^-1 * a^(p-1/2)
	// if a is a square then a^(p-1/2) is 1 else -1
	// so tmp = a^-1 or -a^-1
	f.square(tmp, b)
	// check that a*a^-1 is -1
	f.mul(tmp, tmp, a)
	if f.equal(tmp, negOne) {
		f.copy(c, f.zero)
		return false
	}
	// c = a^(p-3/4) * a = a^(p+1/4)
	f.mul(c, b, a)
	return true
}

func (f *fq) sign(fe fe) int8 {
	neg := f.new()
	f.neg(neg, fe)
	return f.cmp(fe, neg)
}

func (f *fq) isOne(fe fe) bool {
	return f.equal(fe, f.one)
}

func (f *fq) isZero(fe fe) bool {
	return f.equal(fe, f.zero)
}

func (f *fq) isNonResidue(a fe, degree int) bool {
	zero := big.NewInt(0)
	result := f.new()
	exp := new(big.Int).Sub(f.pbig, big.NewInt(1))
	exp, rem := new(big.Int).DivMod(exp, big.NewInt(int64(degree)), zero)
	if rem.Cmp(zero) != 0 {
		return false
	}
	f.exp(result, a, exp)
	if f.equal(result, f.one) {
		return false
	}
	return true
}

func (f *fq) isValid(fe []byte) bool {
	feBig := new(big.Int).SetBytes(fe)
	if feBig.Cmp(f.pbig) != -1 {
		return false
	}
	return true
}

func (f *fq) new() fe {
	fe, err := newFieldElement(f.limbSize)
	if err != nil {
		// panic("this is unexpected")
	}
	return fe
}

func (f *fq) modulus() *big.Int {
	return new(big.Int).Set(f.pbig)
}

func (f *fq) rand(r io.Reader) fe {
	bi := new(big.Int)
	var err error
	for {
		bi, err = rand.Int(r, f.pbig)
		if err != nil {
			panic(err)
		}
		if bi.Cmp(new(big.Int)) != 0 {
			break
		}
	}
	return newFieldElementFromBigUnchecked(f.limbSize, bi)
}

func (f *fq) fromBytesNoTransform(in []byte) (fe, error) {
	if len(in) != f.byteSize() {
		return nil, fmt.Errorf("bad input size")
	}
	fe, limbSize, err := newFieldElementFromBytes(in)
	if err != nil {
		return nil, err
	}
	if limbSize != f.limbSize {
		// panic("this is unexpected")
	}
	return fe, nil
}

func (f *fq) fromBytes(in []byte) (fe, error) {
	if len(in) != f.byteSize() {
		return nil, fmt.Errorf("bad input size %d %d", len(in), f.byteSize())
	}
	if !f.isValid(in) {
		return nil, fmt.Errorf("input is a larger number than modulus")
	}
	fe, limbSize, err := newFieldElementFromBytes(in)
	if err != nil {
		return nil, err
	}
	if limbSize != f.limbSize {
		// panic("this is unexpected")
	}
	f.toMont(fe, fe)
	return fe, nil
}

func (f *fq) fromString(hexStr string) (fe, error) {
	str := hexStr
	if len(str) > 1 && str[:2] == "0x" {
		str = hexStr[2:]
	}
	in, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	if !f.isValid(in) {
		return nil, fmt.Errorf("input is a larger number than modulus")
	}
	if len(in) > f.byteSize() {
		return nil, fmt.Errorf("bad input size")
	}
	fe, limbSize, err := newFieldElementFromBytes(padBytes(in, f.byteSize()))
	if err != nil {
		return nil, err
	}
	if limbSize != f.limbSize {
		// panic("this is unexpected")
	}
	f.toMont(fe, fe)
	return fe, nil
}

func (f *fq) fromBig(a *big.Int) (fe, error) {
	in := a.Bytes()
	if !f.isValid(in) {
		return nil, fmt.Errorf("input is a larger number than modulus")
	}
	if len(in) > f.byteSize() {
		return nil, fmt.Errorf("bad input size %d", len(in))
	}
	fe, limbSize, err := newFieldElementFromBytes(padBytes(in, f.byteSize()))
	if err != nil {
		return nil, err
	}
	if limbSize != f.limbSize {
		// panic("this is unexpected")
	}
	f.toMont(fe, fe)
	return fe, nil
}

func (f *fq) toBytes(in fe) []byte {
	t := f.new()
	f.fromMont(t, in)
	return f.toBytesNoTransform(t)
}

func (f *fq) toBytesDense(in fe) []byte {
	t := f.new()
	f.fromMont(t, in)
	denseLength := f.modulusByteLen
	out := make([]byte, denseLength)
	sparse := f.toBytesNoTransform(t)
	copy(out[:], sparse[f.byteSize()-denseLength:])
	return out
}

func (f *fq) toBytesNoTransform(in fe) []byte {
	switch f.limbSize {
	case 1:
		return toBytes((*[1]uint64)(in)[:])
	case 2:
		return toBytes((*[2]uint64)(in)[:])
	case 3:
		return toBytes((*[3]uint64)(in)[:])
	case 4:
		return toBytes((*[4]uint64)(in)[:])
	case 5:
		return toBytes((*[5]uint64)(in)[:])
	case 6:
		return toBytes((*[6]uint64)(in)[:])
	case 7:
		return toBytes((*[7]uint64)(in)[:])
	case 8:
		return toBytes((*[8]uint64)(in)[:])
	case 9:
		return toBytes((*[9]uint64)(in)[:])
	case 10:
		return toBytes((*[10]uint64)(in)[:])
	case 11:
		return toBytes((*[11]uint64)(in)[:])
	case 12:
		return toBytes((*[12]uint64)(in)[:])
	case 13:
		return toBytes((*[13]uint64)(in)[:])
	case 14:
		return toBytes((*[14]uint64)(in)[:])
	case 15:
		return toBytes((*[15]uint64)(in)[:])
	case 16:
		return toBytes((*[16]uint64)(in)[:])
	default:
		panic("not implemented")
	}
}

func (f *fq) toBig(in fe) *big.Int {
	return new(big.Int).SetBytes(f.toBytes(in))
}

func (f *fq) toBigNoTransform(in fe) *big.Int {
	return new(big.Int).SetBytes(f.toBytesNoTransform(in))
}

func (f *fq) toString(in fe) string {
	return hex.EncodeToString(f.toBytes(in))
}

func (f *fq) toStringNoTransform(in fe) string {
	return hex.EncodeToString(f.toBytesNoTransform(in))
}

func (f *fq) byteSize() int {
	return f.limbSize * 8
}

func toBytes(fe []uint64) []byte {
	size := len(fe)
	byteSize := size * 8
	out := make([]byte, byteSize)
	var a int
	for i := 0; i < size; i++ {
		a = byteSize - i*8
		out[a-1] = byte(fe[i])
		out[a-2] = byte(fe[i] >> 8)
		out[a-3] = byte(fe[i] >> 16)
		out[a-4] = byte(fe[i] >> 24)
		out[a-5] = byte(fe[i] >> 32)
		out[a-6] = byte(fe[i] >> 40)
		out[a-7] = byte(fe[i] >> 48)
		out[a-8] = byte(fe[i] >> 56)
	}
	return out
}

// newFieldElement returns pointer of an uint64 array.
// limbSize is calculated according to size of input slice
func newFieldElementFromBytes(in []byte) (fe, int, error) {
	byteSize := len(in)
	limbSize := byteSize / 8
	if byteSize%8 != 0 {
		return nil, 0, fmt.Errorf("bad input byte size %d", byteSize)
	}
	// TODO: remove after fuzz testing
	if limbSize < 4 {
		limbSize = 4
		in = padBytes(in, 32)
	}
	a, err := newFieldElement(limbSize)
	if err != nil {
		return nil, 0, err
	}
	var data []uint64
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = uintptr(a)
	sh.Len, sh.Cap = limbSize, limbSize
	if err := limbSliceFromBytes(data[:], in); err != nil {
		// panic("this is unexpected")
	}
	return a, limbSize, nil
}

func newFieldElement(limbSize int) (fe, error) {
	switch limbSize {
	case 1:
		return unsafe.Pointer(&[1]uint64{}), nil
	case 2:
		return unsafe.Pointer(&[2]uint64{}), nil
	case 3:
		return unsafe.Pointer(&[3]uint64{}), nil
	case 4:
		return unsafe.Pointer(&[4]uint64{}), nil
	case 5:
		return unsafe.Pointer(&[5]uint64{}), nil
	case 6:
		return unsafe.Pointer(&[6]uint64{}), nil
	case 7:
		return unsafe.Pointer(&[7]uint64{}), nil
	case 8:
		return unsafe.Pointer(&[8]uint64{}), nil
	case 9:
		return unsafe.Pointer(&[9]uint64{}), nil
	case 10:
		return unsafe.Pointer(&[10]uint64{}), nil
	case 11:
		return unsafe.Pointer(&[11]uint64{}), nil
	case 12:
		return unsafe.Pointer(&[12]uint64{}), nil
	case 13:
		return unsafe.Pointer(&[13]uint64{}), nil
	case 14:
		return unsafe.Pointer(&[14]uint64{}), nil
	case 15:
		return unsafe.Pointer(&[15]uint64{}), nil
	case 16:
		return unsafe.Pointer(&[16]uint64{}), nil
	default:
		return nil, fmt.Errorf("limb size %d is not implemented", limbSize)
	}
}

func newFieldElementFromBigUnchecked(limbSize int, bi *big.Int) fe {
	in := bi.Bytes()
	byteSize := limbSize * 8
	fe, _, _ := newFieldElementFromBytes(padBytes(in, byteSize))
	return fe
}

func limbSliceFromBytes(out []uint64, in []byte) error {
	var byteSize = len(in)
	var limbSize = len(out)
	if limbSize*8 != byteSize {
		return fmt.Errorf("(byteSize != limbSize * 8), %d, %d", byteSize, limbSize)
	}
	var a int
	for i := 0; i < limbSize; i++ {
		a = byteSize - i*8
		out[i] = uint64(in[a-1]) | uint64(in[a-2])<<8 |
			uint64(in[a-3])<<16 | uint64(in[a-4])<<24 |
			uint64(in[a-5])<<32 | uint64(in[a-6])<<40 |
			uint64(in[a-7])<<48 | uint64(in[a-8])<<56
	}
	return nil
}

func padBytes(in []byte, size int) []byte {
	out := make([]byte, size)
	if len(in) > size {
		panic("bad input for padding")
	}
	copy(out[size-len(in):], in)
	return out
}

func (f *fq) inverse(inv, e fe) bool {
	if f.equal(e, f.zero) {
		f.copy(inv, f.zero)
		return false
	}
	u, v, s, r := f.new(), f.new(), f.new(), f.new()
	zero := f.new()
	f.copy(u, f.p)
	f.copy(v, e)
	f.copy(r, f.zero)
	f.copy(s, f._one)
	var k int
	var found = false
	byteSize := f.byteSize()
	bitSize := byteSize * 8
	// Phase 1
	for i := 0; i < bitSize*2; i++ {
		if f.equal(v, zero) {
			found = true
			break
		}
		if is_even(u) {
			f.div_two(u)
			f.mul_two(s)
		} else if is_even(v) {
			f.div_two(v)
			f.mul_two(r)
		} else if f.cmp(u, v) == 1 {
			f.subn(u, v)
			f.div_two(u)
			f.addn(r, s)
			f.mul_two(s)
		} else if f.cmp(v, u) != -1 {
			f.subn(v, u)
			f.div_two(v)
			f.addn(s, r)
			f.mul_two(r)
		}
		k += 1
	}
	if !found {
		f.copy(inv, zero)
		return false
	}

	if f.cmp(r, f.p) != -1 {
		f.subn(r, f.p)
	}
	f.copy(u, f.p)
	f.subn(u, r)
	// phase 2
	montPower := f.limbSize * 64
	modulusBitsCeil := f.modulusBitLen
	kInRange := modulusBitsCeil <= k && k <= montPower+modulusBitsCeil
	if !kInRange {
		f.copy(inv, zero)
		return false
	}

	if modulusBitsCeil <= k && k <= montPower {
		f.mul(u, u, f.r2)
		k += montPower
	}

	if k > 2*montPower {
		f.copy(inv, zero)
		return false
	}
	if 2*montPower-k > montPower {
		f.copy(inv, zero)
		return false
	}
	// now we need montgomery(!) multiplication by 2^(2m - k)
	// since 2^(2m - k) < 2^m then we are ok with a multiplication without preliminary reduction
	// of the representation as montgomery multiplication will handle it for us
	xBig := big.NewInt(1)
	xBig = new(big.Int).Lsh(xBig, uint(uint32(2*montPower-k)))
	xBytes := xBig.Bytes()
	if len(xBytes) < f.limbSize*8 {
		xBytes = padBytes(xBig.Bytes(), f.limbSize*8)
	}
	x, err := f.fromBytesNoTransform(xBytes)
	if err != nil {
		f.copy(inv, zero)
		return false
	}
	f.mul(u, u, x)
	f.toMont(inv, u)
	return true
}
