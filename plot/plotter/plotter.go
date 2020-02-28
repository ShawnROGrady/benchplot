package plotter

// NumericData represents basic numeric data to plot.
type NumericData struct {
	X []float64
	Y []float64
}

// Plotter defines the functionality needed to plot a benchmark.
type Plotter interface {
	PlotScatter(data map[string]NumericData, title, xLabel, yLabel string, includeLegend bool) error
	PlotLine(data map[string]NumericData, title, xLabel, yLabel string, includeLegend bool) error
}
