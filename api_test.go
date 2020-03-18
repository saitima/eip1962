package eip

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
)

type vectorAPI struct {
	input       []byte
	expected    []byte
	expectedErr error
	operation   int
	tag         string
}

func newVectorAPI(input, expected []byte, expectedErr error, operation int, tag string) *vectorAPI {
	return &vectorAPI{input, expected, expectedErr, operation, tag}
}

func (b *builder) encodeG1() *bytes.Buffer {
	data := bytes.NewBuffer([]byte{})
	fq := b.fq()
	g1 := b.g1()
	// modulus len
	data.WriteByte(byte(fq.modulusByteLen))
	// modulus
	data.Write(fq.modulus().Bytes())
	// coeff a
	data.Write(fq.toBytesDense(g1.a))
	// coeff b
	data.Write(fq.toBytesDense(g1.b))
	// group order
	q := g1.q.Bytes()
	data.WriteByte(byte(len(q)))
	data.Write(q)
	return data
}

func (b *builder) encodeG1AddInput() *vectorAPI {
	input := b.encodeG1()
	g := b.g1TestInstance()
	one := b.G1()
	p0, p1 := g.new(), g.new()
	g.mulScalar(p0, one, big.NewInt(2))
	g.mulScalar(p1, one, big.NewInt(3))
	// write rhs, lhs inputs
	input.Write(g.toBytesDense(p0))
	input.Write(g.toBytesDense(p1))
	// expected return
	r := g.new()
	g.add(r, p0, p1)
	output := bytes.NewBuffer([]byte{})
	output.Write(g.toBytesDense(r))
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_G1_ADD)
	return newVectorAPI(input.Bytes(), output.Bytes(), nil, OPERATION_G1_ADD, tag)
}

func (b *builder) encodeG1MulInput() *vectorAPI {
	input := b.encodeG1()
	g1 := b.g1TestInstance()
	// operands
	one, s := b.G1(), big.NewInt(11)
	input.Write(g1.toBytesDense(one))
	// scalar operand padded to order length
	orderLen := len(g1.Q().Bytes())
	input.Write(padBytes(s.Bytes(), orderLen))
	// expected return
	output := bytes.NewBuffer([]byte{})
	r := g1.new()
	g1.mulScalar(r, one, s)
	output.Write(g1.toBytesDense(r))
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_G1_MUL)
	return newVectorAPI(input.Bytes(), output.Bytes(), nil, OPERATION_G1_MUL, tag)
}

func (b *builder) encodeG1MultiExpInput() *vectorAPI {
	input := b.encodeG1()
	g := b.g1TestInstance()
	one := b.G1()
	randScalar := func() *big.Int {
		return randScalar(g.Q())
	}
	randPoint := func() point {
		return g.mulScalar(g.new(), one, randScalar())
	}
	// input len
	n := 2
	input.WriteByte(byte(n))
	scalars := make([]*big.Int, n)
	points := make([]point, n)
	orderLen := len(g.Q().Bytes())
	// encode operands
	for i := 0; i < n; i++ {
		s, p := randScalar(), randPoint()
		scalars[i], points[i] = s, p
		input.Write(g.toBytesDense(p))
		input.Write(padBytes(s.Bytes(), orderLen))
	}
	// expected return
	output := bytes.NewBuffer([]byte{})
	r := g.new()
	g.multiExp(r, points, scalars)
	output.Write(g.toBytesDense(r))
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_G1_MULTIEXP)
	return newVectorAPI(input.Bytes(), output.Bytes(), nil, OPERATION_G1_MULTIEXP, tag)
}

func (b *builder) encodeG22() *bytes.Buffer {
	data := bytes.NewBuffer([]byte{})
	fq, fq2, g2 := b.fq(), b.fq2(), b.g22()
	// modulus len
	data.WriteByte(byte(fq.modulusByteLen))
	// modulus
	data.Write(fq.modulus().Bytes())
	// extension degree
	data.WriteByte(byte(2))
	data.Write(fq.toBytesDense(fq2.nonResidue))
	// coeff a
	data.Write(fq2.toBytesDense(g2.a))
	// coeff b
	data.Write(fq2.toBytesDense(g2.b))
	// group order
	q := g2.q.Bytes()
	data.WriteByte(byte(len(q)))
	data.Write(q)
	return data
}

