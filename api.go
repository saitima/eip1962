package eip

import (
	"math/big"
)

var (
	zero                         = []byte{0x00}
	pairingError, pairingSuccess = []byte{0x00}, []byte{0x01}
)

const (
	USE_4LIMBS_FOR_LOWER_LIMBS  = true
	TWIST_M, TWIST_D            = 0x01, 0x02
	NEGATIVE_EXP, POSITIVE_EXP  = 0x01, 0x00
	BOOLEAN_FALSE, BOOLEAN_TRUE = 0x00, 0x01
	OPERATION_G1_ADD            = 0x01
	OPERATION_G1_MUL            = 0x02
	OPERATION_G1_MULTIEXP       = 0x03
	OPERATION_G2_ADD            = 0x04
	OPERATION_G2_MUL            = 0x05
	OPERATION_G2_MULTIEXP       = 0x06
	OPERATION_BLS12PAIR         = 0x07
	OPERATION_BNPAIR            = 0x08
	OPERATION_MNT4PAIR          = 0x09
	OPERATION_MNT6PAIR          = 0x0a
)

type API struct{}

func (api *API) Run(opType int, in []byte) ([]byte, error) {
	switch opType {
	case OPERATION_G1_ADD:
		return new(g1Api).addPoints(in)
	case OPERATION_G1_MUL:
		return new(g1Api).mulPoint(in)
	case OPERATION_G1_MULTIEXP:
		return new(g1Api).multiExp(in)
	case OPERATION_G2_ADD, OPERATION_G2_MUL, OPERATION_G2_MULTIEXP:
		return new(g2Api).run(opType, in)
	case OPERATION_BLS12PAIR:
		return pairBLS(in)
	case OPERATION_BNPAIR:
		return pairBN(in)
	case OPERATION_MNT4PAIR:
		return pairMNT4(in)
	case OPERATION_MNT6PAIR:
		return pairMNT6(in)
	default:
		return apiDecodingErr(ERR_UNKNOWN_OPERATION)
	}
}

type g1Api struct{}

func (api *g1Api) addPoints(in []byte) ([]byte, error) {
	g1, modulusLen, _, _, rest, err := decodeG1(in)
	if err != nil {
		return apiDecodingErr(err)
	}
	p0, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return apiDecodingErr(err)
	}
	p1, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return apiDecodingErr(err)
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if !g1.isOnCurve(p0) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_POINT0_NOT_ON_CURVE)
		}
	}
	if !g1.isOnCurve(p1) {
		if GAS_METERING_MODE {
			return apiExecErr(ERR_POINT1_NOT_ON_CURVE)
		}
	}
	g1.add(p0, p0, p1)

	out := make([]byte, 2*modulusLen)
	encodeG1Point(out, g1.toBytes(p0))
	return out, nil
}

func (api *g1Api) mulPoint(in []byte) ([]byte, error) {
	g1, modulusLen, order, orderLen, rest, err := decodeG1(in)
	if err != nil {
		return apiDecodingErr(err)
	}
	p, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return apiDecodingErr(err)
	}
	scalar, rest, err := decodeScalar(rest, orderLen, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if !g1.isOnCurve(p) {
		if GAS_METERING_MODE {
			return apiExecErr(ERR_POINT_NOT_ON_CURVE)
		}
	}
	g1.mulScalar(p, p, scalar)
	out := make([]byte, 2*modulusLen)
	encodeG1Point(out, g1.toBytes(p))
	return out, nil
}

