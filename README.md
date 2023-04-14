# Wordle Algorithm
This is a better implementation of 3b1b's algorithm for playing wordle optimally, using information theory.

It's made both for simulations (like 3b1b's), but also for interactive use, possibly programmatically for bots. It's usable as a Go package directly, but it had (and will have) an HTTP api for use in bots written in other languages aswell.

## How to use it
I'll write this later :P

## How it works
The basic concepts are the same as 3b1b's algorithm, but the implementation is totally differen.

It goes like this: when looking for the best guess i look at each one of them concurrently, and i calculate for each one how many solutions would be left for any possible feedback i could get, i then apply Shannon's formula on all of those numbers, and that's the expected information for that particular guess. Of course i then look for the guess that maximizes this value.

That's it! That's the algorithm! Except... this is extremely slow, it's `O(AB)` (for A=number of guesses, and B=number of solutions), notably with a very heavy constant factor, operations to find the best guess, so how can i speed this up as well as lighten the CPU load for real-time use? Well, the only thing i can reduce is that `N solutions` part and make some assumptions to eliminate that `N guesses` (That's to say, i hardcoded the first guess).

Basically, i precompute the heaviest calculations, that is, i generate a cache which tells me *"if i guess this word, and i get this feedback, the next best guess is ..."*. You can adjust the depth of this cache in code, but i've found that a 1 layer cache takes a simulation from unbearably long to like 5.6s, 2 layes takes it down to .5s, and 3 layers takes it down to .2s.

3 Layers only weigh 550Kb~, so it's not memory hogging either. It takes around 2s to generate either 2 or 3 laters, and it can be dumped and read from a file.

Notice that all of these times are for a "play-all" simulation, which uses a cache or a chosen first guess to play against the whole dictionary of solutions once. The times for real time use are obviously much lower, for example the simulations that took .2s were using a 3k word dictionary, which means that it took 60μs to play ***a whole game***, that's like 10μs per guess (though it's not actually that linear).
