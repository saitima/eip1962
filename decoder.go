package eip

import (
	"errors"
	"math/big"
)

const (
	BYTES_FOR_LENGTH_ENCODING               = 1
	EXTENSION_DEGREE_LENGTH_ENCODING        = 1
	BOOLEAN_ENCODING_LENGTH                 = 1
	EXTENSION_TWO_DEGREE                    = 2
	EXTENSION_THREE_DEGREE                  = 3
	TWIST_TYPE_LENGTH                       = 1
	MAX_MODULUS_BYTE_LEN                    = 128
	MAX_GROUP_BYTE_LEN                      = 128
	MAX_BN_U_BIT_LENGTH                     = 128
	MAX_BLS12_X_BIT_LENGTH                  = 128
	MAX_ATE_PAIRING_ATE_LOOP_COUNT          = 2032
	SIGN_ENCODING_LENGTH                    = 1
	MAX_BN_SIX_U_PLUS_TWO_HAMMING           = 128
	MAX_BLS12_X_HAMMING                     = 128
	MAX_ATE_PAIRING_ATE_LOOP_COUNT_HAMMING  = 2032
	MAX_ATE_PAIRING_FINAL_EXP_W0_BIT_LENGTH = 2032
	MAX_ATE_PAIRING_FINAL_EXP_W1_BIT_LENGTH = 2032
)

func split(in []byte, offset int) ([]byte, []byte, error) {
	if len(in) < offset {
		return nil, nil, errors.New("cant split at given offset")
	}
	return in[:offset], in[offset:], nil
}

func decodeBaseFieldParams(in []byte) ([]byte, int, []byte, error) {
	modulusLenBuf, rest, err := split(in, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return nil, 0, nil, errors.New("Input is not long enough to get modulus length")
	}
	modulusLen := int(modulusLenBuf[0])
	if modulusLen == 0 {
		return nil, 0, nil, errors.New("Modulus length is zero")
	}
	if modulusLen > MAX_MODULUS_BYTE_LEN {
		return nil, 0, nil, errors.New("Encoded modulus length is too large")
	}
	modulusBuf, rest, err := split(rest, modulusLen)
	if err != nil {
		return nil, 0, nil, errors.New("Input is not long enough to get modulus")
	}
	if int(modulusBuf[0]) == 0 {
		return nil, 0, nil, errors.New("In modulus encoding highest byte is zero")
	}
	modulusBuf = padHex(modulusBuf)
	modulus := new(big.Int).SetBytes(modulusBuf)
	if isBigZero(modulus) {
		return nil, 0, nil, errors.New("Modulus can not be zero")
	}
	if isBigEven(modulus) {
		return nil, 0, nil, errors.New("Modulus is even")
	}

	if isBigLowerThan(modulus, 3) {
		return nil, 0, nil, errors.New("Modulus is less than 3")
	}
	return modulusBuf, modulusLen, rest, nil
}

// G1
func parseBaseFieldFromEncoding(in []byte) (*field, *big.Int, int, []byte, error) {
	modulusBuf, modulusLen, rest, err := decodeBaseFieldParams(in)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	// TODO: remove after for fuzz testing
	var field *field
	if USE_4LIMBS_FOR_LOWER_LIMBS && modulusLen <= 32 {
		field, err = newField(padBytes(modulusBuf, 32))
	} else {
		field, err = newField(modulusBuf)
	}
	if err != nil {
		return nil, nil, 0, nil, errors.New("Failed to create prime field from modulus")
	}
	if len(rest) < modulusLen {
		return nil, nil, 0, nil, errors.New("FInput is not long enough")
	}
	modulus := new(big.Int).Set(field.pbig)
	return field, modulus, modulusLen, rest, nil
}

func decodeFp(in []byte, modulusLen int, field *field) (fieldElement, []byte, error) {
	xBuf, rest, err := split(in, modulusLen)
	if err != nil {
		return nil, nil, err
	}
	x, err := field.newFieldElementFromBytes(padHex(xBuf))
	if err != nil {
		return nil, nil, err
	}
	return x, rest, nil
}

