package tui

import (
	"context"
	"github.com/Orctatech-Engineering-Team/GitMate/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

// --- messages
type gitLineMsg string
type gitErrMsg error
type gitDoneMsg struct{}

// streamStep runs a git command, streams logs/errors into Update, then calls next if success
func streamStep(p *tea.Program, cmd string, args []string, next func()) {
	ctx := context.Background()
	outCh, errCh := git.RunGitWithOutput(ctx, append([]string{cmd}, args...)...)

	// forward stdout
	go func() {
		for line := range outCh {
			p.Send(gitLineMsg(line))
		}
	}()

	// forward stderr/errors
	go func() {
		for err := range errCh {
			if err != nil {
				p.Send(gitErrMsg(err))
				return
			}
		}
		if next != nil {
			next()
		}
	}()
}
