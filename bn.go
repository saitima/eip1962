package fp

// import (
// 	"math/big"
// )

// type (
// 	ellCoeffsBN = []fe6
// )

// // bnInstance ..
// type bnInstance struct {
// 	u                        *big.Int
// 	sixUPlus2                *big.Int
// 	uIsnegative              bool
// 	twistType                uint8
// 	g1                       *g1
// 	g2                       *g22
// 	fq12                     *fq12
// 	nonResidueInPMinus1Over2 *fe2
// 	t2                       [16]*fe2
// 	t12                      [16]*fe12
// }

// // NewBNInstance ..
// func newBNInstance(u, sixUplus2 *big.Int, uIsnegative bool, twistType uint8, g1 *g1, g2 *g22, fq12 *fq12, nonResidue *fe2) bnInstance {
// 	bn := bnInstance{
// 		u:                        u,
// 		sixUPlus2:                sixUplus2,
// 		uIsnegative:              uIsnegative,
// 		twistType:                twistType,
// 		g1:                       g1,
// 		g2:                       g2,
// 		fq12:                     fq12,
// 		nonResidueInPMinus1Over2: nonResidue,
// 	}
// 	for i := 0; i < 16; i++ {
// 		bn.t2[i] = bn.fq12.f.f.newElement()
// 		bn.t12[i] = bn.fq12.newElement()
// 	}
// 	return bn
// }

// func (bn *bnInstance) doublingStep(coeff fe6, r *pointG22, twoInv fieldElement) {
// 	fq2 := bn.fq12.f.f
// 	t := bn.t2
// 	// X*Y/2
// 	fq2.mul(t[0], r[0], r[1])
// 	fq2.mulByFq(t[0], t[0], twoInv)
// 	// Y^2
// 	fq2.square(t[1], r[1])
// 	// Z^2
// 	fq2.square(t[2], r[2])
// 	// 3*Z^2
// 	fq2.double(t[3], t[2])
// 	fq2.add(t[3], t[3], t[2])
// 	// 3*b*Z^2
// 	fq2.mul(t[4], t[3], bn.g2.b)
// 	// 9*b*Z^2
// 	fq2.double(t[5], t[4])
// 	fq2.add(t[5], t[5], t[4])
// 	// (Y^2 + 9*b*Z^2)/2
// 	fq2.add(t[6], t[1], t[5])
// 	fq2.mulByFq(t[6], t[6], twoInv)
// 	// (Y + Z)^2
// 	fq2.add(t[7], r[1], r[2])
// 	fq2.square(t[7], t[7])
// 	// Y^2 + Z^2
// 	fq2.add(t[8], t[1], t[2])
// 	// 2*Y*Z
// 	fq2.sub(t[9], t[7], t[8])
// 	// 3*b*Z^2 - Y^2
// 	fq2.sub(t[10], t[4], t[1])
// 	// X^2
// 	fq2.square(t[11], r[0])
// 	// (3*b*Z^2)^2
// 	fq2.square(t[12], t[4])
// 	// X = (Y^2 - 9*b*Z^2)*(X*Y/2)
// 	fq2.sub(r[0], t[1], t[5])
// 	fq2.mul(r[0], r[0], t[0])
// 	// 27*b^2*Z^4
// 	fq2.double(t[13], t[12])
// 	fq2.add(t[13], t[13], t[12])
// 	// Y = ((Y^2+9*b*Z^2)/2)^2 - 27*b^2*Z^4
// 	fq2.add(r[1], t[1], t[5])
// 	fq2.mulByFq(r[1], r[1], twoInv)
// 	fq2.square(r[1], r[1])
// 	fq2.sub(r[1], r[1], t[13])
// 	// Z = 2*Y^3*Z
// 	fq2.mul(r[2], t[9], t[1])
// 	// 3*X^2
// 	fq2.double(t[14], t[11])
// 	fq2.add(t[14], t[14], t[11])
// 	// -2*Y*Z
// 	fq2.neg(t[15], t[9])

