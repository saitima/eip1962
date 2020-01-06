package fp

import (
	"crypto/rand"
	"math/big"
	"testing"
)

var (
	fq2One = bytes_(48,
		"0x01",
		"0x00",
	)
	fq3One = bytes_(40,
		"0x01",
		"0x00",
		"0x00",
	)
	fq6CubicOne = bytes_(48,
		"0x01", "0x00",
		"0x00", "0x00",
		"0x00", "0x00",
	)
	fq12One = bytes_(48,
		"0x01", "0x00",
		"0x00", "0x00",
		"0x00", "0x00",
		"0x00", "0x00",
		"0x00", "0x00",
		"0x00", "0x00",
	)
	bigZero, bigOne = big.NewInt(0), big.NewInt(1)
)

func TestG1(t *testing.T) {
	// base field
	modulus, ok := new(big.Int).SetString("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab", 16)
	if !ok {
		panic("invalid modulus") // @TODO
	}
	q, ok := new(big.Int).SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)
	if !ok {
		panic("invalid g1 order")
	}
	f := newField(modulus.Bytes())
	a := bytes_(48, "0x00")
	b := bytes_(48, "0x04")
	g, err := newG1(f, a, b, q.Bytes())
	if err != nil {
		panic(err)
	}
	zero := g.zero()
	oneBytes := bytes_(48,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	)
	actual, expected := g.newPoint(), g.zero()
	one := g.newPoint()
	t.Run("FromBytes & ToBytes", func(t *testing.T) {
		one, err = g.fromBytes(oneBytes)
		if err != nil {
			t.Fatal(err)
		}
		q, err := g.fromBytes(
			g.toBytes(one),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !g.equal(one, q) {
			t.Logf("invalid out ")
		}
	})

	t.Run("Is on curve", func(t *testing.T) {
		if !g.isOnCurve(one) {
			t.Fatalf("point is not on the curve")
		}
	})

	t.Run("Copy", func(t *testing.T) {
		q := g.newPoint()
		g.copy(q, one)
		if !g.equal(q, one) {
			t.Fatalf("bad point copy")
		}
	})

	t.Run("Equality", func(t *testing.T) {
		if !g.equal(zero, zero) {
			t.Fatal("bad equality 1")
		}
		if !g.equal(one, one) {
			t.Fatal("bad equality 2")
		}
		if g.equal(one, zero) {
			// t.Fatal("bad equality 3") // TODO: affine equality
		}
	})

	t.Run("Affine", func(t *testing.T) {
		g.double(actual, one)
		g.sub(expected, actual, one)
		g.affine(expected, expected)
		if !g.equal(expected, one) {
			t.Fatal("invalid affine ops")
		}
	})

	t.Run("Addition", func(t *testing.T) {
		g.add(actual, zero, zero)
		if !g.equal(actual, zero) {
			t.Fatal("bad addition 1")
		}
		g.add(actual, one, zero)
		if !g.equal(actual, one) {
			t.Fatal("bad addition 2")
		}
		g.add(actual, zero, one)
		if !g.equal(actual, one) {
			t.Fatal("bad addition 3")
		}
	})
	t.Run("Substraction", func(t *testing.T) {
		g.sub(actual, zero, zero)
		if !g.equal(actual, zero) {
			t.Fatal("bad substraction 1")
		}
		g.sub(actual, one, zero)
		if !g.equal(actual, one) {
			t.Fatal("bad substraction 2")
		}
		g.sub(actual, one, one)
		if !g.equal(actual, zero) {
			t.Fatal("bad substraction 3")
		}
	})
	t.Run("Negation", func(t *testing.T) {
		g.neg(actual, zero)
		if !g.equal(actual, zero) {
			t.Fatal("bad negation 1")
		}
		g.neg(actual, one)
		g.sub(expected, zero, one)
		if !g.equal(actual, expected) {
			t.Fatal("bad negation 2")
		}
	})

	t.Run("Doubling", func(t *testing.T) {
		g.double(actual, zero)
		if !g.equal(actual, zero) {
			t.Fatal("bad doubling 1")
		}
		g.add(expected, one, one)
		g.double(actual, one)
		if !g.equal(actual, expected) {
			t.Fatal("bad addition 2")
		}
	})

	t.Run("Scalar Multiplication", func(t *testing.T) {
		g.mulScalar(actual, zero, bigZero)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 1")
		}
		g.mulScalar(actual, zero, bigOne)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 2")
		}
		g.mulScalar(actual, one, bigZero)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 3")
		}
		g.mulScalar(actual, one, bigOne)
		if !g.equal(actual, one) {
			t.Fatal("bad scalar multiplication 4")
		}
	})

	t.Run("Multi Exponentiation", func(t *testing.T) {
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
			t.Fatalf("bad multi-exponentiation")
		}
	})
}

