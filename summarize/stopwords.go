package summarize

import "strings"

// stopwords are common English words that appear in almost every sentence,
// so they say nothing about what a particular text is about.
var stopwords = set(`
	a about after all also an and any are as at be been but by can could did
	do does for from had has have he her his how i if in into is it its just
	me more my no not of on or our out she so some than that the their them
	then there these they this to up was we were what when where which who
	will with would you your
`)

// set turns a whitespace-separated list of words into a lookup set.
func set(list string) map[string]bool {
	out := make(map[string]bool)
	for _, w := range strings.Fields(list) {
		out[w] = true
	}
	return out
}
