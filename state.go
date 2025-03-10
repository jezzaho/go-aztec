package main

type State struct {
	mode                 int
	token                Token
	binaryShiftByteCount int
	bitCount             int
	binaryShiftCost      int
}

var INITIAL_STATE = State{token: EmptyToken, mode: MODE_UPPER, binaryShiftByteCount: 0, bitCount: 0, binaryShiftCost: 0}