func TestFq2(t *testing.T) {
	modulusBytes := bytes_(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	f := newField(modulusBytes)
	fq2, err := newFq2(f, nil)
	if err != nil {
		panic(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	zero := fq2.zero()
	one := fq2.one()
	actual := fq2.newElement()
	expected := fq2.newElement()

	t.Run("FromBytes & ToBytes", func(t *testing.T) {
		a, err := fq2.fromBytes(fq2One)
		if err != nil {
			t.Fatal(err)
		}
		if !fq2.equal(a, fq2.one()) {
			t.Fatalf("bad fromBytes")
		}
		b, err := fq2.fromBytes(
			fq2.toBytes(a),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !fq2.equal(a, b) {
			t.Fatalf("not equal")
		}
	})

	t.Run("Addition", func(t *testing.T) {
		fq2.add(actual, zero, zero)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad add")
		}
		fq2.add(actual, one, zero)
		if !fq2.equal(actual, one) {
			t.Fatalf("bad add")
		}
		fq2.add(actual, zero, zero)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad add")
		}
	})
	t.Run("Substraction", func(t *testing.T) {
		fq2.sub(actual, zero, zero)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad substraction 1")
		}
		fq2.sub(actual, one, zero)
		if !fq2.equal(actual, one) {
			t.Fatalf("bad substraction 2")
		}
		fq2.sub(actual, one, one)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad substraction 3")
		}
	})

	t.Run("Negation", func(t *testing.T) {
		fq2.sub(expected, zero, one)
		fq2.neg(actual, one)
		if !fq2.equal(expected, actual) {
			t.Fatalf("bad negation")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		fq2.mul(actual, zero, zero)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad multiplication 1")
		}
		fq2.mul(actual, one, zero)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq2.mul(actual, zero, one)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq2.mul(actual, one, one)
		if !fq2.equal(actual, one) {
			t.Fatalf("bad multiplication 2")
		}
	})

	t.Run("Squaring", func(t *testing.T) {
		fq2.square(actual, zero)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad squaring 1")
		}
		fq2.square(actual, one)
		if !fq2.equal(actual, one) {
			t.Fatalf("bad squaring 2")
		}
		fq2.double(expected, one)
		fq2.square(actual, expected)
		fq2.mul(expected, expected, expected)
		if !fq2.equal(expected, actual) {
			t.Fatalf("bad squaring 3")
		}
	})

	t.Run("Inverse", func(t *testing.T) {
		fq2.inverse(actual, zero)
		if !fq2.equal(actual, zero) {
			t.Fatalf("bad inversion 1")
		}
		fq2.inverse(actual, one)
		if !fq2.equal(actual, one) {
			t.Fatalf("bad inversion 2")
		}
		fq2.double(expected, one)
		fq2.inverse(actual, expected)
		fq2.mul(expected, actual, expected)
		if !fq2.equal(expected, one) {
			t.Fatalf("bad inversion 3")
		}
	})

	t.Run("Exponentiation", func(t *testing.T) {
		fq2.exp(actual, zero, bigZero)
		if !fq2.equal(actual, one) {
			t.Fatalf("bad exponentiation 1")
		}
		fq2.exp(actual, zero, bigOne)
		if !fq2.equal(actual, zero) {
			t.Logf("actual %s\n", fq2.toString(actual))
			t.Fatalf("bad exponentiation 2")
		}
		fq2.exp(actual, one, bigZero)
		if !fq2.equal(actual, one) {
			t.Fatalf("bad exponentiation 3")
		}
		fq2.exp(actual, one, bigOne)
		if !fq2.equal(actual, one) {
			t.Fatalf("bad exponentiation 4")
		}
		fq2.double(expected, one)
		fq2.exp(actual, expected, big.NewInt(2))
		fq2.square(expected, expected)
		if !fq2.equal(expected, actual) {
			t.Fatalf("bad exponentiation 4")
		}
	})
}

