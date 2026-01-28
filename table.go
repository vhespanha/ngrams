package ngrams

import (
	"fmt"
	"math"
)

type freq interface {
	uint64 | prob | logprob
}

type table[T freq] struct {
	n        int
	freqs    []T
	total    uint64
	alphabet *Alphabet
}

func newTable[T freq](n int, total uint64, alphabet *Alphabet) *table[T] {
	size := int(math.Pow(float64(alphabet.size), float64(n)))
	return &table[T]{
		n:        n,
		freqs:    make([]T, size),
		alphabet: alphabet,
		total:    total,
	}
}

func (t *table[T]) idx(symbols []symbol) int {
	idx := 0
	for _, s := range symbols {
		idx = idx*t.alphabet.size + int(s)
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

func (t *table[T]) Set(v T, symbols ...symbol) error {
	p, err := t.At(symbols...)
	if err != nil {
		return err
	}
	*p = v
	return nil
}

func (t *table[T]) MustSet(v T, symbols ...symbol) {
	*t.MustAt(symbols...) = v
}
