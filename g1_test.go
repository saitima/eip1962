package fp

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func initG1() *g1 {
	// base field
	p, ok := new(big.Int).SetString("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab", 16)
	if !ok {
		panic("invalid modulus") // @TODO
	}

	q, ok := new(big.Int).SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)
	if !ok {
		panic("invalid g1 order")
	}

	f := newField(p.Bytes())
	a := bytes_(48, "0x00")
	b := bytes_(48, "0x04")
	g1, err := newG1(f, a, b, q.Bytes())
	if err != nil {
		panic(err)
	}
	return g1
}

func TestG1FromBytesToBytes(t *testing.T) {
	in := bytes_(48,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	)

	g1 := initG1()
	p, err := g1.fromBytes(in)
	if err != nil {
		t.Fatal(err)
	}

	if !g1.isOnCurve(p) {
		t.Fatalf("invalid point")
	}
	out := g1.toBytes(p)
	q, err := g1.fromBytes(out)
	if err != nil {
		t.Fatal(err)
	}
	if !g1.equal(p, q) {
		t.Logf("\np: %x\n", p)
		t.Logf("\nq: %x\n", q)
		t.Logf("invalid out ")
	}
}

func TestPointCopy(t *testing.T) {
	g1 := initG1()

	in := bytes_(48,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	)

	p, err := g1.fromBytes(in)
	if err != nil {
		t.Fatal(err)
	}

	q := g1.newPoint()
	g1.copy(q, p)
	if !g1.equal(q, p) {
		t.Fatalf("not equal :/")
	}
}

func TestG1Equality(t *testing.T) {
	g1 := initG1()
	one := g1Generator()
	two := g1.zero()
	g1.double(two, one)
	if g1.equal(two, one) {
		t.Logf("1: %s\n", g1.f.toString(one[0]))
		t.Logf("2: %s\n", g1.f.toString(two[0]))
		t.Fatalf("bad equality")
	}
}

func TestG1IsOnCurve(t *testing.T) {
	g1 := initG1()
	p := g1Generator()
	if !g1.isOnCurve(p) {
		t.Fatalf("point is not on the curve")
	}
}

func TestG1Add(t *testing.T) {
	g1 := initG1()

	one := g1Generator()
	negOne, tmp := g1.newPoint(), g1.newPoint()
	g1.neg(negOne, one)

	g1.add(tmp, one, g1.zero())
	if !g1.equal(tmp, one) {
		t.Fatalf("invalid add([1]P + [0]P != [1]P)")
	}
	g1.add(negOne, negOne, one)
	if !g1.equal(negOne, g1.zero()) {
		t.Fatalf("invalid add([1]P -[1]P != [0]P)")
	}
}

func TestG1Double(t *testing.T) {
	g1 := initG1()
	one := g1Generator()
	two, twoExpected := g1.newPoint(), g1.newPoint()
	g1.add(twoExpected, one, one)
	g1.double(two, one)
	if !g1.equal(two, twoExpected) {
		t.Fatalf("invalid double")
	}
}

func TestG1MulScalar(t *testing.T) {
	g1 := initG1()
	one := g1Generator()
	three, expected := g1.newPoint(), g1.newPoint()
	g1.double(expected, one)
	g1.add(expected, expected, one)

	g1.mulScalar(three, one, big.NewInt(3))
	if !g1.equal(three, expected) {
		t.Fatalf("invalid mul scalar")
	}
}

func TestG1MultiExp(t *testing.T) {
	g := initG1()
	one := g1Generator()

	count := 1000
	bases := make([]*pointG1, count)
	scalars := make([]*big.Int, count)
	// prepare bases
	// bases: S[0]*G, S[1]*G ... S[n-1]*G
	for i, j := 0, count-1; i < count; i, j = i+1, j-1 {
		// TODO : make sure that s is unique
		scalars[j], _ = rand.Int(rand.Reader, big.NewInt(10000))
		bases[i] = g.zero()
		g.mulScalar(bases[i], one, scalars[j])
	}

	// expected
	//  S[n-1]*P[1], S[n-2]*P[2] ... S[0]*P[n-1]
	expected, tmp := g.zero(), g.zero()
	for i := 0; i < count; i++ {
		g.mulScalar(tmp, bases[i], scalars[i])
		g.add(expected, expected, tmp)
	}
	result := g.zero()
	g.multiExp(result, bases, scalars)
	if !g.equal(expected, result) {
		t.Fatalf(":/")
	}
}

func g1Generator() *pointG1 {
	g1 := initG1()
	in := bytes_(48,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	)

	p, err := g1.fromBytes(in)
	if err != nil {
		panic(err)
	}
	if !g1.isOnCurve(p) {
		panic("point is not on the curve")
	}
	return p
}