func TestG22(t *testing.T) {
	// base field
	modulus, ok := new(big.Int).SetString("4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787", 10)
	if !ok {
		panic("invalid modulus") // @TODO
	}
	q, ok := new(big.Int).SetString("52435875175126190479447740508185965837690552500527637822603658699938581184513", 10)
	if !ok {
		panic("invalid g1 order")
	}
	f := newField(modulus.Bytes())
	fq2, err := newFq2(f, nil)
	if err != nil {
		panic(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	g, err := newG22(fq2, nil, nil, q.Bytes())
	if err != nil {
		panic(err)
	}
	zero := g.zero()
	oneBytes := bytes_(48,
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
	)
	actual, expected := g.newPoint(), g.zero()
	one := g.newPoint()
	t.Run("FromBytes & ToBytes", func(t *testing.T) {
		one, err = g.fromBytes(oneBytes)
		if err != nil {
			t.Fatal(err)
		}
		q, err := g.fromBytes(
			g.toBytes(one),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !g.equal(one, q) {
			t.Logf("invalid out ")
		}
	})
	b, err := f.newFieldElementFromBytes(bytes_(48, "0x04"))
	if err != nil {
		t.Fatal(err)
	}
	a2, b2 := fq2.zero(), fq2.newElement()
	f.cpy(b2[0], b)
	f.cpy(b2[1], b)
	fq2.copy(g.a, a2)
	fq2.copy(g.b, b2)

	t.Run("Is on curve", func(t *testing.T) {
		if !g.isOnCurve(one) {
			t.Fatalf("point is not on the curve")
		}
	})

	t.Run("Copy", func(t *testing.T) {
		q := g.newPoint()
		g.copy(q, one)
		if !g.equal(q, one) {
			t.Fatalf("bad point copy")
		}
	})

	t.Run("Equality", func(t *testing.T) {
		if !g.equal(zero, zero) {
			t.Fatal("bad equality 1")
		}
		if !g.equal(one, one) {
			t.Fatal("bad equality 2")
		}
		if g.equal(one, zero) {
			// t.Fatal("bad equality 3") // TODO: affine equality
		}
	})

	t.Run("Affine", func(t *testing.T) {
		g.double(actual, one)
		g.sub(expected, actual, one)
		g.affine(expected, expected)
		if !g.equal(expected, one) {
			t.Fatal("invalid affine ops")
		}
	})

	t.Run("Addition", func(t *testing.T) {
		g.add(actual, zero, zero)
		if !g.equal(actual, zero) {
			t.Fatal("bad addition 1")
		}
		g.add(actual, one, zero)
		if !g.equal(actual, one) {
			t.Fatal("bad addition 2")
		}
		g.add(actual, zero, one)
		if !g.equal(actual, one) {
			t.Fatal("bad addition 3")
		}
	})

	t.Run("Substraction", func(t *testing.T) {
		g.sub(actual, zero, zero)
		if !g.equal(actual, zero) {
			t.Fatal("bad substraction 1")
		}
		g.sub(actual, one, zero)
		if !g.equal(actual, one) {
			t.Fatal("bad substraction 2")
		}
		g.sub(actual, one, one)
		if !g.equal(actual, zero) {
			t.Fatal("bad substraction 3")
		}
	})
	t.Run("Negation", func(t *testing.T) {
		g.neg(actual, zero)
		if !g.equal(actual, zero) {
			t.Fatal("bad negation 1")
		}
		g.neg(actual, one)
		g.sub(expected, zero, one)
		if !g.equal(actual, expected) {
			t.Fatal("bad negation 2")
		}
	})

	t.Run("Doubling", func(t *testing.T) {
		g.double(actual, zero)
		if !g.equal(actual, zero) {
			t.Fatal("bad doubling 1")
		}
		g.add(expected, one, one)
		g.double(actual, one)
		if !g.equal(actual, expected) {
			t.Fatal("bad addition 2")
		}
	})

	t.Run("Scalar Multiplication", func(t *testing.T) {
		g.mulScalar(actual, zero, bigZero)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 1")
		}
		g.mulScalar(actual, zero, bigOne)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 2")
		}
		g.mulScalar(actual, one, bigZero)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 3")
		}
		g.mulScalar(actual, one, bigOne)
		if !g.equal(actual, one) {
			t.Fatal("bad scalar multiplication 4")
		}
	})

	t.Run("Multi Exponentiation", func(t *testing.T) {
		count := 1000
		bases := make([]*pointG22, count)
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
			t.Fatalf("bad multi-exponentiation")
		}
	})
}

