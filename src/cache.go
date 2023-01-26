package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

// FORMAT [2 bytes: word length][2 bytes: number of words][word 1][word 2][word 3]

func store_cache(wordlist []string, path, first_guess string) {
	file, err := os.Create(path)
	if err != nil {
		panic(fmt.Sprintf("Couldn't create cache file at %s", path))
	}
	defer file.Close()

	buf_uint16 := make([]byte, 2)
	word_len := len(wordlist[0])
	binary.LittleEndian.PutUint16(buf_uint16, uint16(word_len))
	file.Write(buf_uint16)
	binary.LittleEndian.PutUint16(buf_uint16, uint16(len(wordlist)))
	file.Write(buf_uint16)

	pattern := bytes.Repeat([]byte{'b'}, word_len)

	color_counter := []int{0, 0}
	for color_counter[0] != word_len {
		is_valid := false
		if color_counter[0] != word_len-1 || color_counter[1] != 1 {
			candidates := get_candidates(wordlist, first_guess, pattern)
			if len(candidates) != 0 {
				is_valid = true
				optimal_guess := get_optimal_guess(candidates, wordlist)
				file.WriteString(optimal_guess.word)
			}
		}

		if !is_valid {
			file.Write([]byte{0})
		}

		for i := 0; i < len(pattern); i++ {
			if pattern[i] == 'g' {
				pattern[i] = 'b'
				color_counter[0]--
				continue
			} else if pattern[i] == 'y' {
				pattern[i] = 'g'
				color_counter[0]++
				color_counter[1]--
			} else {
				pattern[i] = 'y'
				color_counter[1]++
			}
			break
		}
	}
}

func build_cache(path string) func([]byte) string {
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("Couldn't open cache file at %s", path))
	}

	buf_int16 := make([]byte, 2)
	file.Read(buf_int16)
	word_len := int(binary.LittleEndian.Uint16(buf_int16))
	file.Read(buf_int16)
	wordlist_len := int(binary.LittleEndian.Uint16(buf_int16))

	cache := make([]string, wordlist_len)
	word_buf := make([]byte, word_len)
	letter_buf := make([]byte, 1)
	offset := int64(len(buf_int16) * 2)

	for i := 0; i < wordlist_len; i++ {
		if file.ReadAt(letter_buf, int64(offset)); letter_buf[0] == 0 {
			offset++
			continue
		}
		file.ReadAt(word_buf, int64(offset))
		cache[i] = string(word_buf)
		offset += int64(word_len)
	}

	return func(pattern []byte) string {
		idx := 0
		base := 1
		for _, x := range pattern {
			if x == 'y' {
				idx += base
			} else if x == 'g' {
				idx += base * 2
			}
			base *= 3
		}

		return cache[idx]
	}
}
