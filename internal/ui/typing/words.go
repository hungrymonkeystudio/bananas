package typing

import (
	"bufio"
	"embed"
	"math"
	"math/rand"
	"strings"

	"github.com/WarrenWu4/bananatype/internal/foundation/logger"
)

//go:embed resources/word_bank.txt
var wordBankFS embed.FS

const MAXLINES = 3

// max characters per line is either 60
// or if screen is smaller than do
// screen width - 16 (padding on each side)
var MAXCHARPERLINE = 60

var COMMONWORDS = []string{}
var WORD_BANK_MAP = map[int][]string{}

func loadWordsFromFile() {
	if len(COMMONWORDS) > 0 {
		return
	}
	logger.Log(logger.INFO, "Loading words from embedded file: resources/word_bank.txt")
	file, err := wordBankFS.Open("resources/word_bank.txt")
	if err != nil {
		logger.Log(logger.ERROR, "Failed to open embedded word bank: "+err.Error())
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		if len(word) > 0 {
			WORD_BANK_MAP[len(word)] = append(WORD_BANK_MAP[len(word)], word)
		}
		COMMONWORDS = append(COMMONWORDS, word)
	}
}

func createLine() ([]string, []string, []int) {
	loadWordsFromFile()
	var result []string
	var colorResult []string
	var sizeResult []int
	currChars := 0
	prevWord := ""
	for currChars < MAXCHARPERLINE {
		length := int(math.Round(rand.NormFloat64()*1.5 + 5))
		length = min(14, max(length, 1))
		words := WORD_BANK_MAP[length]
		selectedWord := words[rand.Intn(len(words))]
		for selectedWord == prevWord && len(words) > 1 {
			selectedWord = words[rand.Intn(len(words))]
		}
		prevWord = selectedWord
		if currChars+len(selectedWord) <= MAXCHARPERLINE {
			result = append(result, selectedWord)
			colorResult = append(colorResult, strings.Repeat("g", len(selectedWord)))
			sizeResult = append(sizeResult, len(selectedWord))
		}
		currChars += 1 + len(selectedWord)
	}
	return result, colorResult, sizeResult
}

func checkWordCorrect(ogSize int, colors string) bool {
	if len(colors) != ogSize {
		return false
	}
	for i := 0; i < ogSize; i++ {
		if colors[i] != 'w' {
			return false
		}
	}
	return true
}
