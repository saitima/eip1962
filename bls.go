package eip

import (
	"math/big"
)

type blsInstance struct {
	z           *big.Int
	zIsnegative bool
	twistType   int
	g1          *g1
	g2          *g22
	fq12        *fq12
	t2          []*fe2
	t12         []*fe12
	preferNaf   bool
	zNaf        []int8
}

func newBLSInstance(z *big.Int, zIsnegative bool, twistType int, g1 *g1, g2 *g22, fq12 *fq12, forceNoNaf bool) *blsInstance {
	naf := ternaryWnaf(z)
	originalBits := onesCount(z)
	originalHamming := calculateHammingWeight(z)
	nafHamming := calculateNafHammingWeight(naf)
	var preferNaf bool
	if len(naf)+nafHamming < originalBits+originalHamming {
		preferNaf = true
		if forceNoNaf {
			preferNaf = false
		}
	} else {
		preferNaf = false
	}
	bls := &blsInstance{
		z:           z,
		zIsnegative: zIsnegative,
		twistType:   twistType,
		g1:          g1,
		g2:          g2,
		fq12:        fq12,
		preferNaf:   preferNaf,
		zNaf:        naf,
	}
	bls.t2 = make([]*fe2, 17)
	bls.t12 = make([]*fe12, 17)
	for i := 0; i < 17; i++ {
		bls.t2[i] = bls.fq12.f.f.new()
		bls.t12[i] = bls.fq12.new()
	}
	return bls
}

func (bls *blsInstance) gt() *fq12 {
	return bls.fq12
}

func (bls *blsInstance) doublingStep(coeff *fe6C, r *pointG22, twoInv fe) {
	fq2 := bls.fq12.f.f
	t := bls.t2
	// X*Y/2
	fq2.mul(t[0], r[0], r[1])
	fq2.mulByFq(t[0], t[0], twoInv)
	// Y^2
	fq2.square(t[1], r[1])
	// Z^2
	fq2.square(t[2], r[2])
	// 3*Z^2
	fq2.double(t[3], t[2])
	fq2.add(t[3], t[3], t[2])
	// 3*b*Z^2
	fq2.mul(t[4], t[3], bls.g2.b)
	// 9*b*Z^2
	fq2.double(t[5], t[4])
	fq2.add(t[5], t[5], t[4])
	// (Y^2 + 9*b*Z^2)/2
	fq2.add(t[6], t[1], t[5])
	fq2.mulByFq(t[6], t[6], twoInv)
	// (Y + Z)^2
	fq2.add(t[7], r[1], r[2])
	fq2.square(t[7], t[7])
	// Y^2 + Z^2
	fq2.add(t[8], t[1], t[2])
	// 2*Y*Z
	fq2.sub(t[9], t[7], t[8])
	// 3*b*Z^2 - Y^2
	fq2.sub(t[10], t[4], t[1])
	// X^2
	fq2.square(t[11], r[0])
	// (3*b*Z^2)^2
	fq2.square(t[12], t[4])
	// X = (Y^2 - 9*b*Z^2)*(X*Y/2)
	fq2.sub(r[0], t[1], t[5])
	fq2.mul(r[0], r[0], t[0])
	// 27*b^2*Z^4
	fq2.double(t[13], t[12])
	fq2.add(t[13], t[13], t[12])
	// Y = ((Y^2+9*b*Z^2)/2)^2 - 27*b^2*Z^4
	fq2.add(r[1], t[1], t[5])
	fq2.mulByFq(r[1], r[1], twoInv)
	fq2.square(r[1], r[1])
	fq2.sub(r[1], r[1], t[13])
	// Z = 2*Y^3*Z
	fq2.mul(r[2], t[9], t[1])
	// 3*X^2
	fq2.double(t[14], t[11])
	fq2.add(t[14], t[14], t[11])
	// -2*Y*Z
	fq2.neg(t[15], t[9])

	coeff[0] = *fq2.new()
	coeff[1] = *fq2.new()
	coeff[2] = *fq2.new()
	switch bls.twistType {
	case 1: // M
		fq2.copy(&coeff[0], t[10]) // 3*b*Z^2 - Y^2
		fq2.copy(&coeff[1], t[14]) // 3*X^2
		fq2.copy(&coeff[2], t[15]) // -2*Y*Z
		break
	case 2:
		fq2.copy(&coeff[0], t[15]) // -2*Y*Z
		fq2.copy(&coeff[1], t[14]) // 3*X^2
		fq2.copy(&coeff[2], t[10]) // 3*b*Z^2 - Y^2
		break
	}

}

