// session is the root Bubble Tea model that owns the application state machine.
// It orchestrates screen transitions and tracks test progress (timer/word count).

package session

import (
	"fmt"
	"strings"
	"time"

	"github.com/WarrenWu4/bananatype/internal/foundation/logger"
	"github.com/WarrenWu4/bananatype/internal/foundation/theme"
	"github.com/WarrenWu4/bananatype/internal/ui/results"
	"github.com/WarrenWu4/bananatype/internal/ui/settings"
	"github.com/WarrenWu4/bananatype/internal/ui/typing"
	bubbleTimer "github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

// State represents which screen is active
type State int

const (
	StateTyping  State = iota
	StateSettings
	StateResults
)

type Model struct {
	state    State
	settings settings.SettingsModel
	typing   typing.TypingModel
	results  results.ResultsModel

	// session progress tracking
	timer     bubbleTimer.Model
	timerDone bool
	started   bool
	startTime time.Time
	doneTime  time.Time
	done      bool

	// viewport
	width  int
	height int
}

func New() Model {
	s := settings.NewSettingsModel()
	ty := typing.NewTypingModel()
	return Model{
		state:    StateTyping,
		settings: s,
		typing:   ty,
		results:  results.NewResultsModel(),
		timer:    bubbleTimer.NewWithInterval(time.Second*time.Duration(s.ActiveTime), time.Second),
		width:    120,
		height:   8,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global updates
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.state == StateSettings {
				m.state = StateTyping
			} else if m.state == StateTyping {
				m.state = StateSettings
			}
			return m, nil
		}
	case settings.SettingsModel:
		// Settings confirmed — reset session
		m.settings = msg
		m.typing = typing.NewTypingModel()
		m.resetProgress()
		m.state = StateTyping
		return m, nil
	case results.ResultsModel:
		// Retry from results screen
		m.typing = typing.NewTypingModel()
		m.resetProgress()
		m.state = StateTyping
		return m, nil
	}

	// State-specific updates
	switch m.state {
	case StateSettings:
		updatedSettings, cmd := m.settings.Update(msg)
		m.settings = updatedSettings.(settings.SettingsModel)
		return m, cmd

	case StateResults:
		logger.Log(logger.DEBUG, "Updating results with message")
		updatedResults, cmd := m.results.Update(msg)
		m.results = updatedResults.(results.ResultsModel)
		return m, cmd

	case StateTyping:
		var cmds []tea.Cmd

		// Start timer on first keypress
		if !m.started {
			if _, ok := msg.(tea.KeyMsg); ok {
				m.started = true
				m.startTime = time.Now()
				cmds = append(cmds, m.timer.Init())
			}
		}

		// Update timer
		switch msg := msg.(type) {
		case bubbleTimer.TickMsg:
			var cmd tea.Cmd
			m.timer, cmd = m.timer.Update(msg)
			cmds = append(cmds, cmd)
		case bubbleTimer.StartStopMsg:
			var cmd tea.Cmd
			m.timer, cmd = m.timer.Update(msg)
			cmds = append(cmds, cmd)
		case bubbleTimer.TimeoutMsg:
			m.timerDone = true
		}

		// Update typing
		updatedTyping, typerCmd := m.typing.Update(msg)
		m.typing = updatedTyping.(typing.TypingModel)
		if typerCmd != nil {
			cmds = append(cmds, typerCmd)
		}

		// Check completion
		m.checkDone()
		if m.done {
			m.transitionToResults()
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m Model) View() string {
	paddingY := (m.height - typing.MAXLINES + 1) / 2
	paddingX := (m.width - typing.MAXCHARPERLINE) / 2

	output := strings.Repeat("\n", paddingY)

	switch m.state {
	case StateResults:
		outputLines := strings.Split(m.results.View(), "\n")
		for _, line := range outputLines {
			output += strings.Repeat(" ", paddingX) + line + "\n"
		}
	case StateSettings:
		outputLines := strings.Split(m.settings.View(), "\n")
		for _, line := range outputLines {
			output += strings.Repeat(" ", paddingX) + line + "\n"
		}
	case StateTyping:
		output += strings.Repeat(" ", paddingX) + m.progressView() + "\n"
		outputLines := strings.Split(m.typing.View(), "\n")
		for _, line := range outputLines {
			output += strings.Repeat(" ", paddingX) + line + "\n"
		}
	}

	return output
}

func (m *Model) checkDone() {
	switch m.settings.ActiveTyperMode {
	case "timer":
		if m.timerDone {
			m.done = true
			m.doneTime = time.Now()
		}
	case "words":
		if m.typing.TotalWords >= m.settings.ActiveWords {
			m.done = true
			m.doneTime = time.Now()
		}
	}
}

func (m *Model) transitionToResults() {
	switch m.settings.ActiveTyperMode {
	case "timer":
		m.results.Time = float64(m.settings.ActiveTime)
	case "words":
		m.results.Time = m.doneTime.Sub(m.startTime).Seconds()
	}
	m.results.Words = m.typing.TotalWords
	m.results.Correct = m.typing.TotalCorrect
	m.results.Characters = m.typing.TotalTyped
	m.state = StateResults
}

func (m *Model) resetProgress() {
	m.timer = bubbleTimer.NewWithInterval(time.Second*time.Duration(m.settings.ActiveTime), time.Second)
	m.timerDone = false
	m.started = false
	m.startTime = time.Time{}
	m.doneTime = time.Time{}
	m.done = false
}

func (m Model) progressView() string {
	switch m.settings.ActiveTyperMode {
	case "timer":
		seconds := int(m.timer.Timeout.Seconds())
		return theme.Yellow.Render(fmt.Sprintf("%ds", seconds))
	case "words":
		return theme.Yellow.Render(fmt.Sprintf("%d/%d", m.typing.TotalWords, m.settings.ActiveWords))
	default:
		return ""
	}
}

func (m Model) GetTyping() typing.TypingModel {
	return m.typing
}