func TestFq6Cubic(t *testing.T) {
	modulusBytes := bytes_(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	f := newField(modulusBytes)

	fq2, err := newFq2(f, nil)
	if err != nil {
		panic(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		panic(err)
	}
	fq6.calculateFrobeniusCoeffs()
	f.cpy(fq6.nonResidue[0], f.one)
	f.cpy(fq6.nonResidue[1], f.one)

	zero := fq6.zero()
	one := fq6.one()
	actual := fq6.newElement()
	expected := fq6.newElement()

	t.Run("FromBytes & ToBytes", func(t *testing.T) {
		a, err := fq6.fromBytes(fq6CubicOne)
		if err != nil {
			t.Fatal(err)
		}
		if !fq6.equal(a, fq6.one()) {
			t.Fatalf("bad fromBytes")
		}
		b, err := fq6.fromBytes(
			fq6.toBytes(a),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !fq6.equal(a, b) {
			t.Fatalf("not equal")
		}
	})

	t.Run("Addition", func(t *testing.T) {
		fq6.add(actual, zero, zero)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad add")
		}
		fq6.add(actual, one, zero)
		if !fq6.equal(actual, one) {
			t.Fatalf("bad add")
		}
		fq6.add(actual, zero, zero)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad add")
		}
	})
	t.Run("Substraction", func(t *testing.T) {
		fq6.sub(actual, zero, zero)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad substraction 1")
		}
		fq6.sub(actual, one, zero)
		if !fq6.equal(actual, one) {
			t.Fatalf("bad substraction 2")
		}
		fq6.sub(actual, one, one)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad substraction 3")
		}
	})

	t.Run("Negation", func(t *testing.T) {
		fq6.sub(expected, zero, one)
		fq6.neg(actual, one)
		if !fq6.equal(expected, actual) {
			t.Fatalf("bad negation")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		fq6.mul(actual, zero, zero)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad multiplication 1")
		}
		fq6.mul(actual, one, zero)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq6.mul(actual, zero, one)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq6.mul(actual, one, one)
		if !fq6.equal(actual, one) {
			t.Fatalf("bad multiplication 2")
		}
	})

	t.Run("Squaring", func(t *testing.T) {
		fq6.square(actual, zero)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad squaring 1")
		}
		fq6.square(actual, one)
		if !fq6.equal(actual, one) {
			t.Fatalf("bad squaring 2")
		}
		fq6.double(expected, one)
		fq6.square(actual, expected)
		fq6.mul(expected, expected, expected)
		if !fq6.equal(expected, actual) {
			t.Fatalf("bad squaring 3")
		}
	})

	t.Run("Inverse", func(t *testing.T) {
		fq6.inverse(actual, zero)
		if !fq6.equal(actual, zero) {
			t.Fatalf("bad inversion 1")
		}
		fq6.inverse(actual, one)
		if !fq6.equal(actual, one) {
			t.Fatalf("bad inversion 2")
		}
		fq6.double(expected, one)
		fq6.inverse(actual, expected)
		fq6.mul(expected, actual, expected)
		if !fq6.equal(expected, one) {
			t.Fatalf("bad inversion 3")
		}
	})

	t.Run("Exponentiation", func(t *testing.T) {
		fq6.exp(actual, zero, bigZero)
		if !fq6.equal(actual, one) {
			t.Fatalf("bad exponentiation 1")
		}
		fq6.exp(actual, zero, bigOne)
		if !fq6.equal(actual, zero) {
			t.Logf("actual %s\n", fq6.toString(actual))
			t.Fatalf("bad exponentiation 2")
		}
		fq6.exp(actual, one, bigZero)
		if !fq6.equal(actual, one) {
			t.Fatalf("bad exponentiation 3")
		}
		fq6.exp(actual, one, bigOne)
		if !fq6.equal(actual, one) {
			t.Fatalf("bad exponentiation 4")
		}
		fq6.double(expected, one)
		fq6.exp(actual, expected, big.NewInt(2))
		fq6.square(expected, expected)
		if !fq6.equal(expected, actual) {
			t.Fatalf("bad exponentiation 4")
		}
	})
}

