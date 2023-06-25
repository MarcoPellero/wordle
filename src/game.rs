pub const WORD_SIZE: usize = 5;

#[derive(Clone, Copy, PartialEq, Eq, Hash)]
pub enum Feedback {
	Black,
	Yellow,
	Green
}

pub type FeedbackArr = [Feedback; WORD_SIZE];

impl Feedback {
	pub fn generate(guess: &str, solution: &str) -> FeedbackArr {
		let mut alphabet = [0u8; 26];
		let mut fd_chars = [Feedback::Black; WORD_SIZE];
	
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

	pub fn is_solution(feedback: &FeedbackArr) -> bool {
		feedback
			.iter()
			.map(|c| *c == Feedback::Green)
			.reduce(|total, c| total && c)
			.unwrap()
	}

	pub fn cmp(a: &FeedbackArr, b: &FeedbackArr) -> bool {
		a
			.iter()
			.zip(b)
			.map(|(ac, bc)| *ac == *bc)
			.reduce(|total, c| total && c)
			.unwrap()
	}

	pub fn hash(feedback: &FeedbackArr) -> usize {
		let mut acc = 0;
		let mut mul = 1;

		for c in feedback {
			acc += mul * match *c {
				Feedback::Black => 0,
				Feedback::Yellow => 1,
				Feedback::Green => 2
			};
			mul *= 3;
		}

		acc
	}

	pub fn to_str(feedback: &FeedbackArr) -> String {
		feedback
			.iter()
			.map(|c| match c {
				Feedback::Black => 'b',
				Feedback::Yellow => 'y',
				Feedback::Green => 'g'
			})
			.collect()
	}
}


pub trait Algorithm {
	fn init(&mut self);
	fn guess(&mut self) -> String;
	fn update(&mut self, guess: &str, feedback: &FeedbackArr);
}
