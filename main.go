package main

import (
	"errors"
	"fmt"
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

func (b *BitBuffer) SimpleBinary(chunkSize int) string {
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

func (b *BitBuffer) BitsToBinary() string {
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

// ERROR CORRECTNES SR

const fieldSize = 256
const primitivePolynomial = 0x11D // x^8 + x^4 + x^3 + x^2 + 1

var expTable [fieldSize]uint8
var logTable [fieldSize]uint8

func initGaloisField() {
	var value uint16 = 1

	// Initialize log table with invalid values
	for i := range logTable {
		logTable[i] = 255 // Mark as undefined initially
	}

	for i := 0; i < fieldSize-1; i++ {
		expTable[i] = uint8(value)
		logTable[uint8(value)] = uint8(i)

		value <<= 1
		if value&0x100 != 0 {
			value ^= uint16(primitivePolynomial)
		}
	}
	expTable[fieldSize-1] = expTable[0]
}

func gfDivide(a, b uint8) uint8 {
	if b == 0 {
		panic("Division by zero in GF(2⁸)")
	}
	if a == 0 {
		return 0
	}

	logA := int(logTable[a])
	logB := int(logTable[b])

	if logA == 255 || logB == 255 {
		panic("Invalid input value")
	}

	logResult := (logA - logB + 255) % 255
	return expTable[logResult]
}

func gfAdd(a, b int) int {
	return a ^ b
}
func gfMultiply(a, b int) uint8 {
	if a == 0 || b == 0 {
		return 0
	}
	return expTable[(logTable[a]+logTable[b])%(fieldSize-1)]
}

func main() {

	initGaloisField()

	fmt.Println("Addition (5 ⊕ 7):", gfAdd(5, 7))            // Expected: 2
	fmt.Println("Multiplication (5 ⊗ 7):", gfMultiply(5, 7)) // Expected: 35
	fmt.Println("Division (35 ÷ 5):", gfDivide(35, 5))       // Expected: 7

	// text := "HELLO WORLD!"
	// b := EncodeTextToBitsWithMode(text)
	// b.CalculateSize()
	// b.EncodeSize()
	// s := b.BitsToBinary()
	// fmt.Println(s)
	// b.ApplyBitPadding()
	// s = b.SimpleBinary(6)
	// fmt.Println(s)
}
