package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var dictionary = map[string]bool{}
var dictionaryKeys []string
var inputReader = bufio.NewReader(os.Stdin)

const resourcesPath = "resources/"
const dictionaryFileName = "dictionary.txt"

func main() {

	// TODO remove the global variables
	// TODO use colors on the messages text
	// TODO remove the "enter to continue" when the letter is repeated or invalid, just accept a new letter
	// TODO organize the data in structs
	// TODO divide the main function in several functions to better organize and present the code
	// TODO in the function 'selectTwoRandomLetters', if indexs or letters are equal generate different ones
	// TODO in the function 'getDictionary', filter words that are not according to a regex to prevent errors
	// TODO organize the dictionary by categories and present the category when user is trying to guess the word that way it would be more challenging, I think!

	loadDictionary()
	getDictionaryKeys()
	var trys, wins, defeats, hints, hangmanState int
	hints = 1

	for {
		word := getRandomWord()

		guessedLetters := make(map[rune]bool)
		usedLetters := make([]string, 0, 10)
		selectTwoRandomLetters(&word, guessedLetters, &usedLetters)
		msg := ""
		hangmanState = 0
		trys = 0
		hints = 1

		for {
			if hasGuessedAllLetters(&word, guessedLetters) { // You Win

				msg = "You Win! (enter to continue) "
				wins++
				printLayout(msg, word, guessedLetters, hangmanState, usedLetters, trys, wins, defeats, hints)
				readInput()
				break
			} else if hangmanState == 9 { //You lost
				msg = "You lost! (enter to continue) "
				defeats++
				printLayout(msg, word, guessedLetters, hangmanState, usedLetters, trys, wins, defeats, hints)
				readInput()
				break
			} else {
				msg = ""
				printLayout(msg, word, guessedLetters, hangmanState, usedLetters, trys, wins, defeats, hints)
				input := readInput()
				if !validateInput(input) { // not a valid input
					msg = "Invalid input! (enter to continue) "
					printLayout(msg, word, guessedLetters, hangmanState, usedLetters, trys, wins, defeats, hints)
					readInput()
				} else { // Valid input
					if input == "?" && hints > 0 {
						hint := getHint(&word, guessedLetters, &usedLetters, &hints)
						if hint != "" {
							msg = fmt.Sprintf("Hint: %s (enter to continue) ", hint)
							printLayout(msg, word, guessedLetters, hangmanState, usedLetters, trys, wins, defeats, hints)
							readInput()
						} else {
							msg = "It's not possible to give you an hint!"
							printLayout(msg, word, guessedLetters, hangmanState, usedLetters, trys, wins, defeats, hints)
							readInput()
						}
					} else if input == "?" && hints == 0 { // No more hints
						msg = "Sorry, no more hints! (enter to continue) "
						printLayout(msg, word, guessedLetters, hangmanState, usedLetters, trys, wins, defeats, hints)
						readInput()
					} else if input == "0" {
						os.Exit(0)
					} else if isNewGuess(&word, &input, guessedLetters) { // New correct letter
						guessedLetters[rune(input[0])] = true
						usedLetters = append(usedLetters, input)
						trys++
					} else if isRepeatedLetter(&input, &usedLetters) { // Repeated letter, just try again
						msg = "Repeated letter! (enter to continue) "
						printLayout(msg, word, guessedLetters, hangmanState, usedLetters, trys, wins, defeats, hints)
						readInput()
					} else { // Failed letter
						msg = ""
						hangmanState++
						usedLetters = append(usedLetters, input)
						trys++
					}
				}
			}
		}
	}
}

// It's easier to random select a word from a slice than from a map keys!
func getDictionaryKeys() {
	dictionaryKeys = make([]string, len(dictionary))
	i := 0
	for word := range dictionary {
		dictionaryKeys[i] = word
		i++
	}
}

func selectTwoRandomLetters(word *string, guessedLetters map[rune]bool, usedLetters *[]string) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	f := random.Intn(len(*word))
	s := random.Intn(len(*word))

	if f != s {
		guessedLetters[unicode.ToLower(rune((*word)[f]))] = true
		guessedLetters[unicode.ToLower(rune((*word)[s]))] = true
		*usedLetters = append(*usedLetters, string(rune((*word)[f])))
		*usedLetters = append(*usedLetters, string(rune((*word)[s])))
	} else {
		guessedLetters[unicode.ToLower(rune((*word)[f]))] = true
		*usedLetters = append(*usedLetters, string(rune((*word)[f])))
	}
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

func wordToGuessingState(word *string, guessedLetters map[rune]bool) string {

	guessingState := ""
	for _, letter := range *word {
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

func getUsedLetters(letters *[]string) string {
	result := ""
	for _, l := range *letters {
		result += l + " "
	}
	return result
}

func printLayout(msg string, word string, guessedLetters map[rune]bool, hangmanState int, usedLetters []string, trys int, wins int, defeats int, hints int) {

	if msg == "" {
		msg = "Type character (? for hint or 0 to quit): "
	}
	hangman := getHangman(hangmanState)
	firstLine := "  **             " + wordToGuessingState(&word, guessedLetters)
	secondLine := "##############"
	cardinalsLine := "#################################################################################"

	t := "## Trys: " + strconv.Itoa(trys) + "               "
	w := "Wins: " + strconv.Itoa(wins) + "          "
	d := "Defeats: " + strconv.Itoa(defeats) + "                    "
	h := "Hints allowed: " + strconv.Itoa(hints) + " ##"

	clearConsole()
	fmt.Println()
	fmt.Println(hangman)
	fmt.Println(firstLine)
	fmt.Println(secondLine)
	fmt.Println("Used letters: ", getUsedLetters(&usedLetters))
	fmt.Println(cardinalsLine)
	fmt.Println(t[0:24], w[0:11], d[0:24], h)
	fmt.Println(cardinalsLine)
	fmt.Print(msg)
}

func readInput() string {

	input, err := inputReader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(strings.ToLower(input))
}

func validateInput(input string) bool {
	validate, _ := regexp.MatchString(`^[A-Za-z?0]$`, input)

	return validate
}

func clearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println("Clearing the console doesn't work in this OS!")
	}
}

func isNewGuess(word *string, ch *string, guessedLetters map[rune]bool) bool {
	if strings.Contains(strings.ToLower(*word), strings.ToLower(*ch)) {
		if guessedLetters[unicode.ToLower(rune((*ch)[0]))] {
			return false
		} else {
			return true
		}
	}
	return false
}

func isRepeatedLetter(ch *string, usedLetters *[]string) bool {
	for _, l := range *usedLetters {
		if l == *ch {
			return true
		}
	}
	return false
}

func getHint(word *string, guessedLetters map[rune]bool, usedLetters *[]string, hints *int) string {
	for _, l := range *word {
		if !guessedLetters[unicode.ToLower(l)] {
			guessedLetters[unicode.ToLower(l)] = true
			*usedLetters = append(*usedLetters, string(unicode.ToLower(l)))
			*hints--
			return string(unicode.ToLower(l))
		}
	}
	return ""
}

func hasGuessedAllLetters(word *string, guessedLetters map[rune]bool) bool {
	for _, l := range *word {
		if !guessedLetters[unicode.ToLower(l)] {
			return false
		}
	}
	return true
}
