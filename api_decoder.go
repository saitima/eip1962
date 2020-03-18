package eip

import (
	"bytes"
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
	MAX_SCALAR_LEN                          = 128
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

type tape struct {
	data   []byte
	offset int
}

func newTape(in []byte) *tape {
	return &tape{in, 0}
}

func (tape *tape) read(size int) ([]byte, error) {
	if size < 1 {
		return nil, errors.New("size should be positive")
	}
	lower := tape.offset
	upper := lower + size
	if upper > tape.length() {
		return nil, errors.New("exceeds input length")
	}
	r := tape.data[lower:upper]
	tape.offset = upper
	return r, nil
}

func (tape *tape) length() int {
	return len(tape.data)
}

func (tape *tape) remaining() int {
	return tape.length() - tape.offset
}

type cache struct {
	fq            *fq
	modulusLen    int
	groupOrderLen int
}

type context struct {
	willDoPairing bool
}

type decoder struct {
	tape    *tape
	cache   *cache
	context *context
}

func newDecoder(in []byte) *decoder {
	return &decoder{newTape(in), &cache{nil, 0, 0}, &context{false}}
}

func (decoder *decoder) read(size int) ([]byte, error) {
	return decoder.tape.read(size)
}

func (decoder *decoder) modulusLen() int {
	return decoder.cache.modulusLen
}

func (decoder *decoder) groupOrderLen() int {
	return decoder.cache.groupOrderLen
}

func (decoder *decoder) remainingDataLen() int {
	return decoder.tape.remaining()
}

func (decoder *decoder) readFq() (*fq, error) {
	cache := decoder.cache
	if cache.fq != nil {
		return cache.fq, nil
	}
	var err error
	var modulusBuf []byte
	var fq *fq
	// given modulus length decode modulus
	if modulusBuf, err = decoder.readModulus(); err != nil {
		return nil, err
	}
	modulusLen := len(modulusBuf)
	// Use 4 limbs for modulus lower than 256 bit,
	// pad to 32 bytes
	fourLimbBound := 25
	if modulusLen < fourLimbBound {
		modulusBuf = padBytes(modulusBuf, 32)
	} else {
		// otherwise pad to upper 8n bytes
		// fixedLen := (((modulusLen - 1) / 8) + 1) * 8
		modulus := new(big.Int).SetBytes(modulusBuf)
		fixedLen := ((modulus.BitLen() / 64) + 1) * 8
		modulusBuf = padBytes(modulusBuf, fixedLen)
	}
	if fq, err = newField(modulusBuf); err != nil {
		return nil, errors.New(ERR_BASE_FIELD_CONSTRUCTION)
	}
	cache.fq = fq
	return fq, nil
}

func (decoder *decoder) readG1G22Pairs(g1 *g1, g2 *g22) ([]*pointG1, []*pointG22, error) {
	numPairs, err := decoder.readLength()
	if err != nil {
		return nil, nil, err
	}
	if numPairs == 0 {
		return nil, nil, errors.New(ERR_PAIRING_NUM_PAIRS_ZERO)
	}
	var g1Points []*pointG1
	var g2Points []*pointG22
	for i := 0; i < numPairs; i++ {
		checkSubgroup1, err := decoder.readBool()
		if err != nil {
			return nil, nil, err
		}
		p1, err := decoder.readG1Point(g1)
		if err != nil {
			return nil, nil, err
		}
		checkSubgroup2, err := decoder.readBool()
		if err != nil {
			return nil, nil, err
		}
		p2, err := decoder.readG22Point(g2)
		if err != nil {
			return nil, nil, err
		}
		if !g1.isOnCurve(p1) {
			return nil, nil, errors.New(ERR_PAIRING_POINTG1_NOT_ON_CURVE)
		}
		if !g2.isOnCurve(p2) {
			return nil, nil, errors.New(ERR_PAIRING_POINTG2_NOT_ON_CURVE)
		}
		if checkSubgroup1 {
			g1.checkCorrectSubgroup(p1)
		}
		if checkSubgroup2 {
			g2.checkCorrectSubgroup(p2)
		}
		if !g1.isZero(p1) && !g2.isZero(p2) {
			g1Points = append(g1Points, p1)
			g2Points = append(g2Points, p2)
		}
	}
	return g1Points, g2Points, nil
}

func (decoder *decoder) readG1G23Pairs(g1 *g1, g2 *g23) ([]*pointG1, []*pointG23, error) {
	numPairs, err := decoder.readLength()
	if err != nil {
		return nil, nil, err
	}
	if numPairs == 0 {
		return nil, nil, errors.New(ERR_PAIRING_NUM_PAIRS_ZERO)
	}
	var g1Points []*pointG1
	var g2Points []*pointG23
	for i := 0; i < numPairs; i++ {
		checkSubgroup1, err := decoder.readBool()
		if err != nil {
			return nil, nil, err
		}
		p1, err := decoder.readG1Point(g1)
		if err != nil {
			return nil, nil, err
		}
		checkSubgroup2, err := decoder.readBool()
		if err != nil {
			return nil, nil, err
		}
		p2, err := decoder.readG23Point(g2)
		if err != nil {
			return nil, nil, err
		}
		if !g1.isOnCurve(p1) {
			return nil, nil, errors.New(ERR_PAIRING_POINTG1_NOT_ON_CURVE)
		}
		if !g2.isOnCurve(p2) {
			return nil, nil, errors.New(ERR_PAIRING_POINTG2_NOT_ON_CURVE)
		}
		if checkSubgroup1 {
			g1.checkCorrectSubgroup(p1)
		}
		if checkSubgroup2 {
			g2.checkCorrectSubgroup(p2)
		}
		if !g1.isZero(p1) && !g2.isZero(p2) {
			g1Points = append(g1Points, p1)
			g2Points = append(g2Points, p2)
		}
	}
	return g1Points, g2Points, nil
}

func (decoder *decoder) readSign() (bool, error) {
	buf, err := decoder.read(SIGN_ENCODING_LENGTH)
	if err != nil {
		return false, errors.New(ERR_PAIRING_EXP_SIGN_INVALID)
	}
	switch int(buf[0]) {
	case NEGATIVE_EXP:
		return true, nil
	case POSITIVE_EXP:
		return false, nil
	default:
		return false, errors.New(ERR_PAIRING_EXP_SIGN_UNKNWON)
	}
}

func (decoder *decoder) readBool() (bool, error) {
	buf, err := decoder.read(BOOLEAN_ENCODING_LENGTH)
	if err != nil {
		return false, errors.New(ERR_PAIRING_EXP_SIGN_INVALID)
	}
	switch int(buf[0]) {
	case BOOLEAN_FALSE:
		return false, nil
	case BOOLEAN_TRUE:
		return true, nil
	default:
		return false, errors.New(ERR_PAIRING_EXP_SIGN_UNKNWON)
	}
}

func (decoder *decoder) readModulus() ([]byte, error) {
	// read modulus len
	modulusLenBuf, err := decoder.read(BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return nil, errors.New(ERR_BASE_FIELD_MODULUS_LENGTH_NOT_ENOUGH_BYTE)
	}
	modulusLen := int(modulusLenBuf[0])
	if modulusLen > MAX_MODULUS_BYTE_LEN {
		return nil, errors.New(ERR_MODULUS_LENGTH_LARGE)
	}
	if modulusLen <= 0 {
		return nil, errors.New("modulus length must be higher than zero")
	}
	// read modulus
	modulusBuf, err := decoder.read(modulusLen)
	if err != nil {
		return nil, errors.New(ERR_MODULUS_NOT_ENOUGH_BYTE)
	}
	// dense check
	if int(modulusBuf[0]) == 0 {
		return nil, errors.New(ERR_MODULUS_HIGHEST_BYTE)
	}
	lowest := modulusBuf[modulusLen-1]
	// check if even
	if lowest&1 != 1 {
		return nil, errors.New(ERR_MODULUS_EVEN)
	}
	// check if less than 3
	if modulusLen == 1 && lowest > 3 {
		return nil, errors.New(ERR_MODULUS_LESS_THREE)
	}
	decoder.cache.modulusLen = modulusLen
	return modulusBuf, nil
}

func (decoder *decoder) readFieldElement() (fe, error) {
	var buf []byte
	var fe fe
	var err error
	fq, err := decoder.readFq()
	if err != nil {
		return nil, err
	}
	if buf, err = decoder.readFieldElementAsBytes(1); err != nil {
		return nil, err
	}
	if fe, err = fq.fromBytes(buf); err != nil {
		return nil, err
	}
	return fe, nil
}

func (decoder *decoder) readFieldElementAsBytes(size int) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	var err error
	fq, err := decoder.readFq()
	if err != nil {
		return nil, err
	}
	for i := 0; i < size; i++ {
		dat, err := decoder.read(fq.modulusByteLen)
		if err != nil {
			return nil, errors.New(ERR_INPUT_NOT_ENOUGH_FOR_FIELD_ELEMS)
		}
		buf.Write(padBytes(dat, fq.byteSize()))
	}
	return buf.Bytes(), nil
}

