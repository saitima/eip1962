package eip

import (
	"fmt"
	"math"
	"math/big"
)

type pointG23 [3]*fe3

type g23 struct {
	f   *fq3
	a   *fe3
	b   *fe3
	q   *big.Int
	t   [9]*fe3
	inf *pointG23
}

func newG23(f *fq3, a, b, q []byte) (*g23, error) {
	var err error
	a_, b_ := f.newElement(), f.newElement()
	if a != nil {
		a_, err = f.fromBytes(a)
		if err != nil {
			return nil, err
		}
	}
	if b != nil {
		b_, err = f.fromBytes(b)
		if err != nil {
			return nil, err
		}
	}

	t := [9]*fe3{}
	for i := 0; i < 9; i++ {
		t[i] = f.zero()
	}
	g := &g23{
		f:   f,
		a:   a_,
		b:   b_,
		q:   new(big.Int).SetBytes(q),
		t:   t,
		inf: &pointG23{f.zero(), f.one(), f.zero()},
	}
	g.inf = g.newPoint()
	g.f.copy(g.inf[0], g.f.zero())
	g.f.copy(g.inf[1], g.f.one())
	g.f.copy(g.inf[2], g.f.zero())
	return g, nil
}

func (g *g23) newPoint() *pointG23 {
	return &pointG23{g.f.newElement(), g.f.newElement(), g.f.newElement()}
}

func (g *g23) fromBytes(in []byte) (*pointG23, error) {
	byteLen := g.f.f.limbSize * 8 * 3
	if len(in) < 2*byteLen {
		return nil, fmt.Errorf("input string should be equal or larger than 96")
	}
	x, err := g.f.fromBytes(in[:byteLen])
	if err != nil {
		return nil, err
	}
	y, err := g.f.fromBytes(in[byteLen:])
	if err != nil {
		return nil, err
	}
	p := g.newPoint()
	g.f.copy(p[0], x)
	g.f.copy(p[1], y)
	g.f.copy(p[2], g.f.one())
	return p, nil
}

func (g *g23) toBytes(p *pointG23) []byte {
	l := g.f.f.limbSize * 8 * 3
	out := make([]byte, 2*l)
	a := g.newPoint()
	g.affine(a, p)
	copy(out[:l], g.f.toBytes(a[0]))
	copy(out[l:], g.f.toBytes(a[1]))
	return out
}

func (g *g23) copy(q, p *pointG23) *pointG23 {
	g.f.copy(q[0], p[0])
	g.f.copy(q[1], p[1])
	g.f.copy(q[2], p[2])
	return q
}

func (g *g23) affine(q, p *pointG23) *pointG23 {
	if g.isZero(p) {
		g.copy(q, g.inf)
		return q
	}
	t := g.t
	g.f.inverse(t[0], p[2])
	g.f.square(t[1], t[0])
	g.f.mul(q[0], p[0], t[1])
	g.f.mul(t[0], t[0], t[1])
	g.f.mul(q[1], p[1], t[0])
	g.f.copy(q[2], g.f.one())
	return q
}

func (g *g23) toString(p *pointG23) string {
	return fmt.Sprintf(
		"x: %s y: %s, z: %s",
		g.f.toString(p[0]),
		g.f.toString(p[1]),
		g.f.toString(p[2]),
	)
}

func (g *g23) toStringNoTransform(p *pointG23) string {
	return fmt.Sprintf(
		"x: %s y: %s, z: %s",
		g.f.toStringNoTransform(p[0]),
		g.f.toStringNoTransform(p[1]),
		g.f.toStringNoTransform(p[2]),
	)
}

func (g *g23) zero() *pointG23 {
	p := g.newPoint()
	g.f.copy(p[0], g.f.zero())
	g.f.copy(p[1], g.f.one())
	g.f.copy(p[2], g.f.zero())
	return p
}

func (g *g23) isZero(p *pointG23) bool {
	return g.f.isZero(p[2])
}

