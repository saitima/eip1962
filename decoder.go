package fp

import (
	"errors"
	"fmt"
	"math/big"
)

const (
	BYTES_FOR_LENGTH_ENCODING        = 1
	EXTENSION_DEGREE_LENGTH_ENCODING = 1
	EXTENSION_TWO_DEGREE             = 2
	EXTENSION_THREE_DEGREE           = 3
	TWIST_TYPE_LENGTH                = 1
	MAX_BN_U_BIT_LENGTH              = 128
	MAX_BLS12_X_BIT_LENGTH           = 128
	MAX_ATE_PAIRING_ATE_LOOP_COUNT   = 2032
	SIGN_ENCODING_LENGTH             = 1
	MAX_BN_SIX_U_PLUS_TWO_HAMMING    = 128
	MAX_BLS12_X_HAMMING              = 128
)

func split(in []byte, offset int) ([]byte, []byte, error) {
	if len(in) < offset {
		return nil, nil, errors.New(fmt.Sprintf("cant split at given offset %d", offset))
	}
	return in[0:offset], in[offset:], nil
}

func decodeBaseFieldParams(in []byte) (*big.Int, int, []byte, error) {
	modulusLenBuf, rest, err := split(in, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return nil, 0, nil, errors.New("cant decode modulus length")
	}
	modulusLen := int(modulusLenBuf[0])
	modulusBuf, rest, err := split(rest, modulusLen)
	if err != nil {
		return nil, 0, nil, errors.New("cant decode modulus")
	}
	modulus := new(big.Int).SetBytes(modulusBuf)
	return modulus, modulusLen, rest, nil
}

// G1
func parseBaseFieldFromEncoding(in []byte) (*field, *big.Int, int, []byte, error) {
	modulus, modulusLen, rest, err := decodeBaseFieldParams(in)
	if err != nil {
		return nil, nil, 0, nil, err
	}
	if len(rest) < modulusLen {
		return nil, nil, 0, nil, errors.New("Input is not long enough")
	}
	field := newField(modulus.Bytes())
	return field, modulus, modulusLen, rest, nil
}

func decodeFp(in []byte, modulusLen int, field *field) (fieldElement, []byte, error) {
	xBuf, rest, err := split(in, modulusLen)
	if err != nil {
		return nil, nil, err
	}
	x, err := field.newFieldElementFromBytes(xBuf)
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
	c0, err := field.f.newFieldElementFromBytes(c0Buf)
	if err != nil {
		return nil, nil, err
	}
	c1Buf, rest, err := split(rest, modulusLen)
	if err != nil {
		return nil, nil, err
	}
	c1, err := field.f.newFieldElementFromBytes(c1Buf)
	if err != nil {
		return nil, nil, err
	}
	elem := field.newElement()
	field.f.cpy(elem[0], c0)
	field.f.cpy(elem[1], c1)
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

func decodeG1CurveParams(in []byte) ([]byte, int, []byte, error) {
	orderLenBuf, rest, err := split(in, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return nil, 0, nil, err
	}
	orderLen := int(orderLenBuf[0])
	orderBuf, rest, err := split(rest, orderLen)
	if err != nil {
		return nil, 0, nil, err
	}
	return orderBuf, orderLen, rest, nil

}

func parseGroupOrder(in []byte, modulsuLen int) (int, *big.Int, []byte, error) {
	orderBuf, orderLen, rest, err := decodeG1CurveParams(in)
	if err != nil {
		return 0, nil, nil, err
	}
	order := new(big.Int).SetBytes(orderBuf)

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
	p := g1.newPoint()
	g1.f.cpy(p[0], x)
	g1.f.cpy(p[1], y)
	g1.f.cpy(p[2], g1.f.one)
	if !g1.isOnCurve(p) {
		return nil, nil, errors.New("g1 point isn't on the curve")
	}
	return p, rest, nil
}

func decodeG2Point(in []byte, modulusLen int, g2 *g22) (*pointG22, []byte, error) {
	x, rest, err := decodeFp2(in, modulusLen, g2.f)
	if err != nil {
		return nil, nil, err
	}
	y, rest, err := decodeFp2(rest, modulusLen, g2.f)
	if err != nil {
		return nil, nil, err
	}
	q := g2.newPoint()
	g2.f.copy(q[0], x)
	g2.f.copy(q[1], y)
	g2.f.copy(q[2], g2.f.one())
	if !g2.isOnCurve(q) {
		return nil, nil, errors.New("g2 point isn't on the curve")
	}
	return q, rest, nil
}

func decodeScalar(in []byte, orderLen int, order *big.Int) (*big.Int, []byte, error) {
	buf, rest, err := split(in, orderLen)
	if err != nil {
		return nil, nil, err
	}
	s := new(big.Int).SetBytes(buf)
	if s.Cmp(order) != -1 {
		return nil, nil, errors.New("Scalar is larger than the group order")
	}
	return s, rest, nil
}

// G2
func createExtension2FieldParams(in []byte, modulusLen int, field *field, frobenius bool) (*fq2, []byte, error) {
	degreeBuf, rest, err := split(in, EXTENSION_DEGREE_LENGTH_ENCODING)
	if err != nil {
		return nil, nil, errors.New("cant decode extension degree length")
	}
	degree := int(degreeBuf[0])
	if degree != EXTENSION_TWO_DEGREE {
		return nil, nil, errors.New("Extension degree expected to be 2")
	}
	nonResidue, rest, err := decodeFp(rest, modulusLen, field)
	if err != nil {
		return nil, nil, err
	}
	if !isNonNThRoot(field, nonResidue, 2) {
		return nil, nil, errors.New("Non-residue for Fp2 is actually a residue")
	}

	fq2, err := newFq2(field, nil)
	fq2.f.cpy(fq2.nonResidue, nonResidue)
	if err != nil {
		return nil, nil, err
	}
	if frobenius {
		fq2.calculateFrobeniusCoeffs()
	}
	return fq2, rest, nil
}

func isNonNThRoot(field *field, nonResidue fieldElement, power int) bool {
	result := field.newFieldElement()
	field.exp(result, nonResidue, field.pbig)
	if field.equal(result, field.one) {
		return false
	}
	return true
}
func isNonNThRootFp2(fq2 *fq2, nonResidue *fe2, power int) bool {
	result := fq2.newElement()
	fq2.exp(result, nonResidue, fq2.f.pbig)
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

func calculateHammingWeight(s *big.Int) int {
	weight := 0
	for i := range s.Bits() {
		if i == 1 {
			weight++
		}
	}
	return weight
}
