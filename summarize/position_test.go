package summarize

import "testing"

// T6: the first sentence gets 1.5x, the last gets 1.0x.
func TestPositionWeight(t *testing.T) {
	cases := []struct {
		index, count int
		want         float64
	}{
		{0, 5, 1.5},
		{2, 5, 1.25},
		{4, 5, 1.0},
		{0, 1, 1.5},
	}
	for _, c := range cases {
		if got := positionWeight(c.index, c.count); got != c.want {
			t.Errorf("positionWeight(%d, %d) = %v, want %v", c.index, c.count, got, c.want)
		}
	}
}

// T7: the same sentence scores higher at the start of a text than at the end.
func TestEarlierSentenceScoresHigher(t *testing.T) {
	sentences := Split("Merging joins two branches. Something else happens here. Merging joins two branches.")
	freq := frequencies(sentences)
	first := score(sentences[0], freq, len(sentences))
	last := score(sentences[2], freq, len(sentences))
	if first <= last {
		t.Fatalf("first copy scored %v, last copy %v; want first > last", first, last)
	}
}
