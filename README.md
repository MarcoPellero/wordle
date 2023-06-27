# Wordle Algorithm
This is a faster implementation of 3b1b's algorithm for playing wordle optimally, using information theory.

It's made both for simulations (like 3b1b's), but also for interactive use, mainly programmatically. It's usable as a Go package directly, but it had (and might have again) an HTTP api for use in programs written in other languages aswell.

Below are short descriptions on how to use it, if you have doubts you can open an Issue or contact me via Discord at @marco_pellero

## How to use it
### Rust implementation:
This one's waaay faster, i had fun optimizing it but there's still more to be done, i think it can be squeezed to be at least 30-70% faster, but i've barely started using Rust.

Just clone it:
```sh
git clone git@github.com:MarcoPellero/wordle.git --branch rustified
```
Then enter the `src/` folder and run `cargo run --release` (or `cargo r -r`).
You can change the wordlist in `data/wordlist.txt`, and you should be able to use a different word size, if you want to you'll need to change the `game::WORD_SIZE` const.
You'll also probably want to change the first guess used, it's currently hard coded. Just search for "sarti" and replace it.

Technically for now there's no way to have a separate guesses & solutions wordlist, but it'd be really easy to implement.
I also don't know how to write Rust code such that it can be imported from other projects, so that's not done yet.

### Go implementation:
Import the package from your code:
```go
import "github.com/MarcoPellero/wordle/src/solver"
```
you can now use the `solver.chooseGuess()` function by passing it two `[]string`s, the first one being the wordlist of usable guesses, and the second one the wordlist of possible solutions, like this:
```go
guess, expected_information, err := solver.ChooseGuess(usable_guesses, remaining_solutions)
```
After you've gotten the feedback for this guess, you need to parse it to the `solver.Feedback` type, which is a slice of enums with the possible states of `Black`, `Yellow`, and `Green`, like this:
```go
func (response string) Feedback() solver.Feedback {
	fd := make(solver.Feedback, len(data.Word))
	for i, c := range data {
		switch c {
		case "b":
			fd[i] = solver.Black
		case "y":
			fd[i] = solver.Yellow
		case "g":
			fd[i] = solver.Green
		}
	}

	return fd
}
```
Now you can use it to get the remaining solutions:
```go
// signature: ([]string, string, []solver.CharState (the B/Y/G Enum)
solutions = solver.FilterSolutions(solutions, guess, feedback)
```

## How it works
Watch [3b1b's video](https://www.youtube.com/watch?v=v68zYyaEmEA) to understand the algorithm. His code is more mathy and less codey than mine
The only difference is that i haven't implemented behaviour to look forward more than 1 guess (yet? :) ), which improves the guessing performances quite a bit.
