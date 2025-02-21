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
func gfMultiply(a, b uint8) uint8 {
	if a == 0 || b == 0 {
		return 0
	}
	return expTable[(logTable[a]+logTable[b])%(fieldSize-1)]
}

// Generator polynomial
func createGeneratorPolynomial(t int) []uint8 {
	g := []uint8{1}

	for i := 0; i < t; i++ {
		newG := make([]uint8, len(g)+1)

		for j := 0; j < len(g); j++ {
			newG[j] ^= gfMultiply(g[j], expTable[i])
		}

		copy(newG[1:], g)

		g = newG
	}

	return g
}

func calculateParitySymbols(message []uint8, generator []uint8) []uint8 {
	numParity := len(generator) - 1
	// Initialize remainder with the correct size
	remainder := make([]uint8, len(message)+numParity)

	// Copy message into remainder
	copy(remainder, message)

	// Perform polynomial division
	for i := 0; i < len(message); i++ {
		if remainder[i] == 0 {
			continue // Skip if coefficient is zero
		}

		coef := remainder[i]

		// XOR with generator polynomial
		for j := 0; j < len(generator); j++ {
			remainder[i+j] ^= gfMultiply(coef, generator[j])
		}
	}

	// Return only the parity symbols
	return remainder[len(message):]
}

func encodeWithParity(message []uint8, numParity int) []uint8 {
	generator := createGeneratorPolynomial(numParity)
	parity := calculateParitySymbols(message, generator)

	return append(message, parity...)
}

// Helper to be deleted
func readEncodedWithParity(message []uint8) {
	messageBuffer := ""
	for _, b := range message {
		messageBuffer += string(b) + " "
	}

	fmt.Println(messageBuffer)
}

// Constructing Grid Structure

const compactLayerSize = 4

func CreateGrid(layers int, compact bool) [][]bool {
	var size int
	if compact {
		size = 11 + (layers-1)*4
	} else {
		size = 15 + (layers-1)*4
	}

	grid := make([][]bool, size)
	for i := range grid {
		grid[i] = make([]bool, size)
	}
	AddFinderPattern(grid, layers, compact)

	return grid
}

func AddFinderPattern(grid [][]bool, layers int, compact bool) {
	size := len(grid)
	center := size / 2
	numRings := layers + 1

	for i := 0; i < numRings; i++ {
		color := (i % 2) == 0
		for x := center - i; x <= center+i; x++ {
			grid[x][center-i] = color
			grid[x][center+i] = color
		}
		for y := center - i; y <= center+1; y++ {
			grid[center-i][y] = color
			grid[center+i][y] = color
		}
	}
	grid[center][center] = true
}

func PrintGrid(grid [][]bool) {
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func main() {

	grid := CreateGrid(4, false)
	PrintGrid(grid)

	// initGaloisField()
	// g := createGeneratorPolynomial(5)
	// fmt.Println("Generator Polynomial:", g)

	// message := []uint8{72, 69, 76, 76, 79}
	// numParity := 4
	// encodedMessage := encodeWithParity(message, numParity)
	// fmt.Println("Encoded message: ", encodedMessage)
	// readEncodedWithParity(encodedMessage)

	// text := "HELLO WORLD!"
	// b := EncodeTextToBitsWithMode(text)
	// b.CalculateSize()``
	// b.EncodeSize()
	// s := b.BitsToBinary()
	// fmt.Println(s)
	// b.ApplyBitPadding()
	// s = b.SimpleBinary(6)
	// fmt.Println(s)
}
