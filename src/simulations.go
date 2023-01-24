package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func get_pattern(guess, hidden string) []byte {
	pattern := bytes.Repeat([]byte{0}, len(hidden))
	green_counter := make([]int, 26)

	for i := 0; i < len(hidden); i++ {
		if guess[i] == hidden[i] {
			pattern[i] = 'g'
			green_counter[guess[i]-'a']++
		} else if !strings.Contains(hidden, string(guess[i])) {
			pattern[i] = 'b'
		}
	}

	for i := 0; i < len(hidden); i++ {
		if pattern[i] != 0 {
			continue
		}

		if green_counter[guess[i]-'a'] != strings.Count(hidden, string(guess[i])) {
			pattern[i] = 'y'
			green_counter[guess[i]-'a']++
		} else {
			pattern[i] = 'b'
		}
	}

	return pattern
}

func read_wordlist(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("Couldn't open wordlist file at [%s]", path))
	}

	return strings.Split(string(data), "\n")
}

func play_single_word(wordlist []string, first_guess, hidden string) int {
	guesses := 1
	pattern := get_pattern(first_guess, hidden)
	candidates := get_candidates(wordlist, first_guess, pattern)

	for !bytes.Equal(pattern, bytes.Repeat([]byte{'g'}, len(pattern))) {
		guesses++
		next_guess := get_optimal_guess(candidates, wordlist)
		pattern = get_pattern(next_guess.word, hidden)

		fmt.Printf("\t[%d] [%s %f] %s\n", len(candidates), next_guess.word, next_guess.entropy, pattern)
		candidates = get_candidates(candidates, next_guess.word, pattern)
	}

	return guesses
}

func play_dictionary(path, first_guess string) float64 {
	mean := .0
	total_guesses := 0

	wordlist := read_wordlist(path)
	for i, word := range wordlist {
		fmt.Printf("Solving for %s\n", word)

		guesses := play_single_word(wordlist, first_guess, word)
		total_guesses += guesses
		mean = float64(total_guesses) / float64(i+1)

		fmt.Printf("[%d / %d] [%f]  %s: %d\n\n", i, len(wordlist), mean, word, guesses)
	}

	return mean
}
