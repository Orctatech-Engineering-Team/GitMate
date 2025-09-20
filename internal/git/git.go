/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package git

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// FileStatus represents one entry from `git status --porcelain`
type FileStatus struct {
	IndexStatus    byte   // first status char
	WorktreeStatus byte   // second status char
	Path           string // path shown (target path for rename/copy)
	OrigPath       string // original path for rename/copy (empty if not rename)
	RawLine        string // original porcelain line
}

// Run runs `git <args...>` in dir with a default timeout and returns stdout, stderr and error.
func Run(ctx context.Context, dir string, args ...string) (stdout string, stderr string, err error) {
	// If caller didn't provide a context deadline, add a sensible timeout to avoid hanging.
	// Caller can provide ctx with its own deadline to override.
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, "git", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	stdout = strings.TrimRight(outBuf.String(), "\n")
	stderr = strings.TrimRight(errBuf.String(), "\n")

	// Wrap the error with stderr content for better debugging
	if err != nil {
		if stderr != "" {
			err = fmt.Errorf("%w: %s", err, stderr)
		}
	}
	return
}

// RunCombined runs git and returns combined output and error.
func RunCombined(ctx context.Context, dir string, args ...string) (string, error) {
	out, _, err := Run(ctx, dir, args...)
	return out, err
}

// GitStatusPorcelain runs `git status --porcelain` in the provided dir and returns parsed entries.
func GitStatusPorcelain(dir string) ([]FileStatus, error) {
	ctx := context.Background()
	out, err := RunCombined(ctx, dir, "status", "--porcelain")
	if err != nil {
		return nil, err
	}
	return ParsePorcelain(out), nil
}

// ParsePorcelain parses `git status --porcelain` output (porcelain v1).
// It returns a slice of FileStatus preserving Index/Worktree status and path(s).
//
// Format (per-line):
// XY <path> (for regular)
// XY <from> -> <to>  (for rename/copy)
// We split the line into the 2-status chars then the rest after a single space.
func ParsePorcelain(out string) []FileStatus {
	var res []FileStatus
	if strings.TrimSpace(out) == "" {
		return res
	}

	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		fs := FileStatus{RawLine: ln}
		// porcelain v1: first two chars are status, then a space, then path (possibly "from -> to")
		if len(ln) < 3 {
			// malformed, still keep raw
			res = append(res, fs)
			continue
		}
		fs.IndexStatus = ln[0]
		fs.WorktreeStatus = ln[1]
		rest := strings.TrimSpace(ln[3:]) // skip "XY " (2 chars + single space)
		fs.Path = rest
		// check rename pattern "from -> to"
		if idx := strings.Index(rest, "->"); idx != -1 {
			// split on '->', trim spaces around
			parts := strings.SplitN(rest, "->", 2)
			orig := strings.TrimSpace(parts[0])
			newp := strings.TrimSpace(parts[1])
			fs.OrigPath = orig
			fs.Path = newp
		}
		res = append(res, fs)
	}
	return res
}

// RunStream runs `git <args...>` in dir and streams stdout/stderr to callbacks.
// It returns the exit error when the command finishes.
func RunStream(ctx context.Context, dir string, args []string,
	onStdout func(string),
	onStderr func(string)) error {

	cmd := exec.CommandContext(ctx, "git", args...)
	if dir != "" {
		cmd.Dir = dir
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start git: %w", err)
	}

	// Stream stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if onStdout != nil {
				onStdout(line)
			}
		}
	}()

	// Stream stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if onStderr != nil {
				onStderr(line)
			}
		}
	}()

	// Wait for process to exit
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("git failed: %w", err)
	}
	return nil
}

// Fetch runs `git fetch origin main`
func Fetch(dir string) error {
	ctx := context.Background()
	_, err := RunCombined(ctx, dir, "fetch", "origin", "main")
	if err != nil {
		return fmt.Errorf("git fetch failed: %w", err)
	}
	return nil
}

// RebaseOntoMain Rebase runs `git rebase origin/main`
func RebaseOntoMain(dir string) error {
	ctx := context.Background()
	_, err := RunCombined(ctx, dir, "rebase", "origin/main")
	if err != nil {
		return fmt.Errorf("git rebase failed: %w", err)
	}
	return nil
}

// RunGitWithOutput runs `git <args...>` in dir and streams stdout/stderr lines
// back to the caller through channels. The caller must read both channels until
// they are closed. If the command exits with error, it is sent on the error channel.
func RunGitWithOutput(ctx context.Context, args ...string) (<-chan string, <-chan error) {
	outCh := make(chan string)
	errCh := make(chan error, 1) // buffered so goroutine can exit

	go func() {
		defer close(outCh)
		defer close(errCh)

		// apply default timeout if none given
		if _, ok := ctx.Deadline(); !ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, 8*time.Minute)
			defer cancel()
		}

		cmd := exec.CommandContext(ctx, "git", args...)

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			errCh <- fmt.Errorf("stdout pipe: %w", err)
			return
		}
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			errCh <- fmt.Errorf("stderr pipe: %w", err)
			return
		}

		if err := cmd.Start(); err != nil {
			errCh <- fmt.Errorf("start git: %w", err)
			return
		}

		// scan stdout
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				select {
				case outCh <- scanner.Text():
				case <-ctx.Done():
					return
				}
			}
		}()

		// scan stderr
		go func() {
			scanner := bufio.NewScanner(stderrPipe)
			for scanner.Scan() {
				select {
				case outCh <- "[stderr] " + scanner.Text():
				case <-ctx.Done():
					return
				}
			}
		}()

		// wait for command to finish
		if err := cmd.Wait(); err != nil {
			errCh <- fmt.Errorf("git failed: %w", err)
			return
		}
		errCh <- nil
	}()

	return outCh, errCh
}

// IsDirty checks if there are any uncommitted changes in the repo.
// Returns true if there are staged or unstaged changes.
func IsDirty(dir string) (bool, error) {
	ctx := context.Background()
	out, err := RunCombined(ctx, dir, "status", "--porcelain")
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}
