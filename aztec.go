 package main

// import "errors"

// const (
// 	DEFAULT_EC_PERCENT   = 23
// 	DEFAULT_AZTEC_LAYERS = 0
// 	MAX_NB_BITS          = 32
// 	MAX_NB_BITS_COMPACT  = 4
// )

// var WORD_SIZE = []int{4, 6, 6, 8, 8, 8, 8, 8, 8, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
// 	12, 12, 12, 12, 12, 12, 12, 12, 12, 12}

// type AztecCode struct {
// 	compact   bool
// 	size      int
// 	layers    int
// 	codeWords int
// 	// add BitMatrix later
// }

// func NewAztecCode(data string, minECCPercent int) (*AztecCode, error) {
// 	encoder := NewEncoder()
// 	encoder.Encode(data)

// 	bits := encoder.bits.Bytes()

// 	eccBits := len(bits)*DEFAULT_EC_PERCENT/100 + 11
// 	totalSizeBits := len(bits) + eccBits
// 	var compact bool
// 	var layers int
// 	var totalBitsPerLayer int
// 	var wordSize int
	

// 	for i := 0; ; i++ {
// 		if i > MAX_NB_BITS {
// 			return nil, errors.New("Data too large for an Aztec Code")
// 		}
// 		compact = i <= 3
// 		if compact {
// 			layers = i + 1
// 		} else {
// 			layers = i
// 		}

// 		// Change name
// 		totalBitsPerLayer = totalBitsInLayer(layers, compact)
// 		if totalSizeBits > totalBitsPerLayer {
// 			continue
// 		}
// 		// Set word-size
// 		wordSize = WORD_SIZE[layers]
// 		stuffedBits = S(bits, wordSize)

// 	}

// 	return &AztecCode{}
// }

// func totalBitsInLayer(layers int, compact bool) int {
// 	if compact {
// 		return (88 + 16*layers) * layers
// 	} else {
// 		return (112 + 16*layers) * layers
// 	}
// }

// // import (
// // 	"fmt"
// // 	"image"
// // 	"image/color"
// // 	"image/png"
// // 	"os"
// // )

// // type AztecGrid struct {
// // 	Size     int
// // 	Capacity int
// // 	Grid     [][]bool
// // }

// // const compactLayerSize = 4

// // func CreateAztecGrid(layers int, compact bool, bitBuff BitBuffer) *AztecGrid {
// // 	var size int
// // 	if compact {
// // 		size = 11 + (layers-1)*4
// // 	} else {
// // 		size = 15 + (layers-1)*4
// // 	}

// // 	grid := make([][]bool, size)
// // 	for i := range grid {
// // 		grid[i] = make([]bool, size)
// // 	}
// // 	AddFinderPattern(grid, layers, compact)
// // 	AddModeIndicator(grid, layers)
// // 	AddErrorCorrectionLevel(grid, 3)
// // 	PlaceData(grid, bitBuff)

// // 	return &AztecGrid{
// // 		Size:     size,
// // 		Capacity: size * size,
// // 		Grid:     grid,
// // 	}
// // }

// // func AddFinderPattern(grid [][]bool, layers int, compact bool) {
// // 	size := len(grid)
// // 	center := size / 2
// // 	numRings := layers + 1

// // 	for i := 0; i < numRings; i++ {
// // 		color := (i % 2) == 0
// // 		for x := center - i; x <= center+i; x++ {
// // 			grid[x][center-i] = color
// // 			grid[x][center+i] = color
// // 		}
// // 		for y := center - i; y <= center+i; y++ {
// // 			grid[center-i][y] = color
// // 			grid[center+i][y] = color
// // 		}
// // 	}
// // 	grid[center][center] = true
// // }

// // func AddModeIndicator(grid [][]bool, layers int) {
// // 	size := len(grid)
// // 	grid[0][0] = true
// // 	grid[size-1][0] = true
// // 	grid[0][size-1] = true
// // }
// // func AddErrorCorrectionLevel(grid [][]bool, errorCorrectionLevel int) {
// // 	size := len(grid)
// // 	switch errorCorrectionLevel {
// // 	case 1:
// // 		grid[0][size-2] = true // przykładowe umiejscowienie wskaźnika
// // 	case 2:
// // 		grid[size-2][size-2] = true
// // 	case 3:
// // 		grid[size-2][0] = true
// // 	case 4:
// // 		grid[size-1][0] = true
// // 	}
// // }

// // func PrintGrid(grid [][]bool) {
// // 	for _, row := range grid {
// // 		for _, cell := range row {
// // 			if cell {
// // 				fmt.Print("█")
// // 			} else {
// // 				fmt.Print(" ")
// // 			}
// // 		}
// // 		fmt.Println()
// // 	}
// // }

// // // GenerateImage creates a PNG image from the grid with specified pixel size
// // func (az *AztecGrid) GenerateImage(pixelSize int, filename string) error {
// // 	// Calculate the image dimensions
// // 	height := len(az.Grid) * pixelSize
// // 	width := len(az.Grid[0]) * pixelSize

// // 	// Create a new white image
// // 	img := image.NewRGBA(image.Rect(0, 0, width, height))
// // 	white := color.RGBA{255, 255, 255, 255}
// // 	black := color.RGBA{0, 0, 0, 255}

// // 	// Fill the image with white first
// // 	for y := 0; y < height; y++ {
// // 		for x := 0; x < width; x++ {
// // 			img.Set(x, y, white)
// // 		}
// // 	}

// // 	// Draw black pixels for true values
// // 	for y, row := range az.Grid {
// // 		for x, cell := range row {
// // 			if cell {
// // 				// Fill the pixel block
// // 				for py := 0; py < pixelSize; py++ {
// // 					for px := 0; px < pixelSize; px++ {
// // 						img.Set(x*pixelSize+px, y*pixelSize+py, black)
// // 					}
// // 				}
// // 			}
// // 		}
// // 	}

// // 	// Create the output file
// // 	f, err := os.Create(filename)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	defer f.Close()

// // 	// Encode and save the image as PNG
// // 	return png.Encode(f, img)
// // }

// // func PlaceData(grid [][]bool, bitBuff BitBuffer) {
// // 	bitIndex := 0
// // 	size := len(grid)

// // 	for y := 0; y < size; y++ {
// // 		for x := 0; x < size; x++ {
// // 			if grid[y][x] != false {
// // 				continue
// // 			}
// // 			if bitIndex < int(bitBuff.size) {
// // 				grid[y][x] = bitBuff.bits[bitIndex]
// // 				bitIndex++
// // 			}
// // 		}
// // 	}
// // }
