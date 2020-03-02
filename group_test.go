package eip

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
)

// group interface abstacts g1, g22, g23 groups from testing suite
type group interface {
	new() point
	debugPoint(p point)
	fromBytes(in []byte) (point, error)
	fromBytesDense(in []byte) (point, error)
	toBytes(p1 point) []byte
	toBytesDense(p1 point) []byte
	equal(p1, p2 point) bool
	affine(p1, p2 point)
	zero() point
	Q() *big.Int
	fieldModulus() *big.Int
	fieldModulusByteLen() int
	isOnCurve(c point) bool
	checkCorrectSubgroup(c point) bool
	mulScalar(c, a point, e *big.Int) point
	wnafMul(c, a point, e *big.Int) point
	add(c, a, b point) point
	sub(c, a, b point) point
	neg(c, a point) point
	double(c, a point) point
	multiExp(c point, p []point, s []*big.Int) (point, error)
}

// group interface abstacts points at g1, g22, g23 groups from testing suite
type point interface{}

// g1Test wraps g1 to match group interface
type g1Test struct {
	*g1
}

func (g g1Test) Q() *big.Int {
	return g.g1.q
}

func (g g1Test) fieldModulus() *big.Int {
	return g.g1.f.modulus()
}

func (g g1Test) fieldModulusByteLen() int {
	return g.g1.f.modulusByteLen
}

func (g g1Test) fromBytes(in []byte) (point, error) {
	return g.g1.fromBytes(in)
}

func (g g1Test) toBytes(p1 point) []byte {
	return g.g1.toBytes(p1.(*pointG1))
}

func (g g1Test) toBytesDense(p1 point) []byte {
	return g.g1.toBytesDense(p1.(*pointG1))
}

func (g g1Test) fromBytesDense(in []byte) (point, error) {
	byteLen := g.g1.f.byteSize()
	buf := make([]byte, byteLen*2)
	copy(buf[:byteLen], padBytes(in[:len(in)/2], byteLen))
	copy(buf[byteLen:], padBytes(in[len(in)/2:], byteLen))
	return g.g1.fromBytes(buf)
}

func (g g1Test) isOnCurve(p point) bool {
	return g.g1.isOnCurve(p.(*pointG1))
}

func (g g1Test) checkCorrectSubgroup(p point) bool {
	return g.g1.checkCorrectSubgroup(p.(*pointG1))
}

func (g g1Test) affine(r, p point) {
	g.g1.affine(r.(*pointG1), p.(*pointG1))
}

func (g g1Test) new() point {
	return g.g1.newPoint()
}

func (g g1Test) debugPoint(p point) {
	fmt.Println(g.g1.toString(p.(*pointG1)))
}

func (g g1Test) equal(p1, p2 point) bool {
	return g.g1.equal(p1.(*pointG1), p2.(*pointG1))
}

func (g g1Test) zero() point {
	return g.g1.zero()
}

func (g g1Test) add(c, a, b point) point {
	return g.g1.add(c.(*pointG1), a.(*pointG1), b.(*pointG1))
}

func (g g1Test) sub(c, a, b point) point {
	return g.g1.sub(c.(*pointG1), a.(*pointG1), b.(*pointG1))
}

func (g g1Test) neg(c, a point) point {
	return g.g1.neg(c.(*pointG1), a.(*pointG1))
}

func (g g1Test) double(c, a point) point {
	return g.g1.double(c.(*pointG1), a.(*pointG1))
}

func (g g1Test) mulScalar(c point, p point, s *big.Int) point {
	return g.g1.mulScalar(c.(*pointG1), p.(*pointG1), s)
}

func (g g1Test) wnafMul(c point, p point, s *big.Int) point {
	return g.g1.wnafMul(c.(*pointG1), p.(*pointG1), s)
}

func (g g1Test) multiExp(c point, p []point, s []*big.Int) (point, error) {
	p_ := make([]*pointG1, len(p))
	for i := 0; i < len(p); i++ {
		p_[i] = p[i].(*pointG1)
	}
	return g.g1.multiExp(c.(*pointG1), p_, s)
}

// g22Test wraps g22 to match group interface
type g22Test struct {
	*g22
}

func (g g22Test) Q() *big.Int {
	return g.g22.q
}

func (g g22Test) fieldModulus() *big.Int {
	return g.g22.f.modulus()
}

func (g g22Test) fieldModulusByteLen() int {
	return g.g22.f.f.modulusByteLen
}

func (g g22Test) fromBytes(in []byte) (point, error) {
	return g.g22.fromBytes(in)
}

func (g g22Test) toBytes(p1 point) []byte {
	return g.g22.toBytes(p1.(*pointG22))
}

func (g g22Test) toBytesDense(p1 point) []byte {
	return g.g22.toBytesDense(p1.(*pointG22))
}

func (g g22Test) fromBytesDense(in []byte) (point, error) {
	q := len(in) / 4
	p := g.g22.f.byteSize() / 2
	buf := make([]byte, p*4)
	a0 := padBytes(in[:q], p)
	a1 := padBytes(in[q:q*2], p)
	b0 := padBytes(in[q*2:q*3], p)
	b1 := padBytes(in[q*3:], p)
	copy(buf[:p], a0)
	copy(buf[p:p*2], a1)
	copy(buf[p*2:p*3], b0)
	copy(buf[p*3:], b1)
	return g.g22.fromBytes(buf)
}

func (g g22Test) isOnCurve(p point) bool {
	return g.g22.isOnCurve(p.(*pointG22))
}

func (g g22Test) checkCorrectSubgroup(p point) bool {
	return g.g22.checkCorrectSubgroup(p.(*pointG22))
}

func (g g22Test) affine(r, p point) {
	g.g22.affine(r.(*pointG22), p.(*pointG22))
}

func (g g22Test) new() point {
	return g.g22.newPoint()
}

func (g g22Test) debugPoint(p point) {
	fmt.Println(g.g22.toString(p.(*pointG22)))
}

func (g g22Test) equal(p1, p2 point) bool {
	return g.g22.equal(p1.(*pointG22), p2.(*pointG22))
}

func (g g22Test) zero() point {
	return g.g22.zero()
}

func (g g22Test) add(c, a, b point) point {
	return g.g22.add(c.(*pointG22), a.(*pointG22), b.(*pointG22))
}

func (g g22Test) sub(c, a, b point) point {
	return g.g22.sub(c.(*pointG22), a.(*pointG22), b.(*pointG22))
}

func (g g22Test) neg(c, a point) point {
	return g.g22.neg(c.(*pointG22), a.(*pointG22))
}

func (g g22Test) double(c, a point) point {
	return g.g22.double(c.(*pointG22), a.(*pointG22))
}

func (g g22Test) mulScalar(c point, p point, s *big.Int) point {
	return g.g22.mulScalar(c.(*pointG22), p.(*pointG22), s)
}

func (g g22Test) wnafMul(c point, p point, s *big.Int) point {
	return g.g22.wnafMul(c.(*pointG22), p.(*pointG22), s)
}

func (g g22Test) multiExp(c point, p []point, s []*big.Int) (point, error) {
	p_ := make([]*pointG22, len(p))
	for i := 0; i < len(p); i++ {
		p_[i] = p[i].(*pointG22)
	}
	return g.g22.multiExp(c.(*pointG22), p_, s)
}

// g23Test wraps g23 to match group interface
type g23Test struct {
	*g23
}

