# tldr-go

[![test](https://github.com/VidhuSarwal/tldr-go/actions/workflows/test.yml/badge.svg)](https://github.com/VidhuSarwal/tldr-go/actions/workflows/test.yml)

A small command-line text summarizer written in Go. Give it an article, a long
README or your lecture notes, and it prints the few sentences that matter most,
in their original order.

It is **extractive**: it selects existing sentences rather than writing new
ones. Selected sentences can still lose context, so a summary should be
reviewed before it is relied on. It runs offline and uses only the Go
standard library.

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
| `docs/sdlc/` | Requirements, design and report (Incremental SDLC model) |
| `docs/git/` | How this repository was built, including the merge conflict |
| `linux-user-management/` | A separate exercise: user management on Amazon Linux |

## Assignment index

| Question | Where to look |
| --- | --- |
| **Git project with an intentional merge conflict** | This code and its history: branches `feature/position-weighting` and `feature/stopwords`, merge commit `7750acc`. Steps: [docs/git/steps.md](docs/git/steps.md). Reflection: [docs/git/reflection.md](docs/git/reflection.md). Screenshots: [docs/screenshots/](docs/screenshots/) |
| **SDLC design (Incremental model)** | Report: [docs/sdlc/report.md](docs/sdlc/report.md). Requirements: [requirements.md](docs/sdlc/requirements.md). Plan, design and test plan: [design.md](docs/sdlc/design.md) |
| **Linux user management on Amazon Linux** | [linux-user-management/README.md](linux-user-management/README.md) |
