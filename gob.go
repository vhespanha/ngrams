package ngrams

import (
	"bytes"
	"encoding/gob"
)

func (t *table[T]) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(t.n); err != nil {
		return nil, err
	}
	if err := enc.Encode(t.freqs); err != nil {
		return nil, err
	}
	if err := enc.Encode(t.total); err != nil {
		return nil, err
	}
	if err := enc.Encode(t.alphabet); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (t *table[T]) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&t.n); err != nil {
		return err
	}
	if err := dec.Decode(&t.freqs); err != nil {
		return err
	}
	if err := dec.Decode(&t.total); err != nil {
		return err
	}
	if err := dec.Decode(&t.alphabet); err != nil {
		return err
	}

	return nil
}
