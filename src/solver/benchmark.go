package solver

import (
	"fmt"
)

func SimulateGame(guesses, solutions Words, c Cache, hidden string) (int, error) {
	var guess string
	var err error = nil
	var cacheConsumed = false

	for i := 1; true; i++ {
		if !cacheConsumed {
			guess = c.Word
			cacheConsumed = true
		} else if guess, _, err = ChooseGuess(guesses, solutions); err != nil {
			return 0, fmt.Errorf("SimulateGame: %s", err.Error())
		}

		fd := GenerateFeedback(guess, hidden)
		if fd.Won() {
			return i, nil
		}

		if c.NextLayer != nil {
			cacheConsumed = false
			c = (*c.NextLayer)[fd.Hash()]
		}
		solutions = FilterSolutions(solutions, guess, fd)
	}

	panic("reached theoreticlaly unreachable code at SimulateGame")
}

func PlayAll(guesses, solutions Words, c Cache) float64 {
	sum := 0
	results := make(chan int)

	for _, word := range solutions {
		go func(word string) {
			x, err := SimulateGame(guesses, solutions, c, word)
			if err != nil {
				panic(err)
			}
			results <- x
		}(word)
	}

	for range solutions {
		sum += <-results
	}

	return float64(sum) / float64(len(solutions))
}
