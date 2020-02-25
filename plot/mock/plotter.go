package mock

import "github.com/ShawnROGrady/benchplot/plot"

// Plotter is a mock implementation of Plotter
type Plotter struct {
	PlotScatterFn func(data map[string]plot.NumericData, title string, xLabel string, yLabel string, includeLegend bool) error
	PlotLineFn    func(data map[string]plot.NumericData, title string, xLabel string, yLabel string, includeLegend bool) error
}

// PlotScatter returns _m.PlotScatterFn
func (_m *Plotter) PlotScatter(data map[string]plot.NumericData, title string, xLabel string, yLabel string, includeLegend bool) error {
	return _m.PlotScatterFn(data, title, xLabel, yLabel, includeLegend)
}

// PlotLine returns _m.PlotLineFn
func (_m *Plotter) PlotLine(data map[string]plot.NumericData, title string, xLabel string, yLabel string, includeLegend bool) error {
	return _m.PlotLineFn(data, title, xLabel, yLabel, includeLegend)
}
