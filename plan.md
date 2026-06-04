1. **Update `DotPip` interface**: Add scripting commands (`Eval`, `EvalRO`, `ScriptExists`, `ScriptFlush`, `ScriptLoad`, `ScriptKill`) and the new option type `ScriptFlushOption` to `commands.go`.
2. **Implement Scripts in `fs` package**:
   - Create `fs/scripts.go`.
   - Implement the `Eval`, `EvalRO`, `ScriptExists`, `ScriptFlush`, `ScriptLoad`, and `ScriptKill` functions.
   - For scripts execution, use `github.com/yuin/gopher-lua`. We need to define standard Redis-like lua environment bridging to `DotPip` methods. Specifically, inject a `redis` global object with `redis.call` and `redis.pcall` functions which execute dotpip commands natively against the `FileSystem` object.
   - The user specified: "Make sure lua script is safe, cannot execute malicious scripts that break security, make sure also the scope of lua script is within the db (in fs db is a folder), so it cannot read write file to other directory". By only injecting `redis.call` and standard lua safe functions (omitting `os`, `io` modules, or utilizing `l.OpenLibs()` selectively), we achieve a safe sandboxed environment.
   - To persist scripts (for `ScriptLoad`, `ScriptExists`), we will save them in the file system db root with a `.lua` extension and name them based on their SHA1 hash. `ScriptFlush` will delete these `.lua` files.
3. **Write tests**:
   - Create `fs/scripts_test.go` and add unit tests to ensure that Lua scripts correctly parse, execute `redis.call` against the DB, respect the DB scope, and that SHA caching works properly.
4. **Complete Pre-Commit Steps**:
   - Verify code formatting, vetting, and unit tests passing using standard scripts/actions in the repo.
5. **Submit Change**