func (g g23Test) Q() *big.Int {
	return g.g23.q
}

func (g g23Test) fieldModulus() *big.Int {
	return g.g23.f.modulus()
}

func (g g23Test) fieldModulusByteLen() int {
	return g.g23.f.f.modulusByteLen
}

func (g g23Test) fromBytes(in []byte) (point, error) {
	return g.g23.fromBytes(in)
}

func (g g23Test) toBytes(p1 point) []byte {
	return g.g23.toBytes(p1.(*pointG23))
}

func (g g23Test) toBytesDense(p1 point) []byte {
	return g.g23.toBytesDense(p1.(*pointG23))
}

func (g g23Test) fromBytesDense(in []byte) (point, error) {
	byteLen := g.g23.f.byteSize()
	buf := make([]byte, byteLen*3)
	q := len(in) / 6
	a0 := padBytes(in[:q], byteLen)
	a1 := padBytes(in[q:q*2], byteLen)
	a2 := padBytes(in[q*2:q*3], byteLen)
	b0 := padBytes(in[q*3:q*4], byteLen)
	b1 := padBytes(in[q*4:q*5], byteLen)
	b2 := padBytes(in[q*5:], byteLen)
	copy(buf[byteLen:], a0)
	copy(buf[byteLen:byteLen*2], a1)
	copy(buf[byteLen*2:byteLen*3], a2)
	copy(buf[byteLen*3:byteLen*3], b0)
	copy(buf[byteLen*4:byteLen*5], b1)
	copy(buf[byteLen*5:], b2)
	return g.g23.fromBytes(buf)
}

func (g g23Test) isOnCurve(p point) bool {
	return g.g23.isOnCurve(p.(*pointG23))
}

func (g g23Test) checkCorrectSubgroup(p point) bool {
	return g.g23.checkCorrectSubgroup(p.(*pointG23))
}

func (g g23Test) affine(r, p point) {
	g.g23.affine(r.(*pointG23), p.(*pointG23))
}

func (g g23Test) new() point {
	return g.g23.newPoint()
}

func (g g23Test) debugPoint(p point) {
	fmt.Println(g.g23.toString(p.(*pointG23)))
}

func (g g23Test) equal(p1, p2 point) bool {
	return g.g23.equal(p1.(*pointG23), p2.(*pointG23))
}

func (g g23Test) zero() point {
	return g.g23.zero()
}

func (g g23Test) add(c, a, b point) point {
	return g.g23.add(c.(*pointG23), a.(*pointG23), b.(*pointG23))
}

func (g g23Test) sub(c, a, b point) point {
	return g.g23.sub(c.(*pointG23), a.(*pointG23), b.(*pointG23))
}

func (g g23Test) neg(c, a point) point {
	return g.g23.neg(c.(*pointG23), a.(*pointG23))
}

func (g g23Test) double(c, a point) point {
	return g.g23.double(c.(*pointG23), a.(*pointG23))
}

func (g g23Test) mulScalar(c point, p point, s *big.Int) point {
	return g.g23.mulScalar(c.(*pointG23), p.(*pointG23), s)
}

func (g g23Test) wnafMul(c point, p point, s *big.Int) point {
	return g.g23.wnafMul(c.(*pointG23), p.(*pointG23), s)
}

func (g g23Test) multiExp(c point, p []point, s []*big.Int) (point, error) {
	p_ := make([]*pointG23, len(p))
	for i := 0; i < len(p); i++ {
		p_[i] = p[i].(*pointG23)
	}
	return g.g23.multiExp(c.(*pointG23), p_, s)
}

func ceilBitLen(b []byte) int {
	// return (((len(b) - 1) / 8) + 1) * 8
	return (((len(b)) / 64) + 1) * 64
}

// vectorJSON is used as input vector for
// various type of tests. vectorJSON is parsed
// in 'builder' to byte streams.
type vectorJSON struct {
	GroupOrder string `json:"r"`
	FieldOrder string `json:"q"`

	// a and b coefficients for g1
	A string `json:"A"`
	B string `json:"B"`

	// generator of g1
	G1x string `json:"g1_x"`
	G1y string `json:"g1_y"`

	// generator of g2
	G2x0 string `json:"g2_x_0"`
	G2x1 string `json:"g2_x_1"`
	G2x2 string `json:"g2_x_2"`
	G2y0 string `json:"g2_y_0"`
	G2y1 string `json:"g2_y_1"`
	G2y2 string `json:"g2_y_2"`

	// a and b coefficients for g2
	// however we already can calculate these
	// coefficients using a,b coeffs for g1
	// A20          string `json:"A_twist_0"`
	// A21          string `json:"A_twist_1"`
	// B20          string `json:"B_twist_0"`
	// B21          string `json:"B_twist_1"`
	NonResidue   string `json:"non_residue"`
	NonResidue20 string `json:"quadratic_non_residue_0"`
	NonResidue21 string `json:"quadratic_non_residue_1"`
	NonResidue22 string `json:"quadratic_non_residue_2"`

	// loop parameter
	Z string `json:"x"`

	// special fields  for BLS pairing
	IsDType string `json:"is_D_type"`

	// special fields for MNT4 and MNT6 pairing
	ExpW0 string `json:"expW0"`
	ExpW1 string `json:"expW1"`
}

// given json vector we resolve its
// extension degree by checking parameters
func (v *vectorJSON) resolveExtentionDegree() int {
	if v.NonResidue22 != "" {
		return 3
	}
	return 2
}

// in order to test specific component easier
// builder constructs fields, groups and as well as pairing engines
// with given input values namely vectorJSON
type builder struct {
	t             *testing.T
	limbSize      int
	testInput     input
	cache         builderCache
	willDoPairing bool
	tag           string
	family        string
}

type input struct {
	fieldOrder []byte
	groupOrder []byte
	// g1 coefficients
	a []byte
	b []byte
	// generator for g1
	g1one []byte
	// generator for g22
	g2one []byte
	// generator for g23
	g3one      []byte
	nonResidue []byte
	// non residue for fq2 or fq3
	nonResidue2 []byte
	// loop parameter
	negz            bool
	z               []byte
	extensionDegree int
	// BLS pairing inputs
	twistType int
	// MNT4, MNT6 pairign inputs
	expW0neg bool
	expW0    []byte
	expW1    []byte
}

// builderCache stores calculated components in memory
// in order not to calculate these components multiple times
type builderCache struct {
	fq   *fq
	fq2  *fq2
	fq3  *fq3
	fq4  *fq4
	fq6C *fq6C
	fq6Q *fq6Q
	fq12 *fq12
	a    fe
	b    fe
	a2   *fe2
	b2   *fe2
	a3   *fe3
	b3   *fe3
	g1   *g1
	G1   *pointG1
	g22  *g22
	g23  *g23
	G22  *pointG22
	G23  *pointG23
}

type builderOpts struct {
	willDoPairing bool
	family        string
}

func newBuilderOptPairing(family string) *builderOpts {
	return &builderOpts{family: family, willDoPairing: true}
}

func newBuilderOpt(family string) *builderOpts {
	return &builderOpts{family: family}
}