func (decoder *decoder) readFq2() (*fq2, error) {
	fq, err := decoder.readFq()
	if err != nil {
		return nil, err
	}
	nonResidue, err := decoder.readNonResidue1(2)
	if err != nil {
		return nil, err
	}
	fq2, _ := newFq2(fq, nil)
	fq.copy(fq2.nonResidue, nonResidue)
	if decoder.context.willDoPairing {
		if ok := fq2.calculateFrobeniusCoeffs(); !ok {
			return nil, errors.New(ERR_EXT_FIELD_FROBENIUS_FOR_FP2)
		}
	}
	return fq2, nil
}

func (decoder *decoder) readFq3() (*fq3, error) {
	fq, err := decoder.readFq()
	if err != nil {
		return nil, err
	}
	nonResidue, err := decoder.readNonResidue1(3)
	if err != nil {
		return nil, err
	}
	fq3, _ := newFq3(fq, nil)
	fq.copy(fq3.nonResidue, nonResidue)
	if decoder.context.willDoPairing {
		if ok := fq3.calculateFrobeniusCoeffs(); !ok {
			return nil, errors.New(ERR_EXT_FIELD_FROBENIUS_FOR_FP2)
		}
	}
	return fq3, nil
}

func (decoder *decoder) readNonResidue1(degree int) (fe, error) {
	// expecting from cache so no error check here
	fq, _ := decoder.readFq()
	nonResidue, err := decoder.readFieldElement()
	if err != nil {
		return nil, err
	}
	if fq.isZero(nonResidue) {
		return nil, errors.New(ERR_EXT_FIELD_NON_RESIDUE_FP2_ZERO)
	}
	if !fq.isNonResidue(nonResidue, degree) {
		return nil, errors.New(ERR_EXT_FIELD_NON_RESIDUE_FP2_RESIDUE)
	}
	return nonResidue, err
}

