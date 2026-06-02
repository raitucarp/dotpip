package fs

import "dotpip"

func (f *FileSystem) Formatter(formatter dotpip.DataTypeFormatter) {
	f.formatter = &formatter
}
