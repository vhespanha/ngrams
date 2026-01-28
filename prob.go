package ngrams

import "fmt"

type (
	prob      float64
	ProbTable struct {
		*table[prob]
	}
)

func (p prob) min() prob      { return 0 }
func (p prob) max() prob      { return 1 }
func (p prob) validate() bool { return p >= p.min() && p <= p.max() }

func NewProbTable(n int, total uint64, alphabet *Alphabet) *ProbTable {
	return &ProbTable{newTable[prob](n, total, alphabet)}
}

func (pt *ProbTable) SetProbFromCount(v uint64, symbols ...symbol) error {
	p := prob(float64(v) / float64(pt.total))
	if !p.validate() {
		return fmt.Errorf("%T should be [%g,%g], is %g", p, p.min(), p.max(), p)
	}
	return pt.Set(p, symbols...)
}

func (pt *ProbTable) MustSetProbFromCount(v uint64, symbols ...symbol) {
	p := prob(float64(v) / float64(pt.total))
	if !p.validate() {
		panic("operation resulted in invalid prob")
	}
	pt.MustSet(p, symbols...)
}
