override CFLAGS += -std=c99 -Wall -pedantic -Wextra -ggdb -lm -O3 -march=native
override CFLAGS += -Werror

OFILE = main

c:
	$(CC) $(CFLAGS) src/*.c -o $(OFILE)
dbg:
	$(CC) $(CFLAGS) src/*.c -o $(OFILE) -DDEBUG
