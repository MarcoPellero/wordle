known_answers_path = "./solutions"
known_invalid_path = "./ALL_BAD"
testing_against_path = "./guesses"

known_accents = "òìâèóùûîéà"
known_letters = "oiaeouuiea"
banished_letters = "'"
accent_map = {k:v for k, v in zip(known_accents, known_letters)}

dont_look = set(open(known_answers_path).read().split("\n")) | set(open(known_invalid_path).read().split("\n"))

words = open(testing_against_path).read().lower().split("\n")
words = list(set(words))
# replace accents
words = ["".join(map(lambda c: accent_map.get(c, c), x)) for x in words]
print(len(words))

# only 5 letters
words = [x for x in words if len(x) == 5]
print(len(words))

# only untested
words = [x for x in words if x not in dont_look]
print(len(words))

# only letters
words = [x for x in words if all(map(lambda c: c not in banished_letters, x))]
print(len(words))

accents = "".join({c for c in "".join(words) if ord(c) not in range(ord("a"), ord("z")+1)})
open(testing_against_path, "w").write("\n".join(words))
print(accents)
