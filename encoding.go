package main

import "bytes"

type EncodingMode byte

// Mode conts
const (
	Upper EncodingMode = iota
	Lower
	Mixed
	Digit
	Binary
)

const (
	// Shift Operations Codes
	PS_AnyToPunc    = 0
	US_LowerToUpper = 28
	US_DigitToUpper = 15

	// Latch Operation Codes
	LL_UpperToLow   = 28
	ML_UpperToMix   = 29
	DL_UpperToDigit = 30

	ML_LowerToMix   = 29
	DL_LowerToDigit = 30

	UL_MixToUpper = 29
	LL_MixToLower = 28
	PL_MixToPunc  = 30

	UL_PuncToUpper  = 31
	UL_DigitToUpper = 14
)

type Encoder struct {
	bits     bytes.Buffer
	mode     EncodingMode
	prevMode EncodingMode
}

func NewEncoder() *Encoder {
	return &Encoder{
		bits:     *bytes.NewBuffer(make([]byte, 0)),
		mode:     Upper,
		prevMode: Upper,
	}
}

func (e *Encoder) encode(data []byte) ([]byte, error) {
	return nil, nil
}
