package plot

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/ShawnROGrady/benchparse"
	"github.com/ShawnROGrady/benchplot/plot/plotter"
)

// The available plot types.
const (
	ScatterType = "scatter"
	AvgLineType = "avg_line"
)

type plotOptions struct {
	groupBy   []string
	plotTypes []string
}

// Benchmark plots the benchmark.
func Benchmark(b benchparse.Benchmark, p plotter.Plotter, xName, yName string, options ...plotOption) error {
	pltOptions := &plotOptions{
		groupBy:   []string{},
		plotTypes: []string{},
	}
	for _, opt := range options {
		opt.apply(pltOptions)
	}

	grouped := b.GroupResults(pltOptions.groupBy)
	splitGrouped, err := splitGroupedResult(grouped, xName, yName)
	if err != nil {
		return fmt.Errorf("err splitting grouped results: %w", err)
	}

	if len(pltOptions.plotTypes) == 0 {
		plotTypes, err := defaultPlotTypes(splitGrouped)
		if err != nil {
			return err
		}
		pltOptions.plotTypes = plotTypes
	}

	for i, plotType := range pltOptions.plotTypes {
		includeLegend := i == 0
		switch plotType {
		case ScatterType:
			if err := plotScatter(p, b.Name, xName, yName, splitGrouped, includeLegend); err != nil {
				return fmt.Errorf("error creating scatter plot: %w", err)
			}
		case AvgLineType:
			if err := plotAvgLine(p, b.Name, xName, yName, splitGrouped, includeLegend); err != nil {
				return fmt.Errorf("error creating average line plot: %w", err)
			}
		default:
			return fmt.Errorf("unknown plot type: %s", plotType)
		}
	}
	return nil
}

func defaultPlotTypes(splitGrouped map[string][]splitRes) ([]string, error) {
	// just use the first x value
	for _, res := range splitGrouped {
		if len(res) == 0 {
			continue
		}
		xKind := reflect.TypeOf(res[0].x).Kind()
		switch xKind {
		case reflect.Int, reflect.Float64, reflect.Uint64:
			return []string{ScatterType, AvgLineType}, nil
		}
	}
	return []string{}, errors.New("could not determine default plot type")
}

// plotScatter plots the benchmark results as a scatter plot.
func plotScatter(p plotter.Plotter, title, xName, yName string, splitGrouped map[string][]splitRes, includeLegend bool) error {
	var (
		xLabel = xName
		yLabel = yName // TODO: include units
	)

	data, err := splitGroupedPlotData(splitGrouped)
	if err != nil {
		return err
	}
	return p.PlotScatter(data, title, xLabel, yLabel, includeLegend)
}

// plotAvgLine plots the benchmark results as a line where y(x) = avg(f(x)).
func plotAvgLine(p plotter.Plotter, title, xName, yName string, splitGrouped map[string][]splitRes, includeLegend bool) error {
	var (
		xLabel = xName
		yLabel = yName // TODO: include units
	)

	data, err := splitGroupedAvgPlotData(splitGrouped)
	if err != nil {
		return err
	}
	return p.PlotLine(data, title, xLabel, yLabel, includeLegend)
}

func splitGroupedPlotData(splitGrouped map[string][]splitRes) (map[string]plotter.NumericData, error) {
	data := map[string]plotter.NumericData{}
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
		data[groupName] = plotter.NumericData{
			X: xData,
			Y: yData,
		}
	}
	return data, nil
}

func splitGroupedAvgPlotData(splitGrouped map[string][]splitRes) (map[string]plotter.NumericData, error) {
	data := map[string]plotter.NumericData{}
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

		data[groupName] = plotter.NumericData{
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
