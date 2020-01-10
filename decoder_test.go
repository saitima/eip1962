package fp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
)

type TestVectorG1ScalarMultPair struct {
	A      string `json:"a"`
	Gx     string `json:"g_x"`
	Gy     string `json:"g_y"`
	Hx     string `json:"h_x"`
	Hy     string `json:"h_y"`
	Binary string `json:"scalar_mult_binary"`
}

type TestVectorG1MultiExp struct {
	Px     string `json:"expected_x"`
	Py     string `json:"expected_y"`
	Binary string `json:"binary"`
}

type TestVectorG2ScalarMultPair struct {
	A      string `json:"a"`
	Gx0    string `json:"g_x_0"`
	Gx1    string `json:"g_x_1"`
	Gy0    string `json:"g_y_0"`
	Gy1    string `json:"g_y_1"`
	Hx0    string `json:"h_x_0"`
	Hx1    string `json:"h_x_1"`
	Hy0    string `json:"h_y_0"`
	Hy1    string `json:"h_y_1"`
	Binary string `json:"scalar_mult_binary"`
}

type TestVectorJSON struct {
	buf               *bytes.Buffer
	fieldByteLen      int
	groupByteLen      int
	modulusBig        *big.Int
	N                 string                       `json:"n"`
	Q                 string                       `json:"q"`
	NonResidue        string                       `json:"non_residue"`
	NonResidue2_0     string                       `json:"quadratic_non_residue_0"`
	NonResidue2_1     string                       `json:"quadratic_non_residue_1"`
	A                 string                       `json:"A"`
	B                 string                       `json:"B"`
	G1x               string                       `json:"g1_x"`
	G1y               string                       `json:"g1_y"`
	CofactorG1        string                       `json:"cofactor_g1"`
	G1ScalarMultPairs []TestVectorG1ScalarMultPair `json:"g1_scalar_mult_test_vectors"`
	G1MultiExp        TestVectorG1MultiExp         `json:"g1_multiexp_test_vector"`
	IsDType           string                       `json:"is_D_type"`
	A_twist_0         string                       `json:"A_twist_0"`
	A_twist_1         string                       `json:"A_twist_1"`
	B_twist_0         string                       `json:"B_twist_0"`
	B_twist_1         string                       `json:"B_twist_1"`
	G2x0              string                       `json:"g2_x_0"`
	G2x1              string                       `json:"g2_x_1"`
	G2y0              string                       `json:"g2_y_0"`
	G2y1              string                       `json:"g2_y_1"`
	CofactorG2        string                       `json:"cofactor_g2"`
	G2ScalarMultPairs []TestVectorG2ScalarMultPair `json:"g2_scalar_mult_test_vectors"`
	R                 string                       `json:"r"`
	T                 string                       `json:"t"`
	X                 string                       `json:"x"`
}

func ceilBitLen(n int) int {
	return ((big.NewInt(int64(n)).BitLen() / 64) + 1) * 64
}

func newTestVectorJSONFromFile(file string) (*TestVectorJSON, error) {
	v := &TestVectorJSON{}
	v.buf = bytes.NewBuffer([]byte{})
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return nil, err
	}
	v.fieldByteLen = (len(v.N) - 2) / 2
	v.groupByteLen = (len(v.Q) - 2) / 2

	modulusBig, ok := new(big.Int).SetString(v.N[2:], 16)
	if !ok {
		panic("invalid modulus")
	}
	v.modulusBig = modulusBig
	if v.NonResidue[:1] == "-" {
		nonResidue, ok := new(big.Int).SetString(v.NonResidue[3:], 16)
		if !ok {
			panic("invalid nonresidue")
		}
		v.NonResidue = new(big.Int).Sub(v.modulusBig, nonResidue).Text(16)
	}
	return v, nil
}

func (v *TestVectorJSON) encode(str string) []byte {
	return bytes_(v.fieldByteLen, str)
}

func (v *TestVectorJSON) makeBaseFieldBinary() {
	v.buf.WriteByte(byte(v.fieldByteLen))
	// - Field modulus
	v.buf.Write(v.encode(v.N))
}

func (v *TestVectorJSON) makeG1GroupBinary() {
	// - Curve A
	v.buf.Write(v.encode(v.A))
	// - Curve B
	v.buf.Write(v.encode(v.B))
	// - Group order byte length
	v.buf.WriteByte(byte(v.groupByteLen))
	v.buf.Write(bytes_(v.groupByteLen, v.Q))
}

func (v *TestVectorJSON) makeG1PointBinary() {
	// - Point
	v.buf.Write(v.encode(v.G1ScalarMultPairs[0].Gx))
	v.buf.Write(v.encode(v.G1ScalarMultPairs[0].Gy))
}

func (v *TestVectorJSON) makeG1ScalarBinary() {
	// - Scalar
	v.buf.Write(bytes_(v.groupByteLen, v.G1ScalarMultPairs[0].A))
}

