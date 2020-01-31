package eip

import (
	"errors"
	"math/big"
)

var (
	zero                         = []byte{0x00}
	pairingError, pairingSuccess = []byte{0x00}, []byte{0x01}
	OPERATION_G1_ADD             = 0x01
	OPERATION_G1_MUL             = 0x02
	OPERATION_G1_MULTIEXP        = 0x03
	OPERATION_G2_ADD             = 0x04
	OPERATION_G2_MUL             = 0x05
	OPERATION_G2_MULTIEXP        = 0x06
	OPERATION_BLS12PAIR          = 0x07
	OPERATION_BNPAIR             = 0x08
	OPERATION_MNT4PAIR           = 0x09
	OPERATION_MNT6PAIR           = 0x0a
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
		return zero, errors.New("Unknown operation type")
	}
}

type g1Api struct{}

func (api *g1Api) addPoints(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return nil, err
	}
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return nil, err
	}
	_, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}
	g1, err := newG1(field, a, b, order.Bytes())
	if err != nil {
		return nil, err
	}
	p0, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return nil, err
	}
	p1, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("Input contains garbage at the end")
	}
	if !g1.isOnCurve(p0) {
		return nil, errors.New("point 0 isn't on the curve")
	}
	if !g1.isOnCurve(p1) {
		return nil, errors.New("point 1 isn't on the curve")
	}
	g1.add(p0, p0, p1)
	out := make([]byte, 2*modulusLen)
	encodeG1Point(out, g1.toBytes(p0))
	return out, nil
}

func (api *g1Api) mulPoint(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return nil, err
	}

	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return nil, err
	}
	orderLen, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}

	g1, err := newG1(field, a, b, order.Bytes())
	if err != nil {
		return nil, err
	}
	p, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return nil, err
	}
	if !g1.isOnCurve(p) {
		return nil, errors.New("point isn't on the curve")
	}
	s, rest, err := decodeScalar(rest, orderLen, order)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("Input contains garbage at the end")
	}
	g1.mulScalar(p, p, s)
	out := make([]byte, 2*modulusLen)
	encodeG1Point(out, g1.toBytes(p))
	return out, nil
}

func (api *g1Api) multiExp(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return nil, err
	}
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return nil, err
	}
	orderLen, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}
	g1, err := newG1(field, a, b, order.Bytes())
	if err != nil {
		return nil, err
	}
	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get number of pairs")
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return pairingError, errors.New("zero pairs encoded")
	}
	if len(rest) != (2*modulusLen+orderLen)*numPairs {
		return nil, errors.New("Input length is invalid for number of pairs")
	}
	bases := make([]*pointG1, numPairs)
	scalars := make([]*big.Int, numPairs)
	for i := 0; i < numPairs; i++ {
		p, localRest, err := decodeG1Point(rest, modulusLen, g1)
		if err != nil {
			return pairingError, err
		}
		if !g1.isOnCurve(p) {
			return nil, errors.New("point isn't on the curve")
		}
		scalar, localRest, err := decodeScalar(localRest, orderLen, order)
		if err != nil {
			return nil, err
		}
		bases[i], scalars[i] = g1.newPoint(), new(big.Int)
		g1.copy(bases[i], p)
		scalars[i].Set(scalar)
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}

	p := g1.newPoint()
	if len(bases) != len(scalars) || len(bases) == 0 {
		g1.copy(p, g1.zero())
	} else {
		g1.multiExp(p, bases, scalars)
	}
	out := make([]byte, 2*modulusLen)
	encodeG1Point(out, g1.toBytes(p))
	return out, nil
}

type g2Api struct{}

func (api *g2Api) run(opType int, in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return nil, err
	}
	degreeBuf, rest, err := split(rest, EXTENSION_DEGREE_LENGTH_ENCODING)
	if err != nil {
		return nil, errors.New("cant decode extension degree length")
	}
	degree := int(degreeBuf[0])
	// fmt.Printf("[debug] ext degree: %d\n", degree)
	switch degree {
	case EXTENSION_TWO_DEGREE:
		return new(g22Api).run(opType, field, modulusLen, rest)
	case EXTENSION_THREE_DEGREE:
		return new(g23Api).run(opType, field, modulusLen, rest)
	default:
		return nil, errors.New("Extension degree expected to be 2 or 3")
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
		return nil, errors.New("Unknown g22 operation")
	}
}

