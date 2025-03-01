package main

import "fmt"

// import (
// 	"log"
// )

// func main() {
// 	txt := "HELLO KARCIA KOCHAM CIE"

// 	buff, _ := GenerateAztecBitstream(txt, 3)
// 	grid := CreateAztecGrid(4, true, *buff)
// 	err := grid.GenerateImage(10, "aztec.png")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// initGaloisField()
// 	// g := createGeneratorPolynomial(5)
// 	// fmt.Println("Generator Polynomial:", g)

// 	// message := []uint8{72, 69, 76, 76, 79}
// 	// numParity := 4
// 	// encodedMessage := encodeWithParity(message, numParity)
// 	// fmt.Println("Encoded message: ", encodedMessage)
// 	// readEncodedWithParity(encodedMessage)

// 	// text := "HELLO WORLD!"
// 	// b := EncodeTextToBitsWithMode(text)
// 	// b.CalculateSize()``
// 	// b.EncodeSize()
// 	// s := b.BitsToBinary()
// 	// fmt.Println(s)
// 	// b.ApplyBitPadding()
// 	// s = b.SimpleBinary(6)
// 	// fmt.Println(s)
// }

func main() {
	text := "1"
	encoder := NewEncoder()
	encoder.WriteText(text)
	fmt.Printf("%v\n", encoder.bits.Bytes())
	encoder.WriteBit(true)
	fmt.Printf("%v\n", encoder.bits.Bytes())

	// bitarray := NewBitArray(1)
	// bitarray.ToBitArray([]byte(text))
	// var bytesArray []byte
	// bitarray.ToBytes(4, bytesArray, 0, 2)
	// fmt.Printf("%v\n", bitarray.bits)
	// symbolSize, layers, compact := determineAztecParameters(3159, 0.23)
	// fmt.Printf("Symbol Size: %v, Layers: %v, Compact: %v\n", symbolSize, layers, compact)
}
