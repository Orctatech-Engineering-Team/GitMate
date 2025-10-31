# **GitMate**

> *“Your friendly Git workflow companion.”*

## **1. What Is GitMate?**

**GitMate** helps you and your team practice clean, disciplined Git workflows — without getting lost in the commands.
It teaches, guides, and automates common Git tasks like **rebasing**, **squashing**, and **branch prep** through a simple interactive CLI.

Whether you’re cleaning up a messy commit history or preparing for a pull request, GitMate keeps your workflow neat and understandable.

## **2. Why GitMate?**

Working with Git can be powerful — but sometimes confusing. GitMate exists to:

* **Teach** best practices (rebasing vs. merging, when to squash, etc.)
* **Guide** you through step-by-step actions.
* **Automate** the final Git commands safely — only when you confirm.

You learn as you go, while GitMate does the heavy lifting.

## **3. Quick Example**

```bash
gitmate squash
```

GitMate will:

1. Show your recent commits.
2. Ask which to squash (with a short explanation of why squashing matters).
3. Execute the action when you confirm.

Clean history, confident workflow.

## **4. Minimal First Version (MVP)**

* [x] Setup project skeleton (Bubble Tea).
* [x] Add one command: `gitmate squash`.
* [x] Test on local repo.
* [x] Polish the UX (clear prompts, colors, concise explanations).
* [ ] Share with team for feedback.

## **5. Roadmap (Coming Next)**

* [ ] Add `gitmate rebase` and `gitmate branch`.
* [ ] Introduce **Explain Mode** (`--explain`) for deeper learning.
* [ ] Add **Safe Mode** (`--dry`) for simulations.
* [ ] Team feedback → refine UX & add more workflows.

## **6. Contributing**

Feedback, issues, or suggestions are welcome — GitMate’s meant to grow with how *real people* use Git.
Open a PR or start a discussion.
