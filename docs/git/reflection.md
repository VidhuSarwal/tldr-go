# Reflection: the merge conflict in tldr-go

**Project:** tldr-go, a command-line text summarizer in Go
**Repository:** https://github.com/VidhuSarwal/tldr-go

## What caused the conflict

tldr scores every sentence in a text and keeps the best few. The scoring
lives in one function, `score()`. I planned two improvements and built each
on its own branch, both cut from the same commit on `main`:

- `feature/position-weighting` multiplies a sentence's score by a weight
  that favours sentences near the start of the text. To do that it added a
  `count` parameter to `score()` and rewrote the function's `return` line.
- `feature/stopwords` makes `score()` skip filler words such as "the" and
  "and". It added a `counted` variable, changed the loop, and also rewrote
  the `return` line to divide by `counted` instead of the total number of
  words.

Merging the first branch into `main` was automatic, because `main` had not
changed since the branch was cut. The second merge was different. Git
compares each branch with their common ancestor, `a645d5b`. Most of the
stopwords changes (the new variable and the loop) did not overlap anything
the first branch had done, so Git applied them on its own. But both
branches had replaced the same `return` line with different code. Git
cannot know which version is right, or whether both are, so it stopped,
marked the file as "both modified", and put both versions in the file
between `<<<<<<<`, `=======` and `>>>>>>>` markers.

## How I resolved it

I started by reading rather than editing. `git status` showed which file
was unmerged and that the stopwords branch's new files were already staged.
`git diff` showed a combined diff: the lines Git had merged cleanly, plus
the one block it could not decide.

The obvious buttons in an editor ("accept current" or "accept incoming")
would each have silently dropped a whole feature. Accepting both would have
left two `return` statements, the second unreachable. The right answer was
a new line that neither branch contained: keep the stopwords guard
(`if counted == 0 { return 0 }`), divide by `counted`, and multiply by the
position weight. In the combined diff that line is marked `++`, meaning it
differs from both parents.

I then ran `go test` before committing, and it failed to compile. The
stopwords branch had added tests that called `score(s, freq)` with two
arguments, but after the merge `score()` needs three. Git had not flagged
this because the two branches touched different lines in different files:
there was no *textual* conflict, only a *semantic* one. I fixed the two
calls, passing a count of 1 so the position weight is the same for both
sentences each test compares. Then all 14 tests passed, and I ran `git add`.
`git status` confirmed "All conflicts fixed but you are still merging", and
`git commit` created the merge commit `7750acc` with two parents. After the
push, GitHub Actions ran the tests again on the merged code and they passed.

## What I learned about distributed version control

**Every clone is a full repository.** I created the project with
`git init`, and every commit, branch and merge happened locally before
anything was pushed. GitHub was just another copy. To use `git pull`, I
pushed a commit from a second clone and pulled it into the first. Git
fast-forwarded because my copy had no new work of its own. Nothing about
the second clone was special: any copy can be the source of a change.

**A merge conflict is information, not an error.** Git merged everything it
could prove was safe and asked me only about the one place where two people
(here, two branches) made incompatible decisions. That pause is the most
useful thing Git does: it puts the decision where it belongs, with someone
who understands the code.

**Git only sees text; tests see meaning.** The textual conflict was easy to
spot. The dangerous one was the semantic conflict that Git could not see.
Had I committed as soon as the markers were gone, `main` would not have
built. Running the tests between resolving and committing, and having CI run
them again on every push, is what makes a merge trustworthy.

**Small, focused branches keep conflicts small.** Because each branch did
one thing and the scoring code lived in one function, the conflict was one
block in one file and the fix was obvious. The same overlap spread across a
week of unrelated changes would have been much harder to untangle.

**There is always a way back.** While resolving, I knew that
`git merge --abort` would return me to the state before the merge, and that
`git restore --merge` could put the markers back in a file I had edited
badly. Knowing the escape routes made it easier to resolve the conflict
calmly rather than guessing.

If I did this again, I would merge `main` into the second branch (or rebase
it onto `main`) before opening the final merge. The conflict would then be
resolved, and its tests run, on the feature branch, so `main` would only
ever receive code that had already passed.
