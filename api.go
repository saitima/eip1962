package fp

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	zero                         = []byte{0x00}
	pairingError, pairingSuccess = []byte{0x00}, []byte{0x01}
)

type API struct{}

func (api *API) run(in []byte) ([]byte, error) {
	opTypeBuf, rest, err := split(in, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return zero, errors.New("Input should be longer than operation type encoding")
	}
	opType := opTypeBuf[0]
	switch opType {
	case 0x01:
		return new(g1Api).addPoints(rest)
	case 0x02:
		return new(g1Api).mulPoint(rest)
	case 0x03:
		return new(g1Api).multiExp(rest)
	case 0x04:
		return new(g2Api).addPoints(rest)
	case 0x05:
		return new(g2Api).mulPoint(rest)
	case 0x06:
		// TODO: multiexp
		return new(g2Api).mulPoint(rest)
	case 0x07:
		return new(pairingApi).pair(rest)
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
	_, order, rest, err := parseGroupOrder(rest, modulusLen)
	if err != nil {
		return nil, err
	}
	g1, err := newG1(field, nil, nil, order.Bytes())
	if err != nil {
		return nil, err
	}
	g1.a = a
	g1.b = b
	p0, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return nil, err
	}
	p1, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("Input contains garbage at the end")
	}
	if !g1.isOnCurve(p0) {
		return nil, fmt.Errorf("p0 isn't on the curve")
	}
	if !g1.isOnCurve(p1) {
		return nil, fmt.Errorf("p1 isn't on the curve")
	}
	g1.add(p1, p1, p0)
	out := g1.toBytes(p1)
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
	orderLen, order, rest, err := parseGroupOrder(rest, modulusLen)
	if err != nil {
		return nil, err
	}

	g1, err := newG1(field, nil, nil, order.Bytes())
	g1.f.cpy(g1.a, a)
	g1.f.cpy(g1.b, b)
	if err != nil {
		return nil, err
	}
	p, rest, err := decodeG1Point(rest, modulusLen, g1)
	if err != nil {
		return nil, err
	}
	s, rest, err := decodeScalar(rest, orderLen, order)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("Input contains garbage at the end")
	}
	g1.mulScalar(p, p, s)
	out := g1.toBytes(p)
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
	orderLen, order, rest, err := parseGroupOrder(rest, modulusLen)
	if err != nil {
		return nil, err
	}
	g1, err := newG1(field, nil, nil, order.Bytes())
	g1.f.cpy(g1.a, a)
	g1.f.cpy(g1.b, b)
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
		g1Point, localRest, err := decodeG1Point(rest, modulusLen, g1)
		if err != nil {
			return pairingError, err
		}
		scalar, localRest, err := decodeScalar(localRest, orderLen, order)
		g1.copy(bases[i], g1Point)
		scalars[i] = new(big.Int).Set(scalar)
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}
	if len(bases) != len(scalars) || len(bases) == 0 {
		return pairingSuccess, nil // success
	}
	p := g1.newPoint()
	g1.multiExp(p, bases, scalars)
	out := g1.toBytes(p)
	return out, nil
}

type g2Api struct{}

func (api *g2Api) addPoints(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return nil, err
	}
	fq2, rest, err := createExtension2FieldParams(rest, modulusLen, field, false)
	if err != nil {
		return nil, err
	}
	a2, b2, rest, err := decodeBAInExtField2FromEncoding(rest, modulusLen, fq2)
	if err != nil {
		return nil, err
	}
	_, order, rest, err := parseGroupOrder(rest, modulusLen)
	if err != nil {
		return nil, err
	}
	g2, err := newG22(fq2, nil, nil, order.Bytes())
	if err != nil {
		return nil, err
	}
	g2.a = a2
	g2.b = b2
	q0, rest, err := decodeG2Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}
	q1, rest, err := decodeG2Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("Input contains garbage at the end")
	}
	if !g2.isOnCurve(q0) {
		return nil, fmt.Errorf("p0 isn't on the curve")
	}
	if !g2.isOnCurve(q1) {
		return nil, fmt.Errorf("p1 isn't on the curve")
	}
	g2.add(q1, q1, q0)
	out := g2.toBytes(q1)
	return out, nil
}

func (api *g2Api) mulPoint(in []byte) ([]byte, error) {
	field, _, modulusLen, rest, err := parseBaseFieldFromEncoding(in)
	if err != nil {
		return nil, err
	}
	fq2, rest, err := createExtension2FieldParams(rest, modulusLen, field, false)
	if err != nil {
		return nil, err
	}
	a2, b2, rest, err := decodeBAInExtField2FromEncoding(rest, modulusLen, fq2)
	if err != nil {
		return nil, err
	}
	orderLen, order, rest, err := parseGroupOrder(rest, modulusLen)
	if err != nil {
		return nil, err
	}
	g2, err := newG22(fq2, nil, nil, order.Bytes())
	if err != nil {
		return nil, err
	}
	fq2.copy(g2.a, a2)
	fq2.copy(g2.b, b2)
	q, rest, err := decodeG2Point(rest, modulusLen, g2)
	if err != nil {
		return nil, err
	}
	s, rest, err := decodeScalar(rest, orderLen, order)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("Input contains garbage at the end")
	}
	g2.mulScalar(q, q, s)
	out := g2.toBytes(q)
	return out, nil
}

