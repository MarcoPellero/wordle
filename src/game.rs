pub trait Algorithm {
	fn init(&mut self);
	fn guess(&mut self) -> String;
	fn update(&mut self, guess: String, feedback: &Vec<Feedback>);
}

#[derive(Clone, PartialEq, Eq, Hash)]
pub enum Feedback {
	Black,
	Yellow,
	Green
}

impl Feedback {
	pub fn generate(guess: &str, solution: &str) -> Vec<Feedback> {
		let mut alphabet = [0u8; 26];
		let mut fd_chars = vec![Feedback::Black; guess.len()];
	
		const TO_NUM: u8 = 'a' as u8;
		let guess_bytes = guess.as_bytes();
		let solution_bytes = solution.as_bytes();
		
		for i in 0..guess.len() {
			if guess_bytes[i] == solution_bytes[i] {
				fd_chars[i] = Feedback::Green;
			} else {
				alphabet[(solution_bytes[i] - TO_NUM) as usize] += 1;
			}
		}
	
		for i in 0..guess.len() {
			if guess_bytes[i] == solution_bytes[i] {
				continue;
			}
			
			if alphabet[(guess_bytes[i] - TO_NUM) as usize] > 0 {
				fd_chars[i] = Feedback::Yellow;
				alphabet[(guess_bytes[i] - TO_NUM) as usize] -= 1;
			}
		}
	
		fd_chars
	}

	pub fn is_solution(feedback: &Vec<Feedback>) -> bool {
		feedback
			.iter()
			.map(|c| *c == Feedback::Green)
			.reduce(|total, c| total && c)
			.unwrap()
	}

	pub fn cmp(a: &Vec<Feedback>, b: &Vec<Feedback>) -> bool {
		a
			.iter()
			.zip(b)
			.map(|(ac, bc)| *ac == *bc)
			.reduce(|total, c| total && c)
			.unwrap()
	}
}
