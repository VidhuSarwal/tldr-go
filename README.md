# tldr-go

A small command-line text summarizer written in Go. Give it an article, a long
README or your lecture notes, and it prints the few sentences that matter most,
in their original order.

It is **extractive**: it never writes new sentences, so the summary only ever
says what the source says. It runs offline and uses only the Go standard
library.

## How it works

1. Split the text into sentences.
2. Count how often each word appears in the whole text.
3. Score each sentence by the average frequency of its words. Sentences made of
   the words a text repeats most are usually about its main topic.
4. Keep the best `n` sentences and print them in reading order.

## Usage

```bash
go build -o tldr .
./tldr -n 3 testdata/sample.txt
cat notes.txt | ./tldr -n 2
```

`-n` sets how many sentences to keep (default 3). With no file argument, `tldr`
reads standard input.

## Tests

```bash
go test ./...
```

## Project layout

| Path | What it is |
| --- | --- |
| `main.go` | Command-line interface: flags, reading the file or stdin, printing |
| `summarize/` | The summarizer: sentence splitting, word counts, scoring, ranking |
| `testdata/sample.txt` | Sample text to try the tool on |
| `docs/sdlc/` | Requirements and design (Incremental SDLC model) |
