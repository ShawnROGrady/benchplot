package plot

import (
	"sort"

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

// PlotScatter creates a scatter plot of the specified data
func (g *GoNumPlotter) PlotScatter(data map[string]NumericData, title, xLabel, yLabel string, includeLegend bool) error {
	g.p.Title.Text = title
	g.p.X.Label.Text = xLabel
	g.p.Y.Label.Text = yLabel

	// use sorted keys for consistent iteration order
	groupNames := make([]string, len(data))
	j := 0
	for k := range data {
		groupNames[j] = k
		j++
	}
	sort.Strings(groupNames)

	var vs []interface{}
	if includeLegend {
		vs = make([]interface{}, len(data)*2)
		for i, groupName := range groupNames {
			groupData := data[groupName]
			vs[2*i] = groupName
			vs[2*i+1] = numericDataXYs(groupData)
		}
	} else {
		vs = make([]interface{}, len(data))
		for i, groupName := range groupNames {
			groupData := data[groupName]
			vs[i] = numericDataXYs(groupData)
		}
	}
	return plotutil.AddScatters(g.p, vs...)
}

// PlotLine creates a line plot of the specified data
func (g *GoNumPlotter) PlotLine(data map[string]NumericData, title, xLabel, yLabel string, includeLegend bool) error {
	g.p.Title.Text = title
	g.p.X.Label.Text = xLabel
	g.p.Y.Label.Text = yLabel

	// use sorted keys for consistent iteration order
	groupNames := make([]string, len(data))
	j := 0
	for k := range data {
		groupNames[j] = k
		j++
	}
	sort.Strings(groupNames)

	var vs []interface{}
	if includeLegend {
		vs = make([]interface{}, len(data)*2)
		for i, groupName := range groupNames {
			groupData := data[groupName]
			vs[2*i] = groupName
			vs[2*i+1] = numericDataXYs(groupData)
		}
	} else {
		vs = make([]interface{}, len(data))
		for i, groupName := range groupNames {
			groupData := data[groupName]
			vs[i] = numericDataXYs(groupData)
		}
	}
	return plotutil.AddLines(g.p, vs...)
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
