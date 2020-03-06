package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ShawnROGrady/benchparse"
	"github.com/ShawnROGrady/benchplot/gonum"
	"github.com/ShawnROGrady/benchplot/plot"
)

func main() {
	var (
		benchName = flag.String("bench", "", "The name of the benchmark to plot")
		xName     = flag.String("x", "", "The name of the x-axis variable (an input to the benchmark)")
		yName     = flag.String("y", plot.TimeName, "The name of the y-axis variable")
		dstName   = flag.String("o", "", "The output file name with extension (if empty will be set to ${bench}.png)")
		dstWidth  = flag.Float64("width", 500, "The width of the output figure")
		dstHeight = flag.Float64("height", 500, "The height of the output figure")
		help      = flag.Bool("h", false, "Show this help message and exit")
		groupBy   = &stringSliceFlag{}
		plotTypes = &stringSliceFlag{}
		filterBy  = &stringSliceFlag{}
		resFile   *os.File
	)
	flag.Var(groupBy, "group-by", "The variables to group results by (an input to the benchmark)")
	flag.Var(plotTypes, "plots", fmt.Sprintf("The plots to generate (options = %q). If empty will default to %q for numeric data", []string{plot.ScatterType, plot.AvgLineType}, []string{plot.ScatterType, plot.AvgLineType}))
	flag.Var(
		filterBy, "filter-by",
		fmt.Sprintf(
			"Expressions to filter results by. Form: 'var_name==var_value'. Available comparison operations: %q",
			[]benchparse.Comparison{benchparse.Eq, benchparse.Ne, benchparse.Lt, benchparse.Gt, benchparse.Le, benchparse.Ge},
		),
	)

	flag.Parse()
	if help != nil && *help {
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.Usage()
		return
	}
	if xName == nil || *xName == "" {
		log.Fatal("x-axis variable is required")
	}
	if yName == nil || *yName == "" {
		log.Fatal("y-axis variable is required")
	}
	if benchName == nil || *benchName == "" {
		log.Fatal("benchmark name is required")
	}
	if dstName == nil || *dstName == "" {
		*dstName = fmt.Sprintf("%s.png", *benchName)
	}

	args := flag.Args()
	if len(args) == 0 || args[0] == "-" {
		resFile = os.Stdin
	} else {
		var err error
		resFile, err = os.Open(args[0])
		if err != nil {
			log.Fatalf("error opening '%s': %s", args[0], err)
		}
	}

	benches, err := benchparse.ParseBenchmarks(resFile)
	if err != nil {
		log.Fatalf("error parsing input: %s", err)
	}

	bench, err := findBenchmark(benches, *benchName)
	if err != nil {
		log.Fatal(err)
	}

	p := &gonum.Plotter{}
	if err := plot.Benchmark(bench, p, *xName, *yName, plot.WithGroupBy(*groupBy), plot.WithFilterBy(*filterBy), plot.WithPlotTypes(*plotTypes)); err != nil {
		log.Fatalf("error plotting: %s", err)
	}

	if err := p.Save(*dstWidth, *dstHeight, *dstName); err != nil {
		log.Fatalf("error saving figure: %s", err)
	}
}

func findBenchmark(benches []benchparse.Benchmark, benchName string) (benchparse.Benchmark, error) {
	for i := range benches {
		if benches[i].Name == benchName {
			return benches[i], nil
		}
	}
	return benchparse.Benchmark{}, fmt.Errorf("no benches found with name: %s", benchName)
}
