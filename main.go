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

func (b *BitBuffer) BitsToBinary() string {
	if len(b.bits) == 0 {
		return "EMPTY BUFFER"
	}
	var textBuffer []byte
	for i, v := range b.bits {
		if b.mode[i] && i%5 == 0 {
			textBuffer = append(textBuffer, ' ')
		} else if !b.mode[i] && i%8 == 0 {
			textBuffer = append(textBuffer, ' ')
		}

		if v {
			textBuffer = append(textBuffer, '1')
		} else {
			textBuffer = append(textBuffer, '0')
		}
	}
	return string(textBuffer)
}

func (b *BitBuffer) AppendBits(value, size int, alphabetic bool) error {
	if size > 32 {
		return errors.New("can't append more than 32 bits at time")
	}
	// Add validation for value size
	if value >= (1 << size) {
		return fmt.Errorf("value %d cannot be encoded in %d bits", value, size)
	}
	// pad bits if necessary
	for i := size - 1; i >= 0; i-- {
		if (value>>i)&1 == 1 {
			b.bits = append(b.bits, true)

		} else {
			b.bits = append(b.bits, false)
		}
		b.mode = append(b.mode, alphabetic)
	}

	return nil
}
func EncodeTextToBits(text string) *BitBuffer {
	bitBuffer := NewBitBuffer()

	for _, char := range text {
		asciiValue := int(char)
		bitBuffer.AppendBits(asciiValue, 8, false)

	}
	return bitBuffer
}
func EncodeTextToBitsWithMode(text string) *BitBuffer {
	bitBuffer := NewBitBuffer()

	for _, char := range text {
		asciiValue := int(char)
		if (asciiValue >= 65 && asciiValue <= 90) || (asciiValue >= 97 && asciiValue <= 122) {
			bitBuffer.AppendBits(asciiValue, 5, true)
		} else {
			bitBuffer.AppendBits(asciiValue, 8, false)
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

const (
	Alphabetic = iota
	Byte
)
