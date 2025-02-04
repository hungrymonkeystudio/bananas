package main

import (
    "fmt"
    "time"

    tea "github.com/charmbracelet/bubbletea"
)

type TimerModel struct {
    start bool
    countdown int
}

func (tm TimerModel) Init() tea.Cmd {
    return nil
}

func (tm TimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return tm, tea.Quit
        case "q":
            return tm, tea.Quit
        default:
            if tm.start {
                return tm, nil
            } else {
                tm.start = true 
            }
        }
    }

    if tm.start && tm.countdown > 0 {
        tm.countdown--
    }

    if tm.start {
        return tm, tea.Every(time.Second, func(t time.Time) tea.Msg {
            return tea.Tick(time.Second, func(t time.Time) tea.Msg {
                return nil
            }) 
        })
    }
    return tm, nil
}

func (tm TimerModel) View() string {
    timerText := yellow.Render(fmt.Sprintf("Time: %d seconds", tm.countdown))
    return timerText
}
