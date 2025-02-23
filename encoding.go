package main

import (
	"bytes"
	"errors"
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

	segments := SegmentText(data)
	if len(segments) == 0 {
		return nil, errors.New("Data is empty")
	}
	currentSeg := segments[0]
	for _, s := range segments {
		if s.mode != currentSeg.mode {

		}
		s.Encode(&e.bits)

	}
	return nil, nil
}

func findOptimalSequence(segments []Segment) ([]int, error) {
	// Initialize distances array
	dist := make([][]int, len(segments))
	for i := range dist {
		dist[i] = make([]int, len(segments))
	}

	// Initialize previous modes to reconstruct the path
	prev := make([][]int, len(segments))
	for i := range prev {
		prev[i] = make([]int, len(segments))
	}

	// Populate the distance and prev matrices
	for i := 0; i < len(segments)-1; i++ {
		for j := i + 1; j < len(segments); j++ {
			from := segments[i].mode
			to := segments[j].mode
			// Calculate the cost for each segment transition
			shiftCost := changeLen[from][to].Shift
			latchCost := changeLen[from][to].Latch

			// Use the shift or latch cost to populate the dist and prev matrices
			if shiftCost < latchCost {
				dist[i][j] = shiftCost
				prev[i][j] = i
			} else {
				dist[i][j] = latchCost
				prev[i][j] = j
			}
		}
	}

	optimalSeq := []int{0} // Start from the first segment
	current := 0
	for current < len(segments)-1 {
		minCost := E // Start with the highest possible cost
		nextSegment := -1

		// Find the next segment with the minimal cost transition
		for i := current + 1; i < len(segments); i++ {
			cost := dist[current][i]
			if cost < minCost {
				minCost = cost
				nextSegment = i
			}
		}

		if nextSegment == -1 {
			break // No valid path found, break the loop
		}

		// Add the next segment to the optimal sequence
		optimalSeq = append(optimalSeq, nextSegment)
		current = nextSegment
	}

	return optimalSeq, nil

}
