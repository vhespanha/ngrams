package ngrams

import "fmt"

const (
	LowercaseSpec         = "a-z"
	LowercaseAlphanumSpec = "0-9a-z"
	AlphanumericSpec      = "0-9A-Za-z"
	PrintableSpec         = " -/0-9A-Za-z{-~"
)

type symbol int8

type Alphabet struct {
	size   int
	decode [256]symbol
	encode []byte
}

func newAlphabet(spec string) *Alphabet {
	specChars := expandRange(spec)

	a := &Alphabet{
		size:   len(specChars),
		encode: make([]byte, len(specChars)),
	}

	var dec [256]int8
	for i := range dec {
		dec[i] = -1
	}

	for i, r := range specChars {
		b := byte(r)
		a.encode[i] = b
		dec[b] = int8(i)
	}

	return a
}

func NewLowercaseAlphabet() *Alphabet         { return newAlphabet(LowercaseSpec) }
func NewLowercaseAlphanumAlphabet() *Alphabet { return newAlphabet(LowercaseAlphanumSpec) }
func NewAlphanumericAlphabet() *Alphabet      { return newAlphabet(AlphanumericSpec) }
func NewPrintableAlphabet() *Alphabet         { return newAlphabet(PrintableSpec) }

func (a *Alphabet) SymbolFromByte(b byte) (symbol, bool) {
	idx := a.decode[b]
	if idx < 0 {
		return 0, false
	}
	return symbol(idx), true
}

func (a *Alphabet) SymbolsFromString(s string) ([]symbol, error) {
	n := len(s)
	symbols := make([]symbol, n)
	for i := range n {
		symbol, ok := a.SymbolFromByte(s[i])
		if !ok {
			return nil, fmt.Errorf("invalid symbol %q", s[i])
		}
		symbols = append(symbols, symbol)
	}
	return symbols, nil
}

func expandRange(spec string) []rune {
	var out []rune
	for i := 0; i < len(spec); i++ {
		if i+2 < len(spec) && spec[i+1] == '-' {
			for r := rune(spec[i]); r <= rune(spec[i+2]); r++ {
				out = append(out, r)
			}
			i += 2
		} else {
			out = append(out, rune(spec[i]))
		}
	}
	return out
}
