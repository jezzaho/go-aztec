package main

import (
	"bytes"
	"fmt"
)

type CC struct {
	char string
	code int
}

var codeChars = map[EncodingMode][]CC{
	EncodingMode(Upper): {
		{"P/S", 0}, {" ", 1}, {"A", 2}, {"B", 3}, {"C", 4}, {"D", 5}, {"E", 6}, {"F", 7}, {"G", 8}, {"H", 9}, {"I", 10}, {"J", 11}, {"K", 12}, {"L", 13}, {"M", 14}, {"N", 15}, {"O", 16}, {"P", 17}, {"Q", 18}, {"R", 19}, {"S", 20}, {"T", 21}, {"U", 22}, {"V", 23}, {"W", 24}, {"X", 25}, {"Y", 26}, {"Z", 27}, {"L/L", 28}, {"M/L", 29}, {"D/L", 30}, {"B/S", 31},
	},
	EncodingMode(Lower): {
		{"P/S", 0}, {" ", 1}, {"a", 2}, {"b", 3}, {"c", 4}, {"d", 5}, {"e", 6}, {"f", 7}, {"g", 8}, {"h", 9}, {"i", 10}, {"j", 11}, {"k", 12}, {"l", 13}, {"m", 14}, {"n", 15}, {"o", 16}, {"p", 17}, {"q", 18}, {"r", 19}, {"s", 20}, {"t", 21}, {"u", 22}, {"v", 23}, {"w", 24}, {"x", 25}, {"y", 26}, {"z", 27}, {"U/S", 28}, {"M/L", 29}, {"D/L", 30}, {"B/S", 31},
	},
	EncodingMode(Mixed): {
		{"P/S", 0}, {" ", 1}, {"\x01", 2}, {"\x02", 3}, {"\x03", 4}, {"\x04", 5}, {"\x05", 6}, {"\x06", 7}, {"\x07", 8}, {"\x08", 9}, {"\x09", 10}, {"\x0a", 11}, {"\x0b", 12}, {"\x0c", 13}, {"\x0d", 14}, {"\x1b", 15}, {"\x1c", 16}, {"\x1d", 17}, {"\x1e", 18}, {"\x1f", 19}, {"@", 20}, {"\\", 21}, {"^", 22}, {"_", 23}, {"`", 24}, {"|", 25}, {"~", 26}, {"\x7f", 27}, {"L/L", 28}, {"U/L", 29}, {"P/L", 30}, {"B/S", 31},
	},
	EncodingMode(Punct): {
		{"FLG(n)", 0}, {"\r", 1}, {"\r\n", 2}, {". ", 3}, {", ", 4}, {": ", 5}, {"!", 6}, {"\"", 7}, {"\x23", 8}, {"$", 9}, {"%", 10}, {"&", 11}, {"'", 12}, {"(", 13}, {")", 14}, {"*", 15}, {"+", 16}, {",", 17}, {"-", 18}, {".", 19}, {"/", 20}, {":", 21}, {";", 22}, {"<", 23}, {"=", 24}, {">", 25}, {"?}", 26}, {"[", 27}, {"]", 28}, {"{", 29}, {"}", 30}, {"U/L", 31},
	},
	EncodingMode(Digit): {
		{"P/S", 0}, {" ", 1}, {"0", 2}, {"1", 3}, {"2", 4}, {"3", 5}, {"4", 6}, {"5", 7}, {"6", 8}, {"7", 9}, {"8", 10}, {"9", 11}, {",", 12}, {".", 13}, {"U/L", 14}, {"U/S", 15},
	},
}

var (
	upperChars = codeChars[EncodingMode(Upper)][1 : len(codeChars[EncodingMode(Upper)])-4]
	lowerChars = codeChars[EncodingMode(Lower)][1 : len(codeChars[EncodingMode(Lower)])-4]
	mixedChars = codeChars[EncodingMode(Mixed)][1 : len(codeChars[EncodingMode(Mixed)])-4]
	punctChars = codeChars[EncodingMode(Punct)][1 : len(codeChars[EncodingMode(Punct)])-2]
	digitChars = codeChars[EncodingMode(Digit)][1 : len(codeChars[EncodingMode(Digit)])-2]
)

func getCodeForChar(mode EncodingMode, char rune) (int, bool) {

	if chars, ok := codeChars[mode]; ok {
		for _, cc := range chars {
			if cc.char == string(char) {
				return cc.code, true
			}
		}
	}
	return -1, false
}

// TODO: FUNC FOR OPTIMAL SHIFT LATCH
// TODO: TABLES OF COST OF SHIFTING AND LATCHING

type TravelCost struct {
	Shift int
	Latch int
}

const E = 99999

