package main

type State struct {
	mode                 int
	token                Token
	binaryShiftByteCount int
	bitCount             int
	binaryShiftCost      int
}

var INITIAL_STATE = State{token: EmptyToken, mode: MODE_UPPER, binaryShiftByteCount: 0, bitCount: 0, binaryShiftCost: 0}

func NewState(token Token, mode int, binaryBytes int, bitCount int) *State {
	return &State{
		token:                token,
		mode:                 mode,
		binaryShiftByteCount: binaryBytes,
		bitCount:             bitCount,
		binaryShiftCost:      calculateBSC(binaryBytes),
	}
}

func calculateBSC(binaryShiftByteCount int) int {
	if binaryShiftByteCount > 62 {
		return 21
	}
	if binaryShiftByteCount > 31 {
		return 20
	}
	if binaryShiftByteCount > 0 {
		return 10
	}
	return 0
}

func (s *State) GetMode() int {
	return s.mode
}
func (s *State) GetToken() *Token {
	return &s.token
}
func (s *State) GetBinaryShiftByteCount() int {
	return s.binaryShiftByteCount
}
func (s *State) GetBitCount() int {
	return s.bitCount
}
// -=-=-=-=
func(s *State) shiftAndAppend(mode, value int) *State {
	token := s.token
	var thisModeBitCount int 
	if s.mode == MODE_DIGIT {
		thisModeBitCount = 4
	} else {
		thisModeBitCount = 5 
	}
	token = token.Add(TheEncoder.SHIFT_TABLE[s.mode][mode], thisModeBitCount)
	token = token.Add(value, 5)
	return NewState(token, s.mode, 0, s.bitCount + thisModeBitCount + 5)
}
func(s *State) latchAndAppend(mode, value int) *State {
	bitCount := s.bitCount
	token := s.token

	if mode != s.mode {
		latch := TheEncoder.LATCH_TABLE[s.mode][mode]
		token := token.Add(latch & 0xFFFF, latch >> 16)
		bitCount += latch >> 16
	}
	var latchModeBitCount int 
	if mode == MODE_DIGIT {
		latchModeBitCount = 4 
	} else {
		latchModeBitCount = 5
	}
	token = token.Add(value, latchModeBitCount)
	return NewState(token, mode, 0, bitCount + latchModeBitCount)
}


func appendFLGn() *State {
	result := 
}