func (decoder *decoder) readNonResidue2(fq2 *fq2) (*fe2, error) {
	// expecting from cache so no error check here
	buf, err := decoder.readFieldElementAsBytes(2)
	if err != nil {
		return nil, err
	}
	nonResidue, err := fq2.fromBytes(buf)
	if err != nil {
		return nil, err
	}
	if fq2.isZero(nonResidue) {
		return nil, errors.New(ERR_EXT_FIELD_NON_RESIDUE_FP6_ZERO)
	}
	if !fq2.isNonResidue(nonResidue, 6) {
		return nil, errors.New(ERR_EXT_FIELD_NON_RESIDUE_FP6_RESIDUE)
	}
	return nonResidue, err
}

func (decoder *decoder) readLength() (int, error) {
	buf, err := decoder.read(BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return -1, err
	}
	return int(buf[0]), err
}

func (decoder *decoder) readGroupOrder() (*big.Int, error) {
	orderLen, err := decoder.readLength()
	if err != nil {
		return nil, errors.New(ERR_GROUP_ORDER_LENGTH_NOT_ENOUGH_BYTE)
	}
	if orderLen > MAX_GROUP_BYTE_LEN {
		return nil, errors.New(ERR_GROUP_ORDER_LENGTH_LARGE)
	}
	orderBuf, err := decoder.read(orderLen)
	if err != nil {
		return nil, errors.New(ERR_GROUP_ORDER_NOT_ENOUGH_BYTE)
	}
	zero := new(big.Int)
	order := new(big.Int).SetBytes(orderBuf)
	if order.Cmp(zero) == 0 {
		return nil, errors.New(ERR_GROUP_ORDER_ZERO)
	}
	decoder.cache.groupOrderLen = orderLen
	return order, nil
}

func (decoder *decoder) readG1Point(g1 *g1) (*pointG1, error) {
	buf, err := decoder.readFieldElementAsBytes(2)
	if err != nil {
		return nil, err
	}
	return g1.fromBytes(buf)
}

func (decoder *decoder) readG22Point(g22 *g22) (*pointG22, error) {
	buf, err := decoder.readFieldElementAsBytes(4)
	if err != nil {
		return nil, err
	}
	return g22.fromBytes(buf)
}

func (decoder *decoder) readG23Point(g23 *g23) (*pointG23, error) {
	buf, err := decoder.readFieldElementAsBytes(6)
	if err != nil {
		return nil, err
	}
	return g23.fromBytes(buf)
}

func (decoder *decoder) readScalar(size int) (*big.Int, error) {
	buf, err := decoder.read(size)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(buf), nil
}

func (decoder *decoder) readG1() (*g1, error) {
	// read fq
	fq, err := decoder.readFq()
	if err != nil {
		return nil, err
	}
	// read a,b coefficients
	a, b, err := decoder.readAB()
	if err != nil {
		return nil, err
	}
	// read group order length
	order, err := decoder.readGroupOrder()
	if err != nil {
		return nil, err
	}
	// construct g1
	return newG1(fq, a, b, order)
}

func (decoder *decoder) readG22() (*g22, error) {
	// read fq2
	fq2, err := decoder.readFq2()
	if err != nil {
		return nil, err
	}
	// read a,b coefficients
	a, b, err := decoder.readAB2(fq2)
	if err != nil {
		return nil, err
	}
	// read group order length
	order, err := decoder.readGroupOrder()
	if err != nil {
		return nil, err
	}
	// construct g22
	return newG22(fq2, a, b, order)
}

