/*
Copyright © 2025 Bernard Katamanso
*/
package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type SyncModel struct {
	spinner spinner.Model
	logs    []string
	err     error
	done    bool
}

// NewSyncModel Creates a new syncModel
func NewSyncModel() SyncModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return SyncModel{
		spinner: s,
		logs:    []string{},
	}
}

func (m SyncModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SyncModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case spinner.TickMsg:
		if !m.done {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case gitLineMsg:
		m.logs = append(m.logs, string(msg))
		return m, nil

	case gitErrMsg:
		m.err = msg
		m.done = true
		return m, nil

	case gitDoneMsg:
		m.done = true
		return m, nil
	}

	return m, nil
}

// --- View
func (m SyncModel) View() string {
	s := "GitMate: Syncing with main\n\n"
	if m.err != nil {
		s += "Error: " + m.err.Error() + "\n"
	} else if m.done {
		s += "✅ Sync complete.\n\n"
	} else {
		s += m.spinner.View() + " Running git fetch & rebase...\n\n"
	}

	// show last 10 log lines
	start := 0
	if len(m.logs) > 10 {
		start = len(m.logs) - 10
	}
	for _, line := range m.logs[start:] {
		s += line + "\n"
	}

	if m.done {
		s += "\n(press q to quit)"
	}
	return s
}

// --- Orchestration of sync steps
func runSync(p *tea.Program) {
	// Step 1: git fetch --all
	streamStep(p, "fetch", []string{"--all"}, func() {
		// Step 2: git rebase origin/main
		streamStep(p, "rebase", []string{"origin/main"}, func() {
			// Step 3: done
			p.Send(gitDoneMsg{})
		})
	})
}

func RunSyncTUI() error {
	p := tea.NewProgram(NewSyncModel())
	// start orchestration after the program begins
	go runSync(p)
	_, err := p.Run()
	return err
}