func decodeFp2(in []byte, modulusLen int, field *fq2) (*fe2, []byte, error) {
	c0Buf, rest, err := split(in, modulusLen)
	if err != nil {
		return nil, nil, err
	}
	c0, err := field.f.newFieldElementFromBytes(padHex(c0Buf))
	if err != nil {
		return nil, nil, err
	}
	c1Buf, rest, err := split(rest, modulusLen)
	if err != nil {
		return nil, nil, err
	}
	c1, err := field.f.newFieldElementFromBytes(padHex(c1Buf))
	if err != nil {
		return nil, nil, err
	}
	elem := field.newElement()
	field.f.copy(elem[0], c0)
	field.f.copy(elem[1], c1)
	return elem, rest, nil
}

func decodeFp3(in []byte, modulusLen int, field *fq3) (*fe3, []byte, error) {
	c0Buf, rest, err := split(in, modulusLen)
	if err != nil {
		return nil, nil, err
	}
	c0, err := field.f.newFieldElementFromBytes(padHex(c0Buf))
	if err != nil {
		return nil, nil, err
	}
	c1Buf, rest, err := split(rest, modulusLen)
	if err != nil {
		return nil, nil, err
	}
	c1, err := field.f.newFieldElementFromBytes(padHex(c1Buf))
	if err != nil {
		return nil, nil, err
	}
	c2Buf, rest, err := split(rest, modulusLen)
	if err != nil {
		return nil, nil, err
	}
	c2, err := field.f.newFieldElementFromBytes(padHex(c2Buf))
	if err != nil {
		return nil, nil, err
	}
	elem := field.newElement()
	field.f.copy(elem[0], c0)
	field.f.copy(elem[1], c1)
	field.f.copy(elem[2], c2)
	return elem, rest, nil
}

func decodeBAInBaseFieldFromEncoding(in []byte, modulusLen int, field *field) (fieldElement, fieldElement, []byte, error) {
	a, rest, err := decodeFp(in, modulusLen, field)
	if err != nil {
		return nil, nil, nil, err
	}

	b, rest, err := decodeFp(rest, modulusLen, field)
	if err != nil {
		return nil, nil, nil, err
	}

	return a, b, rest, nil
}

func decodeBAInExtField2FromEncoding(in []byte, modulusLen int, field *fq2) (*fe2, *fe2, []byte, error) {
	a2, rest, err := decodeFp2(in, modulusLen, field)
	if err != nil {
		return nil, nil, nil, err
	}
	b2, rest, err := decodeFp2(rest, modulusLen, field)
	if err != nil {
		return nil, nil, nil, err
	}

	return a2, b2, rest, nil
}

func decodeBAInExtField3FromEncoding(in []byte, modulusLen int, field *fq3) (*fe3, *fe3, []byte, error) {
	a3, rest, err := decodeFp3(in, modulusLen, field)
	if err != nil {
		return nil, nil, nil, err
	}
	b3, rest, err := decodeFp3(rest, modulusLen, field)
	if err != nil {
		return nil, nil, nil, err
	}
	return a3, b3, rest, nil
}

func decodeGroupOrder(in []byte) (int, *big.Int, []byte, error) {
	orderLenBuf, rest, err := split(in, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return 0, nil, nil, errors.New("Input is not long enough to get group order length")
	}
	orderLen := int(orderLenBuf[0])
	if orderLen == 0 {
		return 0, nil, nil, errors.New("Encoded group length is zero")
	}
	if orderLen > MAX_GROUP_BYTE_LEN {
		return 0, nil, nil, errors.New("Encoded group length is too large")
	}
	orderBuf, rest, err := split(rest, orderLen)
	if err != nil {
		return 0, nil, nil, errors.New("Input is not long enough to get group order")
	}
	// if orderBuf[0] == 0 {
	// 	return 0, nil, nil, errors.New("Encoded group length has zero top byte")
	// }
	order := new(big.Int).SetBytes(padBytes(orderBuf, orderLen))
	if isBigZero(order) {
		return 0, nil, nil, errors.New("Group order is zero")
	}
	return orderLen, order, rest, nil
}

