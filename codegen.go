package ngrams

import (
	"bytes"
	"fmt"
	"math"
)

const maxWidth = 80

func tabWidth(s string) int {
	w := 0
	for _, r := range s {
		if r == '\t' {
			w += 8
		} else {
			w++
		}
	}
	return w
}

type wrapper struct {
	b            *bytes.Buffer
	lineWidth    int
	maxWidth     int
	indent       string
	indentWidth  int
	lastWasSpace bool
}

func (w *wrapper) write(s string) {
	sWidth := tabWidth(s)
	if w.lineWidth > w.indentWidth && w.lineWidth+sWidth > w.maxWidth {
		if w.lastWasSpace {
			w.b.Truncate(w.b.Len() - 1)
		}
		w.b.WriteString("\n")
		w.lineWidth = 0
		w.lastWasSpace = false
	}
	if w.lineWidth == 0 {
		w.b.WriteString(w.indent)
		w.lineWidth = w.indentWidth
		w.lastWasSpace = false
	}
	w.b.WriteString(s)
	w.lineWidth += sWidth
	w.lastWasSpace = (len(s) > 0 && s[len(s)-1] == ' ')
}

func (w *wrapper) flush() {
	if w.lineWidth > 0 {
		if w.lastWasSpace {
			w.b.Truncate(w.b.Len() - 1)
		}
		w.b.WriteString("\n")
		w.lineWidth = 0
		w.lastWasSpace = false
	}
}

func (t *table[T]) GenGo(pkg string) string {
	var b bytes.Buffer

	w := &wrapper{b: &b, maxWidth: maxWidth, indent: "\t", indentWidth: tabWidth("\t")}

	fmt.Fprintf(&b, "// Generated from %d-gram table, ", t.n)
	fmt.Fprintf(&b, "%d entries, alphabet: %s\n", len(t.freqs), t.alphabet.spec)

	fmt.Fprintf(&b, "package %s\n\n", pkg)
	fmt.Fprintf(&b, "const (\n\tn = %d\n", t.n)
	fmt.Fprintf(&b, "\talphabetSize = %d\n", t.alphabet.size)
	fmt.Fprintf(&b, "\ttotal = %d\n)\n\n", t.total)

	// decode table
	b.WriteString("var decode = [256]int8{\n")
	for i := range 256 {
		w.write(fmt.Sprintf("%d, ", t.alphabet.decode[i]))
	}
	w.flush()
	b.WriteString("}\n\n")

	// encode table
	b.WriteString("var encode = [alphabetSize]byte{\n")
	for i := range t.alphabet.size {
		w.write(fmt.Sprintf("%d, ", t.alphabet.encode[i]))
	}
	w.flush()
	b.WriteString("}\n\n")

	// freqs table
	var verb, primitive string
	switch any(*new(T)).(type) {
	case uint64:
		verb, primitive = "%d", "uint64"
	case prob:
		verb, primitive = "%.17g", "float64"
	case logprob:
		verb, primitive = "%.17g", "float64"
	default:
		verb, primitive = "%.17g", "float64"
	}

	fmt.Fprintf(&b, "var freqs = [%d]%s{\n", len(t.freqs), primitive)
	for i := range t.freqs {
		w.write(fmt.Sprintf(verb+", ", t.freqs[i]))
	}
	w.flush()
	b.WriteString("}\n\n")

	return b.String()
}

func (t *table[T]) GenC() string {
	var b bytes.Buffer

	w := &wrapper{b: &b, maxWidth: maxWidth, indent: "\t", indentWidth: tabWidth("\t")}

	fmt.Fprintf(&b, "// Generated from %d-gram table, ", t.n)
	fmt.Fprintf(&b, "%d entries, alphabet: %s\n", len(t.freqs), t.alphabet.spec)

	b.WriteString("#include <stdint.h>\n")
	b.WriteString("#include <math.h>\n\n")

	fmt.Fprintf(&b, "#define ALPHABET_SIZE %d\n\n", t.alphabet.size)

	// decode table
	b.WriteString("static const int8_t decode[256] = {\n")
	for i := range 256 {
		w.write(fmt.Sprintf("%d, ", t.alphabet.decode[i]))
	}
	w.flush()
	b.WriteString("};\n\n")

	// encode table
	b.WriteString("static const uint8_t encode[ALPHABET_SIZE] = {\n")
	for i := range t.alphabet.size {
		w.write(fmt.Sprintf("%d, ", t.alphabet.encode[i]))
	}
	w.flush()
	b.WriteString("};\n\n")

	var verb, primitive string

	switch any(*new(T)).(type) {
	case uint64:
		verb, primitive = "%d", "uint64_t"
	case prob:
		verb, primitive = "%.17g", "double"
	case logprob:
		verb, primitive = "%.17g", "double"
	default:
		verb, primitive = "%.17g", "double"
	}

	// freqs table
	fmt.Fprintf(&b, "static const %s freqs[%d] = {\n", primitive, len(t.freqs))
	for i := range t.freqs {
		if float64(t.freqs[i]) == math.Inf(-1) {
			w.write("-INFINITY, ")
		} else {
			w.write(fmt.Sprintf(verb+", ", t.freqs[i]))
		}
	}
	w.flush()
	b.WriteString("};\n\n")

	return b.String()
}
