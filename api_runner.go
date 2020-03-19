package eip

import (
	"errors"
	"math/big"
)

type runner interface {
	run() ([]byte, error)
}

type g1AddRunner struct {
	g1 *g1
	p1 *pointG1
	p2 *pointG1
}

func newG1AddRunner(g1 *g1, p1 *pointG1, p2 *pointG1) *g1AddRunner {
	return &g1AddRunner{g1, p1, p2}
}

func (runner *g1AddRunner) run() ([]byte, error) {
	g1 := runner.g1
	p1 := runner.p1
	p2 := runner.p2
	r := g1.newPoint()
	g1.add(r, p1, p2)
	return g1.toBytesDense(r), nil
}

type g1MulRunner struct {
	g1 *g1
	p  *pointG1
	s  *big.Int
}

func newG1MulRunner(g1 *g1, p *pointG1, s *big.Int) *g1MulRunner {
	return &g1MulRunner{g1, p, s}
}

func (runner *g1MulRunner) run() ([]byte, error) {
	g1 := runner.g1
	p := runner.p
	s := runner.s
	r := g1.newPoint()
	g1.mulScalar(r, p, s)
	return g1.toBytesDense(r), nil
}

type g1MultiExpRunner struct {
	g1      *g1
	points  []*pointG1
	scalars []*big.Int
}

func newG1MultiExpRunner(g1 *g1, p []*pointG1, s []*big.Int) *g1MultiExpRunner {
	return &g1MultiExpRunner{g1, p, s}
}

func (runner *g1MultiExpRunner) run() ([]byte, error) {
	g1 := runner.g1
	points := runner.points
	scalars := runner.scalars
	r := g1.newPoint()
	g1.multiExp(r, points, scalars)
	return g1.toBytesDense(r), nil
}

type g22AddRunner struct {
	g22 *g22
	p1  *pointG22
	p2  *pointG22
}

func newG22AddRunner(g22 *g22, p1 *pointG22, p2 *pointG22) *g22AddRunner {
	return &g22AddRunner{g22, p1, p2}
}

func (runner *g22AddRunner) run() ([]byte, error) {
	g22 := runner.g22
	p1 := runner.p1
	p2 := runner.p2
	r := g22.newPoint()
	g22.add(r, p1, p2)
	return g22.toBytesDense(r), nil
}

type g22MulRunner struct {
	g22 *g22
	p   *pointG22
	s   *big.Int
}

func newG22MulRunner(g22 *g22, p *pointG22, s *big.Int) *g22MulRunner {
	return &g22MulRunner{g22, p, s}
}

func (runner *g22MulRunner) run() ([]byte, error) {
	g22 := runner.g22
	p := runner.p
	s := runner.s
	r := g22.newPoint()
	g22.mulScalar(r, p, s)
	return g22.toBytesDense(r), nil
}

type g22MultiExpRunner struct {
	g22     *g22
	points  []*pointG22
	scalars []*big.Int
}

func newG22MultiExpRunner(g22 *g22, p []*pointG22, s []*big.Int) *g22MultiExpRunner {
	return &g22MultiExpRunner{g22, p, s}
}

func (runner *g22MultiExpRunner) run() ([]byte, error) {
	g22 := runner.g22
	points := runner.points
	scalars := runner.scalars
	r := g22.newPoint()
	g22.multiExp(r, points, scalars)
	return g22.toBytesDense(r), nil
}

type g23AddRunner struct {
	g23 *g23
	p1  *pointG23
	p2  *pointG23
}

func newG23AddRunner(g23 *g23, p1 *pointG23, p2 *pointG23) *g23AddRunner {
	return &g23AddRunner{g23, p1, p2}
}

func (runner *g23AddRunner) run() ([]byte, error) {
	g23 := runner.g23
	p1 := runner.p1
	p2 := runner.p2
	r := g23.newPoint()
	g23.add(r, p1, p2)
	return g23.toBytesDense(r), nil
}