func (api *g1Api) multiExp(in []byte) ([]byte, error) {
	g1, modulusLen, order, orderLen, rest, err := decodeG1(in)
	if err != nil {
		return apiDecodingErr(err)
	}
	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIR_LENGTH)
	}
	if len(rest) != (2*modulusLen+orderLen)*numPairs {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIR_INPUT_LENGTH_NOT_MATCH)
	}
	bases := make([]*pointG1, numPairs)
	scalars := make([]*big.Int, numPairs)
	for i := 0; i < numPairs; i++ {
		p, localRest, err := decodeG1Point(rest, modulusLen, g1)
		if err != nil {
			return apiDecodingErr(err)
		}
		if !g1.isOnCurve(p) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_POINT_NOT_ON_CURVE)
			}
		}
		scalar, localRest, err := decodeScalar(localRest, orderLen, order)
		if err != nil {
			return apiDecodingErr(err)
		}
		bases[i], scalars[i] = g1.newPoint(), new(big.Int)
		g1.copy(bases[i], p)
		scalars[i].Set(scalar)
		rest = localRest
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}

	p := g1.newPoint()
	out := make([]byte, 2*modulusLen)
	if len(bases) != len(scalars) || len(bases) == 0 {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_MULTIEXP_EMPTY_INPUT_PAIRS)
		}
		g1.copy(p, g1.inf)
		encodeG1Point(out, g1.toBytes(p))
		return out, nil
	}
	g1.multiExp(p, bases, scalars)
	encodeG1Point(out, g1.toBytes(p))
	return out, nil
}

type g2Api struct{}

func (api *g2Api) run(opType int, in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := decodeBaseFieldFromEncoding(in)
	if err != nil {
		return apiDecodingErr(err)
	}
	degreeBuf, rest, err := split(rest, EXTENSION_DEGREE_LENGTH_ENCODING)
	if err != nil {
		return apiDecodingErr(ERR_G2_CANT_DECODE_EXT_DEGREE_LENGTH)
	}
	degree := int(degreeBuf[0])
	// fmt.Printf("ext degree %d\n", degree)
	switch degree {
	case EXTENSION_TWO_DEGREE:
		return new(g22Api).run(opType, field, modulusLen, rest)
	case EXTENSION_THREE_DEGREE:
		return new(g23Api).run(opType, field, modulusLen, rest)
	default:
		return apiDecodingErr(ERR_G2_UNEXPECTED_EXT_DEGREE)
	}
}

type g22Api struct{}

func (api *g22Api) run(opType int, field *field, modulusLen int, in []byte) ([]byte, error) {
	switch opType {
	case 0x04:
		return api.addPoints(field, modulusLen, in)
	case 0x05:
		return api.mulPoint(field, modulusLen, in)
	case 0x06:
		return api.multiExp(field, modulusLen, in)
	default:
		return apiDecodingErr(ERR_G2_UNKNOWN_OPERATION)
	}
}

func (api *g22Api) addPoints(field *field, modulusLen int, in []byte) ([]byte, error) {
	g2, _, _, rest, err := decodeG22(in, field, modulusLen)
	if err != nil {
		return apiDecodingErr(err)
	}
	q0, rest, err := decodeG22Point(rest, modulusLen, g2)
	if err != nil {
		return apiDecodingErr(err)
	}
	q1, rest, err := decodeG22Point(rest, modulusLen, g2)
	if err != nil {
		return apiDecodingErr(err)
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if !g2.isOnCurve(q0) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_POINT0_NOT_ON_CURVE)
		}
	}
	if !g2.isOnCurve(q1) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_POINT1_NOT_ON_CURVE)
		}
	}
	g2.add(q0, q0, q1)
	out := make([]byte, 4*modulusLen)
	encodeG22Point(out, g2.toBytes(q0))
	return out, nil
}

func (api *g22Api) mulPoint(field *field, modulusLen int, in []byte) ([]byte, error) {
	g2, order, orderLen, rest, err := decodeG22(in, field, modulusLen)
	if err != nil {
		return apiDecodingErr(err)
	}
	q, rest, err := decodeG22Point(rest, modulusLen, g2)
	if err != nil {
		return apiDecodingErr(err)
	}
	scalar, rest, err := decodeScalar(rest, orderLen, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if !g2.isOnCurve(q) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_POINT_NOT_ON_CURVE)
		}
	}
	g2.mulScalar(q, q, scalar)
	out := make([]byte, 4*modulusLen)
	encodeG22Point(out, g2.toBytes(q))
	return out, nil
}

