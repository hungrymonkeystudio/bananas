// setting screen for bananas
// core functionality: change time control
package main

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

const settingInstructions = "UP/DOWN/LEFT/RIGHT to move\nENTER to select"

type SettingsModel struct {
    show bool;
    options []string;
    optionIdx int;
    times []int;
    timeIdx int;
    activeTime int;
}

func NewSettingsModel() SettingsModel {
    return SettingsModel{
        show: false,
        options: []string{"timer", "restart", "quit"},
        optionIdx: 0,
        times: []int{15, 30, 60, 120},
        timeIdx: 1,
        activeTime: 30,
    }
}

func (m SettingsModel) Init() tea.Cmd {
    return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up":
            m.optionIdx = (m.optionIdx-1+len(m.options)) % len(m.options)
        case "down":
            m.optionIdx = (m.optionIdx+1) % len(m.options)
        case "left":
            if m.options[m.optionIdx] == "timer" {
                m.timeIdx = (m.timeIdx-1+len(m.times)) % len(m.times)
            }
        case "right":
            if m.options[m.optionIdx] == "timer" {
                m.timeIdx = (m.timeIdx+1) % len(m.times)
            }
        case "enter":
            if m.options[m.optionIdx] == "quit" {
                return m, tea.Quit
            } else if m.options[m.optionIdx] == "restart" {
                return m, nil
            } else if m.options[m.optionIdx] == "timer" {
                m.activeTime = m.times[m.timeIdx]
            }
        }
    }
    return m, nil
}

func (m SettingsModel) View() string {
    output := yellow.Render("Settings Page") + "\n"
    if m.options[m.optionIdx] == "timer" {
        output += white.Render("timer: ")
        for timeIdx, times := range m.times {
            if timeIdx == m.timeIdx {
                if m.times[m.timeIdx] == m.activeTime {
                    output += white.Underline(true).Render(strconv.Itoa(times)) + " "
                } else {
                    output += gray.Underline(true).Render(strconv.Itoa(times)) + " "
                }
            } else {
                if m.times[timeIdx] == m.activeTime {
                    output += white.Render(strconv.Itoa(times)) + " "
                } else {
                    output += gray.Render(strconv.Itoa(times)) + " "
                }
            }
        }
    } else {
        output += gray.Render("timer: ")
        for timeIdx, times := range m.times {
            if timeIdx == m.timeIdx {
                if m.times[m.timeIdx] == m.activeTime {
                    output += gray.Underline(true).Render(strconv.Itoa(times)) + " "
                } else {
                    output += gray.Render(strconv.Itoa(times)) + " "
                }
            } else {
                output += gray.Render(strconv.Itoa(times)) + " "
            }
        }
    }
    output += "\n"
    if m.options[m.optionIdx] == "restart" {
        output += white.Render("restart")
    } else {
        output += gray.Render("restart")
    }
    output += "\n"
    if m.options[m.optionIdx] == "quit" {
        output += white.Render("quit")
    } else {
        output += gray.Render("quit")
    }
    output += "\n"
    output += "\n" + instructions.Render(settingInstructions)
    return output
}
