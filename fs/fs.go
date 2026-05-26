package fs

import (
	"dotpip"
)

type fileSystem struct {
	pathRoot   string
	formatter  *dotpip.DataTypeFormatter
	encodeType FileEncodeType
}

func FileSystem(pathRoot string) *fileSystem {
	f := fileSystem{
		pathRoot:   pathRoot,
		formatter:  &dotpip.DataTypeFormatter{},
		encodeType: JSON,
	}

	f.formatter.StringEncode = f.stringEncode
	f.formatter.StringDecode = f.stringDecode

	return &f
}
