package plot

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ShawnROGrady/benchparse"
)

var sampleBenchmark = benchparse.Benchmark{
	Name: "BenchmarkMath",
	Results: []benchparse.BenchRes{
		{
			Inputs: benchparse.BenchInputs{
				Subs: []benchparse.BenchSub{{Name: "areaUnder"}},
				VarValues: []benchparse.BenchVarValue{
					{Name: "y", Value: "sin(x)"},
					{Name: "delta", Value: 0.001},
					{Name: "start_x", Value: -2},
					{Name: "end_x", Value: 1},
				},
			},
			Outputs: newTestOutputs(10, withNsPerOp(2000), withAllocsPerOp(0)),
		},
		{
			Inputs: benchparse.BenchInputs{
				Subs: []benchparse.BenchSub{{Name: "areaUnder"}},
				VarValues: []benchparse.BenchVarValue{
					{Name: "y", Value: "sin(x)"},
					{Name: "delta", Value: 0.01},
					{Name: "start_x", Value: -2},
					{Name: "end_x", Value: 1},
				},
			},
			Outputs: newTestOutputs(100, withNsPerOp(200), withAllocsPerOp(0)),
		},
		{
			Inputs: benchparse.BenchInputs{
				Subs: []benchparse.BenchSub{{Name: "areaUnder"}},
				VarValues: []benchparse.BenchVarValue{
					{Name: "y", Value: "2x+3"},
					{Name: "delta", Value: 0.001},
					{Name: "start_x", Value: -2},
					{Name: "end_x", Value: 1},
				},
			},
			Outputs: newTestOutputs(5, withNsPerOp(1000), withAllocsPerOp(0)),
		},
		{
			Inputs: benchparse.BenchInputs{
				Subs: []benchparse.BenchSub{{Name: "areaUnder"}},
				VarValues: []benchparse.BenchVarValue{
					{Name: "y", Value: "2x+3"},
					{Name: "delta", Value: 0.01},
					{Name: "start_x", Value: -2},
					{Name: "end_x", Value: 1},
				},
			},
			Outputs: newTestOutputs(10, withNsPerOp(100), withAllocsPerOp(0)),
		},
	},
}
var splitGroupedResultsTests = map[string]struct {
	groupedResults benchparse.GroupedResults
	xName          string
	yName          string
	expectedSplit  map[string][]splitRes
	expectErr      bool
}{
	"valid_x_valid_y": {
		groupedResults: benchparse.GroupedResults{
			"y=sin(x)": []benchparse.BenchRes{
				sampleBenchmark.Results[0],
				sampleBenchmark.Results[1],
			},
			"y=2x+3": []benchparse.BenchRes{
				sampleBenchmark.Results[2],
				sampleBenchmark.Results[3],
			},
		},
		xName: "delta", yName: TimeName,
		expectedSplit: map[string][]splitRes{
			"y=sin(x)": []splitRes{{x: 0.001, y: float64(2000)}, {x: 0.01, y: float64(200)}},
			"y=2x+3":   []splitRes{{x: 0.001, y: float64(1000)}, {x: 0.01, y: float64(100)}},
		},
	},
	"invalid_x_name": {
		groupedResults: benchparse.GroupedResults{
			"": sampleBenchmark.Results,
		},
		xName: "invalid_var", yName: TimeName,
		expectErr: true,
	},
	"invalid_y_name": {
		groupedResults: benchparse.GroupedResults{
			"": sampleBenchmark.Results,
		},
		xName: "delta", yName: "invalid_var",
		expectErr: true,
	},
	"y_not_measured": {
		groupedResults: benchparse.GroupedResults{
			"": sampleBenchmark.Results,
		},
		xName: "delta", yName: AllocMBytesRate,
		expectErr: true,
	},
}

