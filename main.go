package main

import (
	"errors"
	"fmt"
)

type BitBuffer struct {
	bits []bool
	mode []bool
}

func NewBitBuffer() *BitBuffer {
	return &BitBuffer{
		bits: make([]bool, 0),
		mode: make([]bool, 0),
	}
}

func (b *BitBuffer) SimpleBinary() string {
	var textBuffer []byte

	for _, v := range b.bits {
		if v {
			textBuffer = append(textBuffer, '1')
		} else {
			textBuffer = append(textBuffer, '0')
		}
	}

	return string(textBuffer)
}

func (b *BitBuffer) BitsToBinary() string {
	if len(b.bits) == 0 {
		return "EMPTY BUFFER"
	}

	var textBuffer []byte
	bitIndex := 0
	charCounter := 0
	charLength := 0

	// Start the bit processing loop
	for bitIndex < len(b.bits) {
		// Determine chunk size based on the mode
		chunkSize := 5
		if !b.mode[charCounter] { // If mixed mode (non-alphabetic), use 8 bits
			chunkSize = 8
		}

		// Add space after every chunk of bits (except at the beginning)
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

		// Once we've processed the full chunk for the current character, move to the next character and mode
		if charLength == chunkSize {
			charCounter++
			charLength = 0
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

func main() {

	text := "HELLO WORLD!"
	b := EncodeTextToBitsWithMode(text)
	s := b.BitsToBinary()
	fmt.Print(s)
}