func (api *g22Api) addPoints(field *field, modulusLen int, in []byte) ([]byte, error) {
	fq2, rest, err := createExtension2FieldParams(in, modulusLen, field, false)
	if err != nil {
		return nil, err
	}
	a2, b2, rest, err := decodeBAInExtField2FromEncoding(rest, modulusLen, fq2)
	if err != nil {
		return nil, err
	}
	_, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}
	g2, err := newG22(fq2, a2, b2, order.Bytes())
	if err != nil {
		return nil, err
	}
	q0, rest, err := decodeG22Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}
	q1, rest, err := decodeG22Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}

	if len(rest) != 0 {
		return nil, errors.New("Input contains garbage at the end")
	}
	if !g2.isOnCurve(q0) {
		return nil, errors.New("point 0 isn't on the curve")
	}
	if !g2.isOnCurve(q1) {
		return nil, errors.New("point 1 isn't on the curve")
	}

	g2.add(q0, q0, q1)
	out := make([]byte, 4*modulusLen)
	encodeG22Point(out, g2.toBytes(q0))
	return out, nil
}

func (api *g22Api) mulPoint(field *field, modulusLen int, in []byte) ([]byte, error) {
	fq2, rest, err := createExtension2FieldParams(in, modulusLen, field, false)
	if err != nil {
		return nil, err
	}
	a2, b2, rest, err := decodeBAInExtField2FromEncoding(rest, modulusLen, fq2)
	if err != nil {
		return nil, err
	}
	orderLen, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}
	g2, err := newG22(fq2, a2, b2, order.Bytes())
	if err != nil {
		return nil, err
	}
	fq2.copy(g2.a, a2)
	fq2.copy(g2.b, b2)
	q, rest, err := decodeG22Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}
	s, rest, err := decodeScalar(rest, orderLen, order)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("Input contains garbage at the end")
	}
	if !g2.isOnCurve(q) {
		return nil, errors.New("q1 isn't on the curve")
	}
	g2.mulScalar(q, q, s)
	out := make([]byte, 4*modulusLen)
	encodeG22Point(out, g2.toBytes(q))
	return out, nil
}

func (api *g22Api) multiExp(field *field, modulusLen int, in []byte) ([]byte, error) {
	fq2, rest, err := createExtension2FieldParams(in, modulusLen, field, false)
	if err != nil {
		return nil, err
	}
	a2, b2, rest, err := decodeBAInExtField2FromEncoding(rest, modulusLen, fq2)
	if err != nil {
		return nil, err
	}
	orderLen, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}
	g2, err := newG22(fq2, a2, b2, order.Bytes())
	if err != nil {
		return nil, err
	}

	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return pairingError, errors.New("Invalid number of pairs")
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return pairingError, errors.New("zero pairs encoded")
	}
	if len(rest) != (4*modulusLen+orderLen)*numPairs {
		return nil, errors.New("Input length is invalid for number of pairs")
	}
	bases := make([]*pointG22, numPairs)
	scalars := make([]*big.Int, numPairs)
	for i := 0; i < numPairs; i++ {
		q, localRest, err := decodeG22Point(rest, modulusLen, g2)
		if err != nil {
			return pairingError, err
		}
		if !g2.isOnCurve(q) {
			return nil, errors.New("point isn't on the curve")
		}
		scalar, localRest, err := decodeScalar(localRest, orderLen, order)
		if err != nil {
			return nil, err
		}
		bases[i], scalars[i] = g2.newPoint(), new(big.Int)
		g2.copy(bases[i], q)
		scalars[i].Set(scalar)
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}

	q := g2.newPoint()
	if len(bases) != len(scalars) || len(bases) == 0 {
		g2.copy(q, g2.zero())
	} else {
		g2.multiExp(q, bases, scalars)
	}

	out := make([]byte, 4*modulusLen)
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
		return nil, errors.New("Unknown g23 operation")
	}
}

