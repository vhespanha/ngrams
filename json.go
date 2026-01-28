package ngrams

import "encoding/json"

type aux[T freq] struct {
	N        int       `json:"n"`
	Freqs    []T       `json:"freqs"`
	Total    uint64    `json:"total"`
	Alphabet *Alphabet `json:"alphabet"`
}

func (t *table[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(aux[T]{t.n, t.freqs, t.total, t.alphabet})
}

func (t *table[T]) UnmarshalJSON(data []byte) error {
	aux := aux[T]{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.n = aux.N
	t.freqs = aux.Freqs
	t.total = aux.Total
	t.alphabet = aux.Alphabet

	return nil
}
