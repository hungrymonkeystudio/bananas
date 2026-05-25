// settings screen for bananatype
// core functionality: change time control

package settings

import (
	"github.com/WarrenWu4/bananatype/internal/foundation/theme"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

const settingInstructions = "UP/DOWN/LEFT/RIGHT to move\nENTER to select\nESC to close settings page"

type SettingsModel struct {
	Show            bool
	options         []string
	optionIdx       int
	times           []int
	timeIdx         int
	words           []int
	wordIdx         int
	ActiveTime      int
	ActiveWords     int
	ActiveTyperMode string
}

func NewSettingsModel() SettingsModel {
	s := SettingsModel{
		Show:            false,
		options:         []string{"timer", "words", "restart", "quit"},
		optionIdx:       0,
		times:           []int{15, 30, 60, 120},
		timeIdx:         1,
		words:           []int{10, 25, 50, 100},
		wordIdx:         2,
		ActiveTime:      30,
		ActiveWords:     50,
		ActiveTyperMode: "timer",
	}
	readSettings(&s)
	return s
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.optionIdx = (m.optionIdx - 1 + len(m.options)) % len(m.options)
		case "down":
			m.optionIdx = (m.optionIdx + 1) % len(m.options)
		case "left":
			switch m.options[m.optionIdx] {
			case "timer":
				m.timeIdx = (m.timeIdx - 1 + len(m.times)) % len(m.times)
			case "words":
				m.wordIdx = (m.wordIdx - 1 + len(m.words)) % len(m.words)
			}
		case "right":
			switch m.options[m.optionIdx] {
			case "timer":
				m.timeIdx = (m.timeIdx + 1) % len(m.times)
			case "words":
				m.wordIdx = (m.wordIdx + 1) % len(m.words)
			}
		case "enter":
			switch m.options[m.optionIdx] {
			case "timer":
				m.ActiveTime = m.times[m.timeIdx]
				m.ActiveTyperMode = "timer"
				m.writeSettings()
				return m, func() tea.Msg { return m }
			case "words":
				m.ActiveWords = m.words[m.wordIdx]
				m.ActiveTyperMode = "words"
				m.writeSettings()
				return m, func() tea.Msg { return m }
			case "quit":
				return m, tea.Quit
			case "restart":
				return m, func() tea.Msg { return m }
			}
		}
	}
	return m, nil
}

func (m SettingsModel) View() string {
	output := theme.Yellow.Render("Settings") + "\n"
	if m.options[m.optionIdx] == "timer" {
		output += theme.White.Render("timer: ")
		for timeIdx, times := range m.times {
			if timeIdx == m.timeIdx {
				if m.times[m.timeIdx] == m.ActiveTime && m.ActiveTyperMode == "timer" {
					output += theme.White.Underline(true).Render(strconv.Itoa(times)) + " "
				} else {
					output += theme.Gray.Underline(true).Render(strconv.Itoa(times)) + " "
				}
			} else {
				if m.times[timeIdx] == m.ActiveTime && m.ActiveTyperMode == "timer" {
					output += theme.White.Render(strconv.Itoa(times)) + " "
				} else {
					output += theme.Gray.Render(strconv.Itoa(times)) + " "
				}
			}
		}
	} else {
		output += theme.Gray.Render("timer: ")
		for timeIdx, times := range m.times {
			if m.times[timeIdx] == m.ActiveTime && m.ActiveTyperMode == "timer" {
				output += theme.White.Render(strconv.Itoa(times)) + " "
			} else {
				output += theme.Gray.Render(strconv.Itoa(times)) + " "
			}
		}
	}
	output += "\n"
	if m.options[m.optionIdx] == "words" {
		output += theme.White.Render("words: ")
		for wordIdx, word := range m.words {
			if wordIdx == m.wordIdx {
				if m.words[m.wordIdx] == m.ActiveWords && m.ActiveTyperMode == "words" {
					output += theme.White.Underline(true).Render(strconv.Itoa(word)) + " "
				} else {
					output += theme.Gray.Underline(true).Render(strconv.Itoa(word)) + " "
				}
			} else {
				if m.words[wordIdx] == m.ActiveWords && m.ActiveTyperMode == "words" {
					output += theme.White.Render(strconv.Itoa(word)) + " "
				} else {
					output += theme.Gray.Render(strconv.Itoa(word)) + " "
				}
			}
		}
	} else {
		output += theme.Gray.Render("words: ")
		for wordIdx, word := range m.words {
			if m.words[wordIdx] == m.ActiveWords && m.ActiveTyperMode == "words" {
				output += theme.White.Render(strconv.Itoa(word)) + " "
			} else {
				output += theme.Gray.Render(strconv.Itoa(word)) + " "
			}
		}
	}
	output += "\n"
	if m.options[m.optionIdx] == "restart" {
		output += theme.White.Render("restart")
	} else {
		output += theme.Gray.Render("restart")
	}
	output += "\n"
	if m.options[m.optionIdx] == "quit" {
		output += theme.White.Render("quit")
	} else {
		output += theme.Gray.Render("quit")
	}
	output += "\n"
	output += "\n" + theme.Instructions.Render(settingInstructions)
	return output
}