func (bls *blsInstance) additionStep(coeff *fe6C, r *pointG22, q *pointG22) {
	fq2 := bls.fq12.f.f
	t := bls.t2
	// theta = Y - y*Z
	fq2.mul(t[0], q[1], r[2])
	fq2.sub(t[0], r[1], t[0])
	// lambda = X - x*Z
	fq2.mul(t[1], q[0], r[2])
	fq2.sub(t[1], r[0], t[1])
	// theta^2 = (Y - Y*Z)^2
	fq2.square(t[2], t[0])
	// lambda^2 = (X - X*Z)^2
	fq2.square(t[3], t[1])
	// lambda^3 = (X - X*Z)^3
	fq2.mul(t[4], t[3], t[1])
	// theta^2*Z = (Y - Y*Z)^2 * Z
	fq2.mul(t[5], t[2], r[2])
	// lambda^2*X = (X - X*Z)^2 * X
	fq2.mul(t[6], t[3], r[0])
	// H = lambda^3 + theta^2 * Z - 2*lambda^2 * X
	fq2.double(t[7], t[6])
	fq2.sub(t[7], t[5], t[7])
	fq2.add(t[7], t[7], t[4])
	// X = lambda * H
	fq2.mul(r[0], t[1], t[7])
	// (lambda^2*X - H)*theta
	fq2.sub(t[8], t[6], t[7])
	fq2.mul(t[8], t[8], t[0])
	// Y = (lambda^2*X - H)*theta - lambda^3*Y
	fq2.mul(t[9], t[4], r[1])
	fq2.sub(r[1], t[8], t[9])
	// Z = lambda^3*Z
	fq2.mul(r[2], t[4], r[2])
	// lambda*y
	fq2.mul(t[10], t[1], q[1])
	// theata*x - lambda*y
	fq2.mul(t[11], t[0], q[0])
	fq2.sub(t[11], t[11], t[10])
	// -theta
	fq2.neg(t[0], t[0])
	coeff[0] = *fq2.new()
	coeff[1] = *fq2.new()
	coeff[2] = *fq2.new()
	switch bls.twistType {
	case 1: // M
		fq2.copy(&coeff[0], t[11]) // theata*x - lambda*y
		fq2.copy(&coeff[1], t[0])  // -theta
		fq2.copy(&coeff[2], t[1])  // lambda
		break
	case 2: // D
		fq2.copy(&coeff[0], t[1])  // lambda
		fq2.copy(&coeff[1], t[0])  // -theta
		fq2.copy(&coeff[2], t[11]) // theata*x - lambda*y
		break
	}

}

func (bls *blsInstance) ell(f *fe12, coeffs *fe6C, p *pointG1) {
	fq2 := bls.fq12.f.f
	switch bls.twistType {
	case 1: // M
		fq2.mulByFq(&coeffs[2], &coeffs[2], p[1])
		fq2.mulByFq(&coeffs[1], &coeffs[1], p[0])
		bls.fq12.mulBy014(f, &coeffs[0], &coeffs[1], &coeffs[2])
	case 2: // D
		fq2.mulByFq(&coeffs[0], &coeffs[0], p[1])
		fq2.mulByFq(&coeffs[1], &coeffs[1], p[0])
		bls.fq12.mulBy034(f, &coeffs[0], &coeffs[1], &coeffs[2])
	}
}

func (bls *blsInstance) prepare(coeffs *[]fe6C, Q *pointG22) bool {
	f := bls.fq12.f.f.f
	twoInv := f.new()
	f.double(twoInv, f.one)
	if ok := f.inverse(twoInv, twoInv); !ok {
		return false
	}
	T := bls.g2.newPoint()
	bls.g2.copy(T, Q)

	if bls.preferNaf {
		bls.prepareWithNaf(coeffs, T, Q, twoInv)
	} else {
		bls.prepareWithoutNaf(coeffs, T, Q, twoInv)
	}

	if bls.zIsnegative {
		bls.g2.neg(T, T)
	}
	return true
}

func (bls *blsInstance) prepareWithNaf(coeffs *[]fe6C, T, Q *pointG22, twoInv fe) {
	j := 0
	for i := len(bls.zNaf) - 1; i >= 0; i-- {
		bls.doublingStep(&(*coeffs)[j], T, twoInv)
		j++
		if bls.zNaf[i] != 0 {
			bls.additionStep(&(*coeffs)[j], T, Q)
			j++
		}
	}
}

func (bls *blsInstance) prepareWithoutNaf(coeffs *[]fe6C, T, Q *pointG22, twoInv fe) {
	j := 0
	//  skip first msb bit
	for i := bls.z.BitLen() - 2; i >= 0; i-- {
		bls.doublingStep(&(*coeffs)[j], T, twoInv)
		j++
		if bls.z.Bit(int(i)) != 0 {
			bls.additionStep(&(*coeffs)[j], T, Q)
			j++
		}
	}
}

func (bls *blsInstance) millerLoop(f *fe12, g1Points []*pointG1, g2Points []*pointG22) bool {
	coeffs := make([][]fe6C, len(g1Points))
	coeffsLen := bls.calculateCoeffLength()
	for i := 0; i < len(g1Points); i++ {
		coeffs[i] = make([]fe6C, coeffsLen)
		if ok := bls.prepare(&coeffs[i], g2Points[i]); !ok {
			return false
		}
	}
	if bls.preferNaf {
		bls.millerLoopWithNaf(f, coeffs, g1Points)
	} else {
		bls.millerLoopWithoutNaf(f, coeffs, g1Points)
	}
	if bls.zIsnegative {
		bls.fq12.conjugate(f, f)
	}
	return true
}

