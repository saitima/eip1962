package eip

import (
	"math/big"
)

type mnt4Instance struct {
	z               *big.Int
	zIsnegative     bool
	expW0Isnegative bool
	expW0           *big.Int
	expW1           *big.Int
	fq4             *fq4
	g1              *g1
	g2              *g22
	twist           *fe2
	t               *[13]*fe2
}

type (
	doublingCoeffs struct {
		ch  *fe2
		c4c *fe2
		cj  *fe2
		cl  *fe2
	}

	additionCoeffs struct {
		cl1 *fe2
		crz *fe2
	}

	precomputedG1 struct {
		x  fieldElement
		y  fieldElement
		xx *fe2 // twist of x
		yy *fe2 // twist of y
	}

	precomputedG2 struct {
		x              *fe2
		y              *fe2
		xx             *fe2 // twist of x
		yy             *fe2 // twist of y
		doubleCoeffs   []*doublingCoeffs
		additionCoeffs []*additionCoeffs
	}
	extendedCoordinates struct {
		x *fe2
		y *fe2
		z *fe2
		t *fe2
	}
)

func newMnt4Instance(z *big.Int, zIsnegative bool, expW0, expW1 *big.Int, expW0Isnegative bool, fq4 *fq4, g1 *g1, g2 *g22, twist *fe2) *mnt4Instance {
	mnt4 := &mnt4Instance{
		z:               z,
		zIsnegative:     zIsnegative,
		expW0Isnegative: expW0Isnegative,
		expW0:           expW0,
		expW1:           expW1,
		fq4:             fq4,
		g1:              g1,
		g2:              g2,
		twist:           twist,
		t:               new([13]*fe2),
	}
	for i := 0; i < 13; i++ {
		mnt4.t[i] = mnt4.fq4.f.newElement()
	}
	return mnt4
}

func (mnt4 *mnt4Instance) precomputeG1(precomputedG1 *precomputedG1, g1Point *pointG1) {
	mnt4.fq4.f.f.copy(precomputedG1.x, g1Point[0])
	mnt4.fq4.f.f.copy(precomputedG1.y, g1Point[1])
	mnt4.fq4.f.mulByFq(precomputedG1.xx, mnt4.twist, g1Point[0])
	mnt4.fq4.f.mulByFq(precomputedG1.yy, mnt4.twist, g1Point[1])
}

func (mnt4 *mnt4Instance) doublingStep(coeffs *doublingCoeffs, r *extendedCoordinates) {
	fq2 := mnt4.fq4.f
	t := mnt4.t
	fq2.square(t[0], r.t)  // a
	fq2.square(t[1], r.x)  // b
	fq2.square(t[2], r.y)  // c
	fq2.square(t[3], t[2]) // d

	fq2.add(t[4], r.x, t[2]) // e
	fq2.square(t[4], t[4])
	fq2.sub(t[4], t[4], t[1])
	fq2.sub(t[4], t[4], t[3])

	fq2.mul(t[5], mnt4.g2.a, t[0])
	fq2.add(t[5], t[5], t[1])
	fq2.add(t[5], t[5], t[1])
	fq2.add(t[5], t[5], t[1]) // f

	fq2.square(t[6], t[5]) // g

	fq2.double(t[7], t[3])
	fq2.double(t[7], t[7])
	fq2.double(t[7], t[7]) // d8

	fq2.double(t[8], t[4])
	fq2.double(t[8], t[8]) // t0

	fq2.sub(t[9], t[6], t[8]) // x

	fq2.double(t[10], t[4])
	fq2.sub(t[10], t[10], t[9])
	fq2.mul(t[10], t[10], t[5])
	fq2.sub(t[10], t[10], t[7]) // y

	fq2.square(t[8], r.z) // t0

	fq2.add(t[11], r.y, r.z)
	fq2.square(t[11], t[11])
	fq2.sub(t[11], t[11], t[2])
	fq2.sub(t[11], t[11], t[8]) // z

	fq2.square(t[12], t[11]) // t

	fq2.add(coeffs.ch, t[11], r.t)
	fq2.square(coeffs.ch, coeffs.ch)
	fq2.sub(coeffs.ch, coeffs.ch, t[12])
	fq2.sub(coeffs.ch, coeffs.ch, t[0])

	fq2.double(coeffs.c4c, t[2])
	fq2.double(coeffs.c4c, coeffs.c4c)

	fq2.add(coeffs.cj, r.t, t[5])
	fq2.square(coeffs.cj, coeffs.cj)
	fq2.sub(coeffs.cj, coeffs.cj, t[6])
	fq2.sub(coeffs.cj, coeffs.cj, t[0])

	fq2.add(coeffs.cl, t[5], r.x)
	fq2.square(coeffs.cl, coeffs.cl)
	fq2.sub(coeffs.cl, coeffs.cl, t[6])
	fq2.sub(coeffs.cl, coeffs.cl, t[1])

	fq2.copy(r.x, t[9])
	fq2.copy(r.y, t[10])
	fq2.copy(r.z, t[11])
	fq2.copy(r.t, t[12])
}

