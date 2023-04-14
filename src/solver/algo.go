package solver

import (
	"errors"
	"math"
)

var ErrNoSolutions = errors.New("there are no possible solutions")

func FilterSolutions(solutions Words, guess string, fd Feedback) Words {
	if fd.Won() {
		return Words{}
	}

	filtered := make(Words, 0)
	for _, word := range solutions {
		if word != guess && fd.Match(guess, word) {
			filtered = append(filtered, word)
		}
	}

	return filtered
}

func entropyFormula(oldSolutions, newSolutions int) float64 {
	if newSolutions == 0 {
		return 0
	}

	px := float64(newSolutions) / float64(oldSolutions)
	return -px * math.Log2(px)
}

func RateGuess(guesses Words, oldSolutions Words, guess string) float64 {
	solutionsLeft := make([]int, int(math.Pow(3, float64(len(guess)))))
	for _, word := range oldSolutions {
		fd := GenerateFeedback(guess, word)
		solutionsLeft[fd.Hash()]++
	}

	info := 0.0
	for _, x := range solutionsLeft {
		info += entropyFormula(len(oldSolutions), x)
	}
	return info
}

func ChooseGuess(guesses, solutions Words) (string, float64, error) {
	if len(solutions) == 0 {
		return "", 0.0, ErrNoSolutions
	}
	if len(solutions) == 1 {
		return solutions[0], 0.0, nil
	}
	if len(solutions) == 2 { // coin flip!
		return solutions[0], 1, nil
	}

	type result struct {
		Idx    int
		Rating float64
	}

	results := make(chan result)
	for i := range guesses {
		go func(idx int) {
			results <- result{idx, RateGuess(guesses, solutions, guesses[idx])}
		}(i)
	}

	best := result{0, -1}
	for i := 0; i < len(guesses); i++ {
		x := <-results
		if x.Rating > best.Rating || (x.Rating == best.Rating && x.Idx < best.Idx) {
			best = x
		}
	}

	return guesses[best.Idx], best.Rating, nil
}
