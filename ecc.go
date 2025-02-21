package main

import (
	"fmt"
)

const fieldSize = 256
const primitivePolynomial = 0x11D // x^8 + x^4 + x^3 + x^2 + 1

type GaloisField struct {
	expTable [fieldSize]uint8
	logTable [fieldSize]uint8
}

func NewGaloisField() *GaloisField {
	gf := &GaloisField{}
	gf.init()

	return gf
}

func (gf *GaloisField) init() {
	var value uint16 = 1

	// Initialize log table with invalid values
	for i := range gf.logTable {
		gf.logTable[i] = 255 // Mark as undefined initially
	}

	for i := 0; i < fieldSize-1; i++ {
		gf.expTable[i] = uint8(value)
		gf.logTable[uint8(value)] = uint8(i)

		value <<= 1
		if value&0x100 != 0 {
			value ^= uint16(primitivePolynomial)
		}
	}
	gf.expTable[fieldSize-1] = gf.expTable[0]
}

func (gf *GaloisField) divide(a, b uint8) uint8 {
	if b == 0 {
		panic("Division by zero in GF(2⁸)")
	}
	if a == 0 {
		return 0
	}

	logA := int(gf.logTable[a])
	logB := int(gf.logTable[b])

	if logA == 255 || logB == 255 {
		panic("Invalid input value")
	}

	logResult := (logA - logB + 255) % 255
	return gf.expTable[logResult]
}

func (gf *GaloisField) add(a, b uint8) uint8 {
	return a ^ b
}
func (gf *GaloisField) multiply(a, b uint8) uint8 {
	if a == 0 || b == 0 {
		return 0
	}
	return gf.expTable[(gf.logTable[a]+gf.logTable[b])%(fieldSize-1)]
}

type RSEncoder struct {
	galoisField *GaloisField
	numParity   int
	generator   []uint8
}

func NewRSEncoder(galoisField *GaloisField, numParity int) *RSEncoder {
	encoder := &RSEncoder{
		galoisField: galoisField,
		numParity:   numParity,
	}
	encoder.generator = createGeneratorPolynomial(numParity, galoisField)
	return encoder
}

// Generator polynomial
func createGeneratorPolynomial(t int, gf *GaloisField) []uint8 {
	g := []uint8{1}

	for i := 0; i < t; i++ {
		newG := make([]uint8, len(g)+1)

		for j := 0; j < len(g); j++ {
			newG[j] ^= gf.multiply(g[j], gf.expTable[i])
		}

		copy(newG[1:], g)

		g = newG
	}

	return g
}

func calculateParitySymbols(message []uint8, generator []uint8, gf *GaloisField) []uint8 {
	numParity := len(generator) - 1

	if len(message) == 0 {
		panic("Message cannot be empty")
	}
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
			remainder[i+j] ^= gf.multiply(coef, generator[j])
		}
	}

	// Return only the parity symbols
	return remainder[len(message):]
}
func (rse *RSEncoder) EncodeWithParity(message []uint8) []uint8 {
	return encodeWithParity(message, rse.numParity, rse.generator, rse.galoisField)
}

func encodeWithParity(message []uint8, numParity int, generator []uint8, gf *GaloisField) []uint8 {
	parity := calculateParitySymbols(message, generator, gf)

	return append(message, parity...)
}

func boolsToUint8(bits []bool) []uint8 {
	var uint8Arr []uint8
	var currentByte uint8

	for i, b := range bits {
		if b {
			currentByte |= (1 << (7 - i%8))
		}
		if (i+1)%8 == 0 || i == len(bits)-1 {
			uint8Arr = append(uint8Arr, currentByte)
			currentByte = 0
		}
	}
	return uint8Arr
}

// Helper to convert []uint8 to []bool
func uint8ToBools(uint8Arr []uint8) []bool {
	var bits []bool
	for _, byteValue := range uint8Arr {
		for i := 7; i >= 0; i-- {
			bits = append(bits, (byteValue&(1<<i)) != 0)
		}
	}
	return bits
}

// Helper to be deleted
func readEncodedWithParity(message []uint8) {
	messageBuffer := ""
	for _, b := range message {
		messageBuffer += string(b) + " "
	}

	fmt.Println(messageBuffer)
}