type pairingApi struct{}

func (api *pairingApi) pair(in []byte) ([]byte, error) {
	curveTypeBuf, rest, err := split(in, BYTES_FOR_LENGTH_ENCODING)
	if err != nil {
		return zero, errors.New("Input should be longer than curve type encoding")
	}
	curveType := curveTypeBuf[0]
	switch curveType {
	case 0x01:
		// bls
		return pairBN(rest)
	case 0x02:
		// bn
		return pairBN(rest)
	case 0x03:
		// mnt4
		return pairBN(rest)
	case 0x04:
		// mnt6
		return pairBN(rest)
	default:
		return zero, errors.New("Unknown curve type")
	}
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
	_, order, rest, err := parseGroupOrder(rest, modulusLen)
	if err != nil {
		return pairingError, err
	}
	g1, err := newG1(field, nil, nil, order.Bytes())
	g1.f.cpy(g1.a, a)
	g1.f.cpy(g1.b, b)
	if err != nil {
		return pairingError, err
	}
	// ext2
	nonResidue, rest, err := decodeFp(rest, modulusLen, field)
	if err != nil {
		return pairingError, err
	}
	if !isNonNThRoot(field, nonResidue, 2) {
		return pairingError, errors.New("Non-residue for Fp2 is actually a residue")
	}
	fq2, err := newFq2(field, nil)
	fq2.f.cpy(fq2.nonResidue, nonResidue)
	if err != nil {
		return pairingError, err
	}
	if ok := fq2.calculateFrobeniusCoeffs(); !ok {
		return pairingError, errors.New("Can not calculate Frobenius coefficients for Fp2")
	}
	fq2NonResidue, rest, err := decodeFp2(rest, modulusLen, fq2)
	if err != nil {
		return pairingError, err
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
	// ext6
	fq6, err := newFq6(fq2, nil)
	if err != nil {
		return pairingError, err
	}
	fq2.copy(fq6.nonResidue, fq2NonResidue)
	if ok := fq6.calculateFrobeniusCoeffs(); !ok {
		return pairingError, errors.New("Can not calculate Frobenius coefficients for Fp6")
	}
	// ext12
	fq12, err := newFq12(fq6, nil)
	if err != nil {
		return pairingError, err
	}
	if ok := fq12.calculateFrobeniusCoeffs(); !ok {
		return pairingError, errors.New("Can not calculate Frobenius coefficients for Fp12")
	}
	// g2
	g2, err := newG22(fq2, nil, nil, order.Bytes())
	if err != nil {
		return pairingError, err
	}
	// a2 is pairingError
	fq2.copy(g2.a, fq2.zero())
	if twistType == 0x01 {
		fq2.mulByFq(g2.b, fq6.nonResidue, b)
	} else {
		fq6NonResidueInv := fq2.newElement()
		fq2.inverse(fq6NonResidueInv, fq6.nonResidue)
		fq2.mulByFq(g2.b, fq6NonResidueInv, b)
	}

	// u
	u, rest, err := decodeLoopParameters(rest, MAX_ATE_PAIRING_ATE_LOOP_COUNT)
	if err != nil {
		return pairingError, err
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

	if weight := calculateHammingWeight(u); weight > MAX_BN_SIX_U_PLUS_TWO_HAMMING {
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
	g1zero, g2zero := g1.zero(), g2.zero()
	g1Tmp, g2Tmp := g1.newPoint(), g2.newPoint()
	for i := 0; i < numPairs; i++ {
		g1Point, localRest, err := decodeG1Point(rest, modulusLen, g1)
		if err != nil {
			return pairingError, err
		}
		g2Point, localRest, err := decodeG2Point(localRest, modulusLen, g2)
		if err != nil {
			return pairingError, err
		}
		g1.mulScalar(g1Tmp, g1Point, order)
		if !g1.equal(g1Tmp, g1zero) {
			return pairingError, errors.New("G1 point is not in the expected subgroup")
		}
		g2.mulScalar(g2Tmp, g2Point, order)
		if !g2.equal(g2Tmp, g2zero) {
			return pairingError, errors.New("G2 point is not in the expected subgroup")
		}
		if !g1.equal(g1zero, g1Point) && !g2.equal(g2zero, g2Point) {
			g1Points = append(g1Points, g1Point)
			g2Points = append(g2Points, g2Point)
		}
		rest = localRest
	}
	if len(rest) != 0 {
		return pairingError, errors.New("Input contains garbage at the end")
	}
	if len(g1Points) == 0 {
		return pairingSuccess, nil // success
	}

	// pairs
	engine := newBNInstance(
		u,
		sixUPlus2,
		uIsNegative,
		twistType,
		g1,
		g2,
		fq12,
		nonResidueInPMinus1Over2,
	)
	result := engine.multiPair(g1Points, g2Points)
	if !fq12.equal(result, fq12.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}
