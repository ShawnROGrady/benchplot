package benchmark

import (
	"fmt"
	"strings"
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

// the output names
const (
	RunsName        = "runs"
	TimeName        = "time"
	NumAllocsName   = "mem_allocs"
	AllocBytesName  = "mem_used"
	AllocMBytesRate = "mem_by_time"
)

func (b benchOutputs) valByName(name string) (interface{}, error) {
	switch name {
	case RunsName:
		return b.N, nil
	case TimeName:
		return b.NsPerOp, nil
	case NumAllocsName:
		return b.AllocsPerOp, nil
	case AllocBytesName:
		return b.AllocedBytesPerOp, nil
	case AllocMBytesRate:
		return b.MBPerS, nil
	default:
		return nil, fmt.Errorf("no output found with name: '%s'", name)
	}
}

type splitRes struct {
	x interface{}
	y interface{}
}

type benchRes struct {
	inputs  benchInputs
	outputs benchOutputs
}

func (b benchRes) splitTo(xName, yName string) (splitRes, error) {
	splitRes := splitRes{}
	xFound := false
	for _, varValue := range b.inputs.varValues {
		if varValue.name == xName {
			xFound = true
			splitRes.x = varValue.value
			break
		}
	}

	if !xFound {
		return splitRes, fmt.Errorf("no input found with name: '%s'", xName)
	}

	yVal, err := b.outputs.valByName(yName)
	if err != nil {
		return splitRes, err
	}
	splitRes.y = yVal

	return splitRes, nil
}

type groupedResults map[string][]benchRes

func (g groupedResults) splitTo(xName, yName string) (map[string][]splitRes, error) {
	splitGrouped := map[string][]splitRes{}
	for groupName, results := range g {
		splitResults := make([]splitRes, len(results))
		for i, res := range results {
			split, err := res.splitTo(xName, yName)
			if err != nil {
				return nil, err
			}
			splitResults[i] = split
		}
		splitGrouped[groupName] = splitResults
	}
	return splitGrouped, nil
}
