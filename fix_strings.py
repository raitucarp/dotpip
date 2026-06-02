import re

with open("fs/strings.go", "r") as f:
    content = f.read()

functions = [
    ("Append", "append", "$", "0, err"), # wait, Append calls Set, so Set will emit "set", but reference says APPEND generates "append". Oh! Wait.
]

# Let's check Set first
parts = content.split("func (f *FileSystem) Set(")
if len(parts) > 1:
    part0 = parts[0]
    subparts = parts[1].split('err = f.writeFileByKey(key, finalValue.([]byte))\n\tif err != nil {\n\t\treturn "", err\n\t}', 1)
    if len(subparts) > 1:
        parts[1] = subparts[0] + 'err = f.writeFileByKey(key, finalValue.([]byte))\n\tif err == nil {\n\t\tf.emitKeyspaceEvent(key, "set", \'$\')\n\t}\n\tif err != nil {\n\t\treturn "", err\n\t}' + subparts[1]

    content = part0 + "func (f *FileSystem) Set(" + parts[1]

with open("fs/strings.go", "w") as f:
    f.write(content)
