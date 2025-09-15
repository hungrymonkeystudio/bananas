// analysis screen that shows typing stats

package analysis 

import (
	"fmt"
	colors "bananas/pkg/colors"
	tea "github.com/charmbracelet/bubbletea"
)

const ANALYSIS_INSTRUCTIONS = "CTRL+C to exit\nENTER to retry"

type AnalysisModel struct {
    Time int // amount of time taken
    Words int // number of Correct words
    Correct int // number of Correct characters
    Characters int // total characters typed 
}

func NewAnalysisModel() AnalysisModel{
    return AnalysisModel{
        Time: 0,
        Words: 0,
        Correct: 0,
        Characters: 0,
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
    wpm := int(float64(am.Words) / float64(am.Time) * 60.0)
    accuracy := float64(am.Correct) / float64(am.Characters) * 100
    wpmText := colors.White.Render(fmt.Sprintf("wpm: %d", wpm))
    accuracyText := colors.White.Render(fmt.Sprintf("acc: %.2f", accuracy))
    return wpmText + "\n" + accuracyText + "\n\n" + colors.Instructions.Render(ANALYSIS_INSTRUCTIONS)
}
