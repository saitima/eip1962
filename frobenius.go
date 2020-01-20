package eip

import (
	"errors"
	"math/big"
)

func constructBaseForFq6AndFq12(fq *fq2, nonResidue *fe2) (*fe2, *fe2, error) {
	one, six, rem := big.NewInt(1), big.NewInt(6), big.NewInt(0)
	modulus := fq.f.pbig

	// u^(q-1/6)
	power := new(big.Int).Sub(modulus, one)
	power, rem = new(big.Int).DivMod(power, six, rem)
	if rem.Uint64() != 0 {
		return nil, nil, errors.New("remaining is not zero")
	}
	f1, f2 := fq.newElement(), fq.newElement()
	fq.exp(f1, nonResidue, power)

	// u^(q^2-1/6)
	power = new(big.Int).Mul(modulus, modulus)
	power = new(big.Int).Sub(power, one)
	power, rem = new(big.Int).DivMod(power, six, rem)
	if rem.Uint64() != 0 {
		return nil, nil, errors.New("remaining is not zero")
	}
	fq.exp(f2, nonResidue, power)

	return f1, f2, nil
}

func constructBaseForFq3AndFq6(f *field, nonResidue fieldElement) (fieldElement, error) {
	one, six, rem := big.NewInt(1), big.NewInt(6), big.NewInt(0)
	modulus := f.pbig

	// u^(q-1/6)
	power := new(big.Int).Sub(modulus, one)
	power, rem = new(big.Int).DivMod(power, six, rem)
	if rem.Uint64() != 0 {
		return nil, errors.New("remaining is not zero")
	}
	f1 := f.newFieldElement()
	f.exp(f1, nonResidue, power)

	return f1, nil
}

func constructBaseForFq2AndFq4(f *field, nonResidue fieldElement) fieldElement {
	modulus := f.pbig
	power := new(big.Int).Rsh(modulus, 2)
	f1 := f.newFieldElement()
	f.exp(f1, nonResidue, power)
	return f1
}
