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

func RateGuess(guesses Words, oldSolutions Words, guess string, solsLeft []int) float64 {
	for i := range solsLeft {
		solsLeft[i] = 0
	}

	for _, word := range oldSolutions {
		fd := GenFeedbackHash(guess, word)
		solsLeft[fd]++
	}

	info := 0.0
	for _, x := range solsLeft {
		info += entropyFormula(len(oldSolutions), x)
	}
	return info
}

func ChooseGuess(guesses, solutions Words) (string, float64, error) {
	if len(solutions) == 0 {
		return "", 0, ErrNoSolutions
	} else if len(solutions) == 1 {
		return solutions[0], 0, nil
	} else if len(solutions) <= 2 { // it's "mathematically" better to just do a coin flip on avg.
		return solutions[0], 1, nil
	}

	type Result struct {
		Idx    int
		Rating float64
	}

	jobs := make(chan int, len(guesses)+1)
	results := make(chan Result, len(guesses))
	for i := range guesses {
		jobs <- i
	}
	jobs <- -1

	solsLeft := make([][]int, 100)
	bufSize := int(math.Pow(3, float64(len(guesses[0]))))
	for i := range solsLeft {
		solsLeft[i] = make([]int, bufSize)
	}

	for i := 0; i < 100; i++ {
		go func(workerIdx int) {
			for {
				idx := <-jobs
				if idx == -1 {
					jobs <- -1
					return
				}

				rating := RateGuess(guesses, solutions, guesses[idx], solsLeft[workerIdx])
				results <- Result{idx, rating}
			}
		}(i)
	}

	best := Result{0, -1}
	for i := 0; i < len(guesses); i++ {
		x := <-results
		if x.Rating > best.Rating || (x.Rating == best.Rating && x.Idx < best.Idx) {
			best = x
		}
	}

	return guesses[best.Idx], best.Rating, nil
}
