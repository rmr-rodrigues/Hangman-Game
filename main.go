package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
	"unicode"
)

var dictionary = map[string]bool{}
var dictionaryKeys []string

const resourcesPath = "resources/"
const dictionaryFileName = "dictionary.txt"

func main() {

	// 1 - Select a word to print
	// 2 - Printing the game state
	//		- prints the state of the hangman,
	//		- prints the word to be guessed
	//		- prints the rest of the layout
	// 3 - Read user input
	//		- validate it (only letters)
	// 4 - Is a correct guess or not
	//		- if correct, update the guessed letters
	//		- update the trys
	//		- if not correct, update hangman state
	// 5 - Verify game state
	//		- if word is guessed, you win
	//			- update scores
	//		- if word not guessed
	//			- update scores
	//			- verify hangman state
	//				- if hangman complete, game over (continue yes or no)
	//				- hangman not complete, continue game

	guessedLetters := make(map[rune]bool)
	loadDictionary()
	getDictionaryKeys()
	word := getRandomWord()
	selectTwoLetters(word, guessedLetters)
	printLayout(word, guessedLetters, 9, []string{"a", "b", "l"}, 10, 3, 2, 1)

}

// It's easier to random select a word from a slice than of a map!
func getDictionaryKeys() {
	dictionaryKeys = make([]string, len(dictionary))
	i := 0
	for word := range dictionary {
		dictionaryKeys[i] = word
		i++
	}
}

func selectTwoLetters(word string, guessedLetters map[rune]bool) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	f := random.Intn(len(word))
	s := random.Intn(len(word))

	guessedLetters[unicode.ToLower(rune(word[f]))] = true
	guessedLetters[unicode.ToLower(rune(word[s]))] = true
}

func getRandomWord() string {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	for {
		index := random.Intn(len(dictionaryKeys))
		word := dictionaryKeys[index]

		if !dictionary[word] {
			dictionary[word] = true
			return word
		}
	}
}

func loadDictionary() {
	readFile, err := os.Open(resourcesPath + dictionaryFileName)
	if err != nil {
		panic(err)
	}

	if err != nil {
		fmt.Println()
		fmt.Printf("The dictionary file, \"%s\" wasn't found!\nLoading test dictionary!\n", resourcesPath+dictionaryFileName)
		dictionary = map[string]bool{"apple": false, "house": false, "programming": false, "Portugal": false}
		fmt.Println()
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	// TODO filter words that are not according to a regex
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line != "" {
			dictionary[line] = false
		}
	}

	err = readFile.Close()
	if err != nil {
		panic(err)
	}

}

func getHangman(hangmanState int) string {
	data, err := os.ReadFile(fmt.Sprintf(resourcesPath+"/hangman_states/hangman%d", hangmanState))

	if err != nil {
		panic(err)
	}

	return string(data)
}

func wordToGuessingState(word string, guessedLetters map[rune]bool) string {

	guessingState := ""
	for _, letter := range word {
		if letter == ' ' {
			guessingState += " "
		} else if guessedLetters[unicode.ToLower(letter)] {
			guessingState += string(letter)
		} else {
			guessingState += "_"
		}
		guessingState += " "
	}
	return guessingState
}

func getUsedLetters(letters []string) string {
	result := ""
	for _, l := range letters {
		result += l + " "
	}
	return result
}

func printLayout(word string, guessedLetters map[rune]bool, hangmanState int, usedLetters []string, trys int, wins int, defeats int, hints int) {

	hangman := getHangman(hangmanState)
	firstLine := "  **             " + wordToGuessingState(word, guessedLetters)
	secondLine := "##############"
	cardinalsLine := "#################################################################################"

	t := "## Trys: " + strconv.Itoa(trys) + "               "
	w := "Wins: " + strconv.Itoa(wins) + "          "
	d := "Defeats: " + strconv.Itoa(defeats) + "                    "
	h := "Hints allowed: " + strconv.Itoa(hints) + " ##"

	fmt.Println()
	fmt.Println(hangman)
	fmt.Println(firstLine)
	fmt.Println(secondLine)
	fmt.Println("Used letters: " + getUsedLetters(usedLetters))
	fmt.Println(cardinalsLine)
	fmt.Println(t[0:24], w[0:11], d[0:24], h)
	fmt.Println(cardinalsLine)

}
