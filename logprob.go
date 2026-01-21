package ngrams

import (
	"math"
)

// func isNotPositive[T ~float64](v T) bool {
// 	f := float64(v)
// 	return !math.IsNaN(f) && f <= 0
// }

type logprob float64

// var errNotLogProb = "type %T should not be positive"

// func (l logProb) validate() error {
// 	if !isNotPositive(l) {
// 		return fmt.Errorf(errNotLogProb, l)
// 	}
// 	return nil
// }

func NewLogProbTable(n int, total uint64, alphabet *Alphabet) *Table[logprob] {
	return newTable[logprob](n, total, alphabet)
}

func (t *Table[logProb]) SetLogProb(v uint64, symbols ...symbol) error {
	return t.set(logProb(math.Log(float64(v))-math.Log(float64(t.total))), symbols)
}

func (t *Table[logProb]) MustSetLogProb(v uint64, symbols ...symbol) {
	t.mustSet(logProb(math.Log(float64(v))-math.Log(float64(t.total))), symbols)
}
