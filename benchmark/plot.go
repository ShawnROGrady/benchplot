package benchmark

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/ShawnROGrady/benchplot/plot"
)

// PlotScatter plots the benchmark results as a scatter plot
func (b Benchmark) PlotScatter(p plot.Plotter, groupBy []string, xName, yName string) error {
	var (
		title  = b.Name
		xLabel = xName
		yLabel = yName // TODO: include units
	)

	grouped := b.groupResults(groupBy)
	splitGrouped, err := grouped.splitTo(xName, yName)
	if err != nil {
		return fmt.Errorf("err splitting grouped results: %w", err)
	}

	data, err := splitGroupedPlotData(splitGrouped)
	if err != nil {
		return err
	}
	return p.PlotScatter(data, title, xLabel, yLabel)
}

// PlotAvgLine plots the benchmark results as a line where y(x) = avg(f(x))
func (b Benchmark) PlotAvgLine(p plot.Plotter, groupBy []string, xName, yName string) error {
	var (
		title  = b.Name
		xLabel = xName
		yLabel = yName // TODO: include units
	)

	grouped := b.groupResults(groupBy)
	splitGrouped, err := grouped.splitTo(xName, yName)
	if err != nil {
		return fmt.Errorf("err splitting grouped results: %w", err)
	}

	data, err := splitGroupedAvgPlotData(splitGrouped)
	if err != nil {
		return err
	}
	return p.PlotLine(data, title, xLabel, yLabel)
}

func splitGroupedPlotData(splitGrouped map[string][]splitRes) (map[string]plot.NumericData, error) {
	data := map[string]plot.NumericData{}
	for groupName, splitResults := range splitGrouped {
		var (
			xData = []float64{}
			yData = []float64{}
		)

		for _, res := range splitResults {
			xF, err := getFloat(res.x)
			if err != nil {
				return nil, fmt.Errorf("cannot create scatter plot from x data: %w", err)
			}
			xData = append(xData, xF)

			yF, err := getFloat(res.y)
			if err != nil {
				return nil, fmt.Errorf("cannot create scatter plot from y data: %w", err)
			}
			yData = append(yData, yF)
		}
		data[groupName] = plot.NumericData{
			X: xData,
			Y: yData,
		}
	}
	return data, nil
}

func splitGroupedAvgPlotData(splitGrouped map[string][]splitRes) (map[string]plot.NumericData, error) {
	data := map[string]plot.NumericData{}
	for groupName, splitResults := range splitGrouped {
		// track y values corresponding to each x
		vals := map[float64][]float64{}

		for _, res := range splitResults {
			xF, err := getFloat(res.x)
			if err != nil {
				return nil, fmt.Errorf("cannot create scatter plot from x data: %w", err)
			}

			yF, err := getFloat(res.y)
			if err != nil {
				return nil, fmt.Errorf("cannot create scatter plot from y data: %w", err)
			}

			if xVals, ok := vals[xF]; ok {
				vals[xF] = append(xVals, yF)
			} else {
				vals[xF] = []float64{yF}
			}
		}

		var (
			xData = make([]float64, len(vals))
			yData = make([]float64, len(vals))
		)

		i := 0
		for x := range vals {
			xData[i] = x
			i++
		}
		// keep data sorted wrt x
		sort.Float64s(xData)

		for i, xVal := range xData {
			yVals := vals[xVal]
			var totY float64 = 0
			for _, yVal := range yVals {
				totY += yVal
			}
			yData[i] = totY / float64(len(yVals))
		}

		data[groupName] = plot.NumericData{
			X: xData,
			Y: yData,
		}
	}
	return data, nil
}

func getFloat(data interface{}) (float64, error) {
	val := reflect.ValueOf(data)
	switch val.Type().Kind() {
	case reflect.Int:
		return float64(val.Int()), nil
	case reflect.Float64:
		return val.Float(), nil
	case reflect.Uint64:
		return float64(val.Uint()), nil
	default:
		return 0, fmt.Errorf("unexpected kind: '%s'", val.Type().Kind())
	}
}