var (
	changeLen = map[EncodingMode]map[EncodingMode]TravelCost{
		// Upper
		EncodingMode(Upper): {
			// upper
			// lower
			// mixed
			// punct
			// digit
			// binary
			EncodingMode(Upper): {E, 0},
			EncodingMode(Lower): {E, 5},
			EncodingMode(Mixed): {E, 5},
			// Intermediate Latch
			EncodingMode(Punct):  {5, 10},
			EncodingMode(Digit):  {E, 5},
			EncodingMode(Binary): {E, 10},
		},
		// Lower
		EncodingMode(Lower): {
			// Intermediate Latch
			EncodingMode(Upper): {5, 10},
			EncodingMode(Lower): {E, 0},
			EncodingMode(Mixed): {E, 5},
			// Intermediate Latch
			EncodingMode(Punct):  {5, 10},
			EncodingMode(Digit):  {E, 5},
			EncodingMode(Binary): {E, 10},
		},
		//Mixed
		EncodingMode(Mixed): {
			EncodingMode(Upper): {E, 5},
			EncodingMode(Lower): {E, 5},
			EncodingMode(Mixed): {E, 0},
			EncodingMode(Punct): {5, 5},
			// Intermediate Latch
			EncodingMode(Digit):  {E, 10},
			EncodingMode(Binary): {E, 10},
		},
		// Punct
		EncodingMode(Punct): {
			EncodingMode(Upper): {E, 5},
			// Intermediate Latch
			EncodingMode(Lower): {E, 10},
			// Intermediate Latch
			EncodingMode(Mixed): {E, 10},
			EncodingMode(Punct): {E, 0},
			// Intermediate Latch
			EncodingMode(Digit):  {E, 10},
			EncodingMode(Binary): {E, 15},
		},
		// Digit
		EncodingMode(Digit): {
			EncodingMode(Upper):  {4, 4},
			EncodingMode(Lower):  {E, 9},
			EncodingMode(Mixed):  {E, 9},
			EncodingMode(Punct):  {4, 14},
			EncodingMode(Digit):  {E, 0},
			EncodingMode(Binary): {E, 14},
		},
		// Binary
		EncodingMode(Binary): {
			EncodingMode(Upper):  {E, 0},
			EncodingMode(Lower):  {E, 0},
			EncodingMode(Mixed):  {E, 0},
			EncodingMode(Punct):  {E, 0},
			EncodingMode(Digit):  {E, 0},
			EncodingMode(Binary): {E, 0},
		},
	}
)

var (
	charSize = map[EncodingMode]int{
		EncodingMode(Upper):  5,
		EncodingMode(Lower):  5,
		EncodingMode(Mixed):  5,
		EncodingMode(Punct):  5,
		EncodingMode(Digit):  4,
		EncodingMode(Binary): 8,
	}
)

func belongsTo(char byte, group []CC) bool {
	stringChar := string(char)
	for _, v := range group {
		if stringChar == v.char {
			return true
		}
	}
	return false
}

// func isLower(char byte) bool {
// 	return char >= 'a' && char <= 'z'
// }

// // TODO: Check if this can be compacted or changed
// // Mixed mode contains specific symbols, NOT digits!
// func isMixed(char byte) bool {
// 	return char == ' ' || char == '\r' || char == '\t' ||
// 		char == '"' || char == '#' || char == '$' || char == '%' || char == '&' ||
// 		char == '\'' || char == '*' || char == '+' || char == ',' || char == '-' ||
// 		char == '.' || char == '/' || char == ':' || char == ';' || char == '<' ||
// 		char == '=' || char == '>' || char == '?' || char == '@' || char == '[' ||
// 		char == '\\' || char == ']' || char == '_' || char == '`'
// }
// func isDigit(char byte) bool {
// 	return (char >= '0' && char <= '9') || char == ',' || char == '.'
// }
// func isPunct(char byte) bool {
// 	return char == ';' || char == '<' || char == '>' || char == '@' ||
// 		char == '[' || char == '\\' || char == ']' || char == '_' ||
// 		char == '`' || char == '~' || char == '!' || char == '\n' ||
// 		char == '\r' || char == '\t' || char == '"' || char == '|' ||
// 		char == '{' || char == '}' || char == '\''
// }

func getMode(char byte) EncodingMode {
	switch {
	case belongsTo(char, upperChars):
		return Upper
	case belongsTo(char, lowerChars):
		return Lower
	case belongsTo(char, digitChars):
		return Digit
	case belongsTo(char, mixedChars):
		return Mixed
	case belongsTo(char, punctChars):
		return Punct
	default:
		return Binary
	}
}

type Segment struct {
	mode EncodingMode
	text string
}

func SegmentText(text []byte) []Segment {
	if len(text) == 0 {
		return nil
	}

	var segments []Segment
	start := 0
	currentMode := getMode(text[0])

	for i := 1; i < len(text); i++ {
		charMode := getMode(text[i])
		if charMode != currentMode {
			segments = append(segments, Segment{mode: currentMode, text: string(text[start:i])})
			start = i
			currentMode = charMode
		}
	}

	segments = append(segments, Segment{mode: currentMode, text: string(text[start:])})
	return segments
}

func (s *Segment) Encode(buf *bytes.Buffer) error {

	for _, val := range s.text {
		if code, ok := getCodeForChar(s.mode, val); !ok {
			return fmt.Errorf("Could not find a code for char: %v", val)
		} else {
			buf.WriteByte(byte(code))
		}

	}
	return nil
}
