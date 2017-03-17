import random
import string

def token(l):
    return ''.join(random.choice(string.ascii_lowercase + string.digits) for _ in range(l))

print token(6) + "." + token(16)