func (api *g23Api) addPoints(field *field, modulusLen int, in []byte) ([]byte, error) {
	fq3, rest, err := createExtension3FieldParams(in, modulusLen, field, false)
	if err != nil {
		return nil, err
	}
	a2, b2, rest, err := decodeBAInExtField3FromEncoding(rest, modulusLen, fq3)
	if err != nil {
		return nil, err
	}
	_, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}
	g2, err := newG23(fq3, a2, b2, order.Bytes())
	if err != nil {
		return nil, err
	}
	q0, rest, err := decodeG23Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}
	q1, rest, err := decodeG23Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("Input contains garbage at the end")
	}
	if !g2.isOnCurve(q0) {
		return nil, errors.New("point 0 isn't on the curve")
	}
	if !g2.isOnCurve(q1) {
		return nil, errors.New("point 1 isn't on the curve")
	}
	g2.add(q0, q0, q1)
	out := make([]byte, 6*modulusLen)
	encodeG23Point(out, g2.toBytes(q0))
	return out, nil
}

func (api *g23Api) mulPoint(field *field, modulusLen int, in []byte) ([]byte, error) {
	fq3, rest, err := createExtension3FieldParams(in, modulusLen, field, false)
	if err != nil {
		return nil, err
	}
	a3, b3, rest, err := decodeBAInExtField3FromEncoding(rest, modulusLen, fq3)
	if err != nil {
		return nil, err
	}
	orderLen, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}
	g2, err := newG23(fq3, a3, b3, order.Bytes())
	if err != nil {
		return nil, err
	}
	q, rest, err := decodeG23Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}
	s, rest, err := decodeScalar(rest, orderLen, order)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, errors.New("Input contains garbage at the end")
	}
	if !g2.isOnCurve(q) {
		return nil, errors.New("point isn't on the curve")
	}
	g2.mulScalar(q, q, s)
	out := make([]byte, 6*modulusLen)
	encodeG23Point(out, g2.toBytes(q))
	return out, nil
}

func (api *g23Api) multiExp(field *field, modulusLen int, in []byte) ([]byte, error) {
	fq3, rest, err := createExtension3FieldParams(in, modulusLen, field, false)
	if err != nil {
		return nil, err
	}
	a2, b2, rest, err := decodeBAInExtField3FromEncoding(rest, modulusLen, fq3)
	if err != nil {
		return nil, err
	}
	orderLen, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return nil, err
	}
	g2, err := newG23(fq3, a2, b2, order.Bytes())
	if err != nil {
		return nil, err
	}
	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get number of pairs")
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return pairingError, errors.New("zero pairs encoded")
	}
	// TODO: make handling 6*modulusLen generic
	if len(rest) != (6*modulusLen+orderLen)*numPairs {
		return nil, errors.New("Input length is invalid for number of pairs")
	}
	bases := make([]*pointG23, numPairs)
	scalars := make([]*big.Int, numPairs)
	for i := 0; i < numPairs; i++ {
		q, localRest, err := decodeG23Point(rest, modulusLen, g2)
		if err != nil {
			return pairingError, err
		}
		if !g2.isOnCurve(q) {
			return nil, errors.New("point isn't on the curve")
		}
		scalar, localRest, err := decodeScalar(localRest, orderLen, order)
		if err != nil {
			return nil, err
		}
		bases[i], scalars[i] = g2.newPoint(), new(big.Int)
		g2.copy(bases[i], q)
		scalars[i].Set(scalar)
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}

	q := g2.newPoint()
	if len(bases) != len(scalars) || len(bases) == 0 {
		g2.copy(q, g2.zero())
	} else {
		g2.multiExp(q, bases, scalars)
	}
	out := make([]byte, 6*modulusLen)
	encodeG23Point(out, g2.toBytes(q))
	return out, nil
}

