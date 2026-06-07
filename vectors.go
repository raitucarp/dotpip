package dotpip

// VectorSetElement represents an element in a vector set.
type VectorSetElement struct {
	Vector []float32
	Element string
	Attributes map[string]any
}

// VAddOption configures VAddOptions.
type VAddOption func(*VAddOptions)

// VAddOptions represents options for VADD.
type VAddOptions struct {
	ReduceDim *int
	Cas       bool
	Quant     string // NOQUANT, Q8, BIN
	EF        *int
	M         *int
}

// WithVAddReduceDim sets the dimension reduction for VADD.
func WithVAddReduceDim(dim int) VAddOption {
	return func(o *VAddOptions) {
		o.ReduceDim = &dim
	}
}

// WithVAddCas sets CAS for VADD.
func WithVAddCas(cas bool) VAddOption {
	return func(o *VAddOptions) {
		o.Cas = cas
	}
}

// WithVAddQuant sets quantization for VADD.
func WithVAddQuant(quant string) VAddOption {
	return func(o *VAddOptions) {
		o.Quant = quant
	}
}

// WithVAddEF sets EF for VADD.
func WithVAddEF(ef int) VAddOption {
	return func(o *VAddOptions) {
		o.EF = &ef
	}
}

// WithVAddM sets M for VADD.
func WithVAddM(m int) VAddOption {
	return func(o *VAddOptions) {
		o.M = &m
	}
}

// VSimOption configures VSimOptions.
type VSimOption func(*VSimOptions)

// VSimOptions represents options for VSIM.
type VSimOptions struct {
	WithScores   bool
	WithAttribs  bool
	Count        *int
	Epsilon      *float64
	EF           *int
	Filter       *string
	FilterEF     *int
	Truth        bool
	NoThread     bool
}

// WithVSimWithScores sets WITHSCORES for VSIM.
func WithVSimWithScores(withScores bool) VSimOption {
	return func(o *VSimOptions) {
		o.WithScores = withScores
	}
}

// WithVSimWithAttribs sets WITHATTRIBS for VSIM.
func WithVSimWithAttribs(withAttribs bool) VSimOption {
	return func(o *VSimOptions) {
		o.WithAttribs = withAttribs
	}
}

// WithVSimCount sets COUNT for VSIM.
func WithVSimCount(count int) VSimOption {
	return func(o *VSimOptions) {
		o.Count = &count
	}
}

// WithVSimEpsilon sets EPSILON for VSIM.
func WithVSimEpsilon(epsilon float64) VSimOption {
	return func(o *VSimOptions) {
		o.Epsilon = &epsilon
	}
}

// WithVSimEF sets EF for VSIM.
func WithVSimEF(ef int) VSimOption {
	return func(o *VSimOptions) {
		o.EF = &ef
	}
}

// WithVSimFilter sets FILTER for VSIM.
func WithVSimFilter(filter string) VSimOption {
	return func(o *VSimOptions) {
		o.Filter = &filter
	}
}

// WithVSimFilterEF sets FILTER_EF for VSIM.
func WithVSimFilterEF(filterEF int) VSimOption {
	return func(o *VSimOptions) {
		o.FilterEF = &filterEF
	}
}

// WithVSimTruth sets TRUTH for VSIM.
func WithVSimTruth(truth bool) VSimOption {
	return func(o *VSimOptions) {
		o.Truth = truth
	}
}

// WithVSimNoThread sets NOTHREAD for VSIM.
func WithVSimNoThread(noThread bool) VSimOption {
	return func(o *VSimOptions) {
		o.NoThread = noThread
	}
}

// VSimResult represents a result from VSIM.
type VSimResult struct {
	Element    string
	Score      *float64
	Attributes *string
}

// VRandMemberOption configures VRandMemberOptions.
type VRandMemberOption func(*VRandMemberOptions)

// VRandMemberOptions options for VRANDMEMBER.
type VRandMemberOptions struct {
}

// VRangeOption configures VRangeOptions.
type VRangeOption func(*VRangeOptions)

// VRangeOptions options for VRANGE.
type VRangeOptions struct {
}
