package solver

import (
	"fmt"
)

func SimulateGame(guesses, solutions Words, c Cache, hidden string) (int, error) {
	var guess string
	var err error
	var cacheConsumed = false

	for i := 1; true; i++ {
		if !cacheConsumed {
			guess = c.Word
			cacheConsumed = true
			err = nil
		} else {
			guess, _, err = ChooseGuess(guesses, solutions)
		}

		if err != nil {
			return 0, fmt.Errorf("SimulateGame: %s", err.Error())
		}

		fd := GenerateFeedback(guess, hidden)
		if c.NextLayer != nil {
			cacheConsumed = false
			c = (*c.NextLayer)[fd.Hash()]
		}

		if fd.Won() {
			return i, nil
		}

		solutions = FilterSolutions(solutions, guess, fd)
	}

	panic("reached theoreticlaly unreachable code at SimulateGame")
}

func PlayAll(guesses, solutions Words, c Cache) float64 {
	sum := 0
	for _, word := range solutions {
		x, err := SimulateGame(guesses, solutions, c, word)
		if err != nil {
			panic(err)
		}
		sum += x
	}

	return float64(sum) / float64(len(solutions))
}
