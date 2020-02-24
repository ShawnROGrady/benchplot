package plot

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// GoNumPlotter is a Plotter implementation using gonum/plot/plotter - https://godoc.org/gonum.org/v1/plot/plotter
type GoNumPlotter struct {
	p *plot.Plot
}

// NewGoNumPlotter constructs a new plotter
func NewGoNumPlotter() (*GoNumPlotter, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}

	return &GoNumPlotter{
		p: p,
	}, nil
}

// PlotScatter creates and saves a scatter plot of the specified data
func (g *GoNumPlotter) PlotScatter(data map[string]NumericData, title, xLabel, yLabel string) error {
	g.p.Title.Text = title
	g.p.X.Label.Text = xLabel
	g.p.Y.Label.Text = yLabel
	vs := make([]interface{}, len(data)*2)
	i := 0
	for groupName, groupData := range data {
		vs[i] = groupName
		vs[i+1] = numericDataXYs(groupData)
		i = i + 2
	}
	return plotutil.AddScatters(g.p, vs...)
}

// Save saves the plot to a file
func (g *GoNumPlotter) Save(dstWidth, dstHeight float64, dstName string) error {
	return g.p.Save(vg.Length(dstWidth), vg.Length(dstHeight), dstName)
}

func numericDataXYs(data NumericData) plotter.XYs {
	xys := make(plotter.XYs, len(data.X))
	for i := 0; i < len(data.X); i++ {
		xys[i].X = data.X[i]
		xys[i].Y = data.Y[i]
	}
	return xys
}
