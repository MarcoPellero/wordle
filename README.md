# Wordle Algorithm
This is an alternative implementation to 3b1b's algorithm for playing wordle optimally, using information theory.

It's made not so much for just simulations, but also for interactive use, possibly programmatically for bots. It's usable as a Go package directly, but it had (and will have) an HTTP api for use in bots written in other languages.

## How to use it
I'll write this later :P

## How it works
The basic concepts are the same as 3b1b's algorithm, but the implementation is totally different, mostly because i couldn't understand his code to adapt it.

It goes like this: to choose the best guess i look at each possible one concurrently, and for each guess i look at each possible feedback (the colors), i see how many solutions would be left if i guessed that word and got that feedback, i apply the entropy formula (`-n log n`) on the ratio of solutions left to solutions now, and i sum all of those values for all feedbacks, then i take the guess that maximizes this sum.

That's it! That's the algorithm! Except... this is extremely slow, i perform `(N guesses) * 3^(Word length) * (N solutions)` operations to find the best guess, so how can i speed this up as well as lighten the CPU load for real-time use? Well, the only thing i can reduce is that `N solutions` and make some assumptions to eliminate that `N guesses` (That's to say, i hardcoded the first guess).

Basically, i precompute the heaviest calculations, that is, i generate a cache which tells me *"if i guess this word, and i get this feedback, the next best guess is ..."*. You can adjust the depth of this cache in code, but i've found that a 1 layer cache takes a simulation from unbearably long to like 15s, 2 layes takes it down to .7s, and 3 layers takes it down to .2s.

3 Layers only weigh 550Kb~, so it's not memory hogging either. It takes around 15s to generate it and can be dumped and read from a file.

I think there's LOADS more room for improvement, for example when rating a guess, i'm searching for what solutions are left for each possible feedback, but i think with some smart way one could recycle the solution set for a feedback to use for searching the set for another feedback (for example, when going from BBBBB to BBBBY, you can exclude the previous solution set. But i'm not sure about a more formal correlation).

Another improvement would be being able to run calculations on a GPU, i know that 3b1b does some crazy matrix stuff to precalculate feedbacks that he runs on the GPU with numpy, but 1. i tried it (not with a GPU) and saw no major improvement in time by not needing to calculate them on the fly, and 2. i couldn't understand any of it. Still, he does something right clearly since he can simulate playing against the whole dictionary in just 8s on my machine if we exclude the time to precalculate feedbacks.