func (decoder *decoder) readG23() (*g23, error) {
	// read fq2
	fq3, err := decoder.readFq3()
	if err != nil {
		return nil, err
	}
	// read a,b coefficients
	a, b, err := decoder.readAB3(fq3)
	if err != nil {
		return nil, err
	}
	// read group order length
	order, err := decoder.readGroupOrder()
	if err != nil {
		return nil, err
	}
	// construct g22
	return newG23(fq3, a, b, order)
}

func (decoder *decoder) readAB() (fe, fe, error) {
	var a, b fe
	var err error
	if a, err = decoder.readFieldElement(); err != nil {
		return nil, nil, errors.New(ERR_INPUT_NOT_ENOUGH_FOR_FIELD_ELEMS)
	}
	if b, err = decoder.readFieldElement(); err != nil {
		return nil, nil, errors.New(ERR_INPUT_NOT_ENOUGH_FOR_FIELD_ELEMS)
	}
	return a, b, nil
}

func (decoder *decoder) readAB2(fq2 *fq2) (*fe2, *fe2, error) {
	var a, b *fe2
	var err error
	aBuf, err := decoder.readFieldElementAsBytes(2)
	if err != nil {
		return nil, nil, errors.New(ERR_INPUT_NOT_ENOUGH_FOR_FIELD_ELEMS)
	}
	a, err = fq2.fromBytes(aBuf)
	if err != nil {
		return nil, nil, err
	}
	bBuf, err := decoder.readFieldElementAsBytes(2)
	if err != nil {
		return nil, nil, errors.New(ERR_INPUT_NOT_ENOUGH_FOR_FIELD_ELEMS)
	}
	b, err = fq2.fromBytes(bBuf)
	if err != nil {
		return nil, nil, err
	}
	return a, b, nil
}

func (decoder *decoder) readAB3(fq3 *fq3) (*fe3, *fe3, error) {
	var a, b *fe3
	var err error
	aBuf, err := decoder.readFieldElementAsBytes(3)
	if err != nil {
		return nil, nil, errors.New(ERR_INPUT_NOT_ENOUGH_FOR_FIELD_ELEMS)
	}
	a, err = fq3.fromBytes(aBuf)
	if err != nil {
		return nil, nil, err
	}
	bBuf, err := decoder.readFieldElementAsBytes(3)
	if err != nil {
		return nil, nil, errors.New(ERR_INPUT_NOT_ENOUGH_FOR_FIELD_ELEMS)
	}
	b, err = fq3.fromBytes(bBuf)
	if err != nil {
		return nil, nil, err
	}
	return a, b, nil
}

func (decoder *decoder) readLoopParam(limit int) (*big.Int, error) {
	length, err := decoder.readLength()
	if err != nil {
		return nil, err
	}
	if length == 0 {
		return nil, errors.New(ERR_PAIRING_LOOP_PARAM_LENGTH_ZERO)
	}
	maxLength := (limit + 7) / 8
	if length > maxLength {
		return nil, errors.New(ERR_PAIRING_LOOP_PARAM_LENGTH_LARGE)
	}
	paramBuf, err := decoder.read(length)
	if err != nil {
		return nil, errors.New(ERR_PAIRING_LOOP_PARAM_NOT_ENOUGH_BYTE)
	}
	if paramBuf[0] == 0 {
		return nil, errors.New(ERR_PAIRING_LOOP_PARAM_TOP_BYTE_ZERO)
	}
	param := new(big.Int).SetBytes(paramBuf)
	if param.BitLen() > limit {
		return nil, errors.New(ERR_PAIRING_LOOP_PARAM_LARGE)
	}
	return param, nil
}

func (decoder *decoder) readTwistType() (int, error) {
	twistTypeBuf, err := decoder.read(TWIST_TYPE_LENGTH)
	if err != nil {
		return -1, errors.New(ERR_PAIRING_TWIST_TYPE_NOT_ENOUGH_BYTE)
	}
	twistType := int(twistTypeBuf[0])
	if twistType != TWIST_M && twistType != TWIST_D {
		return -1, errors.New(ERR_PAIRING_TWIST_TYPE_UNKNOWN)
	}
	return twistType, nil
}

func (decoder *decoder) g1AddRunner() (*g1AddRunner, error) {
	g1, err := decoder.readG1()
	if err != nil {
		return nil, err
	}
	p0, err := decoder.readG1Point(g1)
	if err != nil {
		return nil, err
	}
	if !g1.isOnCurve(p0) {
		return nil, errors.New(ERR_POINT0_NOT_ON_CURVE)
	}
	p1, err := decoder.readG1Point(g1)
	if err != nil {
		return nil, err
	}
	if !g1.isOnCurve(p1) {
		return nil, errors.New(ERR_POINT1_NOT_ON_CURVE)
	}
	return newG1AddRunner(g1, p0, p1), nil
}

