package eip

func encodeG1Point(out []byte, in []byte) {
	inPos := len(in) / 2
	length := len(out) / 2
	x := encodeFixedLen(length, in[:inPos])
	y := encodeFixedLen(length, in[inPos:])
	copy(out, append(x, y...))
}

func encodeG22Point(out []byte, in []byte) {
	inPos := len(in) / 2
	outPos := len(out) / 2
	encodeFq2FixedLen(out[:outPos], in[:inPos])
	encodeFq2FixedLen(out[outPos:], in[inPos:])
}

func encodeG23Point(out []byte, in []byte) {
	inPos := len(in) / 2
	outPos := len(out) / 2
	encodeFq3FixedLen(out[:outPos], in[:inPos])
	encodeFq3FixedLen(out[outPos:], in[inPos:])
}

func encodeFq2FixedLen(out []byte, in []byte) {
	pos := len(in) / 2
	length := len(out) / 2
	c0 := encodeFixedLen(length, in[:pos])
	c1 := encodeFixedLen(length, in[pos:])
	copy(out, append(c0, c1...))
}

func encodeFq3FixedLen(out []byte, in []byte) {
	pos := len(in) / 3
	length := len(out) / 3
	c0 := encodeFixedLen(length, in[:pos])
	c1 := encodeFixedLen(length, in[pos:2*pos])
	c2 := encodeFixedLen(length, in[2*pos:])
	copy(out, append(c0, append(c1, c2...)...))
}

func encodeFixedLen(modulusLen int, in []byte) []byte {
	if len(in) > modulusLen {
		// truncate
		tmp := reverse(in)
		return reverse(tmp[:modulusLen])
	} else if len(in) < modulusLen {
		// resize
		return padBytes(in, modulusLen)
	}
	return in
}
