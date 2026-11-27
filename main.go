// Command tldr prints the most important sentences of a file or standard
// input.
//
//	tldr -n 3 article.txt
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/VidhuSarwal/tldr-go/summarize"
)

func main() {
	n := flag.Int("n", 3, "number of sentences to keep")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: tldr [-n sentences] [file]")
		fmt.Fprintln(os.Stderr, "Reads standard input when no file is given.")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(2)
	}

	text, err := readInput(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "tldr:", err)
		os.Exit(1)
	}

	for _, sentence := range summarize.Summarize(text, *n) {
		fmt.Println(sentence)
	}
}

// readInput reads the named file, or standard input when path is empty.
func readInput(path string) (string, error) {
	var r io.Reader = os.Stdin
	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()
		r = f
	}
	data, err := io.ReadAll(r)
	return string(data), err
}
