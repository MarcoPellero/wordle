package main

import (
	"bytes"
	"math"
	"strings"
)

type Guess struct {
	word    string
	entropy float64
}

func is_candidate(word, guess string, pattern []byte) bool {
	letter_counter := make([]int, 26)

	for i := 0; i < len(word); i++ {
		if pattern[i] != 'g' {
			continue
		} else if guess[i] != word[i] {
			return false
		}

		letter_counter[guess[i]-'a']++
	}

	for i := 0; i < len(word); i++ {
		if pattern[i] != 'y' {
			continue
		} else if guess[i] == word[i] {
			return false
		}

		if strings.Count(word, string(guess[i])) <= letter_counter[guess[i]-'a'] {
			return false
		}
		letter_counter[guess[i]-'a']++
	}

	for i := 0; i < len(word); i++ {
		if pattern[i] != 'b' {
			continue
		} else if guess[i] == word[i] {
			return false
		}

		for j := i + 1; j < len(word); j++ {
			if pattern[j] == 'y' && guess[i] == guess[j] {
				return false
			}
		}

		if strings.Count(word, string(guess[i])) > letter_counter[guess[i]-'a'] {
			return false
		}
	}

	return true
}

func get_candidates(wordlist []string, guess string, pattern []byte) []string {
	candidates := make([]string, 0)

	for _, word := range wordlist {
		if is_candidate(word, guess, pattern) {
			candidates = append(candidates, word)
		}
	}

	return candidates
}

func calculate_pattern_entropy(wordlist []string, guess string, pattern []byte) float64 {
	num_of_candidates := 0
	for _, word := range wordlist {
		if is_candidate(word, guess, pattern) {
			num_of_candidates++
		}
	}

	if num_of_candidates == 0 {
		return 0
	}

	px := float64(num_of_candidates) / float64(len(wordlist))
	return -px * math.Log2(px)
}

func calculate_guess_entropy(wordlist []string, guess string) float64 {
	entropy := .0
	color_counter := []int{0, 0} // [G, Y]
	pattern := bytes.Repeat([]byte{'b'}, len(guess))

	is_last := false
	for !is_last {
		// if all letters are green, except for a yellow one, the pattern is illegal and should not be considered
		if color_counter[0] != len(guess)-1 || color_counter[1] != 1 {
			entropy += calculate_pattern_entropy(wordlist, guess, pattern)
		}

		// go to next pattern
		for i := 0; i <= len(pattern); i++ {
			if i == len(pattern) {
				is_last = true
				break
			}

			if pattern[i] == 'b' {
				pattern[i] = 'y'
				color_counter[1]++
			} else if pattern[i] == 'y' {
				pattern[i] = 'g'
				color_counter[1]--
				color_counter[0]++
			} else {
				color_counter[0]--
				pattern[i] = 'b'
				continue
			}
			break
		}
	}

	return entropy
}

func _sync_get_optimal_guess(candidates, wordlist []string) Guess {
	if len(candidates) == 1 {
		return Guess{candidates[0], 0}
	}

	best := Guess{candidates[0], -1}
	for i := 0; i < len(wordlist); i++ {
		entropy := calculate_guess_entropy(candidates, wordlist[i])
		if entropy > best.entropy {
			best = Guess{wordlist[i], entropy}
		}
	}

	return best
}

func get_optimal_guess(candidates, wordlist []string) Guess {
	type Result struct {
		id      int
		entropy float64
	}

	if len(candidates) == 1 {
		return Guess{candidates[0], 0}
	}

	jobs := make(chan int, len(wordlist))
	results := make(chan Result, len(wordlist))
	num_of_workers := 20

	for i := 0; i < num_of_workers; i++ {
		go func() {
			for id := range jobs {
				entropy := calculate_guess_entropy(candidates, wordlist[id])
				results <- Result{id, entropy}
			}
		}()
	}

	for i := 0; i < len(wordlist); i++ {
		jobs <- i
	}

	best := Guess{candidates[0], -1}
	best_id := 0
	for i := 0; i < len(wordlist); i++ {
		res := <-results

		// the reason for taking the one with the lower idx (in wordlist) if 2 guesses have the same entropy
		// is to make sure that the algorithm is deterministic. without this check, the guesses vary,
		// and the avg. guesses in a simulation does aswell, although it does so by a very small margin (< .001, or 0.03%)
		if res.entropy > best.entropy || (res.entropy == best.entropy && res.id < best_id) {
			best = Guess{wordlist[res.id], res.entropy}
			best_id = res.id
		}
	}

	return best
}
