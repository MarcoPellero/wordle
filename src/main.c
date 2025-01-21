#include <stdlib.h>
#include <stdio.h>
#include <assert.h>
#include <stdbool.h>
#include <sys/types.h>
#include <string.h>
#include <math.h>
#include <stdint.h>

#include "types.h"
#include "hashmap.h"

#define WORD_LENGTH 5

struct Wordlist wordlist;
struct SearchSpace search_space;
size_t max_fd; // (3^WORD_LENGTH) - 1
float *log_precalc; // [log2(x) for x in range(wordlist.size)]
size_t sarti_idx;
uint8_t* fd_memo;

#define GET_WORD(idx) (wordlist.raw + (idx)*(WORD_LENGTH + 1))

void read_wordlist(void) {
	FILE *fp = fopen("wordlist.txt", "r");
	assert(fp != NULL);

	fseek(fp, 0, SEEK_END);
	size_t file_size = ftell(fp);
	fseek(fp, 0, SEEK_SET);

	wordlist.size = file_size / (WORD_LENGTH + 1);
	wordlist.raw = (char*)calloc(file_size, sizeof(char));
	
	fread(wordlist.raw, sizeof(char), file_size, fp);
	for (size_t i = 0; i < wordlist.size; i++)
		GET_WORD(i)[WORD_LENGTH] = 0;
}

uint8_t feedback(size_t guess_idx, size_t answer_idx) {
	uint8_t* memo = fd_memo + answer_idx * wordlist.size + guess_idx;
	if (*memo) return *memo-1;
	
	size_t alphabet[26] = { 0 };
	uint8_t result = 0;
	bool green[WORD_LENGTH] = { false };

	char* guess = GET_WORD(guess_idx);
	char* answer = GET_WORD(answer_idx);

	for (size_t i = 0, mult = 1; i < WORD_LENGTH; i++, mult *= 3)
		if (guess[i] == answer[i]) {
			green[i] = true;
			result += mult * 2;
		} else
			alphabet[answer[i] - 'a']++;
	
	for (size_t i = 0, mult = 1; i < WORD_LENGTH; i++, mult *= 3) {
		bool yellow = !green[i] && (alphabet[guess[i] - 'a'] > 0);
		result += mult * yellow;
		alphabet[guess[i] - 'a'] -= yellow;
	}

	*memo = result+1;
	return result;
}

char *pretty_fd(size_t fd) {
	static char buf[64];

	bool green[WORD_LENGTH] = { false };
	bool yellow[WORD_LENGTH] = { false };

	for (ssize_t i = WORD_LENGTH-1, tmp = (max_fd + 1) / 3; i >= 0; i--) {
		green[i] = fd / (tmp * 2);
		yellow[i] = fd / tmp && !green[i];

		fd -= tmp * 2 * green[i];
		fd -= tmp * yellow[i];
		tmp /= 3;
	}

	memset(buf, 0, sizeof(buf));
	for (size_t i = 0; i < WORD_LENGTH; i++) {
		strcat(buf, 
			green[i] ? "\033[42m" :
			yellow[i] ? "\033[43m" :
			"\033[40m"
		);
		
		strcat(buf, " ");
	}
	strcat(buf, "\033[0m");

	return buf;
}

float rate_guess(size_t guess) {
	uint16_t fds[max_fd+1];
	memset(fds, 0, sizeof(fds));

	for (size_t i = 0; i < search_space.size; i++)
		fds[feedback(guess, search_space.idxs[i])]++;
	
	float rating = 0;
	for (size_t i = 0; i <= max_fd; i++) {
		if (!fds[i])
			continue;
		
		float px = ((float)fds[i]) / ((float)search_space.size);

		// rating += -px * log2f(px);
		// log2(a / b) = log2(a) - log2(b)
		rating += -px * (log_precalc[fds[i]] - log_precalc[search_space.size]);
	}

	return rating;
}

size_t choose_guess(void) {
	if (search_space.size == wordlist.size)
		return sarti_idx;
	if (search_space.size <= 2)
		return search_space.idxs[0];

	ssize_t memo_answer = map_get(&search_space);
	if (memo_answer != -1)
		return memo_answer;
	
	float best_rating = 0;
	size_t best_idx = 0;
	for (size_t i = 0; i < wordlist.size; i++) {
		float rating = rate_guess(i);
		if (rating >= best_rating) {
			best_rating = rating;
			best_idx = i;
		}
	}

	map_set(&search_space, best_idx);
	return best_idx;
}

void update(size_t guess, size_t fd) {
	for (size_t i = 0; i < search_space.size;) {
		size_t cur_fd = feedback(guess, search_space.idxs[i]);

		if (cur_fd == fd) {
			i++;
			continue;
		}
		
		search_space.idxs[i] = search_space.idxs[--search_space.size];
	}
}

size_t play(size_t answer) {
	search_space.size = wordlist.size;
	for (size_t i = 0; i < search_space.size; i++)
		search_space.idxs[i] = i;

#ifdef DEBUG
	printf("%s: ", GET_WORD(answer));
#endif

	size_t guesses;
	for (guesses = 1; true; guesses++) {
		size_t guess = choose_guess();
		size_t fd = feedback(guess, answer);

#ifdef DEBUG
		printf("%s %s", GET_WORD(guess), pretty_fd(fd));
		if (fd != max_fd) printf(" -> ");
#endif
		if (fd == max_fd) break;

		update(guess, fd);
		if (!search_space.size) {
			puts("Error: no solutions remaining");
			exit(EXIT_FAILURE);
		}
	}

#ifdef DEBUG
	putchar('\n');
#endif

	return guesses;
}

int main(void) {
	max_fd = pow(3, WORD_LENGTH) - 1;
	map_init();

	read_wordlist();
	search_space.idxs = (size_t*)calloc(wordlist.size, sizeof(size_t));
	printf("Wordlist contains %zu words; first: '%s', last: '%s'\n", wordlist.size, GET_WORD(0), GET_WORD(wordlist.size-1));

	log_precalc = calloc(wordlist.size + 1, sizeof(float));
	for (size_t i = 0; i <= wordlist.size; i++)
		log_precalc[i] = log2f((float)i);
	
	while (strcmp(GET_WORD(sarti_idx), "sarti") != 0)
		sarti_idx++;

	fd_memo = calloc(wordlist.size * wordlist.size, sizeof(uint8_t));
	assert(fd_memo != NULL);

	size_t sum = 0;
	for (size_t i = 0; i < wordlist.size; i++)
		sum += play(i);
	
	printf("Mean: %f\n", ((float)sum) / ((float)wordlist.size));
	return EXIT_SUCCESS;
}
