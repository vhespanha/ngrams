package ngrams

import "encoding/json"

type auxTable[T freq] struct {
	N        int       `json:"n"`
	Freqs    []T       `json:"freqs"`
	Total    uint64    `json:"total"`
	Alphabet *Alphabet `json:"alphabet"`
}

func (t *Table[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(auxTable[T]{t.n, t.freqs, t.total, t.alphabet})
}

func (t *Table[T]) UnmarshalJSON(data []byte) error {
	aux := auxTable[T]{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.n = aux.N
	t.freqs = aux.Freqs
	t.total = aux.Total
	t.alphabet = aux.Alphabet

	return nil
}
