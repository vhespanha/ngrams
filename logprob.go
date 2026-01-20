package ngrams

import (
	"fmt"
	"math"
)

// func isNotPositive[T ~float64](v T) bool {
// 	f := float64(v)
// 	return !math.IsNaN(f) && f <= 0
// }

type logProb float64

// var errNotLogProb = "type %T should not be positive"

// func (l logProb) validate() error {
// 	if !isNotPositive(l) {
// 		return fmt.Errorf(errNotLogProb, l)
// 	}
// 	return nil
// }

func NewLogProbTable(n int, total float64, alphabet *Alphabet) *Table[logProb] {
	return newTable[logProb](n, total, alphabet)
}

func (t *Table[logProb]) SetLogProb(v float64, symbols ...symbol) error {
	if !isWhole(v) {
		return fmt.Errorf(errNotWhole, v)
	}
	return t.set(logProb(math.Log(v)-math.Log(t.total)), symbols)
}

func (t *Table[logProb]) MustSetLogProb(v float64, symbols ...symbol) {
	if !isWhole(v) {
		panic(panicNotWhole)
	}
	t.mustSet(logProb(math.Log(v)-math.Log(t.total)), symbols)
}
