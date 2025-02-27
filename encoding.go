package main

import (
	"bytes"
	"math"
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

// func (e *Encoder) Encode(text string) error {
// 	segments := SegmentText([]byte(text))
// 	optimalChg, _,  := findOptimalSequence(segments)

//		// First segment can be encoded without checking for optimal Sequence
//		segments[0].Encode(&e.bits)
//		for i := 1; i < len(segments); i++ {
//			e.encodeChange(segments[i-1], segments[i], optimalChg[i-1])
//			segments[i].Encode(&e.bits)
//		}
//		return nil
//	}

func (e *Encoder) Encode(text string) (int, error) {
	segments := SegmentText([]byte(text))
	totalCost, changes := findOptimalSequence(segments)

	currentSegment := segments[0]
	currentSegment.Encode(&e.bits)
	for i := 0; i < len(changes); i++ {
		e.encodeChange(changes[i].From, changes[i].To, changes[i].Mode)
		currentSegment = segments[i+1]
		currentSegment.Encode(&e.bits)
	}
	return totalCost, nil
}

func (e *Encoder) encodeChange(from, to EncodingMode, isLatch bool) {
	if isLatch {
		switch from {
		case EncodingMode(Upper):
			switch to {
			case EncodingMode(Upper):
				break
			case EncodingMode(Lower):
				e.bits.WriteByte(LL_UpperToLow)
			case EncodingMode(Mixed):
				e.bits.WriteByte(ML_UpperToMix)
			case EncodingMode(Punct):
				e.bits.WriteByte(ML_UpperToMix)
				e.bits.WriteByte(PL_MixToPunc)
			case EncodingMode(Digit):
				e.bits.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case EncodingMode(Lower):
			switch to {
			case EncodingMode(Upper):
				e.bits.WriteByte(US_LowerToUpper)
			case EncodingMode(Lower):
				break
			case EncodingMode(Mixed):
				e.bits.WriteByte(ML_LowerToMix)
			case EncodingMode(Punct):
				e.bits.WriteByte(ML_LowerToMix)
				e.bits.WriteByte(PL_MixToPunc)
			case EncodingMode(Digit):
				e.bits.WriteByte(DL_LowerToDigit)
			default:
				break
			}
		case EncodingMode(Mixed):
			switch to {
			case EncodingMode(Upper):
				e.bits.WriteByte(UL_MixToUpper)
			case EncodingMode(Lower):
				e.bits.WriteByte(LL_MixToLower)
			case EncodingMode(Mixed):
				break
			case EncodingMode(Punct):
				e.bits.WriteByte(PL_MixToPunc)
			case EncodingMode(Digit):
				e.bits.WriteByte(UL_MixToUpper)
				e.bits.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case EncodingMode(Punct):
			switch to {
			case EncodingMode(Upper):
				e.bits.WriteByte(UL_PuncToUpper)
			case EncodingMode(Lower):
				e.bits.WriteByte(UL_PuncToUpper)
				e.bits.WriteByte(LL_UpperToLow)
			case EncodingMode(Mixed):
				e.bits.WriteByte(UL_PuncToUpper)
				e.bits.WriteByte(ML_UpperToMix)
			case EncodingMode(Punct):
				break
			case EncodingMode(Digit):
				e.bits.WriteByte(UL_PuncToUpper)
				e.bits.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case EncodingMode(Digit):
			switch to {
			case EncodingMode(Upper):
				e.bits.WriteByte(UL_DigitToUpper)
			case EncodingMode(Lower):
				e.bits.WriteByte(DL_LowerToDigit)
			case EncodingMode(Mixed):
				e.bits.WriteByte(UL_DigitToUpper)
				e.bits.WriteByte(ML_UpperToMix)
			case EncodingMode(Punct):
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
		switch from {
		case EncodingMode(Upper):
			switch to {
			case EncodingMode(Punct):
				e.bits.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case EncodingMode(Lower):
			switch to {
			case EncodingMode(Upper):
				e.bits.WriteByte(US_LowerToUpper)
			case EncodingMode(Punct):
				e.bits.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case EncodingMode(Mixed):
			switch to {
			case EncodingMode(Upper):
				e.bits.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case EncodingMode(Digit):
			switch to {
			case EncodingMode(Upper):
				e.bits.WriteByte(US_DigitToUpper)
			case EncodingMode(Punct):
				e.bits.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		default:
			break
		}
	}
}

type Change struct {
	From EncodingMode
	To   EncodingMode
	Mode bool // either latch or shift
}

type dpState struct {
	cost          int
	prevMode      EncodingMode
	changeIsLatch bool
}

func findOptimalSequence(segments []Segment) (int, []Change) {
	if len(segments) == 0 {
		return 0, nil
	}

	shiftCost, latchCost := precomputeMinimalCosts()

	dp := make([]map[EncodingMode]dpState, len(segments))
	for i := range dp {
		dp[i] = make(map[EncodingMode]dpState)
	}

	firstSeg := segments[0]
	firstCost := segCost(firstSeg)
	dp[0][firstSeg.mode] = dpState{
		cost:          firstCost,
		prevMode:      EncodingMode(Upper),
		changeIsLatch: false,
	}

	for i := 1; i < len(segments); i++ {
		currentSeg := segments[i]
		currentMode := currentSeg.mode
		length := len(currentSeg.text)
		encodingCost := length * charSize[currentMode]

		for prevMode, prevState := range dp[i-1] {
			if length == 1 {
				sc := shiftCost[prevMode][currentMode]
				if sc != E {
					totalCost := prevState.cost + sc + encodingCost
					state, exists := dp[i][currentMode]
					if !exists || totalCost < state.cost {
						dp[i][currentMode] = dpState{
							cost:          totalCost,
							prevMode:      prevMode,
							changeIsLatch: false,
						}
					}
				}

				lc := latchCost[prevMode][currentMode]
				if lc != E {
					totalCost := prevState.cost + lc + encodingCost
					state, exists := dp[i][currentMode]
					if !exists || totalCost < state.cost {
						dp[i][currentMode] = dpState{
							cost:          totalCost,
							prevMode:      prevMode,
							changeIsLatch: true,
						}
					}
				}
			} else {
				lc := latchCost[prevMode][currentMode]
				if lc != E {
					totalCost := prevState.cost + lc + encodingCost
					state, exists := dp[i][currentMode]
					if !exists || totalCost < state.cost {
						dp[i][currentMode] = dpState{
							cost:          totalCost,
							prevMode:      prevMode,
							changeIsLatch: true,
						}
					}
				}
			}
		}
	}

	minCost := math.MaxInt32
	var bestMode EncodingMode
	for mode, state := range dp[len(segments)-1] {
		if state.cost < minCost {
			minCost = state.cost
			bestMode = mode
		}
	}

	if minCost == math.MaxInt32 {
		return 0, nil
	}

	changes := []Change{}
	currentMode := bestMode
	for i := len(segments) - 1; i > 0; i-- {
		state := dp[i][currentMode]
		changes = append([]Change{{
			From: state.prevMode,
			To:   currentMode,
			Mode: state.changeIsLatch,
		}}, changes...)
		currentMode = state.prevMode
	}

	return minCost, changes
}

func segLen(seg Segment) int {
	return len(seg.text)
}
func segCost(seg Segment) int {
	return segLen(seg) * charSize[seg.mode]
}
