package ngrams

import (
	"fmt"
	"go/format"
	"strings"
)

const maxWidth = 100

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
	b           *strings.Builder
	lineWidth   int
	maxWidth    int
	indent      string
	indentWidth int
}

func (w *wrapper) write(s string) {
	sWidth := tabWidth(s)
	if w.lineWidth > w.indentWidth && w.lineWidth+sWidth > w.maxWidth {
		w.b.WriteString("\n")
		w.lineWidth = 0
	}
	if w.lineWidth == 0 {
		w.b.WriteString(w.indent)
		w.lineWidth = w.indentWidth
	}
	w.b.WriteString(s)
	w.lineWidth += sWidth
}

func (w *wrapper) flush() {
	if w.lineWidth > 0 {
		w.b.WriteString("\n")
		w.lineWidth = 0
	}
}

func (t *table[T]) CodegenGo(pkg string) string {
	var b strings.Builder

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

	formated, err := format.Source([]byte(b.String()))
	if err != nil {
		panic("generated code has syntax errors")
	}

	return string(formated)
}