// 	switch bn.twistType {
// 	case 1: // M
// 		fq2.copy(coeff[0], t[10]) // 3*b*Z^2 - Y^2
// 		fq2.copy(coeff[1], t[14]) // 3*X^2
// 		fq2.copy(coeff[2], t[15]) // -2*Y*Z
// 		break
// 	case 2:
// 		fq2.copy(coeff[0], t[15])
// 		fq2.copy(coeff[1], t[14])
// 		fq2.copy(coeff[2], t[10])
// 		break
// 	}

// }

// func (bn *bnInstance) additionStep(coeff fe6, r *pointG22, q *pointG22) {
// 	fq2 := bn.fq12.f.f
// 	t := bn.t2
// 	// theta = Y - y*Z
// 	fq2.mul(t[0], q[1], r[2])
// 	fq2.sub(t[0], r[1], t[0])
// 	// lambda = X - x*Z
// 	fq2.mul(t[1], q[0], r[2])
// 	fq2.sub(t[1], r[0], t[1])
// 	// theta^2 = (Y - Y*Z)^2
// 	fq2.square(t[2], t[0])
// 	// lambda^2 = (X - X*Z)^2
// 	fq2.square(t[3], t[1])
// 	// lambda^3 = (X - X*Z)^3
// 	fq2.mul(t[4], t[3], t[1])
// 	// theta^2*Z = (Y - Y*Z)^2 * Z
// 	fq2.mul(t[5], t[2], r[2])
// 	// lambda^2*X = (X - X*Z)^2 * X
// 	fq2.mul(t[6], t[3], r[0])
// 	// H = lambda^3 + theta^2 * Z - 2*lambda^2 * X
// 	fq2.double(t[7], t[6])
// 	fq2.sub(t[7], t[5], t[7])
// 	fq2.add(t[7], t[7], t[4])
// 	// X = lambda * H
// 	fq2.mul(r[0], t[1], t[7])
// 	// (lambda^2*X - H)theta
// 	fq2.sub(t[8], t[6], t[7])
// 	fq2.mul(t[8], t[8], t[0])
// 	// Y = (lambda^2*X - H)theta - lambda^3*Y
// 	fq2.mul(t[9], t[4], r[1])
// 	fq2.sub(r[1], t[8], t[9])
// 	// Z = lambda^3*Z
// 	fq2.mul(r[2], t[4], r[2])
// 	// lambda*y
// 	fq2.mul(t[10], t[1], q[1])
// 	// theata*x - lambda*y
// 	fq2.mul(t[11], t[0], q[0])
// 	fq2.sub(t[11], t[11], t[10])
// 	// -theta
// 	fq2.neg(t[0], t[0])

// 	switch bn.twistType {
// 	case 1: // M
// 		fq2.copy(coeff[0], t[11]) // theata*x - lambda*y
// 		fq2.copy(coeff[1], t[0])  // -theta
// 		fq2.copy(coeff[2], t[1])  // lambda
// 		break
// 	case 2: // D
// 		fq2.copy(coeff[0], t[1])  // lambda
// 		fq2.copy(coeff[1], t[0])  // -theta
// 		fq2.copy(coeff[2], t[11]) // theata*x - lambda*y
// 		break
// 	}

// }

// func (bn *bnInstance) ell(f *fe12, coeffs fe6, p *pointG1) {
// 	// TODO: p needs to be affine/normalized
// 	fq2 := bn.fq12.f.f
// 	switch bn.twistType {
// 	case 1: // M
// 		fq2.mulByFq(coeffs[2], coeffs[2], p[1])
// 		fq2.mulByFq(coeffs[1], coeffs[1], p[0])
// 		bn.fq12.mulBy014(f, coeffs[0], coeffs[1], coeffs[2])
// 	case 2: // D
// 		fq2.mulByFq(coeffs[0], coeffs[0], p[1])
// 		fq2.mulByFq(coeffs[1], coeffs[1], p[0])
// 		bn.fq12.mulBy034(f, coeffs[0], coeffs[1], coeffs[2])
// 	}
// }

