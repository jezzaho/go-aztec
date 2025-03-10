package main

// Modes definitions
const (
	MODE_UPPER = iota
	MODE_LOWER
	MODE_PUNCT
	MODE_DIGIT
	MODE_MIXED
	MODE_BINARY
)

type State struct {
	mode int 
	bitCount int 
	bitBuffer BitArray
	binaryShiftByteCount int 
}