// testBuilderFromVector parses json vector and outputs a test builder
func testBuilderFromVector(t *testing.T, tag string, vectorJSON *vectorJSON, opts *builderOpts) *builder {
	// to parse negative hex values
	maybeNegativeHex := func(hex string) (bool, []byte) {
		if hex[:1] == "-" {
			return true, fromHex(-1, hex[1:])
		}
		return false, fromHex(-1, hex)
	}
	degree := vectorJSON.resolveExtentionDegree()
	var input input
	var builder builder
	// testing context
	builder.tag = tag
	builder.t = t
	// byte size, limb size
	p := new(big.Int).SetBytes(fromHex(-1, vectorJSON.FieldOrder))
	byteLen := ((p.BitLen() / 64) + 1) * 8
	// size := bitLen / 8
	builder.limbSize = byteLen / 8
	// orders
	input.fieldOrder = fromHex(byteLen, vectorJSON.FieldOrder)
	input.groupOrder = fromHex(byteLen, vectorJSON.GroupOrder)
	// a,b coefficients
	input.a = fromHex(byteLen, vectorJSON.A)
	input.b = fromHex(byteLen, vectorJSON.B)
	// g1 generator
	input.g1one = fromHex(byteLen, vectorJSON.G1x, vectorJSON.G1y)
	// non residue 1
	nonResidueIsNeg, nonResidue := maybeNegativeHex(vectorJSON.NonResidue)
	if nonResidueIsNeg {
		q := new(big.Int).SetBytes(input.fieldOrder)
		nonResidueBig := new(big.Int).SetBytes(nonResidue)
		nonResidueBig.Sub(q, nonResidueBig)
		input.nonResidue = padBytes(nonResidueBig.Bytes(), byteLen)
	} else {
		input.nonResidue = padBytes(nonResidue, byteLen)
	}
	if degree == 2 {
		// g22 generator
		// non residue 2
		input.g2one = fromHex(byteLen, vectorJSON.G2x0, vectorJSON.G2x1, vectorJSON.G2y0, vectorJSON.G2y1)
		input.nonResidue2 = fromHex(byteLen, vectorJSON.NonResidue20, vectorJSON.NonResidue21)
	} else if degree == 3 {
		// g23 generator
		// non residue 2
		input.g3one = fromHex(byteLen, vectorJSON.G2x0, vectorJSON.G2x1, vectorJSON.G2x2, vectorJSON.G2y0, vectorJSON.G2y1, vectorJSON.G2y2)
		input.nonResidue2 = fromHex(byteLen, vectorJSON.NonResidue20, vectorJSON.NonResidue21, vectorJSON.NonResidue22)
	} else {
		t.Fatal("unrecognized degree", degree)
	}
	// z for pairing loop
	input.negz, input.z = maybeNegativeHex(vectorJSON.Z)
	if opts != nil {
		builder.willDoPairing = opts.willDoPairing
		builder.family = opts.family
		switch opts.family {
		case "BLS":
			if degree != 2 {
				t.Fatal("bad degree for bls family")
			}
			input.extensionDegree = degree
			if vectorJSON.IsDType == "True" {
				input.twistType = TWIST_D
			} else {
				input.twistType = TWIST_M
			}
		case "BN":
			if degree != 2 {
				t.Fatal("bad degree for bls family")
			}
			input.extensionDegree = degree
			if vectorJSON.IsDType == "True" {
				input.twistType = TWIST_D
			} else {
				input.twistType = TWIST_M
			}
		case "MNT4":
			if degree != 2 {
				t.Fatal("bad degree for mnt4 family")
			}
			input.extensionDegree = degree
			input.expW0neg, input.expW0 = maybeNegativeHex(vectorJSON.ExpW0)
			input.expW1 = fromHex(-1, vectorJSON.ExpW1)
		case "MNT6":
			if degree != 3 {
				t.Fatal("bad degree for mnt6 family")
			}
			input.extensionDegree = degree
			input.expW0neg, input.expW0 = maybeNegativeHex(vectorJSON.ExpW0)
			input.expW1 = fromHex(-1, vectorJSON.ExpW1)
		default:
			t.Fatal("unrecognized family", opts.family)
		}
	}
	builder.testInput = input
	return &builder
}

func testBuilderFromFile(t *testing.T, file string, opts *builderOpts) *builder {
	vectorJSON, _ := readVectorFile(t, file)
	return testBuilderFromVector(t, file, vectorJSON, opts)
}

func readVectorFile(t *testing.T, file string) (*vectorJSON, error) {
	data, err := ioutil.ReadFile("test_vectors/" + file)
	if err != nil {
		t.Fatal(err)
	}
	var vectorJSON vectorJSON
	if err := json.Unmarshal(data, &vectorJSON); err != nil {
		t.Fatal(err)
	}
	return &vectorJSON, nil
}

func (b *builder) input() input {
	return b.testInput
}

func (b *builder) fq() *fq {
	input := b.input()
	cache := b.cache
	if cache.fq != nil {
		return cache.fq
	}
	f, err := newField(input.fieldOrder)
	if err != nil {
		b.t.Fatal(err)
	}
	if f.limbSize != b.limbSize {
		b.t.Fatal(fmt.Errorf("unexpected limb size"))
	}
	b.cache.fq = f
	return f
}

func (b *builder) fq2() *fq2 {
	inputs := b.input()
	cache := b.cache
	if cache.fq2 != nil {
		return cache.fq2
	}
	fq := b.fq()
	fq2, err := newFq2(fq, inputs.nonResidue)
	if err != nil {
		b.t.Fatal(err)
	}
	if b.willDoPairing {
		fq2.calculateFrobeniusCoeffs()
	}
	b.cache.fq2 = fq2
	return fq2
}

func (b *builder) fq3() *fq3 {
	inputs := b.input()
	cache := b.cache
	if cache.fq3 != nil {
		return cache.fq3
	}
	fq := b.fq()
	fq3, err := newFq3(fq, inputs.nonResidue)
	if err != nil {
		b.t.Fatal(err)
	}
	if b.willDoPairing {
		fq3.calculateFrobeniusCoeffs()
	}
	b.cache.fq3 = fq3
	return fq3
}

func (b *builder) fq4() *fq4 {
	input := b.input()
	cache := b.cache
	if cache.fq4 != nil {
		return cache.fq4
	}
	fq2 := b.fq2()
	fq4, err := newFq4(fq2, input.nonResidue2)
	if err != nil {
		b.t.Fatal(err)
	}
	if b.willDoPairing {
		fq4.calculateFrobeniusCoeffs()
	}
	b.cache.fq4 = fq4
	return fq4
}

func (b *builder) fq4TestInstance() field {
	return fq4Test{b.fq4()}
}

func (b *builder) fq6C() *fq6C {
	input := b.input()
	cache := b.cache
	if cache.fq6C != nil {
		return cache.fq6C
	}
	fq2 := b.fq2()
	fq6, err := newFq6Cubic(fq2, input.nonResidue2)
	if err != nil {
		b.t.Fatal(err)
	}
	if b.willDoPairing {
		fq6.calculateFrobeniusCoeffs()
	}
	cache.fq6C = fq6
	return fq6
}

func (b *builder) fq6Q() *fq6Q {
	input := b.input()
	cache := b.cache
	if cache.fq6Q != nil {
		return cache.fq6Q
	}
	fq3 := b.fq3()
	fq6, err := newFq6Quadratic(fq3, input.nonResidue2)
	if err != nil {
		b.t.Fatal(err)
	}
	if b.willDoPairing {
		fq6.calculateFrobeniusCoeffs()
	}
	cache.fq6Q = fq6
	return fq6
}

func (b *builder) fq6QTestInstance() field {
	return fq6QTest{b.fq6Q()}
}

