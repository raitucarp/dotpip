import re

with open("fs/strings.go", "r") as f:
    content = f.read()

# According to spec:
# APPEND generates an append event. (Wait, Append calls Set! Set generates 'set'. Let's override it or ignore? The spec says APPEND -> append event.)
# SET generates set events. SETEX will also generate expire events.
# GetDel -> no spec mention but should probably generate del. Wait, GetDel calls Del internally? No.
# Incr, Decr, Incrby, Decrby generate incrby events!
# IncrByFloat generates incrbyfloat events!
# SetRange generates setrange event.
# MSet generates a separate set event for every key.
# MSetNX generates separate set event for every key.

# BUT! our implementation of Append, Incr, Decr, IncrBy, DecrBy, SetRange, MSet, MSetNX all call f.Set() internally!
# This means calling `Incr` will emit a "set" event instead of "incrby" event!
# We can fix this by adding an internal version of Set or changing how we emit.
