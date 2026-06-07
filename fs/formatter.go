package fs

import "dotpip"

// Formatter sets the formatter for FileSystem.
func (f *FileSystem) Formatter(formatter dotpip.DataTypeFormatter) {
	f.formatter = &formatter
}