func TestSplitGroupedResults(t *testing.T) {
	for testName, testCase := range splitGroupedResultsTests {
		t.Run(testName, func(t *testing.T) {
			splitGrouped, err := splitGroupedResult(testCase.groupedResults, testCase.xName, testCase.yName)
			if err != nil {
				if !testCase.expectErr {
					t.Errorf("unexpected error: %s", err)
				}
				return
			}

			if testCase.expectErr {
				t.Errorf("unexpectedly no error")
			}

			if !reflect.DeepEqual(splitGrouped, testCase.expectedSplit) {
				t.Errorf("unexpected split grouped results\nexpected:\n%#v\nactual:\n%#v", testCase.expectedSplit, splitGrouped)
			}
		})
	}
}

var splitBenchResErr error

func BenchmarkSplitBenchRes(b *testing.B) {
	var (
		allNumValues = []int{1, 2, 3, 4, 5, 10, 20}
	)

	for _, numValues := range allNumValues {
		b.Run(fmt.Sprintf("num_values=%d", numValues), func(b *testing.B) {
			varValues := make([]benchparse.BenchVarValue, numValues)
			for i := 0; i < numValues; i++ {
				varValues[i] = benchparse.BenchVarValue{
					Name:  fmt.Sprintf("var%d", i),
					Value: i,
				}
			}
			res := benchparse.BenchRes{
				Inputs:  benchparse.BenchInputs{VarValues: varValues},
				Outputs: newTestOutputs(10, withNsPerOp(100)),
			}

			var err error
			for i := 0; i < b.N; i++ {
				_, err = splitBenchRes(res, "var0", TimeName)
				if err != nil {
					b.Fatalf("unexpected error: %s", err)
				}
			}
			splitBenchResErr = err
		})
	}
}

var benchOutputValByNameTests = map[string]struct {
	output                       benchparse.BenchOutputs
	expectedIterations           int
	expectedNsPerOp              float64
	expectedNsPerOpErr           error
	expectedAllocedBytesPerOp    uint64
	expectedAllocedBytesPerOpErr error
	expectedAllocsPerOp          uint64
	expectedAllocsPerOpErr       error
	expectedMBPerS               float64
	expectedMBPerSErr            error
}{
	"all_set": {
		output: newTestOutputs(
			21801,
			withNsPerOp(55357),
			withAllocedBytesPerOp(4321),
			withAllocsPerOp(21),
			withMBPerS(0.12),
		),
		expectedIterations:        21801,
		expectedNsPerOp:           55357,
		expectedAllocedBytesPerOp: 4321,
		expectedAllocsPerOp:       21,
		expectedMBPerS:            0.12,
	},
	"benchmem_not_set_with_set_bytes": {
		output: newTestOutputs(
			21801,
			withNsPerOp(55357),
			withMBPerS(0.12),
		),
		expectedIterations:           21801,
		expectedNsPerOp:              55357,
		expectedAllocedBytesPerOpErr: benchparse.ErrNotMeasured,
		expectedAllocsPerOpErr:       benchparse.ErrNotMeasured,
		expectedMBPerS:               0.12,
	},
	"benchmem_set_but_no_allocs": {
		output: newTestOutputs(
			21801,
			withNsPerOp(55357),
			withAllocedBytesPerOp(0),
			withAllocsPerOp(0),
		),
		expectedIterations:        21801,
		expectedNsPerOp:           55357,
		expectedAllocedBytesPerOp: 0,
		expectedAllocsPerOp:       0,
		expectedMBPerSErr:         benchparse.ErrNotMeasured,
	},
	"none_set": {
		output:                       newTestOutputs(1),
		expectedIterations:           1,
		expectedNsPerOpErr:           benchparse.ErrNotMeasured,
		expectedAllocedBytesPerOpErr: benchparse.ErrNotMeasured,
		expectedAllocsPerOpErr:       benchparse.ErrNotMeasured,
		expectedMBPerSErr:            benchparse.ErrNotMeasured,
	},
}

