package main

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var commonWords = []string{
	"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
	"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
	"this", "but", "his", "by", "from", "they", "we", "say", "her",
	"she", "or", "an", "will", "my", "one", "all", "would", "there",
	"their", "what", "so", "up", "out", "if", "about", "who", "get",
	"which", "go", "me", "when", "make", "can", "like", "time", "no",
	"just", "him", "know", "take", "people", "into", "year", "your",
	"good", "some", "could", "them", "see", "other", "than", "then",
	"now", "look", "only", "come", "its", "over", "think", "also",
	"back", "after", "use", "two", "how", "our", "work", "first",
	"well", "way", "even", "new", "want", "because", "any", "these",
	"give", "day", "most", "us", "is", "are", "was", "were", "am",
	"theirs", "how", "shes", "through", "me", "myself", "watch", "find",
	"many", "never", "down", "before", "where", "called", "might", "while",
	"too", "next", "made", "here", "know", "point", "few", "lost",
	"does", "long", "those", "by", "more", "heart", "world", "last",
	"left", "should", "call", "hard", "still", "each", "turn", "too",
	"never", "own", "around", "number", "call", "why",
}

const MAXLINES = 3
// max characters per line is either 60
// or if screen is smaller than do
// screen width - 16 (padding on each side)
var MAXCHARPERLINE = 60

type TyperModel struct {
    // text related parameters
    lines [][]string
    linesColor [][]string
    lineIdx int
    wordIdx int
    charIdx int
    skips [][]int // should specify, line, word, and char idx
    currWordSize int // track to handle extra words
} 

func NewTyper() TyperModel {
    fmt.Println("NewTyper")
    // create MAXLINES number of text
    // each line should have max MAXCHARPERLINE characters
    lines := [][]string{}
    colorLines := [][]string{}
    for i := 0; i < MAXLINES; i++ {
        line, colorLine := createLine()
        lines = append(lines, line)
        colorLines = append(colorLines, colorLine)
    }
    return TyperModel{
        lines: lines,
        linesColor: colorLines,
        lineIdx: 0,
        wordIdx: 0,
        charIdx: 0,
        skips: [][]int{},
        currWordSize: len(lines[0][0]),
    }
}

func createLine() ([]string, []string) {
	var result []string
    var colorResult []string
    currChars := 0
    // while all chars less than MAXCHARPERLINE
    for currChars < MAXCHARPERLINE {
		randomIndex := rand.Intn(len(commonWords))
        if (currChars + len(commonWords[randomIndex]) <= MAXCHARPERLINE) {
            result = append(result, commonWords[randomIndex])
            colorResult = append(colorResult, strings.Repeat("g", len(commonWords[randomIndex])))
        }
        currChars += 1 + len(commonWords[randomIndex])
    }
    return result, colorResult
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
            // if word is completed, reset skips
            // otherwise store skips
            if (tym.charIdx >= tym.currWordSize-1) {
                tym.skips = [][]int{}
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
                    tym.lines[2], tym.linesColor[2] = createLine()
                } else {
                    tym.lineIdx += 1
                }
            } else {
                tym.wordIdx += 1
                tym.charIdx = 0
            }
            tym.currWordSize = len(tym.lines[tym.lineIdx][tym.wordIdx])
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
                tym.currWordSize = len(tym.lines[tym.lineIdx][tym.wordIdx])
            } else if tym.charIdx > tym.currWordSize {
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
            if (tym.charIdx >= tym.currWordSize) {
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
                } else {
                    tempRune := []rune(tym.linesColor[tym.lineIdx][tym.wordIdx])
                    tempRune[tym.charIdx] = 'r'
                    tym.linesColor[tym.lineIdx][tym.wordIdx] = string(tempRune)
                }
            }
            tym.charIdx += 1
        }
    }
    return tym, nil 
}

func (tym TyperModel) View() string {   
    // fmt.Println(tym.lines, tym.linesColor)
    // fmt.Println(tym.linesColor)
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