func (decoder *decoder) g1MulRunner() (*g1MulRunner, error) {
	g1, err := decoder.readG1()
	scalarLen := decoder.groupOrderLen()
	if err != nil {
		return nil, err
	}
	p0, err := decoder.readG1Point(g1)
	if err != nil {
		return nil, err
	}
	if !g1.isOnCurve(p0) {
		return nil, errors.New(ERR_POINT_NOT_ON_CURVE)
	}
	s, err := decoder.readScalar(scalarLen)
	if err != nil {
		return nil, err
	}
	return newG1MulRunner(g1, p0, s), nil
}

func (decoder *decoder) g1MultiExpRunner() (*g1MultiExpRunner, error) {
	g1, err := decoder.readG1()
	if err != nil {
		return nil, err
	}
	scalarLen := decoder.groupOrderLen()
	modulusLen := decoder.modulusLen()
	if err != nil {
		return nil, err
	}
	numPairs, err := decoder.readLength()
	if err != nil {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	if numPairs == 0 {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIR_LENGTH)
	}
	scalars := make([]*big.Int, numPairs)
	points := make([]*pointG1, numPairs)
	if decoder.remainingDataLen() != numPairs*(scalarLen+(2*modulusLen)) {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIR_INPUT_LENGTH_NOT_MATCH)
	}
	for i := 0; i < numPairs; i++ {
		point, err := decoder.readG1Point(g1)
		if err != nil {
			return nil, err
		}
		if !g1.isOnCurve(point) {
			return nil, errors.New(ERR_POINT_NOT_ON_CURVE)
		}
		scalar, err := decoder.readScalar(scalarLen)
		if err != nil {
			return nil, err
		}
		scalars[i], points[i] = scalar, point
	}
	return newG1MultiExpRunner(g1, points, scalars), nil
}

func (decoder *decoder) g2AddRunner() (runner, error) {
	_, err := decoder.readFq()
	if err != nil {
		return nil, err
	}
	degreeBuf, err := decoder.read(EXTENSION_DEGREE_LENGTH_ENCODING)
	if err != nil {
		return nil, errors.New(ERR_G2_CANT_DECODE_EXT_DEGREE_LENGTH)
	}
	degree := int(degreeBuf[0])
	switch degree {
	case EXTENSION_TWO_DEGREE:
		return decoder.g22AddRunner()
	case EXTENSION_THREE_DEGREE:
		return decoder.g23AddRunner()
	default:
		return nil, errors.New(ERR_G2_UNEXPECTED_EXT_DEGREE)
	}
}

func (decoder *decoder) g2MulRunner() (runner, error) {
	_, err := decoder.readFq()
	if err != nil {
		return nil, err
	}
	degreeBuf, err := decoder.read(EXTENSION_DEGREE_LENGTH_ENCODING)
	if err != nil {
		return nil, errors.New(ERR_G2_CANT_DECODE_EXT_DEGREE_LENGTH)
	}
	degree := int(degreeBuf[0])
	switch degree {
	case EXTENSION_TWO_DEGREE:
		return decoder.g22MulRunner()
	case EXTENSION_THREE_DEGREE:
		return decoder.g23MulRunner()
	default:
		return nil, errors.New(ERR_G2_UNEXPECTED_EXT_DEGREE)
	}
}

func (decoder *decoder) g2MultiExpRunner() (runner, error) {
	_, err := decoder.readFq()
	if err != nil {
		return nil, err
	}
	degreeBuf, err := decoder.read(EXTENSION_DEGREE_LENGTH_ENCODING)
	if err != nil {
		return nil, errors.New(ERR_G2_CANT_DECODE_EXT_DEGREE_LENGTH)
	}
	degree := int(degreeBuf[0])
	switch degree {
	case EXTENSION_TWO_DEGREE:
		return decoder.g22MultiExpRunner()
	case EXTENSION_THREE_DEGREE:
		return decoder.g23MultiExpRunner()
	default:
		return nil, errors.New(ERR_G2_UNEXPECTED_EXT_DEGREE)
	}
}

func (decoder *decoder) g22AddRunner() (*g22AddRunner, error) {
	g22, err := decoder.readG22()
	if err != nil {
		return nil, err
	}
	p0, err := decoder.readG22Point(g22)
	if err != nil {
		return nil, err
	}
	if !g22.isOnCurve(p0) {
		return nil, errors.New(ERR_POINT0_NOT_ON_CURVE)
	}
	p1, err := decoder.readG22Point(g22)
	if err != nil {
		return nil, err
	}
	if !g22.isOnCurve(p1) {
		return nil, errors.New(ERR_POINT1_NOT_ON_CURVE)
	}
	return newG22AddRunner(g22, p0, p1), nil
}

func (decoder *decoder) g22MulRunner() (*g22MulRunner, error) {
	g22, err := decoder.readG22()
	if err != nil {
		return nil, err
	}
	scalarLen := decoder.groupOrderLen()
	if err != nil {
		return nil, err
	}
	p0, err := decoder.readG22Point(g22)
	if err != nil {
		return nil, err
	}
	if !g22.isOnCurve(p0) {
		return nil, errors.New(ERR_POINT_NOT_ON_CURVE)
	}
	s, err := decoder.readScalar(scalarLen)
	if err != nil {
		return nil, err
	}
	return newG22MulRunner(g22, p0, s), nil
}

