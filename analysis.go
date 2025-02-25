package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type AnalysisModel struct {
    time int // amount of time taken
    words int // number of correct words
    correct int // number of correct characters
    characters int // total characters typed 
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
        }
    }
    return am, nil
}

func (am AnalysisModel) View() string {
    wpm := int(float64(am.words) / float64(am.time) * 60.0)
    accuracy := float64(am.correct) / float64(am.characters) * 100
    wpmText := white.Render(fmt.Sprintf("wpm: %d", wpm))
    accuracyText := white.Render(fmt.Sprintf("acc: %.2f", accuracy))
    return wpmText + "\n" + accuracyText 
}