func decodeG1Point(in []byte, modulusLen int, g1 *g1) (*pointG1, []byte, error) {
	x, rest, err := decodeFp(in, modulusLen, g1.f)
	if err != nil {
		return nil, nil, err
	}
	y, rest, err := decodeFp(rest, modulusLen, g1.f)
	if err != nil {
		return nil, nil, err
	}
	p := g1.fromXY(x, y)
	return p, rest, nil
}

func decodeG22Point(in []byte, modulusLen int, g2 *g22) (*pointG22, []byte, error) {
	x, rest, err := decodeFp2(in, modulusLen, g2.f)
	if err != nil {
		return nil, nil, err
	}
	y, rest, err := decodeFp2(rest, modulusLen, g2.f)
	if err != nil {
		return nil, nil, err
	}
	q := g2.fromXY(x, y)
	return q, rest, nil
}

func decodeG23Point(in []byte, modulusLen int, g2 *g23) (*pointG23, []byte, error) {
	x, rest, err := decodeFp3(in, modulusLen, g2.f)
	if err != nil {
		return nil, nil, err
	}
	y, rest, err := decodeFp3(rest, modulusLen, g2.f)
	if err != nil {
		return nil, nil, err
	}
	q := g2.fromXY(x, y)
	return q, rest, nil
}

func decodeScalar(in []byte, orderLen int, order *big.Int) (*big.Int, []byte, error) {
	buf, rest, err := split(in, orderLen)
	if err != nil {
		return nil, nil, err
	}
	s := new(big.Int).SetBytes(buf)
	// if s.Cmp(order) != -1 {
	// 	return nil, nil, errors.New("Scalar is larger than the group order")
	// }
	return s, rest, nil
}

func createExtension2FieldParams(in []byte, modulusLen int, field *field, degree int, needFrobenius bool) (*fq2, []byte, error) {
	nonResidue, rest, err := decodeFp(in, modulusLen, field)
	if err != nil {
		return nil, nil, err
	}
	if field.isZero(nonResidue) {
		return nil, nil, errors.New("Non-residue for Fp2 is zero")
	}
	if notSquare := isNonNThRoot(field, nonResidue, degree); !notSquare {
		return nil, nil, errors.New("Non-residue for Fp2 is actually a residue")
	}

	fq2, err := newFq2(field, nil)
	fq2.f.copy(fq2.nonResidue, nonResidue)
	if err != nil {
		return nil, nil, err
	}
	if needFrobenius {
		if ok := fq2.calculateFrobeniusCoeffs(); !ok {
			return nil, nil, errors.New("Can not calculate Frobenius coefficients for Fp2")
		}
	}
	return fq2, rest, nil
}

func createExtension3FieldParams(in []byte, modulusLen int, field *field, degree int, needFrobenius bool) (*fq3, []byte, error) {
	nonResidue, rest, err := decodeFp(in, modulusLen, field)
	if err != nil {
		return nil, nil, err
	}
	if field.isZero(nonResidue) {
		return nil, nil, errors.New("Non-residue for Fp3 is zero")
	}
	if ok := isNonNThRoot(field, nonResidue, degree); !ok {
		return nil, nil, errors.New("Non-residue for Fp3 is actually a residue")
	}

	fq3, err := newFq3(field, nil)
	fq3.f.copy(fq3.nonResidue, nonResidue)
	if err != nil {
		return nil, nil, err
	}
	if needFrobenius {
		if ok := fq3.calculateFrobeniusCoeffs(); !ok {
			return nil, nil, errors.New("Can not calculate Frobenius coefficients for Fp3")
		}
	}
	return fq3, rest, nil
}

func isNonNThRoot(field *field, nonResidue fieldElement, power int) bool {
	result := field.newFieldElement()
	exp := new(big.Int).Sub(field.pbig, big.NewInt(1))
	exp, rem := new(big.Int).DivMod(exp, big.NewInt(int64(power)), big.NewInt(0))
	if !isBigZero(rem) {
		return false
	}
	field.exp(result, nonResidue, exp)
	if field.equal(result, field.one) {
		return false
	}
	return true
}

