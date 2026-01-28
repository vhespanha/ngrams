package ngrams

type (
	RawTable struct {
		*table[uint64]
	}
)

func NewRawTable(n int, total uint64, alphabet *Alphabet) *RawTable {
	return &RawTable{newTable[uint64](n, total, alphabet)}
}

func (rt *RawTable) ToProb() *ProbTable {
	pt := NewProbTable(rt.n, rt.total, rt.alphabet)
	for i, v := range rt.freqs {
		pt.freqs[i] = prob(float64(v) / float64(rt.total))
	}
	return pt
}

func (rt *RawTable) ToLogProb() *LogProbTable {
	lpt := NewLogProbTable(rt.n, rt.total, rt.alphabet)
	for i, v := range rt.freqs {
		lpt.freqs[i] = newLogProb(v, rt.total)
	}
	return lpt
}
