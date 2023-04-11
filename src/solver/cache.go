package solver

import (
	"encoding/gob"
	"fmt"
	"os"
)

type Cache struct {
	Word      string   `json:"word"`
	NextLayer *[]Cache `json:"next"`
}

func (c *Cache) Build(guesses, oldSolutions Words, depth int) {
	if depth == 0 {
		return
	}

	fd := make(Feedback, len(c.Word))
	c.NextLayer = &[]Cache{}

	for i := 0; true; i++ {
		solutions := FilterSolutions(oldSolutions, c.Word, fd)

		guess, _, err := ChooseGuess(guesses, solutions)
		if err == nil {
			*c.NextLayer = append(*c.NextLayer, Cache{guess, nil})
			(*c.NextLayer)[i].Build(guesses, solutions, depth-1)
		} else {
			*c.NextLayer = append(*c.NextLayer, Cache{"", nil})
		}

		if fd.Next() {
			break
		}
	}
}

func (c *Cache) Dump(path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Cache.Dump: %s", err.Error())
	}

	g := gob.NewEncoder(f)
	if err = g.Encode(c); err != nil {
		return fmt.Errorf("Cache.Dump: %s", err.Error())
	}
	return nil
}

func ReadCache(path string) (Cache, error) {
	f, err := os.Open(path)
	if err != nil {
		return Cache{}, fmt.Errorf("ReadCache: %s", err.Error())
	}

	var buf Cache
	g := gob.NewDecoder(f)
	if err = g.Decode(&buf); err != nil {
		return Cache{}, fmt.Errorf("ReadCache: %s", err.Error())
	}
	return buf, nil
}
