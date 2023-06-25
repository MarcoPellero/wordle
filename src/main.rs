use std::cmp::min;
use std::fs::File;
use std::io::{self, BufRead, BufReader};
use std::path::Path;
use std::collections::HashMap;

type Word = u32;

fn gen_feedback(wordlist: &Vec<Word>, guess_idx: u16, solution_idx: u16) -> u8 {
	let mut sol_alpha = [0u8; 26];
	let mut fd = 0;

	let mut mul = 1;
	for i in 0..5 {
		let guess_c = (wordlist[guess_idx as usize] >> (i*5)) & 0b11111;
		let sol_c = (wordlist[solution_idx as usize] >> (i*5)) & 0b11111;

		if guess_c == sol_c {
			fd += mul*2;
		} else {
			sol_alpha[sol_c as usize] += 1;
		}

		mul *= 3;
	}

	mul = 1;
	for i in 0..5 {
		let guess_c = (wordlist[guess_idx as usize] >> (i*5)) & 0b11111;
		let sol_c = (wordlist[solution_idx as usize] >> (i*5)) & 0b11111;

		if guess_c != sol_c && sol_alpha[guess_c as usize] > 0 {
			fd += mul;
			sol_alpha[guess_c as usize] -= 1;
		}

		mul *= 3;
	}

	fd
}

fn filter_solutions(wordlist: &Vec<Word>, solutions: &Vec<u16>, guess_idx: u16, fd: u8) -> Vec<u16> {
	if fd == 242 {
		return vec![];
	}

	solutions
		.iter()
		.filter(|x| fd == gen_feedback(wordlist, guess_idx, **x))
		.map(|x| *x)
		.collect()
}

fn entropy_formula(old_solutions: u16, new_solutions: u16) -> f64 {
	if new_solutions == 0 {
		return 0 as f64;
	}

	let px = (new_solutions as f64) / (old_solutions as f64);
	-px * px.log2()
}

fn rate_guess(wordlist: &Vec<Word>, old_solutions: &Vec<u16>, guess_idx: u16) -> f64 {
	let mut solutions_left = [0 as u16; 3i32.pow(5 as u32) as usize];

	for word in old_solutions {
		let fd = gen_feedback(wordlist, guess_idx, *word);
		solutions_left[fd as usize] += 1;
	}

	solutions_left
		.iter()
		.map(|x| entropy_formula(old_solutions.len() as u16, *x))
		.sum()
}

fn choose_guess(guesses: &Vec<Word>, solutions: &Vec<u16>) -> Result<u16, String> {
	if solutions.len() <= 2 {
		return Ok(solutions[0]);
	}

	let mut best_idx = 0u16;
	let mut best_val = 0f64;
	for i in 0..guesses.len() {
		let val = rate_guess(guesses, solutions, i as u16);
		if val > best_val {
			best_idx = i as u16;
			best_val = val;
		}
	}

	Ok(best_idx)
}

fn read_wordlist(path: &str) -> Result<Vec<Word>, io::Error> {
    let file = File::open(Path::new(path))?;
    let reader = BufReader::new(file);

    let mut wordlist = Vec::new();

    for line in reader.lines() {
        let line = line?;
        wordlist.push(line);
    }

    Ok(fix_wordlist(&wordlist))
}

fn fix_word(word: &String) -> Word {
	word
		.as_bytes()
		.iter()
		.enumerate()
		.map(|(i, c)| ((c - ('a' as u8)) as u32) << (i*5))
		.sum()
}

fn fix_wordlist(bad: &Vec<String>) -> Vec<Word> {
	bad
		.iter()
		.map(fix_word)
		.collect()
}

fn play_game(fd_map: &mut HashMap<Vec<u16>, u16>, guesses: &Vec<Word>, solution_idx: u16, first_guess_idx: u16) -> u8 {
	let mut solutions = (0..guesses.len()).map(|i| i as u16).collect::<Vec<u16>>();
	let mut guess_idx = first_guess_idx;
	for i in 1.. {
		let fd = gen_feedback(guesses, guess_idx, solution_idx);
		if fd == 242 {
			// println!("{} in {}", solution_idx, i);
			return i;
		}
		solutions = filter_solutions(guesses, &solutions, guess_idx, fd);

		let best_next;
		if solutions.len() >= 3 {
			best_next = match fd_map.get(&solutions) {
				Some(v) => *v,
				None => {
					let best_next = choose_guess(guesses, &solutions).unwrap();
					fd_map.insert(solutions.clone(), best_next);
					best_next
				}
			};
		} else {
			best_next = choose_guess(guesses, &solutions).unwrap();
		}

		guess_idx = best_next;
	}

	unreachable!();
}

fn main() {
	let wordlist = read_wordlist("./wordlist.txt").unwrap();
	let first_guess_str = "sarti".to_owned();
	let first_guess = fix_word(&first_guess_str);
	let first_guess_ids = wordlist.iter().position(|x| *x == first_guess).unwrap() as u16;

	let mut fd_map: HashMap<Vec<u16>, u16> = HashMap::new();
	let runs = min(4000, wordlist.len());
	let sum: u32 = (0..runs)
		.map(|i| play_game(&mut fd_map, &wordlist, i as u16, first_guess_ids) as u32)
		.sum();

	let mean = (sum as f64) / (runs as f64);
	println!("Mean: {}", mean);
}