func (decoder *decoder) g22MultiExpRunner() (*g22MultiExpRunner, error) {
	g22, err := decoder.readG22()
	if err != nil {
		return nil, err
	}
	scalarLen := decoder.groupOrderLen()
	modulusLen := decoder.modulusLen()
	if err != nil {
		return nil, err
	}
	numPairs, err := decoder.readLength()
	if err != nil {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	if numPairs == 0 {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIR_LENGTH)
	}
	scalars := make([]*big.Int, numPairs)
	points := make([]*pointG22, numPairs)
	pointLen := modulusLen * 4
	if decoder.remainingDataLen() != 2*(scalarLen+pointLen) {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIR_INPUT_LENGTH_NOT_MATCH)
	}
	for i := 0; i < numPairs; i++ {
		point, err := decoder.readG22Point(g22)
		if err != nil {
			return nil, err
		}
		if !g22.isOnCurve(point) {
			return nil, errors.New(ERR_POINT_NOT_ON_CURVE)
		}
		scalar, err := decoder.readScalar(scalarLen)
		if err != nil {
			return nil, err
		}
		scalars[i], points[i] = scalar, point
	}
	return newG22MultiExpRunner(g22, points, scalars), nil
}

func (decoder *decoder) g23AddRunner() (*g23AddRunner, error) {
	g23, err := decoder.readG23()
	if err != nil {
		return nil, err
	}
	p0, err := decoder.readG23Point(g23)
	if err != nil {
		return nil, err
	}
	if !g23.isOnCurve(p0) {
		return nil, errors.New(ERR_POINT0_NOT_ON_CURVE)
	}
	p1, err := decoder.readG23Point(g23)
	if err != nil {
		return nil, err
	}
	if !g23.isOnCurve(p1) {
		return nil, errors.New(ERR_POINT1_NOT_ON_CURVE)
	}
	return newG23AddRunner(g23, p0, p1), nil
}

func (decoder *decoder) g23MulRunner() (*g23MulRunner, error) {
	g23, err := decoder.readG23()
	if err != nil {
		return nil, err
	}
	scalarLen := decoder.groupOrderLen()
	if err != nil {
		return nil, err
	}
	p0, err := decoder.readG23Point(g23)
	if err != nil {
		return nil, err
	}
	if !g23.isOnCurve(p0) {
		return nil, errors.New(ERR_POINT_NOT_ON_CURVE)
	}
	s, err := decoder.readScalar(scalarLen)
	if err != nil {
		return nil, err
	}
	return newG23MulRunner(g23, p0, s), nil
}

func (decoder *decoder) g23MultiExpRunner() (*g23MultiExpRunner, error) {
	g23, err := decoder.readG23()
	if err != nil {
		return nil, err
	}
	scalarLen := decoder.groupOrderLen()
	modulusLen := decoder.modulusLen()
	if err != nil {
		return nil, err
	}
	numPairs, err := decoder.readLength()
	if err != nil {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	if numPairs == 0 {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIR_LENGTH)
	}
	scalars := make([]*big.Int, numPairs)
	points := make([]*pointG23, numPairs)
	pointLen := modulusLen * 6
	if decoder.remainingDataLen() != 2*(scalarLen+pointLen) {
		return nil, errors.New(ERR_MULTIEXP_NUM_PAIR_INPUT_LENGTH_NOT_MATCH)
	}
	for i := 0; i < numPairs; i++ {
		point, err := decoder.readG23Point(g23)
		if err != nil {
			return nil, err
		}
		if !g23.isOnCurve(point) {
			return nil, errors.New(ERR_POINT_NOT_ON_CURVE)
		}
		scalar, err := decoder.readScalar(scalarLen)
		if err != nil {
			return nil, err
		}
		scalars[i], points[i] = scalar, point
	}
	return newG23MultiExpRunner(g23, points, scalars), nil
}

