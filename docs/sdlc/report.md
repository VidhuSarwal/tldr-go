# SDLC Report: tldr, a command-line text summarizer

**Model:** Incremental. **Supporting documents:** [requirements.md](requirements.md) (purpose, users, features) and [design.md](design.md) (model choice, plan, design, test plan).

## 1. The problem and the chosen model

Students and developers often need the main points of a long text quickly.
tldr is a small offline command-line tool that selects the most important
sentences of a text and prints them in their original order. It serves
three kinds of user: students skimming notes, developers reading long
documentation, and scripts that pipe text in and out.

I chose the **Incremental** model and split the work into three planned
increments (plus one for the future). Each one runs its own design, build,
test and integrate cycle and leaves a working tool behind it:

1. **Core:** split, score by word frequency, print the top *N*, CLI.
2. **Position weighting:** earlier sentences count for more.
3. **Stopword filtering:** filler words are ignored when scoring.

**Why Incremental fits:**

- **Scope:** small, with a clear core and optional refinements. The core
  ships first, and each refinement is added without redesigning it.
- **Risk:** low technical risk overall. The real risks are a scoring rule
  that makes summaries worse, and two changes to the scoring code clashing.
  Testing and integrating one increment at a time isolates both.
- **Time:** one person, about a week, with a hard deadline. A usable tool
  exists after the first increment.

Waterfall would have delivered nothing until the end and left all
integration to one late phase. Spiral's formal risk cycles are too heavy for
a tool this size. RAD depends on user workshops and interface prototyping,
and tldr has no graphical interface.

## 2. Development phases

| Phase | What was produced |
| --- | --- |
| **Requirements** | Purpose, users and scope; functional requirements FR1–FR8 and non-functional requirements NFR1–NFR5; acceptance criteria. |
| **Planning** | Increment plan with estimates and milestones; each increment mapped to requirements and to its own Git branch. Increments 2 and 3 were planned in parallel, with their integration flagged in advance as a risk. |
| **Design** | Pipeline architecture (input → split → count → score → rank → output), a `Sentence` data type, and the scoring formula for each increment. All scoring lives in one function, so each increment changes exactly one place. |
| **Test design** | Test cases T1–T10 traced to increments, at three levels: unit tests, a CLI acceptance check on a sample text, and regression testing on every push through CI. |
| **Integration** | Increments merged into `main` one at a time. Every merge had to pass the full test suite. |

## 3. Testing outcomes

The design was taken through to code as part of the Git exercise, so the
test plan could be executed. Results:

| Tests | Increment | Result |
| --- | --- | --- |
| T1–T5 (splitting, ordering, edge cases, topic choice) | 1 | Pass |
| T6–T7 (position weight values; earlier copy scores higher) | 2 | Pass on its branch |
| T8–T9 (stopword-only sentence scores 0; stopwords do not dilute) | 3 | Pass on its branch |
| T10 (all tests after integrating 2 and 3) | 2 + 3 | **First run failed to build**; passed after the fix in section 4 |

The final state is 9 test functions (14 results including subtests), all
passing locally and on GitHub Actions for every push.

**Acceptance check** (`tldr -n 3` on the sample text about version control):

| After increment | Sentences chosen |
| --- | --- |
| 1 (core) | "Git took a different approach." · "Branching in Git is cheap…" · "That pause is called a merge conflict…" |
| 1 + 2 + 3 | "Version control is a system that records changes to files over time…" · "Git took a different approach." · "In Git, every clone is a complete copy of the repository…" |

The core alone favoured short sentences full of common words. With both
refinements, the summary opens with the text's actual definition of version
control and includes the key idea of distributed version control.

## 4. Conflict-resolution process

**Integration conflict.** Increments 2 and 3 were built in parallel and
both rewrote the last line of `score()`. Increment 2 was merged first. When
increment 3 was merged, Git applied its non-overlapping changes
automatically but stopped on that line with a conflict. The process I
followed:

1. Inspect with `git status` and `git diff` before changing anything.
2. Decide on intent, not on sides. Neither version was wrong, so the
   resolution was a new line combining both rules: skip stopwords, divide
   by the words actually counted, then apply the position weight.
3. Run the full test suite before committing. It exposed a second,
   *semantic* conflict that Git could not see: increment 3's tests called
   the old two-argument `score()`. I fixed those calls.
4. Commit the merge only once everything passed, then let CI confirm it.

**Requirement trade-off.** Position weighting (FR6) and stopword filtering
(FR7) could pull against each other: a short opening sentence made mostly of
stopwords would get a large position boost on almost no content. Applying
the stopword filter first, and returning 0 when nothing is left, settles
this. The boost multiplies real content, never filler.

## 5. Lessons learned

- **Plan for integration, not just for features.** Putting all scoring in one
  function kept the design simple, but it made a conflict between
  increments 2 and 3 certain. Because the plan named it in advance, it cost
  minutes rather than a surprise.
- **Tests are part of the design.** Writing T1–T10 before coding gave each
  increment a clear definition of done. The integration test (T10) caught a
  problem that a clean-looking merge hid.
- **Incremental delivery reduces deadline risk.** A working summarizer
  existed after increment 1. Every later increment was an improvement, never
  a prerequisite.
- **Next time:** bring `main` into a feature branch before the final merge,
  so integration conflicts are resolved and tested on the branch, and `main`
  only ever receives code that has already passed.
