package ngrams

import (
	"fmt"
)

type freq interface {
	raw | prob | logProb
}

type Table[T freq] struct {
	n        int
	freqs    []T
	total    float64
	alphabet *Alphabet
}

func newTable[T freq](n int, total float64, alphabet *Alphabet) *Table[T] {
	size := 1
	return &Table[T]{
		n:        n,
		freqs:    make([]T, size),
		alphabet: alphabet,
		total:    total,
	}
}

func (t *Table[T]) idx(symbols []symbol) int {
	idx := 0
	for _, s := range symbols {
		idx = idx*t.alphabet.size + int(s)
	}
	return idx
}

func (t *Table[T]) At(symbols ...symbol) (*T, error) {
	if len(symbols) != t.n {
		return nil, fmt.Errorf("expected %d symbols, got %d",
			t.n, len(symbols))
	}
	return &t.freqs[t.idx(symbols)], nil
}

func (t *Table[T]) MustAt(symbols ...symbol) *T {
	if len(symbols) != t.n {
		panic("wrong arity for ngram")
	}
	return &t.freqs[t.idx(symbols)]
}

func (t *Table[T]) set(v T, symbols []symbol) error {
	p, err := t.At(symbols...)
	if err != nil {
		return err
	}
	*p = v
	return nil
}

func (t *Table[T]) mustSet(v T, symbols []symbol) {
	*t.MustAt(symbols...) = v
}