func (b *builder) fq12() *fq12 {
	cache := b.cache
	if cache.fq12 != nil {
		return cache.fq12
	}
	fq6 := b.fq6C()
	fq12, _ := newFq12(fq6, nil)
	if b.willDoPairing {
		fq12.calculateFrobeniusCoeffs()
	}
	b.cache.fq12 = fq12
	return fq12
}

func (b *builder) fq12TestInstance() field {
	return fq12Test{b.fq12()}
}

func (b *builder) g1ab() (fe, fe) {
	inputs := b.input()
	cache := b.cache
	if cache.a != nil && cache.b != nil {
		return cache.a, cache.b
	}
	fq := cache.fq
	A, err := fq.fromBytes(inputs.a)
	if err != nil {
		b.t.Fatal(err)
	}
	cache.a = A
	B, err := fq.fromBytes(inputs.b)
	if err != nil {
		b.t.Fatal(err)
	}
	cache.b = B
	return A, B
}

func (b *builder) g1() *g1 {
	inputs := b.input()
	cache := b.cache
	if cache.g1 != nil {
		return cache.g1
	}
	fq := b.fq()
	A, B := b.g1ab()
	q := new(big.Int).SetBytes(inputs.groupOrder)
	g1, err := newG1(fq, A, B, q)
	if err != nil {
		b.t.Fatal(err)
	}
	G1, err := g1.fromBytes(inputs.g1one)
	if err != nil {
		b.t.Fatal(err)
	}
	if !g1.isOnCurve(G1) {
		b.t.Fatalf("g1 is not on curve")
	}
	if !g1.checkCorrectSubgroup(G1) {
		b.t.Fatalf("g1 is not on correct subgroup")
	}
	b.cache.G1 = G1
	b.cache.g1 = g1
	return g1
}

func (b *builder) G1() point {
	return b.cache.G1
}

func (b *builder) g1TestInstance() group {
	return g1Test{b.g1()}
}

func (b *builder) g22() *g22 {
	cache := b.cache
	if cache.g22 != nil {
		return cache.g22
	}
	inputs := b.input()
	fq2 := b.fq2()
	A, B := b.g22ab()
	q := new(big.Int).SetBytes(inputs.groupOrder)
	g22, err := newG22(fq2, A, B, q)
	if err != nil {
		b.t.Fatal(err)
	}
	G2, err := g22.fromBytes(inputs.g2one)
	if err != nil {
		b.t.Fatal(err)
	}
	if !g22.isOnCurve(G2) {
		b.t.Fatalf("g2 is not on curve")
	}
	if !g22.checkCorrectSubgroup(G2) {
		b.t.Fatalf("g2 is not on correct subgroup")
	}
	b.cache.G22 = G2
	cache.g22 = g22
	return g22
}

func (b *builder) g22ab() (*fe2, *fe2) {
	cache := b.cache
	if cache.a2 != nil && cache.b2 != nil {
		return cache.a2, cache.b2
	}
	inputs := b.input()
	fq := b.fq()
	fq2 := b.fq2()
	A, B := b.g1ab()
	B2, A2 := fq2.new(), fq2.new()
	switch b.family {
	case "BLS":
		twist, err := fq2.fromBytes(inputs.nonResidue2)
		if err != nil {
			b.t.Fatal(err)
		}
		if inputs.twistType == TWIST_M {
			fq2.mulByFq(B2, twist, B)
			fq2.mulByFq(A2, twist, A)
		} else {
			fq2.inverse(B2, twist)
			fq2.mulByFq(B2, B2, B)
			fq2.inverse(A2, twist)
			fq2.mulByFq(A2, A2, A)
		}
	case "BN":
		twist, err := fq2.fromBytes(inputs.nonResidue2)
		if err != nil {
			b.t.Fatal(err)
		}
		if inputs.twistType == TWIST_M {
			fq2.mulByFq(B2, twist, B)
			fq2.mulByFq(A2, twist, A)
		} else {
			fq2.inverse(B2, twist)
			fq2.mulByFq(B2, B2, B)
			fq2.inverse(A2, twist)
			fq2.mulByFq(A2, A2, A)
		}
	case "MNT4":
		twist := fq2.new()
		fq.copy(twist[1], fq.one)
		var twistSquare, twistCube = fq2.new(), fq2.new()
		fq2.square(twistSquare, twist)
		fq2.mul(twistCube, twist, twistSquare)
		fq2.mulByFq(A2, twistSquare, A)
		fq2.mulByFq(B2, twistCube, B)
	default:
		b.t.Fatalf("unrecognized family")
	}
	cache.a2 = A2
	cache.b2 = B2
	return A2, B2
}

func (b *builder) g22TestInstance() group {
	return g22Test{b.g22()}
}

func (b *builder) G22() point {
	return b.cache.G22
}

func (b *builder) g23ab() (*fe3, *fe3) {
	cache := b.cache
	if cache.a3 != nil && cache.b3 != nil {
		return cache.a3, cache.b3
	}
	fq := b.fq()
	fq3 := b.fq3()
	A, B := b.g1ab()
	B3, A3 := fq3.new(), fq3.new()
	switch b.family {
	case "MNT6":
		twist := fq3.new()
		fq.copy(twist[1], fq.one)
		var twistSquare, twistCube = fq3.new(), fq3.new()
		fq3.square(twistSquare, twist)
		fq3.mul(twistCube, twist, twistSquare)
		fq3.mulByFq(A3, twistSquare, A)
		fq3.mulByFq(B3, twistCube, B)
	default:
		b.t.Fatalf("unrecognized family")
	}
	cache.a3 = A3
	cache.b3 = B3
	return A3, B3
}

func (b *builder) g23() *g23 {
	inputs := b.input()
	cache := b.cache
	if cache.g23 != nil {
		return cache.g23
	}
	fq3 := b.fq3()
	A, B := b.g23ab()
	q := new(big.Int).SetBytes(inputs.groupOrder)

	g23, err := newG23(fq3, A, B, q)
	if err != nil {
		b.t.Fatal(err)
	}
	G3, err := g23.fromBytes(inputs.g3one)
	if err != nil {
		b.t.Fatal(err)
	}
	if !g23.isOnCurve(G3) {
		b.t.Fatalf("g3 is not on curve")
	}
	if !g23.checkCorrectSubgroup(G3) {
		b.t.Fatalf("g3 is not on correct subgroup")
	}
	b.cache.G23 = G3
	cache.g23 = g23
	return g23
}

func (b *builder) g23TestInstance() group {
	return g23Test{b.g23()}
}

func (b *builder) G23() point {
	return b.cache.G23
}

func (b *builder) bls() pairingEngine {
	inputs := b.input()
	fq12 := b.fq12()
	g1 := b.g1()
	g2 := b.g22()
	isZNegative := inputs.negz
	z := new(big.Int).SetBytes(inputs.z)
	bls := newBLSInstance(z, isZNegative, inputs.twistType, g1, g2, fq12, true)
	return blsTest{bls}
}

func (b *builder) bn() pairingEngine {
	inputs := b.input()
	fq12 := b.fq12()
	g1 := b.g1()
	g2 := b.g22()
	isUNegative := inputs.negz
	u := new(big.Int).SetBytes(inputs.z)
	bn := newBNInstance(u, isUNegative, inputs.twistType, g1, g2, fq12, true)
	return bnTest{bn}
}

