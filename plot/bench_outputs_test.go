package plot

import "github.com/ShawnROGrady/benchparse"

// benchOutputs is a mock implementation of benchparse.BenchOutputs
type benchOutputs struct {
	GetIterationsFn        func() int
	GetNsPerOpFn           func() (float64, error)
	GetAllocedBytesPerOpFn func() (uint64, error)
	GetAllocsPerOpFn       func() (uint64, error)
	GetMBPerSFn            func() (float64, error)
}

// GetIterations returns _m.GetIterationsFn
func (_m *benchOutputs) GetIterations() int {
	return _m.GetIterationsFn()
}

// GetNsPerOp returns _m.GetNsPerOpFn
func (_m *benchOutputs) GetNsPerOp() (float64, error) {
	return _m.GetNsPerOpFn()
}

// GetAllocedBytesPerOp returns _m.GetAllocedBytesPerOpFn
func (_m *benchOutputs) GetAllocedBytesPerOp() (uint64, error) {
	return _m.GetAllocedBytesPerOpFn()
}

// GetAllocsPerOp returns _m.GetAllocsPerOpFn
func (_m *benchOutputs) GetAllocsPerOp() (uint64, error) {
	return _m.GetAllocsPerOpFn()
}

// GetMBPerS returns _m.GetMBPerSFn
func (_m *benchOutputs) GetMBPerS() (float64, error) {
	return _m.GetMBPerSFn()
}

type benchOutOption interface {
	apply(b *benchOutputs)
}

type withNsPerOp float64

func (w withNsPerOp) apply(b *benchOutputs) {
	b.GetNsPerOpFn = func() (float64, error) {
		return float64(w), nil
	}
}

type withAllocedBytesPerOp uint64

func (w withAllocedBytesPerOp) apply(b *benchOutputs) {
	b.GetAllocedBytesPerOpFn = func() (uint64, error) {
		return uint64(w), nil
	}
}

type withAllocsPerOp uint64

func (w withAllocsPerOp) apply(b *benchOutputs) {
	b.GetAllocsPerOpFn = func() (uint64, error) {
		return uint64(w), nil
	}
}

type withMBPerS float64

func (w withMBPerS) apply(b *benchOutputs) {
	b.GetMBPerSFn = func() (float64, error) {
		return float64(w), nil
	}
}

func newTestOutputs(n int, opts ...benchOutOption) benchparse.BenchOutputs {
	b := &benchOutputs{
		GetIterationsFn:        func() int { return n },
		GetNsPerOpFn:           func() (float64, error) { return 0, benchparse.ErrNotMeasured },
		GetAllocedBytesPerOpFn: func() (uint64, error) { return 0, benchparse.ErrNotMeasured },
		GetAllocsPerOpFn:       func() (uint64, error) { return 0, benchparse.ErrNotMeasured },
		GetMBPerSFn:            func() (float64, error) { return 0, benchparse.ErrNotMeasured },
	}
	for _, opt := range opts {
		opt.apply(b)
	}
	return b
}