func pairBN(in []byte) ([]byte, error) {
	// base field
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return pairingError, err
	}
	// g1
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return pairingError, err
	}
	if !field.isZero(a) {
		return pairingError, errors.New("A parameter must be zero for BN curve")
	}
	_, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return pairingError, err
	}
	g1, err := newG1(field, a, b, order.Bytes())
	if err != nil {
		return pairingError, err
	}
	// ext2
	fq2, rest, err := createExtension2FieldParamsForPairing(rest, modulusLen, field, true, 2)
	if err != nil {
		return pairingError, err
	}

	fq2NonResidue, rest, err := decodeFp2(rest, modulusLen, fq2)
	if err != nil {
		return pairingError, err
	}
	if fq2.isZero(fq2NonResidue) {
		return pairingError, errors.New("Non-residue for Fp6 is zero")
	}
	if !isNonNThRootFp2(fq2, fq2NonResidue, 6) {
		return pairingError, errors.New("Non-residue for Fp6 is actually a residue")
	}
	// twist type 0x01: M, 0x02: D
	twistTypeBuf, rest, err := split(rest, TWIST_TYPE_LENGTH)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get twist type")
	}
	twistType := twistTypeBuf[0]
	if twistType != 0x01 && twistType != 0x02 {
		return pairingError, errors.New("Unknown twist type supplied")
	}

	f1, f2, err := constructBaseForFq6AndFq12(fq2, fq2NonResidue)
	if err != nil {
		return pairingError, errors.New("Can not make base precomputations for Fp6/Fp12 frobenius")
	}

	// ext6
	fq6, err := newFq6(fq2, nil)
	if err != nil {
		return pairingError, err
	}
	fq2.copy(fq6.nonResidue, fq2NonResidue)
	if ok := fq6.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return pairingError, errors.New("Can not calculate Frobenius coefficients for Fp6")
	}
	// ext12
	fq12, err := newFq12(fq6, nil)
	if err != nil {
		return pairingError, err
	}
	if ok := fq12.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return pairingError, errors.New("Can not calculate Frobenius coefficients for Fp12")
	}

	b2 := fq2.newElement()
	if twistType == 0x01 {
		fq2.mulByFq(b2, fq6.nonResidue, b)
	} else {
		fq6NonResidueInv := fq2.newElement()
		fq2.inverse(fq6NonResidueInv, fq6.nonResidue)
		fq2.mulByFq(b2, fq6NonResidueInv, b)
	}

	// g2
	g2, err := newG22(fq2, fq2.zero(), b2, order.Bytes())
	if err != nil {
		return pairingError, err
	}
	// u
	u, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return pairingError, err
	}
	if isBigZero(u) {
		return pairingError, errors.New("Loop count parameters can not be zero")
	}
	// u is negative
	uIsNegativeBuf, rest, err := split(rest, SIGN_ENCODING_LENGTH)
	if err != nil {
		return pairingError, errors.New("X is not encoded properly")
	}
	// maybe better? uIsNegativeBuf[0 : SIGN_ENCODING_LENGTH-1]
	var (
		uIsNegative bool
		sixUPlus2   *big.Int
	)

	six, two := big.NewInt(6), big.NewInt(2)
	switch uIsNegativeBuf[0] {
	case 0x01:
		uIsNegative = true
		sixUPlus2 = new(big.Int).Mul(six, u)
		sixUPlus2 = new(big.Int).Sub(sixUPlus2, two)
		break
	case 0x00:
		uIsNegative = false
		sixUPlus2 = new(big.Int).Mul(six, u)
		sixUPlus2 = new(big.Int).Add(sixUPlus2, two)
		break
	default:
		return pairingError, errors.New("Unknown parameter u sign")
	}

	if weight := calculateHammingWeight(sixUPlus2); weight > MAX_BN_SIX_U_PLUS_TWO_HAMMING {
		return pairingError, errors.New("|6*U + 2| has too large hamming weight")
	}
	minus2Inv := new(big.Int).ModInverse(big.NewInt(-2), field.pbig)
	nonResidueInPMinus1Over2 := fq2.newElement()
	fq2.exp(nonResidueInPMinus1Over2, fq6.nonResidue, minus2Inv)

	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get number of pairs")
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return pairingError, errors.New("zero pairs encoded")
	}

	var g1Points []*pointG1
	var g2Points []*pointG22
	for i := 0; i < numPairs; i++ {
		g1Point, localRest, err := decodeG1Point(rest, modulusLen, g1)
		if err != nil {
			return pairingError, err
		}
		g2Point, localRest, err := decodeG22Point(localRest, modulusLen, g2)
		if err != nil {
			return pairingError, err
		}
		if !g1.isOnCurve(g1Point) {
			return pairingError, errors.New("G1 point is not on curve")
		}
		if !g2.isOnCurve(g2Point) {
			return pairingError, errors.New("G2 point is not on curve")
		}
		if ok := g1.checkCorrectSubGroup(g1Point); !ok {
			return pairingError, errors.New("G1 point is not in the expected subgroup")
		}
		if ok := g2.checkCorrectSubGroup(g2Point); !ok {
			return pairingError, errors.New("G2 point is not in the expected subgroup")
		}
		if !g1.isZero(g1Point) && !g2.isZero(g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}
	if len(g1Points) == 0 {
		return pairingError, nil
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
	result, err := engine.multiPair(g1Points, g2Points)
	if err != nil {
		return pairingError, nil
	}
	if !fq12.equal(result, fq12.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}

func pairBLS(in []byte) ([]byte, error) {
	// base field
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return pairingError, err
	}

	// g1
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return pairingError, err
	}
	if !field.isZero(a) {
		return pairingError, errors.New("A parameter must be zero for BLS12 curve")
	}
	_, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return pairingError, err
	}
	g1, err := newG1(field, a, b, order.Bytes())
	if err != nil {
		return pairingError, err
	}
	fq2, rest, err := createExtension2FieldParamsForPairing(rest, modulusLen, field, true, 2)
	if err != nil {
		return pairingError, err
	}

	fq2NonResidue, rest, err := decodeFp2(rest, modulusLen, fq2)
	if err != nil {
		return pairingError, err
	}
	if fq2.isZero(fq2NonResidue) {
		return pairingError, errors.New("Non-residue for Fp6 is zero")
	}
	if !isNonNThRootFp2(fq2, fq2NonResidue, 6) {
		return pairingError, errors.New("Non-residue for Fp6 is actually a residue")
	}
	// twist type 0x01: M, 0x02: D
	twistTypeBuf, rest, err := split(rest, TWIST_TYPE_LENGTH)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get twist type")
	}
	twistType := twistTypeBuf[0]
	if twistType != 0x01 && twistType != 0x02 {
		return pairingError, errors.New("Unknown twist type supplied")
	}

	f1, f2, err := constructBaseForFq6AndFq12(fq2, fq2NonResidue)
	if err != nil {
		return pairingError, errors.New("Can not make base precomputations for Fp6/Fp12 frobenius")
	}

	// ext6
	fq6, err := newFq6(fq2, nil)
	if err != nil {
		return pairingError, err
	}
	fq2.copy(fq6.nonResidue, fq2NonResidue)
	if ok := fq6.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return pairingError, errors.New("Can not calculate Frobenius coefficients for Fp6")
	}
	// ext12
	fq12, err := newFq12(fq6, nil)
	if err != nil {
		return pairingError, err
	}
	if ok := fq12.calculateFrobeniusCoeffsWithPrecomputation(f1, f2); !ok {
		return pairingError, errors.New("Can not calculate Frobenius coefficients for Fp12")
	}

	b2 := fq2.newElement()
	if twistType == 0x01 {
		fq2.mulByFq(b2, fq6.nonResidue, b)
	} else {
		fq6NonResidueInv := fq2.newElement()
		fq2.inverse(fq6NonResidueInv, fq6.nonResidue)
		fq2.mulByFq(b2, fq6NonResidueInv, b)
	}
	// g2
	g2, err := newG22(fq2, fq2.zero(), b2, order.Bytes())
	if err != nil {
		return pairingError, err
	}

	// z
	z, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return pairingError, err
	}
	if z.Cmp(big.NewInt(0)) == 0 {
		return pairingError, errors.New("Loop count parameters can not be zero")
	}
	// u is negative
	zIsNegativeBuf, rest, err := split(rest, SIGN_ENCODING_LENGTH)
	if err != nil {
		return pairingError, errors.New("z is not encoded properly")
	}
	// maybe better? uIsNegativeBuf[0 : SIGN_ENCODING_LENGTH-1]
	var zIsNegative bool
	switch zIsNegativeBuf[0] {
	case 0x01:
		zIsNegative = true
		break
	case 0x00:
		zIsNegative = false
		break
	default:
		return pairingError, errors.New("Unknown parameter z sign")
	}
	if weight := calculateHammingWeight(z); weight > MAX_BLS12_X_HAMMING {
		return pairingError, errors.New("z has too large hamming weight")
	}

	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get number of pairs")
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return pairingError, errors.New("zero pairs encoded")
	}
	var g1Points []*pointG1
	var g2Points []*pointG22
	for i := 0; i < numPairs; i++ {
		g1Point, localRest, err := decodeG1Point(rest, modulusLen, g1)
		if err != nil {
			return pairingError, err
		}
		g2Point, localRest, err := decodeG22Point(localRest, modulusLen, g2)
		if err != nil {
			return pairingError, err
		}
		if !g1.isOnCurve(g1Point) {
			return pairingError, errors.New("G1 point is not on curve")
		}
		if !g2.isOnCurve(g2Point) {
			return pairingError, errors.New("G2 point is not on curve")
		}

		if ok := g1.checkCorrectSubGroup(g1Point); !ok {
			return pairingError, errors.New("G1 point is not in the expected subgroup")
		}

		if ok := g2.checkCorrectSubGroup(g2Point); !ok {
			return pairingError, errors.New("G2 point is not in the expected subgroup")
		}
		if !g1.isZero(g1Point) && !g2.isZero(g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}
	if len(g1Points) == 0 {
		return pairingError, nil
	}

	// pairs
	engine := newBLSInstance(
		z,
		zIsNegative,
		twistType,
		g1,
		g2,
		fq12,
		false,
	)
	result, err := engine.multiPair(g1Points, g2Points)
	if err != nil {
		return pairingError, errors.New("Pairing engine returned no value")
	}
	if !fq12.equal(result, fq12.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}

func pairMNT4(in []byte) ([]byte, error) {
	// base field
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return pairingError, err
	}
	// g1
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return pairingError, err
	}
	_, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return pairingError, err
	}
	g1, err := newG1(field, a, b, order.Bytes())
	if err != nil {
		return pairingError, err
	}
	// ext2
	fq2, rest, err := createExtension2FieldParamsForPairing(rest, modulusLen, field, false, 4)
	if err != nil {
		return pairingError, err
	}

	f1 := constructBaseForFq2AndFq4(field, fq2.nonResidue)
	fq2.calculateFrobeniusCoeffsWithPrecomputation(f1)

	// ext4
	fq4, err := newFq4(fq2, nil)
	if err != nil {
		return pairingError, err
	}
	fq4.calculateFrobeniusCoeffsWithPrecomputation(f1)

	a2, b2 := fq2.newElement(), fq2.newElement()
	twist, twist2, twist3 := fq2.zero(), fq2.newElement(), fq2.newElement()
	fq2.f.copy(twist[1], fq2.f.one)
	fq2.square(twist2, twist)
	fq2.mul(twist3, twist2, twist)
	fq2.mulByFq(a2, twist2, g1.a)
	fq2.mulByFq(b2, twist3, g1.b)

	// g2
	g2, err := newG22(fq2, a2, b2, order.Bytes())
	if err != nil {
		return pairingError, err
	}
	// x
	x, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return pairingError, err
	}
	if isBigZero(x) {
		return pairingError, errors.New("Ate pairing loop count parameters can not be zero")
	}

	if weight := calculateHammingWeight(x); weight > MAX_ATE_PAIRING_ATE_LOOP_COUNT_HAMMING {
		return pairingError, errors.New("X has too large hamming weight")
	}

	// x is negative
	xIsNegativeBuf, rest, err := split(rest, SIGN_ENCODING_LENGTH)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get X sign encoding")
	}
	// maybe better? uIsNegativeBuf[0 : SIGN_ENCODING_LENGTH-1]
	var xIsNegative bool
	switch xIsNegativeBuf[0] {
	case 0x01:
		xIsNegative = true
		break
	case 0x00:
		xIsNegative = false
		break
	default:
		return pairingError, errors.New("X sign is not encoded properly")
	}

	// expW0
	expW0, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_FINAL_EXP_W0_BIT_LENGTH)
	if err != nil {
		return pairingError, err
	}
	if isBigZero(expW0) {
		return pairingError, errors.New("Final exp w0 loop count parameters can not be zero")
	}
	// expW1
	expW1, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_FINAL_EXP_W1_BIT_LENGTH)
	if err != nil {
		return pairingError, err
	}
	if isBigZero(expW1) {
		return pairingError, errors.New("Final exp w1 loop count parameters can not be zero")
	}
	expW0IsNegativeBuf, rest, err := split(rest, SIGN_ENCODING_LENGTH)
	if err != nil {
		return pairingError, errors.New("Exp_w0 sign is not encoded properly")
	}
	var expW0IsNegative bool
	switch expW0IsNegativeBuf[0] {
	case 0x01:
		expW0IsNegative = true
		break
	case 0x00:
		expW0IsNegative = false
		break
	default:
		return pairingError, errors.New("Unknown expW0 sign")
	}

	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get number of pairs")
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return pairingError, errors.New("Zero pairs encoded")
	}

	var g1Points []*pointG1
	var g2Points []*pointG22
	for i := 0; i < numPairs; i++ {
		g1Point, localRest, err := decodeG1Point(rest, modulusLen, g1)
		if err != nil {
			return pairingError, err
		}
		g2Point, localRest, err := decodeG22Point(localRest, modulusLen, g2)
		if err != nil {
			return pairingError, err
		}
		if !g1.isOnCurve(g1Point) {
			return pairingError, errors.New("G1 point is not on curve")
		}
		if !g2.isOnCurve(g2Point) {
			return pairingError, errors.New("G2 point is not on curve")
		}
		if ok := g1.checkCorrectSubGroup(g1Point); !ok {
			return pairingError, errors.New("G1 point is not in the expected subgroup")
		}
		if ok := g2.checkCorrectSubGroup(g2Point); !ok {
			return pairingError, errors.New("G2 point is not in the expected subgroup")
		}
		if !g1.isZero(g1Point) && !g2.isZero(g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}
	if len(g1Points) == 0 {
		return pairingError, nil
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
	result, err := engine.multiPair(g1Points, g2Points)
	if err != nil {
		return pairingError, nil
	}
	if !fq4.equal(result, fq4.one()) {
		return pairingError, nil
	}

	return pairingSuccess, nil
}

func pairMNT6(in []byte) ([]byte, error) {
	// base field
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return pairingError, err
	}

	// g1
	a, b, rest, err := decodeBAInBaseFieldFromEncoding(rest, modulusLen, field)
	if err != nil {
		return pairingError, err
	}
	_, order, rest, err := parseGroupOrder(rest)
	if err != nil {
		return pairingError, err
	}
	g1, err := newG1(field, a, b, order.Bytes())
	if err != nil {
		return pairingError, err
	}

	// ext3
	fq3, rest, err := createExtension3FieldParamsForPairing(rest, modulusLen, field, false, 6)
	if err != nil {
		return pairingError, err
	}
	f1, err := constructBaseForFq3AndFq6(field, fq3.nonResidue)
	if err != nil {
		return pairingError, err
	}
	fq3.calculateFrobeniusCoeffsWithPrecomputation(f1)

	// ext6 quadratic
	fq6, err := newFq6Quadratic(fq3, nil)
	if err != nil {
		return pairingError, err
	}
	fq6.calculateFrobeniusCoeffsWithPrecomputation(f1)

	a3, b3 := fq3.newElement(), fq3.newElement()
	twist, twist2, twist3 := fq3.zero(), fq3.newElement(), fq3.newElement()
	fq3.f.copy(twist[1], fq3.f.one)
	fq3.square(twist2, twist)
	fq3.mul(twist3, twist2, twist)
	fq3.mulByFq(a3, twist2, g1.a)
	fq3.mulByFq(b3, twist3, g1.b)
	// g2
	g2, err := newG23(fq3, a3, b3, order.Bytes())
	if err != nil {
		return pairingError, err
	}
	// x
	x, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return pairingError, err
	}
	if isBigZero(x) {
		return pairingError, errors.New("Ate pairing loop count parameters can not be zero")
	}

	if weight := calculateHammingWeight(x); weight > MAX_ATE_PAIRING_ATE_LOOP_COUNT_HAMMING {
		return pairingError, errors.New("x has too large hamming weight")
	}

	// u is negative
	xIsNegativeBuf, rest, err := split(rest, SIGN_ENCODING_LENGTH)
	if err != nil {
		return pairingError, errors.New("x is not encoded properly")
	}
	// maybe better? uIsNegativeBuf[0 : SIGN_ENCODING_LENGTH-1]
	var xIsNegative bool
	switch xIsNegativeBuf[0] {
	case 0x01:
		xIsNegative = true
		break
	case 0x00:
		xIsNegative = false
		break
	default:
		return pairingError, errors.New("Unknown parameter x sign")
	}

	// expW0
	expW0, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return pairingError, err
	}
	if isBigZero(expW0) {
		return pairingError, errors.New("Final exp w0 loop count parameters can not be zero")
	}
	// expW1
	expW1, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return pairingError, err
	}
	if isBigZero(expW1) {
		return pairingError, errors.New("Final exp w1 loop count parameters can not be zero")
	}

	expW0IsNegativeBuf, rest, err := split(rest, SIGN_ENCODING_LENGTH)
	if err != nil {
		return pairingError, errors.New("Exp_w0 sign is not encoded properly")
	}
	var expW0IsNegative bool
	switch expW0IsNegativeBuf[0] {
	case 0x01:
		expW0IsNegative = true
		break
	case 0x00:
		expW0IsNegative = false
		break
	default:
		return pairingError, errors.New("Unknown expW0 sign")
	}

	numPairsBuf, rest, err := split(rest, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return pairingError, errors.New("Input is not long enough to get number of pairs")
	}
	numPairs := int(numPairsBuf[0])
	if numPairs == 0 {
		return pairingError, errors.New("zero pairs encoded")
	}

	var g1Points []*pointG1
	var g2Points []*pointG23
	for i := 0; i < numPairs; i++ {
		g1Point, localRest, err := decodeG1Point(rest, modulusLen, g1)
		if err != nil {
			return pairingError, err
		}
		g2Point, localRest, err := decodeG23Point(localRest, modulusLen, g2)
		if err != nil {
			return pairingError, err
		}
		if !g1.isOnCurve(g1Point) {
			return pairingError, errors.New("G1 point is not on curve")
		}
		if !g2.isOnCurve(g2Point) {
			return pairingError, errors.New("G2 point is not on curve")
		}

		if ok := g1.checkCorrectSubGroup(g1Point); !ok {
			return pairingError, errors.New("G1 point is not in the expected subgroup")
		}

		if ok := g2.checkCorrectSubGroup(g2Point); !ok {
			return pairingError, errors.New("G2 point is not in the expected subgroup")
		}
		if !g1.isZero(g1Point) && !g2.isZero(g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}
	if len(g1Points) == 0 {
		return pairingError, nil
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
	result, err := engine.multiPair(g1Points, g2Points)
	if err != nil {
		return pairingError, nil
	}
	if !fq6.equal(result, fq6.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}
