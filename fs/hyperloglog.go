package fs

import (
	"dotpip"

	"github.com/axiomhq/hyperloglog"
)

func (f *FileSystem) PFAdd(key dotpip.Key, elements ...string) (int, error) {
	// 1. Get or create the sketch
	var sk *hyperloglog.Sketch
	content, err := f.readFileByKey(key)
	if err != nil {
		sk = hyperloglog.New14()
	} else {
		b, decErr := f.formatter.HyperLogLogDecode(content)
		if decErr != nil {
			sk = hyperloglog.New14() // Or return error? Redis returns WRONGTYPE. We can assume if error, it's not a valid HLL. Actually let's assume valid or create new if not exists.
		} else {
			sk = hyperloglog.New14()
			err = sk.UnmarshalBinary(b)
			if err != nil {
				sk = hyperloglog.New14()
			}
		}
	}

	oldBinary, _ := sk.MarshalBinary()

	for _, element := range elements {
		sk.Insert([]byte(element))
	}

	newBinary, _ := sk.MarshalBinary()

	changed := 0
	if string(oldBinary) != string(newBinary) {
		changed = 1
	}

	encoded, encErr := f.formatter.HyperLogLogEncode(newBinary)
	if encErr != nil {
		return 0, encErr
	}

	writeErr := f.writeFileByKey(key, encoded.([]byte))
	if writeErr != nil {
		return 0, writeErr
	}

	return changed, nil
}

func (f *FileSystem) PFCount(keys ...dotpip.Key) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	mergedSketch := hyperloglog.New14()

	for _, key := range keys {
		content, err := f.readFileByKey(key)
		if err == nil {
			b, decErr := f.formatter.HyperLogLogDecode(content)
			if decErr == nil {
				sk := hyperloglog.New14()
				err = sk.UnmarshalBinary(b)
				if err == nil {
					_ = mergedSketch.Merge(sk)
				}
			}
		}
	}

	return int(mergedSketch.Estimate()), nil
}

func (f *FileSystem) PFMerge(destKey dotpip.Key, sourceKeys ...dotpip.Key) error {
	destSketch := hyperloglog.New14()
	content, err := f.readFileByKey(destKey)
	if err == nil {
		b, decErr := f.formatter.HyperLogLogDecode(content)
		if decErr == nil {
			sk := hyperloglog.New14()
			if err = sk.UnmarshalBinary(b); err == nil {
				destSketch = sk
			}
		}
	}

	for _, sourceKey := range sourceKeys {
		content, err := f.readFileByKey(sourceKey)
		if err == nil {
			b, decErr := f.formatter.HyperLogLogDecode(content)
			if decErr == nil {
				sk := hyperloglog.New14()
				if err = sk.UnmarshalBinary(b); err == nil {
					_ = destSketch.Merge(sk)
				}
			}
		}
	}

	b, _ := destSketch.MarshalBinary()
	encoded, encErr := f.formatter.HyperLogLogEncode(b)
	if encErr != nil {
		return encErr
	}

	return f.writeFileByKey(destKey, encoded.([]byte))
}
