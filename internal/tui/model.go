/*
Copyright © 2025 Bernard Katamanso <bernard@orctatech.com>
*/
package tui

import (
	"github.com/Orctatech-Engineering-Team/GitMate/internal/git"
	tea "github.com/charmbracelet/bubbletea"
)

// Command --- Command Enum
type Command int

const (
	CmdTutor Command = iota
	CmdStart
	CmdSync
	CmdClean
)

// Model --- Model
type Model struct {
	current Command
	status  string
	files   []git.FileStatus
	err     error
}

// --- Constructor
func NewModel(initial Command) Model {
	var status string
	switch initial {
	case CmdTutor:
		status = "Tutor mode ready (future: step-by-step guide)"
	case CmdStart:
		status = "Starting a new feature branch..."
	case CmdSync:
		status = "Syncing your branch with main..."
	case CmdClean:
		status = "Preparing to clean up commits..."
	default:
		status = "Welcome to GitMate"
	}
	return Model{current: initial, status: status}
}

// --- Messages
type tutorDoneMsg string
type startDoneMsg string
type syncDoneMsg []git.FileStatus
type cleanDoneMsg string
type errMsg struct{ err error }

// --- Init
func (m Model) Init() tea.Cmd {
	switch m.current {
	case CmdTutor:
		return tutorCmd
	case CmdStart:
		return startCmd
	case CmdSync:
		return syncCmd
	case CmdClean:
		return cleanCmd
	default:
		return nil
	}
}

// --- Command Functions
func tutorCmd() tea.Msg {
	// stub: later this could walk user through steps
	return tutorDoneMsg("Tutorial started")
}

func startCmd() tea.Msg {
	// stub: later this could run "git checkout -b feature/<name>"
	return startDoneMsg("Created new feature branch (stub)")
}

func syncCmd() tea.Msg {
	dir := "."
	if err := git.Fetch(dir); err != nil {
		return errMsg{err}
	}
	if err := git.RebaseOntoMain(dir); err != nil {
		return errMsg{err}
	}
	files, err := git.GitStatusPorcelain(dir)
	if err != nil {
		return errMsg{err}
	}
	return syncDoneMsg(files)
}

func cleanCmd() tea.Msg {
	// stub: later this could run autosquash or interactive rebase hints
	return cleanDoneMsg("Clean complete (stub)")
}

// --- Update
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// keyboard navigation
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "t":
			m.current = CmdTutor
			m.status = "Tutor mode selected"
			return m, tutorCmd
		case "s":
			m.current = CmdStart
			m.status = "Start command triggered"
			return m, startCmd
		case "y":
			m.current = CmdSync
			m.status = "Syncing..."
			return m, syncCmd
		case "c":
			m.current = CmdClean
			m.status = "Cleaning repo..."
			return m, cleanCmd
		}

	// command results
	case tutorDoneMsg:
		m.status = "Tutor finished: " + string(msg)
	case startDoneMsg:
		m.status = "Branch started: " + string(msg)
	case syncDoneMsg:
		m.files = msg
		m.status = "Sync complete"
	case cleanDoneMsg:
		m.status = "Clean finished: " + string(msg)
	case errMsg:
		m.err = msg.err
	}

	return m, nil
}

// --- View
func (m Model) View() string {
	s := "GitMate UI\n\nStatus: " + m.status + "\n\n"

	if m.err != nil {
		s += "Error: " + m.err.Error() + "\n"
		return s + "(press q to quit)"
	}

	switch m.current {
	case CmdSync:
		if len(m.files) == 0 {
			s += "Working tree clean.\n"
		} else {
			s += "Changes:\n"
			for _, f := range m.files {
				s += string(f.IndexStatus) + string(f.WorktreeStatus) + " " + f.Path + "\n"
			}
		}
	case CmdTutor:
		s += "(Tutor mode: more to come)\n"
	case CmdStart:
		s += "(Feature branch creation goes here)\n"
	case CmdClean:
		s += "(Clean-up logic goes here)\n"
	}

	return s + "\n(press q to quit)"
}
