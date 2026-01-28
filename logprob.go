package ngrams

import (
	"fmt"
	"math"
)

type (
	logprob      float64
	LogProbTable struct {
		*table[logprob]
	}
)

var logProbZeroValue = logprob(math.Inf(-1))

func (lp logprob) min() logprob   { return logProbZeroValue }
func (lp logprob) max() logprob   { return 0 }
func (lp logprob) validate() bool { return !math.IsNaN(float64(lp)) && lp <= lp.max() }

func newLogProb(v, total uint64) logprob {
	return logprob(math.Log(float64(v)) - math.Log(float64(total)))
}

func NewLogProbTable(n int, total uint64, alphabet *Alphabet) *LogProbTable {
	t := newTable[logprob](n, total, alphabet)
	for i := range t.freqs {
		t.freqs[i] = logProbZeroValue
	}
	return &LogProbTable{t}
}

func (lpt *LogProbTable) SetLogProbFromCount(v uint64, symbols ...symbol) error {
	lp := newLogProb(v, lpt.total)
	if !lp.validate() {
		return fmt.Errorf("%T should be [%g,%g], is %g", lp, lp.min(), lp.max(), lp)
	}
	return lpt.Set(lp, symbols...)
}

func (lpt *LogProbTable) MustSetLogProbFromCount(v uint64, symbols ...symbol) {
	lp := newLogProb(v, lpt.total)
	if !lp.validate() {
		panic("operation generated invalid log probability")
	}
	lpt.MustSet(lp, symbols...)
}
