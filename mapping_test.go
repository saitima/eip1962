package eip

import (
	"math/big"
	"testing"
)

func TestMapToG1(t *testing.T) {
	modulus := fromHex(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	f, _ := newField(modulus)

	swuParams := computeSWUParamsForG1(f)
	isoParams := prepareIsogenyParamsForG1(f)

	u, _ := f.fromString("2a") // 42
	x, y := swuMapForG1(u, f, swuParams)
	xx, yy := applyIsogenyMapForG1(f, x, y, isoParams)

	a, _ := f.fromString("0x00")
	b, _ := f.fromString("0x04")
	q := new(big.Int).SetBytes(fromHex(48, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001"))
	g1, err := newG1(f, a, b, q)
	if err != nil {
		t.Fatal(err)
	}

	expected := g1.newPoint()
	ex, _ := f.fromString("0x14ef8eb66c02365fdca133fa506d522dad21919ea3d863d834ba2ddab4f4114273a166eba85d5abf8a17b59cfbec8a66")
	ey, _ := f.fromString("0x11b8e4a01bb5db119a48b9a747f5bd6a9da0be7576d6335f3f9eab0ce258211577c2637ce1707f93c2a4c4cfd0655f90")

	f.copy(expected[0], ex)
	f.copy(expected[1], ey)
	f.copy(expected[2], f.one)

	p := g1.newPoint()
	f.copy(p[0], xx)
	f.copy(p[1], yy)
	f.copy(p[2], f.one)

	if !g1.equal(expected, p) {
		t.Fatal("point is not equal to expected point")
	}
	if !g1.isOnCurve(p) {
		t.Fatal("p is not on the curve")
	}
}

func TestMapToG2(t *testing.T) {
	modulus := fromHex(48, "0x1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab")
	f, _ := newField(modulus)

	el, _ := f.fromString("2a") // 42
	// construct fq2
	nonResidue := f.new()
	f.neg(nonResidue, f.one)
	fq2, err := newFq2(f, f.toBytes(nonResidue))
	if err != nil {
		t.Fatal(err)
	}
	fq2.calculateFrobeniusCoeffs()
	u := fq2.new()
	f.copy(u[0], el)
	f.copy(u[1], el)

	swuParams := computeSWUParamsForG2(fq2)
	isoParams := prepareIsogenyParamsForG2(fq2)

	x, y := swuMapForG2(u, fq2, swuParams)
	xx, yy := applyIsogenyMapForG2(fq2, x, y, isoParams)

	a0, _ := f.fromString("0x00")
	a1, _ := f.fromString("0x00")
	b0, _ := f.fromString("0x04")
	b1, _ := f.fromString("0x04")
	a, b := fq2.new(), fq2.new()
	f.copy(a[0], a0)
	f.copy(a[1], a1)
	f.copy(b[0], b0)
	f.copy(b[1], b1)

	q := new(big.Int).SetBytes(fromHex(48, "0x73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001"))
	g2, err := newG22(fq2, a, b, q)
	if err != nil {
		t.Fatal(err)
	}

	expected := g2.newPoint()
	x0, _ := f.fromString("0x0a1915f637cd6361401db304517c1c6350008968d1cc26451e4b6aa5667ca6003450f94b7a526d01a0cfb447f2b89a74")
	x1, _ := f.fromString("0x06c08e47ee2e91fbabaafbca941bb901296ecb501fe130f5466e6b9cd29142983157b194d36246f7c7ea730e88ded395")
	y0, _ := f.fromString("0x18850ac995a5d8533e96d6cbec300ae158fada29912577c6aaba729f111f76abbb9bca56116d58e634ca174ffcf402eb")
	y1, _ := f.fromString("0x0d5b1f22b35c7bc66a34e65c03f09057cd38a0e035d50b4bc35f703f7ffc285fb5ad29cfde3a7cd29b5f1a8cd332a970")

	f.copy(expected[0][0], x0)
	f.copy(expected[0][1], x1)
	f.copy(expected[1][0], y0)
	f.copy(expected[1][1], y1)
	fq2.copy(expected[2], fq2.one())

	p := g2.newPoint()
	fq2.copy(p[0], xx)
	fq2.copy(p[1], yy)
	fq2.copy(p[2], fq2.one())

	if !g2.equal(expected, p) {
		t.Fatal("point is not equal to expected point")
	}
	if !g2.isOnCurve(p) {
		t.Fatal("p is not on the curve")
	}
}
