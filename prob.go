package ngrams

// func isFactor[T ~float64](v T) bool {
// 	f := float64(v)
// 	return f >= 0 && f <= 1
// }

type prob float64

// var errNotProb = "type %T must be within [0, 1]"

// func (p prob) validate() error {
// 	if !isFactor(p) {
// 		return fmt.Errorf(errNotProb, p)
// 	}
// 	return nil
// }

func NewProbTable(n int, total uint64, alphabet *Alphabet) *Table[prob] {
	return newTable[prob](n, total, alphabet)
}

func (t *Table[prob]) SetProbFromCount(v uint64, symbols ...symbol) error {
	return t.Set(prob(v/t.total), symbols)
}

func (t *Table[prob]) MustSetProbFromCount(v uint64, symbols ...symbol) {
	t.MustSet(prob(v/t.total), symbols)
}
