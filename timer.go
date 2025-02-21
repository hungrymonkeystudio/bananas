package main

import (
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/timer"
)

// custom wrapper around bubbbles timer

type TimerModel struct {
    timer timer.Model
    started bool
    done bool
}

func NewTimerModel(startTime time.Duration) TimerModel {
    return TimerModel{
        timer: timer.NewWithInterval(startTime, time.Second),
        started: false,
        done: false,
    }
}

func (m TimerModel) Init() tea.Cmd {
    return nil
}

func (m TimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case timer.TickMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd
    case timer.StartStopMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd
    case timer.TimeoutMsg:
        m.done = true
        return m, nil
    case tea.KeyMsg:
        switch msg.String() {
        default:
            if !m.started {
                m.started = true
                return m, m.timer.Init()
            }
        }
    }
    return m, nil
}

func (m TimerModel) View() string {
    return yellow.Render(m.timer.View())
}
