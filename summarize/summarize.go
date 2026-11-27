// Package summarize selects the highest-scoring sentences of a text and
// returns them in their original order.
package summarize

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Sentence is one sentence of the input text.
type Sentence struct {
	Text  string   // as it appeared in the input
	Words []string // lower-case words without punctuation
	Index int      // position in the input, from 0
}

// sentenceEnd matches . ! or ? followed by whitespace.
var sentenceEnd = regexp.MustCompile(`([.!?]+)\s+`)

// Split breaks text into sentences. Line breaks inside a paragraph count as
// spaces.
func Split(text string) []Sentence {
	text = strings.Join(strings.Fields(text), " ")
	if text == "" {
		return nil
	}

	var out []Sentence
	for _, part := range strings.Split(sentenceEnd.ReplaceAllString(text, "$1\n"), "\n") {
		if part = strings.TrimSpace(part); part != "" {
			out = append(out, Sentence{Text: part, Words: words(part), Index: len(out)})
		}
	}
	return out
}

// words returns the lower-case words of s, keeping apostrophes inside words
// such as "don't".
func words(s string) []string {
	fields := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '\''
	})
	out := fields[:0]
	for _, f := range fields {
		if f = strings.Trim(f, "'"); f != "" {
			out = append(out, f)
		}
	}
	return out
}

func frequencies(sentences []Sentence) map[string]int {
	freq := make(map[string]int)
	for _, s := range sentences {
		for _, w := range s.Words {
			freq[w]++
		}
	}
	return freq
}

// score returns the average text-wide frequency of the sentence's
// non-stopwords, weighted by position. count is the number of sentences.
func score(s Sentence, freq map[string]int, count int) float64 {
	if len(s.Words) == 0 {
		return 0
	}
	total := 0.0
	counted := 0
	for _, w := range s.Words {
		if stopwords[w] {
			continue
		}
		total += float64(freq[w])
		counted++
	}
	if counted == 0 {
		return 0
	}
	// Opening sentences usually state the main idea, so they get a boost.
	return total / float64(counted) * positionWeight(s.Index, count)
}

// Summarize returns the n highest-scoring sentences of text in their original
// order. If the text has n sentences or fewer, all of them are returned.
func Summarize(text string, n int) []string {
	sentences := Split(text)
	if n <= 0 || len(sentences) == 0 {
		return nil
	}
	if n >= len(sentences) {
		return texts(sentences)
	}

	freq := frequencies(sentences)
	scores := make([]float64, len(sentences))
	for i, s := range sentences {
		scores[i] = score(s, freq, len(sentences))
	}

	// SliceStable keeps the earlier sentence first on a tie, so the same input
	// always gives the same summary.
	ranked := append([]Sentence(nil), sentences...)
	sort.SliceStable(ranked, func(i, j int) bool {
		return scores[ranked[i].Index] > scores[ranked[j].Index]
	})

	// Keep the best n, then restore reading order.
	top := ranked[:n]
	sort.Slice(top, func(i, j int) bool { return top[i].Index < top[j].Index })
	return texts(top)
}

func texts(sentences []Sentence) []string {
	out := make([]string, len(sentences))
	for i, s := range sentences {
		out[i] = s.Text
	}
	return out
}
