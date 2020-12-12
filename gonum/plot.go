// Package gonum contains a Plotter implementation using
// the gonum/plot package(https://godoc.org/gonum.org/v1/plot).
package gonum

import (
	"fmt"
	"sort"

	"github.com/ShawnROGrady/benchplot/plot/plotter"
	gonumplot "gonum.org/v1/plot"
	gonumplotter "gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Plotter wraps a gonum/plot.Plot to implement Plotter.
type Plotter struct {
	TopLegend  bool
	LeftLegend bool
	p          *gonumplot.Plot
}

func (g *Plotter) init() error {
	// NOTE: might be worth switching to sync.Once if it's safe to plot multiple things concurrently
	if g.p != nil {
		return nil
	}

	var err error
	g.p, err = gonumplot.New()

	if err != nil {
		return fmt.Errorf("error initializing plotter: %w", err)
	}

	g.p.Legend.Top = g.TopLegend
	g.p.Legend.Left = g.LeftLegend

	return nil
}

// PlotScatter creates a scatter plot of the specified data.
func (g *Plotter) PlotScatter(data map[string]plotter.NumericData, title, xLabel, yLabel string, includeLegend bool) error {
	if err := g.init(); err != nil {
		return err
	}
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

// PlotLine creates a line plot of the specified data.
func (g *Plotter) PlotLine(data map[string]plotter.NumericData, title, xLabel, yLabel string, includeLegend bool) error {
	if err := g.init(); err != nil {
		return err
	}
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
func (g *Plotter) Save(dstWidth, dstHeight float64, dstName string) error {
	if err := g.init(); err != nil {
		return err
	}
	return g.p.Save(vg.Length(dstWidth), vg.Length(dstHeight), dstName)
}

func numericDataXYs(data plotter.NumericData) gonumplotter.XYs {
	xys := make(gonumplotter.XYs, len(data.X))
	for i := 0; i < len(data.X); i++ {
		xys[i].X = data.X[i]
		xys[i].Y = data.Y[i]
	}
	return xys
}