func TestFq3(t *testing.T) {
	modulusBytes := bytes_(40, "0x3bcf7bcd473a266249da7b0548ecaeec9635cf44194fb494c07925d6ad3bb4334a400000001")
	f := newField(modulusBytes)
	fq3, err := newFq3(f, nil)
	if err != nil {
		panic(err)
	}
	nonResidue, err := f.newFieldElementFromBytes(bytes_(40, "0x05"))
	if err != nil {
		panic(err)
	}
	f.neg(fq3.nonResidue, nonResidue)
	fq3.calculateFrobeniusCoeffs()

	zero := fq3.zero()
	one := fq3.one()
	actual := fq3.newElement()
	expected := fq3.newElement()

	t.Run("FromBytes & ToBytes", func(t *testing.T) {
		a, err := fq3.fromBytes(fq3One)
		if err != nil {
			t.Fatal(err)
		}
		if !fq3.equal(a, fq3.one()) {
			t.Fatalf("bad fromBytes")
		}
		b, err := fq3.fromBytes(
			fq3.toBytes(a),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !fq3.equal(a, b) {
			t.Fatalf("not equal")
		}
	})

	t.Run("Addition", func(t *testing.T) {
		fq3.add(actual, zero, zero)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad add")
		}
		fq3.add(actual, one, zero)
		if !fq3.equal(actual, one) {
			t.Fatalf("bad add")
		}
		fq3.add(actual, zero, zero)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad add")
		}
	})

	t.Run("Substraction", func(t *testing.T) {
		fq3.sub(actual, zero, zero)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad substraction 1")
		}
		fq3.sub(actual, one, zero)
		if !fq3.equal(actual, one) {
			t.Fatalf("bad substraction 2")
		}
		fq3.sub(actual, one, one)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad substraction 3")
		}
	})

	t.Run("Negation", func(t *testing.T) {
		fq3.sub(expected, zero, one)
		fq3.neg(actual, one)
		if !fq3.equal(expected, actual) {
			t.Fatalf("bad negation")
		}
	})

	t.Run("Multiplication", func(t *testing.T) {
		fq3.mul(actual, zero, zero)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad multiplication 1")
		}
		fq3.mul(actual, one, zero)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq3.mul(actual, zero, one)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq3.mul(actual, one, one)
		if !fq3.equal(actual, one) {
			t.Fatalf("bad multiplication 2")
		}
	})

	t.Run("Squaring", func(t *testing.T) {
		fq3.square(actual, zero)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad squaring 1")
		}
		fq3.square(actual, one)
		if !fq3.equal(actual, one) {
			t.Fatalf("bad squaring 2")
		}
		fq3.double(expected, one)
		fq3.square(actual, expected)
		fq3.mul(expected, expected, expected)
		if !fq3.equal(expected, actual) {
			t.Fatalf("bad squaring 3")
		}
	})

	t.Run("Inverse", func(t *testing.T) {
		fq3.inverse(actual, zero)
		if !fq3.equal(actual, zero) {
			t.Fatalf("bad inversion 1")
		}
		fq3.inverse(actual, one)
		if !fq3.equal(actual, one) {
			t.Fatalf("bad inversion 2")
		}
		fq3.double(expected, one)
		fq3.inverse(actual, expected)
		fq3.mul(expected, actual, expected)
		if !fq3.equal(expected, one) {
			t.Fatalf("bad inversion 3")
		}
	})

	t.Run("Exponentiation", func(t *testing.T) {
		fq3.exp(actual, zero, bigZero)
		if !fq3.equal(actual, one) {
			t.Fatalf("bad exponentiation 1")
		}
		fq3.exp(actual, zero, bigOne)
		if !fq3.equal(actual, zero) {
			t.Logf("actual %s\n", fq3.toString(actual))
			t.Fatalf("bad exponentiation 2")
		}
		fq3.exp(actual, one, bigZero)
		if !fq3.equal(actual, one) {
			t.Fatalf("bad exponentiation 3")
		}
		fq3.exp(actual, one, bigOne)
		if !fq3.equal(actual, one) {
			t.Fatalf("bad exponentiation 4")
		}
		fq3.double(expected, one)
		fq3.exp(actual, expected, big.NewInt(2))
		fq3.square(expected, expected)
		if !fq3.equal(expected, actual) {
			t.Fatalf("bad exponentiation 4")
		}
	})
}

