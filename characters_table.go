package main

var codeChars = map[string][]string{
	"upper": {
		"P/S", " ", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "L/L", "M/L", "D/L", "B/S",
	},
	"lower": {
		"P/S", " ", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "U/S", "M/L", "D/L", "B/S",
	},
	"mixed": {
		"P/S", " ", "\x01", "\x02", "\x03", "\x04", "\x05", "\x06", "\x07", "\x08", "\t", "\n", "\x0b", "\x0c", "\r", "\x1b", "\x1c", "\x1d", "\x1e", "\x1f", "@", "\\", "^", "_", "`", "|", "~", "\x7f", "L/L", "U/L", "P/L", "B/S",
	},
	"punct": {
		"FLG(n)", "\r", "\r\n", ". ", ", ", ": ", "!", "\"", "H", "$", "%", "&", "'", "(", ")", "*", "+", ",", "-", ".", "/", ":", ";", "<", "=", ">", "?", "[", "]", "{", "}", "U/L",
	},
	"digit": {
		"P/S", " ", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", ",", ".", "U/L", "U/S",
	},
}

var (
	upperChars = codeChars["upper"][1 : len(codeChars["upper"])-4]
	lowerChars = codeChars["lower"][1 : len(codeChars["lower"])-4]
	mixedChars = codeChars["mixed"][1 : len(codeChars["mixed"])-4]
	punctChars = codeChars["punct"][1 : len(codeChars["punct"])-4]
	digitChars = codeChars["digit"][1 : len(codeChars["digit"])-4]
)

// TODO: FUNC FOR OPTIMAL SHIFT LATCH
// TODO: TABLES OF COST OF SHIFTING AND LATCHING
