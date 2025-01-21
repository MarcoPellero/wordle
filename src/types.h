#include <stddef.h>

struct Wordlist {
	size_t size;
	char* raw;
};

struct SearchSpace {
	size_t size;
	size_t* idxs;
};

struct Guesser {
	struct Wordlist wordlist;
	struct SearchSpace search_space;
};
