package main

import (
    "bufio"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var COMMONWORDS []string

type MainModel struct {
    settings SettingsModel
    timer TimerModel 
    typer TyperModel
    analysis AnalysisModel
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
            m.settings.show = !m.settings.show
            return m, nil
        }
    case SettingsModel:
        tea.Println("Settings Updated")
        m.typer = NewTyper()
        m.timer = NewTimerModel(time.Second*time.Duration(m.settings.activeTime))
        m.settings.show = !m.settings.show
    case AnalysisModel:
        m.typer = NewTyper()
        m.timer = NewTimerModel(time.Second*time.Duration(m.settings.activeTime))
    }
    // local updates that are dependent on which view is active
    if m.settings.show { // only update settings when settings show
        updatedSettings, settingsCmd := m.settings.Update(msg)
        m.settings = updatedSettings.(SettingsModel)
        return m, settingsCmd
    } else if (m.timer.done) { // only update analysis when timer is done
        updatedAnalysis, analysisCmd := m.analysis.Update(msg)
        m.analysis = updatedAnalysis.(AnalysisModel)
        return m, analysisCmd
    }
    // otherwise update timer and typer
    updatedTimer, timerCmd := m.timer.Update(msg)
    m.timer = updatedTimer.(TimerModel)
    updatedTyper, typerCmd := m.typer.Update(msg)
    m.typer = updatedTyper.(TyperModel)
    return m, tea.Batch(timerCmd, typerCmd)
}

func (m MainModel) View() string {
    output := ""
    paddingY := (m.height - MAXLINES+1) / 2
	paddingX := (m.width - MAXCHARPERLINE) / 2
    // top padding
    output += strings.Repeat("\n", paddingY)
    // left padding
    if (m.timer.done) {
        m.analysis.time = int(m.settings.times[m.settings.timeIdx])
        m.analysis.words = m.typer.totalWords
        m.analysis.correct = m.typer.totalCorrect
        m.analysis.characters = m.typer.totalTyped
        outputLines := strings.Split(m.analysis.View(), "\n")
        for i := 0; i < len(outputLines) ; i++ {
            output += strings.Repeat(" ", paddingX) + outputLines[i] + "\n"
        }
    } else if m.settings.show {
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

func (m MainModel) Restart() tea.Cmd {
    return func() tea.Msg {
        s := NewSettingsModel()
        ti := NewTimerModel(time.Second * time.Duration(s.times[s.timeIdx]))
        ty := NewTyper()
        a := NewAnalysisModel()
        m = MainModel{
            timer: ti,
            typer: ty,
            analysis: a,
            settings: s,
            width: 120,
            height: 8,
        }
        return ""
    }
} 

func setup() MainModel {
    // get common words from file
    file, _ := os.Open("common-words.txt")
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        COMMONWORDS = append(COMMONWORDS, scanner.Text())
    }
    // initialize main model
    s := NewSettingsModel()
    ti := NewTimerModel(time.Second * time.Duration(s.activeTime))
    ty := NewTyper()
    a := NewAnalysisModel()
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

