package eip

import "errors"

var (
	GAS_METERING_MODE = false
)

func _err(msg string) error {
	return errors.New(msg)
}

func _e(msgOrErr interface{}) ([]byte, error) {
	switch val := msgOrErr.(type) {
	case string:
		return nil, errors.New(val)
	case error:
		return nil, val
	default:
		return nil, errors.New(ERR_UNKNOWN)
	}
}

func apiDecodingErr(msgOrErr interface{}) ([]byte, error) {
	return _e(msgOrErr)
}

func apiExecErr(msgOrErr interface{}) ([]byte, error) {
	return _e(msgOrErr)
}

const (
	// base field
	ERR_BASE_FIELD_MODULUS_LENGTH_NOT_ENOUGH_BYTE = "Input is not long enough to get modulus length"
	ERR_MODULUS_LENGTH_ZERO                       = "Modulus length is zero"
	ERR_MODULUS_LENGTH_LARGE                      = "Encoded modulus length is too large"
	ERR_MODULUS_NOT_ENOUGH_BYTE                   = "Input is not long enough to get modulus"
	ERR_MODULUS_HIGHEST_BYTE                      = "In modulus encoding highest byte is zero"
	ERR_MODULUS_ZERO                              = "Modulus can not be zero"
	ERR_MODULUS_EVEN                              = "Modulus is even"
	ERR_MODULUS_LESS_THREE                        = "Modulus is less than 3"
	ERR_BASE_FIELD_CONSTRUCTION                   = "Failed to create prime field from modulus"
	ERR_INPUT_NOT_ENOUGH_FOR_FIELD_ELEMS          = "Input is not long enough"
	ERR_INPUT_NOT_ENOUGH_FOR_SCAKAR               = "Input is not long enough to get scalar"
	ERR_GROUP_ORDER_LENGTH_NOT_ENOUGH_BYTE        = "Input is not long enough to get group order length"
	ERR_GROUP_ORDER_LENGTH_ZERO                   = "Encoded group length is zero"
	ERR_GROUP_ORDER_LENGTH_LARGE                  = "Encoded group length is too large"
	ERR_GROUP_ORDER_NOT_ENOUGH_BYTE               = "Input is not long enough to get group order"
	ERR_GROUP_ORDER_ZERO                          = "Group order is zero"
	ERR_UNKNOWN                                   = "Unknown error"
	ERR_UNKNOWN_OPERATION                         = "Unknown operation type"
	ERR_GARBAGE_INPUT                             = "Input contains garbage at the end"
	ERR_POINT_NOT_ON_CURVE                        = "point  isn't on the curve"
	ERR_POINT0_NOT_ON_CURVE                       = "point 0 isn't on the curve"
	ERR_POINT1_NOT_ON_CURVE                       = "point 1 isn't on the curve"
	ERR_MULTIEXP_NUM_PAIRS_NOT_ENOUGH_BYTE        = "Input is not long enough to get number of pairs"
	ERR_MULTIEXP_NUM_PAIR_LENGTH                  = "Invalid number of pairs"
	ERR_MULTIEXP_NUM_PAIR_INPUT_LENGTH_NOT_MATCH  = "Input length is invalid for number of pairs"
	ERR_MULTIEXP_EMPTY_INPUT_PAIRS                = "Multiexp with empty input pairs"
	// g2
	ERR_G2_CANT_DECODE_EXT_DEGREE_LENGTH = "cant decode extension degree length"
	ERR_G2_UNEXPECTED_EXT_DEGREE         = "Extension degree expected to be 2 or 3"
	ERR_G2_UNKNOWN_OPERATION             = "Unknown G2 operation"
	// PAIRING
	ERR_EXT_FIELD_NON_RESIDUE_FP2_ZERO         = "Non-residue for Fp2 is zero"
	ERR_EXT_FIELD_NON_RESIDUE_FP2_RESIDUE      = "Non-residue for Fp2 is actually residue"
	ERR_EXT_FIELD_NON_RESIDUE_FP3_ZERO         = "Non-residue for Fp3 is zero"
	ERR_EXT_FIELD_NON_RESIDUE_FP3_RESIDUE      = "Non-residue for Fp3 is actually residue"
	ERR_EXT_FIELD_NON_RESIDUE_FP6_ZERO         = "Non-residue for Fp6 is zero"
	ERR_EXT_FIELD_NON_RESIDUE_FP6_RESIDUE      = "Non-residue for Fp6 is actually residue"
	ERR_EXT_FIELD_BASE_FROBENIUS_FOR_FP612     = "Can not make base precomputations for Fp6/Fp12 frobenius"
	ERR_EXT_FIELD_FROBENIUS_FOR_FP2            = "Can not calculate Frobenius coefficients for Fp2"
	ERR_EXT_FIELD_FROBENIUS_FOR_FP3            = "Can not calculate Frobenius coefficients for Fp3"
	ERR_EXT_FIELD_FROBENIUS_FOR_FP6            = "Can not calculate Frobenius coefficients for Fp6"
	ERR_EXT_FIELD_FROBENIUS_FOR_FP12           = "Can not calculate Frobenius coefficients for Fp12"
	ERR_PAIRING_FP2_NON_RESIDUE_NOT_INVERTIBLE = "Fp2 non-residue must be invertible"
	ERR_PAIRING_LOOP_COUNT_PARAM_ZERO          = "Loop count parameters can not be zero"
	ERR_PAIRING_NUM_PAIRS_NOT_ENOUGH_BYTE      = "Input is not long enough to get number of pairs"
	ERR_PAIRING_NUM_PAIRS_ZERO                 = "Zero pairs encoded"
	ERR_PAIRING_POINTG1_NOT_ON_CURVE           = "G1 point is not on curve"
	ERR_PAIRING_POINTG2_NOT_ON_CURVE           = "G2 point is not on curve"
	ERR_PAIRING_POINTG1_NOT_IN_SUBGROUP        = "G1 point is not in the expected subgroup"
	ERR_PAIRING_POINTG2_NOT_IN_SUBGROUP        = "G2 point is not in the expected subgroup"
	ERR_PAIRING_NO_RETURN_VALUE                = "Pairing engine returned no value"
	ERR_PAIRING_LOOP_PARAM_LENGTH              = "cant decode loop parameter's length"
	ERR_PAIRING_LOOP_PARAM_LENGTH_ZERO         = "Loop parameter scalar has zero length"
	ERR_PAIRING_LOOP_PARAM_LENGTH_LARGE        = "Loop parameter length scalar is too large for bit length"
	ERR_PAIRING_LOOP_PARAM_NOT_ENOUGH_BYTE     = "Input is not long enough to get loop parameter"
	ERR_PAIRING_LOOP_PARAM_TOP_BYTE_ZERO       = "Encoded loop parameter has zero top byte"
	ERR_PAIRING_LOOP_PARAM_LARGE               = "Loop parameter scalar is too large for bit length"
	ERR_PAIRING_TWIST_TYPE_NOT_ENOUGH_BYTE     = "Input is not long enough to get twist type"
	ERR_PAIRING_TWIST_TYPE_UNKNOWN             = "Unknown twist type supplied"
	ERR_PAIRING_EXP_SIGN_INVALID               = "exp is not encoded properly"
	ERR_PAIRING_EXP_SIGN_UNKNWON               = "Unknown parameter exp sign"
	ERR_PAIRING_BOOL_NOT_ENOUGH_BYTE           = "Input is not long enough to get boolean"
	ERR_PAIRING_BOOL_INVALID                   = "Boolean is not encoded properly"
	// Family specific
	ERR_BN_PAIRING_LOW_HAMMING_WEIGHT    = "|6*U + 2| has too large hamming weight"
	ERR_BN_PAIRING_A_PARAMETER_NOT_ZERO  = "A parameter must be zero for BN curve"
	ERR_BLS_PAIRING_A_PARAMETER_NOT_ZERO = "A parameter must be zero for BLS curve"
	ERR_BLS_PAIRING_LOW_HAMMING_WEIGHT   = "z has too large hamming weight"
	ERR_MNT_PAIRING_LOW_HAMMING_WEIGHT   = "x has too large hamming weight"
	ERR_MNT_EXP_NOT_ZERO                 = "loop count parameters can not be zero"
	ERR_MNT_EXPW0_NOT_ZERO               = "Final exp w0 loop count parameters can not be zero"
	ERR_MNT_EXPW1_NOT_ZERO               = "Final exp w1 loop count parameters can not be zero"
	ERR_MNT_INVALID_EXPW0                = "Exp_w0 sign is not encoded properly"
)
