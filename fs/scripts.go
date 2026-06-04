package fs

import (
	"crypto/sha1"
	"dotpip"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func getScriptsDir(pathRoot string) string {
	return filepath.Join(pathRoot, ".scripts")
}

func getScriptPath(pathRoot, hash string) string {
	return filepath.Join(getScriptsDir(pathRoot), hash+".lua")
}

func (f *FileSystem) initLuaState() *lua.LState {
	L := lua.NewState(lua.Options{
		SkipOpenLibs: true,
	})

	for _, pair := range []struct {
		n string
		f lua.LGFunction
	}{
		{lua.BaseLibName, lua.OpenBase},
		{lua.TabLibName, lua.OpenTable},
		{lua.StringLibName, lua.OpenString},
		{lua.MathLibName, lua.OpenMath},
	} {
		if err := L.CallByParam(lua.P{
			Fn:      L.NewFunction(pair.f),
			NRet:    0,
			Protect: true,
		}, lua.LString(pair.n)); err != nil {
			panic(err)
		}
	}

	redisTable := L.NewTable()
	L.SetField(redisTable, "call", L.NewFunction(f.luaRedisCall))
	L.SetField(redisTable, "pcall", L.NewFunction(f.luaRedisPCall))
	L.SetGlobal("redis", redisTable)

	// Inject global functions for all exported methods on DotPip interface
	var dp dotpip.DotPip
	v := reflect.TypeOf(&dp).Elem()
	fVal := reflect.ValueOf(f)

	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		mName := method.Name

		fn := func(mName string, mType reflect.Type) func(L *lua.LState) int {
			return func(L *lua.LState) int {
				return f.invokeGoMethod(L, fVal, mName, mType)
			}
		}(mName, method.Type)

		L.SetGlobal(mName, L.NewFunction(fn))
	}

	return L
}

func (f *FileSystem) luaRedisCall(l *lua.LState) int {
	return f.invokeRedisCommand(l, false)
}

func (f *FileSystem) luaRedisPCall(l *lua.LState) int {
	return f.invokeRedisCommand(l, true)
}

func (f *FileSystem) invokeRedisCommand(l *lua.LState, pcall bool) int {
	cmdName := l.ToString(1)
	if cmdName == "" {
		if pcall {
			l.Push(lua.LString("Please specify at least one argument for redis.call()"))
			return 1
		}
		l.RaiseError("Please specify at least one argument for redis.call()")
		return 0
	}

	// For standard redis commands, map uppercase redis command to Go method name.
	// We'll look up method by name on the interface.
	cmdName = strings.ToUpper(cmdName)

	// Fast rudimentary mapping for common ones (can add more as needed):
	methodName, mapped := redisCommandMapping[cmdName]
	if !mapped {
		methodName = l.ToString(1) // Try exact case if not found
	}

	var dp dotpip.DotPip
	v := reflect.TypeOf(&dp).Elem()
	method, ok := v.MethodByName(methodName)

	if !ok {
		if pcall {
			// standard redis.pcall returns a table with an err field, or just a string.
			// Let's just return a string error. Actually, it's conventional in our simple implementation
			// to return an error string. Or if it raises an error, pcall in Lua catches it.
			// If we push a string, it gets returned as a string. But a real redis.pcall returns a lua table {err="message"}.
			errTbl := l.NewTable()
			l.SetField(errTbl, "err", lua.LString(fmt.Sprintf("Unknown Redis command called from Lua: %s", cmdName)))
			l.Push(errTbl)
			return 1
		}
		l.RaiseError("%s", fmt.Sprintf("Unknown Redis command called from Lua: %s", cmdName))
		return 0
	}

	fVal := reflect.ValueOf(f)

	// Remove the first argument (cmdName) before processing args for the Go method
	numArgs := l.GetTop() - 1
	args := make([]lua.LValue, numArgs)
	for i := 0; i < numArgs; i++ {
		args[i] = l.Get(i + 2)
	}

	return f.invokeGoMethodArgs(l, fVal, method.Name, method.Type, args, pcall)
}

func (f *FileSystem) invokeGoMethod(l *lua.LState, fVal reflect.Value, methodName string, methodType reflect.Type) int {
	numArgs := l.GetTop()
	args := make([]lua.LValue, numArgs)
	for i := 0; i < numArgs; i++ {
		args[i] = l.Get(i + 1)
	}
	return f.invokeGoMethodArgs(l, fVal, methodName, methodType, args, false)
}

