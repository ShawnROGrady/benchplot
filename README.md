# benchplot
This is a tool to plot the results of go benchmarks. It assumes that each sub-benchmark is named as `var_name=var_value`.
The input is the output of a go benchmark.

## Usage
`benchplot -bench ${bench} -x ${x_var} ${FILE}`
Where `${FILE}` is the path to a file containing the output of a go benchmark (if empty or `"-"` stdin is used), `${bench}` is the name of the benchmark to plot, and `${x_var}` is the name of the variable to use for the x-axis of the plot.

Full flag set:
```
 -bench string
        the name of the benchmark to plot
  -group-by value
        the variables to group results by (an input to the benchmark)
  -h float
        the height of the output figure (default 500)
  -o string
        the output file name with extension (if empty will be set to ${bench}.png)
  -plots value
        the plots to generate (options = ["scatter" "avg_line"]). If empty will default to ["scatter" "avg_line"] for numeric data
  -w float
        the width of the output figure (default 500)
  -x string
        the name of the x-axis variable (an input to the benchmark)
  -y string
        the name of the y-axis variable (default "time")
```

## Examples
Plotting the results of `BenchmarkGroupResults` (in `benchmark/benchmark_test.go`):
```
go test ./benchmark -run ! -bench BenchmarkGroupResults -count 3 -race | benchplot -bench BenchmarkGroupResults -group-by group_by_count -x num_results
```
![benchgroupres](https://github.com/ShawnROGrady/benchplot/blob/master/assets/BenchmarkGroupResults.png)

## Next Steps
Right now the main focus is bringing the feature set to parity with my [initial implementation of this tool](https://github.com/ShawnROGrady/gobenchplot) which used Python and matplotlib. This includes supporting bar charts (which would be the default if the provided `${x_var}` had non-numeric data) and filtering of data.