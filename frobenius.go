package eip

import (
	"errors"
	"math/big"
)

func constructBaseForFq6AndFq12(fq6 *fq6C) (*fe2, *fe2, error) {
	fq2 := fq6.fq2()
	zero, one, six, rem := big.NewInt(0), big.NewInt(1), big.NewInt(6), big.NewInt(0)
	modulus := fq2.modulus()
	power := new(big.Int).Sub(modulus, one)
	power, rem = new(big.Int).DivMod(power, six, rem)
	if rem.Cmp(zero) != 0 {
		return nil, nil, errors.New("remaining is not zero")
	}
	f1, f2 := fq2.new(), fq2.new()
	fq2.exp(f1, fq6.nonResidue, power)
	power = new(big.Int).Mul(modulus, modulus)
	power = new(big.Int).Sub(power, one)
	power, rem = new(big.Int).DivMod(power, six, rem)
	if rem.Cmp(zero) != 0 {
		return nil, nil, errors.New("remaining is not zero")
	}
	fq2.exp(f2, fq6.nonResidue, power)
	return f1, f2, nil
}

func constructBaseForFq3AndFq6(fq3 *fq3) (fe, error) {
	fq := fq3.fq()
	zero, one, six, rem := big.NewInt(0), big.NewInt(1), big.NewInt(6), big.NewInt(0)
	power := new(big.Int).Sub(fq.modulus(), one)
	power, rem = new(big.Int).DivMod(power, six, rem)
	if rem.Cmp(zero) != 0 {
		return nil, errors.New("remaining is not zero")
	}
	f1 := fq.new()
	fq.exp(f1, fq3.nonResidue, power)
	return f1, nil
}

func constructBaseForFq2AndFq4(fq2 *fq2) fe {
	fq := fq2.fq()
	power := new(big.Int).Rsh(fq.modulus(), 2)
	f1 := fq.new()
	fq.exp(f1, fq2.nonResidue, power)
	return f1
}
