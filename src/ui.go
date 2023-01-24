package main

import "fmt"

func main() {
	mean := play_dictionary("./../wordlists/small", "sarti")
	fmt.Printf("Avg. guesses per word: %f\n", mean)
}
