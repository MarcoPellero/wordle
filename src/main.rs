use std::cmp::min;
use std::fs::File;
use std::io::{self, BufRead};
use std::path::Path;

mod algo;
mod game;

const LOG_LEVEL: u8 = 2;
const MAX_RUNS: usize = 300;

fn read_wordlist(path: &str) -> Result<Vec<String>, io::Error> {
    let file = File::open(Path::new(path))?;
    let reader = io::BufReader::new(file);

	reader
		.lines()
		.collect()
}

fn word_simulation(guesser: &mut impl game::Algorithm, solution: &str) -> u64 {
	guesser.init();

	for i in 1.. {
		let next_guess = guesser.guess();
		let fd = game::Feedback::generate(next_guess.as_str(), solution);
		if LOG_LEVEL >= 2 {
			println!("{} | {} = {}", next_guess, solution, game::Feedback::to_str(&fd));
		}

		if game::Feedback::is_solution(&fd) {
			if LOG_LEVEL >= 1 {
				println!("{} done in {}g", solution, i);
			}
			return i;
		}

		guesser.update(next_guess, &fd);
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
	let wordlist = read_wordlist("../data/wordlist.txt").unwrap();
	println!("Read {} words from file", wordlist.len());

	let mut guesser = algo::BaseAlgo::new(&wordlist);

	let mean = dictionary_simulation(&mut guesser, &wordlist);
	println!("Mean guesses: {}", mean);
}
