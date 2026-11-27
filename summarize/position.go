package summarize

// positionWeight boosts sentences near the start of a text. The first
// sentence gets 1.5x its score, falling in a straight line to 1.0x for the
// last sentence. index is the sentence's position and count is the number of
// sentences in the text.
func positionWeight(index, count int) float64 {
	if count <= 1 {
		return 1.5
	}
	return 1.5 - 0.5*float64(index)/float64(count-1)
}
