package main

// Modes definitions
const (
	MODE_UPPER = iota
	MODE_LOWER
	MODE_PUNCT
	MODE_DIGIT
	MODE_MIXED
	MODE_BINARY
)

var LATCH_TABLE = [][]int{
	{
		0,
		(5 << 16) + 28,              // UPPER -> LOWER
		(5 << 16) + 30,              // UPPER -> DIGIT
		(5 << 16) + 29,              // UPPER -> MIXED
		(10 << 16) + (29 << 5) + 30, // UPPER -> MIXED -> PUNCT
	},
	{
		(9 << 16) + (30 << 4) + 14, // LOWER -> DIGIT -> UPPER
		0,
		(5 << 16) + 30,              // LOWER -> DIGIT
		(5 << 16) + 29,              // LOWER -> MIXED
		(10 << 16) + (29 << 5) + 30, // LOWER -> MIXED -> PUNCT
	},
	{
		(4 << 16) + 14,             // DIGIT -> UPPER
		(9 << 16) + (14 << 5) + 28, // DIGIT -> UPPER -> LOWER
		0,
		(9 << 16) + (14 << 5) + 29, // DIGIT -> UPPER -> MIXED
		(14 << 16) + (14 << 10) + (29 << 5) + 30,
		// DIGIT -> UPPER -> MIXED -> PUNCT
	},
	{
		(5 << 16) + 29,              // MIXED -> UPPER
		(5 << 16) + 28,              // MIXED -> LOWER
		(10 << 16) + (29 << 5) + 30, // MIXED -> UPPER -> DIGIT
		0,
		(5 << 16) + 30, // MIXED -> PUNCT
	},
	{
		(5 << 16) + 31,              // PUNCT -> UPPER
		(10 << 16) + (31 << 5) + 28, // PUNCT -> UPPER -> LOWER
		(10 << 16) + (31 << 5) + 30, // PUNCT -> UPPER -> DIGIT
		(10 << 16) + (31 << 5) + 29, // PUNCT -> UPPER -> MIXED
		0,
	},
}
var SHIFT_TABLE = [6][6]int{}

func initST() [6][6]int {
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			SHIFT_TABLE[i][j] = -1
		}
	}
	SHIFT_TABLE[MODE_UPPER][MODE_PUNCT] = 0

	SHIFT_TABLE[MODE_LOWER][MODE_PUNCT] = 0
	SHIFT_TABLE[MODE_LOWER][MODE_UPPER] = 28

	SHIFT_TABLE[MODE_MIXED][MODE_PUNCT] = 0

	SHIFT_TABLE[MODE_DIGIT][MODE_PUNCT] = 0
	SHIFT_TABLE[MODE_DIGIT][MODE_UPPER] = 15

	return SHIFT_TABLE
}

var CHAR_MAP = [5][256]int{}

func initCM() [5][256]int {
	CHAR_MAP[MODE_UPPER][' '] = 1
	for c := 'A'; c <= 'Z'; c++ {
		CHAR_MAP[MODE_UPPER][c] = int(c - 'A' + 2)
	}
	CHAR_MAP[MODE_LOWER][' '] = 1
	for c := 'a'; c <= 'z'; c++ {
		CHAR_MAP[MODE_LOWER][c] = int(c - 'a' + 2)
	}
	CHAR_MAP[MODE_DIGIT][' '] = 1
	for c := '0'; c <= '9'; c++ {
		CHAR_MAP[MODE_DIGIT][c] = int(c - '0' + 2)
	}
	CHAR_MAP[MODE_DIGIT][','] = 12
	CHAR_MAP[MODE_DIGIT]['.'] = 13

	mixedTable := []rune{
		0, ' ', 1, 2, 3, 4, 5, 6, 7,
		'\b', '\t', '\n', 11, '\f', '\r',
		27, 28, 29, 30, 31,
		'@', '\\', '^', '_', '`', '|', '~', 127,
	}
	for i, char := range mixedTable {
		CHAR_MAP[MODE_MIXED][char] = i
	}
	punctTable := []rune{
		0, '\r', 0, 0, 0, 0,
		'!', '\'', '#', '$', '%', '&', '\'',
		'(', ')', '*', '+', ',', '-', '.',
		'/', ':', ';', '<', '=', '>', '?',
		'[', ']', '{', '}',
	}
	for i, char := range punctTable {
		CHAR_MAP[MODE_PUNCT][char] = i
	}

	return CHAR_MAP
}

type Encoder struct {
	LATCH_TABLE [][]int
	SHIFT_TABLE [6][6]int
	CHAR_MAP    [5][256]int
}

func NewEncoder() *Encoder {
	return &Encoder{
		LATCH_TABLE: LATCH_TABLE,
		SHIFT_TABLE: initST(),
		CHAR_MAP:    initCM(),
	}
}