func (mnt4 *mnt4Instance) additionStep(coeffs *additionCoeffs, x, y *fe2, r *extendedCoordinates) {
	fq2 := mnt4.fq4.f
	t := mnt4.t

	fq2.square(t[0], y)   // a
	fq2.mul(t[1], r.t, x) // b

	fq2.add(t[2], r.z, y)
	fq2.square(t[2], t[2])
	fq2.sub(t[2], t[2], t[0])
	fq2.sub(t[2], t[2], r.t)
	fq2.mul(t[2], t[2], r.t) // d

	fq2.sub(t[3], t[1], r.x) // h

	fq2.square(t[4], t[3]) // i

	fq2.double(t[5], t[4])
	fq2.double(t[5], t[5]) // e

	fq2.mul(t[6], t[3], t[5]) // j

	fq2.mul(t[7], r.x, t[5]) // v

	fq2.sub(coeffs.cl1, t[2], r.y)
	fq2.sub(coeffs.cl1, coeffs.cl1, r.y)

	fq2.square(r.x, coeffs.cl1)
	fq2.sub(r.x, r.x, t[6])
	fq2.sub(r.x, r.x, t[7])
	fq2.sub(r.x, r.x, t[7]) // r.x

	fq2.double(t[8], r.y)
	fq2.mul(t[8], t[8], t[6]) // t0

	fq2.sub(r.y, t[7], r.x)
	fq2.mul(r.y, r.y, coeffs.cl1)
	fq2.sub(r.y, r.y, t[8]) // r.r.y

	fq2.add(r.z, r.z, t[3])
	fq2.square(r.z, r.z)
	fq2.sub(r.z, r.z, r.t)
	fq2.sub(r.z, r.z, t[4]) // z

	fq2.square(r.t, r.z)
	fq2.copy(coeffs.crz, r.z)
}