func (g *g23) equal(p1, p2 *pointG23) bool {
	// TODO: Affine equality ?
	// TODO: P and -P equals why?
	if g.isZero(p1) {
		return g.isZero(p2)
	}
	if g.isZero(p2) {
		return g.isZero(p1)
	}
	t := g.t
	// X1 * Z2^2 == X2 * Z1^2
	// &&
	// Y1 * Z2^3 == Y2 * Z1^3
	g.f.square(t[0], p1[2])
	g.f.square(t[1], p2[2])
	g.f.mul(t[2], t[0], p2[0])
	g.f.mul(t[3], t[1], p1[0])
	g.f.mul(t[0], t[0], p1[2])
	g.f.mul(t[1], t[1], p2[2])
	g.f.mul(t[1], t[1], p1[1])
	g.f.mul(t[0], t[0], p2[1])
	return g.f.equal(t[0], t[1]) && g.f.equal(t[2], t[3])
}

func (g *g23) isOnCurve(p *pointG23) bool {
	if g.isZero(p) {
		return true
	}
	t := g.t
	// Y^2 = X^3 + a Z^4 + b Z^6
	g.f.square(t[0], p[1])    // Y2
	g.f.square(t[1], p[0])    // X2
	g.f.mul(t[1], t[1], p[0]) // X3
	g.f.square(t[2], p[2])    // Z2
	g.f.square(t[3], t[2])    // Z4
	if !g.f.isZero(g.a) {
		g.f.mul(t[4], g.a, t[3])  // aZ4
		g.f.mul(t[4], t[4], p[0]) // aXZ4
		g.f.add(t[1], t[1], t[4]) // X3 + aXZ4
	}
	g.f.mul(t[2], t[2], t[3])    // Z6
	g.f.mul(t[2], g.b, t[2])     // bZ6
	g.f.add(t[1], t[1], t[2])    // X3 + aXZ4 + bZ6
	return g.f.equal(t[0], t[1]) // Y2 == X3 + aZ4 + bZ6
}

func (g *g23) add(r, p1, p2 *pointG23) *pointG23 {
	// http://www.hyperelliptic.org/EFD/gp/auto-shortw-jacobian-0.html#addition-add-2007-bl
	if g.isZero(p1) {
		g.copy(r, p2)
		return r
	}
	if g.isZero(p2) {
		g.copy(r, p1)
		return r
	}
	t := g.t
	g.f.square(t[7], p1[2])    // z1z1
	g.f.mul(t[1], p2[0], t[7]) // u2 = x2 * z1z1
	g.f.mul(t[2], p1[2], t[7]) // z1z1 * z1
	g.f.mul(t[0], p2[1], t[2]) // s2 = y2 * z1z1 * z1
	g.f.square(t[8], p2[2])    // z2z2
	g.f.mul(t[3], p1[0], t[8]) // u1 = x1 * z2z2
	g.f.mul(t[4], p2[2], t[8]) // z2z2 * z2
	g.f.mul(t[2], p1[1], t[4]) // s1 = y1 * z2z2 * z2
	if g.f.equal(t[1], t[3]) {
		if g.f.equal(t[0], t[2]) {
			return g.double(r, p1)
		} else {
			return g.copy(r, g.inf)
		}
	}
	g.f.sub(t[1], t[1], t[3])   // h = u2 - u1
	g.f.double(t[4], t[1])      // 2h
	g.f.square(t[4], t[4])      // i = 2h^2
	g.f.mul(t[5], t[1], t[4])   // j = h*i
	g.f.sub(t[0], t[0], t[2])   // s2 - s1
	g.f.double(t[0], t[0])      // r = 2*(s2 - s1)
	g.f.square(t[6], t[0])      // r^2
	g.f.sub(t[6], t[6], t[5])   // r^2 - j
	g.f.mul(t[3], t[3], t[4])   // v = u1 * i
	g.f.double(t[4], t[3])      // 2*v
	g.f.sub(r[0], t[6], t[4])   // x3 = r^2 - j - 2*v
	g.f.sub(t[4], t[3], r[0])   // v - x3
	g.f.mul(t[6], t[2], t[5])   // s1 * j
	g.f.double(t[6], t[6])      // 2 * s1 * j
	g.f.mul(t[0], t[0], t[4])   // r * (v - x3)
	g.f.sub(r[1], t[0], t[6])   // y3 = r * (v - x3) - (2 * s1 * j)
	g.f.add(t[0], p1[2], p2[2]) // z1 + z2
	g.f.square(t[0], t[0])      // (z1 + z2)^2
	g.f.sub(t[0], t[0], t[7])   // (z1 + z2)^2 - z1z1
	g.f.sub(t[0], t[0], t[8])   // (z1 + z2)^2 - z1z1 - z2z2
	g.f.mul(r[2], t[0], t[1])   // z3 = ((z1 + z2)^2 - z1z1 - z2z2) * h
	return r
}

