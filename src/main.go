package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/MarcoPellero/wordle/src/solver"
	"golang.org/x/exp/slices"
)

func interactive(guesses, solutions solver.Words) int {
	wordLen := len(guesses[0])
	var guess, fdBuf string
	for i := 1; true; i++ {
		for {
			fmt.Printf("(%d) Guess: ", i)
			fmt.Scanln(&guess)
			if len(guess) != wordLen {
				fmt.Printf("Invalid length %d; the word length is %d\n", len(fdBuf), wordLen)
				continue
			}

			valid := true
			for j, c := range guess {
				if c < 'a' || c > 'z' {
					fmt.Printf("Invalid symbol '%c' at index %d; please only use letters\n", c, j)
					valid = false
					break
				}
			}

			if valid {
				break
			}
		}

		var fd solver.Feedback
		for {
			fmt.Printf("(%d) Feedback: ", i)
			fmt.Scanln(&fdBuf)
			if len(fdBuf) != wordLen {
				fmt.Printf("Invalid length %d; the word length is %d\n", len(fdBuf), wordLen)
				continue
			}

			var err error
			if fd, err = solver.FdFromString(fdBuf); err != nil {
				fmt.Println(err.Error())
				continue
			}

			break
		}

		if fd.Won() {
			fmt.Printf("Congratulations! You won in %d guesses!", i)
			return i
		}

		solutions = solver.FilterSolutions(solutions, guess, fd)
		fmt.Printf("(%d) There's %d possible solutions left\n", i, len(solutions))

		word, rating, err := solver.ChooseGuess(guesses, solutions)
		if err != nil {
			panic(err)
		}
		fmt.Printf("(%d) The best guess is %s, with a rating of %.3f\n", i, word, rating)
	}

	panic("reached theoreticlaly unreachable code at main.interactive")
}

func loadBalancer(port int, guesses solver.Words, outFile string) []float64 {
	// i'm assuming all of the remote workers are using the same wordlists..
	results := make([]float64, len(guesses))
	running := 0
	idx := 0

	lock := sync.Mutex{}

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		lock.Lock()
		defer lock.Unlock()
		defer r.Body.Close()

		var job int
		if idx != len(guesses) {
			job = idx
			idx++
		}

		w.Write([]byte(fmt.Sprint(job)))
		running++
	})

	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("/result: %s\n", err.Error())
		}

		var job int
		var result float64
		_, err = fmt.Sscanf(string(buf), "%d %f", &job, &result)
		if err != nil {
			fmt.Printf("/result: %s\n", err.Error())
			return
		}

		results[job] = result
		running--

		fmt.Printf("/result [%s | %.3f]\n", guesses[job], result)
	})

	go http.ListenAndServe("0.0.0.0:8080", nil)
	for idx < len(guesses) || running > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("loadBalancer: %s\n", err.Error())
	}

	for i := range results {
		fmt.Fprintf(f, "%d %s %f", i, guesses[i], results[i])
	}

	return results
}

func remoteSlave(endpoint string, guesses, solutions solver.Words) {
	for {
		res, err := http.Get(endpoint + "/work")
		if err != nil {
			fmt.Printf("/work GET: %s\n", err.Error())
			time.Sleep(50 * time.Millisecond)
			continue
		}

		buf, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("/work read: %s\n", err.Error())
			time.Sleep(50 * time.Millisecond)
			continue
		}

		var job int
		_, err = fmt.Sscanf(string(buf), "%d", &job)
		if err != nil {
			fmt.Printf("/work: %s", err.Error())
			time.Sleep(50 * time.Millisecond)
			continue
		}
		fmt.Printf("Working on [%d] | %s\n", job, guesses[job])

		var cache solver.Cache
		cache.Word = guesses[job]
		cache.Build(guesses, solutions, 3)
		rating := solver.PlayAll(guesses, solutions, cache)
		fmt.Printf("Rating: %.3f\n", rating)

		bodyData := fmt.Sprintf("%d %f", job, rating)
		http.Post(endpoint+"/result", "text/plain", bytes.NewBufferString(bodyData))
	}
}

func main() {
	guesses, _ := solver.ReadWords("../data/wordlists/length6.txt")

	if slices.Contains(os.Args, "serve") {
		var port int
		fmt.Sscanf(os.Args[len(os.Args)-1], "%d", &port)
		outFile := fmt.Sprintf("../data/results/out_%d.txt", len(guesses[0]))

		fmt.Printf("Running load balancer on port %d\n", port)
		results := loadBalancer(8080, guesses, outFile)

		best := 0
		for i := range results {
			if results[i] > results[best] {
				best = i
			}
		}

		fmt.Printf("Best: %d %s | %.3f\n", best, guesses[best], results[best])
	} else if slices.Contains(os.Args, "slave") {
		url := os.Args[len(os.Args)-1]
		fmt.Printf("Running slave for %s\n", url)
		remoteSlave(url, guesses, guesses)
	} else if slices.Contains(os.Args, "simulate") {
		path := fmt.Sprintf("../data/caches/%d", len(guesses[0]))
		c, err := solver.ReadCache(path)
		if err != nil {
			fmt.Printf("Generating cache")
			c.Word = os.Args[len(os.Args)-1]
			c.Build(guesses, guesses, 3)
			c.Dump(path)
		}

		fmt.Println("Running simulation")
		mean := solver.PlayAll(guesses, guesses, c)
		fmt.Printf("Solved any one word in on avg %.2f guesses\n", mean)
	} else {
		fmt.Printf("Unknown command")
	}
}
