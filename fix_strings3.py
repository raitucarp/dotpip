with open("fs/strings.go", "r") as f:
    content = f.read()

content = content.replace('func (f *FileSystem) Append(key dotpip.Key, value string) (appendedString int) {', 'func (f *FileSystem) Append(key dotpip.Key, value string) (appendedString int) {\n\tf.subMutex.Lock()\n\tf.suppressSetEvent = True\n\tf.subMutex.Unlock()\n\tdefer func() { f.subMutex.Lock(); f.suppressSetEvent = False; f.subMutex.Unlock(); if appendedString > 0 { f.emitKeyspaceEvent(key, "append", \'$\') } }()\n')

# But we need it for IncrBy, IncrByFloat, SetRange
content = content.replace('func (f *FileSystem) IncrBy(key dotpip.Key, increment int) (int, error) {', 'func (f *FileSystem) IncrBy(key dotpip.Key, increment int) (ret int, err error) {\n\tf.subMutex.Lock()\n\tf.suppressSetEvent = True\n\tf.subMutex.Unlock()\n\tdefer func() { f.subMutex.Lock(); f.suppressSetEvent = False; f.subMutex.Unlock(); if err == nil { f.emitKeyspaceEvent(key, "incrby", \'$\') } }()\n')

content = content.replace('func (f *FileSystem) IncrByFloat(key dotpip.Key, increment float64) (float64, error) {', 'func (f *FileSystem) IncrByFloat(key dotpip.Key, increment float64) (ret float64, err error) {\n\tf.subMutex.Lock()\n\tf.suppressSetEvent = True\n\tf.subMutex.Unlock()\n\tdefer func() { f.subMutex.Lock(); f.suppressSetEvent = False; f.subMutex.Unlock(); if err == nil { f.emitKeyspaceEvent(key, "incrbyfloat", \'$\') } }()\n')

content = content.replace('func (f *FileSystem) SetRange(key dotpip.Key, offset int, value string) (int, error) {', 'func (f *FileSystem) SetRange(key dotpip.Key, offset int, value string) (ret int, err error) {\n\tf.subMutex.Lock()\n\tf.suppressSetEvent = True\n\tf.subMutex.Unlock()\n\tdefer func() { f.subMutex.Lock(); f.suppressSetEvent = False; f.subMutex.Unlock(); if err == nil { f.emitKeyspaceEvent(key, "setrange", \'$\') } }()\n')

with open("fs/strings.go", "w") as f:
    f.write(content)
