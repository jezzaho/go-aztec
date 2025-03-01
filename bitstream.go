package main

// import "errors"

// type BitArray struct {
// 	bits []int
// 	size int
// }

// func makeArray(size int) []int {
// 	return make([]int, (size+31)/32)
// }

// func NewBitArray(size int) *BitArray {
// 	return &BitArray{bits: makeArray(size), size: size}
// }

// func (b *BitArray) GetSize() int {
// 	return b.size
// }
// func (b *BitArray) GetSizeInBytes() int {
// 	return (b.size + 7) / 8
// }

// func (b *BitArray) grow(newSize int) {
// 	if newSize > len(b.bits)*32 {
// 		newBits := makeArray(newSize)
// 		copy(newBits, b.bits)
// 		b.bits = newBits
// 	}
// }

// func (b *BitArray) GetBit(i int) bool {
// 	return (b.bits[i/32] & (1 << (i & 0x1F))) != 0
// }
// func (b *BitArray) SetBit(i int) {
// 	b.bits[i/32] |= 1 << (i & 0x1F)
// }
// func (b *BitArray) FlipBit(i int) {
// 	b.bits[i/32] ^= 1 << (i & 0x1F)
// }

// func (b *BitArray) AppendBit(bit bool) {
// 	b.grow(b.size + 1)
// 	if bit {
// 		b.bits[b.size/32] |= 1 << (b.size & 0x1F)
// 	}
// 	b.size++
// }
// func (b *BitArray) AppendBits(value, numBits int) error {
// 	if numBits < 0 || numBits > 32 {
// 		return errors.New("Num bits must be between 0 and 32")
// 	}
// 	nextSize := b.size
// 	b.grow(nextSize + numBits)
// 	for numBitsLeft := numBits - 1; numBitsLeft >= 0; numBitsLeft-- {
// 		if (value & (1 << numBitsLeft)) != 0 {
// 			b.bits[nextSize/32] |= 1 << (nextSize & 0x1F)
// 		}
// 		nextSize++
// 	}
// 	b.size = nextSize
// 	return nil
// }

// func (b *BitArray) ToBytes(bitOffset int, array []byte, offset, numBytes int) {
// 	for i := 0; i < numBytes; i++ {
// 		value := 0
// 		for j := 0; j < 8; j++ {
// 			if b.GetBit(bitOffset) {
// 				value |= 1 << (7 - j)
// 			}
// 			bitOffset++
// 		}
// 		array[offset+i] = byte(value)
// 	}
// }

// func (b *BitArray) ToBitArray(text []byte) {
// 	for _, c := range text {
// 		for i := 7; i >= 0; i-- {
// 			if (c & (1 << i)) != 0 {
// 				b.AppendBit(true)
// 			} else {
// 				b.AppendBit(false)
// 			}
// 		}
// 	}
// }
