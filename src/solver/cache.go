package solver

import (
	"encoding/gob"
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
		valid := fd.Legal()
		if valid {
			solutions := FilterSolutions(oldSolutions, c.Word, fd)
			guess, _, err := ChooseGuess(guesses, solutions)

			valid = err == nil
			if valid {
				*c.NextLayer = append(*c.NextLayer, Cache{guess, nil})
				(*c.NextLayer)[i].Build(guesses, solutions, depth-1)
			}
		}

		if !valid {
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
		return err
	}

	g := gob.NewEncoder(f)
	if err = g.Encode(c); err != nil {
		return err
	}
	return nil
}

func ReadCache(path string) (Cache, error) {
	f, err := os.Open(path)
	if err != nil {
		return Cache{}, err
	}

	var buf Cache
	g := gob.NewDecoder(f)
	if err = g.Decode(&buf); err != nil {
		return Cache{}, err
	}
	return buf, nil
}
