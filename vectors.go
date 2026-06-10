package dotpip

type VectorSetElement struct {
	Vector     []float32
	Element    string
	Attributes map[string]any
}

type VAddOption func(*VAddOptions)

type VAddOptions struct {
	ReduceDim *int
	Cas       bool
	Quant     string // NOQUANT, Q8, BIN
	EF        *int
	M         *int
}

func WithVAddReduceDim(dim int) VAddOption {
	return func(o *VAddOptions) {
		o.ReduceDim = &dim
	}
}

func WithVAddCas(cas bool) VAddOption {
	return func(o *VAddOptions) {
		o.Cas = cas
	}
}

func WithVAddQuant(quant string) VAddOption {
	return func(o *VAddOptions) {
		o.Quant = quant
	}
}

func WithVAddEF(ef int) VAddOption {
	return func(o *VAddOptions) {
		o.EF = &ef
	}
}

func WithVAddM(m int) VAddOption {
	return func(o *VAddOptions) {
		o.M = &m
	}
}

type VSimOption func(*VSimOptions)

type VSimOptions struct {
	WithScores  bool
	WithAttribs bool
	Count       *int
	Epsilon     *float64
	EF          *int
	Filter      *string
	FilterEF    *int
	Truth       bool
	NoThread    bool
}

func WithVSimWithScores(withScores bool) VSimOption {
	return func(o *VSimOptions) {
		o.WithScores = withScores
	}
}

func WithVSimWithAttribs(withAttribs bool) VSimOption {
	return func(o *VSimOptions) {
		o.WithAttribs = withAttribs
	}
}

func WithVSimCount(count int) VSimOption {
	return func(o *VSimOptions) {
		o.Count = &count
	}
}

func WithVSimEpsilon(epsilon float64) VSimOption {
	return func(o *VSimOptions) {
		o.Epsilon = &epsilon
	}
}

func WithVSimEF(ef int) VSimOption {
	return func(o *VSimOptions) {
		o.EF = &ef
	}
}

func WithVSimFilter(filter string) VSimOption {
	return func(o *VSimOptions) {
		o.Filter = &filter
	}
}

func WithVSimFilterEF(filterEF int) VSimOption {
	return func(o *VSimOptions) {
		o.FilterEF = &filterEF
	}
}

func WithVSimTruth(truth bool) VSimOption {
	return func(o *VSimOptions) {
		o.Truth = truth
	}
}

func WithVSimNoThread(noThread bool) VSimOption {
	return func(o *VSimOptions) {
		o.NoThread = noThread
	}
}

type VSimResult struct {
	Element    string
	Score      *float64
	Attributes *string
}

type VRandMemberOption func(*VRandMemberOptions)

type VRandMemberOptions struct {
}

type VRangeOption func(*VRangeOptions)

type VRangeOptions struct {
}
