package ngrams

import (
	"encoding/json"
	"math"
	"testing"
)

// Tests for JSON marshalling - only tests marshalling since unmarshalling
// has limitations with the Alphabet struct's unexported fields.

// jsonTableAux is used for parsing JSON output to verify structure
type jsonTableAux struct {
	N        int             `json:"n"`
	Total    uint64          `json:"total"`
	Freqs    json.RawMessage `json:"freqs"`
	Alphabet json.RawMessage `json:"alphabet"`
}

func TestRawTableJSONMarshal(t *testing.T) {
	a := NewLowercaseAlphabet()
	rt := NewRawTable(2, 1000, a)
	rt.MustSet(100, 0, 0)
	rt.MustSet(200, 1, 2)

	// Marshal
	data, err := json.Marshal(rt)
	if err != nil {
		t.Fatalf("JSON marshal error: %v", err)
	}

	// Parse JSON to verify structure
	var parsed jsonTableAux
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if parsed.N != 2 {
		t.Errorf("JSON n = %d, want 2", parsed.N)
	}
	if parsed.Total != 1000 {
		t.Errorf("JSON total = %d, want 1000", parsed.Total)
	}
	if parsed.Freqs == nil {
		t.Error("JSON freqs is nil")
	}
	if parsed.Alphabet == nil {
		t.Error("JSON alphabet is nil")
	}
}

func TestProbTableJSONMarshal(t *testing.T) {
	a := NewLowercaseAlphabet()
	pt := NewProbTable(1, 100, a)
	pt.MustSetProbFromCount(50, 0)

	// Marshal
	data, err := json.Marshal(pt)
	if err != nil {
		t.Fatalf("JSON marshal error: %v", err)
	}

	// Parse JSON to verify structure
	var parsed jsonTableAux
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if parsed.N != 1 {
		t.Errorf("JSON n = %d, want 1", parsed.N)
	}
	if parsed.Total != 100 {
		t.Errorf("JSON total = %d, want 100", parsed.Total)
	}
	if parsed.Freqs == nil {
		t.Error("JSON freqs is nil")
	}
}

func TestLogProbTableJSONMarshal(t *testing.T) {
	a := NewLowercaseAlphabet()
	lpt := NewLogProbTable(1, 100, a)
	// Set a value so there's at least some non-Inf data
	lpt.MustSetLogProbFromCount(50, 0)

	// Note: JSON marshalling will fail when there are -Inf values
	// (which are the default for unset entries in LogProbTable)
	// This is a known limitation of JSON encoding.
	_, err := json.Marshal(lpt)
	if err == nil {
		t.Fatalf("Expected JSON marshal to fail with -Inf values, but it succeeded")
	}
}

func TestTableMarshalJSON(t *testing.T) {
	a := NewLowercaseAlphabet()

	t.Run("raw table roundtrip via table", func(t *testing.T) {
		// Create original table
		rt := NewRawTable(1, 100, a)
		rt.MustSet(50, 0) // 'a' = 50
		rt.MustSet(30, 1) // 'b' = 30

		// Marshal
		data, err := rt.table.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON error: %v", err)
		}

		// Create new table for unmarshal
		newRT := NewRawTable(1, 0, a)
		err = newRT.table.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("UnmarshalJSON error: %v", err)
		}

		// Verify data was restored
		if newRT.n != 1 {
			t.Errorf("decoded.n = %d, want 1", newRT.n)
		}
		if newRT.total != 100 {
			t.Errorf("decoded.total = %d, want 100", newRT.total)
		}
		if *newRT.MustAt(0) != 50 {
			t.Errorf("decoded[0] = %d, want 50", *newRT.MustAt(0))
		}
		if *newRT.MustAt(1) != 30 {
			t.Errorf("decoded[1] = %d, want 30", *newRT.MustAt(1))
		}
	})

	t.Run("prob table roundtrip via table", func(t *testing.T) {
		// Create original table
		pt := NewProbTable(1, 100, a)
		pt.MustSetProbFromCount(50, 0)

		// Marshal
		data, err := pt.table.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON error: %v", err)
		}

		// Create new table for unmarshal
		newPT := NewProbTable(1, 0, a)
		err = newPT.table.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("UnmarshalJSON error: %v", err)
		}

		// Verify data was restored
		got := float64(*newPT.MustAt(0))
		if math.Abs(got-0.5) > 1e-10 {
			t.Errorf("decoded[0] = %v, want 0.5", got)
		}
	})

	t.Run("logprob table with set values roundtrip", func(t *testing.T) {
		// Create original table
		lpt := NewLogProbTable(1, 100, a)

		// Set all values to avoid -Inf (which JSON cannot encode)
		for i := 0; i < a.size; i++ {
			lpt.MustSetLogProbFromCount(uint64(i+1), symbol(i))
		}

		// Marshal
		data, err := lpt.table.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON error: %v", err)
		}

		// Create new table for unmarshal
		newLPT := NewLogProbTable(1, 0, a)
		err = newLPT.table.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("UnmarshalJSON error: %v", err)
		}

		// Verify data was restored
		got := float64(*newLPT.MustAt(0))
		expected := newLogProb(1, 100)
		if math.Abs(got-float64(expected)) > 1e-10 {
			t.Errorf("decoded[0] = %v, want %v", got, expected)
		}
	})
}
