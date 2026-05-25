package main

import (
	"fmt"
	"os"

	"github.com/WarrenWu4/bananatype/internal/foundation/logger"
	"github.com/WarrenWu4/bananatype/internal/session"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	Build        = "dev"
	wordPath     = "./resources/word_bank.txt"
	settingsPath = "./resources/settings.json"
)

func main() {
	if Build == "prod" {
		wordPath = "/usr/share/bananatype/word_bank.txt"
		settingsPath = os.Getenv("HOME") + "/.local/state/bananatype/settings.json"
		logger.InitLogger(os.Getenv("HOME") + "/.local/state/banantype/log.txt")
	}
	os.OpenFile(settingsPath, os.O_APPEND|os.O_CREATE, 0644)
	initialModel := session.New()
	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting game:", err)
		os.Exit(1)
	}
}
