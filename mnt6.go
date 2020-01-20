package eip

import (
	"math/big"
)

type mnt6Instance struct {
	z               *big.Int
	zIsnegative     bool
	expW0Isnegative bool
	expW0           *big.Int
	expW1           *big.Int
	fq6             *fq6q
	g1              *g1
	g2              *g23
	twist           *fe3
	t               *[13]*fe3
}

type (
	doublingCoeffs_6 struct {
		ch  *fe3
		c4c *fe3
		cj  *fe3
		cl  *fe3
	}

	additionCoeffs_6 struct {
		cl1 *fe3
		crz *fe3
	}

	precomputedG1_6 struct {
		x  fieldElement
		y  fieldElement
		xx *fe3 // twist of x
		yy *fe3 // twist of y
	}

	precomputedG2_6 struct {
		x              *fe3
		y              *fe3
		xx             *fe3 // twist of x
		yy             *fe3 // twist of y
		doubleCoeffs   []*doublingCoeffs_6
		additionCoeffs []*additionCoeffs_6
	}
	extendedCoordinates6 struct {
		x *fe3
		y *fe3
		z *fe3
		t *fe3
	}
)

func newMNT6Instance(z *big.Int, zIsnegative bool, expW0, expW1 *big.Int, expW0Isnegative bool, fq6 *fq6q, g1 *g1, g2 *g23, twist *fe3) *mnt6Instance {
	mnt6 := &mnt6Instance{
		z:               z,
		zIsnegative:     zIsnegative,
		expW0Isnegative: expW0Isnegative,
		expW0:           expW0,
		expW1:           expW1,
		fq6:             fq6,
		g1:              g1,
		g2:              g2,
		twist:           twist,
		t:               new([13]*fe3),
	}
	for i := 0; i < 13; i++ {
		mnt6.t[i] = mnt6.fq6.f.newElement()
	}
	return mnt6
}

func (mnt6 *mnt6Instance) precomputeG1(precomputedG1_6 *precomputedG1_6, g1Point *pointG1) {
	mnt6.fq6.f.f.copy(precomputedG1_6.x, g1Point[0])
	mnt6.fq6.f.f.copy(precomputedG1_6.y, g1Point[1])
	mnt6.fq6.f.mulByFq(precomputedG1_6.xx, mnt6.twist, g1Point[0])
	mnt6.fq6.f.mulByFq(precomputedG1_6.yy, mnt6.twist, g1Point[1])
}

func (mnt6 *mnt6Instance) doublingStep(coeffs *doublingCoeffs_6, r *extendedCoordinates6) {
	fq3 := mnt6.fq6.f
	t := mnt6.t
	fq3.square(t[0], r.t)  // a
	fq3.square(t[1], r.x)  // b
	fq3.square(t[2], r.y)  // c
	fq3.square(t[3], t[2]) // d

	fq3.add(t[4], r.x, t[2]) // e
	fq3.square(t[4], t[4])
	fq3.sub(t[4], t[4], t[1])
	fq3.sub(t[4], t[4], t[3])

	fq3.mul(t[5], mnt6.g2.a, t[0])
	fq3.add(t[5], t[5], t[1])
	fq3.add(t[5], t[5], t[1])
	fq3.add(t[5], t[5], t[1]) // f

	fq3.square(t[6], t[5]) // g

	fq3.double(t[7], t[3])
	fq3.double(t[7], t[7])
	fq3.double(t[7], t[7]) // d8

	fq3.double(t[8], t[4])
	fq3.double(t[8], t[8]) // t0

	fq3.sub(t[9], t[6], t[8]) // x

	fq3.double(t[10], t[4])
	fq3.sub(t[10], t[10], t[9])
	fq3.mul(t[10], t[10], t[5])
	fq3.sub(t[10], t[10], t[7]) // y

	fq3.square(t[8], r.z) // t0

	fq3.add(t[11], r.y, r.z)
	fq3.square(t[11], t[11])
	fq3.sub(t[11], t[11], t[2])
	fq3.sub(t[11], t[11], t[8]) // z

	fq3.square(t[12], t[11]) // t

	fq3.add(coeffs.ch, t[11], r.t)
	fq3.square(coeffs.ch, coeffs.ch)
	fq3.sub(coeffs.ch, coeffs.ch, t[12])
	fq3.sub(coeffs.ch, coeffs.ch, t[0])

	fq3.double(coeffs.c4c, t[2])
	fq3.double(coeffs.c4c, coeffs.c4c)

	fq3.add(coeffs.cj, r.t, t[5])
	fq3.square(coeffs.cj, coeffs.cj)
	fq3.sub(coeffs.cj, coeffs.cj, t[6])
	fq3.sub(coeffs.cj, coeffs.cj, t[0])

	fq3.add(coeffs.cl, t[5], r.x)
	fq3.square(coeffs.cl, coeffs.cl)
	fq3.sub(coeffs.cl, coeffs.cl, t[6])
	fq3.sub(coeffs.cl, coeffs.cl, t[1])

	fq3.copy(r.x, t[9])
	fq3.copy(r.y, t[10])
	fq3.copy(r.z, t[11])
	fq3.copy(r.t, t[12])
}

