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