func (v *TestVectorJSON) makeG1MulBinary() ([]byte, []byte, error) {
	v.makeBaseFieldBinary()
	v.makeG1GroupBinary()
	v.makeG1PointBinary()
	v.makeG1ScalarBinary()
	input := make([]byte, len(v.buf.Bytes()))
	copy(input, v.buf.Bytes())
	if len(input) < (5*v.fieldByteLen)+(2*v.groupByteLen)+2 {
		return nil, nil, fmt.Errorf("cant assemble input data for g1 scalar mul")
	}
	v.buf.Reset()

	// - Output
	v.buf.Write(v.encode(v.G1ScalarMultPairs[0].Hx))
	v.buf.Write(v.encode(v.G1ScalarMultPairs[0].Hy))
	output := make([]byte, len(v.buf.Bytes()))
	copy(output, v.buf.Bytes())
	if len(output) < (2 * 32) {
		return nil, nil, fmt.Errorf("cant assemble output data for g1 scalar mul")
	}
	return input, output, nil
}

func (v *TestVectorJSON) makeExtension2Field() {
	v.buf.WriteByte(byte(EXTENSION_TWO_DEGREE))
	v.buf.Write(v.encode(v.NonResidue))
	v.buf.Write(v.encode(v.A_twist_0))
	v.buf.Write(v.encode(v.A_twist_1))
	v.buf.Write(v.encode(v.B_twist_0))
	v.buf.Write(v.encode(v.B_twist_1))
	g2GroupOrderLen := (len(v.N) - 2) / 2
	v.buf.WriteByte(byte(g2GroupOrderLen))
	v.buf.Write(v.encode(v.R))

}

func (v *TestVectorJSON) makeG2PointBinary() {
	// - Point
	v.buf.Write(v.encode(v.G2ScalarMultPairs[0].Gx0))
	v.buf.Write(v.encode(v.G2ScalarMultPairs[0].Gx1))
	v.buf.Write(v.encode(v.G2ScalarMultPairs[0].Gy0))
	v.buf.Write(v.encode(v.G2ScalarMultPairs[0].Gy1))
}

func (v *TestVectorJSON) makeG2ScalarBinary() {
	// - Scalar
	v.buf.Write(bytes_(v.groupByteLen, v.G2ScalarMultPairs[0].A))
}

func (v *TestVectorJSON) makeG2MulBinary() ([]byte, []byte, error) {
	v.makeBaseFieldBinary()
	v.makeExtension2Field()
	v.makeG2PointBinary()
	v.makeG2ScalarBinary()
	input := make([]byte, len(v.buf.Bytes()))
	copy(input, v.buf.Bytes())
	if len(input) < (5*v.fieldByteLen)+(2*v.groupByteLen)+2 {
		return nil, nil, fmt.Errorf("cant assemble input data for g2 scalar mul")
	}
	v.buf.Reset()
	// - Output
	v.buf.Write(v.encode(v.G2ScalarMultPairs[0].Hx0))
	v.buf.Write(v.encode(v.G2ScalarMultPairs[0].Hx1))
	v.buf.Write(v.encode(v.G2ScalarMultPairs[0].Hy0))
	v.buf.Write(v.encode(v.G2ScalarMultPairs[0].Hy1))
	output := make([]byte, len(v.buf.Bytes()))
	copy(output, v.buf.Bytes())
	if len(output) < (4 * 32) {
		return nil, nil, fmt.Errorf("cant assemble output data for g2 scalar mul")
	}
	return input, output, nil
}

func (v *TestVectorJSON) makeBNPairingBinary() ([]byte, []byte, error) {
	v.buf.WriteByte(0x01)                  // curve type
	v.makeBaseFieldBinary()                // base field
	v.makeG1GroupBinary()                  // g1
	v.buf.Write(v.encode(v.NonResidue))    // non residue
	v.buf.Write(v.encode(v.NonResidue2_0)) // quadratic non residue
	v.buf.Write(v.encode(v.NonResidue2_1))
	if v.IsDType == "True" {
		v.buf.WriteByte(0x02) // twist type D
	} else {
		v.buf.WriteByte(0x01) // twist type M
	}

	if v.X[:1] == "-" {
		length := len(v.X[3:]) / 2
		v.buf.WriteByte(byte(length))        // x length
		v.buf.Write(bytes_(length, v.X[1:])) // x encoded
		v.buf.WriteByte(0x01)                // sign of x
	} else {
		length := len(v.X[2:]) / 2
		v.buf.WriteByte(byte(length))    // x length
		v.buf.Write(bytes_(length, v.X)) // x encoded
		v.buf.WriteByte(0x00)
	}
	v.buf.WriteByte(byte(0x02)) // num pairs
	// e(P, Q)*e(-P, Q)=1
	// first pair
	v.buf.Write(v.encode(v.G1x))
	v.buf.Write(v.encode(v.G1y))
	v.buf.Write(v.encode(v.G2x0))
	v.buf.Write(v.encode(v.G2x1))
	v.buf.Write(v.encode(v.G2y0))
	v.buf.Write(v.encode(v.G2y1))

	// second pair
	pYBuf := v.encode(v.G1y)
	pY := new(big.Int).SetBytes(pYBuf)
	pYStr := new(big.Int).Sub(v.modulusBig, pY).Text(16) // -p.y
	v.buf.Write(v.encode(v.G1x))
	v.buf.Write(v.encode(pYStr))
	v.buf.Write(v.encode(v.G2x0))
	v.buf.Write(v.encode(v.G2x1))
	v.buf.Write(v.encode(v.G2y0))
	v.buf.Write(v.encode(v.G2y1))

	if len(v.buf.Bytes()) < (16*32)+(6*1)+(len(v.X)-2)/2 {
		return nil, nil, errors.New("can't assemble pairing binary data")
	}
	input := v.buf.Bytes()
	output := []byte{0x01}
	return input, output, nil
}
