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
func TestG2MulPoint(t *testing.T) {
	file := "test_vectors/custom/256.json"
	v, err := newTestVectorJSONFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	in, expected, err := v.makeG2MulBinary()
	if err != nil {
		t.Fatal(err)
	}

	api := new(g2Api)
	actual, err := api.mulPoint(in)
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

	api := new(pairingApi)
	actual, err := api.pair(in)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Logf("actual %x\n", actual)
		t.Logf("expected %x\n", expected)
		t.Fatalf("not equal")
	}
}
