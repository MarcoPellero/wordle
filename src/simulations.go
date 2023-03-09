package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
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

func play_single_word(solutions, guesses []string, cache func([]byte) string, first_guess, hidden string) int {
	guess := Guess{first_guess, -1}
	var pattern []byte
	candidates := make([]string, len(solutions))
	copy(candidates, solutions)

	for i := 1; true; i++ {
		pattern = get_pattern(guess.word, hidden)
		if bytes.Equal(pattern, bytes.Repeat([]byte{'g'}, len(pattern))) {
			return i
		}

		candidates = get_candidates(candidates, guess.word, pattern)
		if i == 1 {
			guess.word = cache(pattern)
		}

		if i != 1 || len(guess.word) == 0 {
			guess, _ = get_optimal_guess(candidates, guesses)
		}
	}

	panic("Unreachable return statement")
}

func interactive_game(solutions, guesses []string, cache func([]byte) string) {
	candidates := make([]string, len(solutions))
	copy(candidates, solutions)

	var guess string
	var pattern []byte
	var n_guesses int
	for n_guesses = 0; ; n_guesses++ {
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
			n_guesses += 2
			break
		}

		suggestion, _ := get_optimal_guess(candidates, guesses)
		fmt.Printf("You should guess %s\n", suggestion.word)
	}

	fmt.Printf("Success! solved in %d guesses\n", n_guesses)
}

func play_dictionary(solutions, guesses []string, cache func([]byte) string, first_guess string) float64 {
	total_guesses := 0

	for i, word := range solutions {
		guesses := play_single_word(solutions, guesses, cache, first_guess, word)
		total_guesses += guesses
		fmt.Printf("[%d / %d] %s: %d\n", i, len(solutions), word, guesses)
	}

	return float64(total_guesses) / float64(len(solutions))
}

func main() {
	first_guess := "sarti"
	solutions_path := "../wordlists/guesses"
	guesses_path := "../wordlists/guesses"

	solutions := read_wordlist(solutions_path)
	guesses := read_wordlist(guesses_path)

	cache_path := "../data/cache1"
	if _, err := os.Stat(cache_path); err != nil {
		fmt.Println("Generating cache, wait...")
		store_cache(solutions, guesses, cache_path, first_guess)
	}
	cache := build_cache(cache_path)

	if len(os.Args) <= 1 {
		fmt.Println("You need to invoke this command with an extra parameter! (-dictionary or -interactive)")
	} else if os.Args[1] == "-dictionary" {
		if len(os.Args) >= 3 {
			first_guess = os.Args[2]
		}

		num_of_runs := 1
		if len(os.Args) >= 4 {
			num_of_runs, _ = strconv.Atoi(os.Args[3])
		}

		for i := 0; i < num_of_runs; i++ {
			mean := play_dictionary(solutions, guesses, cache, first_guess)
			fmt.Printf("[%d] Solved words in, on average, %f guesses\n", i, mean)
		}
	} else if os.Args[1] == "-interactive" {
		interactive_game(solutions, guesses, cache)
	} else if os.Args[1] == "-api" {
		bot_server(solutions_path, guesses_path, cache_path)
	} else if os.Args[1] == "-filter-wordlist" {
		filter_wordlist_server(solutions_path)
	} else {
		fmt.Println("Invalid parameter! Use dictionary or interactive")
	}
}