func (b *builder) mnt4() pairingEngine {
	inputs := b.input()
	fq4, fq2, fq := b.fq4(), b.fq2(), b.fq()
	g1 := b.g1()
	g2 := b.g22()
	zneg := inputs.negz
	z := new(big.Int).SetBytes(inputs.z)
	expW0 := new(big.Int).SetBytes(inputs.expW0)
	expW0neg := inputs.expW0neg
	expW1 := new(big.Int).SetBytes(inputs.expW1)
	twist := fq2.new()
	fq.copy(twist[1], fq.one)
	mnt4 := newMNT4Instance(z, zneg, expW0, expW1, expW0neg, fq4, g1, g2, twist)
	return mnt4Test{mnt4}
}

func (b *builder) mnt6() pairingEngine {
	inputs := b.input()
	fq := b.fq()
	fq3 := b.fq3()
	fq6 := b.fq6Q()
	g1 := b.g1()
	g2 := b.g23()
	zneg := inputs.negz
	z := new(big.Int).SetBytes(inputs.z)
	expW0 := new(big.Int).SetBytes(inputs.expW0)
	expW0neg := inputs.expW0neg
	expW1 := new(big.Int).SetBytes(inputs.expW1)

	twist := fq3.new()
	fq.copy(twist[1], fq.one)

	mnt6 := newMNT6Instance(z, zneg, expW0, expW1, expW0neg, fq6, g1, g2, twist)
	return mnt6Test{mnt6}
}

func TestG1(t *testing.T) {

	vectors := []*builder{
		testBuilderFromFile(t, "bls12/256.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/320.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/384.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/448.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/512.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/576.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/640.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/704.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/768.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/832.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/896.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/960.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/1024.json", newBuilderOpt("BLS")),
	}
	for _, v := range vectors {
		g := v.g1TestInstance()
		one := v.G1()
		testG(t, g, one, v.tag)
	}
}

func TestG2(t *testing.T) {
	vectors := []*builder{
		testBuilderFromFile(t, "bls12/256.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/320.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/384.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/448.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/512.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/576.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/640.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/704.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/768.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/832.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/896.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/960.json", newBuilderOpt("BLS")),
		testBuilderFromFile(t, "bls12/1024.json", newBuilderOpt("BLS")),
	}
	for _, v := range vectors {
		g := v.g22TestInstance()
		one := v.G22()
		testG(t, g, one, v.tag)
	}
}

func testG(t *testing.T, g group, one point, tag string) {
	zero := g.zero()
	randPoint := func() point {
		k, err := rand.Int(rand.Reader, g.Q())
		if err != nil {
			panic(err)
		}
		return g.mulScalar(g.new(), one, k)
	}
	t0, t1 := g.new(), g.new()
	testName := tag + "_" + "generator"
	t.Run(testName, func(t *testing.T) {
		if !g.isOnCurve(one) {
			t.Fatalf("generator is not on curve")
		}
	})
	testName = tag + "_" + "serialize"
	t.Run(testName, func(t *testing.T) {
		for i := 0; i < fuz; i++ {
			a := randPoint()
			buf := g.toBytes(a)
			b, err := g.fromBytes(buf)
			if err != nil {
				t.Fatal(err)
			}
			if !g.equal(a, b) {
				t.Fatalf("bad serialization")
			}
			buf = g.toBytesDense(a)
			b, err = g.fromBytesDense(buf)
			if err != nil {
				t.Fatal(err)
			}
			if !g.equal(a, b) {
				t.Fatalf("bad serialization, dense")
			}
		}
	})
	testName = tag + "_" + "addition_properties"
	t.Run(testName, func(t *testing.T) {
		for i := 0; i < fuz; i++ {
			a, b := randPoint(), randPoint()
			g.add(t0, a, zero)
			if !g.equal(t0, a) {
				t.Fatalf("a + 0 == a")
			}
			g.add(t0, zero, zero)
			if !g.equal(t0, zero) {
				t.Fatalf("0 + 0 == 0")
			}
			g.sub(t0, a, zero)
			if !g.equal(t0, a) {
				t.Fatalf("a - 0 == a")
			}
			g.sub(t0, zero, zero)
			if !g.equal(t0, zero) {
				t.Fatalf("0 - 0 == 0")
			}
			g.neg(t0, zero)
			if !g.equal(t0, zero) {
				t.Fatalf("- 0 == 0")
			}
			g.sub(t0, zero, a)
			g.neg(t0, t0)
			if !g.equal(t0, a) {
				t.Fatalf(" - (0 - a) == a")
			}
			g.double(t0, zero)
			if !g.equal(t0, zero) {
				t.Fatalf("2 * 0 == 0")
			}
			g.double(t0, a)
			g.sub(t0, t0, a)
			if !g.equal(t0, a) || !g.isOnCurve(t0) {
				t.Fatalf(" (2 * a) - a == a")
			}
			g.add(t0, a, b)
			g.add(t1, b, a)
			if !g.equal(t0, t1) {
				t.Fatalf("a + b == b + a")
			}
			g.sub(t0, a, b)
			g.sub(t1, b, a)
			g.neg(t1, t1)
			if !g.equal(t0, t1) {
				t.Fatalf("a - b == - ( b - a )")
			}
			c := randPoint()
			g.add(t0, a, b)
			g.add(t0, t0, c)
			g.add(t1, a, c)
			g.add(t1, t1, b)
			if !g.equal(t0, t1) {
				t.Fatalf("(a + b) + c == (a + c ) + b")
			}
			g.sub(t0, a, b)
			g.sub(t0, t0, c)
			g.sub(t1, a, c)
			g.sub(t1, t1, b)
			if !g.equal(t0, t1) {
				t.Fatalf("(a - b) - c == (a - c) -b")
			}
		}
	})
	testName = tag + "_" + "multiplication_properties"
	t.Run(testName, func(t *testing.T) {
		for i := 0; i < fuz; i++ {
			a := randPoint()
			s1, s2, s3 := randScalar(g.Q()), randScalar(g.Q()), randScalar(g.Q())
			sone := big.NewInt(1)
			g.mulScalar(t0, zero, s1)
			if !g.equal(t0, zero) {
				t.Fatalf(" 0 ^ s == 0")
			}
			g.mulScalar(t0, a, sone)
			if !g.equal(t0, a) {
				t.Fatalf(" a ^ 1 == a")
			}
			g.mulScalar(t0, zero, s1)
			if !g.equal(t0, zero) {
				t.Fatalf(" 0 ^ s == a")
			}
			g.mulScalar(t0, a, s1)
			g.mulScalar(t0, t0, s2)
			s3.Mul(s1, s2)
			g.mulScalar(t1, a, s3)
			if !g.equal(t0, t1) {
				t.Errorf(" (a ^ s1) ^ s2 == a ^ (s1 * s2)")
			}
			g.mulScalar(t0, a, s1)
			g.mulScalar(t1, a, s2)
			g.add(t0, t0, t1)
			s3.Add(s1, s2)
			g.mulScalar(t1, a, s3)
			if !g.equal(t0, t1) {
				t.Errorf(" (a ^ s1) + (a ^ s2) == a ^ (s1 + s2)")
			}
		}
	})
	testName = tag + "_" + "wnaf_mul"
	t.Run(testName, func(t *testing.T) {
		a := randPoint()
		s1, s2, s3 := randScalar(g.Q()), randScalar(g.Q()), randScalar(g.Q())
		sone := big.NewInt(1)
		g.wnafMul(t0, zero, s1)
		if !g.equal(t0, zero) {
			t.Fatalf(" 0 ^ s == 0")
		}
		g.wnafMul(t0, a, sone)
		if !g.equal(t0, a) {
			t.Fatalf(" a ^ 1 == a")
		}
		g.wnafMul(t0, zero, s1)
		if !g.equal(t0, zero) {
			t.Fatalf(" 0 ^ s == a")
		}
		g.wnafMul(t0, a, s1)
		g.wnafMul(t0, t0, s2)
		s3.Mul(s1, s2)
		g.wnafMul(t1, a, s3)
		if !g.equal(t0, t1) {
			t.Errorf(" (a ^ s1) ^ s2 == a ^ (s1 * s2)")
		}
		g.wnafMul(t0, a, s1)
		g.wnafMul(t1, a, s2)
		g.add(t0, t0, t1)
		s3.Add(s1, s2)
		g.wnafMul(t1, a, s3)
		if !g.equal(t0, t1) {
			t.Errorf(" (a ^ s1) + (a ^ s2) == a ^ (s1 + s2)")
		}
	})
	testName = tag + "_" + "multi_exp"
	t.Run(testName, func(t *testing.T) {
		count := 1000
		bases := make([]point, count)
		scalars := make([]*big.Int, count)
		for i, j := 0, count-1; i < count; i, j = i+1, j-1 {
			scalars[j] = new(big.Int)
			s, _ := rand.Int(rand.Reader, big.NewInt(10000))
			scalars[j].Set(s)
			bases[i] = g.zero()
			g.mulScalar(bases[i], one, scalars[j])
		}
		expected, tmp := g.zero(), g.zero()
		for i := 0; i < count; i++ {
			g.mulScalar(tmp, bases[i], scalars[i])
			g.add(expected, expected, tmp)
		}
		result := g.zero()
		_, err := g.multiExp(result, bases, scalars)
		if err != nil {
			t.Fatal(err)
		}
		if !g.equal(expected, result) {
			t.Fatalf("bad multi exponentiation")
		}
	})
}

