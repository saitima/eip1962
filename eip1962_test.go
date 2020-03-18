package eip

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"
	"testing"
)

var fuz int = 1

var targetNumberOfLimb int = -1

var from = 1
var to = 16

func TestArch(t *testing.T) {
	answer := "Yes."
	if nonADXBMI2 {
		answer = "No."
	}
	fmt.Printf("Is using ADX backend extension? %s\n", answer)
}

func TestMain(m *testing.M) {
	_fuz := flag.Int("fuzz", 1, "# of iters")
	nol := flag.Int("nol", 0, "backend bit size")
	flag.Parse()
	fuz = *_fuz
	if *nol > 0 {
		targetNumberOfLimb = *nol
		if !(targetNumberOfLimb >= from && targetNumberOfLimb <= to) {
			panic(fmt.Sprintf("limb size %d not supported", targetNumberOfLimb))
		}
		from = targetNumberOfLimb
		to = targetNumberOfLimb
	}
	m.Run()
}

func randBytes(max *big.Int) []byte {
	return padBytes(randBig(max).Bytes(), resolveLimbSize(max.BitLen())*8)
}

func randBig(max *big.Int) *big.Int {
	bi, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return bi
}

func randScalar(max *big.Int) *big.Int {
	return randBig(max)
}

func debugBytes(a ...[]byte) {
	for _, b := range a {
		for i := (len(b) / 8) - 1; i > -1; i-- {
			fmt.Printf("0x%16.16x,\n", b[i*8:i*8+8])
		}
		fmt.Println()
	}
}

func fromHex(size int, hexStrs ...string) []byte {
	var out []byte
	if size > 0 {
		out = make([]byte, size*len(hexStrs))
	}
	for i := 0; i < len(hexStrs); i++ {
		hexStr := hexStrs[i]
		if hexStr[:2] == "0x" {
			hexStr = hexStr[2:]
		}
		if len(hexStr)%2 == 1 {
			hexStr = "0" + hexStr
		}
		bytes, err := hex.DecodeString(hexStr)
		if err != nil {
			panic(err)
		}
		if size <= 0 {
			out = append(out, bytes...)
		} else {
			if len(bytes) > size {
				panic(fmt.Sprintf("bad input string\ninput: %x\nsize: %d\nlenght: %d\n", bytes, size, len(bytes)))
			}
			offset := i*size + (size - len(bytes))
			copy(out[offset:], bytes)
		}
	}
	return out
}
