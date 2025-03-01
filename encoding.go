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

func (e *Encoder) WriteText(text string) {
	for _, c := range text {
		e.bits.WriteByte(byte(c))
	}
}

func (e *Encoder) WriteBit(bit bool) error {
	byteIdx := e.bits.Len() - 1
	bitPos := 7 - e.bits.Len()%8
	if e.bits.Len()%8 == 0 {
		e.bits.WriteByte(0)
	}
	data := e.bits.Bytes()
	if bit {
		data[byteIdx] |= 1 << bitPos
	} else {
		data[byteIdx] &= ^(1 << bitPos)
	}
	return nil
}

	

func (e *Encoder) Encode(text string) (int, error) {
	segments := SegmentText([]byte(text))
	totalCost, changes := findOptimalSequence(segments)

	// Add a change to Upper if necessary at beggining
	firstSeg := segments[0]
	newChange := Change{}
	if firstSeg.mode != EncodingMode(Upper) {
		if segLen(firstSeg) <= 1 && (firstSeg.mode == EncodingMode(Lower) || firstSeg.mode == EncodingMode(Digit)) {
			newChange = Change{From: EncodingMode(Upper), To: firstSeg.mode, Mode: false}
		} else {
			newChange = Change{From: EncodingMode(Upper), To: firstSeg.mode, Mode: true}
		}
	}

	if newChange != (Change{}) {
		e.encodeChange(newChange.From, newChange.To, newChange.Mode)
	}
	for i := 0; i < len(changes); i++ {
		segments[i].Encode(&e.bits)
		e.encodeChange(changes[i].From, changes[i].To, changes[i].Mode)
	}
	segments[len(changes)].Encode(&e.bits)
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

func findOptimalSequence(segments []Segment) (int, []Change) {
	if len(segments) == 0 {
		return 0, nil
	}
	shiftCost, latchCost := precomputeMinimalCosts()
	dp := make([]map[EncodingMode]dpState, len(segments))
	for i := range dp {
		dp[i] = make(map[EncodingMode]dpState)
	}

	// Initialize the first segment
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
			// Try latch operation
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

			// Try shift operation for single-character segments
			if length == 1 {
				sc := shiftCost[prevMode][currentMode]
				if sc != E {
					totalCost := prevState.cost + sc + encodingCost

					// For a shift, the effective mode after encoding remains the same as before
					// This is critical - we're storing this state in prevMode, not currentMode
					state, exists := dp[i][prevMode]
					if !exists || totalCost < state.cost {
						dp[i][prevMode] = dpState{
							cost:          totalCost,
							prevMode:      prevMode,
							changeIsLatch: false,
							// Record that we used a shift for reconstruction
							usedShift:   true,
							shiftToMode: currentMode,
						}
					}
				}
			}
		}
	}

	// Find minimum cost ending mode
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

	// Reconstruct the changes
	changes := []Change{}
	currentMode := bestMode

	for i := len(segments) - 1; i > 0; i-- {
		state := dp[i][currentMode]

		// If we used a shift at this position
		if state.usedShift {
			changes = append([]Change{{
				From: state.prevMode,
				To:   state.shiftToMode,
				Mode: false, // false for shift
			}}, changes...)
		} else if state.changeIsLatch {
			// If we used a latch at this position
			changes = append([]Change{{
				From: state.prevMode,
				To:   currentMode,
				Mode: true, // true for latch
			}}, changes...)
		}

		currentMode = state.prevMode
	}

	return minCost, changes
}

type dpState struct {
	cost          int
	prevMode      EncodingMode
	changeIsLatch bool
	usedShift     bool
	shiftToMode   EncodingMode
}

func segLen(seg Segment) int {
	return len(seg.text)
}
func segCost(seg Segment) int {
	return segLen(seg) * charSize[seg.mode]
}
