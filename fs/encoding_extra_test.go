package fs

import (
	"dotpip"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtraEncodingCoverage(t *testing.T) {
	for _, format := range []FileEncodeType{JSON, YAML, TOML, RAW} {
		f := FileSystem{
			encodeType: format,
		}

		// bitmapDecode coverage
		encBitmap, err := f.bitmapEncode([]uint{1, 2, 3})
		assert.NoError(t, err)
		decBitmap, err := f.bitmapDecode(encBitmap)
		assert.NoError(t, err)
		assert.Equal(t, []uint{1, 2, 3}, decBitmap)
	}
}

func TestExtraEncodingCoverage2(t *testing.T) {
	for _, format := range []FileEncodeType{JSON, YAML, TOML, RAW} {
		f := FileSystem{
			encodeType: format,
		}

		enc, err := f.stringEncode("test")
		assert.NoError(t, err)
		dec, err := f.stringDecode(enc)
		assert.NoError(t, err)
		assert.Equal(t, "test", dec)

		enc2, err := f.arrayEncode([]string{"a", "b"})
		assert.NoError(t, err)
		dec2, err := f.arrayDecode(enc2)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, dec2)

		enc3, err := f.listEncode([]any{"a", "b"})
		assert.NoError(t, err)
		dec3, err := f.listDecode(enc3)
		assert.NoError(t, err)
		assert.Equal(t, "a", dec3[0])

		enc4, err := f.hashEncode(map[string]string{"a": "b"})
		assert.NoError(t, err)
		dec4, err := f.hashDecode(enc4)
		assert.NoError(t, err)
		assert.Equal(t, "b", dec4["a"])

		enc5, err := f.setEncode(map[string]any{"a": struct{}{}})
		assert.NoError(t, err)
		dec5, err := f.setDecode(enc5)
		assert.NoError(t, err)
		_, ok := dec5["a"]
		assert.True(t, ok)

		enc6, err := f.sortedSetEncode(map[string]float64{"a": 1.0})
		assert.NoError(t, err)
		dec6, err := f.sortedSetDecode(enc6)
		assert.NoError(t, err)
		assert.Equal(t, 1.0, dec6["a"])

		enc7, err := f.hyperLogLogEncode([]byte("test"))
		assert.NoError(t, err)
		dec7, err := f.hyperLogLogDecode(enc7)
		assert.NoError(t, err)
		assert.Equal(t, []byte("test"), dec7)

		streamData := dotpip.Stream{
			Entries: []dotpip.StreamEntry{
				{ID: "1-0", Values: map[string]string{"k": "v"}},
			},
		}
		enc8, err := f.streamEncode(streamData)
		assert.NoError(t, err)
		dec8, err := f.streamDecode(enc8)
		assert.NoError(t, err)
		assert.Equal(t, "1-0", dec8.Entries[0].ID)

		enc9, err := f.JSONEncode(map[string]any{"a": "b"})
		assert.NoError(t, err)
		dec9, err := f.JSONDecode(enc9)
		assert.NoError(t, err)
		m := dec9.(map[string]any)
		assert.Equal(t, "b", m["a"])
	}
}

func TestVectorEncoding(t *testing.T) {
	for _, format := range []FileEncodeType{JSON, YAML, TOML, RAW} {
		f := FileSystem{
			encodeType: format,
		}

		vs := make(map[string]dotpip.VectorSetElement)
		vs["test"] = dotpip.VectorSetElement{
			Vector:  []float32{1.0, 2.0},
			Element: "test",
		}

		encoded, err := f.vectorSetEncode(vs)
		assert.NoError(t, err)

		decoded, err := f.vectorSetDecode(encoded)
		assert.NoError(t, err)
		assert.Equal(t, vs["test"].Element, decoded["test"].Element)

		// Invalid decode
		_, err = f.vectorSetDecode([]byte("invalid"))
		assert.Error(t, err)
	}

	// Invalid format
	f := FileSystem{
		encodeType: FileEncodeType("INVALID"),
	}
	_, err := f.vectorSetEncode(nil)
	assert.Error(t, err)

	_, err = f.vectorSetDecode([]byte(""))
	assert.Error(t, err)
}