func (api *g22Api) multiExp(field *field, modulusLen int, in []byte) ([]byte, error) {
	g2, order, orderLen, rest, err := decodeG22(in, field, modulusLen)
	if err != nil {
		return apiDecodingErr(err)
	}
	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIR_LENGTH)
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIR_LENGTH)
	}
	if len(rest) != (4*modulusLen+orderLen)*numPairs {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIR_INPUT_LENGTH_NOT_MATCH)
	}
	bases := make([]*pointG22, numPairs)
	scalars := make([]*big.Int, numPairs)
	for i := 0; i < numPairs; i++ {
		q, localRest, err := decodeG22Point(rest, modulusLen, g2)
		if err != nil {
			return apiDecodingErr(err)
		}
		if !g2.isOnCurve(q) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_POINT_NOT_ON_CURVE)
			}
		}
		scalar, localRest, err := decodeScalar(localRest, orderLen, order)
		if err != nil {
			return apiDecodingErr(err)
		}
		bases[i], scalars[i] = g2.newPoint(), new(big.Int)
		g2.copy(bases[i], q)
		scalars[i].Set(scalar)
		rest = localRest
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}

	q := g2.newPoint()
	out := make([]byte, 4*modulusLen)
	if len(bases) != len(scalars) || len(bases) == 0 {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_MULTIEXP_EMPTY_INPUT_PAIRS)
		}
		g2.copy(q, g2.inf)
		encodeG22Point(out, g2.toBytes(q))
		return out, nil
	}
	g2.multiExp(q, bases, scalars)
	encodeG22Point(out, g2.toBytes(q))
	return out, nil
}

type g23Api struct{}

func (api *g23Api) run(opType int, field *field, modulusLen int, in []byte) ([]byte, error) {
	switch opType {
	case 0x04:
		return api.addPoints(field, modulusLen, in)
	case 0x05:
		return api.mulPoint(field, modulusLen, in)
	case 0x06:
		return api.multiExp(field, modulusLen, in)
	default:
		return apiDecodingErr(ERR_G2_UNKNOWN_OPERATION)
	}
}

func (api *g23Api) addPoints(field *field, modulusLen int, in []byte) ([]byte, error) {
	g2, _, _, rest, err := decodeG23(in, field, modulusLen)
	if err != nil {
		return apiDecodingErr(err)
	}
	q0, rest, err := decodeG23Point(rest, modulusLen, g2)
	if err != nil {
		return apiDecodingErr(err)
	}
	q1, rest, err := decodeG23Point(rest, modulusLen, g2)
	if err != nil {
		return apiDecodingErr(err)
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if !g2.isOnCurve(q0) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_POINT0_NOT_ON_CURVE)
		}
	}
	if !g2.isOnCurve(q1) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_POINT1_NOT_ON_CURVE)
		}
	}
	g2.add(q0, q0, q1)
	out := make([]byte, 6*modulusLen)
	encodeG23Point(out, g2.toBytes(q0))
	return out, nil
}

func (api *g23Api) mulPoint(field *field, modulusLen int, in []byte) ([]byte, error) {
	g2, order, orderLen, rest, err := decodeG23(in, field, modulusLen)
	if err != nil {
		return apiDecodingErr(err)
	}
	q, rest, err := decodeG23Point(rest, modulusLen, g2)
	if err != nil {
		return apiDecodingErr(err)
	}
	s, rest, err := decodeScalar(rest, orderLen, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if !g2.isOnCurve(q) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_POINT_NOT_ON_CURVE)
		}
	}
	g2.mulScalar(q, q, s)
	out := make([]byte, 6*modulusLen)
	encodeG23Point(out, g2.toBytes(q))
	return out, nil
}

func (api *g23Api) multiExp(field *field, modulusLen int, in []byte) ([]byte, error) {
	g2, order, orderLen, rest, err := decodeG23(in, field, modulusLen)
	if err != nil {
		return apiDecodingErr(err)
	}
	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIR_LENGTH)
	}

	if len(rest) != (6*modulusLen+orderLen)*numPairs {
		return apiDecodingErr(ERR_MULTIEXP_NUM_PAIR_INPUT_LENGTH_NOT_MATCH)
	}
	bases := make([]*pointG23, numPairs)
	scalars := make([]*big.Int, numPairs)
	for i := 0; i < numPairs; i++ {
		q, localRest, err := decodeG23Point(rest, modulusLen, g2)
		if err != nil {
			return apiDecodingErr(err)
		}
		if !g2.isOnCurve(q) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_POINT_NOT_ON_CURVE)
			}
		}
		scalar, localRest, err := decodeScalar(localRest, orderLen, order)
		if err != nil {
			return apiDecodingErr(err)
		}
		bases[i], scalars[i] = g2.newPoint(), new(big.Int)
		g2.copy(bases[i], q)
		scalars[i].Set(scalar)
		rest = localRest
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}

	q := g2.newPoint()
	out := make([]byte, 6*modulusLen)
	if len(bases) != len(scalars) || len(bases) == 0 {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_MULTIEXP_EMPTY_INPUT_PAIRS)
		}
		g2.copy(q, g2.inf)
		encodeG23Point(out, g2.toBytes(q))
		return out, nil
	}

	g2.multiExp(q, bases, scalars)
	encodeG23Point(out, g2.toBytes(q))
	return out, nil
}