func (b *builder) encodeG22AddInput() *vectorAPI {
	input := b.encodeG22()
	g := b.g22TestInstance()
	one := b.G22()
	p0, p1 := g.new(), g.new()
	g.mulScalar(p0, one, big.NewInt(2))
	g.mulScalar(p1, one, big.NewInt(3))
	// write rhs, lhs inputs
	input.Write(g.toBytesDense(p0))
	input.Write(g.toBytesDense(p1))
	// expected return
	r := g.new()
	g.add(r, p0, p1)
	output := bytes.NewBuffer([]byte{})
	output.Write(g.toBytesDense(r))
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_G2_ADD)
	return newVectorAPI(input.Bytes(), output.Bytes(), nil, OPERATION_G2_ADD, tag)
}

func (b *builder) encodeG22MulInput() *vectorAPI {
	input := b.encodeG22()
	g := b.g22TestInstance()
	// operands
	one, s := b.G22(), big.NewInt(11)
	input.Write(g.toBytesDense(one))
	// scalar operand padded to order length
	orderLen := len(g.Q().Bytes())
	input.Write(padBytes(s.Bytes(), orderLen))
	// expected return
	output := bytes.NewBuffer([]byte{})
	r := g.new()
	g.mulScalar(r, one, s)
	output.Write(g.toBytesDense(r))
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_G2_MUL)
	return newVectorAPI(input.Bytes(), output.Bytes(), nil, OPERATION_G2_MUL, tag)
}

func (b *builder) encodeG22MultiExpInput() *vectorAPI {
	input := b.encodeG22()
	g := b.g22TestInstance()
	one := b.G22()
	randScalar := func() *big.Int {
		return randScalar(g.Q())
	}
	randPoint := func() point {
		return g.mulScalar(g.new(), one, randScalar())
	}
	// input len
	n := 10
	input.WriteByte(byte(n))
	// inputs
	scalars := make([]*big.Int, n)
	points := make([]point, n)
	orderLen := len(g.Q().Bytes())
	// encode operands
	for i := 0; i < n; i++ {
		s, p := randScalar(), randPoint()
		scalars[i], points[i] = s, p
		input.Write(g.toBytesDense(p))
		input.Write(padBytes(s.Bytes(), orderLen))
	}
	// expected return
	output := bytes.NewBuffer([]byte{})
	r := g.new()
	g.multiExp(r, points, scalars)
	output.Write(g.toBytesDense(r))
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_G2_MULTIEXP)
	return newVectorAPI(input.Bytes(), output.Bytes(), nil, OPERATION_G2_MULTIEXP, tag)
}

func (b *builder) encodeBLSInput() *vectorAPI {
	// encode g1
	input := b.encodeG1()
	// fetch fq, fq2, fq6
	fq, fq2, fq6 := b.fq(), b.fq2(), b.fq6C()
	// fp2_non_residue
	input.Write(fq.toBytesDense(fq2.nonResidue))
	// fp6_non_residue
	input.Write(fq2.toBytesDense(fq6.nonResidue))
	// twist type
	vector := b.input()
	if vector.twistType == TWIST_D {
		input.WriteByte(byte(2))
	} else {
		input.WriteByte(byte(1))
	}
	// loop param
	input.WriteByte(byte(len(vector.z)))
	input.Write(vector.z)
	if vector.negz {
		input.WriteByte(byte(1))
	} else {
		input.WriteByte(byte(0))
	}
	// pair len
	n := 5
	input.WriteByte(byte(n))
	// compute pairs
	g1, g2 := b.g1(), b.g22()
	targetExp := new(big.Int)
	A1, A2 := make([]*pointG1, n), make([]*pointG22, n)
	G1, G2 := b.G1().(*pointG1), b.G22().(*pointG22)
	q := g1.q
	P1, P2 := g1.newPoint(), g2.newPoint()
	for i := 0; i < n-1; i++ {
		A1[i], A2[i] = g1.newPoint(), g2.newPoint()
		a1, a2 := randScalar(q), randScalar(q)
		g1.mulScalar(P1, G1, a1)
		g2.mulScalar(P2, G2, a2)
		g1.affine(P1, P1)
		g2.affine(P2, P2)
		g1.copy(A1[i], P1)
		g2.copy(A2[i], P2)
		a1.Mul(a1, a2)
		targetExp.Add(targetExp, a1)
	}
	A1[n-1], A2[n-1] = g1.newPoint(), g2.newPoint()
	targetExp.Mod(targetExp, q).Sub(q, targetExp)
	g1.mulScalar(P1, G1, targetExp)
	g1.affine(P1, P1)
	g1.copy(A1[n-1], P1)
	g2.copy(A2[n-1], G2)
	// write pairs
	// subGroupCheck := false
	for i := 0; i < n; i++ {
		input.WriteByte(byte(0)) // subgroup check for g1
		input.Write(g1.toBytesDense(A1[i]))
		input.WriteByte(byte(0)) // subgroup check for g2
		input.Write(g2.toBytesDense(A2[i]))
	}
	output := pairingSuccess
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_BLS12PAIR)
	return newVectorAPI(input.Bytes(), output, nil, OPERATION_BLS12PAIR, tag)
}

