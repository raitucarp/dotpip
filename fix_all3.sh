#!/bin/bash

sed -i 's/type VectorSetElement struct {/\/\/ VectorSetElement represents an element in a vector set.\ntype VectorSetElement struct {/' vectors.go
sed -i 's/type VAddOption func(\*VAddOptions)/\/\/ VAddOption configures VAddOptions.\ntype VAddOption func(\*VAddOptions)/' vectors.go
sed -i 's/type VAddOptions struct {/\/\/ VAddOptions represents options for VADD.\ntype VAddOptions struct {/' vectors.go
sed -i 's/func WithVAddReduceDim(dim int) VAddOption {/\/\/ WithVAddReduceDim sets the dimension reduction for VADD.\nfunc WithVAddReduceDim(dim int) VAddOption {/' vectors.go
sed -i 's/func WithVAddCas(cas bool) VAddOption {/\/\/ WithVAddCas sets CAS for VADD.\nfunc WithVAddCas(cas bool) VAddOption {/' vectors.go
sed -i 's/func WithVAddQuant(quant string) VAddOption {/\/\/ WithVAddQuant sets quantization for VADD.\nfunc WithVAddQuant(quant string) VAddOption {/' vectors.go
sed -i 's/func WithVAddEF(ef int) VAddOption {/\/\/ WithVAddEF sets EF for VADD.\nfunc WithVAddEF(ef int) VAddOption {/' vectors.go
sed -i 's/func WithVAddM(m int) VAddOption {/\/\/ WithVAddM sets M for VADD.\nfunc WithVAddM(m int) VAddOption {/' vectors.go
sed -i 's/type VSimOption func(\*VSimOptions)/\/\/ VSimOption configures VSimOptions.\ntype VSimOption func(\*VSimOptions)/' vectors.go
sed -i 's/type VSimOptions struct {/\/\/ VSimOptions represents options for VSIM.\ntype VSimOptions struct {/' vectors.go
sed -i 's/func WithVSimWithScores(withScores bool) VSimOption {/\/\/ WithVSimWithScores sets WITHSCORES for VSIM.\nfunc WithVSimWithScores(withScores bool) VSimOption {/' vectors.go
sed -i 's/func WithVSimWithAttribs(withAttribs bool) VSimOption {/\/\/ WithVSimWithAttribs sets WITHATTRIBS for VSIM.\nfunc WithVSimWithAttribs(withAttribs bool) VSimOption {/' vectors.go
sed -i 's/func WithVSimCount(count int) VSimOption {/\/\/ WithVSimCount sets COUNT for VSIM.\nfunc WithVSimCount(count int) VSimOption {/' vectors.go
sed -i 's/func WithVSimEpsilon(epsilon float64) VSimOption {/\/\/ WithVSimEpsilon sets EPSILON for VSIM.\nfunc WithVSimEpsilon(epsilon float64) VSimOption {/' vectors.go
sed -i 's/func WithVSimEF(ef int) VSimOption {/\/\/ WithVSimEF sets EF for VSIM.\nfunc WithVSimEF(ef int) VSimOption {/' vectors.go
sed -i 's/func WithVSimFilter(filter string) VSimOption {/\/\/ WithVSimFilter sets FILTER for VSIM.\nfunc WithVSimFilter(filter string) VSimOption {/' vectors.go
sed -i 's/func WithVSimFilterEF(filterEF int) VSimOption {/\/\/ WithVSimFilterEF sets FILTER_EF for VSIM.\nfunc WithVSimFilterEF(filterEF int) VSimOption {/' vectors.go
sed -i 's/func WithVSimTruth(truth bool) VSimOption {/\/\/ WithVSimTruth sets TRUTH for VSIM.\nfunc WithVSimTruth(truth bool) VSimOption {/' vectors.go
sed -i 's/func WithVSimNoThread(noThread bool) VSimOption {/\/\/ WithVSimNoThread sets NOTHREAD for VSIM.\nfunc WithVSimNoThread(noThread bool) VSimOption {/' vectors.go
sed -i 's/type VSimResult struct {/\/\/ VSimResult represents a result from VSIM.\ntype VSimResult struct {/' vectors.go
sed -i 's/type VRandMemberOption func(\*VRandMemberOptions)/\/\/ VRandMemberOption configures VRandMemberOptions.\ntype VRandMemberOption func(\*VRandMemberOptions)/' vectors.go
sed -i 's/type VRandMemberOptions struct {/\/\/ VRandMemberOptions options for VRANDMEMBER.\ntype VRandMemberOptions struct {/' vectors.go
sed -i 's/type VRangeOption func(\*VRangeOptions)/\/\/ VRangeOption configures VRangeOptions.\ntype VRangeOption func(\*VRangeOptions)/' vectors.go
sed -i 's/type VRangeOptions struct {/\/\/ VRangeOptions options for VRANGE.\ntype VRangeOptions struct {/' vectors.go
