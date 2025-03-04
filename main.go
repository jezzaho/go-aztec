package main

import "fmt"

func main() {
	ba := NewBitArray(0)
	ba.AppendBits(75, 8)
	ba.AppendBits(65, 8)
	bt := ba.ToBytes()
	fmt.Printf("ba: %v\n", string(bt))
}