func (b *builder) encodeBNInput() *vectorAPI {
	// encode g1
	input := b.encodeG1()
	// fetch fq, fq2, fq6
	fq, fq2, fq6 := b.fq(), b.fq2(), b.fq6C()
	// fp2_non_residue
	input.Write(fq.toBytesDense(fq2.nonResidue))
	// fp6_non_residue
	input.Write(fq2.toBytesDense(fq6.nonResidue))
	// twist type
	vector := b.input()
	if vector.twistType == TWIST_D {
		input.WriteByte(byte(2))
	} else {
		input.WriteByte(byte(1))
	}
	// loop param
	input.WriteByte(byte(len(vector.z)))
	input.Write(vector.z)
	if vector.negz {
		input.WriteByte(byte(1))
	} else {
		input.WriteByte(byte(0))
	}
	// pair len
	n := 5
	input.WriteByte(byte(n))
	// compute pairs
	g1, g2 := b.g1(), b.g22()
	targetExp := new(big.Int)
	A1, A2 := make([]*pointG1, n), make([]*pointG22, n)
	G1, G2 := b.G1().(*pointG1), b.G22().(*pointG22)
	q := g1.q
	P1, P2 := g1.newPoint(), g2.newPoint()
	for i := 0; i < n-1; i++ {
		A1[i], A2[i] = g1.newPoint(), g2.newPoint()
		a1, a2 := randScalar(q), randScalar(q)
		g1.mulScalar(P1, G1, a1)
		g2.mulScalar(P2, G2, a2)
		g1.affine(P1, P1)
		g2.affine(P2, P2)
		g1.copy(A1[i], P1)
		g2.copy(A2[i], P2)
		a1.Mul(a1, a2)
		targetExp.Add(targetExp, a1)
	}
	A1[n-1], A2[n-1] = g1.newPoint(), g2.newPoint()
	targetExp.Mod(targetExp, q).Sub(q, targetExp)
	g1.mulScalar(P1, G1, targetExp)
	g1.affine(P1, P1)
	g1.copy(A1[n-1], P1)
	g2.copy(A2[n-1], G2)
	// write pairs
	// subGroupCheck := false
	for i := 0; i < n; i++ {
		input.WriteByte(byte(0)) // subgroup check for g1
		input.Write(g1.toBytesDense(A1[i]))
		input.WriteByte(byte(0)) // subgroup check for g2
		input.Write(g2.toBytesDense(A2[i]))
	}
	output := pairingSuccess
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_BNPAIR)
	return newVectorAPI(input.Bytes(), output, nil, OPERATION_BNPAIR, tag)
}

