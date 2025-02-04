package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type AnalysisModel struct {
    wpm int
    accuracy float64
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
        case "q":
            return am, tea.Quit
        }
    }
    return am, nil
}

func (am AnalysisModel) View() string {
    wpmText := white.Render(fmt.Sprintf("wpm: %d", am.wpm))
    accuracyText := white.Render(fmt.Sprintf("acc: %.2f", am.accuracy))
    return wpmText + "\n" + accuracyText 
}
