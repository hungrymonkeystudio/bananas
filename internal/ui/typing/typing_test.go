package typing_test

import (
	"testing"

	"github.com/WarrenWu4/bananatype/internal/ui/typing"
	tea "github.com/charmbracelet/bubbletea"
)

func TestTypingInputLogic(t *testing.T) {
	ty := typing.NewTypingModel()

	lines := ty.GetLines()
	if len(lines) == 0 || len(lines[0]) == 0 || len(lines[0][0]) == 0 {
		t.Fatal("Typing lines not initialized correctly")
	}
	firstWord := lines[0][0]
	firstChar := string(firstWord[0])

	msg := tea.KeyMsg{Runes: []rune(firstChar), Type: tea.KeyRunes}
	updatedModel, _ := ty.Update(msg)
	ty = updatedModel.(typing.TypingModel)

	if ty.TotalTyped != 1 {
		t.Errorf("Expected TotalTyped to be 1, got %d", ty.TotalTyped)
	}
	if ty.TotalCorrect != 1 {
		t.Errorf("Expected TotalCorrect to be 1, got %d", ty.TotalCorrect)
	}

	msgWrong := tea.KeyMsg{Runes: []rune("~"), Type: tea.KeyRunes}
	updatedModel, _ = ty.Update(msgWrong)
	ty = updatedModel.(typing.TypingModel)

	if ty.TotalTyped != 2 {
		t.Errorf("Expected TotalTyped to be 2, got %d", ty.TotalTyped)
	}
	if ty.TotalCorrect != 1 {
		t.Errorf("Expected TotalCorrect to be 1, still")
	}
}