func TestFq12(t *testing.T) {
	modulusBytes := bytes_(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	f := newField(modulusBytes)
	fq2, err := newFq2(f, nil)
	if err != nil {
		panic(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		panic(err)
	}
	fq6.calculateFrobeniusCoeffs()
	f.cpy(fq6.nonResidue[0], f.one)
	f.cpy(fq6.nonResidue[1], f.one)

	fq12, err := newFq12(fq6, nil)
	if err != nil {
		panic(err)
	}
	// f.neg(fq12.nonResidue, f.one)
	fq12.calculateFrobeniusCoeffs()

	zero := fq12.zero()
	one := fq12.one()
	actual := fq12.newElement()
	expected := fq12.newElement()

	t.Run("FromBytes & ToBytes", func(t *testing.T) {
		a, err := fq12.fromBytes(fq12One)
		if err != nil {
			t.Fatal(err)
		}
		if !fq12.equal(a, fq12.one()) {
			t.Fatalf("bad fromBytes")
		}
		b, err := fq12.fromBytes(
			fq12.toBytes(a),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !fq12.equal(a, b) {
			t.Fatalf("not equal")
		}
	})

	t.Run("Addition", func(t *testing.T) {
		fq12.add(actual, zero, zero)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad add")
		}
		fq12.add(actual, one, zero)
		if !fq12.equal(actual, one) {
			t.Fatalf("bad add")
		}
		fq12.add(actual, zero, zero)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad add")
		}
	})
	t.Run("Substraction", func(t *testing.T) {
		fq12.sub(actual, zero, zero)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad substraction 1")
		}
		fq12.sub(actual, one, zero)
		if !fq12.equal(actual, one) {
			t.Fatalf("bad substraction 2")
		}
		fq12.sub(actual, one, one)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad substraction 3")
		}
	})

	t.Run("Negation", func(t *testing.T) {
		fq12.sub(expected, zero, one)
		fq12.neg(actual, one)
		if !fq12.equal(expected, actual) {
			t.Fatalf("bad negation")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		fq12.mul(actual, zero, zero)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad multiplication 1")
		}
		fq12.mul(actual, one, zero)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq12.mul(actual, zero, one)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq12.mul(actual, one, one)
		if !fq12.equal(actual, one) {
			t.Fatalf("bad multiplication 2")
		}
	})

	t.Run("Squaring", func(t *testing.T) {
		fq12.square(actual, zero)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad squaring 1")
		}
		fq12.square(actual, one)
		if !fq12.equal(actual, one) {
			t.Fatalf("bad squaring 2")
		}
		fq12.double(expected, one)
		fq12.square(actual, expected)
		fq12.mul(expected, expected, expected)
		if !fq12.equal(expected, actual) {
			t.Fatalf("bad squaring 3")
		}
	})

	t.Run("Inverse", func(t *testing.T) {
		fq12.inverse(actual, zero)
		if !fq12.equal(actual, zero) {
			t.Fatalf("bad inversion 1")
		}
		fq12.inverse(actual, one)
		if !fq12.equal(actual, one) {
			t.Fatalf("bad inversion 2")
		}
		fq12.double(expected, one)
		fq12.inverse(actual, expected)
		fq12.mul(expected, actual, expected)
		if !fq12.equal(expected, one) {
			t.Fatalf("bad inversion 3")
		}
	})

	t.Run("Exponentiation", func(t *testing.T) {
		fq12.exp(actual, zero, bigZero)
		if !fq12.equal(actual, one) {
			t.Fatalf("bad exponentiation 1")
		}
		fq12.exp(actual, zero, bigOne)
		if !fq12.equal(actual, zero) {
			t.Logf("actual %s\n", fq12.toString(actual))
			t.Fatalf("bad exponentiation 2")
		}
		fq12.exp(actual, one, bigZero)
		if !fq12.equal(actual, one) {
			t.Fatalf("bad exponentiation 3")
		}
		fq12.exp(actual, one, bigOne)
		if !fq12.equal(actual, one) {
			t.Fatalf("bad exponentiation 4")
		}
		fq12.double(expected, one)
		fq12.exp(actual, expected, big.NewInt(2))
		fq12.square(expected, expected)
		if !fq12.equal(expected, actual) {
			t.Fatalf("bad exponentiation 4")
		}
	})
}

