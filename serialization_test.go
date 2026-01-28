package ngrams

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
)

// Tests for JSON marshalling - only tests marshalling since unmarshalling
// has limitations with the Alphabet struct's unexported fields.

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

	// Verify JSON structure contains expected fields
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"n":2`) {
		t.Error("JSON output missing n field")
	}
	if !strings.Contains(jsonStr, `"total":1000`) {
		t.Error("JSON output missing total field")
	}
	if !strings.Contains(jsonStr, `"freqs"`) {
		t.Error("JSON output missing freqs field")
	}
	if !strings.Contains(jsonStr, `"alphabet"`) {
		t.Error("JSON output missing alphabet field")
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

	// Verify JSON structure contains expected fields
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"n":1`) {
		t.Error("JSON output missing n field")
	}
	if !strings.Contains(jsonStr, `"total":100`) {
		t.Error("JSON output missing total field")
	}
	if !strings.Contains(jsonStr, `"freqs"`) {
		t.Error("JSON output missing freqs field")
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
		t.Skip("Expected JSON marshal to fail with -Inf values, but it succeeded")
	}
	// Verify it's the expected error
	if !strings.Contains(err.Error(), "-Inf") && !strings.Contains(err.Error(), "unsupported value") {
		t.Errorf("Unexpected error: %v", err)
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
