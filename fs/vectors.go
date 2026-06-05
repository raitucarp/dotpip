package fs

import (
	"dotpip"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"

	"github.com/viterin/vek"
)

func (f *FileSystem) readVectorSet(key dotpip.Key) (map[string]dotpip.VectorSetElement, error) {
	content, err := f.readFileByKey(key)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]dotpip.VectorSetElement), nil
		}
		return nil, err
	}

	return f.formatter.VectorSetDecode(content)
}

func (f *FileSystem) writeVectorSet(key dotpip.Key, vectorSet map[string]dotpip.VectorSetElement) error {
	if len(vectorSet) == 0 {
		err := f.removeFileByKey(key)
		if err == nil {
			f.emitKeyspaceEvent(key, "del", 'g')
		}
		return err
	}

	encoded, err := f.formatter.VectorSetEncode(vectorSet)
	if err != nil {
		return err
	}

	return f.writeFileByKey(key, encoded.([]byte))
}

func (f *FileSystem) VAdd(key dotpip.Key, element string, vector []float32, options ...dotpip.VAddOption) (int, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return 0, err
	}

	if len(vs) > 0 {
		// check dim
		for _, el := range vs {
			if len(el.Vector) != len(vector) {
				return 0, errors.New("ERR vector dimension mismatch")
			}
			break
		}
	}

	_, exists := vs[element]
	vs[element] = dotpip.VectorSetElement{
		Vector: vector,
		Element: element,
	}

	err = f.writeVectorSet(key, vs)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 1, nil
	}
	return 0, nil
}

func (f *FileSystem) VCard(key dotpip.Key) (int, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return 0, err
	}
	return len(vs), nil
}

func (f *FileSystem) VDim(key dotpip.Key) (int, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return 0, err
	}
	for _, el := range vs {
		return len(el.Vector), nil
	}
	return 0, errors.New("ERR no elements")
}

func (f *FileSystem) VEmb(key dotpip.Key, elements ...string) ([][]float32, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return nil, err
	}
	res := make([][]float32, len(elements))
	for i, el := range elements {
		if e, ok := vs[el]; ok {
			res[i] = e.Vector
		} else {
			res[i] = nil
		}
	}
	return res, nil
}

func (f *FileSystem) VGetAttr(key dotpip.Key, elements ...string) ([]*string, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return nil, err
	}
	res := make([]*string, len(elements))
	for i, el := range elements {
		if e, ok := vs[el]; ok && e.Attributes != nil {
			b, _ := json.Marshal(e.Attributes)
			s := string(b)
			res[i] = &s
		} else {
			res[i] = nil
		}
	}
	return res, nil
}

func (f *FileSystem) VInfo(key dotpip.Key) (map[string]any, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return nil, err
	}
	if len(vs) == 0 {
		return nil, errors.New("ERR no such key")
	}
	dim := 0
	for _, el := range vs {
		dim = len(el.Vector)
		break
	}
	return map[string]any{
		"dimensions": dim,
		"count": len(vs),
	}, nil
}

func (f *FileSystem) VIsMember(key dotpip.Key, elements ...string) ([]bool, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return nil, err
	}
	res := make([]bool, len(elements))
	for i, el := range elements {
		_, res[i] = vs[el]
	}
	return res, nil
}

func (f *FileSystem) VLinks(key dotpip.Key, element string) (map[int][]string, error) {
	return make(map[int][]string), nil
}

func (f *FileSystem) VRandMember(key dotpip.Key, count int, options ...dotpip.VRandMemberOption) ([]string, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(vs))
	for el := range vs {
		res = append(res, el)
	}
	if count > len(res) {
		count = len(res)
	}
	return res[:count], nil
}

func (f *FileSystem) VRange(key dotpip.Key, start string, end string, options ...dotpip.VRangeOption) ([]string, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(vs))
	for el := range vs {
		res = append(res, el)
	}
	sort.Strings(res)

	filtered := make([]string, 0)
	for _, el := range res {
		if strings.Compare(el, start) >= 0 && strings.Compare(el, end) <= 0 {
			filtered = append(filtered, el)
		}
	}

	return filtered, nil
}

func (f *FileSystem) VRem(key dotpip.Key, elements ...string) (int, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, el := range elements {
		if _, ok := vs[el]; ok {
			delete(vs, el)
			count++
		}
	}
	if count > 0 {
		err = f.writeVectorSet(key, vs)
	}
	return count, err
}

func (f *FileSystem) VSetAttr(key dotpip.Key, element string, attributes string) (int, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return 0, err
	}
	e, ok := vs[element]
	if !ok {
		return 0, errors.New("ERR unknown element")
	}
	var attr map[string]any
	if err := json.Unmarshal([]byte(attributes), &attr); err != nil {
		return 0, err
	}
	e.Attributes = attr
	vs[element] = e
	err = f.writeVectorSet(key, vs)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

func (f *FileSystem) VSim(key dotpip.Key, reference any, options ...dotpip.VSimOption) ([]dotpip.VSimResult, error) {
	vs, err := f.readVectorSet(key)
	if err != nil {
		return nil, err
	}

	var refVector []float32
	switch ref := reference.(type) {
	case string:
		e, ok := vs[ref]
		if !ok {
			return nil, errors.New("ERR unknown element")
		}
		refVector = e.Vector
	case []float32:
		refVector = ref
	default:
		return nil, errors.New("ERR invalid reference type")
	}

	res := make([]dotpip.VSimResult, 0, len(vs))

    // Perform bulk conversion once for the reference vector
    ref64 := vek.FromFloat32(refVector)

	for _, el := range vs {
        el64 := vek.FromFloat32(el.Vector)
		sim := vek.CosineSimilarity(ref64, el64)
		score := float64(sim)
		res = append(res, dotpip.VSimResult{
			Element: el.Element,
			Score: &score,
		})
	}
	sort.Slice(res, func(i, j int) bool {
		return *res[i].Score > *res[j].Score
	})

	return res, nil
}
