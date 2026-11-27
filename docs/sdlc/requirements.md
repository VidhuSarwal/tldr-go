# tldr: Requirements

## 1. Purpose

Students and developers regularly face long texts (articles, lecture notes,
documentation, issue threads) when they only need the main points. **tldr** is
a command-line tool that picks the most important sentences of a text and
prints them as a short summary.

It is extractive: it selects sentences from the source rather than writing new
ones. The summary therefore cannot state anything the source does not, and the
tool needs no AI model, no network and no account.

## 2. Users

| User | What they need from tldr |
| --- | --- |
| Student | Skim lecture notes or a reading before class or an exam |
| Developer | Get the gist of a long README, changelog or issue thread without leaving the terminal |
| Scripts and other programs | Pipe text in and get a summary out (standard input and output) |

## 3. Scope

**In scope:** English plain text; a command-line interface; fully offline use.

**Out of scope:** AI-written (abstractive) summaries; PDF, Word or web-page
input; a graphical interface; languages other than English.

## 4. Functional requirements

| ID | Requirement | Increment |
| --- | --- | --- |
| FR1 | Read text from a file given as an argument, or from standard input when no file is given. | 1 |
| FR2 | Split text into sentences at `.`, `!` or `?` followed by whitespace. Line breaks inside a paragraph count as spaces. | 1 |
| FR3 | Score each sentence by how often its words appear in the whole text. | 1 |
| FR4 | Output the top *N* sentences (default 3, set with `-n`) in their original order. | 1 |
| FR5 | If the text has *N* sentences or fewer, output all of them unchanged. | 1 |
| FR6 | Give sentences near the start of the text extra weight, since opening sentences usually state the main idea. | 2 |
| FR7 | Ignore common filler words ("the", "and", "of") when scoring. | 3 |
| FR8 | If the file cannot be read, print a clear error and exit with a non-zero status. | 1 |

## 5. Non-functional requirements

| ID | Requirement |
| --- | --- |
| NFR1 | **Performance:** summarize a 10,000-word text in under one second on a laptop. |
| NFR2 | **Portability:** a single binary for macOS, Linux and Windows, built from the Go standard library only. |
| NFR3 | **Determinism:** the same input always gives the same output. Ties go to the earlier sentence. |
| NFR4 | **Maintainability:** every scoring rule has unit tests, and the tests run automatically on every push. |
| NFR5 | **Privacy:** works offline; the text never leaves the machine. |

## 6. Constraints and assumptions

- One developer, about one week of part-time work, with a fixed deadline.
- Go 1.26, standard library only.
- Sentence splitting is rule-based, so an abbreviation such as "e.g." may end a
  sentence early. This is an accepted limitation for this version.

## 7. Acceptance criteria

- `tldr -n 3 testdata/sample.txt` prints exactly three sentences from the
  sample, in their original order.
- `go test ./...` passes, locally and on GitHub Actions, for every merge into
  `main`.