func (decoder *decoder) blsRunner() (*blsRunner, error) {
	decoder.context.willDoPairing = true
	g1, err := decoder.readG1()
	if err != nil {
		return nil, err
	}
	fq2, err := decoder.readFq2()
	if err != nil {
		return nil, err
	}
	nonResidue2, err := decoder.readNonResidue2(fq2)
	if err != nil {
		return nil, err
	}
	twistType, err := decoder.readTwistType()
	if err != nil {
		return nil, err
	}
	fq6, err := newFq6Cubic(fq2, nil)
	if err != nil {
		return nil, err
	}
	fq2.copy(fq6.nonResidue, nonResidue2)
	f1, f2, err := constructBaseForFq6AndFq12(fq6)
	if err != nil {
		return nil, errors.New(ERR_EXT_FIELD_BASE_FROBENIUS_FOR_FP612)
	}
	if ok := fq6.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return nil, errors.New(ERR_EXT_FIELD_FROBENIUS_FOR_FP6)
	}
	fq12, err := newFq12(fq6, nil)
	if err != nil {
		return nil, err
	}
	if ok := fq12.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return nil, errors.New(ERR_EXT_FIELD_FROBENIUS_FOR_FP12)
	}
	nonResidueInv := fq2.new()
	if hasInverse := fq2.inverse(nonResidueInv, nonResidue2); !hasInverse {
		return nil, errors.New(ERR_PAIRING_FP2_NON_RESIDUE_NOT_INVERTIBLE)
	}
	a2, b2 := fq2.new(), fq2.new()
	if twistType == TWIST_M {
		fq2.mulByFq(a2, nonResidue2, g1.a)
		fq2.mulByFq(b2, nonResidue2, g1.b)
	} else {
		fq2.mulByFq(a2, nonResidueInv, g1.a)
		fq2.mulByFq(b2, nonResidueInv, g1.b)
	}
	g2, err := newG22(fq2, a2, b2, g1.q)
	if err != nil {
		return nil, err
	}
	z, err := decoder.readLoopParam(MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return nil, err
	}
	if weight := calculateHammingWeight(z); weight > MAX_BLS12_X_HAMMING {
		return nil, errors.New(ERR_BLS_PAIRING_LOW_HAMMING_WEIGHT)
	}
	zIsNegative, err := decoder.readSign()
	if err != nil {
		return nil, err
	}
	g1Points, g2Points, err := decoder.readG1G22Pairs(g1, g2)
	if err != nil {
		return nil, err
	}
	e := newBLSInstance(z, zIsNegative, twistType, g1, g2, fq12, false)
	return newBLSRunner(e, g1Points, g2Points), nil
}

func (decoder *decoder) bnRunner() (*bnRunner, error) {
	decoder.context.willDoPairing = true
	g1, err := decoder.readG1()
	if err != nil {
		return nil, err
	}
	fq2, err := decoder.readFq2()
	if err != nil {
		return nil, err
	}
	nonResidue2, err := decoder.readNonResidue2(fq2)
	if err != nil {
		return nil, err
	}
	twistType, err := decoder.readTwistType()
	if err != nil {
		return nil, err
	}
	fq6, err := newFq6Cubic(fq2, nil)
	if err != nil {
		return nil, err
	}
	fq2.copy(fq6.nonResidue, nonResidue2)
	f1, f2, err := constructBaseForFq6AndFq12(fq6)
	if err != nil {
		return nil, errors.New(ERR_EXT_FIELD_BASE_FROBENIUS_FOR_FP612)
	}
	if ok := fq6.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return nil, errors.New(ERR_EXT_FIELD_FROBENIUS_FOR_FP6)
	}
	fq12, err := newFq12(fq6, nil)
	if err != nil {
		return nil, err
	}
	if ok := fq12.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return nil, errors.New(ERR_EXT_FIELD_FROBENIUS_FOR_FP12)
	}
	nonResidueInv := fq2.new()
	if hasInverse := fq2.inverse(nonResidueInv, nonResidue2); !hasInverse {
		return nil, errors.New(ERR_PAIRING_FP2_NON_RESIDUE_NOT_INVERTIBLE)
	}
	a2, b2 := fq2.new(), fq2.new()
	if twistType == TWIST_M {
		fq2.mulByFq(a2, nonResidue2, g1.a)
		fq2.mulByFq(b2, nonResidue2, g1.b)
	} else {
		fq2.mulByFq(a2, nonResidueInv, g1.a)
		fq2.mulByFq(b2, nonResidueInv, g1.b)
	}
	g2, err := newG22(fq2, a2, b2, g1.q)
	if err != nil {
		return nil, err
	}
	u, err := decoder.readLoopParam(MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return nil, err
	}
	uIsNegative, err := decoder.readSign()
	if err != nil {
		return nil, err
	}
	sixUPlus2 := new(big.Int)
	six, two := big.NewInt(6), big.NewInt(2)
	if uIsNegative {
		sixUPlus2.Mul(six, u)
		sixUPlus2.Sub(sixUPlus2, two)
	} else {
		sixUPlus2.Mul(six, u)
		sixUPlus2.Add(sixUPlus2, two)
	}
	if weight := calculateHammingWeight(sixUPlus2); weight > MAX_BN_SIX_U_PLUS_TWO_HAMMING {
		return nil, errors.New(ERR_BLS_PAIRING_LOW_HAMMING_WEIGHT)
	}
	minus2Inv := new(big.Int).ModInverse(big.NewInt(-2), fq2.modulus())
	nonResidueInPMinus1Over2 := fq2.new()
	fq2.exp(nonResidueInPMinus1Over2, fq6.nonResidue, minus2Inv)

	g1Points, g2Points, err := decoder.readG1G22Pairs(g1, g2)
	if err != nil {
		return nil, err
	}
	e := newBNInstance(u, uIsNegative, twistType, g1, g2, fq12, true)
	return newBNRunner(e, g1Points, g2Points), nil
}

