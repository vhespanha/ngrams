package ngrams

import (
	"fmt"
	"math"
)

func isWhole[T ~float64](v T) bool {
	f := float64(v)
	return f >= 0 && !math.IsInf(f, 0) && !math.IsNaN(f) && f == math.Trunc(f)
}

type raw float64

var (
	errNotWhole   = "type %T should be a whole number"
	panicNotWhole = "not a whole number"
)

func (r raw) validate() error {
	if !isWhole(r) {
		return fmt.Errorf(errNotWhole, r)
	}
	return nil
}

func NewRawTable(n int, total float64, alphabet *Alphabet) *Table[raw] {
	return newTable[raw](n, total, alphabet)
}

func (t *Table[raw]) SetRaw(v float64, symbols ...symbol) error {
	if !isWhole(v) {
		return fmt.Errorf(errNotWhole, v)
	}
	return t.set(raw(v), symbols)
}

func (t *Table[raw]) MustSetRaw(v float64, symbols ...symbol) {
	if !isWhole(v) {
		panic(panicNotWhole)
	}
	t.mustSet(raw(v), symbols)
}

func (t *Table[raw]) ToProb() *Table[prob] {
	pt := NewProbTable(t.n, t.total, t.alphabet)
	for i, v := range t.freqs {
		pt.freqs[i] = prob(float64(v) / t.total)
	}
	return pt
}

func (t *Table[raw]) ToLogProb() *Table[logProb] {
	lpt := NewLogProbTable(t.n, t.total, t.alphabet)
	for i, v := range t.freqs {
		lpt.freqs[i] = logProb(math.Log(float64(v)) - math.Log(t.total))
	}
	return lpt
}
