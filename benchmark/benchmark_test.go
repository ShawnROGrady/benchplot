package benchmark

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/tools/benchmark/parse"
)

var parseBenchmarksTests = map[string]struct {
	resultSet          string
	expectedBenchmarks []Benchmark
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
		expectedBenchmarks: []Benchmark{
			{
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
			benchmarks, err := ParseBenchmarks(b)
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
	benchmark              Benchmark
	groupBy                []string
	expectedGroupedResults groupedResults
}{
	"group_by_1_string_var": {
		benchmark: Benchmark{
			Name: "BenchmarkMath",
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
	"no_group_by": {
		benchmark: Benchmark{
			Name: "BenchmarkMath",
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
			},
		},
		groupBy: []string{},
		expectedGroupedResults: map[string][]benchRes{
			"": []benchRes{
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
		},
	},
	"group_by_2_vars": {
		benchmark: Benchmark{
			Name: "BenchmarkMath",
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

var parseBenchmarksErr error

func BenchmarkParseBenchmarks(b *testing.B) {
	var (
		allNumBenchmarks     = []int{1, 2, 3, 4, 5}
		allCasesPerBenchmark = []int{5, 10, 15, 20, 25}
	)

	for _, numBenchmarks := range allNumBenchmarks {
		b.Run(fmt.Sprintf("num_benchmarks=%d", numBenchmarks), func(b *testing.B) {
			for _, casesPerBench := range allCasesPerBenchmark {
				b.Run(fmt.Sprintf("cases_per_bench=%d", casesPerBench), func(b *testing.B) {
					benchmarkParseBenchmarks(b, numBenchmarks, casesPerBench)
				})
			}
		})
	}
}

func benchmarkParseBenchmarks(b *testing.B, numBenchmarks, casesPerBench int) {
	b.Helper()
	newReader := func() io.Reader {
		var buf bytes.Buffer
		for i := 0; i < numBenchmarks; i++ {
			benchName := fmt.Sprintf("BenchmarkMethod%d", i)
			for j := 0; j < casesPerBench; j++ {
				bench := &parse.Benchmark{
					Name:    fmt.Sprintf("%s/var1=%d/var2=%d", benchName, j, j),
					N:       j,
					NsPerOp: float64(j),
				}
				if _, err := buf.WriteString(fmt.Sprintf("%s\n", bench)); err != nil {
					b.Fatalf("err constructing input: %s", err)
				}
			}
		}
		return &buf
	}

	var err error
	var benches []Benchmark
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		r := newReader()
		b.StartTimer()

		benches, err = ParseBenchmarks(r)
		if err != nil {
			b.Fatalf("unexpected error: %s", err)
		}
		if len(benches) != numBenchmarks {
			b.Fatalf("unexpected number of benchmarks (expected=%d, actual=%d)", numBenchmarks, len(benches))
		}
	}
	parseBenchmarksErr = err
}

var benchGroupResultsRes groupedResults

func BenchmarkGroupResults(b *testing.B) {
	var (
		allNumValues    = []int{1, 2, 3, 4, 5, 10, 20}
		allGroupByCount = []int{1, 2, 3}
		allNumResults   = []int{5, 10, 15, 20, 25}
	)

	for _, numResults := range allNumResults {
		b.Run(fmt.Sprintf("num_results=%d", numResults), func(b *testing.B) {
			for _, numValues := range allNumValues {
				b.Run(fmt.Sprintf("num_values=%d", numValues), func(b *testing.B) {
					for _, groupByCount := range allGroupByCount {
						b.Run(fmt.Sprintf("group_by_count=%d", groupByCount), func(b *testing.B) {
							if groupByCount > numValues {
								b.Skip("skipping due to groupByCount>numValues")
							}
							benchmarkGroupResults(b, numResults, numValues, groupByCount)
						})
					}
				})
			}
		})
	}
}

func benchmarkGroupResults(b *testing.B, numResults, numValues, groupByCount int) {
	b.Helper()
	var (
		bench   = newTestBenchmark(numResults, numValues)
		groupBy = make([]string, groupByCount)
	)
	for i := 0; i < groupByCount; i++ {
		groupBy[i] = fmt.Sprintf("var%d", i)
	}

	var groupRes groupedResults
	for i := 0; i < b.N; i++ {
		groupRes = bench.groupResults(groupBy)
	}
	benchGroupResultsRes = groupRes
}

func newTestBenchmark(numResults, numValues int) Benchmark {
	res := make([]benchRes, numResults)
	for i := 0; i < numResults; i++ {
		var (
			outputs   = benchOutputs{N: i, NsPerOp: float64(i)}
			varValues = make([]benchVarValue, numValues)
		)
		for j := 0; j < numValues; j++ {
			varValues[j] = benchVarValue{
				name:  fmt.Sprintf("var%d", j),
				value: j,
			}
		}
		res[i] = benchRes{
			inputs: benchInputs{
				varValues: varValues,
			},
			outputs: outputs,
		}
	}
	return Benchmark{
		Name:    "BenchmarkSomeMethod",
		results: res,
	}
}

var parseInfoErr error

func BenchmarkParseInfo(b *testing.B) {
	var (
		dTypes = map[string]func(varName string) string{
			"int": func(varName string) string {
				return fmt.Sprintf("%s=%d", varName, 1)
			},
			"float64": func(varName string) string {
				return fmt.Sprintf("%s=%f", varName, 1.1)
			},
			"bool": func(varName string) string {
				return fmt.Sprintf("%s=%t", varName, true)
			},
			"string": func(varName string) string {
				return fmt.Sprintf("%s=%s", varName, "foo")
			},
		}
		allNumValues = []int{1, 2, 3, 4, 5, 10, 20}
	)

	for _, numValues := range allNumValues {
		b.Run(fmt.Sprintf("num_values=%d", numValues), func(b *testing.B) {
			for dtype, fn := range dTypes {
				b.Run(fmt.Sprintf("dtype=%s", dtype), func(b *testing.B) {
					s := make([]string, numValues+1)
					s[0] = "BenchmarkSomeMethod"
					for i := 1; i <= numValues; i++ {
						s[i] = fn(fmt.Sprintf("var%d", i))
					}
					input := strings.Join(s, "/")

					var err error
					for i := 0; i < b.N; i++ {
						_, _, err = parseInfo(input)
						if err != nil {
							b.Fatalf("unexpected error: %s", err)
						}
					}
					parseInfoErr = err
				})
			}
		})
	}
}
