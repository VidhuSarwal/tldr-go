# SDLC Report: tldr, a command-line text summarizer

**Model:** Incremental. **Supporting documents:** [requirements.md](requirements.md) (purpose, users, features) and [design.md](design.md) (model choice, plan, design, test plan).

## 1. Chosen model and why

tldr is a small offline command-line tool that picks the most important
sentences of a text, for students skimming notes, developers reading long
documentation, and scripts that pipe text through it. I chose the
**Incremental** model with three increments, each with its own design, build,
test and integrate cycle:

1. **Core:** split the text, score by word frequency, print the top *N*.
2. **Position weighting:** earlier sentences count for more.
3. **Stopword filtering:** filler words are ignored when scoring.

- **Scope:** small, with a clear core and optional refinements. The core
  ships first, and each refinement is added without redesigning it.
- **Risk:** low technical risk. The real risks are a rule that makes
  summaries worse, and two changes to the scoring code clashing. Testing and
  integrating one increment at a time isolates both.
- **Time:** one person, one week, a hard deadline. A usable tool exists
  after the first increment.

Waterfall delivers nothing until the end, Spiral's formal risk cycles are
too heavy for a tool this size, and RAD needs user workshops and interface
prototypes that a command-line tool does not have.

## 2. Development phases

1. **Requirements:** purpose, users, scope, FR1–FR8, NFR1–NFR5, acceptance criteria.
2. **Planning:** increments mapped to requirements, estimates, milestones and one Git branch each. The integration of increments 2 and 3 was flagged as a risk in advance.
3. **Design:** a pipeline (input → split → count → score → rank → output), a `Sentence` type, and a scoring formula per increment, all in one `score()` function.
4. **Test design:** cases T1–T10 traced to increments, run as unit tests, a CLI acceptance check, and CI on every push.
5. **Integration:** increments merged into `main` one at a time, each only after the full suite passed.

## 3. Testing outcomes

The design was later built for the Git exercise, so the test plan was run:

- **T1–T5** (increment 1) passed. **T6–T7** (increment 2) and **T8–T9**
  (increment 3) passed on their own branches.
- **T10** (all tests after integrating 2 and 3) **failed to build on its
  first run**, then passed after the fix in section 4. The final suite has 9
  tests (14 results with subtests), all passing locally and on CI.
- **Acceptance:** on the sample text, the core alone picked short sentences
  full of common words ("Git took a different approach."). With all three
  increments, the summary opens with the text's actual definition of version
  control.

## 4. Conflict-resolution process

Increments 2 and 3 both rewrote the last line of `score()`. After increment 2
was merged, merging increment 3 produced a conflict on that line. Git had
already merged everything else. I resolved it in four steps:

1. Inspect with `git status` and `git diff` before editing.
2. Decide on intent, not on sides. Neither version was wrong, so the fix was a new line combining both rules: skip stopwords, divide by the words counted, then apply the position weight.
3. Run the full test suite before committing. It exposed a second, *semantic* conflict that Git could not see: increment 3's tests still called the old two-argument `score()`.
4. Fix those calls, commit the merge once everything passed, and let CI confirm it.

The ordering also settles a trade-off between the two requirements. Filtering
stopwords before applying the position boost means a short opening sentence
of filler words cannot win on position alone.

## 5. Lessons learned

- **Plan for integration, not just features.** Keeping all scoring in one
  function made the clash between increments 2 and 3 certain. Naming it in
  the plan made it cheap.
- **Tests are part of the design.** T10 caught a problem that a clean-looking
  merge hid.
- **Incremental delivery lowers deadline risk.** Every increment after the
  first was an improvement, never a prerequisite.
- **Next time:** merge `main` into the feature branch first, so conflicts are
  resolved and tested there before reaching `main`.
