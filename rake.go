package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

func getLinesFromFile(filename string) []string {
	content, _ := ioutil.ReadFile(filename)
	return strings.Split(string(content), "\n")
}

func splitIntoWords(text string) []string {
	words := []string{}
	wordSplitRegex := regexp.MustCompile("[\\p{L}\\d_]+")
	splitWords := wordSplitRegex.FindAllString(text, -1)
	for _, word := range splitWords {
		currentWord := strings.ToLower(strings.TrimSpace(word))
		if currentWord != "" {
			words = append(words, currentWord)
		}
	}
	return words
}

func getStopWordRegex() string {
	stopwords := getLinesFromFile("SmartStoplist.txt")
	stopwordRegexPattern := []string{}
	for _, word := range stopwords {
		wordRegex := fmt.Sprintf(`(?:\A|\z|\s)%s(?:\A|\z|\s)`, word)
		stopwordRegexPattern = append(stopwordRegexPattern, wordRegex)
	}
	return `(?i)` + strings.Join(stopwordRegexPattern, "|")
}

func generateCandidatePhrases(text string) []string {
	stopWordRegex := regexp.MustCompile(getStopWordRegex())
	temp := stopWordRegex.ReplaceAllString(text, "|")
	multipleWhitespaceRegex := regexp.MustCompile(`\s\s+`)
	temp = multipleWhitespaceRegex.ReplaceAllString(strings.TrimSpace(temp), " ")

	phraseList := []string{}
	phrases := strings.Split(temp, "|")
	for _, phrase := range phrases {
		phrase = strings.ToLower(phrase)
		if phrase != "" {
			phraseList = append(phraseList, phrase)
		}
	}
	return phraseList
}

func splitIntoSentences(text string) []string {
	splitPattern := regexp.MustCompile(`[.,\/#!$%\^&\*;:{}=\-_~()]`)
	return splitPattern.Split(text, -1)
}

func calculateWordScores(phraseList []string) map[string]float64 {
	frequencies := map[string]int{}
	degrees := map[string]int{}
	for _, phrase := range phraseList {
		words := splitIntoWords(phrase)
		length := len(words)
		degree := length - 1

		for _, word := range words {
			frequencies[word]++
			degrees[word] += degree
		}
	}
	for key := range frequencies {
		degrees[key] = degrees[key] + frequencies[key]
	}

	score := map[string]float64{}

	for key := range frequencies {
		score[key] += (float64(degrees[key]) / float64(frequencies[key]))
	}

	return score
}

func main() {
	sentences := splitIntoSentences("Compatibility of systems of linear constraints over the set of natural numbers. Criteria of compatibility of a system of linear Diophantine equations, strict inequations, and nonstrict inequations are considered. Upper bounds for components of a minimal set of solutions and algorithms of construction of minimal generating sets of solutions for all types of systems are given. These criteria and the corresponding algorithms for constructing a minimal supporting set of solutions can be used in solving all the considered types of systems and systems of mixed types.")
	phraseList := []string{}
	for _, sentence := range sentences {
		phraseList = append(phraseList, generateCandidatePhrases(sentence)...)
		wordScores := calculateWordScores(phraseList)
		fmt.Println(wordScores)
	}
}
