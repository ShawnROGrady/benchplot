package plot

type plotOption interface {
	apply(*plotOptions)
}

// WithPlotTypes is an option to specify the plots to create.
type WithPlotTypes []string

func (w WithPlotTypes) apply(p *plotOptions) {
	p.plotTypes = []string(w)
}

// WithGroupBy is an option to specify how plot data should be grouped.
type WithGroupBy []string

func (w WithGroupBy) apply(p *plotOptions) {
	p.groupBy = []string(w)
}

// WithFilterBy is an option to specify any expressions to filter
// the data to be plotted.
type WithFilterBy []string

func (w WithFilterBy) apply(p *plotOptions) {
	p.filterExprs = []string(w)
}
