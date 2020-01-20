package eip

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
	fq4One = bytes_(40,
		"0x01", "0x00",
		"0x00", "0x00",
	)
	fq6CubicOne = bytes_(48,
		"0x01", "0x00",
		"0x00", "0x00",
		"0x00", "0x00",
	)
	fq6QuadraticOne = bytes_(40,
		"0x01", "0x00", "0x00",
		"0x00", "0x00", "0x00",
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

func TestFq2(t *testing.T) {
	modulusBytes := bytes_(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	f := newField(modulusBytes)
	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
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

func TestFq6Cubic(t *testing.T) {
	modulusBytes := bytes_(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	f := newField(modulusBytes)

	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(fq6.nonResidue[0], f.one)
	f.copy(fq6.nonResidue[1], f.one)
	fq6.calculateFrobeniusCoeffs()

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
		t.Fatal(err)
	}
	nonResidue, err := f.newFieldElementFromBytes(bytes_(40, "0x05"))
	if err != nil {
		t.Fatal(err)
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

func TestFq4(t *testing.T) {
	byteLen := 40
	modulusBytes := bytes_(byteLen, "0x3bcf7bcd473a266249da7b0548ecaeec9635d1330ea41a9e35e51200e12c90cd65a71660001")
	f := newField(modulusBytes)

	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	nonResidue, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x11")) // decimal: 17
	if err != nil {
		t.Fatal(err)
	}
	f.neg(fq2.nonResidue, nonResidue)
	fq2.calculateFrobeniusCoeffs()

	fq4, err := newFq4(fq2, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq4.nonResidue = fq2.zero()
	fq4.f.f.copy(fq4.nonResidue[0], fq2.nonResidue)
	fq4.calculateFrobeniusCoeffs()

	zero := fq4.zero()
	one := fq4.one()
	actual := fq4.newElement()
	expected := fq4.newElement()

	t.Run("FromBytes & ToBytes", func(t *testing.T) {
		a, err := fq4.fromBytes(fq4One)
		if err != nil {
			t.Fatal(err)
		}
		if !fq4.equal(a, fq4.one()) {
			t.Fatalf("bad fromBytes")
		}
		b, err := fq4.fromBytes(
			fq4.toBytes(a),
		)
		if err != nil {
			t.Fatal(err)
		}
		if !fq4.equal(a, b) {
			t.Fatalf("not equal")
		}
	})

	t.Run("Addition", func(t *testing.T) {
		fq4.add(actual, zero, zero)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad add")
		}
		fq4.add(actual, one, zero)
		if !fq4.equal(actual, one) {
			t.Fatalf("bad add")
		}
		fq4.add(actual, zero, zero)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad add")
		}
	})
	t.Run("Substraction", func(t *testing.T) {
		fq4.sub(actual, zero, zero)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad substraction 1")
		}
		fq4.sub(actual, one, zero)
		if !fq4.equal(actual, one) {
			t.Fatalf("bad substraction 2")
		}
		fq4.sub(actual, one, one)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad substraction 3")
		}
	})

	t.Run("Negation", func(t *testing.T) {
		fq4.sub(expected, zero, one)
		fq4.neg(actual, one)
		if !fq4.equal(expected, actual) {
			t.Fatalf("bad negation")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		fq4.mul(actual, zero, zero)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad multiplication 1")
		}
		fq4.mul(actual, one, zero)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq4.mul(actual, zero, one)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad multiplication 2")
		}
		fq4.mul(actual, one, one)
		if !fq4.equal(actual, one) {
			t.Fatalf("bad multiplication 2")
		}
	})

	t.Run("Squaring", func(t *testing.T) {
		fq4.square(actual, zero)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad squaring 1")
		}
		fq4.square(actual, one)
		if !fq4.equal(actual, one) {
			t.Fatalf("bad squaring 2")
		}
		fq4.double(expected, one)
		fq4.square(actual, expected)
		fq4.mul(expected, expected, expected)
		if !fq4.equal(expected, actual) {
			t.Fatalf("bad squaring 3")
		}
	})

	t.Run("Inverse", func(t *testing.T) {
		fq4.inverse(actual, zero)
		if !fq4.equal(actual, zero) {
			t.Fatalf("bad inversion 1")
		}
		fq4.inverse(actual, one)
		if !fq4.equal(actual, one) {
			t.Fatalf("bad inversion 2")
		}
		fq4.double(expected, one)
		fq4.inverse(actual, expected)
		fq4.mul(expected, actual, expected)
		if !fq4.equal(expected, one) {
			t.Fatalf("bad inversion 3")
		}
	})

	t.Run("Exponentiation", func(t *testing.T) {
		fq4.exp(actual, zero, bigZero)
		if !fq4.equal(actual, one) {
			t.Fatalf("bad exponentiation 1")
		}
		fq4.exp(actual, zero, bigOne)
		if !fq4.equal(actual, zero) {
			t.Logf("actual %s\n", fq4.toString(actual))
			t.Fatalf("bad exponentiation 2")
		}
		fq4.exp(actual, one, bigZero)
		if !fq4.equal(actual, one) {
			t.Fatalf("bad exponentiation 3")
		}
		fq4.exp(actual, one, bigOne)
		if !fq4.equal(actual, one) {
			t.Fatalf("bad exponentiation 4")
		}
		fq4.double(expected, one)
		fq4.exp(actual, expected, big.NewInt(2))
		fq4.square(expected, expected)
		if !fq4.equal(expected, actual) {
			t.Fatalf("bad exponentiation 4")
		}
	})
}

