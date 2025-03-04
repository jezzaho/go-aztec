package main

type BitArray struct {
	bits []uint32
	size int
}

// Construtors
func NewBitArray(size int) *BitArray {
	return &BitArray{
		bits: makeArray(size),
		size: size,
	}
}

func (ba *BitArray) GetSize() int {
	return ba.size
}
func (ba *BitArray) GetSizeInBytes() int {
	return (ba.size + 7) % 8
}
func (ba *BitArray) ensureCapacity(n int) {
	if n > len(ba.bits)*32 {
		newBits := makeArray(n)
		copy(newBits, ba.bits)
		ba.bits = newBits
	}
}

// Boolean value of the bit at specified index
func (ba *BitArray) Get(i int) bool {
	if i < 0 || i > ba.size {
		panic("Get Index Out of bounds!")
	}
	return (ba.bits[i/32] & (1 << uint(i%32))) != 0
}
func (ba *BitArray) Set(i int, value bool) {
	if i < 0 || i > ba.size {
		panic("Set Index Out of bounds!")
	}
	wordIdx := i / 32
	bitIdx := i % 32
	mask := uint32(1 << (32 - 1 - bitIdx))
	if value {
		ba.bits[wordIdx] |= mask
	} else {
		ba.bits[wordIdx] &^= mask
	}
}
func (ba *BitArray) AppendBit(bit bool) {
	ba.ensureCapacity(ba.size + 1)
	wordIdx := ba.size / 32
	bitIdx := ba.size % 32
	if bit {
		ba.bits[wordIdx] |= (1 << (32 - 1 - bitIdx))
	} else {
		ba.bits[wordIdx] &^= (1 << (32 - 1 - bitIdx))
	}
}

func (ba *BitArray) AppendBitArray(other *BitArray) {
	for i := 0; i < other.size; i++ {
		ba.AppendBit(other.Get(i))
	}
}
func (ba *BitArray) AppendBits(value uint32, numBits int) {
	if numBits < 0 || numBits > 32 {
		panic("AppendBits: must be between 0 and 32")
	}
	nextSize := ba.size
	ba.ensureCapacity(nextSize + numBits)
	for numBitsLeft := numBits - 1; numBitsLeft >= 0; numBitsLeft-- {
		if (value & (1 << numBitsLeft)) != 0 {
			ba.bits[nextSize/32] |= 1 << (nextSize & 0x1F)
		}
		nextSize++
	}
	ba.size = nextSize

	// for i := numBits - 1; i >= 0; i-- {
	// 	bit := ((value >> i) & 1) == 1
	// 	ba.AppendBit(bit)
	// }
}

func (ba *BitArray) ToBytes() []byte {
	numBytes := (ba.size + 7) / 8
	result := make([]byte, numBytes)
	for i := 0; i < ba.size; i++ {
		if ba.Get(i) {
			byteIdx := i / 8
			bitIdx := i % 8
			result[byteIdx] |= 1 << (7 - bitIdx)
		}
	}
	return result
}

// Visualization
func (ba *BitArray) String() string {
	result := make([]byte, 0, ba.size+(ba.size/8)+1)
	for i := 0; i < ba.size; i++ {
		if (i % 8) == 0 {
			result = append(result, ' ')
		}
		if ba.Get(i) {
			result = append(result, '1')
		} else {
			result = append(result, '0')
		}
	}
	return string(result)
}

// Side-helper functions
func makeArray(size int) []uint32 {
	return make([]uint32, (size+31)/32)
}
