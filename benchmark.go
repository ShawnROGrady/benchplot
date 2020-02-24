package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/benchmark/parse"
)

type benchVarValue struct {
	name  string
	value interface{}
}

func (b benchVarValue) String() string {
	return fmt.Sprintf("%s=%v", b.name, b.value)
}

type benchVarValues []benchVarValue

func (b benchVarValues) String() string {
	s := make([]string, len(b))
	for i, val := range b {
		s[i] = val.String()
	}
	return strings.Join(s, ",")
}

type benchInputs struct {
	varValues []benchVarValue
	subs      []string
}

type benchOutputs struct {
	N                 int     // number of iterations
	NsPerOp           float64 // nanoseconds per iteration
	AllocedBytesPerOp uint64  // bytes allocated per iteration
	AllocsPerOp       uint64  // allocs per iteration
	MBPerS            float64 // MB processed per second
}

type benchRes struct {
	inputs  benchInputs
	outputs benchOutputs
}

type benchmark struct {
	name    string
	results []benchRes
}

func (b benchmark) groupResults(groupBy []string) map[string][]benchRes {
	groupedResults := map[string][]benchRes{}
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

func parseBenchmarks(r io.Reader) ([]benchmark, error) {
	var (
		scanner    = bufio.NewScanner(r)
		benchmarks = map[string]benchmark{}
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
			bench = benchmark{name: benchName, results: []benchRes{}}
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

	parsedBenchmarks := make([]benchmark, len(benchmarks))
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
