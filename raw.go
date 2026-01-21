package ngrams

import (
	"math"
)

func NewRawTable(n int, total uint64, alphabet *Alphabet) *Table[uint64] {
	return newTable[uint64](n, total, alphabet)
}

func (t *Table[uint64]) ToProb() *Table[prob] {
	pt := NewProbTable(t.n, t.total, t.alphabet)
	for i, v := range t.freqs {
		pt.freqs[i] = prob(float64(v) / float64(t.total))
	}
	return pt
}

func (t *Table[uint64]) ToLogProb() *Table[logprob] {
	lpt := NewLogProbTable(t.n, t.total, t.alphabet)
	for i, v := range t.freqs {
		lpt.freqs[i] = logprob(math.Log(float64(v)) - math.Log(float64(t.total)))
	}
	return lpt
}
