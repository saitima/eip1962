package eip

import "errors"

const (
	// api operations
	OPERATION_G1_ADD      = 0x01
	OPERATION_G1_MUL      = 0x02
	OPERATION_G1_MULTIEXP = 0x03
	OPERATION_G2_ADD      = 0x04
	OPERATION_G2_MUL      = 0x05
	OPERATION_G2_MULTIEXP = 0x06
	OPERATION_BLS12PAIR   = 0x07
	OPERATION_BNPAIR      = 0x08
	OPERATION_MNT4PAIR    = 0x09
	OPERATION_MNT6PAIR    = 0x0a
	// flags
	USE_4LIMBS_FOR_LOWER_LIMBS  = true
	TWIST_M, TWIST_D            = 0x01, 0x02
	NEGATIVE_EXP, POSITIVE_EXP  = 0x01, 0x00
	BOOLEAN_FALSE, BOOLEAN_TRUE = 0x00, 0x01
)

type API struct{}

func newAPI() *API {
	return &API{}
}

func (api *API) Run(opType int, in []byte) ([]byte, error) {
	decoder := newDecoder(in)
	var runner runner
	var err error
	switch opType {
	case OPERATION_G1_ADD:
		runner, err = decoder.g1AddRunner()
	case OPERATION_G1_MUL:
		runner, err = decoder.g1MulRunner()
	case OPERATION_G1_MULTIEXP:
		runner, err = decoder.g1MultiExpRunner()
	case OPERATION_G2_ADD:
		runner, err = decoder.g2AddRunner()
	case OPERATION_G2_MUL:
		runner, err = decoder.g2MulRunner()
	case OPERATION_G2_MULTIEXP:
		runner, err = decoder.g2MultiExpRunner()
	case OPERATION_BLS12PAIR:
		runner, err = decoder.blsRunner()
	case OPERATION_BNPAIR:
		runner, err = decoder.bnRunner()
	case OPERATION_MNT4PAIR:
		runner, err = decoder.mnt4Runner()
	case OPERATION_MNT6PAIR:
		runner, err = decoder.mnt6Runner()
	default:
		err = errors.New(ERR_UNKNOWN_OPERATION)
	}
	if err != nil {
		return nil, err
	}
	return runner.run()
}
