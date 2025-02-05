package main

import (
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type TyperModel struct {
    wordList []string
    wordIdx int
    currWord []string
    charIdx int
    completeWords []string
} 

func NewTyper() TyperModel {
    return TyperModel{
        wordList: createTest(),
        wordIdx: 0,
        currWord: []string{},
        charIdx: 0,
        completeWords: []string{},
    }
}

func createTest() []string {
	var result []string
	for i := 0; i < 20; i++ {
		randomIndex := rand.Intn(len(commonWords))
		result = append(result, commonWords[randomIndex])
	}
    return result
}

func typerComp(userLetter, msgLetter string) bool {
    return userLetter == msgLetter
}

func (tym TyperModel) Init() tea.Cmd {
    return nil
}

func (tym TyperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case " ":
            if tym.wordIdx < len(tym.wordList) {
                tym.wordIdx++
                tym.charIdx = 0
                tym.completeWords = append(tym.completeWords, strings.Join(tym.currWord, ""))
                tym.currWord = []string{}
            } else {
                return tym, tea.Quit
            }
        case "backspace":
            if (tym.charIdx >= 0) {
                tym.currWord = tym.currWord[:len(tym.currWord)-1]
                tym.charIdx--
            }
        default:
            currentWord := tym.wordList[tym.wordIdx]
            if (tym.charIdx < len(currentWord)) {
                if typerComp(msg.String(), string(currentWord[tym.charIdx])) {
                    tym.currWord = append(tym.currWord, white.Render(msg.String()))
                } else {
                    tym.currWord = append(tym.currWord, red.Render(msg.String()))
                }
                tym.charIdx++
            } else {
                tym.currWord = append(tym.currWord, red.Render(msg.String()))
            }
        }
    }
    return tym, nil
}

func (tym TyperModel) View() string {   
    wordListStr := gray.Render(strings.Join(tym.wordList, " "))
    currWordText := strings.Join(tym.currWord, "")
    completeWordsText := strings.Join(tym.completeWords, " ") 
    return lipgloss.NewStyle().Render(wordListStr + "\n" + completeWordsText + " " + currWordText)
}





