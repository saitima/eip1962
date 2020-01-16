package fp

import (
	"bytes"
	"testing"
)

func TestG1MulPoint(t *testing.T) {
	file := "test_vectors/custom/256.json"
	v, err := newTestVectorJSONFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	in, expected, err := v.makeG1MulBinary()
	if err != nil {
		t.Fatal(err)
	}

	api := new(g1Api)
	actual, err := api.mulPoint(in)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Logf("actual %x\n", actual)
		t.Logf("expected %x\n", expected)
	}
}
func TestG22MulPoint(t *testing.T) {
	file := "test_vectors/custom/256.json"
	v, err := newTestVectorJSONFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	in, expected, err := v.makeG22MulBinary()
	if err != nil {
		t.Fatal(err)
	}

	api := new(API)
	actual, err := api.run(in)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Logf("actual %x\n", actual)
		t.Logf("expected %x\n", expected)
		t.Fatalf("not equal")
	}
}
func TestG23MulPoint(t *testing.T) {
	file := "test_vectors/custom/320_cubic.json"
	v, err := newTestVectorJSONFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	in, expected, err := v.makeG23MulBinary()
	if err != nil {
		t.Fatal(err)
	}

	api := new(API)
	actual, err := api.run(in)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Logf("actual %x\n", actual)
		t.Logf("expected %x\n", expected)
		t.Fatalf("not equal")
	}
}

func TestBNPairing(t *testing.T) {
	file := "test_vectors/custom/256.json"
	v, err := newTestVectorJSONFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	in, expected, err := v.makeBNPairingBinary()
	if err != nil {
		t.Fatal(err)
	}
	actual, err := new(API).run(in)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, expected) {
		t.Logf("actual %x\n", actual)
		t.Logf("expected %x\n", expected)
		t.Fatalf("not equal")
	}
}
func TestBLSPairing(t *testing.T) {
	file := "test_vectors/custom/384.json"
	v, err := newTestVectorJSONFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	in, expected, err := v.makeBLSPairingBinary()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := new(API).run(in)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Logf("actual %x\n", actual)
		t.Logf("expected %x\n", expected)
		t.Fatalf("not equal")
	}
}

func TestMNT4Pairing(t *testing.T) {
	file := "test_vectors/custom/320.json"
	v, err := newTestVectorJSONFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	in, expected, err := v.makeMNT4PairingBinary()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := new(API).run(in)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Logf("actual %x\n", actual)
		t.Logf("expected %x\n", expected)
		t.Fatalf("not equal")
	}
}

func TestMNT6Pairing(t *testing.T) {
	file := "test_vectors/custom/320_mnt6.json"
	v, err := newTestVectorJSONFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	in, expected, err := v.makeMNT6PairingBinary()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := new(API).run(in)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Logf("actual %x\n", actual)
		t.Logf("expected %x\n", expected)
		t.Fatalf("not equal")
	}
}
