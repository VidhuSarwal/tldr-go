package summarize

import (
	"reflect"
	"testing"
)

// T1: sentences end at . ! or ? followed by whitespace.
func TestSplitSentences(t *testing.T) {
	got := texts(Split("Go is fast. It compiles quickly! Does it scale? Yes."))
	want := []string{"Go is fast.", "It compiles quickly!", "Does it scale?", "Yes."}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Split = %q, want %q", got, want)
	}
}

// T2: a line break inside a sentence is just a space.
func TestSplitJoinsWrappedLines(t *testing.T) {
	got := texts(Split("This sentence is\nwrapped across lines.\n\nThis one is not."))
	want := []string{"This sentence is wrapped across lines.", "This one is not."}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Split = %q, want %q", got, want)
	}
}

// T3: the summary has exactly n sentences, in the order of the original text.
func TestSummarizeKeepsOriginalOrder(t *testing.T) {
	text := "Git records history. Branches let work happen in parallel. " +
		"Merging joins branches back together. Git merges most branches automatically. " +
		"The office has a coffee machine."
	got := Summarize(text, 3)
	if len(got) != 3 {
		t.Fatalf("got %d sentences, want 3: %q", len(got), got)
	}

	position := map[string]int{}
	for _, s := range Split(text) {
		position[s.Text] = s.Index
	}
	for i := 1; i < len(got); i++ {
		if position[got[i-1]] > position[got[i]] {
			t.Fatalf("summary is out of order: %q", got)
		}
	}
}

// T4: empty text, n <= 0, and n larger than the text.
func TestSummarizeEdgeCases(t *testing.T) {
	short := "One sentence. Two sentences."
	cases := []struct {
		name string
		text string
		n    int
		want []string
	}{
		{"empty text", "", 3, nil},
		{"whitespace only", " \n\t ", 3, nil},
		{"n is zero", short, 0, nil},
		{"n is negative", short, -1, nil},
		{"n larger than text", short, 5, []string{"One sentence.", "Two sentences."}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Summarize(c.text, c.n); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Summarize(%q, %d) = %q, want %q", c.text, c.n, got, c.want)
			}
		})
	}
}

// T5: a sentence about the text's main topic beats an off-topic one.
func TestSummarizePrefersTopicSentence(t *testing.T) {
	text := "Git stores history. Git branches are cheap. Git merges branches. Lunch was nice."
	got := Summarize(text, 1)
	if len(got) != 1 || got[0] == "Lunch was nice." {
		t.Fatalf("Summarize picked %q, want a sentence about Git", got)
	}
}