func isNonNThRootFp2(fq2 *fq2, nonResidue *fe2, power int) bool {
	exp := new(big.Int).Mul(fq2.f.pbig, fq2.f.pbig)
	exp = new(big.Int).Sub(exp, big.NewInt(1))
	exp, rem := new(big.Int).DivMod(exp, big.NewInt(int64(power)), big.NewInt(0))
	if !isBigZero(rem) {
		return false
	}
	result := fq2.newElement()
	fq2.exp(result, nonResidue, exp)
	if fq2.equal(result, fq2.one()) {
		return false
	}
	return true
}

func decodeLoopParameters(in []byte, limit int) (*big.Int, []byte, error) {
	lengthBuf, rest, err := split(in, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return nil, nil, errors.New("cant decode modulus length")
	}
	maxLength := (limit + 7) / 8
	length := int(lengthBuf[0])
	if length == 0 {
		return nil, nil, errors.New("Loop parameter scalar has zero length")
	}

	if length > maxLength {
		return nil, nil, errors.New("Scalar is too large for bit length")
	}
	paramBuf, rest, err := split(rest, length)
	if err != nil {
		return nil, nil, errors.New("Input is not long enough to get loop parameter")
	}
	param := new(big.Int).SetBytes(paramBuf)
	if param.BitLen() > limit {
		return nil, nil, errors.New("Scalar is too large for bit length")
	}
	return param, rest, nil
}

func decodeTwistType(in []byte) (int, []byte, error) {
	// twist type 0x01: M, 0x02: D
	twistTypeBuf, rest, err := split(in, TWIST_TYPE_LENGTH)
	if err != nil {
		return 0, nil, errors.New("Input is not long enough to get twist type")
	}
	twistType := int(twistTypeBuf[0])
	if twistType != TWIST_M && twistType != TWIST_D {
		return 0, nil, errors.New("Unknown twist type supplied")
	}
	return twistType, rest, nil
}

func decodePairingExpSign(in []byte) (bool, []byte, error) {
	expIsNegativeBuf, rest, err := split(in, SIGN_ENCODING_LENGTH)
	if err != nil {
		return false, nil, errors.New("exp is not encoded properly")
	}
	switch int(expIsNegativeBuf[0]) {
	case NEGATIVE_EXP:
		return true, rest, nil
	case POSITIVE_EXP:
		return false, rest, nil
	default:
		return false, nil, errors.New("Unknown parameter exp sign")
	}
}

func decodeBoolean(in []byte) (bool, []byte, error) {
	booleanBuf, rest, err := split(in, BOOLEAN_ENCODING_LENGTH)
	if err != nil {
		return false, nil, errors.New("Input is not long enough to get boolean")
	}
	switch int(booleanBuf[0]) {
	case BOOLEAN_FALSE:
		return false, rest, nil
	case BOOLEAN_TRUE:
		return true, rest, nil
	default:
		return false, nil, errors.New("Boolean is not encoded properly")
	}
}

func decodeG1(in []byte) (*g1, int, *big.Int, int, []byte, error) {
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return nil, 0, nil, 0, nil, err
	}
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return nil, 0, nil, 0, nil, err
	}
	orderLen, order, rest, err := decodeGroupOrder(rest)
	if err != nil {
		return nil, 0, nil, 0, nil, err
	}
	g1, err := newG1(field, a, b, order)
	if err != nil {
		return nil, 0, nil, 0, nil, err
	}
	return g1, modulusLen, order, orderLen, rest, nil
}

func decodeG22(in []byte, field *field, modulusLen int) (*g22, *big.Int, int, []byte, error) {
	fq2, rest, err := createExtension2FieldParams(in, modulusLen, field, 2, false)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	a2, b2, rest, err := decodeBAInExtField2FromEncoding(rest, modulusLen, fq2)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	orderLen, order, rest, err := decodeGroupOrder(rest)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	g22, err := newG22(fq2, a2, b2, order)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	return g22, order, orderLen, rest, nil
}

func decodeG23(in []byte, field *field, modulusLen int) (*g23, *big.Int, int, []byte, error) {
	fq3, rest, err := createExtension3FieldParams(in, modulusLen, field, 3, false)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	a2, b2, rest, err := decodeBAInExtField3FromEncoding(rest, modulusLen, fq3)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	orderLen, order, rest, err := decodeGroupOrder(rest)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	g23, err := newG23(fq3, a2, b2, order)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	return g23, order, orderLen, rest, nil
}