func (decoder *decoder) mnt4Runner() (*mnt4Runner, error) {
	decoder.context.willDoPairing = true
	g1, err := decoder.readG1()
	if err != nil {
		return nil, err
	}
	fq2, err := decoder.readFq2()
	if err != nil {
		return nil, err
	}
	f1 := constructBaseForFq2AndFq4(fq2)
	fq2.calculateFrobeniusCoeffsWithPrecomputation(f1)
	fq4, err := newFq4(fq2, nil)
	if err != nil {
		return nil, err
	}
	fq4.calculateFrobeniusCoeffsWithPrecomputation(f1)
	a2, b2 := fq2.new(), fq2.new()
	twist, twistSquare, twistCube := fq2.twistOne(), fq2.new(), fq2.new()
	fq2.square(twistSquare, twist)
	fq2.mul(twistCube, twistSquare, twist)
	fq2.mulByFq(a2, twistSquare, g1.a)
	fq2.mulByFq(b2, twistCube, g1.b)
	g2, err := newG22(fq2, a2, b2, g1.q)
	if err != nil {
		return nil, err
	}
	x, err := decoder.readLoopParam(MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return nil, err
	}
	xIsNegative, err := decoder.readSign()
	if err != nil {
		return nil, err
	}
	if weight := calculateHammingWeight(x); weight > MAX_ATE_PAIRING_ATE_LOOP_COUNT_HAMMING {
		return nil, errors.New(ERR_MNT_PAIRING_LOW_HAMMING_WEIGHT)
	}
	expW0, err := decoder.readLoopParam(MAX_ATE_PAIRING_FINAL_EXP_W0_BIT_LENGTH)
	if err != nil {
		return nil, err
	}
	expW1, err := decoder.readLoopParam(MAX_ATE_PAIRING_FINAL_EXP_W1_BIT_LENGTH)
	if err != nil {
		return nil, err
	}
	expW0IsNegative, err := decoder.readSign()
	if err != nil {
		return nil, err
	}
	g1Points, g2Points, err := decoder.readG1G22Pairs(g1, g2)
	if err != nil {
		return nil, err
	}
	mnt4 := newMNT4Instance(x, xIsNegative, expW0, expW1, expW0IsNegative, fq4, g1, g2, twist)
	return newMNT4Runner(mnt4, g1Points, g2Points), nil
}

func (decoder *decoder) mnt6Runner() (*mnt6Runner, error) {
	decoder.context.willDoPairing = true
	g1, err := decoder.readG1()
	if err != nil {
		return nil, err
	}
	fq3, err := decoder.readFq3()
	if err != nil {
		return nil, err
	}
	f1, err := constructBaseForFq3AndFq6(fq3)
	if err != nil {
		return nil, err
	}
	fq3.calculateFrobeniusCoeffsWithPrecomputation(f1)
	fq6, err := newFq6Quadratic(fq3, nil)
	if err != nil {
		return nil, err
	}
	fq6.calculateFrobeniusCoeffsWithPrecomputation(f1)
	a2, b2 := fq3.new(), fq3.new()
	twist, twistSquare, twistCube := fq3.twistOne(), fq3.new(), fq3.new()
	fq3.square(twistSquare, twist)
	fq3.mul(twistCube, twistSquare, twist)
	fq3.mulByFq(a2, twistSquare, g1.a)
	fq3.mulByFq(b2, twistCube, g1.b)
	g2, err := newG23(fq3, a2, b2, g1.q)
	if err != nil {
		return nil, err
	}
	x, err := decoder.readLoopParam(MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return nil, err
	}
	xIsNegative, err := decoder.readSign()
	if err != nil {
		return nil, err
	}
	if weight := calculateHammingWeight(x); weight > MAX_ATE_PAIRING_ATE_LOOP_COUNT_HAMMING {
		return nil, errors.New(ERR_MNT_PAIRING_LOW_HAMMING_WEIGHT)
	}
	expW0, err := decoder.readLoopParam(MAX_ATE_PAIRING_FINAL_EXP_W0_BIT_LENGTH)
	if err != nil {
		return nil, err
	}
	expW1, err := decoder.readLoopParam(MAX_ATE_PAIRING_FINAL_EXP_W1_BIT_LENGTH)
	if err != nil {
		return nil, err
	}
	expW0IsNegative, err := decoder.readSign()
	if err != nil {
		return nil, err
	}
	g1Points, g2Points, err := decoder.readG1G23Pairs(g1, g2)
	if err != nil {
		return nil, err
	}
	mnt6 := newMNT6Instance(x, xIsNegative, expW0, expW1, expW0IsNegative, fq6, g1, g2, twist)
	return newMNT6Runner(mnt6, g1Points, g2Points), nil
}
