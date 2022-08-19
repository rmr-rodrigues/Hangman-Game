# Hangman Game

Hangman written in Go (Golang) for the console. Just trying to learn Go!

Essentially this was a good exercise to learn how to use strings, runes, files and pointers in Go.

## Improvements to make

* Remove the global variables
* Use colors on the text messages
* Remove the "enter to continue" when the letter is repeated or invalid, just accept a new letter 
* Organize the variables in structs
* Divide the main function in several functions to better organize and present the code
* In the function 'selectTwoRandomLetters', if indexs or letters are equal generate different ones
* When loading the dictionary, filter words that are not according to a regex to prevent errors 
* Organize the dictionary by categories and present the category when user is trying to guess the word that way it would be more challenging, I think
* Create a more complete dictionary