func TestBLS1238Pairing(t *testing.T) {
	modulusBytes := bytes_(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	groupBytes := bytes_(48, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	f := newField(modulusBytes)

	// G1
	a, err := f.newFieldElementFromBytes(bytes_(48, "0x00"))
	if err != nil {
		t.Fatal(err)
	}

	b, err := f.newFieldElementFromBytes(bytes_(48, "0x04"))
	if err != nil {
		t.Fatal(err)
	}

	g1, err := newG1(f, nil, nil, groupBytes)
	if err != nil {
		panic(err)
	}
	f.cpy(g1.a, a)
	f.cpy(g1.b, b)

	fq2, err := newFq2(f, nil)
	if err != nil {
		panic(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	// G2
	g2, err := newG22(fq2, nil, nil, groupBytes)
	if err != nil {
		panic(err)
	}
	a2, b2 := fq2.zero(), fq2.newElement()
	f.cpy(b2[0], b)
	f.cpy(b2[1], b)
	fq2.copy(g2.a, a2)
	fq2.copy(g2.b, b2)

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		panic(err)
	}
	f.cpy(fq6.nonResidue[0], f.one)
	f.cpy(fq6.nonResidue[1], f.one)
	fq6.calculateFrobeniusCoeffs()

	fq12, err := newFq12(fq6, nil)
	if err != nil {
		panic(err)
	}
	fq12.calculateFrobeniusCoeffs()

	z, ok := new(big.Int).SetString("d201000000010000", 16)
	if !ok {
		panic("invalid exponent")
	}

	bls := newBLSInstance(z, true, 1, g1, g2, fq12)

	bytesLen := 48
	generatorBytes := bytes_(bytesLen,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	)
	g1One, err := bls.g1.fromBytes(generatorBytes)
	if err != nil {
		panic(err)
	}
	if !bls.g1.isOnCurve(g1One) {
		panic("p is not on curve\n")
	}
	generatorBytes = bytes_(bytesLen,
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
	)
	g2One, err := bls.g2.fromBytes(generatorBytes)
	if err != nil {
		panic(err)
	}
	if !bls.g2.isOnCurve(g2One) {
		panic("q is not on curve\n")
	}
	expectedBytes := bytes_(bytesLen,
		"0x1250ebd871fc0a92a7b2d83168d0d727272d441befa15c503dd8e90ce98db3e7b6d194f60839c508a84305aaca1789b6",
		"0x089a1c5b46e5110b86750ec6a532348868a84045483c92b7af5af689452eafabf1a8943e50439f1d59882a98eaa0170f",
		"0x1368bb445c7c2d209703f239689ce34c0378a68e72a6b3b216da0e22a5031b54ddff57309396b38c881c4c849ec23e87",
		"0x193502b86edb8857c273fa075a50512937e0794e1e65a7617c90d8bd66065b1fffe51d7a579973b1315021ec3c19934f",
		"0x01b2f522473d171391125ba84dc4007cfbf2f8da752f7c74185203fcca589ac719c34dffbbaad8431dad1c1fb597aaa5",
		"0x018107154f25a764bd3c79937a45b84546da634b8f6be14a8061e55cceba478b23f7dacaa35c8ca78beae9624045b4b6",
		"0x19f26337d205fb469cd6bd15c3d5a04dc88784fbb3d0b2dbdea54d43b2b73f2cbb12d58386a8703e0f948226e47ee89d",
		"0x06fba23eb7c5af0d9f80940ca771b6ffd5857baaf222eb95a7d2809d61bfe02e1bfd1b68ff02f0b8102ae1c2d5d5ab1a",
		"0x11b8b424cd48bf38fcef68083b0b0ec5c81a93b330ee1a677d0d15ff7b984e8978ef48881e32fac91b93b47333e2ba57",
		"0x03350f55a7aefcd3c31b4fcb6ce5771cc6a0e9786ab5973320c806ad360829107ba810c5a09ffdd9be2291a0c25a99a2",
		"0x04c581234d086a9902249b64728ffd21a189e87935a954051c7cdba7b3872629a4fafc05066245cb9108f0242d0fe3ef",
		"0x0f41e58663bf08cf068672cbd01a7ec73baca4d72ca93544deff686bfd6df543d48eaa24afe47e1efde449383b676631",
	)
	expected, err := bls.fq12.fromBytes(expectedBytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Expected", func(t *testing.T) {
		actual := bls.pair(g1One, g2One)
		if !bls.fq12.equal(expected, actual) {
			t.Fatalf("bad pairing-1")
		}
	})

	t.Run("Bilinearity", func(t *testing.T) {
		a, _ := rand.Int(rand.Reader, big.NewInt(100))
		b, _ := rand.Int(rand.Reader, big.NewInt(100))
		c := new(big.Int).Mul(a, b)
		G, H := bls.g1.newPoint(), bls.g2.newPoint()
		bls.g1.mulScalar(G, g1One, a)
		bls.g2.mulScalar(H, g2One, b)
		if !bls.g1.isOnCurve(G) {
			t.Fatal("G isnt on the curve")
		}
		if !bls.g2.isOnCurve(H) {
			t.Fatal("H isnt on the curve")
		}

		var f1, f2 *fe12
		// e(a*G1, b*G2) = e(G1, G2)^c
		t.Run("First", func(t *testing.T) {
			bls.g1.affine(G, G)
			bls.g2.affine(H, H)
			f1 = bls.pair(G, H)
			f2 = bls.pair(g1One, g2One)
			bls.fq12.exp(f2, f2, c)
			if !bls.fq12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(c*G1, G2)
		t.Run("Second", func(t *testing.T) {
			G = bls.g1.mulScalar(G, g1One, c)
			bls.g1.affine(G, G)
			f2 = bls.pair(G, g2One)
			if !bls.fq12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(G1, c*G2)
		t.Run("Third", func(t *testing.T) {
			H = bls.g2.mulScalar(H, g2One, c)
			bls.g2.affine(H, H)
			f2 = bls.pair(g1One, H)
			if !bls.fq12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
	})

}

func BenchmarkBLS(t *testing.B) {
	modulusBytes := bytes_(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	groupBytes := bytes_(48, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	f := newField(modulusBytes)
	// G1
	a, err := f.newFieldElementFromBytes(bytes_(48, "0x00"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := f.newFieldElementFromBytes(bytes_(48, "0x04"))
	if err != nil {
		t.Fatal(err)
	}
	g1, err := newG1(f, nil, nil, groupBytes)
	if err != nil {
		panic(err)
	}
	f.cpy(g1.a, a)
	f.cpy(g1.b, b)

	fq2, err := newFq2(f, nil)
	if err != nil {
		panic(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()
	// G2
	g2, err := newG22(fq2, nil, nil, groupBytes)
	if err != nil {
		panic(err)
	}
	a2, b2 := fq2.zero(), fq2.newElement()
	f.cpy(b2[0], b)
	f.cpy(b2[1], b)
	fq2.copy(g2.a, a2)
	fq2.copy(g2.b, b2)

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		panic(err)
	}
	f.cpy(fq6.nonResidue[0], f.one)
	f.cpy(fq6.nonResidue[1], f.one)
	fq6.calculateFrobeniusCoeffs()

	fq12, err := newFq12(fq6, nil)
	if err != nil {
		panic(err)
	}
	fq12.calculateFrobeniusCoeffs()

	z, ok := new(big.Int).SetString("d201000000010000", 16)
	if !ok {
		panic("invalid exponent")
	}

	bls := newBLSInstance(z, true, 1, g1, g2, fq12)

	bytesLen := 48
	generatorBytes := bytes_(bytesLen,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	)
	g1One, err := bls.g1.fromBytes(generatorBytes)
	if err != nil {
		panic(err)
	}
	if !bls.g1.isOnCurve(g1One) {
		panic("p is not on curve\n")
	}
	generatorBytes = bytes_(bytesLen,
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
	)
	g2One, err := bls.g2.fromBytes(generatorBytes)
	if err != nil {
		panic(err)
	}
	if !bls.g2.isOnCurve(g2One) {
		panic("q is not on curve\n")
	}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		bls.pair(g1One, g2One)
	}
}
