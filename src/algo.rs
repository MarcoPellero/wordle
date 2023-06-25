use std::collections::HashMap;

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

		let feedback2 = game::generate_feedback(guess, possible_solution);
		return game::Feedback::cmp(feedback, &feedback2);
	}

	fn rate_guess(&self, guess: &str) -> f64 {
		let mut remaining_solutions: HashMap<Vec<game::Feedback>, u64> = HashMap::new();

		for solution in self.possible_solutions.iter() {
			let feedback = game::generate_feedback(guess, solution);
			let key = remaining_solutions.get_mut(&feedback);
			match key {
				Some(v) => { *v += 1; },
				None => { remaining_solutions.insert(feedback, 1); }
			};
		}

		remaining_solutions
			.values()
			.map(|x| {
				if *x == 0 {
					0f64
				} else {
					let px = (*x as f64) / (self.possible_solutions.len() as f64);
					-px * px.log2()
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
		let mut best_rating = 0f64;
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