func pairBN(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := decodeBaseFieldFromEncoding(in)
	if err != nil {
		return apiDecodingErr(err)
	}
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return apiDecodingErr(err)
	}
	if !field.isZero(a) {
		return apiDecodingErr(ERR_BN_PAIRING_A_PARAMETER_NOT_ZERO)
	}
	_, order, rest, err := decodeGroupOrder(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	g1, err := newG1(field, a, b, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq2, rest, err := createExtension2FieldParams(rest, modulusLen, field, 2, true)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq2NonResidue, rest, err := decodeFp2(rest, modulusLen, fq2)
	if err != nil {
		return apiDecodingErr(err)
	}
	if fq2.isZero(fq2NonResidue) {
		return apiDecodingErr(ERR_EXT_FIELD_NON_RESIDUE_FP6_ZERO)
	}
	if !isNonNThRootFp2(fq2, fq2NonResidue, 6) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_EXT_FIELD_NON_RESIDUE_FP6_RESIDUE)
		}
	}
	twistType, rest, err := decodeTwistType(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	f1, f2, err := constructBaseForFq6AndFq12(fq2, fq2NonResidue)
	if err != nil {
		return apiDecodingErr(ERR_EXT_FIELD_BASE_FROBENIUS_FOR_FP612)
	}
	fq6, err := newFq6(fq2, nil)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq2.copy(fq6.nonResidue, fq2NonResidue)
	if ok := fq6.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return apiDecodingErr(ERR_EXT_FIELD_FROBENIUS_FOR_FP6)
	}
	fq12, err := newFq12(fq6, nil)
	if err != nil {
		return apiDecodingErr(err)
	}
	if ok := fq12.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return apiDecodingErr(ERR_EXT_FIELD_FROBENIUS_FOR_FP12)
	}
	fq2NonResidueInv := fq2.newElement()
	if hasInverse := fq2.inverse(fq2NonResidueInv, fq2NonResidue); !hasInverse {
		return apiDecodingErr(ERR_PAIRING_FP2_NON_RESIDUE_NOT_INVERTIBLE)
	}
	b2 := fq2.newElement()
	if twistType == TWIST_M {
		fq2.mulByFq(b2, fq6.nonResidue, b)
	} else {
		fq2.mulByFq(b2, fq2NonResidueInv, b)
	}
	g2, err := newG22(fq2, fq2.zero(), b2, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	u, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return apiDecodingErr(err)
	}
	if isBigZero(u) {
		return apiDecodingErr(ERR_PAIRING_LOOP_COUNT_PARAM_ZERO)
	}
	uIsNegative, rest, err := decodePairingExpSign(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	var sixUPlus2 *big.Int
	six, two := big.NewInt(6), big.NewInt(2)
	if uIsNegative {
		sixUPlus2 = new(big.Int).Mul(six, u)
		sixUPlus2 = new(big.Int).Sub(sixUPlus2, two)
	} else {
		sixUPlus2 = new(big.Int).Mul(six, u)
		sixUPlus2 = new(big.Int).Add(sixUPlus2, two)
	}
	if weight := calculateHammingWeight(sixUPlus2); weight > MAX_BN_SIX_U_PLUS_TWO_HAMMING {
		return apiDecodingErr(ERR_BN_PAIRING_LOW_HAMMING_WEIGHT)
	}
	minus2Inv := new(big.Int).ModInverse(big.NewInt(-2), field.pbig)
	nonResidueInPMinus1Over2 := fq2.newElement()
	fq2.exp(nonResidueInPMinus1Over2, fq6.nonResidue, minus2Inv)

	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return apiDecodingErr(ERR_PAIRING_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_PAIRING_NUM_PAIRS_ZERO)
		}
	}
	var g1Points []*pointG1
	var g2Points []*pointG22
	for i := 0; i < numPairs; i++ {
		needG1SubGroupCheck, localRest, err := decodeBoolean(rest)
		if err != nil {
			return apiDecodingErr(err)
		}
		g1Point, localRest, err := decodeG1Point(localRest, modulusLen, g1)
		if err != nil {
			return apiDecodingErr(err)
		}
		needG2SubGroupCheck, localRest, err := decodeBoolean(localRest)
		if err != nil {
			return apiDecodingErr(err)
		}
		g2Point, localRest, err := decodeG22Point(localRest, modulusLen, g2)
		if err != nil {
			return apiDecodingErr(err)
		}
		if !g1.isOnCurve(g1Point) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_PAIRING_POINTG1_NOT_ON_CURVE)
			}
		}
		if !g2.isOnCurve(g2Point) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_PAIRING_POINTG2_NOT_ON_CURVE)
			}
		}
		if needG1SubGroupCheck {
			if ok := g1.checkCorrectSubGroup(g1Point); !ok {
				if !GAS_METERING_MODE {
					return apiExecErr(ERR_PAIRING_POINTG1_NOT_IN_SUBGROUP)
				}
			}
		}
		if needG2SubGroupCheck {
			if ok := g2.checkCorrectSubGroup(g2Point); !ok {
				if !GAS_METERING_MODE {
					return apiExecErr(ERR_PAIRING_POINTG2_NOT_IN_SUBGROUP)
				}
			}
		}
		if !g1.isZero(g1Point) && !g2.isZero(g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if len(g1Points) == 0 {
		return pairingSuccess, nil
	}

	engine := newBNInstance(
		u,
		sixUPlus2,
		uIsNegative,
		twistType,
		g1,
		g2,
		fq12,
		nonResidueInPMinus1Over2,
		true,
	)
	result, hasValue := engine.multiPair(g1Points, g2Points)
	if !hasValue {
		return apiDecodingErr(ERR_PAIRING_NO_RETURN_VALUE)
	}
	if !fq12.equal(result, fq12.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}

func pairBLS(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := decodeBaseFieldFromEncoding(in)
	if err != nil {
		return apiDecodingErr(err)
	}
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return apiDecodingErr(err)
	}
	if !field.isZero(a) {
		return apiDecodingErr(ERR_BLS_PAIRING_A_PARAMETER_NOT_ZERO)
	}
	_, order, rest, err := decodeGroupOrder(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	g1, err := newG1(field, a, b, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq2, rest, err := createExtension2FieldParams(rest, modulusLen, field, 2, true)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq2NonResidue, rest, err := decodeFp2(rest, modulusLen, fq2)
	if err != nil {
		return apiDecodingErr(err)
	}
	if fq2.isZero(fq2NonResidue) {
		return apiDecodingErr(ERR_EXT_FIELD_NON_RESIDUE_FP6_ZERO)
	}
	if !isNonNThRootFp2(fq2, fq2NonResidue, 6) {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_EXT_FIELD_NON_RESIDUE_FP6_RESIDUE)
		}
	}
	twistType, rest, err := decodeTwistType(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	f1, f2, err := constructBaseForFq6AndFq12(fq2, fq2NonResidue)
	if err != nil {
		return apiDecodingErr(ERR_EXT_FIELD_BASE_FROBENIUS_FOR_FP612)
	}
	fq6, err := newFq6(fq2, nil)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq2.copy(fq6.nonResidue, fq2NonResidue)
	if ok := fq6.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return apiDecodingErr(ERR_EXT_FIELD_FROBENIUS_FOR_FP6)
	}
	fq12, err := newFq12(fq6, nil)
	if err != nil {
		return apiDecodingErr(err)
	}
	if ok := fq12.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return apiDecodingErr(ERR_EXT_FIELD_FROBENIUS_FOR_FP12)
	}
	fq2NonResidueInv := fq2.newElement()
	if hasInverse := fq2.inverse(fq2NonResidueInv, fq2NonResidue); !hasInverse {
		return apiDecodingErr(ERR_PAIRING_FP2_NON_RESIDUE_NOT_INVERTIBLE)
	}
	b2 := fq2.newElement()
	if twistType == TWIST_M {
		fq2.mulByFq(b2, fq6.nonResidue, b)
	} else {
		fq2.mulByFq(b2, fq2NonResidueInv, b)
	}
	g2, err := newG22(fq2, fq2.zero(), b2, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	z, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return apiDecodingErr(err)
	}
	if z.Cmp(big.NewInt(0)) == 0 {
		return apiDecodingErr(ERR_PAIRING_LOOP_COUNT_PARAM_ZERO)
	}

	if weight := calculateHammingWeight(z); weight > MAX_BLS12_X_HAMMING {
		return apiDecodingErr(ERR_BLS_PAIRING_LOW_HAMMING_WEIGHT)
	}
	zIsNegative, rest, err := decodePairingExpSign(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return apiDecodingErr(ERR_PAIRING_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_PAIRING_NUM_PAIRS_ZERO)
		}
	}
	var g1Points []*pointG1
	var g2Points []*pointG22
	for i := 0; i < numPairs; i++ {
		needG1SubGroupCheck, localRest, err := decodeBoolean(rest)
		if err != nil {
			return apiDecodingErr(err)
		}
		g1Point, localRest, err := decodeG1Point(localRest, modulusLen, g1)
		if err != nil {
			return apiDecodingErr(err)
		}
		needG2SubGroupCheck, localRest, err := decodeBoolean(localRest)
		if err != nil {
			return apiDecodingErr(err)
		}
		g2Point, localRest, err := decodeG22Point(localRest, modulusLen, g2)
		if err != nil {
			return apiDecodingErr(err)
		}
		if !g1.isOnCurve(g1Point) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_PAIRING_POINTG1_NOT_ON_CURVE)
			}
		}
		if !g2.isOnCurve(g2Point) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_PAIRING_POINTG2_NOT_ON_CURVE)
			}
		}

		if needG1SubGroupCheck {
			if ok := g1.checkCorrectSubGroup(g1Point); !ok {
				if !GAS_METERING_MODE {
					return apiExecErr(ERR_PAIRING_POINTG1_NOT_IN_SUBGROUP)
				}
			}
		}
		if needG2SubGroupCheck {
			if ok := g2.checkCorrectSubGroup(g2Point); !ok {
				if !GAS_METERING_MODE {
					return apiExecErr(ERR_PAIRING_POINTG2_NOT_IN_SUBGROUP)
				}
			}
		}
		if !g1.isZero(g1Point) && !g2.isZero(g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if len(g1Points) == 0 {
		return pairingSuccess, nil
	}
	engine := newBLSInstance(
		z,
		zIsNegative,
		twistType,
		g1,
		g2,
		fq12,
		false,
	)
	result, hasValue := engine.multiPair(g1Points, g2Points)
	if !hasValue {
		return apiDecodingErr(ERR_PAIRING_NO_RETURN_VALUE)
	}
	if !fq12.equal(result, fq12.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}

func pairMNT4(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := decodeBaseFieldFromEncoding(in)
	if err != nil {
		return apiDecodingErr(err)
	}
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return apiDecodingErr(err)
	}
	_, order, rest, err := decodeGroupOrder(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	g1, err := newG1(field, a, b, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq2, rest, err := createExtension2FieldParams(rest, modulusLen, field, 4, false)
	if err != nil {
		return apiDecodingErr(err)
	}
	f1 := constructBaseForFq2AndFq4(field, fq2.nonResidue)
	fq2.calculateFrobeniusCoeffsWithPrecomputation(f1)
	fq4, err := newFq4(fq2, nil)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq4.calculateFrobeniusCoeffsWithPrecomputation(f1)

	a2, b2 := fq2.newElement(), fq2.newElement()
	twist, twist2, twist3 := fq2.zero(), fq2.newElement(), fq2.newElement()
	fq2.f.copy(twist[1], fq2.f.one)
	fq2.square(twist2, twist)
	fq2.mul(twist3, twist2, twist)
	fq2.mulByFq(a2, twist2, g1.a)
	fq2.mulByFq(b2, twist3, g1.b)
	g2, err := newG22(fq2, a2, b2, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	x, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return apiDecodingErr(err)
	}
	if isBigZero(x) {
		return apiDecodingErr(ERR_PAIRING_LOOP_COUNT_PARAM_ZERO)
	}

	if weight := calculateHammingWeight(x); weight > MAX_ATE_PAIRING_ATE_LOOP_COUNT_HAMMING {
		return apiDecodingErr(ERR_MNT_PAIRING_LOW_HAMMING_WEIGHT)
	}
	xIsNegative, rest, err := decodePairingExpSign(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	expW0, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_FINAL_EXP_W0_BIT_LENGTH)
	if err != nil {
		return apiDecodingErr(err)
	}
	if isBigZero(expW0) {
		return apiDecodingErr(ERR_MNT_EXPW0_NOT_ZERO)
	}
	expW1, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_FINAL_EXP_W1_BIT_LENGTH)
	if err != nil {
		return apiDecodingErr(err)
	}
	if isBigZero(expW1) {
		return apiDecodingErr(ERR_MNT_EXPW1_NOT_ZERO)
	}
	expW0IsNegative, rest, err := decodePairingExpSign(rest)
	if err != nil {
		return apiDecodingErr(ERR_MNT_INVALID_EXPW0)
	}
	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return apiDecodingErr(ERR_PAIRING_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_PAIRING_NUM_PAIRS_ZERO)
		}
	}

	var g1Points []*pointG1
	var g2Points []*pointG22
	for i := 0; i < numPairs; i++ {
		needG1SubGroupCheck, localRest, err := decodeBoolean(rest)
		if err != nil {
			return apiDecodingErr(err)
		}
		g1Point, localRest, err := decodeG1Point(localRest, modulusLen, g1)
		if err != nil {
			return apiDecodingErr(err)
		}
		needG2SubGroupCheck, localRest, err := decodeBoolean(localRest)
		if err != nil {
			return apiDecodingErr(err)
		}
		g2Point, localRest, err := decodeG22Point(localRest, modulusLen, g2)
		if err != nil {
			return apiDecodingErr(err)
		}
		if !g1.isOnCurve(g1Point) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_PAIRING_POINTG1_NOT_ON_CURVE)
			}
		}
		if !g2.isOnCurve(g2Point) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_PAIRING_POINTG2_NOT_ON_CURVE)
			}
		}
		if needG1SubGroupCheck {
			if ok := g1.checkCorrectSubGroup(g1Point); !ok {
				if !GAS_METERING_MODE {
					return apiExecErr(ERR_PAIRING_POINTG1_NOT_IN_SUBGROUP)
				}
			}
		}
		if needG2SubGroupCheck {
			if ok := g2.checkCorrectSubGroup(g2Point); !ok {
				if !GAS_METERING_MODE {
					return apiExecErr(ERR_PAIRING_POINTG2_NOT_IN_SUBGROUP)
				}
			}
		}
		if !g1.isZero(g1Point) && !g2.isZero(g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if len(g1Points) == 0 {
		return pairingSuccess, nil
	}

	engine := newMnt4Instance(
		x,
		xIsNegative,
		expW0,
		expW1,
		expW0IsNegative,
		fq4,
		g1,
		g2,
		twist,
	)
	result, hasValue := engine.multiPair(g1Points, g2Points)
	if !hasValue {
		return apiDecodingErr(ERR_PAIRING_NO_RETURN_VALUE)
	}
	if !fq4.equal(result, fq4.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}

func pairMNT6(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := decodeBaseFieldFromEncoding(in)
	if err != nil {
		return apiDecodingErr(err)
	}
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return apiDecodingErr(err)
	}
	_, order, rest, err := decodeGroupOrder(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	g1, err := newG1(field, a, b, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq3, rest, err := createExtension3FieldParams(rest, modulusLen, field, 6, false)
	if err != nil {
		return apiDecodingErr(err)
	}
	f1, err := constructBaseForFq3AndFq6(field, fq3.nonResidue)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq3.calculateFrobeniusCoeffsWithPrecomputation(f1)
	fq6, err := newFq6Quadratic(fq3, nil)
	if err != nil {
		return apiDecodingErr(err)
	}
	fq6.calculateFrobeniusCoeffsWithPrecomputation(f1)

	a3, b3 := fq3.newElement(), fq3.newElement()
	twist, twist2, twist3 := fq3.zero(), fq3.newElement(), fq3.newElement()
	fq3.f.copy(twist[1], fq3.f.one)
	fq3.square(twist2, twist)
	fq3.mul(twist3, twist2, twist)
	fq3.mulByFq(a3, twist2, g1.a)
	fq3.mulByFq(b3, twist3, g1.b)
	g2, err := newG23(fq3, a3, b3, order)
	if err != nil {
		return apiDecodingErr(err)
	}
	x, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return apiDecodingErr(err)
	}
	if isBigZero(x) {
		return apiDecodingErr(ERR_PAIRING_LOOP_COUNT_PARAM_ZERO)
	}
	if weight := calculateHammingWeight(x); weight > MAX_ATE_PAIRING_ATE_LOOP_COUNT_HAMMING {
		return apiDecodingErr(ERR_MNT_PAIRING_LOW_HAMMING_WEIGHT)
	}
	xIsNegative, rest, err := decodePairingExpSign(rest)
	if err != nil {
		return apiDecodingErr(err)
	}
	expW0, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return apiDecodingErr(err)
	}
	if isBigZero(expW0) {
		return apiDecodingErr(ERR_MNT_EXPW0_NOT_ZERO)
	}
	expW1, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return apiDecodingErr(err)
	}
	if isBigZero(expW1) {
		return apiDecodingErr(ERR_MNT_EXPW1_NOT_ZERO)
	}
	expW0IsNegative, rest, err := decodePairingExpSign(rest)
	if err != nil {
		return apiDecodingErr(ERR_MNT_INVALID_EXPW0)
	}
	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return apiDecodingErr(ERR_PAIRING_NUM_PAIRS_NOT_ENOUGH_BYTE)
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		if !GAS_METERING_MODE {
			return apiExecErr(ERR_PAIRING_NUM_PAIRS_ZERO)
		}
	}

	var g1Points []*pointG1
	var g2Points []*pointG23
	for i := 0; i < numPairs; i++ {
		needG1SubGroupCheck, localRest, err := decodeBoolean(rest)
		if err != nil {
			return apiDecodingErr(err)
		}
		g1Point, localRest, err := decodeG1Point(localRest, modulusLen, g1)
		if err != nil {
			return apiDecodingErr(err)
		}
		needG2SubGroupCheck, localRest, err := decodeBoolean(localRest)
		if err != nil {
			return apiDecodingErr(err)
		}
		g2Point, localRest, err := decodeG23Point(localRest, modulusLen, g2)
		if err != nil {
			return apiDecodingErr(err)
		}
		if !g1.isOnCurve(g1Point) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_PAIRING_POINTG1_NOT_ON_CURVE)
			}
		}
		if !g2.isOnCurve(g2Point) {
			if !GAS_METERING_MODE {
				return apiExecErr(ERR_PAIRING_POINTG2_NOT_ON_CURVE)
			}
		}
		if needG1SubGroupCheck {
			if ok := g1.checkCorrectSubGroup(g1Point); !ok {
				if !GAS_METERING_MODE {
					return apiExecErr(ERR_PAIRING_POINTG1_NOT_IN_SUBGROUP)
				}
			}
		}
		if needG2SubGroupCheck {
			if ok := g2.checkCorrectSubGroup(g2Point); !ok {
				if !GAS_METERING_MODE {
					return apiExecErr(ERR_PAIRING_POINTG2_NOT_IN_SUBGROUP)
				}
			}
		}
		if !g1.isZero(g1Point) && !g2.isZero(g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return apiDecodingErr(ERR_GARBAGE_INPUT)
	}
	if len(g1Points) == 0 {
		return pairingSuccess, nil
	}

	engine := newMNT6Instance(
		x,
		xIsNegative,
		expW0,
		expW1,
		expW0IsNegative,
		fq6,
		g1,
		g2,
		twist,
	)
	result, hasValue := engine.multiPair(g1Points, g2Points)
	if !hasValue {
		return apiDecodingErr(ERR_PAIRING_NO_RETURN_VALUE)
	}
	if !fq6.equal(result, fq6.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}
