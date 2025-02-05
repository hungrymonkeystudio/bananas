package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const timeout = time.Second * 30

type MainModel struct {
    timer TimerModel 
    typer TyperModel
    analysis AnalysisModel
}

func (m MainModel) Init() tea.Cmd {
    return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "q":
            return m, tea.Quit
        case "enter":
            m.typer.words = createTest()
            m.typer.idx = 0
        }
    }
    updatedTimer, timerCmd := m.timer.Update(msg)
    m.timer = updatedTimer.(TimerModel)
    updatedTyper, typerCmd := m.typer.Update(msg)
    m.typer = updatedTyper.(TyperModel)
    updatedAnalysis, analysisCmd := m.analysis.Update(msg)
    m.analysis = updatedAnalysis.(AnalysisModel)
    return m, tea.Batch(timerCmd, typerCmd, analysisCmd)
}

func (m MainModel) View() string {
    if (m.timer.done) {
        return m.analysis.View()
    }
    return fmt.Sprintf("%s\n%s", m.timer.View(), m.typer.View())
}


func main() {
    initialModel := MainModel{
        timer: NewTimerModel(timeout),
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

