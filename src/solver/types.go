package solver

import (
	"fmt"
)

// letter state enums
const (
	Black uint8 = iota
	Yellow
	Green
)

type CharState = uint8
type Feedback []CharState
type Words = []string

func (fd *Feedback) Next() bool {
	var i int
	for i = 0; i < len(*fd); i++ {
		if (*fd)[i] != Green {
			(*fd)[i]++
			break
		}
		(*fd)[i] = Black
	}

	return i == len(*fd)
}

func (fd *Feedback) Won() bool {
	for _, c := range *fd {
		if c != Green {
			return false
		}
	}

	return true
}

func (fd *Feedback) String() string {
	buf := make([]byte, len(*fd))
	for i, c := range *fd {
		switch c {
		case Black:
			buf[i] = 'B'
		case Yellow:
			buf[i] = 'Y'
		case Green:
			buf[i] = 'G'
		}
	}

	return string(buf)
}

func (fd *Feedback) Hash() int {
	total := 0
	mul := 1
	for _, x := range *fd {
		total += mul * int(x)
		mul *= 3
	}

	return total
}

func (fd *Feedback) Match(guess, solution string) bool {
	solAlpha := make([]uint8, 26)
	for i := range solution {
		solAlpha[solution[i]-'a']++
	}

	for i := range solution {
		if guess[i] != solution[i] {
			continue
		}

		if (*fd)[i] != Green {
			return false
		}

		solAlpha[guess[i]-'a']--
	}

	for i := range solution {
		if guess[i] == solution[i] {
			continue
		}

		if solAlpha[guess[i]-'a'] > 0 {
			if (*fd)[i] != Yellow {
				return false
			}
			solAlpha[guess[i]-'a']--
		} else if (*fd)[i] != Black {
			return false
		}
	}

	return true
}

func FdFromString(buf string) (Feedback, error) {
	fd := make(Feedback, len(buf))
	for i, c := range buf {
		switch c {
		case 'b':
			fd[i] = Black
		case 'y':
			fd[i] = Yellow
		case 'g':
			fd[i] = Green
		default:
			return nil, fmt.Errorf("FdFromString: Invalid symbol '%c' at index %d", c, i)
		}
	}

	return fd, nil
}

func GenFeedback(guess, solution string) Feedback {
	solAlpha := make([]uint8, 26)
	for i := range solution {
		solAlpha[solution[i]-'a']++
	}

	fd := make(Feedback, len(solution))
	for i := range solution {
		if guess[i] == solution[i] {
			fd[i] = Green
			solAlpha[guess[i]-'a']--
		}
	}

	for i := range solution {
		if fd[i] == Green {
			continue
		}

		if solAlpha[guess[i]-'a'] > 0 {
			fd[i] = Yellow
			solAlpha[guess[i]-'a']--
		} else {
			fd[i] = Black
		}
	}

	return fd
}

func GenFeedbackHash(guess, solution string) int {
	solAlpha := make([]uint8, 26)
	for i := range solution {
		solAlpha[solution[i]-'a']++
	}

	fd := 0
	mul := 1
	for i := range solution {
		if guess[i] == solution[i] {
			fd += mul * 2
			solAlpha[guess[i]-'a']--
		}
		mul *= 3
	}

	mul = 1
	for i := range solution {
		if guess[i] == solution[i] {
			mul *= 3
			continue
		}

		if solAlpha[guess[i]-'a'] > 0 {
			fd += mul
			solAlpha[guess[i]-'a']--
		}
		mul *= 3
	}

	return fd
}
