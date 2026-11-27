// Package summarize picks the most important sentences out of a block of text.
//
// It is an extractive summarizer: it never writes new sentences. It scores the
// sentences that are already there and returns the best few, in the order they
// appeared.
package summarize

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// Sentence is one sentence of the input text.
type Sentence struct {
	Text  string   // the sentence exactly as it appeared in the input
	Words []string // lower-cased words with punctuation removed
	Index int      // position of the sentence in the input, starting at 0
}

// sentenceEnd matches the end of a sentence: one or more of . ! ? followed
// by whitespace.
var sentenceEnd = regexp.MustCompile(`([.!?]+)\s+`)

// Split breaks text into sentences. Line breaks inside a paragraph are
// treated as spaces, so wrapped text splits the same way as unwrapped text.
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

// words lower-cases s and returns its words, dropping punctuation.
// Apostrophes inside a word ("don't") are kept.
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

// frequencies counts how often each word appears across all sentences.
func frequencies(sentences []Sentence) map[string]int {
	freq := make(map[string]int)
	for _, s := range sentences {
		for _, w := range s.Words {
			freq[w]++
		}
	}
	return freq
}

// score rates a sentence by the average frequency of its words. Sentences
// built from the words a text repeats most are usually about its main topic.
// count is the number of sentences in the whole text.
func score(s Sentence, freq map[string]int, count int) float64 {
	if len(s.Words) == 0 {
		return 0
	}
	total := 0.0
	counted := 0
	for _, w := range s.Words {
		if stopwords[w] {
			continue // "the", "and", "of" say nothing about the topic
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

// Summarize returns the n highest-scoring sentences of text, in the order they
// appear in the text. If the text has n sentences or fewer, all of them are
// returned unchanged.
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

	// Rank a copy by score, best first. SliceStable keeps the earlier sentence
	// ahead on a tie, so the same input always gives the same summary.
	ranked := append([]Sentence(nil), sentences...)
	sort.SliceStable(ranked, func(i, j int) bool {
		return scores[ranked[i].Index] > scores[ranked[j].Index]
	})

	// Keep the best n, then put them back in reading order.
	top := ranked[:n]
	sort.Slice(top, func(i, j int) bool { return top[i].Index < top[j].Index })
	return texts(top)
}

// texts returns the Text of each sentence.
func texts(sentences []Sentence) []string {
	out := make([]string, len(sentences))
	for i, s := range sentences {
		out[i] = s.Text
	}
	return out
}
