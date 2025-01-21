#include <stddef.h>
#include <sys/types.h>

void map_init(void);
ssize_t map_get(struct SearchSpace* key);
void map_set(struct SearchSpace* key, size_t val);
