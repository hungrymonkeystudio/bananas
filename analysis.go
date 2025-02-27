package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const ANALYSIS_INSTRUCTIONS = "CTRL+C to exit\nENTER to retry"

type AnalysisModel struct {
    time int // amount of time taken
    words int // number of correct words
    correct int // number of correct characters
    characters int // total characters typed 
}

func NewAnalysisModel() AnalysisModel{
    return AnalysisModel{
        time: 0,
        words: 0,
        correct: 0,
        characters: 0,
    }
}


func (am AnalysisModel) Init() tea.Cmd {
    return nil
}

func (am AnalysisModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return am, tea.Quit
        case "enter":
            return am, func() tea.Msg { return am }
        }
    }
    return am, nil
}

func (am AnalysisModel) View() string {
    wpm := int(float64(am.words) / float64(am.time) * 60.0)
    accuracy := float64(am.correct) / float64(am.characters) * 100
    wpmText := white.Render(fmt.Sprintf("wpm: %d", wpm))
    accuracyText := white.Render(fmt.Sprintf("acc: %.2f", accuracy))
    return wpmText + "\n" + accuracyText + "\n\n" + instructions.Render(ANALYSIS_INSTRUCTIONS)
}
