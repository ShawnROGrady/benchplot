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
				PlotScatterFn: func(data map[string]plot.NumericData, title string, xLabel string, yLabel string, includeLegend bool) error {
					// validate args
					if !includeLegend {
						t.Errorf("unexpectedly not including legend")
					}
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

			err := testCase.benchmark.Plot(p, testCase.xName, testCase.yName, WithGroupBy(testCase.groupBy), WithPlotTypes([]string{ScatterType}))
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

var plotAvgLineTests = map[string]struct {
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
	"x=float64,y=float64,no_dups": {
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
	"x=int,y=int,2_dups": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"y"},
		xName:     "start_x", yName: RunsName,
		expectedData: map[string]plot.NumericData{
			"y=sin(x)": plot.NumericData{
				X: []float64{-2},
				Y: []float64{55},
			},
			"y=2x+3": plot.NumericData{
				X: []float64{-2},
				Y: []float64{7.5},
			},
		},
		expectedTitle:  "BenchmarkMath",
		expectedXLabel: "start_x",
		expectedYLabel: RunsName,
	},
	"x=int,y=int,4_dups": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"end_x"},
		xName:     "start_x", yName: RunsName,
		expectedData: map[string]plot.NumericData{
			"end_x=1": plot.NumericData{
				X: []float64{-2},
				Y: []float64{31.25},
			},
		},
		expectedTitle:  "BenchmarkMath",
		expectedXLabel: "start_x",
		expectedYLabel: RunsName,
	},
	"x=float64,y=int,4_dups": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"end_x"},
		xName:     "delta", yName: RunsName,
		expectedData: map[string]plot.NumericData{
			"end_x=1": plot.NumericData{
				X: []float64{0.001, 0.01},
				Y: []float64{7.5, 55},
			},
		},
		expectedTitle:  "BenchmarkMath",
		expectedXLabel: "delta",
		expectedYLabel: RunsName,
	},
	"x=float64,y=uint64,no_dups": {
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

func TestPlotAvgLine(t *testing.T) {
	for testName, testCase := range plotAvgLineTests {
		t.Run(testName, func(t *testing.T) {
			p := &mock.Plotter{
				PlotLineFn: func(data map[string]plot.NumericData, title string, xLabel string, yLabel string, includeLegend bool) error {
					// validate args
					if !includeLegend {
						t.Errorf("unexpectedly not including legend")
					}
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

			err := testCase.benchmark.Plot(p, testCase.xName, testCase.yName, WithGroupBy(testCase.groupBy), WithPlotTypes([]string{AvgLineType}))
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

type plotFnInput struct {
	data          map[string]plot.NumericData
	includeLegend bool
	title         string
	xLabel        string
	yLabel        string
}

var plotTests = map[string]struct {
	benchmark            Benchmark
	groupBy              []string
	plots                []string
	xName                string
	yName                string
	expectedScatterInput plotFnInput
	expectedLineInput    plotFnInput
	expectErr            bool
}{
	"x=float64,avg_line+scatter": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"y"},
		plots:     []string{ScatterType, AvgLineType},
		xName:     "delta", yName: TimeName,
		expectedScatterInput: plotFnInput{
			data: map[string]plot.NumericData{
				"y=sin(x)": plot.NumericData{
					X: []float64{0.001, 0.01},
					Y: []float64{2000, 200},
				},
				"y=2x+3": plot.NumericData{
					X: []float64{0.001, 0.01},
					Y: []float64{1000, 100},
				},
			},
			title:         "BenchmarkMath",
			xLabel:        "delta",
			yLabel:        TimeName,
			includeLegend: true,
		},
		expectedLineInput: plotFnInput{
			data: map[string]plot.NumericData{
				"y=sin(x)": plot.NumericData{
					X: []float64{0.001, 0.01},
					Y: []float64{2000, 200},
				},
				"y=2x+3": plot.NumericData{
					X: []float64{0.001, 0.01},
					Y: []float64{1000, 100},
				},
			},
			title:         "BenchmarkMath",
			xLabel:        "delta",
			yLabel:        TimeName,
			includeLegend: false,
		},
	},
	"x=float64,default_plots": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"y"},
		xName:     "delta", yName: TimeName,
		expectedScatterInput: plotFnInput{
			data: map[string]plot.NumericData{
				"y=sin(x)": plot.NumericData{
					X: []float64{0.001, 0.01},
					Y: []float64{2000, 200},
				},
				"y=2x+3": plot.NumericData{
					X: []float64{0.001, 0.01},
					Y: []float64{1000, 100},
				},
			},
			title:         "BenchmarkMath",
			xLabel:        "delta",
			yLabel:        TimeName,
			includeLegend: true,
		},
		expectedLineInput: plotFnInput{
			data: map[string]plot.NumericData{
				"y=sin(x)": plot.NumericData{
					X: []float64{0.001, 0.01},
					Y: []float64{2000, 200},
				},
				"y=2x+3": plot.NumericData{
					X: []float64{0.001, 0.01},
					Y: []float64{1000, 100},
				},
			},
			title:         "BenchmarkMath",
			xLabel:        "delta",
			yLabel:        TimeName,
			includeLegend: false,
		},
	},
	"invalid_plot_type": {
		benchmark: sampleBenchmark,
		groupBy:   []string{"y"},
		xName:     "delta", yName: TimeName,
		plots:     []string{"invalid"},
		expectErr: true,
	},
}

func TestPlot(t *testing.T) {
	for testName, testCase := range plotTests {
		t.Run(testName, func(t *testing.T) {
			p := &mock.Plotter{
				PlotScatterFn: func(data map[string]plot.NumericData, title string, xLabel string, yLabel string, includeLegend bool) error {
					// validate args
					if includeLegend != testCase.expectedScatterInput.includeLegend {
						t.Errorf("unexpected includeLegend\nexpected:%t\nactual:%t", testCase.expectedScatterInput.includeLegend, includeLegend)
					}
					if !reflect.DeepEqual(data, testCase.expectedScatterInput.data) {
						t.Errorf("unexpected plot data\nexpected:\n%v\nactual:\n%v", testCase.expectedScatterInput.data, data)
					}
					if title != testCase.expectedScatterInput.title {
						t.Errorf("unexpected title\nexpected:\n%s\nactual:\n%s", testCase.expectedScatterInput.title, title)
					}
					if xLabel != testCase.expectedScatterInput.xLabel {
						t.Errorf("unexpected xLabel\nexpected:\n%s\nactual:\n%s", testCase.expectedScatterInput.xLabel, xLabel)
					}
					if yLabel != testCase.expectedScatterInput.yLabel {
						t.Errorf("unexpected yLabel\nexpected:\n%s\nactual:\n%s", testCase.expectedScatterInput.yLabel, yLabel)
					}
					return nil
				},
				PlotLineFn: func(data map[string]plot.NumericData, title string, xLabel string, yLabel string, includeLegend bool) error {
					// validate args
					if includeLegend != testCase.expectedLineInput.includeLegend {
						t.Errorf("unexpected includeLegend\nexpected:%t\nactual:%t", testCase.expectedLineInput.includeLegend, includeLegend)
					}
					if !reflect.DeepEqual(data, testCase.expectedLineInput.data) {
						t.Errorf("unexpected plot data\nexpected:\n%v\nactual:\n%v", testCase.expectedLineInput.data, data)
					}
					if title != testCase.expectedLineInput.title {
						t.Errorf("unexpected title\nexpected:\n%s\nactual:\n%s", testCase.expectedLineInput.title, title)
					}
					if xLabel != testCase.expectedLineInput.xLabel {
						t.Errorf("unexpected xLabel\nexpected:\n%s\nactual:\n%s", testCase.expectedLineInput.xLabel, xLabel)
					}
					if yLabel != testCase.expectedLineInput.yLabel {
						t.Errorf("unexpected yLabel\nexpected:\n%s\nactual:\n%s", testCase.expectedLineInput.yLabel, yLabel)
					}
					return nil
				},
			}

			err := testCase.benchmark.Plot(p, testCase.xName, testCase.yName, WithGroupBy(testCase.groupBy), WithPlotTypes(testCase.plots))
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