func (b *builder) encodeMNT4Input() *vectorAPI {
	// encode g1
	input := b.encodeG1()
	// fetch fq, fq2, fq6
	fq, fq2 := b.fq(), b.fq2()
	// fp2_non_residue
	input.Write(fq.toBytesDense(fq2.nonResidue))
	// loop params
	vector := b.input()
	// ate loop
	input.WriteByte(byte(len(vector.z)))
	input.Write(vector.z)
	// sign of ate loop
	if vector.negz {
		input.WriteByte(byte(1))
	} else {
		input.WriteByte(byte(0))
	}
	// exp w0
	input.WriteByte(byte(len(vector.expW0)))
	input.Write(vector.expW0)
	// exp w1
	input.WriteByte(byte(len(vector.expW1)))
	input.Write(vector.expW1)
	// sign of exp w0
	if vector.expW0neg {
		input.WriteByte(byte(1))
	} else {
		input.WriteByte(byte(0))
	}
	// pair len
	n := 5
	input.WriteByte(byte(n))
	// compute pairs
	g1, g2 := b.g1(), b.g22()
	targetExp := new(big.Int)
	A1, A2 := make([]*pointG1, n), make([]*pointG22, n)
	G1, G2 := b.G1().(*pointG1), b.G22().(*pointG22)
	q := g1.q
	P1, P2 := g1.newPoint(), g2.newPoint()
	for i := 0; i < n-1; i++ {
		A1[i], A2[i] = g1.newPoint(), g2.newPoint()
		a1, a2 := randScalar(q), randScalar(q)
		g1.mulScalar(P1, G1, a1)
		g2.mulScalar(P2, G2, a2)
		g1.affine(P1, P1)
		g2.affine(P2, P2)
		g1.copy(A1[i], P1)
		g2.copy(A2[i], P2)
		a1.Mul(a1, a2)
		targetExp.Add(targetExp, a1)
	}
	A1[n-1], A2[n-1] = g1.newPoint(), g2.newPoint()
	targetExp.Mod(targetExp, q).Sub(q, targetExp)
	g1.mulScalar(P1, G1, targetExp)
	g1.affine(P1, P1)
	g1.copy(A1[n-1], P1)
	g2.copy(A2[n-1], G2)
	// write pairs
	// subGroupCheck := false
	for i := 0; i < n; i++ {
		input.WriteByte(byte(0)) // subgroup check for g1
		input.Write(g1.toBytesDense(A1[i]))
		input.WriteByte(byte(0)) // subgroup check for g2
		input.Write(g2.toBytesDense(A2[i]))
	}
	output := pairingSuccess
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_MNT4PAIR)
	return newVectorAPI(input.Bytes(), output, nil, OPERATION_MNT4PAIR, tag)
}

func (b *builder) encodeMNT6Input() *vectorAPI {
	// encode g1
	input := b.encodeG1()
	// fetch fq, fq2, fq6
	fq, fq3 := b.fq(), b.fq3()
	// fp3_non_residue
	input.Write(fq.toBytesDense(fq3.nonResidue))
	// loop params
	vector := b.input()
	// ate loop
	input.WriteByte(byte(len(vector.z)))
	input.Write(vector.z)
	// sign of ate loop
	if vector.negz {
		input.WriteByte(byte(1))
	} else {
		input.WriteByte(byte(0))
	}
	// exp w0
	input.WriteByte(byte(len(vector.expW0)))
	input.Write(vector.expW0)
	// exp w1
	input.WriteByte(byte(len(vector.expW1)))
	input.Write(vector.expW1)
	// sign of exp w0
	if vector.expW0neg {
		input.WriteByte(byte(1))
	} else {
		input.WriteByte(byte(0))
	}
	// pair len
	n := 5
	input.WriteByte(byte(n))
	// compute pairs
	g1, g2 := b.g1(), b.g23()
	targetExp := new(big.Int)
	A1, A2 := make([]*pointG1, n), make([]*pointG23, n)
	G1, G2 := b.G1().(*pointG1), b.G23().(*pointG23)
	q := g1.q
	P1, P2 := g1.newPoint(), g2.newPoint()
	for i := 0; i < n-1; i++ {
		A1[i], A2[i] = g1.newPoint(), g2.newPoint()
		a1, a2 := randScalar(q), randScalar(q)
		g1.mulScalar(P1, G1, a1)
		g2.mulScalar(P2, G2, a2)
		g1.affine(P1, P1)
		g2.affine(P2, P2)
		g1.copy(A1[i], P1)
		g2.copy(A2[i], P2)
		a1.Mul(a1, a2)
		targetExp.Add(targetExp, a1)
	}
	A1[n-1], A2[n-1] = g1.newPoint(), g2.newPoint()
	targetExp.Mod(targetExp, q).Sub(q, targetExp)
	g1.mulScalar(P1, G1, targetExp)
	g1.affine(P1, P1)
	g1.copy(A1[n-1], P1)
	g2.copy(A2[n-1], G2)
	// write pairs
	// subGroupCheck := false
	for i := 0; i < n; i++ {
		input.WriteByte(byte(0)) // subgroup check for g1
		input.Write(g1.toBytesDense(A1[i]))
		input.WriteByte(byte(0)) // subgroup check for g2
		input.Write(g2.toBytesDense(A2[i]))
	}
	output := pairingSuccess
	tag := fmt.Sprintf("%s_%d", b.tag, OPERATION_MNT6PAIR)
	return newVectorAPI(input.Bytes(), output, nil, OPERATION_MNT6PAIR, tag)
}

