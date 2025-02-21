package main

import "fmt"

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