func (g *g23) double(r, p *pointG23) *pointG23 {
	if g.f.equal(g.a, g.f.zero()) {
		return g.doubleZeroA(r, p)
	}
	return g.doubleNonZeroA(r, p)
}

func (g *g23) doubleNonZeroA(r, p *pointG23) *pointG23 {
	// http://www.hyperelliptic.org/EFD/gp/auto-shortw-jacobian.html#doubling-dbl-2007-bl
	if g.isZero(p) {
		g.copy(r, p)
		return r
	}
	t := g.t
	g.f.square(t[0], p[0])    // xx
	g.f.square(t[1], p[1])    // yy
	g.f.square(t[3], p[2])    // zz
	g.f.add(t[2], p[1], p[2]) // y1 + z1
	g.f.square(t[2], t[2])    // (y1 + z1)^2
	g.f.sub(t[2], t[2], t[1]) // (y1 + z1)^2-yy
	g.f.sub(r[2], t[2], t[3]) // z3=(y1 + z1)^2 - yy - zz
	g.f.add(t[2], p[0], t[1]) // x1 + yy
	g.f.square(t[1], t[1])    // yyyy
	g.f.square(t[2], t[2])    // (x1 + yy)^2
	g.f.sub(t[2], t[2], t[0]) // (x1 + yy)^2-xx
	g.f.sub(t[2], t[2], t[1]) // (x2 + yy)^2 - xx - yyyy
	g.f.double(t[2], t[2])    // s = 2((x2 + yy)^2 - xx - yyyy)
	g.f.double(t[4], t[0])    // 2xx
	g.f.add(t[0], t[0], t[4]) // 3xx
	g.f.square(t[3], t[3])    // zz^2
	g.f.mul(t[3], g.a, t[3])  // zz^2a
	g.f.add(t[0], t[3], t[0]) // m = 3xx + zz^2a
	g.f.square(t[3], t[0])    // m^2
	g.f.double(t[4], t[2])    // 2s
	g.f.sub(t[3], t[3], t[4]) // t = m^2 - 2s
	g.f.copy(r[0], t[3])      // x3 = t
	g.f.sub(t[2], t[2], t[3]) // s - t
	g.f.mul(t[0], t[0], t[2]) // m * (s - t)
	g.f.double(t[1], t[1])    //
	g.f.double(t[1], t[1])    //
	g.f.double(t[1], t[1])    // 8yyyy
	g.f.sub(r[1], t[0], t[1]) // y3 = m * (s - t) - 8yyyy

	return r
}

func (g *g23) doubleZeroA(r, p *pointG23) *pointG23 {
	// http://www.hyperelliptic.org/EFD/gp/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	if g.isZero(p) {
		g.copy(r, p)
		return r
	}
	t := g.t
	g.f.square(t[0], p[0])    // a = x^2
	g.f.square(t[1], p[1])    // b = y^2
	g.f.square(t[2], t[1])    // c = b^2
	g.f.add(t[1], p[0], t[1]) // b + x1
	g.f.square(t[1], t[1])    // (b + x1)^2
	g.f.sub(t[1], t[1], t[0]) // (b + x1)^2 - a
	g.f.sub(t[1], t[1], t[2]) // (b + x1)^2 - a - c
	g.f.double(t[1], t[1])    // d = 2((b+x1)^2 - a - c)
	g.f.double(t[3], t[0])    // 2a
	g.f.add(t[0], t[3], t[0]) // e = 3a
	g.f.square(t[4], t[0])    // f = e^2
	g.f.double(t[3], t[1])    // 2d
	g.f.sub(r[0], t[4], t[3]) // x3 = f - 2d
	g.f.sub(t[1], t[1], r[0]) // d-x3
	g.f.double(t[2], t[2])    //
	g.f.double(t[2], t[2])    //
	g.f.double(t[2], t[2])    // 8c
	g.f.mul(t[0], t[0], t[1]) // e * (d - x3)
	g.f.sub(t[1], t[0], t[2]) // x3 = e * (d - x3) - 8c
	g.f.mul(t[0], p[1], p[2]) // y1 * z1
	g.f.copy(r[1], t[1])      //
	g.f.double(r[2], t[0])    // z3 = 2(y1 * z1)
	return r
}