func (bls *blsInstance) millerLoopWithNaf(f *fe12, coeffs [][]fe6C, g1Points []*pointG1) {
	j := 0
	for i := len(bls.zNaf) - 1; i >= 0; i-- {
		bls.fq12.square(f, f)
		for k, point := range g1Points {
			bls.ell(f, &(coeffs[k])[j], point)
		}
		j++
		if bls.zNaf[i] != 0 {
			for k, point := range g1Points {
				bls.ell(f, &(coeffs)[k][j], point)
			}
			j++
		}
	}
}

func (bls *blsInstance) millerLoopWithoutNaf(f *fe12, coeffs [][]fe6C, g1Points []*pointG1) {
	j := 0
	for i := bls.z.BitLen() - 2; i >= 0; i-- {
		if j > 0 {
			bls.fq12.square(f, f)
		}
		for k, point := range g1Points {
			bls.ell(f, &(coeffs[k])[j], point)
		}
		j++
		if bls.z.Bit(int(i)) != 0 {
			for k, point := range g1Points {
				bls.ell(f, &(coeffs)[k][j], point)
			}
			j++
		}
	}
}

func (bls *blsInstance) expByZ(c, a *fe12) {
	bls.fq12.cyclotomicExp(c, a, bls.z)
	if bls.zIsnegative {
		bls.fq12.conjugate(c, c)
	}
}

func (bls *blsInstance) finalExp(f *fe12) bool {
	fq := bls.fq12

	f1 := fq.new()
	fq.frobeniusMap(f1, f, 6)
	f2 := fq.new()
	if ok := fq.inverse(f2, f); !ok {
		return false
	}

	r := fq.new()
	fq.mul(r, f1, f2)
	fq.frobeniusMap(f2, r, 2)
	fq.mul(r, r, f2)

	// hard part
	y0 := fq.new()
	fq.cyclotomicSquare(y0, r)
	fq.conjugate(y0, y0)

	y5 := fq.new()
	bls.expByZ(y5, r)

	y1 := fq.new()
	fq.cyclotomicSquare(y1, y5)

	y3 := fq.new()
	fq.mul(y3, y0, y5)

	y2 := fq.new()
	bls.expByZ(y0, y3)
	bls.expByZ(y2, y0)

	y4 := fq.new()
	bls.expByZ(y4, y2)
	fq.mul(y4, y4, y1)

	bls.expByZ(y1, y4)

	fq.conjugate(y3, y3)
	fq.mul(y1, y1, y3)
	fq.mul(y1, y1, r)

	fq.conjugate(y3, r)
	fq.mul(y0, y0, r)
	fq.frobeniusMap(y0, y0, 3)

	fq.mul(y4, y4, y3)
	fq.frobeniusMap(y4, y4, 1)

	fq.mul(y5, y5, y2)
	fq.frobeniusMap(y5, y5, 2)

	fq.mul(y5, y5, y0)
	fq.mul(y5, y5, y4)
	fq.mul(y5, y5, y1)

	fq.copy(f, y5)
	return true
}

func (bls *blsInstance) pair(g1Point *pointG1, g2Point *pointG22) (*fe12, bool) {
	f := bls.fq12.one()
	if bls.g1.isZero(g1Point) || bls.g2.isZero(g2Point) {
		return f, true
	}
	if ok := bls.millerLoop(f, []*pointG1{g1Point}, []*pointG22{g2Point}); !ok {
		return nil, false
	}
	if ok := bls.finalExp(f); !ok {
		return nil, false
	}
	return f, true
}

func (bls *blsInstance) multiPair(g1Points []*pointG1, g2Points []*pointG22) (*fe12, bool) {
	if len(g1Points) != len(g2Points) {
		return nil, false
	}
	if !GAS_METERING_MODE {
		if len(g1Points) == 0 {
			return nil, false
		}
	}
	var _g1Points []*pointG1
	var _g2Points []*pointG22
	for i := 0; i < len(g1Points); i++ {
		if !bls.g1.isZero(g1Points[i]) && !bls.g2.isZero(g2Points[i]) {
			_g1Points = append(_g1Points, g1Points[i])
			_g2Points = append(_g2Points, g2Points[i])
		}
	}
	f := bls.fq12.one()
	if len(_g1Points) == 0 {
		return f, true
	}
	if ok := bls.millerLoop(f, _g1Points, _g2Points); !ok {
		return nil, false
	}
	if ok := bls.finalExp(f); !ok {
		return nil, false
	}
	return f, true
}

func (bls *blsInstance) calculateCoeffLength() int {
	j := 0
	if bls.preferNaf {
		for i := len(bls.zNaf) - 1; i >= 0; i-- {
			if bls.zNaf[i] != 0 {
				j++
			}
			j++
		}
	} else {
		for i := bls.z.BitLen() - 2; i >= 0; i-- {
			if bls.z.Bit(i) != 0 {
				j++
			}
			j++
		}
	}
	return j
}
