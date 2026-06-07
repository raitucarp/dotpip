package fs_test

import (
	"dotpip"
	"dotpip/fs"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTOML(t *testing.T) {
	encodings := []fs.FileEncodeType{fs.JSON, fs.YAML, fs.TOML, fs.RAW}
	for _, enc := range encodings {
		t.Run(string(enc), func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "dotpip_toml_test_")
			assert.NoError(t, err)
			defer func() { _ = os.RemoveAll(tmpDir) }()

			dotfs := fs.NewFileSystem(tmpDir)
			defer dotfs.Close()
			dotfs.EncodeType(enc)

			key := dotpip.NewKey("mytoml")

			res, err := dotfs.TOMLSet(key, "$", map[string]any{"a": 1, "b": []any{1, 2, 3}})
			assert.NoError(t, err)
			assert.Equal(t, string(dotpip.StatusOK), res)

			val, err := dotfs.TOMLGet(key, "$.b")
			assert.NoError(t, err)

			b, _ := json.Marshal(val)
			var v []any
			_ = json.Unmarshal(b, &v)

			assert.Equal(t, []any{float64(1), float64(2), float64(3)}, v)

			resArr, err := dotfs.TOMLArrAppend(key, "$.b", 4)
			assert.NoError(t, err)
			assert.Equal(t, []any{4}, resArr)
		})
	}
}
