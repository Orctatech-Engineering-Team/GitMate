/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package tui

import (
	"fmt"
	"github.com/Orctatech-Engineering-Team/GitMate/internal/git"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------- Tutor Step ----------------

type tutorStep struct {
	Title       string
	Description string
	Command     string
	Action      func(*tea.Program) // optional live execution
	Completed   bool
}

// ---------------- Tutor Model ----------------

type tutorModel struct {
	spinner   spinner.Model
	steps     []tutorStep
	current   int
	logs      []string
	done      bool
	repoReady bool
}

// ---------------- Constructor ----------------

func NewTutorModel() tutorModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	isGitRepo, _ := git.IsRepo(".")

	steps := []tutorStep{
		{
			Title:       "Start Command",
			Description: "Learn how to start a new feature branch with `gitmate start <feature>`.",
			Command:     "gitmate start login-api",
			Action: func(p *tea.Program) {
				if isGitRepo {
					_ = RunStartTUI("login-api")
				} else {
					p.Send(tutorMsg("Repository not initialized. Cannot run start command."))
				}
			},
		},
		{
			Title:       "Sync Command",
			Description: "Keep your branch up-to-date with main using `gitmate sync`.",
			Command:     "gitmate sync",
			Action: func(p *tea.Program) {
				if isGitRepo {
					_ = RunSyncTUI()
				} else {
					p.Send(tutorMsg("Repository not initialized. Cannot run sync command."))
				}
			},
		},
	}

	repoReady := isGitRepo
	return tutorModel{
		spinner:   s,
		steps:     steps,
		current:   0,
		logs:      []string{},
		repoReady: repoReady,
	}
}

// ---------------- Init ----------------

func (m tutorModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// ---------------- Update ----------------

func (m tutorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n": // next
			if m.current < len(m.steps)-1 {
				m.current++
			}
		case "p": // previous
			if m.current > 0 {
				m.current--
			}
		case "r": // run current step
			step := m.steps[m.current]
			if step.Action != nil {
				go step.Action(nil) // run in goroutine
				m.logs = append(m.logs, fmt.Sprintf("Running step: %s", step.Title))
			} else {
				m.logs = append(m.logs, "Step cannot be run directly.")
			}
		case "c": // mark completed
			m.steps[m.current].Completed = true
			m.logs = append(m.logs, fmt.Sprintf("Step marked completed: %s", m.steps[m.current].Title))
		}

	case spinner.TickMsg:
		if !m.done {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tutorMsg:
		m.logs = append(m.logs, string(msg))
	case tutorErrMsg:
		m.logs = append(m.logs, "Error: "+msg.Error())
	case tutorDoneMsg:
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

// ---------------- View ----------------

func (m tutorModel) View() string {
	step := m.steps[m.current]
	s := fmt.Sprintf("Tutor: Step %d/%d\n\n", m.current+1, len(m.steps))
	s += fmt.Sprintf("Title: %s\n", step.Title)
	s += fmt.Sprintf("Description:\n%s\n\n", step.Description)
	if step.Command != "" {
		s += fmt.Sprintf("Command: %s\n", step.Command)
	}
	if step.Completed {
		s += "✅ Completed\n\n"
	}

	if !m.done {
		s += m.spinner.View() + " Navigate with (n)ext/(p)revious, (r)un step, (c)omplete, (q)uit\n\n"
	}

	start := 0
	if len(m.logs) > 10 {
		start = len(m.logs) - 10
	}
	for _, log := range m.logs[start:] {
		s += log + "\n"
	}

	if m.done {
		s += "\n(press q to quit)"
	}

	return s
}

// ---------------- Public Entry ----------------

func RunTutorTUI() error {
	p := tea.NewProgram(NewTutorModel())
	_, err := p.Run()
	return err
}