// // TODO: twistPoint should be a affine point
// func (bn *bnInstance) prepare(coeffs ellCoeffsBN, Q *pointG22) {
// 	f := bn.fq12.f.f.f
// 	twoInw := f.newFieldElement()
// 	f.double(twoInw, f.one)
// 	f.inverse(twoInw, twoInw)
// 	if bn.g2.isZero(Q) {
// 		// TODO: mark this point as infinity
// 		return
// 	}

// 	T := bn.g2.newPoint()
// 	bn.g2.copy(T, Q)

// 	j := 0
// 	//  skip first msb bit
// 	for i := bn.sixUPlus2.BitLen() - 2; i >= 0; i-- {
// 		coeffs[j] = *bn.fq12.f.newElement()
// 		bn.doublingStep(coeffs[j], T, twoInw)
// 		j++
// 		if bn.sixUPlus2.Bit(int(i)) != 0 {
// 			coeffs[j] = *bn.fq12.f.newElement()
// 			bn.additionStep(coeffs[j], T, Q)
// 			j++
// 		}
// 	}
// 	if bn.uIsnegative {
// 		bn.g2.neg(T, T)
// 	}

// 	// Q1 = π(Q)
// 	Q1 := bn.g2.newPoint()
// 	bn.fq12.f.f.conjugate(Q1[0], Q[0])
// 	bn.fq12.f.f.conjugate(Q1[1], Q[1])
// 	bn.fq12.f.f.mul(Q1[0], Q1[0], bn.fq12.f.frobeniusCoeffs[0][1])
// 	bn.fq12.f.f.mul(Q1[1], Q1[1], bn.nonResidueInPMinus1Over2)
// 	coeffs[j] = *bn.fq12.f.newElement()
// 	bn.additionStep(coeffs[j], T, Q1)
// 	j++

// 	// -Q2 = -π(π(Q))
// 	Q2 := bn.g2.newPoint()
// 	bn.g2.copy(Q2, Q)
// 	bn.fq12.f.f.mul(Q2[0], Q2[0], bn.fq12.f.frobeniusCoeffs[0][2])
// 	coeffs[j] = *bn.fq12.f.newElement()
// 	bn.additionStep(coeffs[j], T, Q2)
// }

// func (bn *bnInstance) millerLoop(f *fe12, g1Points []*pointG1, g2Points []*pointG22) {
// 	coeffs := make([]ellCoeffsBN, len(g1Points))

// 	// prepare and collect miller lines for each pair (Pi,Qi)
// 	// TODO: check points that are normalized/affine form?
// 	// for i := 0; i < len(g1Points); i++ {
// 	for i, _ := range g1Points {
// 		coeffs[i] = make(ellCoeffsBN, bn.calculateCoeffLength())
// 		bn.prepare(coeffs[i], g2Points[i])
// 	}
// 	j := 0
// 	for i := bn.sixUPlus2.BitLen() - 2; i >= 0; i-- {
// 		if j > 0 {
// 			bn.fq12.square(f, f)
// 		}
// 		// doubling coeffs
// 		for k, point := range g1Points {
// 			bn.ell(f, coeffs[k][j], point)
// 		}
// 		j++
// 		// addition coeffs
// 		if bn.sixUPlus2.Bit(int(i)) != 0 {
// 			for k, point := range g1Points {
// 				bn.ell(f, coeffs[k][j], point)
// 			}
// 			j++
// 		}

// 	}

// 	if bn.uIsnegative {
// 		bn.fq12.conjugate(f, f)
// 	}
// 	// Q1 = π(Q)
// 	for k, point := range g1Points {
// 		bn.ell(f, coeffs[k][j], point)
// 	}
// 	j++

// 	// -Q2 = -π(π(Q))
// 	for k, point := range g1Points {
// 		bn.ell(f, coeffs[k][j], point)
// 	}
// }

// func (bn *bnInstance) expByU(c, a *fe12) {
// 	bn.fq12.cyclotomicExp(c, a, bn.u)
// 	if bn.uIsnegative {
// 		bn.fq12.conjugate(c, c)
// 	}
// }

