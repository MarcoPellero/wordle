use std::collections::HashMap;

use fast_math::log2_raw;

use crate::game;

pub struct Guesser<'a> {
	wordlist: &'a Vec<String>,
	possible_solutions: Vec<String>,
	best_guess_cache: HashMap<Vec<String>, String>,
	round: u64
}

impl Guesser<'_> {
	pub fn new(wordlist: &Vec<String>) -> Guesser {
		Guesser {
			wordlist: wordlist,
			possible_solutions: vec![],
			best_guess_cache: HashMap::new(),
			round: 0
		}
	}

	fn filter_solution(guess: &str, feedback: usize, possible_solution: &str) -> bool {
		if guess == possible_solution {
			return feedback == game::FDHASH_MAX;
		}

		let feedback2 = game::generate_feedback_hash(guess, possible_solution);
		return feedback == feedback2;
	}

	fn rate_guess(&self, guess: &str) -> f32 {
		let mut remaining_solutions = [0u64; 3usize.pow(game::WORD_SIZE as u32)];

		for solution in self.possible_solutions.iter() {
			let feedback = game::generate_feedback_hash(guess, solution);
			remaining_solutions[feedback] += 1;
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

impl game::Algorithm for Guesser<'_> {
	fn init(&mut self) {
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

		match self.best_guess_cache.get(&self.possible_solutions) {
			Some(v) => return v.to_owned(),
			None => {}
		};

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

		self.best_guess_cache.insert(self.possible_solutions.clone(), self.wordlist[best_idx].clone());

		self.wordlist[best_idx].clone()
	}

	fn update(&mut self, guess: &str, feedback: usize) {
		let old_solutions = if self.round <= 1 { self.wordlist } else { &self.possible_solutions };
		self.possible_solutions = old_solutions
			.iter()
			.filter(|word| Guesser::filter_solution(&guess, feedback, *word))
			.map(|word| word.to_owned())
			.collect();
	}
}
