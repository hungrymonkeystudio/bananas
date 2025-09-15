package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	settings "bananas/pkg/settings"
	timer "bananas/pkg/timer"
	typer "bananas/pkg/typer"
	analysis "bananas/pkg/analysis"
	tea "github.com/charmbracelet/bubbletea"
)

type MainModel struct {
    settings settings.SettingsModel
    timer timer.TimerModel 
    typer typer.TyperModel
    analysis analysis.AnalysisModel
    width int
    height int
}

func (m MainModel) Init() tea.Cmd {
    return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // global updates that happen regardless of current view 
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "esc":
            m.settings.Show = !m.settings.Show
            return m, nil
        }
    case settings.SettingsModel:
        tea.Println("Settings Updated")
        m.typer = typer.NewTyper()
        m.timer = timer.NewTimerModel(time.Second*time.Duration(m.settings.ActiveTime))
        m.settings.Show = !m.settings.Show
    case analysis.AnalysisModel:
        m.typer = typer.NewTyper()
        m.timer = timer.NewTimerModel(time.Second*time.Duration(m.settings.ActiveTime))
    }
    // local updates that are dependent on which view is active
    if m.settings.Show { // only update settings when settings show
        updatedSettings, settingsCmd := m.settings.Update(msg)
        m.settings = updatedSettings.(settings.SettingsModel)
        return m, settingsCmd
    } else if (m.timer.Done) { // only update analysis when timer is done
        updatedAnalysis, analysisCmd := m.analysis.Update(msg)
        m.analysis = updatedAnalysis.(analysis.AnalysisModel)
        return m, analysisCmd
    }
    // otherwise update timer and typer
    updatedTimer, timerCmd := m.timer.Update(msg)
    m.timer = updatedTimer.(timer.TimerModel)
    updatedTyper, typerCmd := m.typer.Update(msg)
    m.typer = updatedTyper.(typer.TyperModel)
    return m, tea.Batch(timerCmd, typerCmd)
}

func (m MainModel) View() string {
    output := ""
    paddingY := (m.height - typer.MAXLINES+1) / 2
	paddingX := (m.width - typer.MAXCHARPERLINE) / 2
    // top padding
    output += strings.Repeat("\n", paddingY)
    // left padding
    if (m.timer.Done) {
        m.analysis.Time = m.settings.ActiveTime
        m.analysis.Words = m.typer.TotalWords
        m.analysis.Correct = m.typer.TotalCorrect
        m.analysis.Characters = m.typer.TotalTyped
        outputLines := strings.Split(m.analysis.View(), "\n")
        for i := 0; i < len(outputLines) ; i++ {
            output += strings.Repeat(" ", paddingX) + outputLines[i] + "\n"
        }
    } else if m.settings.Show {
        outputLines := strings.Split(m.settings.View(), "\n")
        for _, line := range outputLines {
            output += strings.Repeat(" ", paddingX) + line + "\n"
        }
    } else {
        output += strings.Repeat(" ", paddingX) + m.timer.View() + "\n"
        outputLines := strings.Split(m.typer.View(), "\n")
        for i := 0; i < len(outputLines); i++ {
            output += strings.Repeat(" ", paddingX) + outputLines[i] + "\n"
        }
    }
    return output
}

func setup() MainModel {
    // initialize main model
    s := settings.NewSettingsModel()
    ti := timer.NewTimerModel(time.Second * time.Duration(s.ActiveTime))
    ty := typer.NewTyper()
    a := analysis.NewAnalysisModel()
    return MainModel{
        timer: ti,
        typer: ty,
        analysis: a,
        settings: s,
        width: 120,
        height: 8,
    }
}

func main() {
    // run this function to initialize important shit
    initialModel := setup()
    p := tea.NewProgram(initialModel, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println("Error starting game:", err)
        os.Exit(1)
    }

}

