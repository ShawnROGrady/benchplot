package benchmark

import (
	"fmt"
	"reflect"

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

	data := map[string]plot.NumericData{}
	for groupName, splitResults := range splitGrouped {
		var (
			xData = []float64{}
			yData = []float64{}
		)

		for _, res := range splitResults {
			var (
				xVal = reflect.ValueOf(res.x)
				yVal = reflect.ValueOf(res.y)
			)

			switch xVal.Type().Kind() {
			case reflect.Int:
				xData = append(xData, float64(xVal.Int()))
			case reflect.Float64:
				xData = append(xData, xVal.Float())
			default:
				return fmt.Errorf("cannot create scatter plot from x data of kind: '%s'", xVal.Type().Kind())
			}

			switch yVal.Type().Kind() {
			case reflect.Int:
				yData = append(yData, float64(yVal.Int()))
			case reflect.Float64:
				yData = append(yData, yVal.Float())
			case reflect.Uint64:
				yData = append(yData, float64(yVal.Uint()))
			default:
				return fmt.Errorf("cannot create scatter plot from y data of kind: '%s'", yVal.Type().Kind())
			}
		}
		data[groupName] = plot.NumericData{
			X: xData,
			Y: yData,
		}
	}

	return p.PlotScatter(data, title, xLabel, yLabel)
}