func TestAPIG1Add(t *testing.T) {
	vectors := []*vectorAPI{
		testBuilderFromFile(t, "bls12/256.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/320.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/384.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/448.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/512.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/576.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/640.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/704.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/768.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/832.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/896.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/960.json", newBuilderOpt("BLS")).encodeG1AddInput(),
		testBuilderFromFile(t, "bls12/1024.json", newBuilderOpt("BLS")).encodeG1AddInput(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error")
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result")
		}
	}
}

func TestAPIG1Mul(t *testing.T) {
	vectors := []*vectorAPI{
		testBuilderFromFile(t, "bls12/256.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/320.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/384.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/448.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/512.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/576.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/640.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/704.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/768.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/832.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/896.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/960.json", newBuilderOpt("BLS")).encodeG1MulInput(),
		testBuilderFromFile(t, "bls12/1024.json", newBuilderOpt("BLS")).encodeG1MulInput(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error")
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result")
		}
	}
}

func TestAPIG1MultiExp(t *testing.T) {
	vectors := []*vectorAPI{
		testBuilderFromFile(t, "bls12/256.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/320.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/384.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/448.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/512.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/576.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/640.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/704.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/768.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/832.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/896.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/960.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
		testBuilderFromFile(t, "bls12/1024.json", newBuilderOpt("BLS")).encodeG1MultiExpInput(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Fatal("not have expected error", v.tag)
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result", v.tag)
		}
	}
}

func TestAPIG2Add(t *testing.T) {
	vectors := []*vectorAPI{
		testBuilderFromFile(t, "bls12/256.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/320.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/384.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/448.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/512.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/576.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/640.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/704.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/768.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/832.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/896.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/960.json", newBuilderOpt("BLS")).encodeG22AddInput(),
		testBuilderFromFile(t, "bls12/1024.json", newBuilderOpt("BLS")).encodeG22AddInput(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error")
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result")
		}
	}
}

func TestAPIG2Mul(t *testing.T) {
	vectors := []*vectorAPI{
		testBuilderFromFile(t, "bls12/256.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/320.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/384.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/448.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/512.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/576.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/640.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/704.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/768.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/832.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/896.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/960.json", newBuilderOpt("BLS")).encodeG22MulInput(),
		testBuilderFromFile(t, "bls12/1024.json", newBuilderOpt("BLS")).encodeG22MulInput(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error")
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result")
		}
	}
}

func TestAPIG2MultiExp(t *testing.T) {
	vectors := []*vectorAPI{
		testBuilderFromFile(t, "bls12/256.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/320.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/384.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/448.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/512.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/576.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/640.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/704.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/768.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/832.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/896.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/960.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
		testBuilderFromFile(t, "bls12/1024.json", newBuilderOpt("BLS")).encodeG22MultiExpInput(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error", v.tag)
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result", v.tag)
		}
	}
}

func TestAPIBLS(t *testing.T) {
	opts := newBuilderOptPairing("BLS")
	vectors := []*vectorAPI{
		testBuilderFromFile(t, "bls12/256.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/320.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/384.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/448.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/512.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/576.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/640.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/704.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/768.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/832.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/896.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/960.json", opts).encodeBLSInput(),
		testBuilderFromFile(t, "bls12/1024.json", opts).encodeBLSInput(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error", v.tag)
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result", v.tag)
		}
	}
}

func TestAPIBN(t *testing.T) {
	opts := newBuilderOptPairing("BN")
	vectors := []*vectorAPI{
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
			},
			opts,
		).encodeBNInput(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error", v.tag)
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result", v.tag)
		}
	}
}

func TestAPIMNT4(t *testing.T) {
	opts := newBuilderOptPairing("MNT4")
	vectors := []*vectorAPI{
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
		).encodeMNT4Input(),
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
		).encodeMNT4Input(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error", v.tag)
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result", v.tag)
		}
	}
}

func TestAPIMNT6(t *testing.T) {
	opts := newBuilderOptPairing("MNT6")
	vectors := []*vectorAPI{
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
		).encodeMNT6Input(),
	}
	api := newAPI()
	for _, v := range vectors {
		result, err := api.Run(v.operation, v.input)
		if err != v.expectedErr {
			t.Log(err)
			t.Fatal("not have expected error", v.tag)
		}
		if !bytes.Equal(result, v.expected) {
			t.Fatal("not have expected result", v.tag)
		}
	}
}
