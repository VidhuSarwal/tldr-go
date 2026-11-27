package summarize

import "testing"

// T8: a sentence made only of stopwords scores zero.
func TestStopwordOnlySentenceScoresZero(t *testing.T) {
	s := Sentence{Words: []string{"and", "so", "it", "is"}}
	freq := map[string]int{"and": 9, "so": 9, "it": 9, "is": 9}
	if got := score(s, freq, 1); got != 0 {
		t.Fatalf("score = %v, want 0", got)
	}
}

// T9: adding stopwords to a sentence does not change its score.
func TestStopwordsDoNotDiluteScore(t *testing.T) {
	freq := map[string]int{"branches": 2, "merge": 3, "the": 9, "of": 9, "and": 9}
	plain := Sentence{Words: []string{"branches", "merge"}}
	padded := Sentence{Words: []string{"the", "branches", "of", "merge", "and"}}
	// Both sentences sit at index 0, so the position weight is the same for each.
	if a, b := score(plain, freq, 1), score(padded, freq, 1); a != b {
		t.Fatalf("plain scored %v, padded scored %v; want equal", a, b)
	}
}
