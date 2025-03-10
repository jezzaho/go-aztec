package main

type Token interface {
	GetPrevious() Token
	Add(int, int) Token
	AddBinaryShift(int, int) Token
	AppendTo(*BitArray, []byte)
}

var EmptyToken Token = &SimpleToken{}

type SimpleToken struct {
	previous Token
	value    int
	bitCount int
}

func (s *SimpleToken) GetPrevious() Token {
	return s.previous
}
func (s *SimpleToken) Add(value, bitCount int) Token {
	return &SimpleToken{
		previous: s,
		value:    value,
		bitCount: bitCount,
	}
}
func (s *SimpleToken) AddBinaryShift(start, byteCount int) Token {
	return &BinaryShiftToken{
		previous:  s,
		start:     start,
		byteCount: byteCount,
	}
}
func (s *SimpleToken) AppendTo(bitArray *BitArray, text []byte) {
	if bitArray != nil {
		bitArray.AppendBits(uint32(s.value), s.bitCount)
	}
	if s.previous != nil {
		s.previous.AppendTo(bitArray, text)
	}
}

type BinaryShiftToken struct {
	previous  Token
	start     int
	byteCount int
}

func (b *BinaryShiftToken) GetPrevious() Token {
	return b.previous
}

func (b *BinaryShiftToken) Add(value, bitCount int) Token {
	return &SimpleToken{
		previous: b,
		value:    value,
		bitCount: bitCount,
	}
}

func (b *BinaryShiftToken) AddBinaryShift(start, byteCount int) Token {
	return &BinaryShiftToken{
		previous:  b,
		start:     start,
		byteCount: byteCount,
	}
}

func (b *BinaryShiftToken) AppendTo(bitArray *BitArray, text []byte) {
	// Implement binary shift logic here
	bitCount := b.byteCount * 8
	if b.byteCount <= 31 {
		bitCount += 10
	} else if b.byteCount <= 62 {
		bitCount += 20
	} else {
		bitCount += 21
	}

	// Add binary shift header
	bitArray.AppendBits(uint32(b.start), bitCount)

	// Add actual bytes
	for i := 0; i < b.byteCount; i++ {
		bitArray.AppendBits(uint32(int(text[b.start+i])), 8)
	}

	if b.previous != nil {
		b.previous.AppendTo(bitArray, text)
	}
}
