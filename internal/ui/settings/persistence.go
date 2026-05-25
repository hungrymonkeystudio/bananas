package settings

import (
	"encoding/json"
	"os"
	"slices"

	"github.com/WarrenWu4/bananatype/internal/foundation/paths"
)

func (m SettingsModel) writeSettings() {
	basePath := paths.GetResourcePath()
	file, err := os.Create(basePath + "/settings.json")
	if err != nil {
		return
	}
	defer file.Close()

	data := map[string]any{
		"activeTime":      m.ActiveTime,
		"activeWords":     m.ActiveWords,
		"activeTyperMode": m.ActiveTyperMode,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(data)
}

func readSettings(m *SettingsModel) {
	basePath := paths.GetResourcePath()
	file, err := os.Open(basePath + "/settings.json")
	if err != nil {
		return
	}
	defer file.Close()

	var data struct {
		ActiveTime      int    `json:"activeTime"`
		ActiveWords     int    `json:"activeWords"`
		ActiveTyperMode string `json:"activeTyperMode"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return
	}

	if slices.Contains(m.times, data.ActiveTime) {
		m.ActiveTime = data.ActiveTime
		m.timeIdx = slices.Index(m.times, data.ActiveTime)
	}

	if slices.Contains(m.words, data.ActiveWords) {
		m.ActiveWords = data.ActiveWords
		m.wordIdx = slices.Index(m.words, data.ActiveWords)
	}

	if data.ActiveTyperMode == "timer" || data.ActiveTyperMode == "words" {
		m.ActiveTyperMode = data.ActiveTyperMode
	}
}