// // p(x) = 36x^4 +36x^3 +24x^2 +6x+1 ⌘ 1 mod 6
// // so degree 12 extension is squaring friendly
// // x^q =x^((q-1))·x= x^(3(q-1)/3)·x = i^((q-1)/6) ·x
// // w = i^(q-1)/6 is a primitive sixth root of unity in Fq
// // x^q = w·x
// // x^(q^2) = (x^q)^2 = (wx)^2=  w^2·x
// func (bn *bnInstance) finalExp(f *fe12) {
// 	fq12 := bn.fq12

// 	f1 := fq12.newElement()
// 	fq12.frobeniusMap(f1, f, 6)

// 	f2 := fq12.newElement()
// 	fq12.inverse(f2, f)
// 	// TODO: handle case when parameter f doesnt have inverse

// 	r := fq12.newElement()
// 	fq12.mul(r, f1, f2)

// 	fq12.copy(f2, r)
// 	fq12.frobeniusMap(r, r, 2)
// 	fq12.mul(r, r, f2)

// 	fp := fq12.newElement()
// 	fq12.frobeniusMap(fp, r, 1)

// 	fp2 := fq12.newElement()
// 	fq12.frobeniusMap(fp2, r, 2)

// 	fp3 := fq12.newElement()
// 	fq12.frobeniusMap(fp3, fp2, 1)

// 	fu := fq12.newElement()
// 	bn.expByU(fu, r)

// 	fu2 := fq12.newElement()
// 	bn.expByU(fu2, fu)

// 	fu3 := fq12.newElement()
// 	bn.expByU(fu3, fu2)

// 	y3 := fq12.newElement()
// 	fq12.frobeniusMap(y3, fu, 1)

// 	fu2p := fq12.newElement()
// 	fq12.frobeniusMap(fu2p, fu2, 1)

// 	fu3p := fq12.newElement()
// 	fq12.frobeniusMap(fu3p, fu3, 1)

// 	y2 := fq12.newElement()
// 	fq12.frobeniusMap(y2, fu2, 2)

// 	y0 := fq12.newElement()
// 	fq12.mul(y0, fp, fp2)
// 	fq12.mul(y0, y0, fp3)

// 	y1 := fq12.newElement()
// 	fq12.conjugate(y1, r)

// 	y5 := fq12.newElement()
// 	fq12.conjugate(y5, fu2)

// 	fq12.conjugate(y3, y3)

// 	y4 := fq12.newElement()
// 	fq12.mul(y4, fu, fu2p)
// 	fq12.conjugate(y4, y4)

// 	y6 := fq12.newElement()
// 	fq12.mul(y6, fu3, fu3p)
// 	fq12.conjugate(y6, y6)

// 	fq12.square(y6, y6)
// 	fq12.mul(y6, y6, y4)
// 	fq12.mul(y6, y6, y5)

// 	t1 := fq12.newElement()
// 	fq12.mul(t1, y3, y5)
// 	fq12.mul(t1, t1, y6)

// 	fq12.mul(y6, y6, y2)

// 	fq12.square(t1, t1)
// 	fq12.mul(t1, t1, y6)
// 	fq12.square(t1, t1)

// 	t0 := fq12.newElement()
// 	fq12.mul(t0, t1, y1)

// 	fq12.mul(t1, t1, y0)

// 	fq12.square(t0, t0)
// 	fq12.mul(t0, t0, t1)

// 	fq12.copy(f, t0)
// }

// // Pair ..
// func (bn *bnInstance) pair(point *pointG1, twistPoint *pointG22) *fe12 {
// 	f := bn.fq12.one()
// 	bn.millerLoop(f, []*pointG1{point}, []*pointG22{twistPoint})
// 	bn.finalExp(f)
// 	return f
// }

// func (bn *bnInstance) multiPair(points []*pointG1, twistPoints []*pointG22) *fe12 {
// 	f := bn.fq12.one()
// 	bn.millerLoop(f, points, twistPoints)
// 	bn.finalExp(f)
// 	return f
// }

// func (bn *bnInstance) calculateCoeffLength() int {
// 	j := 0
// 	for i := bn.sixUPlus2.BitLen() - 2; i >= 0; i-- {
// 		if bn.sixUPlus2.Bit(i) != 0 {
// 			j++
// 		}
// 		j++
// 	}
// 	j = j + 2
// 	return j
// }