func (f *FileSystem) invokeGoMethodArgs(l *lua.LState, fVal reflect.Value, methodName string, methodType reflect.Type, args []lua.LValue, pcall bool) int {
	methodVal := fVal.MethodByName(methodName)
	if !methodVal.IsValid() {
		if pcall {
			errTbl := l.NewTable()
			l.SetField(errTbl, "err", lua.LString(fmt.Sprintf("Method %s not implemented", methodName)))
			l.Push(errTbl)
			return 1
		}
		l.RaiseError("%s", fmt.Sprintf("Method %s not implemented", methodName))
		return 0
	}

	numIn := methodType.NumIn()
	isVariadic := methodType.IsVariadic()

	var in []reflect.Value

	// We need to match arguments.
	for i := 0; i < len(args); i++ {
		var argType reflect.Type
		switch {
		case isVariadic && i >= numIn-1:
			argType = methodType.In(numIn - 1).Elem()
		case i < numIn:
			argType = methodType.In(i)
		default:
			break // Too many arguments passed
		}

		if argType == nil {
			break
		}

		val, err := f.convertLuaArgToGo(args[i], argType)
		if err != nil {
			if pcall {
				errTbl := l.NewTable()
				l.SetField(errTbl, "err", lua.LString(err.Error()))
				l.Push(errTbl)
				return 1
			}
			l.RaiseError("%s", err.Error())
			return 0
		}
		in = append(in, val)
	}

	// Pad with zero values for missing optional/variadic arguments?
	// Variadic functions can take 0 arguments for the variadic part, but non-variadic must have exactly numIn
	if !isVariadic && len(in) < numIn {
		for i := len(in); i < numIn; i++ {
			in = append(in, reflect.Zero(methodType.In(i)))
		}
	} else if isVariadic && len(in) < numIn-1 {
		for i := len(in); i < numIn-1; i++ {
			in = append(in, reflect.Zero(methodType.In(i)))
		}
	}

	out := methodVal.Call(in)

	// Check error which is typically the last return value
	if len(out) > 0 {
		lastOut := out[len(out)-1]
		if lastOut.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) && !lastOut.IsNil() {
			errStr := lastOut.Interface().(error).Error()
			if pcall {
				errTbl := l.NewTable()
				l.SetField(errTbl, "err", lua.LString(errStr))
				l.Push(errTbl)
				return 1
			}
			l.RaiseError("%s", errStr)
			return 0
		}
	}

	// Push results
	numOut := len(out)
	if numOut > 0 && out[numOut-1].Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		numOut--
	}

	for i := 0; i < numOut; i++ {
		l.Push(f.convertGoToLua(l, out[i].Interface()))
	}

	return numOut
}

func (f *FileSystem) convertLuaArgToGo(arg lua.LValue, argType reflect.Type) (reflect.Value, error) {
	if argType == reflect.TypeOf(dotpip.Key{}) {
		return reflect.ValueOf(dotpip.NewKey(arg.String())), nil
	}

	switch argType.Kind() {
	case reflect.String:
		return reflect.ValueOf(arg.String()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, ok := arg.(lua.LNumber); ok {
			return reflect.ValueOf(int(num)).Convert(argType), nil
		}
		return reflect.Zero(argType), fmt.Errorf("expected number for %v", argType)
	case reflect.Float32, reflect.Float64:
		if num, ok := arg.(lua.LNumber); ok {
			return reflect.ValueOf(float64(num)).Convert(argType), nil
		}
		return reflect.Zero(argType), fmt.Errorf("expected float for %v", argType)
	case reflect.Bool:
		if b, ok := arg.(lua.LBool); ok {
			return reflect.ValueOf(bool(b)), nil
		}
		return reflect.Zero(argType), fmt.Errorf("expected bool for %v", argType)
	case reflect.Slice:
		if argType.Elem().Kind() == reflect.String {
			if tbl, ok := arg.(*lua.LTable); ok {
				var res []string
				maxN := tbl.MaxN()
				for i := 1; i <= maxN; i++ {
					res = append(res, tbl.RawGetInt(i).String())
				}
				return reflect.ValueOf(res), nil
			}
		}
		if tbl, ok := arg.(*lua.LTable); ok {
			var res = reflect.MakeSlice(argType, 0, 0)
			maxN := tbl.MaxN()
			for i := 1; i <= maxN; i++ {
				item, err := f.convertLuaArgToGo(tbl.RawGetInt(i), argType.Elem())
				if err != nil {
					return reflect.Zero(argType), err
				}
				res = reflect.Append(res, item)
			}
			return res, nil
		}
		// Fallback for any slice
		return reflect.Zero(argType), fmt.Errorf("unsupported slice type conversion")
	case reflect.Struct:
		// Map struct using reflection and Lua tables
		if tbl, ok := arg.(*lua.LTable); ok {
			val := reflect.New(argType).Elem()
			tbl.ForEach(func(k, v lua.LValue) {
				if k.Type() == lua.LTString {
					// find field case insensitively
					for i := 0; i < argType.NumField(); i++ {
						field := argType.Field(i)
						if strings.EqualFold(field.Name, k.String()) {
							fieldVal, err := f.convertLuaArgToGo(v, field.Type)
							if err == nil {
								val.Field(i).Set(fieldVal)
							}
						}
					}
				}
			})
			return val, nil
		}
		return reflect.Zero(argType), fmt.Errorf("unsupported struct type conversion")
	case reflect.Interface:
		// Convert directly as primitive interface{}
		switch v := arg.(type) {
		case lua.LString:
			return reflect.ValueOf(string(v)), nil
		case lua.LNumber:
			return reflect.ValueOf(float64(v)), nil
		case lua.LBool:
			return reflect.ValueOf(bool(v)), nil
		default:
			return reflect.Zero(argType), fmt.Errorf("unsupported generic interface type")
		}
	}

	return reflect.Zero(argType), fmt.Errorf("unsupported type conversion: %v", argType)
}

func (f *FileSystem) convertGoToLua(l *lua.LState, val any) lua.LValue {
	if val == nil {
		return lua.LNil
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.String:
		return lua.LString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(v.Uint())
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(v.Float())
	case reflect.Bool:
		return lua.LBool(v.Bool())
	case reflect.Slice, reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return lua.LString(v.Bytes())
		}
		tbl := l.NewTable()
		for i := 0; i < v.Len(); i++ {
			l.RawSetInt(tbl, i+1, f.convertGoToLua(l, v.Index(i).Interface()))
		}
		return tbl
	case reflect.Map:
		tbl := l.NewTable()
		for _, key := range v.MapKeys() {
			k := f.convertGoToLua(l, key.Interface())
			val := f.convertGoToLua(l, v.MapIndex(key).Interface())
			l.SetTable(tbl, k, val)
		}
		return tbl
	default:
		return lua.LString(fmt.Sprintf("%v", val))
	}
}

