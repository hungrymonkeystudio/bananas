package session_test

import (
	"testing"

	"github.com/WarrenWu4/bananatype/internal/session"
	tea "github.com/charmbracelet/bubbletea"
)

func TestSession_EscTogglesSettings(t *testing.T) {
	model := session.New()

	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updated.(session.Model)

	// After ESC, should be in settings state (view should contain "Settings")
	view := model.View()
	if !contains(view, "Settings") {
		t.Error("ESC should show settings screen")
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updated.(session.Model)

	view = model.View()
	if contains(view, "Settings") {
		t.Error("second ESC should hide settings screen")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
