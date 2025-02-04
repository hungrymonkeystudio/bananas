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

// creates a string of random words separated by spaces
func createTest() string {
	var result []string

	// Generate the random string by picking words from the commonWords list
	for i := 0; i < 20; i++ {
		// Pick a random word from the commonWords slice
		randomIndex := rand.Intn(len(commonWords))
		result = append(result, commonWords[randomIndex])
	}

	// Join the words into a single string, separated by spaces
	return strings.Join(result, " ")
}

type TyperModel struct {
    words string
    idx int
} 

func (tym TyperModel) Init() tea.Cmd {
    return nil
}

func (tym TyperModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c":
            return tym, tea.Quit
        case "q":
            return tym, tea.Quit
        default:
            if tym.idx < len(tym.words) {
                tym.idx++
            }
            return tym, nil
        }
    }
    return tym, nil
}

func (tym TyperModel) View() string {   
    prev := white.Render(tym.words[:tym.idx])
    next := gray.Render(tym.words[tym.idx:])
    return lipgloss.NewStyle().Render(prev + next)
}