func (mnt6 *mnt6Instance) additionStep(coeffs *additionCoeffs_6, x, y *fe3, r *extendedCoordinates6) {
	fq3 := mnt6.fq6.f
	t := mnt6.t

	fq3.square(t[0], y)   // a
	fq3.mul(t[1], r.t, x) // b

	fq3.add(t[2], r.z, y)
	fq3.square(t[2], t[2])
	fq3.sub(t[2], t[2], t[0])
	fq3.sub(t[2], t[2], r.t)
	fq3.mul(t[2], t[2], r.t) // d

	fq3.sub(t[3], t[1], r.x) // h

	fq3.square(t[4], t[3]) // i

	fq3.double(t[5], t[4])
	fq3.double(t[5], t[5]) // e

	fq3.mul(t[6], t[3], t[5]) // j

	fq3.mul(t[7], r.x, t[5]) // v

	fq3.sub(coeffs.cl1, t[2], r.y)
	fq3.sub(coeffs.cl1, coeffs.cl1, r.y)

	fq3.square(r.x, coeffs.cl1)
	fq3.sub(r.x, r.x, t[6])
	fq3.sub(r.x, r.x, t[7])
	fq3.sub(r.x, r.x, t[7]) // r.x

	fq3.double(t[8], r.y)
	fq3.mul(t[8], t[8], t[6]) // t0

	fq3.sub(r.y, t[7], r.x)
	fq3.mul(r.y, r.y, coeffs.cl1)
	fq3.sub(r.y, r.y, t[8]) // r.r.y

	fq3.add(r.z, r.z, t[3])
	fq3.square(r.z, r.z)
	fq3.sub(r.z, r.z, r.t)
	fq3.sub(r.z, r.z, t[4]) // z

	fq3.square(r.t, r.z)
	fq3.copy(coeffs.crz, r.z)
}

