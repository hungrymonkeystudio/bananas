package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

const timeout = time.Second * 30

type MainModel struct {
    timer timer.Model 
    typer TyperModel
    analysis AnalysisModel
}

func (m MainModel) Init() tea.Cmd {
    return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case timer.TickMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd
    case timer.StartStopMsg:
        var cmd tea.Cmd
        m.timer, cmd = m.timer.Update(msg)
        return m, cmd
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "q":
            return m, tea.Quit
        case "enter":
            m.typer.words = createTest()
            m.typer.idx = 0
        default:
            return m, m.timer.Init()
        }
    }
    updatedTyper, typerCmd := m.typer.Update(msg)
    m.typer = updatedTyper.(TyperModel)
    updatedAnalysis, analysisCmd := m.analysis.Update(msg)
    m.analysis = updatedAnalysis.(AnalysisModel)
    return m, tea.Batch(typerCmd, analysisCmd)
}

func (m MainModel) View() string {
    if (m.timer.Timedout()) {
        return m.analysis.View()
    }
    return fmt.Sprintf("%s\n%s", m.timer.View(), m.typer.View())
}


func main() {
    initialModel := MainModel{
        timer: timer.NewWithInterval(timeout, time.Second),
        typer: TyperModel{
            words: createTest(),
            idx: 0,
        },
        analysis: AnalysisModel{
            wpm: 0,
            accuracy: 0.0,
        },
    }
    p := tea.NewProgram(initialModel)
    if _, err := p.Run(); err != nil {
        fmt.Println("Error starting game:", err)
        os.Exit(1)
    }

}