func TestFq6Quadratic(t *testing.T) {
	byteLen := 40
	modulusBytes := bytes_(byteLen, "0x3bcf7bcd473a266249da7b0548ecaeec9635cf44194fb494c07925d6ad3bb4334a400000001")
	f := newField(modulusBytes)

	fq3, err := newFq3(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	nonResidue, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x05"))
	if err != nil {
		t.Fatal(err)
	}
	f.neg(fq3.nonResidue, nonResidue)
	fq3.calculateFrobeniusCoeffs()

	fq6, err := newFq6Quadratic(fq3, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq6.nonResidue = fq3.zero()
	fq6.f.f.copy(fq6.nonResidue[0], fq3.nonResidue)
	fq6.calculateFrobeniusCoeffs()

	zero := fq6.zero()
	one := fq6.one()
	actual := fq6.newElement()
	expected := fq6.newElement()

	t.Run("FromBytes & ToBytes", func(t *testing.T) {
		a, err := fq6.fromBytes(fq6QuadraticOne)
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

func TestFq12(t *testing.T) {
	modulusBytes := bytes_(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	f := newField(modulusBytes)
	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(fq6.nonResidue[0], f.one)
	f.copy(fq6.nonResidue[1], f.one)
	fq6.calculateFrobeniusCoeffs()

	fq12, err := newFq12(fq6, nil)
	if err != nil {
		t.Fatal(err)
	}
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

func TestG1(t *testing.T) {
	// base field
	byteLen := 48
	modulusBytes := bytes_(byteLen, "1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	groupBytes := bytes_(byteLen, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")

	f := newField(modulusBytes)
	a := bytes_(byteLen, "0x00")
	b := bytes_(byteLen, "0x04")
	g, err := newG1(f, a, b, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	zero := g.zero()
	oneBytes := bytes_(byteLen,
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
			t.Fatal("bad equality 3")
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
	t.Run("Wnaf Multiplication", func(t *testing.T) {
		g.wnafMul(actual, zero, bigZero)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 1")
		}
		g.wnafMul(actual, zero, bigOne)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 2")
		}
		g.wnafMul(actual, one, bigZero)
		if !g.equal(actual, zero) {
			t.Fatal("bad scalar multiplication 3")
		}
		g.wnafMul(actual, one, bigOne)
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

func TestG22(t *testing.T) {
	byteLen := 48
	modulusBytes := bytes_(byteLen, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	groupBytes := bytes_(byteLen, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")

	f := newField(modulusBytes)
	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	g, err := newG22(fq2, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	zero := g.zero()
	oneBytes := bytes_(byteLen,
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
	f.copy(b2[0], b)
	f.copy(b2[1], b)
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
			t.Fatal("bad equality 3")
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

func TestG23(t *testing.T) {
	byteLen := 40
	modulusBytes := bytes_(byteLen, "0x3bcf7bcd473a266249da7b0548ecaeec9635cf44194fb494c07925d6ad3bb4334a400000001")
	groupBytes := bytes_(byteLen, "0x3bcf7bcd473a266249da7b0548ecaeec9635d1330ea41a9e35e51200e12c90cd65a71660001")

	f := newField(modulusBytes)

	fq3, err := newFq3(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq3.nonResidue, err = f.newFieldElementFromBytes(bytes_(byteLen, "0x05"))
	if err != nil {
		t.Fatal(err)
	}
	fq3.calculateFrobeniusCoeffs()

	g, err := newG23(fq3, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}

	zero := g.zero()
	oneBytes := bytes_(byteLen,
		"0x34f7320a12b56ce532bccb3b44902cbaa723cd60035ada7404b743ad2e644ad76257e4c6813",
		"0xcf41620baa52eec50e61a70ab5b45f681952e0109340fec84f1b2890aba9b15cac5a0c80fa",
		"0x11f99170e10e326433cccb8032fb48007ca3c4e105cf31b056ac767e2cb01258391bd4917ce",
		"0x3a65968f03cc64d62ad05c79c415e07ebd38b363ec48309487c0b83e1717a582c1b60fecc91",
		"0xca5e8427e5db1506c1a24cefc2451ab3accaea5db82dcb0c7117cc74402faa5b2c37685c6e",
		"0xf75d2dd88302c9a4ef941307629a1b3e197277d83abb715f647c2e55a27baf782f5c60e7f7",
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
	a, err := f.newFieldElementFromBytes(bytes_(byteLen, "0xb"))
	if err != nil {
		t.Fatal(err)
	}
	b, err := f.newFieldElementFromBytes(bytes_(byteLen, "0xd68c7b1dc5dd042e957b71c44d3d6c24e683fc09b420b1a2d263fde47ddba59463d0c65282"))
	if err != nil {
		t.Fatal(err)
	}

	twist, twist2, twist3 := fq3.newElement(), fq3.newElement(), fq3.newElement()
	f.copy(twist[0], f.zero)
	f.copy(twist[1], f.one)
	fq3.square(twist2, twist)
	fq3.mul(twist3, twist2, twist)
	fq3.mulByFq(g.a, twist2, a)
	fq3.mulByFq(g.b, twist3, b)

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
			t.Fatal("bad equality 3")
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

	t.Run("Multiplication", func(t *testing.T) {
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
		bases := make([]*pointG23, count)
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

func TestBLS12384Pairing(t *testing.T) {
	byteLen := 48
	modulusBytes := bytes_(byteLen, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	groupBytes := bytes_(byteLen, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001")
	f := newField(modulusBytes)

	// G1
	a, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x00"))
	if err != nil {
		t.Fatal(err)
	}

	b, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x04"))
	if err != nil {
		t.Fatal(err)
	}

	g1, err := newG1(f, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(g1.a, a)
	f.copy(g1.b, b)

	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	// G2
	g2, err := newG22(fq2, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	a2, b2 := fq2.zero(), fq2.newElement()
	f.copy(b2[0], b)
	f.copy(b2[1], b)
	fq2.copy(g2.a, a2)
	fq2.copy(g2.b, b2)

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(fq6.nonResidue[0], f.one)
	f.copy(fq6.nonResidue[1], f.one)
	fq6.calculateFrobeniusCoeffs()

	fq12, err := newFq12(fq6, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq12.calculateFrobeniusCoeffs()

	z, ok := new(big.Int).SetString("d201000000010000", 16)
	if !ok {
		t.Fatal("invalid exponent")
	}

	bls := newBLSInstance(z, true, 1, g1, g2, fq12, true)

	generatorBytes := bytes_(byteLen,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	)
	g1One, err := bls.g1.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !bls.g1.isOnCurve(g1One) {
		t.Fatal("p is not on curve\n")
	}
	generatorBytes = bytes_(byteLen,
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
	)
	g2One, err := bls.g2.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !bls.g2.isOnCurve(g2One) {
		t.Fatal("q is not on curve\n")
	}
	expectedBytes := bytes_(byteLen,
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
		t.Fatal(err)
	}
	f.copy(g1.a, a)
	f.copy(g1.b, b)

	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()
	// G2
	g2, err := newG22(fq2, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	a2, b2 := fq2.zero(), fq2.newElement()
	f.copy(b2[0], b)
	f.copy(b2[1], b)
	fq2.copy(g2.a, a2)
	fq2.copy(g2.b, b2)

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(fq6.nonResidue[0], f.one)
	f.copy(fq6.nonResidue[1], f.one)
	fq6.calculateFrobeniusCoeffs()

	fq12, err := newFq12(fq6, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq12.calculateFrobeniusCoeffs()

	z, ok := new(big.Int).SetString("d201000000010000", 16)
	if !ok {
		t.Fatal("invalid exponent")
	}

	bls := newBLSInstance(z, true, 1, g1, g2, fq12, true)

	bytesLen := 48
	generatorBytes := bytes_(bytesLen,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	)
	g1One, err := bls.g1.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !bls.g1.isOnCurve(g1One) {
		t.Fatal("p is not on curve\n")
	}
	generatorBytes = bytes_(bytesLen,
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
	)
	g2One, err := bls.g2.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !bls.g2.isOnCurve(g2One) {
		t.Fatal("q is not on curve\n")
	}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		bls.pair(g1One, g2One)
	}
}

func TestMNT4320Pairing(t *testing.T) {
	byteLen := 40
	modulusBytes := bytes_(byteLen, "0x3bcf7bcd473a266249da7b0548ecaeec9635d1330ea41a9e35e51200e12c90cd65a71660001")
	groupBytes := bytes_(byteLen, "0x3bcf7bcd473a266249da7b0548ecaeec9635cf44194fb494c07925d6ad3bb4334a400000001")
	f := newField(modulusBytes)

	// G1
	a, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x02"))
	if err != nil {
		t.Fatal(err)
	}

	b, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x03545a27639415585ea4d523234fc3edd2a2070a085c7b980f4e9cd21a515d4b0ef528ec0fd5"))
	if err != nil {
		t.Fatal(err)
	}

	g1, err := newG1(f, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(g1.a, a)
	f.copy(g1.b, b)

	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}

	nonResidue, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x011")) // decimal 17
	if err != nil {
		t.Fatal(err)
	}
	f.copy(fq2.nonResidue, nonResidue)
	fq2.calculateFrobeniusCoeffs()

	fq4, err := newFq4(fq2, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq4.nonResidue = fq2.zero()
	fq4.f.f.copy(fq4.nonResidue[0], fq2.nonResidue)
	fq4.calculateFrobeniusCoeffs()

	// G2
	g2, err := newG22(fq2, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	// y^2 = x^3 + b/(9+u)
	twist, twist2, twist3 := fq2.newElement(), fq2.newElement(), fq2.newElement()
	f.copy(twist[0], f.zero)
	f.copy(twist[1], f.one)
	fq2.square(twist2, twist)
	fq2.mul(twist3, twist2, twist)

	a2, b2 := fq2.newElement(), fq2.newElement()
	fq2.mulByFq(a2, twist2, a)
	fq2.mulByFq(b2, twist3, b)
	fq2.copy(g2.a, a2)
	fq2.copy(g2.b, b2)

	// mnt4 instance)
	z, ok := new(big.Int).SetString("689871209842287392837045615510547309923794944", 10)
	if !ok {
		t.Fatal("invalid value")
	}

	expW0, ok := new(big.Int).SetString("689871209842287392837045615510547309923794945", 10)
	if !ok {
		t.Fatal("invalid expW0")
	}
	expW1 := big.NewInt(1)

	mnt4 := newMnt4Instance(z, false, expW0, expW1, false, fq4, g1, g2, twist)
	generatorBytes := bytes_(byteLen,
		"0x7a2caf82a1ba85213fe6ca3875aee86aba8f73d69060c4079492b948dea216b5b9c8d2af46",
		"0x2db619461cc82672f7f159fec2e89d0148dcc9862d36778c1afd96a71e29cba48e710a48ab2",
	)
	g1One, err := mnt4.g1.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !mnt4.g1.isOnCurve(g1One) {
		t.Fatal("p is not on curve\n")
	}
	generatorBytes = bytes_(byteLen,
		"0x371780491c5660571ff542f2ef89001f205151e12a72cb14f01a931e72dba7903df6c09a9a4",
		"0x4ba59a3f72da165def838081af697c851f002f576303302bb6c02c712c968be32c0ae0a989",
		"0x4b471f33ffaad868a1c47d6605d31e5c4b3b2e0b60ec98f0f610a5aafd0d9522bca4e79f22",
		"0x355d05a1c69a5031f3f81a5c100cb7d982f78ec9cfc3b5168ed8d75c7c484fb61a3cbf0e0f1",
	)

	g2One, err := mnt4.g2.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !mnt4.g2.isOnCurve(g2One) {
		t.Fatal("q is not on curve\n")
	}
	expectedBytes := bytes_(byteLen,
		"0x000003653498b90d54a52c420cc4a73ad1882feb23bf2ae451037a96e17babd70402dd237238b101",
		"0x000000cd0a4994729a71440144fedc4378511a7febdf4cdb0499253bcbea9e023c6cfa9cf9682784",
		"0x000002532341e5b711a9f8f7049a99af28177e51d7a0c384d19cb7547352a7e65c44417babfe0089",
		"0x0000030c15f867b221786f818a8e96ffa041ea4366fee9bdc9b2d845d6a9b9aded3d34d24b1a34b8",
	)
	expected, err := mnt4.fq4.fromBytes(expectedBytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Expected", func(t *testing.T) {
		actual := mnt4.pair(g1One, g2One)
		if !mnt4.fq4.equal(expected, actual) {
			t.Logf("\nexpected: %s\b", mnt4.fq4.toString(expected))
			t.Logf("\actual: %s\b", mnt4.fq4.toString(actual))
			t.Fatalf("bad pairing-1")
		}
	})

	t.Run("Bilinearity", func(t *testing.T) {
		a, _ := rand.Int(rand.Reader, big.NewInt(100))
		b, _ := rand.Int(rand.Reader, big.NewInt(100))
		c := new(big.Int).Mul(a, b)
		G, H := mnt4.g1.newPoint(), mnt4.g2.newPoint()
		mnt4.g1.mulScalar(G, g1One, a)
		mnt4.g2.mulScalar(H, g2One, b)
		if !mnt4.g1.isOnCurve(G) {
			t.Fatal("G isnt on the curve")
		}
		if !mnt4.g2.isOnCurve(H) {
			t.Fatal("H isnt on the curve")
		}

		var f1, f2 *fe4
		// e(a*G1, b*G2) = e(G1, G2)^c
		t.Run("First", func(t *testing.T) {
			mnt4.g1.affine(G, G)
			mnt4.g2.affine(H, H)
			f1 = mnt4.pair(G, H)
			f2 = mnt4.pair(g1One, g2One)
			mnt4.fq4.exp(f2, f2, c)
			if !mnt4.fq4.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(c*G1, G2)
		t.Run("Second", func(t *testing.T) {
			G = mnt4.g1.mulScalar(G, g1One, c)
			mnt4.g1.affine(G, G)
			f2 = mnt4.pair(G, g2One)
			if !mnt4.fq4.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(G1, c*G2)
		t.Run("Third", func(t *testing.T) {
			H = mnt4.g2.mulScalar(H, g2One, c)
			mnt4.g2.affine(H, H)
			f2 = mnt4.pair(g1One, H)
			if !mnt4.fq4.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
	})

}

func TestMNT4753Pairing(t *testing.T) {
	byteLen := 96
	modulusBytes := bytes_(byteLen, "0x1c4c62d92c41110229022eee2cdadb7f997505b8fafed5eb7e8f96c97d87307fdb925e8a0ed8d99d124d9a15af79db117e776f218059db80f0da5cb537e38685acce9767254a4638810719ac425f0e39d54522cdd119f5e9063de245e8001")
	groupBytes := bytes_(byteLen, "0x1c4c62d92c41110229022eee2cdadb7f997505b8fafed5eb7e8f96c97d87307fdb925e8a0ed8d99d124d9a15af79db26c5c28c859a99b3eebca9429212636b9dff97634993aa4d6c381bc3f0057974ea099170fa13a4fd90776e240000001")
	f := newField(modulusBytes)

	// G1
	a, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x02"))
	if err != nil {
		t.Fatal(err)
	}

	b, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x1373684a8c9dcae7a016ac5d7748d3313cd8e39051c596560835df0c9e50a5b59b882a92c78dc537e51a16703ec9855c77fc3d8bb21c8d68bb8cfb9db4b8c8fba773111c36c8b1b4e8f1ece940ef9eaad265458e06372009c9a0491678ef4"))
	if err != nil {
		t.Fatal(err)
	}

	g1, err := newG1(f, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(g1.a, a)
	f.copy(g1.b, b)

	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}

	nonResidue, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x0d")) // decimal 13
	if err != nil {
		t.Fatal(err)
	}
	f.copy(fq2.nonResidue, nonResidue)
	fq2.calculateFrobeniusCoeffs()

	fq4, err := newFq4(fq2, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq4.nonResidue = fq2.zero()
	// fq4.f.f.copy(fq4.nonResidue[0], fq2.nonResidue)
	fq4.calculateFrobeniusCoeffs()

	// G2
	g2, err := newG22(fq2, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	// y^2 = x^3 + b/(9+u)
	twist, twist2, twist3 := fq2.newElement(), fq2.newElement(), fq2.newElement()
	f.copy(twist[0], f.zero)
	f.copy(twist[1], f.one)
	fq2.square(twist2, twist)
	fq2.mul(twist3, twist2, twist)
	fq2.mulByFq(g2.a, twist2, a)
	fq2.mulByFq(g2.b, twist3, b)

	// mnt4 instance)
	z, ok := new(big.Int).SetString("15474b1d641a3fd86dcbcee5dcda7fe51852c8cbe26e600733b714aa43c31a66b0344c4e2c428b07a7713041ba18000", 16)
	if !ok {
		t.Fatal("bad exp")
	}
	expW0, ok := new(big.Int).SetString("15474b1d641a3fd86dcbcee5dcda7fe51852c8cbe26e600733b714aa43c31a66b0344c4e2c428b07a7713041ba17fff", 16)
	if !ok {
		t.Fatal("bad exp w0")
	}
	expW1 := big.NewInt(1)

	mnt4 := newMnt4Instance(z, true, expW0, expW1, true, fq4, g1, g2, twist)
	generatorBytes := bytes_(byteLen,
		"0x1013b42397c8b004d06f0e98fbc12e8ee65adefcdba683c5630e6b58fb69610b02eab1d43484ddfab28213098b562d799243fb14330903aa64878cfeb34a45d1285da665f5c3f37eb76b86209dcd081ccaef03e65f33d490de480bfee06db",
		"0xe3eb479d308664381e7942d6c522c0833f674296169420f1dd90680d0ba6686fc27549d52e4292ea5d611cb6b0df32545b07f281032d0a71f8d485e6907766462e17e8dd55a875bd36fe4cd42cac31c0629fb26c333fe091211d0561d10e",
	)
	g1One, err := mnt4.g1.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !mnt4.g1.isOnCurve(g1One) {
		t.Fatalf("p is not on curve\n")
	}
	generatorBytes = bytes_(byteLen,
		"0xf1b7155ed4e903332835a5de0f327aa11b2d74eb8627e3a7b833be42c11d044b5cf0ae49850eeb07d90c77c67256474b2febf924aca0bfa2e4dacb821c91a04fd0165ac8debb2fc1e763a5c32c2c9f572caa85a91c5243ec4b2981af8904",
		"0xd49c264ec663e731713182a88907b8e979ced82ca592777ad052ec5f4b95dc78dc2010d74f82b9e6d066813ed67f3af1de0d5d425da7a19916cf103f102adf5f95b6b62c24c7d186d60b4a103e157e5667038bb2e828a3374d6439526272",
		"0x4b0e2fef08096ebbaddd2d7f288c4acf17b2267e21dc5ce0f925cd5d02209e34d8b69cc94aef5d90af34d3cd98287ace8f1162079cd2d3d7e6c6c2c073c24a359437e75638a1458f4b2face11f8d2a5200b14d6f9dd0fdd407f04be620ee",
		"0xbc1925e7fcb64f6f8697cd5e45fae22f5688e51b30bd984c0acdc67d2962520e80d31966e3ec477909ecca358be2eee53c75f55a6f7d9660dd6f3d4336ad50e8bfa5375791d73b863d59c422c3ea006b013e7afb186f2eaa9df68f4d6098",
	)

	g2One, err := mnt4.g2.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !mnt4.g2.isOnCurve(g2One) {
		t.Fatalf("g2 one is not on the curve")
	}
	expectedBytes := bytes_(byteLen,
		"0x000003653498b90d54a52c420cc4a73ad1882feb23bf2ae451037a96e17babd70402dd237238b101",
		"0x000000cd0a4994729a71440144fedc4378511a7febdf4cdb0499253bcbea9e023c6cfa9cf9682784",
		"0x000002532341e5b711a9f8f7049a99af28177e51d7a0c384d19cb7547352a7e65c44417babfe0089",
		"0x0000030c15f867b221786f818a8e96ffa041ea4366fee9bdc9b2d845d6a9b9aded3d34d24b1a34b8",
	)
	expected, err := mnt4.fq4.fromBytes(expectedBytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Expected", func(t *testing.T) {
		actual := mnt4.pair(g1One, g2One)
		if !mnt4.fq4.equal(expected, actual) {
			// t.Logf("\nexpected: %s\b", mnt4.fq4.toString(expected))
			// t.Logf("\actual: %s\b", mnt4.fq4.toString(actual))
			// t.Fatalf("bad pairing-1")
		}
	})

	t.Run("Bilinearity", func(t *testing.T) {
		a, _ := rand.Int(rand.Reader, big.NewInt(100))
		b, _ := rand.Int(rand.Reader, big.NewInt(100))
		c := new(big.Int).Mul(a, b)
		G, H := mnt4.g1.newPoint(), mnt4.g2.newPoint()
		mnt4.g1.mulScalar(G, g1One, a)
		mnt4.g2.mulScalar(H, g2One, b)
		if !mnt4.g1.isOnCurve(G) {
			t.Fatal("G isnt on the curve")
		}
		if !mnt4.g2.isOnCurve(H) {
			t.Fatal("H isnt on the curve")
		}

		var f1, f2 *fe4
		// e(a*G1, b*G2) = e(G1, G2)^c
		t.Run("First", func(t *testing.T) {
			mnt4.g1.affine(G, G)
			mnt4.g2.affine(H, H)
			f1 = mnt4.pair(G, H)
			f2 = mnt4.pair(g1One, g2One)
			mnt4.fq4.exp(f2, f2, c)
			if !mnt4.fq4.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(c*G1, G2)
		t.Run("Second", func(t *testing.T) {
			G = mnt4.g1.mulScalar(G, g1One, c)
			mnt4.g1.affine(G, G)
			f2 = mnt4.pair(G, g2One)
			if !mnt4.fq4.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(G1, c*G2)
		t.Run("Third", func(t *testing.T) {
			H = mnt4.g2.mulScalar(H, g2One, c)
			mnt4.g2.affine(H, H)
			f2 = mnt4.pair(g1One, H)
			if !mnt4.fq4.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
	})

}

func TestBN254Pairing(t *testing.T) {
	byteLen := 32
	modulusBytes := bytes_(byteLen, "0x30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47")
	groupBytes := bytes_(byteLen, "0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001")
	f := newField(modulusBytes)

	// G1
	a, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x00"))
	if err != nil {
		t.Fatal(err)
	}

	b, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x03"))
	if err != nil {
		t.Fatal(err)
	}

	g1, err := newG1(f, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(g1.a, a)
	f.copy(g1.b, b)

	fq2, err := newFq2(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.neg(fq2.nonResidue, f.one)
	fq2.calculateFrobeniusCoeffs()

	fq6, err := newFq6(fq2, nil)
	if err != nil {
		t.Fatal(err)
	}
	nine, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x09"))
	if err != nil {
		t.Fatal(err)
	}
	f.copy(fq6.nonResidue[0], nine)
	f.copy(fq6.nonResidue[1], f.one)
	fq6.calculateFrobeniusCoeffs()

	fq12, err := newFq12(fq6, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq12.calculateFrobeniusCoeffs()

	// G2
	g2, err := newG22(fq2, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	// y^2 = x^3 + b/(9+u)
	a2, b2 := fq2.zero(), fq2.newElement()
	fq2.inverse(b2, fq6.nonResidue)
	fq2.mulByFq(b2, b2, b)
	fq2.copy(g2.a, a2)
	fq2.copy(g2.b, b2)

	minus2Inv := new(big.Int).ModInverse(big.NewInt(-2), f.pbig)
	nonResidueInPMinus1Over2 := fq2.newElement()
	fq2.exp(nonResidueInPMinus1Over2, fq6.nonResidue, minus2Inv)
	u := new(big.Int).SetUint64(4965661367192848881)
	sixUPlus2 := new(big.Int).Mul(u, big.NewInt(6))
	sixUPlus2 = new(big.Int).Add(sixUPlus2, big.NewInt(2))

	bn := newBNInstance(u, sixUPlus2, false, 2, g1, g2, fq12, nonResidueInPMinus1Over2, true)

	generatorBytes := bytes_(byteLen,
		"0x01",
		"0x02",
	)
	g1One, err := bn.g1.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !bn.g1.isOnCurve(g1One) {
		t.Fatal("p is not on curve\n")
	}
	generatorBytes = bytes_(byteLen,
		"0x1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed",
		"0x198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2",
		"0x12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa",
		"0x90689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b",
	)
	g2One, err := bn.g2.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !bn.g2.isOnCurve(g2One) {
		t.Fatal("q is not on curve\n")
	}
	expectedBytes := bytes_(byteLen,
		"0x12c70e90e12b7874510cd1707e8856f71bf7f61d72631e268fca81000db9a1f5",
		"0x084f330485b09e866bc2f2ea2b897394deaf3f12aa31f28cb0552990967d4704",
		"0x0e841c2ac18a4003ac9326b9558380e0bc27fdd375e3605f96b819a358d34bde",
		"0x2067586885c3318eeffa1938c754fe3c60224ee5ae15e66af6b5104c47c8c5d8",
		"0x01676555de427abc409c4a394bc5426886302996919d4bf4bdd02236e14b3636",
		"0x2b03614464f04dd772d86df88674c270ffc8747ea13e72da95e3594468f222c4",
		"0x2c53748bcd21a7c038fb30ddc8ac3bf0af25d7859cfbc12c30c866276c565909",
		"0x27ed208e7a0b55ae6e710bbfbd2fd922669c026360e37cc5b2ab862411536104",
		"0x1ad9db1937fd72f4ac462173d31d3d6117411fa48dba8d499d762b47edb3b54a",
		"0x279db296f9d479292532c7c493d8e0722b6efae42158387564889c79fc038ee3",
		"0x0dc26f240656bbe2029bd441d77c221f0ba4c70c94b29b5f17f0f6d08745a069",
		"0x108c19d15f9446f744d0f110405d3856d6cc3bda6c4d537663729f5257628417",
	)
	expected, err := bn.fq12.fromBytes(expectedBytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Expected", func(t *testing.T) {
		actual := bn.pair(g1One, g2One)
		if !bn.fq12.equal(expected, actual) {
			t.Logf("\nexpected: %s\b", bn.fq12.toString(expected))
			t.Logf("\actual: %s\b", bn.fq12.toString(actual))
			t.Fatalf("bad pairing-1")
		}
	})

	t.Run("Bilinearity", func(t *testing.T) {
		a, _ := rand.Int(rand.Reader, big.NewInt(100))
		b, _ := rand.Int(rand.Reader, big.NewInt(100))
		c := new(big.Int).Mul(a, b)
		G, H := bn.g1.newPoint(), bn.g2.newPoint()
		bn.g1.mulScalar(G, g1One, a)
		bn.g2.mulScalar(H, g2One, b)
		if !bn.g1.isOnCurve(G) {
			t.Fatal("G isnt on the curve")
		}
		if !bn.g2.isOnCurve(H) {
			t.Fatal("H isnt on the curve")
		}

		var f1, f2 *fe12
		// e(a*G1, b*G2) = e(G1, G2)^c
		t.Run("First", func(t *testing.T) {
			bn.g1.affine(G, G)
			bn.g2.affine(H, H)
			f1 = bn.pair(G, H)
			f2 = bn.pair(g1One, g2One)
			bn.fq12.exp(f2, f2, c)
			if !bn.fq12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(c*G1, G2)
		t.Run("Second", func(t *testing.T) {
			G = bn.g1.mulScalar(G, g1One, c)
			bn.g1.affine(G, G)
			f2 = bn.pair(G, g2One)
			if !bn.fq12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(G1, c*G2)
		t.Run("Third", func(t *testing.T) {
			H = bn.g2.mulScalar(H, g2One, c)
			bn.g2.affine(H, H)
			f2 = bn.pair(g1One, H)
			if !bn.fq12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
	})

}

func TestMNT6320Pairing(t *testing.T) {
	byteLen := 40
	modulusBytes := bytes_(byteLen, "0x3bcf7bcd473a266249da7b0548ecaeec9635cf44194fb494c07925d6ad3bb4334a400000001")
	groupBytes := bytes_(byteLen, "0x3bcf7bcd473a266249da7b0548ecaeec9635d1330ea41a9e35e51200e12c90cd65a71660001")
	f := newField(modulusBytes)

	// G1
	a, err := f.newFieldElementFromBytes(bytes_(byteLen, "0xb"))
	if err != nil {
		t.Fatal(err)
	}

	b, err := f.newFieldElementFromBytes(bytes_(byteLen, "0xd68c7b1dc5dd042e957b71c44d3d6c24e683fc09b420b1a2d263fde47ddba59463d0c65282"))
	if err != nil {
		t.Fatal(err)
	}

	g1, err := newG1(f, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}
	f.copy(g1.a, a)
	f.copy(g1.b, b)

	fq3, err := newFq3(f, nil)
	if err != nil {
		t.Fatal(err)
	}
	nonResidue, err := f.newFieldElementFromBytes(bytes_(byteLen, "0x05"))
	if err != nil {
		t.Fatal(err)
	}
	f.copy(fq3.nonResidue, nonResidue)
	fq3.calculateFrobeniusCoeffs()

	fq6, err := newFq6Quadratic(fq3, nil)
	if err != nil {
		t.Fatal(err)
	}
	fq6.nonResidue = fq3.zero()
	fq6.f.f.copy(fq6.nonResidue[0], fq3.nonResidue)
	fq6.calculateFrobeniusCoeffs()

	// G2
	g2, err := newG23(fq3, nil, nil, groupBytes)
	if err != nil {
		t.Fatal(err)
	}

	twist, twist2, twist3 := fq3.newElement(), fq3.newElement(), fq3.newElement()
	f.copy(twist[0], f.zero)
	f.copy(twist[1], f.one)
	fq3.square(twist2, twist)
	fq3.mul(twist3, twist2, twist)
	fq3.mulByFq(g2.a, twist2, a)
	fq3.mulByFq(g2.b, twist3, b)

	// mnt6 instance)
	z, ok := new(big.Int).SetString("1eef5546609756bec2a33f0dc9a1b671660000", 16)
	if !ok {
		t.Fatal("invalid value")
	}

	expW0, ok := new(big.Int).SetString("1eef5546609756bec2a33f0dc9a1b671660000", 16)
	if !ok {
		t.Fatal("invalid expW0")
	}
	expW1 := big.NewInt(1)

	mnt6 := newMNT6Instance(z, true, expW0, expW1, true, fq6, g1, g2, twist)

	generatorBytes := bytes_(byteLen,
		"0x2a4feee24fd2c69d1d90471b2ba61ed56f9bad79b57e0b4c671392584bdadebc01abbc0447d",
		"0x32986c245f6db2f82f4e037bf7afd69cbfcbff07fc25d71e9c75e1b97208a333d73d91d3028",
	)
	g1One, err := mnt6.g1.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !mnt6.g1.isOnCurve(g1One) {
		t.Fatal("p is not on curve\n")
	}
	generatorBytes = bytes_(byteLen,
		"0x34f7320a12b56ce532bccb3b44902cbaa723cd60035ada7404b743ad2e644ad76257e4c6813",
		"0xcf41620baa52eec50e61a70ab5b45f681952e0109340fec84f1b2890aba9b15cac5a0c80fa",
		"0x11f99170e10e326433cccb8032fb48007ca3c4e105cf31b056ac767e2cb01258391bd4917ce",
		"0x3a65968f03cc64d62ad05c79c415e07ebd38b363ec48309487c0b83e1717a582c1b60fecc91",
		"0xca5e8427e5db1506c1a24cefc2451ab3accaea5db82dcb0c7117cc74402faa5b2c37685c6e",
		"0xf75d2dd88302c9a4ef941307629a1b3e197277d83abb715f647c2e55a27baf782f5c60e7f7",
	)
	g2One, err := mnt6.g2.fromBytes(generatorBytes)
	if err != nil {
		t.Fatal(err)
	}
	if !mnt6.g2.isOnCurve(g2One) {
		t.Fatal("q is not on curve\n")
	}
	expectedBytes := bytes_(byteLen,
		"0x0000014ac12149eebffe74a1c75a7225deb91ca243c49eef01392080ff519ab6209431f81b50ec03",
		"0x000001ba8ab5bc93186b5bc2b1936fee360528228ab953fbce3c7b84f71d6c0e87b293d0de36eb93",
		"0x00000323a5728ce32f5a04635ca9f84857882e9c13a2b415a021985921c79f303f1f0b69557c5c3d",
		"0x0000032e067f62de41a786c2a43da960855694f3e0da14a964377a32ddad42cf9dd6b80bdc8d4300",
		"0x000000bf02fd56dcd4f6b1d132c8b56a9f8801696d77cdb911a35335360f07eba30bc3083ecaa394",
		"0x0000028a449b7699751b6bf17003c141307311241614b886c0fb6ffaf5b39896e182bddd85859e9c",
	)
	expected, err := mnt6.fq6.fromBytes(expectedBytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Expected", func(t *testing.T) {
		actual := mnt6.pair(g1One, g2One)
		if !mnt6.fq6.equal(expected, actual) {
			t.Fatalf("bad pairing-1")
		}
	})

	t.Run("Bilinearity", func(t *testing.T) {
		a, _ := rand.Int(rand.Reader, big.NewInt(100))
		b, _ := rand.Int(rand.Reader, big.NewInt(100))
		c := new(big.Int).Mul(a, b)
		G, H := mnt6.g1.newPoint(), mnt6.g2.newPoint()
		mnt6.g1.mulScalar(G, g1One, a)
		mnt6.g2.mulScalar(H, g2One, b)
		if !mnt6.g1.isOnCurve(G) {
			t.Fatal("G isnt on the curve")
		}
		if !mnt6.g2.isOnCurve(H) {
			t.Fatal("H isnt on the curve")
		}

		var f1, f2 *fe6q
		// e(a*G1, b*G2) = e(G1, G2)^c
		t.Run("First", func(t *testing.T) {
			mnt6.g1.affine(G, G)
			mnt6.g2.affine(H, H)
			f1 = mnt6.pair(G, H)
			f2 = mnt6.pair(g1One, g2One)
			mnt6.fq6.exp(f2, f2, c)
			if !mnt6.fq6.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(c*G1, G2)
		t.Run("Second", func(t *testing.T) {
			G = mnt6.g1.mulScalar(G, g1One, c)
			mnt6.g1.affine(G, G)
			f2 = mnt6.pair(G, g2One)
			if !mnt6.fq6.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(G1, c*G2)
		t.Run("Third", func(t *testing.T) {
			H = mnt6.g2.mulScalar(H, g2One, c)
			mnt6.g2.affine(H, H)
			f2 = mnt6.pair(g1One, H)
			if !mnt6.fq6.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
	})

}