func (mnt6 *mnt6Instance) precomputeG2(precomputedG2 *precomputedG2_6, g2Point *pointG23, twistInv *fe3) {
	fq6 := mnt6.fq6

	xTwist, yTwist := mnt6.fq6.f.newElement(), mnt6.fq6.f.newElement()

	fq6.f.mul(xTwist, g2Point[0], twistInv)
	fq6.f.mul(yTwist, g2Point[1], twistInv)

	fq6.f.copy(precomputedG2.x, g2Point[0])
	fq6.f.copy(precomputedG2.y, g2Point[1])
	fq6.f.copy(precomputedG2.xx, xTwist)
	fq6.f.copy(precomputedG2.yy, yTwist)

	r := &extendedCoordinates6{
		fq6.f.newElement(),
		fq6.f.newElement(),
		fq6.f.newElement(),
		fq6.f.newElement(),
	}

	fq6.f.copy(r.x, g2Point[0])
	fq6.f.copy(r.y, g2Point[1])
	fq6.f.copy(r.z, fq6.f.one())
	fq6.f.copy(r.t, fq6.f.one())
	d, a := 0, 0
	for i := mnt6.z.BitLen() - 2; i >= 0; i-- {
		mnt6.doublingStep(precomputedG2.doubleCoeffs[d], r)
		d++
		if mnt6.z.Bit(i) != 0 {
			mnt6.additionStep(precomputedG2.additionCoeffs[a], g2Point[0], g2Point[1], r)
			a++
		}
	}

	if mnt6.zIsnegative {
		rzInv, rzInv2, rzInv3 := mnt6.fq6.f.newElement(), mnt6.fq6.f.newElement(), mnt6.fq6.f.newElement()
		fq6.f.inverse(rzInv, r.z)
		fq6.f.square(rzInv2, rzInv)
		fq6.f.mul(rzInv3, rzInv2, rzInv)

		// affine forms
		minusRxAffine, minusRyAffine := mnt6.fq6.f.newElement(), mnt6.fq6.f.newElement()
		fq6.f.mul(minusRxAffine, rzInv2, r.x)
		fq6.f.mul(minusRyAffine, rzInv3, r.y)
		fq6.f.neg(minusRyAffine, minusRyAffine)
		a = len(precomputedG2.additionCoeffs) - 1 // hack
		mnt6.additionStep(precomputedG2.additionCoeffs[a], minusRxAffine, minusRyAffine, r)
	}
}

func (mnt6 *mnt6Instance) atePairingLoop(f *fe6q, g1Point *pointG1, g2Point *pointG23) {
	// TODO: check that points are in affine form
	fq6 := mnt6.fq6
	twistInv := mnt6.fq6.f.newElement()
	fq6.f.inverse(twistInv, mnt6.twist)

	p := &precomputedG1_6{
		fq6.f.f.newFieldElement(),
		fq6.f.f.newFieldElement(),
		fq6.f.newElement(),
		fq6.f.newElement(),
	}
	mnt6.precomputeG1(p, g1Point)

	doubleCount, addCount := mnt6.calculateCoeffSize()
	q := &precomputedG2_6{
		x:              fq6.f.newElement(),
		y:              fq6.f.newElement(),
		xx:             fq6.f.newElement(),
		yy:             fq6.f.newElement(),
		doubleCoeffs:   make([]*doublingCoeffs_6, doubleCount),
		additionCoeffs: make([]*additionCoeffs_6, addCount),
	}
	for i := 0; i < doubleCount; i++ {
		q.doubleCoeffs[i] = &doublingCoeffs_6{
			fq6.f.newElement(),
			fq6.f.newElement(),
			fq6.f.newElement(),
			fq6.f.newElement(),
		}
	}
	for i := 0; i < addCount; i++ {
		q.additionCoeffs[i] = &additionCoeffs_6{
			fq6.f.newElement(),
			fq6.f.newElement(),
		}
	}

	mnt6.precomputeG2(q, g2Point, twistInv)

	l1Coeff := fq6.f.zero()
	fq6.f.f.copy(l1Coeff[0], p.x)
	fq6.f.sub(l1Coeff, l1Coeff, q.xx)

	ff := fq6.one()
	d, a := 0, 0
	dc, ac := &doublingCoeffs_6{
		fq6.f.newElement(),
		fq6.f.newElement(),
		fq6.f.newElement(),
		fq6.f.newElement(),
	}, &additionCoeffs_6{
		fq6.f.newElement(),
		fq6.f.newElement(),
	}
	gRR, gRQ := mnt6.fq6.newElement(), mnt6.fq6.newElement()
	t := mnt6.fq6.f.newElement()
	for i := mnt6.z.BitLen() - 2; i >= 0; i-- {
		dc = q.doubleCoeffs[d]
		d++
		fq6.f.mul(gRR[0], dc.cj, p.xx)
		fq6.f.neg(gRR[0], gRR[0])
		fq6.f.add(gRR[0], gRR[0], dc.cl)
		fq6.f.sub(gRR[0], gRR[0], dc.c4c)
		fq6.f.mul(gRR[1], dc.ch, p.yy)
		fq6.square(ff, ff)
		fq6.mul(ff, ff, gRR)
		if mnt6.z.Bit(i) != 0 {
			ac = q.additionCoeffs[a]
			a++
			fq6.f.mul(t, l1Coeff, ac.cl1)
			fq6.f.mul(gRQ[0], ac.crz, p.yy)
			fq6.f.mul(gRQ[1], ac.crz, q.yy)
			fq6.f.add(gRQ[1], gRQ[1], t)
			fq6.f.neg(gRQ[1], gRQ[1])
			fq6.mul(ff, ff, gRQ)
		}
	}

	if mnt6.zIsnegative {
		ac = q.additionCoeffs[a]
		fq6.f.mul(t, l1Coeff, ac.cl1)
		fq6.f.mul(gRQ[0], ac.crz, p.yy)
		fq6.f.mul(gRQ[1], ac.crz, q.yy)
		fq6.f.add(gRQ[1], gRQ[1], t)
		fq6.f.neg(gRQ[1], gRQ[1])
		fq6.mul(ff, ff, gRQ)
		fq6.inverse(ff, ff) // TODO: check that f has inverse
	}

	fq6.mul(f, f, ff)
}

