// Simulate a typing test with a sequence of inputs and delays

package performance

import (
	"github.com/WarrenWu4/bananatype/internal/session"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var export = flag.Bool("export", false, "save results to a file")

type TestTypingModel struct {
	main      session.Model
	frameChan chan bool // Channel to signal a frame update
}

func (m TestTypingModel) Init() tea.Cmd {
	return nil
}

func (m TestTypingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Signal that an update (frame) occurred
	if m.frameChan != nil {
		select {
		case m.frameChan <- true:
		default:
		}
	}
	updatedMain, mainCmd := m.main.Update(msg)
	m.main = updatedMain.(session.Model)
	return m, tea.Batch(mainCmd)
}

func (m TestTypingModel) View() string {
	return m.main.View()
}

type Collector interface {
	Name() string
	Unit() string
	Collect(t time.Time) float64
	Reset()
	SetChannel(chan bool)
}

type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type SummaryStats struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`
}

func NewSummaryStats() SummaryStats {
	return SummaryStats{
		Min: math.MaxFloat64,
		Max: 0.0,
		Avg: 0.0,
	}
}

type ExportData struct {
	Metadata struct {
		Metric string `json:"metric"`
		Title  string `json:"title"`
	} `json:"metadata"`
	Data []DataPoint `json:"data"`
}

func simulate(inputs []string, delayMs int, testTypingModel TestTypingModel, collector Collector) []DataPoint {
	collector.SetChannel(testTypingModel.frameChan)
	collector.Reset()
	p := tea.NewProgram(testTypingModel, tea.WithInput(nil), tea.WithoutRenderer(), tea.WithFPS(60))
	data := []DataPoint{}
	stopMonitoring := make(chan bool)
	// start monitoring in goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-stopMonitoring:
				return
			case t := <-ticker.C:
				data = append(data, DataPoint{
					Timestamp: t,
					Value:     collector.Collect(t),
				})
			}
		}
	}()
	// run TUI in goroutine
	go func() {
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()
	// simulate inputs
	time.Sleep(500 * time.Millisecond)
	for _, input := range inputs {
		p.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(input)})
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
	// exit
	p.Quit()
	time.Sleep(200 * time.Millisecond)
	stopMonitoring <- true
	return data
}

func summarize(data []DataPoint) SummaryStats {
	summaryStats := NewSummaryStats()
	for _, d := range data {
		summaryStats.Avg += d.Value
		summaryStats.Min = min(summaryStats.Min, d.Value)
		summaryStats.Max = max(summaryStats.Max, d.Value)
	}
	summaryStats.Avg /= float64(len(data))
	return summaryStats
}

func exportData(data []DataPoint, collector Collector) {
	exportObj := ExportData{
		Data: data,
	}
	exportObj.Metadata.Metric = collector.Unit()
	exportObj.Metadata.Title = collector.Name()
	jsonData, err := json.MarshalIndent(exportObj, "", "  ")
	if err == nil {
		_ = os.MkdirAll("test_results", 0755)
		filename := fmt.Sprintf("%s_%d.json", collector.Name(), time.Now().Unix())
		filePath := filepath.Join("test_results", filename)
		_ = os.WriteFile(filePath, jsonData, 0644)
	}
}

func runTypingTest(wpm int, accuracy float64, collector Collector) []DataPoint {
	testTypingModel := TestTypingModel{
		main:      session.New(),
		frameChan: make(chan bool, 100),
	}
	delayMs := int(60000 / (wpm * 5))
	inputs := []string{}
	for _, line := range testTypingModel.main.GetTyping().GetLines() {
		for _, word := range line {
			for _, char := range word {
				inputs = append(inputs, string(char))
			}
			inputs = append(inputs, " ")
		}
	}
	return simulate(inputs, delayMs, testTypingModel, collector)
}
