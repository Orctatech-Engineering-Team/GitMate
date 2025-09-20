# GITMATE
### 1. Identity & Purpose

* **Name**: GitMate
* **Tagline**: *“Your friendly Git workflow companion.”*
* **Core Purpose**: Help individuals and teams practice clean, disciplined Git workflows (rebasing, squashing, branching, PR prep) without confusion.

### 2. Minimal First Version

* Teach (through prompts/explanations).
* Guide (step-by-step workflows like squash or rebase).
* Automate (run the actual Git commands when user confirms).

Example:

```bash
gitmate squash
```

→ Shows you your last commits, asks which to squash, explains why squashing matters, then executes.

### 3. Roadmap (tiny steps for your ADHD flow)

1. **Setup project skeleton** with Bubble Tea.
2. **Add one command**: `gitmate squash`.
3. Test with your own repo.
4. **Polish the UX** (clear prompts, colors, short explanations).
5. Share with your team for feedback.

Later, we can add:

* `gitmate rebase` (guided flow).
* `gitmate clean-branch`.
* `gitmate prepare-pr`.

Do you want me to draft the **exact step 1 setup (repo + bubble tea skeleton)** so you can immediately open your editor and start GitMate today?
