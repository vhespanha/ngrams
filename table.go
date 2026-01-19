package ngrams

import (
	"fmt"
)

type freq interface {
	raw | prob | logProb
}

type table[T freq] struct {
	n     int
	freqs []T
	base  int
	total float64
}

func newTable[T freq](n, base int, total float64) *table[T] {
	size := 1
	for i := 0; i < n; i++ {
		size *= base
	}
	return &table[T]{
		n:     n,
		freqs: make([]T, size),
		base:  base,
		total: total,
	}
}

func (t *table[T]) idx(symbols []symbol) int {
	idx := 0
	for _, s := range symbols {
		idx = idx*t.base + int(s)
	}
	return idx
}

func (t *table[T]) At(symbols ...symbol) (*T, error) {
	if len(symbols) != t.n {
		return nil, fmt.Errorf("expected %d symbols, got %d",
			t.n, len(symbols))
	}
	return &t.freqs[t.idx(symbols)], nil
}

func (t *table[T]) MustAt(symbols ...symbol) *T {
	if len(symbols) != t.n {
		panic("wrong arity for ngram")
	}
	return &t.freqs[t.idx(symbols)]
}

func (t *table[T]) set(v T, symbols []symbol) error {
	p, err := t.At(symbols...)
	if err != nil {
		return err
	}
	*p = v
	return nil
}

func (t *table[T]) mustSet(v T, symbols []symbol) {
	*t.MustAt(symbols...) = v
}
