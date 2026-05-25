package typing

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func FuzzTypingUpdate(f *testing.F) {
	f.Add([]byte{0x01, 0x02, 0x03})
	f.Add([]byte{0x00, 0x00, 0x00, 0x00, 0x00})
	f.Add([]byte{0x01, 0x01, 0x01, 0x01, 0x01})
	f.Add([]byte{2, 3, 4, 0, 2, 3, 4, 0, 2, 3, 4, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		MAXCHARPERLINE = 60
		model := NewTypingModel()

		for _, b := range data {
			var msg tea.Msg
			switch {
			case b == 0:
				msg = tea.KeyMsg{Type: tea.KeySpace}
			case b == 1:
				msg = tea.KeyMsg{Type: tea.KeyBackspace}
			case b >= 2 && b <= 27:
				ch := rune('a' + b - 2)
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}}
			default:
				continue
			}
			updated, _ := model.Update(msg)
			model = updated.(TypingModel)
		}

		if model.LineIdx < 0 || model.LineIdx >= MAXLINES {
			t.Fatalf("LineIdx out of bounds: %d", model.LineIdx)
		}
		if model.WordIdx < 0 || model.WordIdx >= len(model.Lines[model.LineIdx]) {
			t.Fatalf("WordIdx out of bounds: %d (line has %d words)",
				model.WordIdx, len(model.Lines[model.LineIdx]))
		}
		if model.CharIdx < 0 {
			t.Fatalf("CharIdx negative: %d", model.CharIdx)
		}
		if len(model.Lines) != MAXLINES {
			t.Fatalf("lines length changed: got %d, want %d", len(model.Lines), MAXLINES)
		}
		for i := range model.Lines {
			if len(model.Lines[i]) != len(model.LinesColor[i]) {
				t.Fatalf("lines/linesColor length mismatch at line %d: %d vs %d",
					i, len(model.Lines[i]), len(model.LinesColor[i]))
			}
			for j := range model.Lines[i] {
				if len(model.Lines[i][j]) != len(model.LinesColor[i][j]) {
					t.Fatalf("word/color length mismatch at [%d][%d]: %d vs %d",
						i, j, len(model.Lines[i][j]), len(model.LinesColor[i][j]))
				}
			}
		}
		if model.TotalCorrect > model.TotalTyped {
			t.Fatalf("TotalCorrect (%d) > TotalTyped (%d)", model.TotalCorrect, model.TotalTyped)
		}
	})
}
