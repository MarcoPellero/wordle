use std::collections::HashMap;
use fast_math::log2_raw;

use crate::game;

pub struct BaseAlgo<'a> {
	wordlist: &'a Vec<String>,
	possible_solutions: Vec<String>,
	round: u64
}

impl BaseAlgo<'_> {
	pub fn new(wordlist: &Vec<String>) -> BaseAlgo {
		BaseAlgo { wordlist: wordlist, possible_solutions: vec![], round: 0 }
	}

	fn filter_solution(guess: &str, feedback: &Vec<game::Feedback>, possible_solution: &str) -> bool {
		if guess == possible_solution {
			return game::Feedback::is_solution(feedback);
		}

		let feedback2 = game::Feedback::generate(guess, possible_solution);
		return game::Feedback::cmp(feedback, &feedback2);
	}

	fn rate_guess(&self, guess: &str) -> f32 {
		let mut remaining_solutions = vec![0u64; 3usize.pow(guess.len() as u32)];

		for solution in self.possible_solutions.iter() {
			let feedback = game::Feedback::generate(guess, solution);
			remaining_solutions[game::Feedback::hash(&feedback)] += 1;
		}

		remaining_solutions
			.iter()
			.map(|x| {
				if *x == 0 {
					0f32
				} else {
					let px = (*x as f32) / (self.possible_solutions.len() as f32);
					-px * log2_raw(px)
				}
			})
			.sum()
	}
}

impl game::Algorithm for BaseAlgo<'_> {
	fn init(&mut self) {
		self.possible_solutions = self.wordlist.clone();
		self.round = 0;
	}

	fn guess(&mut self) -> String {
		self.round += 1;
		if self.round == 1 {
			return "sarti".to_owned();
		}
		
		if self.possible_solutions.len() <= 2 {
			return self.possible_solutions[0].to_owned();
		}

		let ratings = self.wordlist
			.iter()
			.map(|guess| self.rate_guess(guess));

		let mut best_idx = 0;
		let mut best_rating = 0f32;
		for (i, rating) in ratings.enumerate() {
			if rating > best_rating {
				best_idx = i;
				best_rating = rating;
			}
		}

		self.wordlist[best_idx].clone()
	}

	fn update(&mut self, guess: String, feedback: &Vec<game::Feedback>) {
		self.possible_solutions = self.possible_solutions
			.iter()
			.filter(|word| BaseAlgo::filter_solution(&guess, &feedback, *word))
			.map(|word| word.to_owned())
			.collect();
	}
}
