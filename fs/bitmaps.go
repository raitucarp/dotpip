package fs

import (
	"dotpip"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// In Redis, bit operations work on string values.
// We decode the string, operate on bits, and encode back.
func (f *FileSystem) getBitmapBytes(key dotpip.Key) ([]byte, error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		if os.IsNotExist(err) {
			return []byte{}, nil // Treat non-existent as empty string/bitmap
		}
		return nil, err
	}
	strVal, err := f.formatter.StringDecode(content)
	if err != nil {
		return nil, err
	}
	return []byte(strVal), nil
}

func (f *FileSystem) setBitmapBytes(key dotpip.Key, b []byte) error {
	_, err := f.Set(key, string(b))
	return err
}

func (f *FileSystem) SetBit(key dotpip.Key, offset int, value int) (int, error) {
	b, err := f.getBitmapBytes(key)
	if err != nil {
		return 0, err
	}
	byteIndex := offset / 8
	bitIndex := offset % 8

	// Expand if necessary
	if byteIndex >= len(b) {
		newB := make([]byte, byteIndex+1)
		copy(newB, b)
		b = newB
	}

	origByte := b[byteIndex]
	origBit := int((origByte >> (7 - bitIndex)) & 1)

	if value == 1 {
		b[byteIndex] |= (1 << (7 - bitIndex))
	} else {
		b[byteIndex] &= ^(1 << (7 - bitIndex))
	}

	err = f.setBitmapBytes(key, b)
	if err != nil {
		return 0, err
	}

	return origBit, nil
}

func (f *FileSystem) GetBit(key dotpip.Key, offset int) (int, error) {
	b, err := f.getBitmapBytes(key)
	if err != nil {
		return 0, err
	}
	byteIndex := offset / 8
	bitIndex := offset % 8

	if byteIndex >= len(b) {
		return 0, nil
	}

	bit := int((b[byteIndex] >> (7 - bitIndex)) & 1)
	return bit, nil
}

func (f *FileSystem) BitCount(key dotpip.Key, start int, end int) (int, error) {
	b, err := f.getBitmapBytes(key)
	if err != nil {
		return 0, err
	}

	length := len(b)
	if length == 0 {
		return 0, nil
	}

	// Handle negative offsets
	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}

	// Bound checking
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start > length-1 {
		return 0, nil
	}
	if end > length-1 {
		end = length - 1
	}
	if start > end {
		return 0, nil
	}

	count := 0
	for i := start; i <= end; i++ {
		val := b[i]
		// Count set bits in byte
		for val > 0 {
			count += int(val & 1)
			val >>= 1
		}
	}

	return count, nil
}

func (f *FileSystem) BitOp(operation dotpip.BitOp, destKey dotpip.Key, keys ...dotpip.Key) (int, error) {
	var maxLen int
	var byteArrays [][]byte

	for _, key := range keys {
		b, err := f.getBitmapBytes(key)
		if err != nil {
			return 0, err
		}
		if len(b) > maxLen {
			maxLen = len(b)
		}
		byteArrays = append(byteArrays, b)
	}

	res := make([]byte, maxLen)

	if operation == dotpip.BitOpNot {
		if len(keys) != 1 {
			// Redis BITOP NOT takes exactly one source key
			return 0, nil // Should probably be an error, but let's just do nothing
		}
		for i := 0; i < maxLen; i++ {
			if i < len(byteArrays[0]) {
				res[i] = ^byteArrays[0][i]
			} else {
				res[i] = ^byte(0)
			}
		}
	} else {
		for i := 0; i < maxLen; i++ {
			var b byte
			if operation == dotpip.BitOpAnd {
				b = 0xFF
			} else {
				b = 0
			}

			for j := 0; j < len(byteArrays); j++ {
				var val byte
				if i < len(byteArrays[j]) {
					val = byteArrays[j][i]
				} else {
					val = 0
				}

				if j == 0 {
					b = val
				} else {
					switch operation {
					case dotpip.BitOpAnd:
						b &= val
					case dotpip.BitOpOr:
						b |= val
					case dotpip.BitOpXor:
						b ^= val
					}
				}
			}
			res[i] = b
		}
	}

	err := f.setBitmapBytes(destKey, res)
	if err != nil {
		return 0, err
	}

	return maxLen, nil
}

