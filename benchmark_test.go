package main

import (
	"bytes"
	"reflect"
	"testing"
)

var parseBenchmarksTests = map[string]struct {
	resultSet          string
	expectedBenchmarks []benchmark
	expectErr          bool
}{
	"1_bench_2_subs": {
		resultSet: `
			goos: darwin
			goarch: amd64
			BenchmarkMath/areaUnder/y=sin(x)/delta=0.001000/start_x=-2/end_x=1/abs_val=true-4         	   21801	     55357 ns/op	       0 B/op	       0 allocs/op
			BenchmarkMath/areaUnder/y=2x+3/delta=1.000000/start_x=-1/end_x=2/abs_val=false-4          	88335925	        13.3 ns/op	       0 B/op	       0 allocs/op
			BenchmarkMath/max/y=2x+3/delta=0.001000/start_x=-2/end_x=1-4                              	   56282	     20361 ns/op	       0 B/op	       0 allocs/op
			BenchmarkMath/max/y=sin(x)/delta=1.000000/start_x=-1/end_x=2-4                            	16381138	        62.7 ns/op	       0 B/op	       0 allocs/op
			PASS
			`,
		expectedBenchmarks: []benchmark{
			{
				name: "BenchmarkMath",
				results: []benchRes{
					{
						inputs: benchInputs{
							subs: []string{"areaUnder"},
							varValues: []benchVarValue{
								{name: "y", value: "sin(x)"},
								{name: "delta", value: 0.001},
								{name: "start_x", value: -2},
								{name: "end_x", value: 1},
								{name: "abs_val", value: true},
							},
						},
						outputs: benchOutputs{N: 21801, NsPerOp: 55357},
					},
					{
						inputs: benchInputs{
							subs: []string{"areaUnder"},
							varValues: []benchVarValue{
								{name: "y", value: "2x+3"},
								{name: "delta", value: 1.0},
								{name: "start_x", value: -1},
								{name: "end_x", value: 2},
								{name: "abs_val", value: false},
							},
						},
						outputs: benchOutputs{N: 88335925, NsPerOp: 13.3},
					},
					{
						inputs: benchInputs{
							subs: []string{"max"},
							varValues: []benchVarValue{
								{name: "y", value: "2x+3"},
								{name: "delta", value: 0.001},
								{name: "start_x", value: -2},
								{name: "end_x", value: 1},
							},
						},
						outputs: benchOutputs{N: 56282, NsPerOp: 20361},
					},
					{
						inputs: benchInputs{
							subs: []string{"max"},
							varValues: []benchVarValue{
								{name: "y", value: "sin(x)"},
								{name: "delta", value: 1.0},
								{name: "start_x", value: -1},
								{name: "end_x", value: 2},
							},
						},
						outputs: benchOutputs{N: 16381138, NsPerOp: 62.7},
					},
				},
			},
		},
	},
}

func TestParseBencharks(t *testing.T) {
	for testName, testCase := range parseBenchmarksTests {
		t.Run(testName, func(t *testing.T) {
			b := bytes.NewReader([]byte(testCase.resultSet))
			benchmarks, err := parseBenchmarks(b)
			if err != nil {
				if !testCase.expectErr {
					t.Errorf("unexpected error: %s", err)
				}
				return
			}

			if testCase.expectErr {
				t.Fatalf("unexpectedly no error")
			}

			if !reflect.DeepEqual(benchmarks, testCase.expectedBenchmarks) {
				t.Errorf("unexpected parsed benchmarks\nexpected:\n%v\nactual:\n%v", testCase.expectedBenchmarks, benchmarks)
			}
		})
	}
}

var groupResultsTests = map[string]struct {
	benchmark              benchmark
	groupBy                []string
	expectedGroupedResults map[string][]benchRes
}{
	"group_by_1_string_var": {
		benchmark: benchmark{
			name: "BenchmarkMath",
			results: []benchRes{
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
		groupBy: []string{"y"},
		expectedGroupedResults: map[string][]benchRes{
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
	},
	"group_by_2_vars": {
		benchmark: benchmark{
			name: "BenchmarkMath",
			results: []benchRes{
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
		groupBy: []string{"y", "delta"},
		expectedGroupedResults: map[string][]benchRes{
			"y=sin(x),delta=0.001": []benchRes{{
				inputs: benchInputs{
					subs: []string{"areaUnder"},
					varValues: []benchVarValue{
						{name: "y", value: "sin(x)"},
						{name: "delta", value: 0.001},
					},
				},
				outputs: benchOutputs{N: 10, NsPerOp: 2000},
			}},
			"y=sin(x),delta=0.01": []benchRes{{
				inputs: benchInputs{
					subs: []string{"areaUnder"},
					varValues: []benchVarValue{
						{name: "y", value: "sin(x)"},
						{name: "delta", value: 0.01},
					},
				},
				outputs: benchOutputs{N: 100, NsPerOp: 200},
			}},
			"y=2x+3,delta=0.001": []benchRes{{
				inputs: benchInputs{
					subs: []string{"areaUnder"},
					varValues: []benchVarValue{
						{name: "y", value: "2x+3"},
						{name: "delta", value: 0.001},
					},
				},
				outputs: benchOutputs{N: 5, NsPerOp: 1000},
			}},
			"y=2x+3,delta=0.01": []benchRes{{
				inputs: benchInputs{
					subs: []string{"areaUnder"},
					varValues: []benchVarValue{
						{name: "y", value: "2x+3"},
						{name: "delta", value: 0.01},
					},
				},
				outputs: benchOutputs{N: 10, NsPerOp: 100},
			}},
		},
	},
}

func TestGroupResults(t *testing.T) {
	for testName, testCase := range groupResultsTests {
		t.Run(testName, func(t *testing.T) {
			grouped := testCase.benchmark.groupResults(testCase.groupBy)
			if !reflect.DeepEqual(grouped, testCase.expectedGroupedResults) {
				t.Errorf("unexpected grouped results\nexpected:\n%v\nactual:\n%v", testCase.expectedGroupedResults, grouped)
			}
		})
	}
}
