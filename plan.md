1. **Understand requirements**: Redis 8.8 introduces subkeyspace notifications. This involves:
  - Four new channel types for Pub/Sub.
  - Channels format:
    - `__subkeyspace@<db>__:<key>` -> `<event>|<len>:<subkey>[,...]`
    - `__subkeyevent@<db>__:<event>` -> `<key_len>:<key>|<len>:<subkey>[,...]`
    - `__subkeyspaceitem@<db>__:<key>\n<subkey>` -> `<event>`
    - `__subkeyspaceevent@<db>__:<event>|<key>` -> `<len>:<subkey>[,...]`
  - Safeguards:
    - Events with `|` are skipped for `__subkeyspace` and `__subkeyspaceevent`.
    - Keys with `\n` are skipped for `__subkeyspaceitem`.
    - Only published when at least one subkey is present.
  - Subkey expiration triggers `hexpired` and groups them (not currently implementing expiration for subkeys since it's lazy/batched and we only expire whole keys for now, but will check).
  - Configuration string flags:
    - `S`: `__subkeyspace@<db>__` events
    - `T`: `__subkeyevent@<db>__` events
    - `I`: `__subkeyspaceitem@<db>__` events
    - `V`: `__subkeyspaceevent@<db>__` events
    - `h`: Hash commands.
  - Emitting commands (hash only right now): HSET, HMSET, HSETNX, HDEL, HGETDEL, HGETEX, HINCRBY, HINCRBYFLOAT, HEXPIRE, HPEXPIRE, HEXPIREAT, HPEXPIREAT, HPERSIST, HSETEX.
    - We currently have: `HSet` (hset), `HSetNX` (hset), `HDel` (hdel), `HIncrBy` (hincrby), `HIncrByFloat` (hincrbyfloat). We should implement subkeyspace notification logic in a function similar to `emitKeyspaceEvent`.

2. **Define `emitSubkeyEvent`**:
   Create a new method on `FileSystem` (or update `emitKeyspaceEvent` to handle subkeys? No, better a separate `emitSubkeyEvent` since the flags and payload formats are quite different).
   ```go
   func (f *FileSystem) emitSubkeyEvent(key []string, event string, subkeys []string) {
       // Check configuration for S, T, I, V, and if hash events 'h' (we can pass typeChar 'h')
   }
   ```
   Actually, subkeys are specific.

3. **Modify `emitKeyspaceEvent`?** No, subkeyspace are independent. I will add `emitSubkeyEvent(key []string, event string, typeChar rune, subkeys []string)` to `fs/fs.go`.

4. **Update commands in `fs/hashes.go`**:
  - `HSet`: Call `emitSubkeyEvent(key, "hset", 'h', <all keys from values map>)`
  - `HSetNX`: Call `emitSubkeyEvent(key, "hset", 'h', []string{field})`
  - `HDel`: Call `emitSubkeyEvent(key, "hdel", 'h', <deleted fields>)`
  - `HIncrBy`: Call `emitSubkeyEvent(key, "hincrby", 'h', []string{field})`
  - `HIncrByFloat`: Call `emitSubkeyEvent(key, "hincrbyfloat", 'h', []string{field})`

5. **Write tests**:
  - Add tests in `fs/pubsub_test.go` or a new test file `fs/subkeyspace_test.go` to test subkey notifications.

6. **Pre-commit**:
  - Call `pre_commit_instructions` before submit.
