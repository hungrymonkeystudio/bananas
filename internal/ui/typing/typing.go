package typing

import (
	"github.com/WarrenWu4/bananatype/internal/foundation/theme"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const TYPER_INSTRUCTIONS = "CTRL+C to quit program\nESC to open settings"

type TypingModel struct {
	// text related parameters
	Lines      [][]string
	LinesColor [][]string
	WordSizes  [][]int // store word sizes for each line to handle complex extra character cases
	LineIdx    int
	WordIdx    int
	CharIdx    int
	Skips      [][]int // should specify, line, word, and char idx
	// analytics
	TotalWords   int
	TotalCorrect int
	TotalTyped   int
}

func NewTypingModel() TypingModel {
	lines := [][]string{}
	colorLines := [][]string{}
	wordSizes := [][]int{}
	for i := 0; i < MAXLINES; i++ {
		line, colorLine, wordSize := createLine()
		lines = append(lines, line)
		colorLines = append(colorLines, colorLine)
		wordSizes = append(wordSizes, wordSize)
	}
	return TypingModel{
		Lines:        lines,
		LinesColor:   colorLines,
		WordSizes:    wordSizes,
		LineIdx:      0,
		WordIdx:      0,
		CharIdx:      0,
		Skips:        [][]int{},
		TotalWords:   0,
		TotalCorrect: 0,
		TotalTyped:   0,
	}
}

func (tym TypingModel) GetLines() [][]string {
	return tym.Lines
}

func (tym TypingModel) Init() tea.Cmd {
	return nil
}

func (tym TypingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		MAXCHARPERLINE = min(msg.Width-16, 60)
		return tym, nil
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			currWordSize := tym.WordSizes[tym.LineIdx][tym.WordIdx]
			if tym.CharIdx >= currWordSize && checkWordCorrect(currWordSize, tym.LinesColor[tym.LineIdx][tym.WordIdx]) {
				tym.Skips = [][]int{}
				tym.TotalWords += 1
			} else {
				tym.Skips = append(tym.Skips, []int{tym.LineIdx, tym.WordIdx, tym.CharIdx})
			}
			if tym.WordIdx == len(tym.Lines[tym.LineIdx])-1 {
				tym.WordIdx = 0
				tym.CharIdx = 0
				if tym.LineIdx != 0 {
					tym.Lines[0] = tym.Lines[1]
					tym.LinesColor[0] = tym.LinesColor[1]
					tym.Lines[1] = tym.Lines[2]
					tym.LinesColor[1] = tym.LinesColor[2]
					tym.WordSizes[0] = tym.WordSizes[1]
					tym.WordSizes[1] = tym.WordSizes[2]
					tym.Lines[2], tym.LinesColor[2], tym.WordSizes[2] = createLine()
				} else {
					tym.LineIdx += 1
				}
			} else {
				tym.WordIdx += 1
				tym.CharIdx = 0
			}
		case "backspace":
			if tym.CharIdx <= 0 && len(tym.Skips) <= 0 {
				return tym, nil
			}
			if tym.CharIdx <= 0 && len(tym.Skips) > 0 {
				tym.LineIdx = tym.Skips[len(tym.Skips)-1][0]
				tym.WordIdx = tym.Skips[len(tym.Skips)-1][1]
				tym.CharIdx = tym.Skips[len(tym.Skips)-1][2]
				tym.Skips = tym.Skips[:len(tym.Skips)-1]
			} else if tym.CharIdx > tym.WordSizes[tym.LineIdx][tym.WordIdx] {
				tym.CharIdx -= 1
				tym.Lines[tym.LineIdx][tym.WordIdx] = tym.Lines[tym.LineIdx][tym.WordIdx][:tym.CharIdx]
				tym.LinesColor[tym.LineIdx][tym.WordIdx] = tym.LinesColor[tym.LineIdx][tym.WordIdx][:tym.CharIdx]
			} else {
				tym.CharIdx -= 1
				tempRune := []rune(tym.LinesColor[tym.LineIdx][tym.WordIdx])
				tempRune[tym.CharIdx] = 'g'
				tym.LinesColor[tym.LineIdx][tym.WordIdx] = string(tempRune)
			}
		default:
			if len(msg.String()) > 1 {
				return tym, nil
			}
			if tym.CharIdx >= tym.WordSizes[tym.LineIdx][tym.WordIdx] {
				tempRune := []rune(tym.Lines[tym.LineIdx][tym.WordIdx])
				tempRune = append(tempRune, []rune(msg.String())...)
				tym.Lines[tym.LineIdx][tym.WordIdx] = string(tempRune)
				tempRune = []rune(tym.LinesColor[tym.LineIdx][tym.WordIdx])
				tempRune = append(tempRune, 'r')
				tym.LinesColor[tym.LineIdx][tym.WordIdx] = string(tempRune)
			} else {
				if msg.String() == string(tym.Lines[tym.LineIdx][tym.WordIdx][tym.CharIdx]) {
					tempRune := []rune(tym.LinesColor[tym.LineIdx][tym.WordIdx])
					tempRune[tym.CharIdx] = 'w'
					tym.LinesColor[tym.LineIdx][tym.WordIdx] = string(tempRune)
					tym.TotalCorrect += 1
				} else {
					tempRune := []rune(tym.LinesColor[tym.LineIdx][tym.WordIdx])
					tempRune[tym.CharIdx] = 'r'
					tym.LinesColor[tym.LineIdx][tym.WordIdx] = string(tempRune)
				}
			}
			tym.CharIdx += 1
			tym.TotalTyped += 1
		}
	}
	return tym, nil
}

func (tym TypingModel) View() string {
	output := ""
	cursorOnSpace := true
	for i := 0; i < MAXLINES; i++ {
		for j := 0; j < len(tym.Lines[i]); j++ {
			for k := 0; k < len(tym.Lines[i][j]); k++ {
				color := tym.LinesColor[i][j][k]
				letter := tym.Lines[i][j][k]
				if k == tym.CharIdx && i == tym.LineIdx && j == tym.WordIdx {
					cursorOnSpace = false
					output += theme.Cursor.Render(string(letter))
				} else {
					switch color {
					case 'r':
						output += theme.Red.Render(string(letter))
					case 'g':
						output += theme.Gray.Render(string(letter))
					case 'w':
						output += theme.White.Render(string(letter))
					default:
						output += string(letter)
					}
				}
			}
			if i == tym.LineIdx && j == tym.WordIdx && cursorOnSpace {
				output += theme.Cursor.Render(" ")
			} else {
				output += " "
			}
		}
		output += "\n"
	}
	output += "\n" + theme.Instructions.Render(TYPER_INSTRUCTIONS)
	return output
}

func centerView(content string, paddingX int) string {
	var output string
	for _, line := range strings.Split(content, "\n") {
		output += strings.Repeat(" ", paddingX) + line + "\n"
	}
	return output
}