func (g *g23) neg(r, p *pointG23) *pointG23 {
	g.f.copy(r[0], p[0])
	g.f.neg(r[1], p[1])
	g.f.copy(r[2], p[2])
	return r
}

func (g *g23) sub(c, a, b *pointG23) *pointG23 {
	d := g.newPoint()
	g.neg(d, b)
	g.add(c, a, d)
	return c
}

func (g *g23) mulScalar(c, p *pointG23, e *big.Int) *pointG23 {
	q, n := g.newPoint(), g.newPoint()
	g.copy(n, p)
	l := e.BitLen()
	for i := 0; i < l; i++ {
		if e.Bit(i) == 1 {
			g.add(q, q, n)
		}
		g.double(n, n)
	}
	g.copy(c, q)
	return c
}

func (g *g23) checkCorrectSubGroup(c, p *pointG23) *pointG23 {
	return g.wnafMul(c, p, g.q)
}

func (g *g23) wnafMul(c, p *pointG23, e *big.Int) *pointG23 {
	windowSize := uint(3)
	precompTable := make([]*pointG23, (1 << (windowSize - 1)))
	for i := 0; i < len(precompTable); i++ {
		precompTable[i] = g.newPoint()
	}
	var indexForPositive uint64
	indexForPositive = (1 << (windowSize - 2))
	g.copy(precompTable[indexForPositive], p)
	g.neg(precompTable[indexForPositive-1], p)
	doubled, precomp := g.newPoint(), g.newPoint()
	g.double(doubled, p)
	g.copy(precomp, p)
	for i := uint64(1); i < indexForPositive; i++ {
		g.add(precomp, precomp, doubled)
		g.copy(precompTable[indexForPositive+i], precomp)
		g.neg(precompTable[indexForPositive-1-i], precomp)
	}
	wnaf := wnaf(e, windowSize)
	q := g.zero()
	l := len(wnaf)
	found := false
	var idx uint64
	for i := l - 1; i >= 0; i-- {
		if found {
			g.double(q, q)
		}
		if wnaf[i] != 0 {
			found = true
			if wnaf[i] > 0 {
				idx = uint64(wnaf[i] >> 1)
				g.add(q, q, precompTable[indexForPositive+idx])
			} else {
				idx = uint64(((0 - wnaf[i]) >> 1))
				g.add(q, q, precompTable[indexForPositive-1-idx])
			}
		}
	}
	g.copy(c, q)
	return c
}

func (g *g23) multiExp(r *pointG23, points []*pointG23, powers []*big.Int) (*pointG23, error) {
	if len(points) != len(powers) {
		return nil, fmt.Errorf("point and scalar vectors should be in same length")
	}
	var c uint = 3
	if len(powers) > 32 {
		c = uint(math.Ceil(math.Log10(float64(len(powers)))))
	}
	bucket_size, numBits := (1<<c)-1, g.q.BitLen()
	windows := make([]*pointG23, numBits/int(c)+1)
	bucket := make([]*pointG23, bucket_size)
	acc, sum, zero := g.zero(), g.zero(), g.zero()
	s := new(big.Int)
	for i, m := 0, 0; i <= numBits; i, m = i+int(c), m+1 {
		for i := 0; i < bucket_size; i++ {
			bucket[i] = g.newPoint() // TODO: do it in a make or new func
		}
		for j := 0; j < len(powers); j++ {
			s = powers[j]
			index := s.Uint64() & uint64(bucket_size)
			if index != 0 {
				g.add(bucket[index-1], bucket[index-1], points[j])
			}
			s.Rsh(s, c)
		}
		g.copy(acc, zero)
		g.copy(sum, zero)
		for k := bucket_size - 1; k >= 0; k-- {
			g.add(sum, sum, bucket[k])
			g.add(acc, acc, sum)
		}
		windows[m] = g.zero()
		g.copy(windows[m], acc)
	}
	g.copy(acc, zero)
	for i := len(windows) - 1; i >= 0; i-- {
		for j := 0; j < int(c); j++ {
			g.double(acc, acc)
		}
		g.add(acc, acc, windows[i])
	}
	g.copy(r, acc)
	return r, nil
}
