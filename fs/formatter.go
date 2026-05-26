package fs

import "dotpip"

func (f *fileSystem) Formatter(formatter dotpip.DataTypeFormatter) {
	f.formatter = &formatter
}
