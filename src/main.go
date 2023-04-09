package main

import (
	"fmt"

	"github.com/MarcoPellero/wordle-gopher/src/solver"
)

func interactive(guesses, solutions solver.Words) int {
	wordLen := 5

	var buf int
	fmt.Printf("Enter word length, or nothing for %d: ", wordLen)
	if n, _ := fmt.Scanln(&buf); n > 0 {
		wordLen = buf
	}

	var guess, fdBuf string
	for i := 1; true; i++ {
		for {
			fmt.Printf("(%d) Guess: ", i)
			fmt.Scanln(&guess)
			if len(guess) != wordLen {
				fmt.Printf("Invalid length %d; the word length is %d\n", len(fdBuf), wordLen)
				continue
			}

			valid := true
			for j, c := range guess {
				if c < 'a' || c > 'z' {
					fmt.Printf("Invalid symbol '%c' at index %d; please only use letters\n", c, j)
					valid = false
					break
				}
			}

			if valid {
				break
			}
		}

		var fd solver.Feedback
		for {
			fmt.Printf("(%d) Feedback: ", i)
			fmt.Scanln(&fdBuf)
			if len(fdBuf) != wordLen {
				fmt.Printf("Invalid length %d; the word length is %d\n", len(fdBuf), wordLen)
				continue
			}

			var err error
			if fd, err = solver.FdFromString(fdBuf); err != nil {
				fmt.Println(err.Error())
				continue
			}

			break
		}

		if fd.Won() {
			fmt.Printf("Congratulations! You won in %d guesses!", i)
			return i
		}

		solutions = solver.FilterSolutions(solutions, guess, fd)
		fmt.Printf("(%d) There's %d possible solutions left\n", i, len(solutions))

		word, rating, err := solver.ChooseGuess(guesses, solutions)
		if err != nil {
			panic(err)
		}
		fmt.Printf("(%d) The best guess is %s, with a rating of %.3f\n", i, word, rating)
	}

	panic("reached theoreticlaly unreachable code at main.interactive")
}

func main() {
	guessesPath := "../data/wordlists/guesses.txt"
	answersPath := "../data/wordlists/answers.txt"
	cachePath := "../data/caches/main"

	guesses, err := solver.ReadWords(guessesPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d guesses\n", len(guesses))

	answers, err := solver.ReadWords(answersPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read %d answers\n", len(guesses))

	cache, err := solver.ReadCache(cachePath)
	if err != nil {
		fmt.Println("Generating cache")
		cache.Word = "sarti"
		cache.Build(guesses, answers, 3)
		if err := cache.Dump(cachePath); err != nil {
			panic(err)
		}
	}

	// interactive(gWords, gWords)

	fmt.Println("Simulating games")
	mean := solver.PlayAll(guesses, answers, cache)
	fmt.Printf("Solved any one word in on avg %.2f guesses\n", mean)
}
