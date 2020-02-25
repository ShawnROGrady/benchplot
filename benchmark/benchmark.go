package benchmark

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/benchmark/parse"
)

// Benchmark represents a single benchmark and it's results
type Benchmark struct {
	Name    string
	results []benchRes
}

func (b Benchmark) groupResults(groupBy []string) groupedResults {
	groupedResults := map[string][]benchRes{}
	if len(groupBy) == 0 {
		res := make([]benchRes, len(b.results))
		copy(res, b.results)
		groupedResults[""] = res
		return groupedResults
	}
	for _, result := range b.results {
		groupVals := benchVarValues{}
		for _, varValue := range result.inputs.varValues {
			for _, groupName := range groupBy {
				if varValue.name == groupName {
					groupVals = append(groupVals, varValue)
				}
			}
		}
		if len(groupVals) == len(groupBy) {
			k := groupVals.String()
			if existingResults, ok := groupedResults[k]; ok {
				groupedResults[k] = append(existingResults, result)
			} else {
				groupedResults[k] = []benchRes{result}
			}
		}
	}
	return groupedResults
}

// ParseBenchmarks extracts a list of Benchmarks from testing.B output
func ParseBenchmarks(r io.Reader) ([]Benchmark, error) {
	var (
		scanner    = bufio.NewScanner(r)
		benchmarks = map[string]Benchmark{}
	)
	for scanner.Scan() {
		parsed, err := parse.ParseLine(scanner.Text())
		if err != nil {
			// TODO: this is what ParseSet does but feels awkward - https://github.com/golang/tools/blob/master/benchmark/parse/parse.go#L114
			continue
		}

		benchName, inputs, err := parseInfo(parsed.Name)
		if err != nil {
			return nil, err
		}
		bench, ok := benchmarks[benchName]
		if !ok {
			bench = Benchmark{Name: benchName, results: []benchRes{}}
		}

		outputs := benchOutputs{
			N:                 parsed.N,
			NsPerOp:           parsed.NsPerOp,
			AllocedBytesPerOp: parsed.AllocedBytesPerOp,
			AllocsPerOp:       parsed.AllocsPerOp,
			MBPerS:            parsed.MBPerS,
		}

		bench.results = append(bench.results, benchRes{
			inputs:  inputs,
			outputs: outputs,
		})

		benchmarks[benchName] = bench
	}

	parsedBenchmarks := make([]Benchmark, len(benchmarks))
	i := 0
	for _, v := range benchmarks {
		parsedBenchmarks[i] = v
		i++
	}

	return parsedBenchmarks, nil
}

// used to trim unnecessary trailing chars from benchname
var benchInfoExpr = regexp.MustCompile(`^(Benchmark.+?)(?:\-[0-9])?$`)

func parseInfo(s string) (string, benchInputs, error) {
	submatches := benchInfoExpr.FindStringSubmatch(s)
	if len(submatches) < 1 {
		return "", benchInputs{}, fmt.Errorf("info string '%s' didn't match regex", s)
	}
	info := submatches[1]
	var (
		name      string
		varValues = []benchVarValue{}
		subs      = []string{}
		bySub     = strings.Split(info, "/")
	)

	for i, sub := range bySub {
		if i == 0 {
			name = sub
			continue
		}

		split := strings.Split(sub, "=")
		if len(split) == 2 {
			varValues = append(varValues, benchVarValue{
				name:  split[0],
				value: value(split[1]),
			})
		} else {
			subs = append(subs, sub)
		}
	}

	return name, benchInputs{varValues: varValues, subs: subs}, nil
}

func value(s string) interface{} {
	convs := []func(str string) (interface{}, error){
		func(str string) (interface{}, error) {
			return strconv.Atoi(str)
		},
		func(str string) (interface{}, error) {
			return strconv.ParseFloat(str, 64)
		},
		func(str string) (interface{}, error) {
			return strconv.ParseBool(str)
		},
	}

	for _, conv := range convs {
		if res, err := conv(s); err == nil {
			return res
		}
	}

	return s
}