func (mnt6 *mnt6Instance) millerLoop(f *fe6q, g1Points []*pointG1, g2Points []*pointG23) {
	for i := 0; i < len(g1Points); i++ {
		mnt6.atePairingLoop(f, g1Points[i], g2Points[i])
	}
}

func (mnt6 *mnt6Instance) finalexp(f *fe6q) {
	fInv, first, firstInv := mnt6.fq6.newElement(), mnt6.fq6.newElement(), mnt6.fq6.newElement()
	mnt6.fq6.inverse(fInv, f)
	mnt6.finalexpPart1(first, f, fInv)
	mnt6.finalexpPart1(firstInv, fInv, f)
	mnt6.finalexpPart2(f, first, firstInv)
}

func (mnt6 *mnt6Instance) finalexpPart1(f, elt, eltInv *fe6q) {
	t := mnt6.fq6.newElement()
	mnt6.fq6.frobeniusMap(f, elt, 3)
	mnt6.fq6.mul(t, f, eltInv)
	mnt6.fq6.frobeniusMap(f, t, 1)
	mnt6.fq6.mul(f, f, t)
}

func (mnt6 *mnt6Instance) finalexpPart2(f, elt, eltInv *fe6q) {
	w0Part, w1Part := mnt6.fq6.newElement(), mnt6.fq6.newElement()
	mnt6.fq6.frobeniusMap(w1Part, elt, 1)
	mnt6.fq6.exp(w1Part, w1Part, mnt6.expW1)
	if mnt6.zIsnegative {
		mnt6.fq6.exp(w0Part, eltInv, mnt6.expW0)
	} else {
		mnt6.fq6.exp(w0Part, elt, mnt6.expW0)
	}
	mnt6.fq6.mul(f, w0Part, w1Part)
}

func (mnt6 *mnt6Instance) calculateCoeffSize() (int, int) {
	d, a := 0, 0
	for i := 0; i < mnt6.z.BitLen()-1; i++ {
		d++
		if mnt6.z.Bit(i) != 0 {
			a++
		}
	}
	if mnt6.zIsnegative {
		a++
	}
	return d, a
}

// Pair ..
func (mnt6 *mnt6Instance) pair(g1Point *pointG1, g2Point *pointG23) *fe6q {
	f := mnt6.fq6.one()
	mnt6.millerLoop(f, []*pointG1{g1Point}, []*pointG23{g2Point})
	mnt6.finalexp(f)
	return f
}

func (mnt6 *mnt6Instance) multiPair(g1Points []*pointG1, g2Points []*pointG23) *fe6q {
	f := mnt6.fq6.one()
	mnt6.millerLoop(f, g1Points, g2Points)
	mnt6.finalexp(f)
	return f
}
