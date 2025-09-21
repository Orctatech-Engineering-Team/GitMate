/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package tui

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Orctatech-Engineering-Team/GitMate/internal/git"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------- Prompt Model ----------------

type confirmModel struct {
	list list.Model
	done bool
	yes  bool
}

func NewConfirmModel() confirmModel {
	items := []list.Item{
		listItem("Yes, run autosquash rebase"),
		listItem("No, cancel"),
	}
	l := list.New(items, list.NewDefaultDelegate(), 40, 6)
	l.Title = "Noisy commits detected. Run 'git rebase -i --autosquash'?"
	return confirmModel{list: l}
}

func (m confirmModel) Init() tea.Cmd { return nil }

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				if i == "Yes, run autosquash rebase" {
					m.yes = true
				}
				m.done = true
				return m, tea.Quit
			}
		case "q", "ctrl+c":
			m.done = true
			m.yes = false
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m confirmModel) View() string {
	if m.done {
		return ""
	}
	return m.list.View()
}

// ---------------- Clean Model ----------------

type cleanModel struct {
	spinner spinner.Model
	logs    []string
	err     error
	done    bool
	noisy   []string
}

func NewCleanModel(noisyCommits []string) cleanModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return cleanModel{
		spinner: s,
		noisy:   noisyCommits,
		logs:    []string{},
	}
}

func (m cleanModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m cleanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case gitErrMsg:
		m.err = msg
		m.done = true
		return m, tea.Quit
	case gitDoneMsg:
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m cleanModel) View() string {
	s := "GitMate: Cleaning noisy commits\n\n"
	if m.err != nil {
		s += "Error: " + m.err.Error() + "\n"
	} else if m.done {
		s += "✅ Clean operation complete.\n\n"
	} else {
		s += m.spinner.View() + " Preparing...\n\n"
	}
	if len(m.noisy) > 0 {
		s += "Noisy commits:\n"
		for _, c := range m.noisy {
			s += fmt.Sprintf(" - %s\n", c)
		}
	}
	// show last 10 logs
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

// ---------------- Public Entry ----------------

func RunCleanTUI() error {
	// 1. Get recent commits
	out, err := git.RunCombined(context.Background(), ".", "log", "--oneline", "-n", "20")
	if err != nil {
		return err
	}
	lines := strings.Split(out, "\n")
	noisy := []string{}
	re := regexp.MustCompile(`\bfix(e[sd])?\b|\btypo\b|\bdebug\b|\boops\b`)
	for _, l := range lines {
		parts := strings.SplitN(l, " ", 2)
		if len(parts) < 2 {
			continue
		}
		msg := parts[1]
		if re.MatchString(strings.ToLower(msg)) {
			noisy = append(noisy, l)
		}
	}

	if len(noisy) == 0 {
		fmt.Println("No noisy commits detected. Nothing to clean.")
		return nil
	}

	// 2. Prompt user for confirmation
	pm := tea.NewProgram(NewConfirmModel())
	final, err := pm.Run()
	if err != nil {
		return err
	}
	if p, ok := final.(confirmModel); ok && p.yes {
		// 3. Run interactive autosquash rebase with live logs
		cm := tea.NewProgram(NewCleanModel(noisy))
		go runClean(cm)
		_, err := cm.Run()
		return err
	}

	return nil
}

// ---------------- Orchestration ----------------

func runClean(p *tea.Program) {
	streamStep(p, "rebase", []string{"-i", "--autosquash", "HEAD~20"}, func() {
		p.Send(gitDoneMsg{})
	})
}
