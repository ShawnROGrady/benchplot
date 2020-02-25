package benchmark

import (
	"reflect"
	"testing"

	"github.com/ShawnROGrady/benchplot/plot"
	"github.com/ShawnROGrady/benchplot/plot/mock"
)

var sampleBenchmark = Benchmark{
	Name: "BenchmarkMath",
	results: []benchRes{
		{
			inputs: benchInputs{
				subs: []string{"areaUnder"},
				varValues: []benchVarValue{
					{name: "y", value: "sin(x)"},
					{name: "delta", value: 0.001},
					{name: "start_x", value: -2},
					{name: "end_x", value: 1},
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
					{name: "start_x", value: -2},
					{name: "end_x", value: 1},
				},
			},
			outputs: benchOutputs{N: 100, NsPerOp: 200},
		},
		{
			inputs: benchInputs{
				subs: []string{"areaUnder"},
				varValues: []benchVarValue{
					{name: "y", value: "2x+3"},
					{name: "delta", value: 0.001},
					{name: "start_x", value: -2},
					{name: "end_x", value: 1},
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
					{name: "start_x", value: -2},
					{name: "end_x", value: 1},
				},
			},
			outputs: benchOutputs{N: 10, NsPerOp: 100},
		},
	},
}

var plotScatterTests = map[string]struct {
	benchmark      Benchmark
	groupBy        []string
	xName          string
	yName          string
	expectedData   map[string]plot.NumericData
	expectedTitle  string
	expectedXLabel string
	expectedYLabel string
	expectErr      bool
}{
	"x=float64,y=float64": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"y"},
		xName:     "delta", yName: TimeName,
		expectedData: map[string]plot.NumericData{
			"y=sin(x)": plot.NumericData{
				X: []float64{0.001, 0.01},
				Y: []float64{2000, 200},
			},
			"y=2x+3": plot.NumericData{
				X: []float64{0.001, 0.01},
				Y: []float64{1000, 100},
			},
		},
		expectedTitle:  "BenchmarkMath",
		expectedXLabel: "delta",
		expectedYLabel: TimeName,
	},
	"x=int,y=int": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"y"},
		xName:     "start_x", yName: RunsName,
		expectedData: map[string]plot.NumericData{
			"y=sin(x)": plot.NumericData{
				X: []float64{-2, -2},
				Y: []float64{10, 100},
			},
			"y=2x+3": plot.NumericData{
				X: []float64{-2, -2},
				Y: []float64{5, 10},
			},
		},
		expectedTitle:  "BenchmarkMath",
		expectedXLabel: "start_x",
		expectedYLabel: RunsName,
	},
	"x=float64,y=uint64": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"y"},
		xName:     "delta", yName: NumAllocsName,
		expectedData: map[string]plot.NumericData{
			"y=sin(x)": plot.NumericData{
				X: []float64{0.001, 0.01},
				Y: []float64{0, 0},
			},
			"y=2x+3": plot.NumericData{
				X: []float64{0.001, 0.01},
				Y: []float64{0, 0},
			},
		},
		expectedTitle:  "BenchmarkMath",
		expectedXLabel: "delta",
		expectedYLabel: NumAllocsName,
	},
	"x=string,y=float64": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"start_x"},
		xName:     "y", yName: TimeName,
		expectErr: true,
	},
	"invalid_x_name": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"start_x"},
		xName:     "invalid_name", yName: TimeName,
		expectErr: true,
	},
	"invalid_y_name": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"start_x"},
		xName:     "delta", yName: "invalid_name",
		expectErr: true,
	},
}

func TestPlotScatter(t *testing.T) {
	for testName, testCase := range plotScatterTests {
		t.Run(testName, func(t *testing.T) {
			p := &mock.Plotter{
				PlotScatterFn: func(data map[string]plot.NumericData, title string, xLabel string, yLabel string) error {
					// validate args
					if !reflect.DeepEqual(data, testCase.expectedData) {
						t.Errorf("unexpected plot data\nexpected:\n%v\nactual:\n%v", testCase.expectedData, data)
					}
					if title != testCase.expectedTitle {
						t.Errorf("unexpected title\nexpected:\n%s\nactual:\n%s", testCase.expectedTitle, title)
					}
					if xLabel != testCase.expectedXLabel {
						t.Errorf("unexpected xLabel\nexpected:\n%s\nactual:\n%s", testCase.expectedXLabel, xLabel)
					}
					if yLabel != testCase.expectedYLabel {
						t.Errorf("unexpected yLabel\nexpected:\n%s\nactual:\n%s", testCase.expectedYLabel, yLabel)
					}
					return nil
				},
			}

			err := testCase.benchmark.PlotScatter(p, testCase.groupBy, testCase.xName, testCase.yName)
			if err != nil {
				if !testCase.expectErr {
					t.Errorf("unexpected error: %s", err)
				}
				return
			}
			if testCase.expectErr {
				t.Error("unexpectedly no error")
			}
		})
	}
}
