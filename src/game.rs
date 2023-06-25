pub trait Algorithm {
	fn init(&mut self);
	fn guess(&mut self) -> String;
	fn update(&mut self, guess: String, feedback: String);
}

pub fn generate_feedback(guess: &str, solution: &str) -> String {
	let mut alphabet = [0u8; 26];
	let mut fd_chars = vec!['b'; guess.len()];

	const TO_NUM: u8 = 'a' as u8;
	let guess_bytes = guess.as_bytes();
	let solution_bytes = solution.as_bytes();
	
	for i in 0..guess.len() {
		if guess_bytes[i] == solution_bytes[i] {
			fd_chars[i] = 'g';
		} else {
			alphabet[(solution_bytes[i] - TO_NUM) as usize] += 1;
		}
	}

	for i in 0..guess.len() {
		if guess_bytes[i] == solution_bytes[i] {
			continue;
		}
		
		if alphabet[(guess_bytes[i] - TO_NUM) as usize] > 0 {
			fd_chars[i] = 'y';
			alphabet[(guess_bytes[i] - TO_NUM) as usize] -= 1;
		}
	}

	fd_chars.iter().collect()
}
