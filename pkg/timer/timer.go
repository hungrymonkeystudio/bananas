package timer 

import (
	"strconv"
	"time"
	colors "bananas/pkg/colors"
	bubbleTimer "github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

// custom wrapper around bubbbles timer

type TimerModel struct {
    timer bubbleTimer.Model
    started bool
    Done bool
}

func NewTimerModel(startTime time.Duration) TimerModel {
    return TimerModel{
        timer: bubbleTimer.NewWithInterval(startTime, time.Second),
        started: false,
        Done: false,
    }
}

func (m TimerModel) Init() tea.Cmd {
    return nil
}

func (m TimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case bubbleTimer.TickMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd
    case bubbleTimer.StartStopMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd
    case bubbleTimer.TimeoutMsg:
        m.Done = true
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
    secondsDuration := strconv.Itoa(int(m.timer.Timeout.Seconds()))
    return colors.Yellow.Render(secondsDuration + "s")
}
