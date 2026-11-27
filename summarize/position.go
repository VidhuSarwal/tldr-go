package summarize

// positionWeight falls linearly from 1.5 for the first sentence to 1.0 for
// the last.
func positionWeight(index, count int) float64 {
	if count <= 1 {
		return 1.5
	}
	return 1.5 - 0.5*float64(index)/float64(count-1)
}
