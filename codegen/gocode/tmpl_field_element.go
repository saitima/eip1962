// +build ignore

package main

var fieldElementTemplates = []string{
	feTmplMarshal,
	feTmplUnmarshal,
	feTmplSetBig,
	feTmplSetUint,
	feTmplSetString,
	feTmplSet,
	feTmplBig,
	feTmplString,
	feTmplIsOdd,
	feTmplIsEven,
	feTmplIsZero,
	feTmplIsOne,
	feTmplCompare,
	feTmplEquals,
	feTmplRightSh,
	feTmplLeftSh,
	feTmplBit,
	feTmplBitLen,
	feTmplRand,
}

const (
	// 	feTmplDefine = `
	// type Fe{{ $N_BIT }} [ {{ $N_LIMB }} ]uint64
	// `

	feTmplMarshal = `
func (fe *Fe{{ $N_BIT }}) Marshal(out []byte) []byte {
var a int
for i := 0; i < {{ $N_LIMB }}; i++ {
a = {{ $N_LIMB }}*8 - i*8
out[a-1] = byte(fe[i])
out[a-2] = byte(fe[i] >> 8)
out[a-3] = byte(fe[i] >> 16)
out[a-4] = byte(fe[i] >> 24)
out[a-5] = byte(fe[i] >> 32)
out[a-6] = byte(fe[i] >> 40)
out[a-7] = byte(fe[i] >> 48)
out[a-8] = byte(fe[i] >> 56)}
return out}
`

	feTmplUnmarshal = `
func (fe *Fe{{ $N_BIT }}) Unmarshal(in []byte) *Fe{{ $N_BIT }} {
size := {{ $N_LIMB }} * 8
padded := make([]byte, size)
l := len(in)
if l >= size {
l = size}
copy(padded[size-l:], in[:])
var a int
for i := 0; i < {{ $N_LIMB }}; i++ {
a = size - i*8
fe[i] = uint64(padded[a-1]) | uint64(padded[a-2])<<8 |
uint64(padded[a-3])<<16 | uint64(padded[a-4])<<24 |
uint64(padded[a-5])<<32 | uint64(padded[a-6])<<40 |
uint64(padded[a-7])<<48 | uint64(padded[a-8])<<56}
return fe}
`

	feTmplSetBig = `
func (fe *Fe{{ $N_BIT }}) SetBig(a *big.Int) *Fe{{ $N_BIT }} {
return fe.Unmarshal(a.Bytes())}	
`

	// todo : add util func div by 4 or mul by
	feTmplBig = `
func (fe *Fe{{ $N_BIT }}) Big() *big.Int {
h := [{{ $N_LIMB }} * 8 ]byte{}
return new(big.Int).SetBytes(fe.Marshal(h[:]))
}`

	feTmplSetUint = `
func (fe *Fe{{ $N_BIT }}) SetUint(a uint64) *Fe{{ $N_BIT }} {
fe[0] = a
{{ range $x := iterUp 1 $N_LIMB }} fe[ {{ $x }} ] = 0 
{{ end }} return fe }
`

	feTmplSetString = `
func (fe *Fe{{ $N_BIT }}) SetString(s string) (*Fe{{ $N_BIT }}, error) {
if s[:2] == "0x" {
s = s[2:]}
h, err := hex.DecodeString(s)
if err != nil {
return nil, err}
return fe.Unmarshal(h), nil}
`

	feTmplSet = `
func (fe *Fe{{ $N_BIT }}) Set(fe2 *Fe{{ $N_BIT }}) *Fe{{ $N_BIT }} {
{{ range $x := iterUp 0 $N_LIMB }} fe[ {{ $x }} ] = fe2[ {{ $x }} ]
{{ end }} return fe }
`

	feTmplString = `
func (fe Fe{{ $N_BIT }}) String() (s string) {
for i := {{ decr $N_LIMB }}; i >= 0; i-- {
s = fmt.Sprintf("%s%16.16x", s, fe[i]) }
return "0x" + s }
`

	feTmplCompare = `
func (fe *Fe{{ $N_BIT }}) Cmp(fe2 *Fe{{ $N_BIT }}) int64 {
{{range $i := iterDown $N_LIMB }} 
if fe[{{ $i }}] > fe2[{{ $i }}] {
return 1
} else if fe[{{ $i }}] < fe2[{{ $i }}] {
return -1
}{{end}}
return 0 }
`

	feTmplIsEven = `
func (fe *Fe{{ $N_BIT }}) IsEven() bool {
var mask uint64 = 1
return fe[0]&mask == 0 }
`

	feTmplIsOdd = `
func (fe *Fe{{ $N_BIT }}) IsOdd() bool {
var mask uint64 = 1
return fe[0]&mask != 0 }
`

	feTmplIsOne = `
func (fe *Fe{{$N_BIT}}) IsOne() bool {
return 1 == fe[0] {{ range $x := iterUp 1 $N_LIMB }} && 0 == fe[ {{ $x }} ] {{ end }} }
`

	feTmplIsZero = `
func (fe *Fe{{$N_BIT}}) IsZero() bool {
return 0 == fe[0] {{ range $x := iterUp 1 $N_LIMB }} && 0 == fe[ {{ $x }} ] {{ end }} }
`

	feTmplEquals = `
func (fe *Fe{{$N_BIT}}) Equals(fe2 *Fe{{$N_BIT}}) bool {
return fe2[0] == fe[0] {{ range $x := iterUp 1 $N_LIMB }} && fe2[ {{ $x }} ] == fe[ {{ $x }} ] {{ end }} }
`

	feTmplRightSh = `
func (fe *Fe{{$N_BIT}}) div2(e uint64) {
{{ range $x := iterUp 1 $N_LIMB }}; fe[{{ decr $x }}] = fe[ {{ decr $x }} ]>>1 | fe[{{$x}}]<<63 ; {{ end }}
fe[{{ decr $N_LIMB }}] = fe[{{ decr $N_LIMB }}] >> 1 | e << 63 }
`

	feTmplLeftSh = `
func (fe *Fe{{$N_BIT}}) mul2() uint64 {
e := fe[{{ decr $N_LIMB }}] >> 63
{{ range $i := iterDown $N_LIMB }} ; {{if $i}} fe[ {{$i}} ] = fe[ {{$i}} ]<<1 | fe[ {{decr $i}} ]>>63 {{else}}fe[0] = fe[0] << 1{{end}}; {{ end }}
return e }
`

	feTmplBit = `
func (fe *Fe{{$N_BIT}}) bit(i int) bool {
k := i >> 6
i = i - k<<6
b := (fe[k] >> uint(i)) & 1
return b != 0 } 
`
	feTmplBitLen = `
func (fe *Fe{{$N_BIT}}) bitLen() int {
for i := len(fe) - 1; i >= 0; i-- {
if len := bits.Len64(fe[i]); len != 0 {
return len + 64*i}}
return 0}
`

	feTmplRand = `
func (f *Fe{{ $N_BIT }}) rand(max *Fe{{ $N_BIT }}, r io.Reader) error {
bitLen := bits.Len64(max[{{ decr $N_LIMB }}]) + ({{ $N_LIMB }} -1)*64
k := (bitLen + 7) / 8
b := uint(bitLen % 8)
if b == 0 {
b = 8 }
bytes := make([]byte, k)
for {
_, err := io.ReadFull(r, bytes)
if err != nil {
return err }
bytes[0] &= uint8(int(1<<b) - 1)
f.Unmarshal(bytes)
if f.Cmp(max) < 0 {
break } }
return nil }
`
)