type pairingEngine interface {
	pair(a0, a1 point) (fieldElement, bool)
	multiPair(a1, a2 []point) (fieldElement, bool)
	gt() field // should return the target group
}

type blsTest struct {
	bls *blsInstance
}

func (e blsTest) pair(a0, a1 point) (fieldElement, bool) {
	f, v := e.bls.pair(a0.(*pointG1), a1.(*pointG22))
	return f, v
}

func (e blsTest) multiPair(a1, a2 []point) (fieldElement, bool) {
	numOfPair := len(a1)
	A1, A2 := make([]*pointG1, numOfPair), make([]*pointG22, numOfPair)
	for i := 0; i < numOfPair; i++ {
		A1[i], A2[i] = a1[i].(*pointG1), a2[i].(*pointG22)
	}
	f, v := e.bls.multiPair(A1, A2)
	return f, v
}

func (e blsTest) gt() field {
	return fq12Test{e.bls.fq12}
}

type bnTest struct {
	bn *bnInstance
}

func (e bnTest) pair(a0, a1 point) (fieldElement, bool) {
	f, v := e.bn.pair(a0.(*pointG1), a1.(*pointG22))
	return f, v
}

func (e bnTest) multiPair(a1, a2 []point) (fieldElement, bool) {
	numOfPair := len(a1)
	A1, A2 := make([]*pointG1, numOfPair), make([]*pointG22, numOfPair)
	for i := 0; i < numOfPair; i++ {
		A1[i], A2[i] = a1[i].(*pointG1), a2[i].(*pointG22)
	}
	f, v := e.bn.multiPair(A1, A2)
	return f, v
}

func (e bnTest) gt() field {
	return fq12Test{e.bn.fq12}
}

type mnt4Test struct {
	mnt4 *mnt4Instance
}

func (e mnt4Test) pair(a0, a1 point) (fieldElement, bool) {
	f, v := e.mnt4.pair(a0.(*pointG1), a1.(*pointG22))
	return f, v
}

func (e mnt4Test) multiPair(a1, a2 []point) (fieldElement, bool) {
	numOfPair := len(a1)
	A1, A2 := make([]*pointG1, numOfPair), make([]*pointG22, numOfPair)
	for i := 0; i < numOfPair; i++ {
		A1[i], A2[i] = a1[i].(*pointG1), a2[i].(*pointG22)
	}
	f, v := e.mnt4.multiPair(A1, A2)
	return f, v
}

func (e mnt4Test) gt() field {
	return fq4Test{e.mnt4.fq4}
}

type mnt6Test struct {
	mnt6 *mnt6Instance
}

func (e mnt6Test) pair(a0, a1 point) (fieldElement, bool) {
	f, v := e.mnt6.pair(a0.(*pointG1), a1.(*pointG23))
	return f, v
}

func (e mnt6Test) multiPair(a1, a2 []point) (fieldElement, bool) {
	numOfPair := len(a1)
	A1, A2 := make([]*pointG1, numOfPair), make([]*pointG23, numOfPair)
	for i := 0; i < numOfPair; i++ {
		A1[i], A2[i] = a1[i].(*pointG1), a2[i].(*pointG23)
	}
	f, v := e.mnt6.multiPair(A1, A2)
	return f, v
}

func (e mnt6Test) gt() field {
	return fq6QTest{e.mnt6.fq6}
}

