use std::cmp::min;
use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;

mod algo;
mod game;

const LOG_LEVEL: u8 = 2;
const MAX_RUNS: usize = 4000;

fn read_wordlist(path: &str) -> Vec<String> {
    let file = File::open(Path::new(path)).unwrap();
    let reader = io::BufReader::new(file);

	let wordlist: Vec<String> = reader
		.lines()
		.map(|line| line.unwrap())
		.collect();

	for word in wordlist.iter() {
		assert_eq!(word.len(), game::WORD_SIZE);
	}

	wordlist
}

fn word_simulation(guesser: &mut impl game::Algorithm, solution: &str) -> u64 {
	guesser.init();

	for i in 1.. {
		let next_guess = guesser.guess();
		let fd = game::generate_feedback_hash(next_guess.as_str(), solution);
		if LOG_LEVEL >= 2 {
			println!("{} | {} = {}", next_guess, solution, game::to_str(fd));
		}

		if fd == game::FDHASH_MAX {
			if LOG_LEVEL >= 1 {
				println!("{} done in {}g", solution, i);
			}
			return i;
		}

		guesser.update(&next_guess, fd);
	}

	unreachable!()
}

fn dictionary_simulation(guesser: &mut impl game::Algorithm, wordlist: &Vec<String>) -> f64 {
	let runs = min(MAX_RUNS, wordlist.len());
	let score_sum: u64 = (0..runs)
		.map(|i| &wordlist[i])
		.map(|solution| word_simulation(guesser, solution))
		.sum();

	(score_sum as f64) / (runs as f64)
}

fn main() {
	let wordlist = read_wordlist("../data/wordlist.txt");
	println!("Read {} words from file", wordlist.len());

	let mut guesser = algo::Guesser::new(&wordlist);

	let mean = dictionary_simulation(&mut guesser, &wordlist);
	println!("Mean guesses: {}", mean);
}
