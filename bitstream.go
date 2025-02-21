package main

import (
	"errors"
	"math/bits"
)

type BitBuffer struct {
	bits         []bool
	mode         []bool
	size         uint32
	bitsRequired uint32
}

func NewBitBuffer() *BitBuffer {
	return &BitBuffer{
		bits:         make([]bool, 0),
		mode:         make([]bool, 0),
		size:         0,
		bitsRequired: 0,
	}
}

func (b *BitBuffer) CalculateSize() {
	size := uint32(0)
	for _, v := range b.mode {
		if v {
			size += 5
		} else {
			size += 8
		}
	}

	b.size = size
}

func (b *BitBuffer) PrintBufferWithChunkSize(chunkSize int) string {
	var textBuffer []byte

	bitIndex := 0
	charCounter := 0
	charLength := 0

	for bitIndex < int(b.size) {
		if bitIndex > 0 && charLength == 0 {
			textBuffer = append(textBuffer, ' ')
		}

		if b.bits[bitIndex] {
			textBuffer = append(textBuffer, '1')
		} else {
			textBuffer = append(textBuffer, '0')
		}

		bitIndex++
		charLength++

		if charLength == chunkSize {
			charCounter++
			charLength = 0
		}
	}

	return string(textBuffer)
}

func (b *BitBuffer) PrintBuffer() string {
	if len(b.bits) == 0 {
		return "EMPTY BUFFER"
	}

	var textBuffer []byte
	bitIndex := 0
	charCounter := 0
	charLength := 0

	// Process text encoding (excluding last bits for size)
	textEnd := len(b.bits) - int(b.bitsRequired) // Ignore size bits for now

	for bitIndex < textEnd {
		// Determine chunk size based on the mode
		chunkSize := 5
		if !b.mode[charCounter] { // If mixed mode (non-alphabetic), use 8 bits
			chunkSize = 8
		}

		// Add space before each new chunk (except at the start)
		if bitIndex > 0 && charLength == 0 {
			textBuffer = append(textBuffer, ' ')
		}

		// Append the current bit (1 or 0)
		if b.bits[bitIndex] {
			textBuffer = append(textBuffer, '1')
		} else {
			textBuffer = append(textBuffer, '0')
		}

		// Move to the next bit
		bitIndex++
		charLength++

		// Once we've processed the full chunk for the current character, move to the next character
		if charLength == chunkSize {
			charCounter++
			charLength = 0
		}
	}

	// Append space before adding size encoding
	textBuffer = append(textBuffer, ' ')

	// Add the size encoding bits at the end
	for i := textEnd; i < len(b.bits); i++ {
		if b.bits[i] {
			textBuffer = append(textBuffer, '1')
		} else {
			textBuffer = append(textBuffer, '0')
		}
	}

	return string(textBuffer)
}

func (b *BitBuffer) AppendBits(value, size int) error {
	if size > 32 {
		return errors.New("can't append more than 32 bits at time")
	}
	// pad bits if necessary
	for i := size - 1; i >= 0; i-- {
		if (value>>i)&1 == 1 {
			b.bits = append(b.bits, true)

		} else {
			b.bits = append(b.bits, false)
		}

	}

	return nil
}
func EncodeTextToBits(text string) *BitBuffer {
	bitBuffer := NewBitBuffer()

	for _, char := range text {
		asciiValue := int(char)
		bitBuffer.AppendBits(asciiValue, 8)

	}
	return bitBuffer
}
func EncodeTextToBitsWithMode(text string) *BitBuffer {
	bitBuffer := NewBitBuffer()

	for _, char := range text {
		asciiValue := int(char)
		if (asciiValue >= 65 && asciiValue <= 90) || (asciiValue >= 97 && asciiValue <= 122) {
			bitBuffer.AppendBits(asciiValue, 5)
			bitBuffer.mode = append(bitBuffer.mode, true)
		} else {
			bitBuffer.AppendBits(asciiValue, 8)
			bitBuffer.mode = append(bitBuffer.mode, false)
		}
	}
	return bitBuffer
}
func EncodeIntWithMode(data []uint8) *BitBuffer {
	bitBuffer := NewBitBuffer()

	for _, char := range data {
		if (char >= 65 && char <= 90) || (char >= 97 && char <= 122) {
			bitBuffer.AppendBits(int(char), 5)
			bitBuffer.mode = append(bitBuffer.mode, true)
		} else {
			bitBuffer.AppendBits(int(char), 8)
			bitBuffer.mode = append(bitBuffer.mode, false)
		}
	}
	bitBuffer.CalculateSize()
	return bitBuffer
}

func (b *BitBuffer) EncodeSize() error {
	if b.size <= 0 {
		return errors.New("size is not correct or not calculated")
	}

	bitsRequired := uint32(bits.Len32(b.size))
	b.bitsRequired = bitsRequired

	return b.AppendBits(int(b.size), int(bitsRequired))

}

func (b *BitBuffer) ApplyBitPadding() {
	padding := (6 - (b.size % 6)) % 6
	for i := 0; i < int(padding); i++ {
		b.bits = append(b.bits, false)
	}
	b.size += uint32(padding)
}
func GenerateAztecBitstream(data string, numParity int) (*BitBuffer, error) {
	bitBuffer := EncodeTextToBitsWithMode(data)

	err := bitBuffer.EncodeSize()
	if err != nil {
		return nil, err
	}

	bitBuffer.ApplyBitPadding()
	uint8message := boolsToUint8(bitBuffer.bits)
	gf := NewGaloisField()

	rse := RSEncoder{
		numParity:   numParity,
		galoisField: gf,
	}
	parity := rse.EncodeWithParity(uint8message)
	parityBits := uint8ToBools(parity)

	bitBuffer.bits = append(bitBuffer.bits, parityBits...)

	// TODO: Finish
	return bitBuffer, nil
}