func TestMNT4Pairing(t *testing.T) {
	opts := newBuilderOptPairing("MNT4")
	vectors := []*builder{
		testBuilderFromVector(t, "mnt4_320",
			&vectorJSON{
				FieldOrder:   "0x3bcf7bcd473a266249da7b0548ecaeec9635d1330ea41a9e35e51200e12c90cd65a71660001",
				GroupOrder:   "0x3bcf7bcd473a266249da7b0548ecaeec9635cf44194fb494c07925d6ad3bb4334a400000001",
				A:            "0x02",
				B:            "0x03545a27639415585ea4d523234fc3edd2a2070a085c7b980f4e9cd21a515d4b0ef528ec0fd5",
				G1x:          "0x7a2caf82a1ba85213fe6ca3875aee86aba8f73d69060c4079492b948dea216b5b9c8d2af46",
				G1y:          "0x2db619461cc82672f7f159fec2e89d0148dcc9862d36778c1afd96a71e29cba48e710a48ab2",
				G2x0:         "0x371780491c5660571ff542f2ef89001f205151e12a72cb14f01a931e72dba7903df6c09a9a4",
				G2x1:         "0x4ba59a3f72da165def838081af697c851f002f576303302bb6c02c712c968be32c0ae0a989",
				G2y0:         "0x4b471f33ffaad868a1c47d6605d31e5c4b3b2e0b60ec98f0f610a5aafd0d9522bca4e79f22",
				G2y1:         "0x355d05a1c69a5031f3f81a5c100cb7d982f78ec9cfc3b5168ed8d75c7c484fb61a3cbf0e0f1",
				NonResidue:   "0x11",
				NonResidue20: "0x11",
				NonResidue21: "0x00",
				Z:            "0x1eef5546609756bec2a33f0dc9a1b671660000",
				ExpW0:        "0x1eef5546609756bec2a33f0dc9a1b671660001",
				ExpW1:        "0x01",
			},
			opts,
		),
		testBuilderFromVector(t, "mnt4_753",
			&vectorJSON{
				FieldOrder:   "0x1c4c62d92c41110229022eee2cdadb7f997505b8fafed5eb7e8f96c97d87307fdb925e8a0ed8d99d124d9a15af79db117e776f218059db80f0da5cb537e38685acce9767254a4638810719ac425f0e39d54522cdd119f5e9063de245e8001",
				GroupOrder:   "0x1c4c62d92c41110229022eee2cdadb7f997505b8fafed5eb7e8f96c97d87307fdb925e8a0ed8d99d124d9a15af79db26c5c28c859a99b3eebca9429212636b9dff97634993aa4d6c381bc3f0057974ea099170fa13a4fd90776e240000001",
				A:            "0x02",
				B:            "0x1373684a8c9dcae7a016ac5d7748d3313cd8e39051c596560835df0c9e50a5b59b882a92c78dc537e51a16703ec9855c77fc3d8bb21c8d68bb8cfb9db4b8c8fba773111c36c8b1b4e8f1ece940ef9eaad265458e06372009c9a0491678ef4",
				G1x:          "0x1013b42397c8b004d06f0e98fbc12e8ee65adefcdba683c5630e6b58fb69610b02eab1d43484ddfab28213098b562d799243fb14330903aa64878cfeb34a45d1285da665f5c3f37eb76b86209dcd081ccaef03e65f33d490de480bfee06db",
				G1y:          "0xe3eb479d308664381e7942d6c522c0833f674296169420f1dd90680d0ba6686fc27549d52e4292ea5d611cb6b0df32545b07f281032d0a71f8d485e6907766462e17e8dd55a875bd36fe4cd42cac31c0629fb26c333fe091211d0561d10e",
				G2x0:         "0xf1b7155ed4e903332835a5de0f327aa11b2d74eb8627e3a7b833be42c11d044b5cf0ae49850eeb07d90c77c67256474b2febf924aca0bfa2e4dacb821c91a04fd0165ac8debb2fc1e763a5c32c2c9f572caa85a91c5243ec4b2981af8904",
				G2x1:         "0xd49c264ec663e731713182a88907b8e979ced82ca592777ad052ec5f4b95dc78dc2010d74f82b9e6d066813ed67f3af1de0d5d425da7a19916cf103f102adf5f95b6b62c24c7d186d60b4a103e157e5667038bb2e828a3374d6439526272",
				G2y0:         "0x4b0e2fef08096ebbaddd2d7f288c4acf17b2267e21dc5ce0f925cd5d02209e34d8b69cc94aef5d90af34d3cd98287ace8f1162079cd2d3d7e6c6c2c073c24a359437e75638a1458f4b2face11f8d2a5200b14d6f9dd0fdd407f04be620ee",
				G2y1:         "0xbc1925e7fcb64f6f8697cd5e45fae22f5688e51b30bd984c0acdc67d2962520e80d31966e3ec477909ecca358be2eee53c75f55a6f7d9660dd6f3d4336ad50e8bfa5375791d73b863d59c422c3ea006b013e7afb186f2eaa9df68f4d6098",
				NonResidue:   "0x0d",
				NonResidue20: "0x0d",
				NonResidue21: "0x00",
				Z:            "-0x15474b1d641a3fd86dcbcee5dcda7fe51852c8cbe26e600733b714aa43c31a66b0344c4e2c428b07a7713041ba18000",
				ExpW0:        "-0x15474b1d641a3fd86dcbcee5dcda7fe51852c8cbe26e600733b714aa43c31a66b0344c4e2c428b07a7713041ba17fff",
				ExpW1:        "0x01",
			},
			opts,
		),
	}
	for _, v := range vectors {
		testName := v.tag
		t.Run(testName, func(t *testing.T) {
			// construct the suite
			mnt4 := v.mnt4()
			g1, g2 := v.g1TestInstance(), v.g22TestInstance()
			G1, G2 := v.G1(), v.G22()
			// run tests
			testNonDegeneracy(t, mnt4, g1, g2, G1, G2)
			testBilinearity(t, mnt4, g1, g2, G1, G2)
			testMultiPair(t, mnt4, g1, g2, G1, G2)
		})
	}
}

func TestMNT6Pairing(t *testing.T) {
	opts := newBuilderOptPairing("MNT6")
	vectors := []*builder{
		testBuilderFromVector(t, "mnt6_320",
			&vectorJSON{
				FieldOrder:   "0x3bcf7bcd473a266249da7b0548ecaeec9635cf44194fb494c07925d6ad3bb4334a400000001",
				GroupOrder:   "0x3bcf7bcd473a266249da7b0548ecaeec9635d1330ea41a9e35e51200e12c90cd65a71660001",
				A:            "0x0b",
				B:            "0xd68c7b1dc5dd042e957b71c44d3d6c24e683fc09b420b1a2d263fde47ddba59463d0c65282",
				G1x:          "0x2a4feee24fd2c69d1d90471b2ba61ed56f9bad79b57e0b4c671392584bdadebc01abbc0447d",
				G1y:          "0x32986c245f6db2f82f4e037bf7afd69cbfcbff07fc25d71e9c75e1b97208a333d73d91d3028",
				G2x0:         "0x34f7320a12b56ce532bccb3b44902cbaa723cd60035ada7404b743ad2e644ad76257e4c6813",
				G2x1:         "0xcf41620baa52eec50e61a70ab5b45f681952e0109340fec84f1b2890aba9b15cac5a0c80fa",
				G2x2:         "0x11f99170e10e326433cccb8032fb48007ca3c4e105cf31b056ac767e2cb01258391bd4917ce",
				G2y0:         "0x3a65968f03cc64d62ad05c79c415e07ebd38b363ec48309487c0b83e1717a582c1b60fecc91",
				G2y1:         "0xca5e8427e5db1506c1a24cefc2451ab3accaea5db82dcb0c7117cc74402faa5b2c37685c6e",
				G2y2:         "0xf75d2dd88302c9a4ef941307629a1b3e197277d83abb715f647c2e55a27baf782f5c60e7f7",
				NonResidue:   "0x05",
				NonResidue20: "0x05",
				NonResidue21: "0x00",
				NonResidue22: "0x00",
				Z:            "-0x1eef5546609756bec2a33f0dc9a1b671660000",
				ExpW0:        "-0x1eef5546609756bec2a33f0dc9a1b671660000",
				ExpW1:        "0x01",
			},
			opts,
		),
	}
	for _, v := range vectors {
		testName := v.tag
		t.Run(testName, func(t *testing.T) {
			// construct the suite
			mnt6 := v.mnt6()
			g1, g2 := v.g1TestInstance(), v.g23TestInstance()
			G1, G2 := v.G1(), v.G23()
			// run tests
			testNonDegeneracy(t, mnt6, g1, g2, G1, G2)
			testBilinearity(t, mnt6, g1, g2, G1, G2)
			testMultiPair(t, mnt6, g1, g2, G1, G2)
		})
	}
}

