pub const WORD_SIZE: usize = 5;

pub fn generate_hash(guess: &str, solution: &str) -> usize {
	let mut alphabet = [0u8; 26];
	let mut feedback = 0usize;
	let mut is_green = [false; WORD_SIZE];
	
	const TO_NUM: u8 = 'a' as u8;
	let guess_bytes = guess.as_bytes();
	let solution_bytes = solution.as_bytes();
		
	let mut mul = 1;
	for i in 0..guess.len() {
		if guess_bytes[i] == solution_bytes[i] {
			is_green[i] = true;
			feedback += mul*2;
		} else {
			alphabet[(solution_bytes[i] - TO_NUM) as usize] += 1;
		}

		mul *= 3;
	}
	
	mul = 1;
	for i in 0..guess.len() {
		if !is_green[i] && alphabet[(guess_bytes[i] - TO_NUM) as usize] > 0 {
			feedback += mul;
			alphabet[(guess_bytes[i] - TO_NUM) as usize] -= 1;
		}

		mul *= 3;
	}
	
	feedback
}

pub fn to_str(feedback: usize) -> String {
	let mut chars = ['b'; WORD_SIZE];

	let mut mul = 3usize.pow(WORD_SIZE as u32 - 1);
	let mut fd = feedback;
	for i in (0..WORD_SIZE).rev() {
		if fd >= mul*2 {
			chars[i] = 'g';
			fd -= mul*2;
		} else if fd >= mul {
			chars[i] = 'y';
			fd -= mul;
		}

		mul /= 3;
	}

	chars
		.iter()
		.collect()
}


pub trait Algorithm {
	fn init(&mut self);
	fn guess(&mut self) -> String;
	fn update(&mut self, guess: &str, feedback: usize);
}