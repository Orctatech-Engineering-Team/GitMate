/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package tui

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"regexp"
	"strings"

	"github.com/Orctatech-Engineering-Team/GitMate/internal/git"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------- Prompt Model ----------------

type promptChoice string

const (
	choiceStash   promptChoice = "Stash changes"
	choiceCommit  promptChoice = "Commit all changes"
	choiceDiscard promptChoice = "Discard changes"
	choiceQuit    promptChoice = "Quit"
)

type promptModel struct {
	list   list.Model
	done   bool
	choice promptChoice
}

func NewPromptModel() promptModel {
	items := []list.Item{
		listItem(choiceStash),
		listItem(choiceCommit),
		listItem(choiceDiscard),
		listItem(choiceQuit),
	}
	l := list.New(items, list.NewDefaultDelegate(), 30, 10)
	l.Title = "Uncommitted changes detected. What do you want to do?"
	return promptModel{list: l}
}

func (m promptModel) Init() tea.Cmd { return nil }

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(listItem)
			if ok {
				m.choice = promptChoice(i)
				m.done = true
				return m, tea.Quit
			}
		case "q", "ctrl+c":
			m.choice = choiceQuit
			m.done = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m promptModel) View() string {
	if m.done {
		return fmt.Sprintf("Selected: %s\n", m.choice)
	}
	return m.list.View()
}

type listItem string

func (i listItem) FilterValue() string { return string(i) }

// ---------------- Branch Name Input Model ----------------

type branchInputModel struct {
	textInput textinput.Model
	done      bool
	branch    string
}

func NewBranchInputModel() branchInputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter feature branch name..."
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 30
	return branchInputModel{textInput: ti}
}

func (m branchInputModel) Init() tea.Cmd { return textinput.Blink }

func (m branchInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			m.branch = m.textInput.Value()
			return m, tea.Quit
		case "ctrl+c":
			m.done = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m branchInputModel) View() string {
	return fmt.Sprintf(
		"Feature branch name:\n%s\n\n(press enter to continue, ctrl+c to quit)",
		m.textInput.View(),
	)
}

// ---------------- Start Model ----------------

type startModel struct {
	spinner spinner.Model
	logs    []string
	err     error
	done    bool
	branch  string
}

func NewStartModel(branch string) startModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return startModel{
		spinner: s,
		branch:  branch,
	}
}

func (m startModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m startModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m startModel) View() string {
	s := fmt.Sprintf("GitMate: Starting new feature branch '%s'\n\n", m.branch)
	if m.err != nil {
		s += "Error: " + m.err.Error() + "\n"
	} else if m.done {
		s += fmt.Sprintf("✅ Branch feature/%s created and checked out.\n\n", m.branch)
	} else {
		s += m.spinner.View() + " Running git commands...\n\n"
	}
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

// ---------------- Helpers ----------------

func sanitizeBranchName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	re := regexp.MustCompile(`[^a-z0-9._-]+`)
	return re.ReplaceAllString(name, "-")
}

// ---------------- Orchestration ----------------

// runStart orchestrates checkout main → pull → create feature branch with live logs
func runStart(p *tea.Program, branch string) {
	streamStep(p, "checkout", []string{"main"}, func() {
		streamStep(p, "pull", []string{"origin", "main"}, func() {
			safe := sanitizeBranchName(branch)
			streamStep(p, "checkout", []string{"-b", "feature/" + safe}, func() {
				p.Send(gitDoneMsg{})
			})
		})
	})
}

// ---------------- Public Entry ----------------

func RunStartTUI(featureName string) error {
	// 1. Prompt for branch name if none provided
	if featureName == "" {
		tiProgram := tea.NewProgram(NewBranchInputModel())
		final, err := tiProgram.Run()
		if err != nil {
			return err
		}
		if m, ok := final.(branchInputModel); ok {
			featureName = m.branch
		}
		if featureName == "" {
			return fmt.Errorf("no branch name provided")
		}
	}

	// 2. Check if repo is dirty
	dirty, err := git.IsDirty(".")
	if err != nil {
		return err
	}
	if dirty {
		// Run prompt for stash/commit/discard
		pm := tea.NewProgram(NewPromptModel())
		go runStart(pm, featureName)
		final, err := pm.Run()
		if err != nil {
			return err
		}
		if p, ok := final.(promptModel); ok {
			switch p.choice {
			case choiceStash:
				_, err = git.RunCombined(context.Background(), ".", "stash", "push", "-u")
			case choiceCommit:
				_, err = git.RunCombined(context.Background(), ".", "add", "-A")
				if err == nil {
					// Open Git editor for commit message
					_, err = git.RunCombined(context.Background(), ".", "commit")
				}
			case choiceDiscard:
				_, err = git.RunCombined(context.Background(), ".", "reset", "--hard")
			case choiceQuit:
				return nil
			}
			if err != nil {
				return err
			}
		}
	}

	// 3. Run main start model with live logs
	p := tea.NewProgram(NewStartModel(featureName))
	runStart(p, featureName)
	_, err = p.Run()
	return err
}
