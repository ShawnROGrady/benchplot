package plot

import (
	"fmt"

	"github.com/ShawnROGrady/benchparse"
)

// The benchmark output names. These are the available values
// for the y-axis.
const (
	RunsName        = "runs"
	TimeName        = "time"
	NumAllocsName   = "mem_allocs"
	AllocBytesName  = "mem_used"
	AllocMBytesRate = "mem_by_time"
)

func benchOutputValByName(b benchparse.BenchOutputs, name string) (interface{}, error) {
	switch name {
	case RunsName:
		return b.GetIterations(), nil
	case TimeName:
		return b.GetNsPerOp()
	case NumAllocsName:
		return b.GetAllocsPerOp()
	case AllocBytesName:
		return b.GetAllocedBytesPerOp()
	case AllocMBytesRate:
		return b.GetMBPerS()
	default:
		return nil, fmt.Errorf("no output found with name: '%s'", name)
	}
}

type splitRes struct {
	x interface{}
	y interface{}
}

func splitBenchRes(b benchparse.BenchRes, xName, yName string) (splitRes, error) {
	splitRes := splitRes{}
	xFound := false
	for _, varValue := range b.Inputs.VarValues {
		if varValue.Name == xName {
			xFound = true
			splitRes.x = varValue.Value
			break
		}
	}

	if !xFound {
		return splitRes, fmt.Errorf("no input found with name: '%s'", xName)
	}

	yVal, err := benchOutputValByName(b.Outputs, yName)
	if err != nil {
		return splitRes, fmt.Errorf("error getting y value: %w", err)
	}
	splitRes.y = yVal

	return splitRes, nil
}

func splitGroupedResult(g benchparse.GroupedResults, xName, yName string) (map[string][]splitRes, error) {
	splitGrouped := map[string][]splitRes{}
	for groupName, results := range g {
		splitResults := make([]splitRes, len(results))
		for i, res := range results {
			split, err := splitBenchRes(res, xName, yName)
			if err != nil {
				return nil, err
			}
			splitResults[i] = split
		}
		splitGrouped[groupName] = splitResults
	}
	return splitGrouped, nil
}