func TestBLSPairing(t *testing.T) {
	opts := newBuilderOptPairing("BLS")
	vectors := []*builder{
		testBuilderFromVector(t, "bls12_381",
			&vectorJSON{
				FieldOrder:   "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab",
				GroupOrder:   "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001",
				A:            "0x00",
				B:            "0x04",
				G1x:          "0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
				G1y:          "0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
				G2x0:         "0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
				G2x1:         "0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
				G2y0:         "0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
				G2y1:         "0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
				NonResidue:   "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaaa",
				NonResidue20: "0x01",
				NonResidue21: "0x01",
				IsDType:      "False",
				Z:            "-0xd201000000010000",
			}, opts),
		testBuilderFromFile(t, "bls12/256.json", opts),
		testBuilderFromFile(t, "bls12/320.json", opts),
		testBuilderFromFile(t, "bls12/384.json", opts),
		testBuilderFromFile(t, "bls12/448.json", opts),
		testBuilderFromFile(t, "bls12/512.json", opts),
		testBuilderFromFile(t, "bls12/576.json", opts),
		testBuilderFromFile(t, "bls12/640.json", opts),
		testBuilderFromFile(t, "bls12/704.json", opts),
		testBuilderFromFile(t, "bls12/768.json", opts),
		testBuilderFromFile(t, "bls12/832.json", opts),
		testBuilderFromFile(t, "bls12/896.json", opts),
		testBuilderFromFile(t, "bls12/960.json", opts),
		testBuilderFromFile(t, "bls12/1024.json", opts),
	}
	for _, v := range vectors {
		testName := v.tag
		t.Run(testName, func(t *testing.T) {
			// construct the suite
			bls := v.bls()
			g1, g2 := v.g1TestInstance(), v.g22TestInstance()
			G1, G2 := v.G1(), v.G22()
			// run tests
			testNonDegeneracy(t, bls, g1, g2, G1, G2)
			testBilinearity(t, bls, g1, g2, G1, G2)
			testMultiPair(t, bls, g1, g2, G1, G2)
		})
	}
}

func TestBNPairing(t *testing.T) {
	opts := newBuilderOptPairing("BN")
	vectors := []*builder{
		testBuilderFromVector(t, "bn254",
			&vectorJSON{
				FieldOrder:   "0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47",
				GroupOrder:   "0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001",
				A:            "0x00",
				B:            "0x03",
				G1x:          "0x01",
				G1y:          "0x02",
				G2x0:         "0x1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed",
				G2x1:         "0x198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2",
				G2y0:         "0x12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
				G2y1:         "0x90689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b",
				NonResidue:   "-0x01",
				NonResidue20: "0x09",
				NonResidue21: "0x01",
				IsDType:      "True",
				Z:            "0x44e992b44a6909f1",
			}, opts),
	}
	for _, v := range vectors {
		testName := v.tag
		t.Run(testName, func(t *testing.T) {
			// construct the suite
			bn := v.bn()
			g1, g2 := v.g1TestInstance(), v.g22TestInstance()
			G1, G2 := v.G1(), v.G22()
			// run tests
			testNonDegeneracy(t, bn, g1, g2, G1, G2)
			testBilinearity(t, bn, g1, g2, G1, G2)
			testMultiPair(t, bn, g1, g2, G1, G2)
		})
	}
}

func testNonDegeneracy(t *testing.T, e pairingEngine, g1, g2 group, G1, G2 point) {
	gt := e.gt()
	// e(g1^a, g2^b) != 1
	// e(g1^a, g2^b) != 0
	{
		f0, ok := e.pair(G1, G2)
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		if gt.isOne(f0) {
			t.Fatalf("pairing result is not expected to be one")
		}
		if gt.isZero(f0) {
			t.Fatalf("pairing result is not expected to be zero")
		}
	}
	// e(g1^a, 0) == 1
	{
		f0, ok := e.pair(g1.zero(), G2)
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		if !gt.isOne(f0) {
			t.Fatalf("pairing result is expected to be one")
		}
	}
	// e(0, g2^b) == 1
	{
		f0, ok := e.pair(G1, g2.zero())
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		if !gt.isOne(f0) {
			t.Fatalf("pairing result is expected to be one")
		}
	}
}

func testBilinearity(t *testing.T, e pairingEngine, g1, g2 group, G1, G2 point) {
	gt := e.gt()
	// e(g1^a, g2^b) == e(g1, g2)^(a * b)
	{
		a, b := big.NewInt(17), big.NewInt(117)
		c := new(big.Int).Mul(a, b)
		f0, ok := e.pair(G1, G2)
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		P1, P2 := g1.new(), g2.new()
		g1.mulScalar(P1, G1, a)
		g2.mulScalar(P2, G2, b)
		g1.affine(P1, P1)
		g2.affine(P2, P2)
		f1, ok := e.pair(P1, P2)
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		gt.exp(f0, f0, c)
		if !gt.equal(f1, f0) {
			t.Fatalf("bad pairing")
		}
	}
	// e(g1^a, g2^b) == e(g1^(a*b), g2)
	{
		a, b := big.NewInt(17), big.NewInt(117)
		c := new(big.Int).Mul(a, b)
		Q1 := g1.new()
		g1.mulScalar(Q1, G1, c)
		g1.affine(Q1, Q1)
		f0, ok := e.pair(Q1, G2)
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		P1, P2 := g1.new(), g2.new()
		g1.mulScalar(P1, G1, a)
		g2.mulScalar(P2, G2, b)
		g1.affine(P1, P1)
		g2.affine(P2, P2)
		f1, ok := e.pair(P1, P2)
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		if !gt.equal(f1, f0) {
			t.Fatalf("bad pairing")
		}
	}
	// e(g1^a, g2^b) == e(g1, g2^(a*b))
	{
		a, b := big.NewInt(17), big.NewInt(117)
		c := new(big.Int).Mul(a, b)
		Q2 := g2.new()
		g2.mulScalar(Q2, G2, c)
		g2.affine(Q2, Q2)
		f0, ok := e.pair(G1, Q2)
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		P1, P2 := g1.new(), g2.new()
		g1.mulScalar(P1, G1, a)
		g2.mulScalar(P2, G2, b)
		g1.affine(P1, P1)
		g2.affine(P2, P2)
		f1, ok := e.pair(P1, P2)
		if !ok {
			t.Fatalf("pairing engine returned no value")
		}
		if !gt.equal(f1, f0) {
			t.Fatalf("bad pairing")
		}
	}
}

func testMultiPair(t *testing.T, e pairingEngine, g1, g2 group, G1, G2 point) {
	// e(g1, g2) ^ S == e(g1^a01, g2^a02) * e(g1^a11, g2^a12) * ... * e(g1^an1, g2^an2)
	// where S = sum(ai1 * ai2)
	numOfPair := 100
	targetExp := new(big.Int)
	q := g1.Q()
	A1, A2 := []point{}, []point{}
	for i := 0; i < numOfPair; i++ {
		a1, a2 := randScalar(q), randScalar(q)
		P1, P2 := g1.new(), g2.new()
		g1.mulScalar(P1, G1, a1)
		g2.mulScalar(P2, G2, a2)
		g1.affine(P1, P1)
		g2.affine(P2, P2)
		A1, A2 = append(A1, P1), append(A2, P2)
		a1.Mul(a1, a2)
		targetExp.Add(targetExp, a1)
	}
	gt := e.gt()
	// LHS
	f0, ok := e.pair(G1, G2)
	if !ok {
		t.Fatalf("pairing engine returned no value")
	}
	gt.exp(f0, f0, targetExp)
	// RHS
	f1, ok := e.multiPair(A1, A2)
	if !ok {
		t.Fatalf("pairing engine returned no value")
	}
	if !gt.equal(f0, f1) {
		t.Fatalf("bad multi pairing")
	}
}
