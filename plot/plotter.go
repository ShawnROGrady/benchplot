package plot

// NumericData represents basic numeric data to plot
type NumericData struct {
	X []float64
	Y []float64
}

// Plotter defines the functionality needed to plot a benchmark
type Plotter interface {
	PlotScatter(data map[string]NumericData, title, xLabel, yLabel string) error
}
