// results screen that shows typing stats

package results

import (
	"fmt"

	"github.com/WarrenWu4/bananatype/internal/foundation/theme"
	tea "github.com/charmbracelet/bubbletea"
)

const RESULTS_INSTRUCTIONS = "CTRL+C to exit\nENTER to retry"

type ResultsModel struct {
	Time       float64 // amount of time taken in seconds
	Words      int     // number of correct words
	Correct    int     // number of correct characters
	Characters int     // total characters typed
}

func NewResultsModel() ResultsModel {
	return ResultsModel{
		Time:       0,
		Words:      0,
		Correct:    0,
		Characters: 0,
	}
}

func (rm ResultsModel) Init() tea.Cmd {
	return nil
}

func (rm ResultsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return rm, tea.Quit
		case "enter":
			return rm, func() tea.Msg { return rm }
		}
	}
	return rm, nil
}

func (rm ResultsModel) View() string {
	wpm := 0
	if rm.Time > 0 {
		wpm = int(float64(rm.Words) / rm.Time * 60.0)
	}
	accuracy := 0.0
	if rm.Characters > 0 {
		accuracy = float64(rm.Correct) / float64(rm.Characters) * 100
	}
	wpmText := theme.White.Render(fmt.Sprintf("wpm: %d", wpm))
	accuracyText := theme.White.Render(fmt.Sprintf("acc: %.2f", accuracy))
	return wpmText + "\n" + accuracyText + "\n\n" + theme.Instructions.Render(RESULTS_INSTRUCTIONS)
}
