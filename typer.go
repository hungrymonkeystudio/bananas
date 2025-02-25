package main

import (
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const MAXLINES = 3
// max characters per line is either 60
// or if screen is smaller than do
// screen width - 16 (padding on each side)
var MAXCHARPERLINE = 60

type TyperModel struct {
    // text related parameters
    lines [][]string
    linesColor [][]string
    wordSizes [][]int // store word sizes for each line to handle complex extra character cases
    lineIdx int
    wordIdx int
    charIdx int
    skips [][]int // should specify, line, word, and char idx
    // analytics
    totalWords int
    totalCorrect int
    totalTyped int
} 

func NewTyper() TyperModel {
    // create MAXLINES number of text
    // each line should have max MAXCHARPERLINE characters
    lines := [][]string{}
    colorLines := [][]string{}
    wordSizes := [][]int{}
    for i := 0; i < MAXLINES; i++ {
        line, colorLine, wordSize := createLine()
        lines = append(lines, line)
        colorLines = append(colorLines, colorLine)
        wordSizes = append(wordSizes, wordSize)
    }
    return TyperModel{
        lines: lines,
        linesColor: colorLines,
        wordSizes: wordSizes,
        lineIdx: 0,
        wordIdx: 0,
        charIdx: 0,
        skips: [][]int{},
        totalWords: 0,
        totalCorrect: 0,
        totalTyped: 0,
    }
}

func createLine() ([]string, []string, []int) {
	var result []string
    var colorResult []string
    var sizeResult []int
    currChars := 0
    // while all chars less than MAXCHARPERLINE
    for currChars < MAXCHARPERLINE {
		randomIndex := rand.Intn(len(COMMONWORDS))
        if (currChars + len(COMMONWORDS[randomIndex]) <= MAXCHARPERLINE) {
            result = append(result, COMMONWORDS[randomIndex])
            colorResult = append(colorResult, strings.Repeat("g", len(COMMONWORDS[randomIndex])))
            sizeResult = append(sizeResult, len(COMMONWORDS[randomIndex]))
        }
        currChars += 1 + len(COMMONWORDS[randomIndex])
    }
    return result, colorResult, sizeResult
}

func checkWordCorrect(ogSize int, colors string) bool {
    if (len(colors) != ogSize) {
        return false
    }
    for i := 0; i < ogSize; i++ {
        if (colors[i] != 'w') {
            return false
        }
    }
    return true
}

func (tym TyperModel) Init() tea.Cmd {
    return nil
}

func (tym TyperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        MAXCHARPERLINE = min(msg.Width-16, 60)
        return tym, nil
    case tea.KeyMsg:
        switch msg.String() {
        case " ":
            // if word is completely correct, reset skips
            // otherwise store skips
            currWordSize := tym.wordSizes[tym.lineIdx][tym.wordIdx]
            if (tym.charIdx >= currWordSize && checkWordCorrect(currWordSize, tym.linesColor[tym.lineIdx][tym.wordIdx])) {
                tym.skips = [][]int{}
                tym.totalWords += 1
            } else {
                tym.skips = append(tym.skips, []int{tym.lineIdx, tym.wordIdx, tym.charIdx})
            }
            // if word is last word in the line
            // lineIdx + 1, wordIdx = 0, charIdx = 0
            // otherwise move to next line
            // wordIdx + 1, charIdx = 0 
            if (tym.wordIdx == len(tym.lines[tym.lineIdx]) - 1) {
                tym.wordIdx = 0
                tym.charIdx = 0
                // if not first line, push all lines up and create new line words
                if (tym.lineIdx != 0) {
                    tym.lines[0] = tym.lines[1]
                    tym.linesColor[0] = tym.linesColor[1]
                    tym.lines[1] = tym.lines[2]
                    tym.linesColor[1] = tym.linesColor[2]
                    tym.wordSizes[0] = tym.wordSizes[1]
                    tym.wordSizes[1] = tym.wordSizes[2]
                    tym.lines[2], tym.linesColor[2], tym.wordSizes[2] = createLine()
                } else {
                    tym.lineIdx += 1
                }
            } else {
                tym.wordIdx += 1
                tym.charIdx = 0
            }
        case "backspace":
            // if charIdx is 0 and skip isn' t empty
            // start at skipped position
            // if charIdx is 0 and skip is empty
            // do nothing
            // if current charIdx is at an extra character
            // delete from word and color
            // otherwise simply change color back to gray
            if tym.charIdx <= 0 && len(tym.skips) <= 0 {
                return tym, nil 
            } 
            if tym.charIdx <= 0 && len(tym.skips) > 0 {
                tym.lineIdx = tym.skips[len(tym.skips)-1][0]
                tym.wordIdx = tym.skips[len(tym.skips)-1][1]
                tym.charIdx = tym.skips[len(tym.skips)-1][2]
                tym.skips = tym.skips[:len(tym.skips)-1]
            } else if tym.charIdx > tym.wordSizes[tym.lineIdx][tym.wordIdx] {
                tym.charIdx -= 1
                tym.lines[tym.lineIdx][tym.wordIdx] = tym.lines[tym.lineIdx][tym.wordIdx][:tym.charIdx]
                tym.linesColor[tym.lineIdx][tym.wordIdx] = tym.linesColor[tym.lineIdx][tym.wordIdx][:tym.charIdx]
            } else {
                tym.charIdx -= 1
                tempRune := []rune(tym.linesColor[tym.lineIdx][tym.wordIdx])
                tempRune[tym.charIdx] = 'g'
                tym.linesColor[tym.lineIdx][tym.wordIdx] = string(tempRune)
            }
        default:
            // if there are extra characters inputed
            // append to end of word and color word
            if (tym.charIdx >= tym.wordSizes[tym.lineIdx][tym.wordIdx]) {
                tempRune := []rune(tym.lines[tym.lineIdx][tym.wordIdx])
                tempRune = append(tempRune, []rune(msg.String())...)
                tym.lines[tym.lineIdx][tym.wordIdx] = string(tempRune)
                tempRune = []rune(tym.linesColor[tym.lineIdx][tym.wordIdx])
                tempRune = append(tempRune, 'r')
                tym.linesColor[tym.lineIdx][tym.wordIdx] = string(tempRune)
            } else {
                // if correct change color to white
                // if wrong change color to red
                if (msg.String() == string(tym.lines[tym.lineIdx][tym.wordIdx][tym.charIdx])) {
                    tempRune := []rune(tym.linesColor[tym.lineIdx][tym.wordIdx])
                    tempRune[tym.charIdx] = 'w'
                    tym.linesColor[tym.lineIdx][tym.wordIdx] = string(tempRune)
                    tym.totalCorrect += 1
                } else {
                    tempRune := []rune(tym.linesColor[tym.lineIdx][tym.wordIdx])
                    tempRune[tym.charIdx] = 'r'
                    tym.linesColor[tym.lineIdx][tym.wordIdx] = string(tempRune)
                }
            }
            tym.charIdx += 1
            tym.totalTyped += 1
        }
    }
    return tym, nil 
}

func (tym TyperModel) View() string {   
    output := ""
    cursorOnSpace := true 
    // print out each line based on color
    for i := 0; i < MAXLINES; i++ { // line num
        for j := 0; j < len(tym.lines[i]); j++ { // word num
            for k := 0; k < len(tym.lines[i][j]); k++ {
                color := tym.linesColor[i][j][k]
                letter := tym.lines[i][j][k]
                if (k == tym.charIdx && i == tym.lineIdx && j == tym.wordIdx) {
                    cursorOnSpace = false 
                    output += cursor.Render(string(letter))
                } else {
                    switch color {
                    case 'r':
                        output += red.Render(string(letter))
                    case 'g':
                        output += gray.Render(string(letter))
                    case 'w':
                        output += white.Render(string(letter))
                    default:
                        output += string(letter)
                    }
                }
            }
            if (i == tym.lineIdx && j == tym.wordIdx && cursorOnSpace) {
                output += cursor.Render(" ")
            } else {
                output += " "
            }
        }
        output += "\n"
    }
    return output 
}