func (f *FileSystem) convertLuaToGo(val lua.LValue) any {
	switch val.Type() {
	case lua.LTString:
		return val.String()
	case lua.LTNumber:
		return float64(val.(lua.LNumber))
	case lua.LTBool:
		return bool(val.(lua.LBool))
	case lua.LTTable:
		tbl := val.(*lua.LTable)
		maxN := tbl.MaxN()
		if maxN > 0 {
			res := make([]any, maxN)
			for i := 1; i <= maxN; i++ {
				res[i-1] = f.convertLuaToGo(tbl.RawGetInt(i))
			}
			return res
		}
		res := make(map[string]any)
		tbl.ForEach(func(k, v lua.LValue) {
			res[k.String()] = f.convertLuaToGo(v)
		})
		return res
	}
	return nil
}

func (f *FileSystem) Eval(script string, _ int, keys []string, args []string) (any, error) {
	L := f.initLuaState()
	defer L.Close()

	keysTable := L.NewTable()
	for i, k := range keys {
		L.RawSetInt(keysTable, i+1, lua.LString(k))
	}
	L.SetGlobal("KEYS", keysTable)

	argsTable := L.NewTable()
	for i, a := range args {
		L.RawSetInt(argsTable, i+1, lua.LString(a))
	}
	L.SetGlobal("ARGV", argsTable)

	err := L.DoString(script)
	if err != nil {
		return nil, err
	}

	ret := L.Get(-1)
	if ret == lua.LNil {
		return nil, nil
	}
	return f.convertLuaToGo(ret), nil
}

func (f *FileSystem) EvalSha(sha1Str string, numkeys int, keys []string, args []string) (any, error) {
	path := getScriptPath(f.pathRoot, sha1Str)
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("NOSCRIPT No matching script. Please use EVAL")
		}
		return nil, err
	}
	return f.Eval(string(b), numkeys, keys, args)
}

func (f *FileSystem) EvalRO(script string, numkeys int, keys []string, args []string) (any, error) {
	return f.Eval(script, numkeys, keys, args)
}

func (f *FileSystem) EvalShaRO(sha1Str string, numkeys int, keys []string, args []string) (any, error) {
	return f.EvalSha(sha1Str, numkeys, keys, args)
}

func (f *FileSystem) ScriptExists(scripts ...string) ([]bool, error) {
	dir := getScriptsDir(f.pathRoot)
	res := make([]bool, len(scripts))
	for i, hash := range scripts {
		path := filepath.Join(dir, hash+".lua")
		_, err := os.Stat(path)
		res[i] = !os.IsNotExist(err)
	}
	return res, nil
}

func (f *FileSystem) ScriptFlush(_ ...dotpip.ScriptFlushOption) error {
	dir := getScriptsDir(f.pathRoot)
	return os.RemoveAll(dir)
}

func (f *FileSystem) ScriptLoad(script string) (string, error) {
	hash := sha1.Sum([]byte(script))
	hashStr := hex.EncodeToString(hash[:])

	dir := getScriptsDir(f.pathRoot)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	path := filepath.Join(dir, hashStr+".lua")
	if err := os.WriteFile(path, []byte(script), 0644); err != nil {
		return "", err
	}

	return hashStr, nil
}

func (f *FileSystem) ScriptKill() error {
	return nil
}
