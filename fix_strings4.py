with open("fs/strings.go", "r") as f:
    content = f.read()

content = content.replace("f.emitKeyspaceEventFunc", "f.emitKeyspaceEvent")

with open("fs/strings.go", "w") as f:
    f.write(content)

with open("fs/fs.go", "r") as f:
    content = f.read()

content = content.replace('func (f *FileSystem) emitKeyspaceEvent(key dotpip.Key, event string, typeChar rune) {', 'func (f *FileSystem) emitKeyspaceEvent(key dotpip.Key, event string, typeChar rune) {\n\tif key == nil {\n\t\treturn\n\t}')

with open("fs/fs.go", "w") as f:
    f.write(content)