func TestGetOutputValByName(t *testing.T) {
	for testName, testCase := range benchOutputValByNameTests {
		t.Run(testName, func(t *testing.T) {
			t.Run("num_iterations", func(t *testing.T) {
				num, err := benchOutputValByName(testCase.output, RunsName)
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				if num != testCase.expectedIterations {
					t.Errorf("unexpected iterations (expected=%d, actual=%d)", testCase.expectedIterations, num)
				}
			})
			t.Run("ns_per_op", func(t *testing.T) {
				testNsPerOp(t, testCase.output, testCase.expectedNsPerOp, testCase.expectedNsPerOpErr)
			})
			t.Run("allocated_bytes_per_op", func(t *testing.T) {
				testAllocedBytesPerOp(t, testCase.output, testCase.expectedAllocedBytesPerOp, testCase.expectedAllocedBytesPerOpErr)
			})
			t.Run("allocs_per_op", func(t *testing.T) {
				testAllocsPerOp(t, testCase.output, testCase.expectedAllocsPerOp, testCase.expectedAllocsPerOpErr)
			})
			t.Run("MB_per_s", func(t *testing.T) {
				testMBPerS(t, testCase.output, testCase.expectedMBPerS, testCase.expectedMBPerSErr)
			})
		})
	}
}

func testNsPerOp(t *testing.T, b benchparse.BenchOutputs, expectedV float64, expectedErr error) {
	t.Helper()
	ns, err := benchOutputValByName(b, TimeName)
	if err != nil {
		if expectedErr != nil {
			if err != expectedErr {
				t.Errorf("unexpected error received (expected=%s, actual=%s)", expectedErr, err)
			}
		} else {
			t.Errorf("unexpected error: %s", err)
		}
		return
	}

	if expectedErr != nil {
		t.Errorf("unexpectedly no error")
	}

	if expectedV != ns {
		t.Errorf("unexpected NsPerOp (expected=%v, actual=%v)", expectedV, ns)
	}
}

func testAllocedBytesPerOp(t *testing.T, b benchparse.BenchOutputs, expectedV uint64, expectedErr error) {
	t.Helper()
	v, err := benchOutputValByName(b, AllocBytesName)
	if err != nil {
		if expectedErr != nil {
			if err != expectedErr {
				t.Errorf("unexpected error received (expected=%s, actual=%s)", expectedErr, err)
			}
		} else {
			t.Errorf("unexpected error: %s", err)
		}
		return
	}

	if expectedErr != nil {
		t.Errorf("unexpectedly no error")
	}

	if expectedV != v {
		t.Errorf("unexpected AllocedBytesPerOp (expected=%v, actual=%v)", expectedV, v)
	}
}

func testAllocsPerOp(t *testing.T, b benchparse.BenchOutputs, expectedV uint64, expectedErr error) {
	t.Helper()
	v, err := benchOutputValByName(b, NumAllocsName)
	if err != nil {
		if expectedErr != nil {
			if err != expectedErr {
				t.Errorf("unexpected error received (expected=%s, actual=%s)", expectedErr, err)
			}
		} else {
			t.Errorf("unexpected error: %s", err)
		}
		return
	}

	if expectedErr != nil {
		t.Errorf("unexpectedly no error")
	}

	if expectedV != v {
		t.Errorf("unexpected AllocsPerOp (expected=%v, actual=%v)", expectedV, v)
	}
}

func testMBPerS(t *testing.T, b benchparse.BenchOutputs, expectedV float64, expectedErr error) {
	t.Helper()
	v, err := benchOutputValByName(b, AllocMBytesRate)
	if err != nil {
		if expectedErr != nil {
			if err != expectedErr {
				t.Errorf("unexpected error received (expected=%s, actual=%s)", expectedErr, err)
			}
		} else {
			t.Errorf("unexpected error: %s", err)
		}
		return
	}

	if expectedErr != nil {
		t.Errorf("unexpectedly no error")
	}

	if expectedV != v {
		t.Errorf("unexpected MBPerS (expected=%v, actual=%v)", expectedV, v)
	}
}