func (mnt4 *mnt4Instance) precomputeG2(precomputedG2 *precomputedG2, g2Point *pointG22, twistInv *fe2) {
	fq4 := mnt4.fq4

	xTwist, yTwist := mnt4.fq4.f.newElement(), mnt4.fq4.f.newElement()

	fq4.f.mul(xTwist, g2Point[0], twistInv)
	fq4.f.mul(yTwist, g2Point[1], twistInv)

	fq4.f.copy(precomputedG2.x, g2Point[0])
	fq4.f.copy(precomputedG2.y, g2Point[1])
	fq4.f.copy(precomputedG2.xx, xTwist)
	fq4.f.copy(precomputedG2.yy, yTwist)

	r := &extendedCoordinates{
		fq4.f.newElement(),
		fq4.f.newElement(),
		fq4.f.newElement(),
		fq4.f.newElement(),
	}

	fq4.f.copy(r.x, g2Point[0])
	fq4.f.copy(r.y, g2Point[1])
	fq4.f.copy(r.z, fq4.f.one())
	fq4.f.copy(r.t, fq4.f.one())
	d, a := 0, 0
	for i := mnt4.z.BitLen() - 2; i >= 0; i-- {
		mnt4.doublingStep(precomputedG2.doubleCoeffs[d], r)
		d++
		if mnt4.z.Bit(i) != 0 {
			mnt4.additionStep(precomputedG2.additionCoeffs[a], g2Point[0], g2Point[1], r)
			a++
		}
	}

	if mnt4.zIsnegative {
		rzInv, rzInv2, rzInv3 := mnt4.fq4.f.newElement(), mnt4.fq4.f.newElement(), mnt4.fq4.f.newElement()
		fq4.f.inverse(rzInv, r.z)
		fq4.f.square(rzInv2, rzInv)
		fq4.f.mul(rzInv3, rzInv2, rzInv)

		// affine forms
		minusRxAffine, minusRyAffine := mnt4.fq4.f.newElement(), mnt4.fq4.f.newElement()
		fq4.f.mul(minusRxAffine, rzInv2, r.x)
		fq4.f.mul(minusRyAffine, rzInv3, r.y)
		fq4.f.neg(minusRyAffine, minusRyAffine)
		a = len(precomputedG2.additionCoeffs) - 1 // hack
		mnt4.additionStep(precomputedG2.additionCoeffs[a], minusRxAffine, minusRyAffine, r)
	}
}

func (mnt4 *mnt4Instance) atePairingLoop(f *fe4, g1Point *pointG1, g2Point *pointG22) {
	// TODO: check that points are in affine form
	fq4 := mnt4.fq4
	twistInv := mnt4.fq4.f.newElement()
	fq4.f.inverse(twistInv, mnt4.twist)

	p := &precomputedG1{
		fq4.f.f.newFieldElement(),
		fq4.f.f.newFieldElement(),
		fq4.f.newElement(),
		fq4.f.newElement(),
	}
	mnt4.precomputeG1(p, g1Point)

	doubleCount, addCount := mnt4.calculateCoeffSize()
	q := &precomputedG2{
		x:              fq4.f.newElement(),
		y:              fq4.f.newElement(),
		xx:             fq4.f.newElement(),
		yy:             fq4.f.newElement(),
		doubleCoeffs:   make([]*doublingCoeffs, doubleCount),
		additionCoeffs: make([]*additionCoeffs, addCount),
	}
	for i := 0; i < doubleCount; i++ {
		q.doubleCoeffs[i] = &doublingCoeffs{
			fq4.f.newElement(),
			fq4.f.newElement(),
			fq4.f.newElement(),
			fq4.f.newElement(),
		}
	}
	for i := 0; i < addCount; i++ {
		q.additionCoeffs[i] = &additionCoeffs{
			fq4.f.newElement(),
			fq4.f.newElement(),
		}
	}

	mnt4.precomputeG2(q, g2Point, twistInv)

	l1Coeff := fq4.f.zero()
	fq4.f.f.copy(l1Coeff[0], p.x)
	fq4.f.sub(l1Coeff, l1Coeff, q.xx)

	ff := fq4.one()
	d, a := 0, 0
	dc, ac := &doublingCoeffs{
		fq4.f.newElement(),
		fq4.f.newElement(),
		fq4.f.newElement(),
		fq4.f.newElement(),
	}, &additionCoeffs{
		fq4.f.newElement(),
		fq4.f.newElement(),
	}
	gRR, gRQ := mnt4.fq4.newElement(), mnt4.fq4.newElement()
	t := mnt4.fq4.f.newElement()
	for i := mnt4.z.BitLen() - 2; i >= 0; i-- {
		dc = q.doubleCoeffs[d]

		d++
		fq4.f.mul(gRR[0], dc.cj, p.xx)
		fq4.f.neg(gRR[0], gRR[0])
		fq4.f.add(gRR[0], gRR[0], dc.cl)
		fq4.f.sub(gRR[0], gRR[0], dc.c4c)
		fq4.f.mul(gRR[1], dc.ch, p.yy)
		fq4.square(ff, ff)
		fq4.mul(ff, ff, gRR)

		if mnt4.z.Bit(i) != 0 {
			ac = q.additionCoeffs[a]
			a++
			fq4.f.mul(t, l1Coeff, ac.cl1)
			fq4.f.mul(gRQ[0], ac.crz, p.yy)
			fq4.f.mul(gRQ[1], ac.crz, q.yy)
			fq4.f.add(gRQ[1], gRQ[1], t)
			fq4.f.neg(gRQ[1], gRQ[1])
			fq4.mul(ff, ff, gRQ)
		}
	}

	if mnt4.zIsnegative {
		ac = q.additionCoeffs[a]
		fq4.f.mul(t, l1Coeff, ac.cl1)
		fq4.f.mul(gRQ[0], ac.crz, p.yy)
		fq4.f.mul(gRQ[1], ac.crz, q.yy)
		fq4.f.add(gRQ[1], gRQ[1], t)
		fq4.f.neg(gRQ[1], gRQ[1])
		fq4.mul(ff, ff, gRQ)
		fq4.inverse(ff, ff) // TODO: check that f has inverse
	}

	fq4.mul(f, f, ff)
}

