package main

import (
	"bytes"
	"errors"
	"fmt"
)

type EncodingMode byte

// Mode conts
const (
	Upper EncodingMode = iota
	Lower
	Mixed
	Punct
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
	bits bytes.Buffer
}

func NewEncoder() *Encoder {
	return &Encoder{
		bits: *bytes.NewBuffer(make([]byte, 0)),
	}
}

func (e *Encoder) Encode(text string) error {
	segments := SegmentText([]byte(text))
	optimalChg, _, err := findOptimalSequence(segments)
	if err != nil {
		fmt.Println(err)
	}
	// First segment can be encoded without checking for optimal Sequence
	segments[0].Encode(&e.bits)
	for i := 1; i < len(segments); i++ {
		e.encodeChange(segments[i-1], segments[i], optimalChg[i-1])
		segments[i].Encode(&e.bits)
	}
	return nil
}
func (e *Encoder) encodeChange(from, to Segment, isLatch bool) {
	if isLatch {
		switch from.mode {
		case 0:
			switch to.mode {
			case 0:
				break
			case 1:
				e.bits.WriteByte(LL_UpperToLow)
			case 2:
				e.bits.WriteByte(ML_UpperToMix)
			case 3:
				e.bits.WriteByte(ML_UpperToMix)
				e.bits.WriteByte(PL_MixToPunc)
			case 4:
				e.bits.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case 1:
			switch to.mode {
			case 0:
				e.bits.WriteByte(US_LowerToUpper)
			case 1:
				break
			case 2:
				e.bits.WriteByte(ML_LowerToMix)
			case 3:
				e.bits.WriteByte(ML_LowerToMix)
				e.bits.WriteByte(PL_MixToPunc)
			case 4:
				e.bits.WriteByte(DL_LowerToDigit)
			default:
				break
			}
		case 2:
			switch to.mode {
			case 0:
				e.bits.WriteByte(UL_MixToUpper)
			case 1:
				e.bits.WriteByte(LL_MixToLower)
			case 2:
				break
			case 3:
				e.bits.WriteByte(PL_MixToPunc)
			case 4:
				e.bits.WriteByte(UL_MixToUpper)
				e.bits.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case 3:
			switch to.mode {
			case 0:
				e.bits.WriteByte(UL_PuncToUpper)
			case 1:
				e.bits.WriteByte(UL_PuncToUpper)
				e.bits.WriteByte(LL_UpperToLow)
			case 2:
				e.bits.WriteByte(UL_PuncToUpper)
				e.bits.WriteByte(ML_UpperToMix)
			case 3:
				break
			case 4:
				e.bits.WriteByte(UL_PuncToUpper)
				e.bits.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case 4:
			switch to.mode {
			case 0:
				e.bits.WriteByte(UL_DigitToUpper)
			case 1:
				e.bits.WriteByte(DL_LowerToDigit)
			case 2:
				e.bits.WriteByte(UL_DigitToUpper)
				e.bits.WriteByte(ML_UpperToMix)
			case 3:
				e.bits.WriteByte(UL_DigitToUpper)
				e.bits.WriteByte(ML_UpperToMix)
				e.bits.WriteByte(PL_MixToPunc)
			default:
				break
			}
		default:
			break
		}
		// SHIFTS
	} else {
		switch from.mode {
		case 0:
			switch to.mode {
			case 3:
				e.bits.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case 1:
			switch to.mode {
			case 0:
				e.bits.WriteByte(US_LowerToUpper)
			case 3:
				e.bits.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case 2:
			switch to.mode {
			case 3:
				e.bits.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case 4:
			switch to.mode {
			case 0:
				e.bits.WriteByte(US_DigitToUpper)
			case 3:
				e.bits.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		default:
			break
		}
	}
}

func findOptimalSequence(segments []Segment) ([]bool, int, error) {
	n := len(segments)
	if n == 0 {
		return nil, 0, errors.New("no segments provided")
	}

	var changeSeq []bool
	totalCost := 0

	// For each consecutive pair, determine whether to Shift or Latch.
	for i := 0; i < n-1; i++ {
		from := segments[i].mode
		to := segments[i+1].mode

		tc := changeLen[from][to]
		// We decide: if latch cost is less than or equal to shift cost,
		// then we use latch; otherwise, we use shift.
		if len(segments[i+1].text) > 1 {
			changeSeq = append(changeSeq, true) // Latch
			totalCost += tc.Latch
		} else if tc.Shift < tc.Latch {
			changeSeq = append(changeSeq, false) // Shift
			totalCost += tc.Shift
		} else if tc.Shift == tc.Latch && len(segments[i+1].text) <= 1 {
			changeSeq = append(changeSeq, false) // Shift
			totalCost += tc.Shift
		} else {
			changeSeq = append(changeSeq, true) // Latch
			totalCost += tc.Latch
		}

	}

	return changeSeq, totalCost, nil
}
