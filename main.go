package main

import (
    "bufio"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const timeout = time.Second * 30
var COMMONWORDS []string

type MainModel struct {
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
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return m, tea.Quit
        case "enter":
            m.typer = NewTyper()
            m.timer = NewTimerModel(timeout)
            return m, nil
        }
    }
    if (m.timer.done) {
        updatedAnalysis, analysisCmd := m.analysis.Update(msg)
        m.analysis = updatedAnalysis.(AnalysisModel)
        return m, analysisCmd
    }
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
        m.analysis.time = int(timeout.Seconds())
        m.analysis.words = m.typer.totalWords
        m.analysis.correct = m.typer.totalCorrect
        m.analysis.characters = m.typer.totalTyped
        outputLines := strings.Split(m.analysis.View(), "\n")
        for i := 0; i < len(outputLines) ; i++ {
            output += strings.Repeat(" ", paddingX) + outputLines[i] + "\n"
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

func setup() {
    // get common words from file
    file, _ := os.Open("common-words.txt")
    defer file.Close()
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        COMMONWORDS = append(COMMONWORDS, scanner.Text())
    }
}

func main() {
    // run this function to initialize important shit
    setup()
    initialModel := MainModel{
        timer: NewTimerModel(timeout),
        typer: NewTyper(),
        analysis: AnalysisModel{
            time: 0,
            words: 0,
            correct: 0,
            characters: 0,
        },
        width: 120,
        height: 8,
    }
    p := tea.NewProgram(initialModel, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println("Error starting game:", err)
        os.Exit(1)
    }

}