func (mnt4 *mnt4Instance) millerLoop(f *fe4, g1Points []*pointG1, g2Points []*pointG22) {
	for i := 0; i < len(g1Points); i++ {
		mnt4.atePairingLoop(f, g1Points[i], g2Points[i])
	}
}

func (mnt4 *mnt4Instance) finalexp(f *fe4) {
	fInv, first, firstInv := mnt4.fq4.newElement(), mnt4.fq4.newElement(), mnt4.fq4.newElement()
	mnt4.fq4.inverse(fInv, f)
	mnt4.finalexpPart1(first, f, fInv)
	mnt4.finalexpPart1(firstInv, fInv, f)
	mnt4.finalexpPart2(f, first, firstInv)
}

func (mnt4 *mnt4Instance) finalexpPart1(f, elt, eltInv *fe4) {
	mnt4.fq4.frobeniusMap(f, elt, 2)
	mnt4.fq4.mul(f, f, eltInv)

}

func (mnt4 *mnt4Instance) finalexpPart2(f, elt, eltInv *fe4) {
	w0Part, w1Part := mnt4.fq4.newElement(), mnt4.fq4.newElement()
	mnt4.fq4.frobeniusMap(w1Part, elt, 1)
	mnt4.fq4.exp(w1Part, w1Part, mnt4.expW1)
	if mnt4.zIsnegative {
		mnt4.fq4.exp(w0Part, eltInv, mnt4.expW0)
	} else {
		mnt4.fq4.exp(w0Part, elt, mnt4.expW0)
	}
	mnt4.fq4.mul(f, w0Part, w1Part)
}

func (mnt4 *mnt4Instance) calculateCoeffSize() (int, int) {
	d, a := 0, 0
	for i := 0; i < mnt4.z.BitLen()-1; i++ {
		d++
		if mnt4.z.Bit(i) != 0 {
			a++
		}
	}
	if mnt4.zIsnegative {
		a++
	}
	return d, a
}

// Pair ..
func (mnt4 *mnt4Instance) pair(g1Point *pointG1, g2Point *pointG22) *fe4 {
	f := mnt4.fq4.one()
	mnt4.millerLoop(f, []*pointG1{g1Point}, []*pointG22{g2Point})
	mnt4.finalexp(f)
	return f
}

func (mnt4 *mnt4Instance) multiPair(g1Points []*pointG1, g2Points []*pointG22) *fe4 {
	f := mnt4.fq4.one()
	mnt4.millerLoop(f, g1Points, g2Points)
	mnt4.finalexp(f)
	return f
}
