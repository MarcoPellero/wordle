#include <stdlib.h>
#include <assert.h>
#include <sys/types.h>

#include "types.h"

struct BucketNode {
	struct BucketNode* next;
	size_t hash;
	size_t val;
};

const size_t n_buckets = 500;

struct BucketNode** buckets;

size_t hashfn(struct SearchSpace* key) {
	// fnv1
	size_t hash = 0xcbf29ce484222325;

	for (size_t i = 0; i < key->size; i++) {
		hash ^= key->idxs[i];
		hash *= 1099511628211;
	}

	return hash;
}

void map_init(void) {
	buckets = calloc(n_buckets, sizeof(struct BucketNode*));
	assert(buckets != NULL);
}

struct BucketNode** map_walk(size_t hash) {
	size_t bucket_idx = hash % n_buckets;

	struct BucketNode** walk = &buckets[bucket_idx];
	while (*walk && (*walk)->hash != hash)
		walk = &(*walk)->next;
	
	return walk;
}

ssize_t map_get(struct SearchSpace* key) {
	size_t hash = hashfn(key);
	struct BucketNode** node = map_walk(hash);

	if (*node == NULL)
		return -1;
	
	return (*node)->val;
}

void map_set(struct SearchSpace* key, size_t val) {
	size_t hash = hashfn(key);
	struct BucketNode** node = map_walk(hash);

	if (*node != NULL) return;

	*node = calloc(1, sizeof(struct BucketNode));
	assert(*node != NULL);
	
	(*node)->hash = hash;
	(*node)->val = val;
}
