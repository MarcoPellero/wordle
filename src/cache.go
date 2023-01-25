package main

import (
	"bufio"
	"bytes"
	"math"
	"os"
	"strings"
)

func build_starting_cache(initial_guess string, wordlist []string, path string) {
	// this will try all possile patterns, for each one decide the best next guess, and save all of these to a file
	// since the first 1 or 2 guesses are EXTREMELY slow to find cause there's a lot of candidates, this will be a huge speedup
	pattern := make([]byte, len(initial_guess))
	for i := 0; i < len(pattern); i++ {
		pattern[i] = 'b'
	}

	precalc_logs := make([]float64, len(wordlist))
	for i := 0; i < len(precalc_logs); i++ {
		precalc_logs[i] = math.Log2(float64(i))
	}

	gy_counter := []int{0, 0}

	file, err := os.Create(path)
	if err != nil {
		panic("Couldn't open initial cache file")
	}
	defer file.Close()

	for !bytes.Equal(pattern, []byte(strings.Repeat("g", len(pattern)))) {
		// if all letters are green, and the last one is yellow, the pattern is illegal
		if gy_counter[0] != 4 || gy_counter[1] != 1 {
			candidates := get_candidates(wordlist, initial_guess, pattern)
			if len(candidates) != 0 {
				best_guess := get_optimal_guess(candidates, wordlist)
				file.Write([]byte(best_guess.word))
			}
		}

		file.WriteString("\n")

		for i := 0; i < len(pattern); i++ {
			if pattern[i] == 'g' {
				pattern[i] = 'b'
				gy_counter[0]--
				continue
			} else if pattern[i] == 'y' {
				pattern[i] = 'g'
				gy_counter[0]++
				gy_counter[1]--
			} else {
				pattern[i] = 'y'
				gy_counter[1]++
			}
			break
		}
	}
}

func read_from_cache(pattern []byte) string {
	file, err := os.Open("../data/cache1")
	if err != nil {
		panic("Couldn't open initial cache to read from")
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	idx := 0
	for i, x := range pattern {
		base := int(math.Pow(3, float64(i)))
		if x == 'b' {
			base = 0
		} else if x == 'g' {
			base *= 2
		}
		idx += base
	}

	for i := 0; i < idx; i++ {
		reader.ReadBytes('\n')
	}

	line, _ := reader.ReadBytes('\n')
	return string(line[:len(line)-1])
}
