package main

import (
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
	Input  *ByteStream
	Output *ByteStream
}

func NewEncoder() *Encoder {
	return &Encoder{
		Output: NewByteStream(),
	}
}
func (e *Encoder) addData(data *ByteStream) {
	e.Input = data
}

func (e *Encoder) Encode(stream *ByteStream) (int, error) {
	e.addData(stream)
	segments := SegmentText([]byte(e.Input.Bytes.Bytes()))
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
		segments[i].Encode(&e.Output.Bytes)
		e.encodeChange(changes[i].From, changes[i].To, changes[i].Mode)
	}
	segments[len(changes)].Encode(&e.Output.Bytes)
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
				e.Output.Bytes.WriteByte(LL_UpperToLow)
			case EncodingMode(Mixed):
				e.Output.Bytes.WriteByte(ML_UpperToMix)
			case EncodingMode(Punct):
				e.Output.Bytes.WriteByte(ML_UpperToMix)
				e.Output.Bytes.WriteByte(PL_MixToPunc)
			case EncodingMode(Digit):
				e.Output.Bytes.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case EncodingMode(Lower):
			switch to {
			case EncodingMode(Upper):
				e.Output.Bytes.WriteByte(US_LowerToUpper)
			case EncodingMode(Lower):
				break
			case EncodingMode(Mixed):
				e.Output.Bytes.WriteByte(ML_LowerToMix)
			case EncodingMode(Punct):
				e.Output.Bytes.WriteByte(ML_LowerToMix)
				e.Output.Bytes.WriteByte(PL_MixToPunc)
			case EncodingMode(Digit):
				e.Output.Bytes.WriteByte(DL_LowerToDigit)
			default:
				break
			}
		case EncodingMode(Mixed):
			switch to {
			case EncodingMode(Upper):
				e.Output.Bytes.WriteByte(UL_MixToUpper)
			case EncodingMode(Lower):
				e.Output.Bytes.WriteByte(LL_MixToLower)
			case EncodingMode(Mixed):
				break
			case EncodingMode(Punct):
				e.Output.Bytes.WriteByte(PL_MixToPunc)
			case EncodingMode(Digit):
				e.Output.Bytes.WriteByte(UL_MixToUpper)
				e.Output.Bytes.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case EncodingMode(Punct):
			switch to {
			case EncodingMode(Upper):
				e.Output.Bytes.WriteByte(UL_PuncToUpper)
			case EncodingMode(Lower):
				e.Output.Bytes.WriteByte(UL_PuncToUpper)
				e.Output.Bytes.WriteByte(LL_UpperToLow)
			case EncodingMode(Mixed):
				e.Output.Bytes.WriteByte(UL_PuncToUpper)
				e.Output.Bytes.WriteByte(ML_UpperToMix)
			case EncodingMode(Punct):
				break
			case EncodingMode(Digit):
				e.Output.Bytes.WriteByte(UL_PuncToUpper)
				e.Output.Bytes.WriteByte(DL_UpperToDigit)
			default:
				break
			}
		case EncodingMode(Digit):
			switch to {
			case EncodingMode(Upper):
				e.Output.Bytes.WriteByte(UL_DigitToUpper)
			case EncodingMode(Lower):
				e.Output.Bytes.WriteByte(DL_LowerToDigit)
			case EncodingMode(Mixed):
				e.Output.Bytes.WriteByte(UL_DigitToUpper)
				e.Output.Bytes.WriteByte(ML_UpperToMix)
			case EncodingMode(Punct):
				e.Output.Bytes.WriteByte(UL_DigitToUpper)
				e.Output.Bytes.WriteByte(ML_UpperToMix)
				e.Output.Bytes.WriteByte(PL_MixToPunc)
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
				e.Output.Bytes.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case EncodingMode(Lower):
			switch to {
			case EncodingMode(Upper):
				e.Output.Bytes.WriteByte(US_LowerToUpper)
			case EncodingMode(Punct):
				e.Output.Bytes.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case EncodingMode(Mixed):
			switch to {
			case EncodingMode(Upper):
				e.Output.Bytes.WriteByte(PS_AnyToPunc)
			default:
				break
			}
		case EncodingMode(Digit):
			switch to {
			case EncodingMode(Upper):
				e.Output.Bytes.WriteByte(US_DigitToUpper)
			case EncodingMode(Punct):
				e.Output.Bytes.WriteByte(PS_AnyToPunc)
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

// Stuff bits
type BitWriter struct {
	buffer      ByteStream
	currentByte byte
	bitCount    int
}

func (bw *BitWriter) WriteBit(bit int) {
	if bit == 1 {
		bw.currentByte |= (1 << (7 - bw.bitCount))
	}
	bw.bitCount++
	if bw.bitCount == 8 {
		bw.buffer.WriteByte(bw.currentByte)
		bw.currentByte = 0
		bw.bitCount = 0
	}
}
func (bw *BitWriter) WriteBits(value uint, length int) {
	for i := length; i >= 0; i-- {
		bw.WriteBit(int((value >> uint(i)) & 1))
	}
}
func (bw *BitWriter) Flush() {
	if bw.bitCount > 0 {
		bw.buffer.WriteByte(bw.currentByte)
	}
}

func (bw *BitWriter) Bytes() []byte {
	return bw.buffer.Bytes.Bytes()
}

func (e *Encoder) StuffBits(wordSize int) {
	bw := BitWriter{}
	n := e.Output.Bytes.Len() * 8
	mask := (1 << wordSize) - 2

	for i := 0; i < n; i += wordSize {
		word := 0
		for j := 0; j < wordSize; j++ {
			bitIndex := (i + j) % 8
			byteIndex := (i + j) / 8
			if byteIndex < e.Output.Bytes.Len() && ((e.Output.Bytes.Bytes()[byteIndex]>>(7-bitIndex))&1) == 1 {
				word |= 1 << (wordSize - 1 - j)
			}
		}

		if (word & mask) == mask {
			bw.WriteBits(uint(word&mask), wordSize)
			i--
		} else if (word & mask) == 0 {
			bw.WriteBits(uint(word|1), wordSize)
			i--
		} else {
			bw.WriteBits(uint(word), wordSize)
		}
	}

	bw.Flush()
	e.Output = NewByteStream()
	e.Output.WriteBytes(bw.Bytes())

}
