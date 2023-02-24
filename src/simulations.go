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

func play_single_word(wordlist []string, cache func([]byte) string, first_guess, hidden string) int {
	guess := Guess{first_guess, -1}
	var pattern []byte
	candidates := make([]string, len(wordlist))
	copy(candidates, wordlist)

	for i := 1; true; i++ {
		pattern = get_pattern(guess.word, hidden)
		if bytes.Equal(pattern, bytes.Repeat([]byte{'g'}, len(pattern))) {
			return i
		}

		candidates = get_candidates(candidates, guess.word, pattern)
		if i == 1 {
			guess.word = cache(pattern)
		} else {
			guess = get_optimal_guess(candidates, wordlist)
		}
	}

	panic("Unreachable return statement")
}

func interactive_game(wordlist []string, cache func([]byte) string) {
	candidates := make([]string, len(wordlist))
	copy(candidates, wordlist)

	var guess string
	var pattern []byte
	var guesses int
	for guesses = 0; ; guesses++ {
		fmt.Print("Guess: ")
		fmt.Scanln(&guess)
		fmt.Print("Pattern: ")
		fmt.Scanln(&pattern)

		if bytes.Equal(pattern, bytes.Repeat([]byte{'g'}, len(pattern))) {
			break
		}

		candidates = get_candidates(candidates, guess, pattern)
		if len(candidates) == 1 {
			fmt.Printf("The hidden word is %s!\n", candidates[0])
			guesses += 2
			break
		}

		suggestion := get_optimal_guess(candidates, wordlist)
		fmt.Printf("You should guess %s\n", suggestion.word)
	}

	fmt.Printf("Success! solved in %d guesses\n", guesses)
}

func play_dictionary(wordlist []string, cache func([]byte) string, first_guess string) float64 {
	total_guesses := 0

	for i, word := range wordlist {
		guesses := play_single_word(wordlist, cache, first_guess, word)
		total_guesses += guesses
		fmt.Printf("[%d / %d] %s: %d\n", i, len(wordlist), word, guesses)
	}

	return float64(total_guesses) / float64(len(wordlist))
}

func main() {
	first_guess := "sarti"

	wordlist_path := "../wordlists/adaptive"
	wordlist := read_wordlist(wordlist_path)

	cache_path := "../data/cache1"
	if _, err := os.Stat(cache_path); err != nil {
		fmt.Println("Generating cache, wait...")
		store_cache(wordlist, cache_path, first_guess)
	}
	cache := build_cache(cache_path)

	if len(os.Args) <= 1 {
		fmt.Println("You need to invoke this command with an extra parameter! (-dictionary or -interactive)")
	} else if os.Args[1] == "-dictionary" {
		if len(os.Args) >= 3 {
			first_guess = os.Args[2]
		}

		mean := play_dictionary(wordlist, cache, first_guess)
		fmt.Printf("Solved words in, on average, %f guesses\n", mean)
	} else if os.Args[1] == "-interactive" {
		interactive_game(wordlist, cache)
	} else if os.Args[1] == "-api" {
		bot_server(wordlist_path)
	} else if os.Args[1] == "-filter-wordlist" {
		filter_wordlist_server(wordlist_path)
	} else {
		fmt.Println("Invalid parameter! Use dictionary or interactive")
	}
}
