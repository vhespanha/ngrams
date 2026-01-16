package ngrams

import (
	"fmt"
)

type nGramTable struct {
	N    int
	Data []float64
	base int
}

func NewNGramTable(n, base int) *nGramTable {
	size := 1
	for i := 0; i < n; i++ {
		size *= base
	}
	return &nGramTable{N: n, Data: make([]float64, size), base: base}
}

func (t *nGramTable) idx(symbols []symbol) int {
	idx := 0
	for _, s := range symbols {
		idx = idx*t.base + int(s)
	}
	return idx
}

func (t *nGramTable) At(symbols ...symbol) (*float64, error) {
	if len(symbols) != t.N {
		return nil, fmt.Errorf("expected %d symbols, got %d",
			t.N, len(symbols))
	}
	return &t.Data[t.idx(symbols)], nil
}

func (t *nGramTable) MustAt(symbols ...symbol) *float64 {
	if len(symbols) != t.N {
		panic("wrong arity for ngram")
	}
	return &t.Data[t.idx(symbols)]
}

func (t *nGramTable) Set(v float64, symbols ...symbol) error {
	p, err := t.At(symbols...)
	if err != nil {
		return err
	}
	*p = v
	return nil
}

func (t *nGramTable) MustSet(v float64, symbols ...symbol) {
	*t.MustAt(symbols...) = v
}
