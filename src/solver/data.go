package solver

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrInconsistentWords = errors.New("not all words were of the same length")
var ErrEmptyWords = errors.New("the words file was empty")

func ReadWords(path string) (Words, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadWords: failed to open file: %s", err.Error())
	}

	rawB, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("ReadWords: failed to read from file: %s", err.Error())
	}

	for len(rawB) > 0 && rawB[len(rawB)-1] == '\n' {
		rawB = rawB[:len(rawB)-1]
	}

	if len(rawB) == 0 {
		return nil, ErrEmptyWords
	}

	rawS := string(rawB)
	words := strings.Split(rawS, "\n")
	for _, word := range words {
		if len(word) != len(words[0]) {
			return nil, ErrInconsistentWords
		}
	}

	return words, nil
}

func WriteWords(wordlist Words, path string) error {
	// mainly for use by external packages
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("WriteWords: %s", err.Error())
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	for _, word := range wordlist {
		writer.WriteString(word)
		writer.WriteByte('\n')
	}

	return nil
}
