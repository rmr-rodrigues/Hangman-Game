package main

import (
	"bufio"
	"fmt"
	"os"
)

var dictionary = map[string]bool{}

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

	loadDictionary()
	for word := range dictionary {
		fmt.Println(word)
	}

}

func generateWordForGuessing() {

}

func loadDictionary() {
	readFile, err := os.Open(resourcesPath + dictionaryFileName)
	defer readFile.Close()

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
}
