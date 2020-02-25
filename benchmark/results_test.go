package benchmark

import (
	"fmt"
	"reflect"
	"testing"
)

var splitGroupedResultsTests = map[string]struct {
	groupedResults groupedResults
	xName          string
	yName          string
	expectedSplit  map[string][]splitRes
	expectErr      bool
}{
	"valid_x_valid_y": {
		groupedResults: groupedResults{
			"y=sin(x)": []benchRes{
				{
					inputs: benchInputs{
						subs: []string{"areaUnder"},
						varValues: []benchVarValue{
							{name: "y", value: "sin(x)"},
							{name: "delta", value: 0.001},
						},
					},
					outputs: benchOutputs{N: 10, NsPerOp: 2000},
				},
				{
					inputs: benchInputs{
						subs: []string{"areaUnder"},
						varValues: []benchVarValue{
							{name: "y", value: "sin(x)"},
							{name: "delta", value: 0.01},
						},
					},
					outputs: benchOutputs{N: 100, NsPerOp: 200},
				},
			},
			"y=2x+3": []benchRes{
				{
					inputs: benchInputs{
						subs: []string{"areaUnder"},
						varValues: []benchVarValue{
							{name: "y", value: "2x+3"},
							{name: "delta", value: 0.001},
						},
					},
					outputs: benchOutputs{N: 5, NsPerOp: 1000},
				},
				{
					inputs: benchInputs{
						subs: []string{"areaUnder"},
						varValues: []benchVarValue{
							{name: "y", value: "2x+3"},
							{name: "delta", value: 0.01},
						},
					},
					outputs: benchOutputs{N: 10, NsPerOp: 100},
				},
			},
		},
		xName: "delta", yName: TimeName,
		expectedSplit: map[string][]splitRes{
			"y=sin(x)": []splitRes{{x: 0.001, y: float64(2000)}, {x: 0.01, y: float64(200)}},
			"y=2x+3":   []splitRes{{x: 0.001, y: float64(1000)}, {x: 0.01, y: float64(100)}},
		},
	},
	"invalid_x": {
		groupedResults: groupedResults{
			"y=sin(x)": []benchRes{
				{
					inputs: benchInputs{
						subs: []string{"areaUnder"},
						varValues: []benchVarValue{
							{name: "y", value: "sin(x)"},
							{name: "delta", value: 0.001},
						},
					},
					outputs: benchOutputs{N: 10, NsPerOp: 2000},
				},
			},
		},
		xName: "invalid_var", yName: TimeName,
		expectErr: true,
	},
	"invalid_y": {
		groupedResults: groupedResults{
			"y=sin(x)": []benchRes{
				{
					inputs: benchInputs{
						subs: []string{"areaUnder"},
						varValues: []benchVarValue{
							{name: "y", value: "sin(x)"},
							{name: "delta", value: 0.001},
						},
					},
					outputs: benchOutputs{N: 10, NsPerOp: 2000},
				},
			},
		},
		xName: "delta", yName: "invalid_var",
		expectErr: true,
	},
}

func TestSplitGroupedResults(t *testing.T) {
	for testName, testCase := range splitGroupedResultsTests {
		t.Run(testName, func(t *testing.T) {
			splitGrouped, err := testCase.groupedResults.splitTo(testCase.xName, testCase.yName)
			if err != nil {
				if !testCase.expectErr {
					t.Errorf("unexpected error: %s", err)
				}
				return
			}

			if testCase.expectErr {
				t.Fatalf("unexpectedly no error")
			}

			if !reflect.DeepEqual(splitGrouped, testCase.expectedSplit) {
				t.Errorf("unexpected split grouped results\nexpected:\n%#v\nactual:\n%#v", testCase.expectedSplit, splitGrouped)
			}
		})
	}
}

var splitGroupedResultsErr error

func BenchmarkSplitGroupedResults(b *testing.B) {
	var (
		allNumValues = []int{1, 2, 3, 4, 5, 10, 20}
	)

	for _, numValues := range allNumValues {
		b.Run(fmt.Sprintf("num_values=%d", numValues), func(b *testing.B) {
			varValues := make([]benchVarValue, numValues)
			for i := 0; i < numValues; i++ {
				varValues[i] = benchVarValue{
					name:  fmt.Sprintf("var%d", i),
					value: i,
				}
			}
			res := benchRes{
				inputs:  benchInputs{varValues: varValues},
				outputs: benchOutputs{N: 10, NsPerOp: 100},
			}

			var err error
			for i := 0; i < b.N; i++ {
				_, err = res.splitTo("var0", TimeName)
				if err != nil {
					b.Fatalf("unexpected error: %s", err)
				}
			}
			splitGroupedResultsErr = err
		})
	}
}