type g23MulRunner struct {
	g23 *g23
	p   *pointG23
	s   *big.Int
}

func newG23MulRunner(g23 *g23, p *pointG23, s *big.Int) *g23MulRunner {
	return &g23MulRunner{g23, p, s}
}

func (runner *g23MulRunner) run() ([]byte, error) {
	g23 := runner.g23
	p := runner.p
	s := runner.s
	r := g23.newPoint()
	g23.mulScalar(r, p, s)
	return g23.toBytesDense(r), nil
}

type g23MultiExpRunner struct {
	g23     *g23
	points  []*pointG23
	scalars []*big.Int
}

func newG23MultiExpRunner(g23 *g23, p []*pointG23, s []*big.Int) *g23MultiExpRunner {
	return &g23MultiExpRunner{g23, p, s}
}

func (runner *g23MultiExpRunner) run() ([]byte, error) {
	g23 := runner.g23
	points := runner.points
	scalars := runner.scalars
	r := g23.newPoint()
	g23.multiExp(r, points, scalars)
	return g23.toBytesDense(r), nil
}

var (
	pairingError, pairingSuccess = []byte{0x00}, []byte{0x01}
)

type blsRunner struct {
	e  *blsInstance
	P1 []*pointG1
	P2 []*pointG22
}

func newBLSRunner(e *blsInstance, P1 []*pointG1, P2 []*pointG22) *blsRunner {
	return &blsRunner{e, P1, P2}
}

func (runner *blsRunner) run() ([]byte, error) {
	e := runner.e
	gt := e.gt()
	P1 := runner.P1
	P2 := runner.P2
	result, hasValue := e.multiPair(P1, P2)
	if !hasValue {
		return nil, errors.New(ERR_PAIRING_NO_RETURN_VALUE)
	}
	if !gt.equal(result, gt.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}

type bnRunner struct {
	e  *bnInstance
	P1 []*pointG1
	P2 []*pointG22
}

func newBNRunner(e *bnInstance, P1 []*pointG1, P2 []*pointG22) *bnRunner {
	return &bnRunner{e, P1, P2}
}

func (runner *bnRunner) run() ([]byte, error) {
	e := runner.e
	gt := e.gt()
	P1 := runner.P1
	P2 := runner.P2
	result, hasValue := e.multiPair(P1, P2)
	if !hasValue {
		return nil, errors.New(ERR_PAIRING_NO_RETURN_VALUE)
	}
	if !gt.equal(result, gt.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}

type mnt4Runner struct {
	e  *mnt4Instance
	P1 []*pointG1
	P2 []*pointG22
}

func newMNT4Runner(e *mnt4Instance, P1 []*pointG1, P2 []*pointG22) *mnt4Runner {
	return &mnt4Runner{e, P1, P2}
}

func (runner *mnt4Runner) run() ([]byte, error) {
	e := runner.e
	gt := e.gt()
	P1 := runner.P1
	P2 := runner.P2
	result, hasValue := e.multiPair(P1, P2)
	if !hasValue {
		return nil, errors.New(ERR_PAIRING_NO_RETURN_VALUE)
	}
	if !gt.equal(result, gt.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}

type mnt6Runner struct {
	e  *mnt6Instance
	P1 []*pointG1
	P2 []*pointG23
}

func newMNT6Runner(e *mnt6Instance, P1 []*pointG1, P2 []*pointG23) *mnt6Runner {
	return &mnt6Runner{e, P1, P2}
}

func (runner *mnt6Runner) run() ([]byte, error) {
	e := runner.e
	gt := e.gt()
	P1 := runner.P1
	P2 := runner.P2
	result, hasValue := e.multiPair(P1, P2)
	if !hasValue {
		return nil, errors.New(ERR_PAIRING_NO_RETURN_VALUE)
	}
	if !gt.equal(result, gt.one()) {
		return pairingError, nil
	}
	return pairingSuccess, nil
}
