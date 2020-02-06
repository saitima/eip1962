package eip

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

func bytes_(size int, hexStrs ...string) []byte {
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

func padHex(value []byte) []byte {
	requiredPad := len(value) % 8
	if requiredPad != 0 {
		padLen := (8 - requiredPad)
		value = append(make([]byte, padLen), value...)
		for i := 0; i < padLen; i++ {
			value[i] = 0x00
		}
	}
	// TODO: remove after for fuzz testing
	if USE_4LIMBS_FOR_LOWER_LIMBS && len(value) <= 32 {
		value = padBytes(value, 32)
	}
	return value
}

func reverse(in []byte) []byte {
	l := len(in)
	out := make([]byte, l)
	for i := l - 1; i >= 0; i-- {
		a := in[i]
		b := ((a & 0xf0) >> 4) | ((a & 0x0f) << 4)
		out[l-1-i] = b
	}
	return out
}

func isBigZero(a *big.Int) bool {
	if big.NewInt(0).Cmp(a) == 0 {
		return true
	}
	return false
}

func isBigLowerThan(a *big.Int, b int) bool {
	if a.Cmp(big.NewInt(int64(b))) == -1 {
		return true
	}
	return false
}

func isBigEven(a *big.Int) bool {
	if a.Bit(0) == 0 {
		return true
	}
	return false
}

func _error(str string) error {
	return errors.New(str)
}
