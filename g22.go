package eip

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

type pointG22 [3]*fe2

type g22 struct {
	f   *fq2
	a   *fe2
	b   *fe2
	q   *big.Int
	t   [9]*fe2
	inf *pointG22
}

func newG22(f *fq2, a, b *fe2, q *big.Int) (*g22, error) {
	t := [9]*fe2{}
	for i := 0; i < 9; i++ {
		t[i] = f.zero()
	}
	g := &g22{
		f:   f,
		a:   f.new(),
		b:   f.new(),
		q:   new(big.Int).Set(q),
		t:   t,
		inf: &pointG22{f.zero(), f.one(), f.zero()},
	}
	if a == nil || b == nil {
		return nil, errors.New("a and b values should be provided")
	}
	if g.f.isZero(b) || (g.f.isZero(a) && g.f.isZero(b)) {
		return nil, errors.New("Curve shape is not supported")
	}
	g.f.copy(g.a, a)
	g.f.copy(g.b, b)
	return g, nil
}

func (g *g22) newPoint() *pointG22 {
	return &pointG22{g.f.new(), g.f.new(), g.f.new()}
}

func (g *g22) fromBytes(in []byte) (*pointG22, error) {
	byteLen := g.f.byteSize()
	if len(in) < 2*byteLen {
		return nil, fmt.Errorf("input string should be equal or larger than %d given %d", 2*byteLen, len(in))
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

func (g *g22) toBytes(p *pointG22) []byte {
	l := g.f.f.limbSize * 8 * 2
	out := make([]byte, 2*l)
	a := g.newPoint()
	g.affine(a, p)
	copy(out[:l], g.f.toBytes(a[0]))
	copy(out[l:], g.f.toBytes(a[1]))
	return out
}

func (g *g22) toBytesDense(p *pointG22) []byte {
	a := g.newPoint()
	g.affine(a, p)
	a0 := g.f.toBytesDense(a[0])
	a1 := g.f.toBytesDense(a[1])
	byteLen := len(a0)
	out := make([]byte, 2*byteLen)
	copy(out[:byteLen], a0)
	copy(out[byteLen:], a1)
	return out
}

func (g *g22) copy(q, p *pointG22) *pointG22 {
	g.f.copy(q[0], p[0])
	g.f.copy(q[1], p[1])
	g.f.copy(q[2], p[2])
	return q
}

func (g *g22) affine(r, p *pointG22) *pointG22 {
	if g.isZero(p) {
		g.f.copy(r[0], g.f.zero())
		g.f.copy(r[1], g.f.zero())
		g.f.copy(r[2], g.f.zero())
		return r
	}
	if g.f.equal(p[2], g.f.one()) {
		g.copy(r, p)
		return r
	}
	t := g.t
	if ok := g.f.inverse(t[0], p[2]); !ok {
		g.f.copy(r[0], g.f.zero())
		g.f.copy(r[1], g.f.zero())
		g.f.copy(r[2], g.f.zero())
		return r
	}
	q := g.newPoint()
	g.f.square(t[1], t[0])
	g.f.mul(q[0], p[0], t[1])
	g.f.mul(t[0], t[0], t[1])
	g.f.mul(q[1], p[1], t[0])
	g.f.copy(q[2], g.f.one())
	g.copy(r, q)
	return q
}

func (g *g22) toString(p *pointG22) string {
	return fmt.Sprintf("%s\n%s\n%s", g.f.toString(p[0]), g.f.toString(p[1]), g.f.toString(p[2]))
}

func (g *g22) toStringNoTransform(p *pointG22) string {
	return fmt.Sprintf("%s\n%s\n%s", g.f.toStringNoTransform(p[0]), g.f.toStringNoTransform(p[1]), g.f.toStringNoTransform(p[2]))
}

func (g *g22) zero() *pointG22 {
	p := g.newPoint()
	g.copy(p, g.inf)
	return p
}

func (g *g22) isZero(p *pointG22) bool {
	return g.f.isZero(p[2])
}

func (g *g22) equal(p1, p2 *pointG22) bool {
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

func (g *g22) isOnCurve(p *pointG22) bool {
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

func (g *g22) add(r, p1, p2 *pointG22) *pointG22 {
	// http://www.hyperelliptic.org/EFD/gp/auto-shortw-jacobian-0.html#addition-add-2007-bl
	if g.isZero(p2) {
		if g.isZero(p1) {
			g.copy(r, g.inf)
		} else {
			g.copy(r, p1)
		}
		return r
	}
	if g.isZero(p1) {
		if g.isZero(p2) {
			g.copy(r, g.inf)
		} else {
			g.copy(r, p2)
		}
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

func (g *g22) double(r, p *pointG22) *pointG22 {
	if g.f.equal(g.a, g.f.zero()) {
		return g.doubleZeroA(r, p)
	}
	return g.doubleNonZeroA(r, p)
}

func (g *g22) doubleNonZeroA(r, p *pointG22) *pointG22 {
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

func (g *g22) doubleZeroA(r, p *pointG22) *pointG22 {
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

func (g *g22) neg(r, p *pointG22) *pointG22 {
	if g.isZero(p) {
		g.copy(r, g.inf)
		return r
	}
	g.copy(r, p)
	g.f.neg(r[1], p[1])
	return r
}

func (g *g22) sub(c, a, b *pointG22) *pointG22 {
	d := g.newPoint()
	g.neg(d, b)
	g.add(c, a, d)
	return c
}

func (g *g22) mulScalar(c, p *pointG22, e *big.Int) *pointG22 {
	q, n := g.zero(), g.newPoint()
	g.copy(n, p)
	l := e.BitLen()
	found := false
	for i := l - 1; i >= 0; i-- {
		if found {
			g.double(q, q)
		} else {
			found = e.Bit(i) == 1
		}
		if e.Bit(i) == 1 {
			g.add(q, q, n)
		}
	}
	g.copy(c, q)
	return c
}

func (g *g22) checkCorrectSubgroup(p *pointG22) bool {
	c := g.newPoint()
	g.wnafMul(c, p, g.q)
	if g.equal(c, g.inf) {
		return true
	}
	return false
}

func (g *g22) wnafMul(c, p *pointG22, e *big.Int) *pointG22 {
	windowSize := uint(3)
	precompTable := make([]*pointG22, (1 << (windowSize - 1)))
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

func (g *g22) multiExp(r *pointG22, points []*pointG22, powers []*big.Int) (*pointG22, error) {
	if len(points) != len(powers) {
		return nil, fmt.Errorf("point and scalar vectors should be in same length")
	}
	var c uint32 = 3
	if len(powers) >= 32 {
		c = uint32(math.Ceil(math.Log10(float64(len(powers)))))
	}
	bucket_size, numBits := (1<<c)-1, uint32(g.q.BitLen())
	windows := make([]*pointG22, numBits/c+1)
	bucket := make([]*pointG22, bucket_size)
	acc, sum := g.newPoint(), g.newPoint()
	for i := 0; i < bucket_size; i++ {
		bucket[i] = g.newPoint()
	}
	mask := (uint64(1) << c) - 1
	j := 0
	var cur uint32
	for cur <= numBits {
		g.copy(acc, g.inf)
		bucket = make([]*pointG22, (1<<c)-1)
		for i := 0; i < len(bucket); i++ {
			bucket[i] = g.newPoint()
		}
		for i := 0; i < len(powers); i++ {
			s := newRepr(powers[i])
			index := uint(s[0] & mask)
			if index != 0 {
				g.add(bucket[index-1], bucket[index-1], points[i])
			}
			powers[i] = new(big.Int).Rsh(powers[i], uint(c))
		}

		g.copy(sum, g.inf)
		for i := len(bucket) - 1; i >= 0; i-- {
			g.add(sum, sum, bucket[i])
			g.add(acc, acc, sum)
		}
		windows[j] = g.newPoint()
		g.copy(windows[j], acc)
		j++
		cur += c
	}
	g.copy(acc, g.inf)
	for i := len(windows) - 1; i >= 0; i-- {
		for j := uint32(0); j < c; j++ {
			g.double(acc, acc)
		}
		g.add(acc, acc, windows[i])
	}
	g.copy(r, acc)
	return r, nil
}