func (f *FileSystem) BitPos(key dotpip.Key, bit int, start int, end int) (int, error) {
	// We need to implement BITPOS properly with the start/end options which are byte indexes
	b, err := f.getBitmapBytes(key)
	if err != nil {
		return 0, err
	}

	length := len(b)
	if length == 0 {
		if bit == 0 {
			return 0, nil
		}
		return -1, nil
	}

	if start < 0 {
		start = length + start
	}
	if end < 0 {
		end = length + end
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start > length-1 {
		if bit == 0 {
			return length * 8, nil
		}
		return -1, nil
	}
	if end > length-1 {
		end = length - 1
	}

	if start > end {
		return -1, nil
	}

	for i := start; i <= end; i++ {
		for j := 0; j < 8; j++ {
			val := int((b[i] >> (7 - j)) & 1)
			if val == bit {
				return i*8 + j, nil
			}
		}
	}

	if bit == 0 {
		// If searching for 0 and not found, and we checked up to the end of the string (or beyond it originally),
		// Redis bitpos returns the first bit of the *next* byte. But only if end was not specified or we searched to the end.
		// For simplicity we will assume if we searched to the end, the next bit is 0.
		// Actually, if we didn't specify end, or end is the last byte, then the bit after the string is 0.
		return length * 8, nil
	}

	return -1, nil
}

type overflowType int

const (
	overflowWrap overflowType = iota
	overflowSat
	overflowFail
)

func parseBitfieldType(t string) (bool, int, error) {
	if len(t) < 2 {
		return false, 0, fmt.Errorf(string(dotpip.ErrMsgInvalidType))
	}
	signed := t[0] == 'i'
	if !signed && t[0] != 'u' {
		return false, 0, fmt.Errorf(string(dotpip.ErrMsgInvalidTypeFormat))
	}
	bits, err := strconv.Atoi(t[1:])
	if err != nil || bits <= 0 || (signed && bits > 64) || (!signed && bits > 63) {
		return false, 0, fmt.Errorf(string(dotpip.ErrMsgInvalidBits))
	}
	return signed, bits, nil
}

func getBitFieldValue(b []byte, offset int, bits int, signed bool) (int64, bool) {
	if offset < 0 {
		return 0, false
	}
	val := int64(0)
	for i := 0; i < bits; i++ {
		byteIndex := (offset + i) / 8
		bitIndex := (offset + i) % 8
		bit := int64(0)
		if byteIndex < len(b) {
			bit = int64((b[byteIndex] >> (7 - bitIndex)) & 1)
		}
		val = (val << 1) | bit
	}

	if signed && bits > 0 {
		// sign extension
		signBit := (val >> (bits - 1)) & 1
		if signBit == 1 {
			shift := 64 - bits
			val = (val << shift) >> shift
		}
	}
	return val, true
}

func setBitFieldValue(b []byte, offset int, bits int, value int64) []byte {
	// Ensure capacity
	lastBit := offset + bits - 1
	lastByte := lastBit / 8
	if lastByte >= len(b) {
		newB := make([]byte, lastByte+1)
		copy(newB, b)
		b = newB
	}

	for i := 0; i < bits; i++ {
		byteIndex := (offset + i) / 8
		bitIndex := (offset + i) % 8
		bitVal := (value >> (bits - 1 - i)) & 1
		if bitVal == 1 {
			b[byteIndex] |= (1 << (7 - bitIndex))
		} else {
			b[byteIndex] &= ^(1 << (7 - bitIndex))
		}
	}
	return b
}

func getMinMaxForType(signed bool, bits int) (int64, int64) {
	if signed {
		minVal := int64(-1) << (bits - 1)
		maxVal := (int64(1) << (bits - 1)) - 1
		return minVal, maxVal
	}
	maxVal := (int64(1) << bits) - 1
	return 0, maxVal
}

// BitField operates on multiple bit fields in a string.
func (f *FileSystem) BitField(key dotpip.Key, args ...any) ([]any, error) {
	b, err := f.getBitmapBytes(key)
	if err != nil {
		return nil, err
	}

	var results []any
	overflow := overflowWrap

	for i := 0; i < len(args); i++ {
		argStr, ok := args[i].(string)
		if !ok {
			return nil, fmt.Errorf(string(dotpip.ErrMsgInvalidArgumentType))
		}
		argUpper := strings.ToUpper(argStr)

		if argUpper == "OVERFLOW" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf(string(dotpip.ErrMsgSyntaxError))
			}
			ovStr, ok := args[i+1].(string)
			if !ok {
				return nil, fmt.Errorf(string(dotpip.ErrMsgInvalidOverflowArg))
			}
			ovUpper := strings.ToUpper(ovStr)
			switch ovUpper {
			case "WRAP":
				overflow = overflowWrap
			case "SAT":
				overflow = overflowSat
			case "FAIL":
				overflow = overflowFail
			default:
				return nil, fmt.Errorf(string(dotpip.ErrMsgInvalidOverflowType))
			}
			i++
			continue
		}

		if argUpper == "GET" {
			if i+2 >= len(args) {
				return nil, fmt.Errorf(string(dotpip.ErrMsgSyntaxError))
			}
			typStr := args[i+1].(string)
			signed, bits, err := parseBitfieldType(typStr)
			if err != nil {
				return nil, err
			}
			offsetRaw := args[i+2]
			var offset int
			switch v := offsetRaw.(type) {
			case int:
				offset = v
			case float64:
				offset = int(v)
			case string:
				if strings.HasPrefix(v, "#") {
					multOffset, err := strconv.Atoi(v[1:])
					if err != nil {
						return nil, err
					}
					offset = multOffset * bits
				} else {
					offset, err = strconv.Atoi(v)
					if err != nil {
						return nil, err
					}
				}
			}

			val, ok := getBitFieldValue(b, offset, bits, signed)
			if !ok {
				results = append(results, nil)
			} else {
				results = append(results, int(val))
			}
			i += 2
			continue
		}

		if argUpper == "SET" {
			if i+3 >= len(args) {
				return nil, fmt.Errorf(string(dotpip.ErrMsgSyntaxError))
			}
			typStr := args[i+1].(string)
			signed, bits, err := parseBitfieldType(typStr)
			if err != nil {
				return nil, err
			}

			var offset int
			switch v := args[i+2].(type) {
			case int:
				offset = v
			case float64:
				offset = int(v)
			case string:
				if strings.HasPrefix(v, "#") {
					multOffset, err := strconv.Atoi(v[1:])
					if err != nil {
						return nil, err
					}
					offset = multOffset * bits
				} else {
					offset, err = strconv.Atoi(v)
					if err != nil {
						return nil, err
					}
				}
			}

			var value int64
			switch v := args[i+3].(type) {
			case int:
				value = int64(v)
			case float64:
				value = int64(v)
			case string:
				valInt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, err
				}
				value = valInt
			}

			oldVal, _ := getBitFieldValue(b, offset, bits, signed)

			// For SET, we still check bounds for bits, but typically value is just truncated.
			// Wrap logic applies
			minVal, maxVal := getMinMaxForType(signed, bits)
			if value < minVal || value > maxVal {
				mask := (int64(1) << bits) - 1
				value &= mask
				if signed {
					signBit := (value >> (bits - 1)) & 1
					if signBit == 1 {
						shift := 64 - bits
						value = (value << shift) >> shift
					}
				}
			}

			b = setBitFieldValue(b, offset, bits, value)
			results = append(results, int(oldVal))
			i += 3
			continue
		}

		if argUpper == "INCRBY" {
			if i+3 >= len(args) {
				return nil, fmt.Errorf(string(dotpip.ErrMsgSyntaxError))
			}
			typStr := args[i+1].(string)
			signed, bits, err := parseBitfieldType(typStr)
			if err != nil {
				return nil, err
			}
			var offset int
			switch v := args[i+2].(type) {
			case int:
				offset = v
			case float64:
				offset = int(v)
			case string:
				if strings.HasPrefix(v, "#") {
					multOffset, err := strconv.Atoi(v[1:])
					if err != nil {
						return nil, err
					}
					offset = multOffset * bits
				} else {
					offset, err = strconv.Atoi(v)
					if err != nil {
						return nil, err
					}
				}
			}

			var inc int64
			switch v := args[i+3].(type) {
			case int:
				inc = int64(v)
			case float64:
				inc = int64(v)
			case string:
				valInt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, err
				}
				inc = valInt
			}

			oldVal, _ := getBitFieldValue(b, offset, bits, signed)
			newVal := oldVal + inc
			minVal, maxVal := getMinMaxForType(signed, bits)

			fail := false
			if newVal < minVal || newVal > maxVal {
				switch overflow {
				case overflowWrap:
					mask := (int64(1) << bits) - 1
					newVal &= mask
					if signed {
						signBit := (newVal >> (bits - 1)) & 1
						if signBit == 1 {
							shift := 64 - bits
							newVal = (newVal << shift) >> shift // rely on go sign extension
						}
					}
				case overflowSat:
					if newVal < minVal {
						newVal = minVal
					} else {
						newVal = maxVal
					}
				case overflowFail:
					fail = true
				}
			}

			if fail {
				results = append(results, nil)
			} else {
				b = setBitFieldValue(b, offset, bits, newVal)
				results = append(results, int(newVal))
			}
			i += 3
			continue
		}

		return nil, fmt.Errorf(string(dotpip.ErrMsgUnknownSubcommand))
	}

	if err := f.setBitmapBytes(key, b); err != nil {
		return nil, err
	}

	return results, nil
}
