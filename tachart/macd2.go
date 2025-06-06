package tachart

import (
	"strings"
)

type macd2 struct {
	nm   string
	dif  []float64 // 快线
	eda  []float64 // 慢线
	macd []float64 // 柱状图
	ci   int
}

func NewMACD(dif, eda, macds []float64, macd指标值 string) Indicator {
	return &macd2{
		nm:   macd指标值,
		dif:  dif,
		eda:  eda,
		macd: macds,
	}
}

func (c macd2) name() string {
	return c.nm
}

func (c macd2) yAxisLabel() string {
	return strings.Replace(yLabelFormatterFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (c macd2) yAxisMin() string {
	return strings.Replace(minRoundFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (c macd2) yAxisMax() string {
	return strings.Replace(maxRoundFuncTpl, "__DECIMAL_PLACES__", "0", -1)
}

func (c macd2) getNumColors() int {
	return 2
}

func (c *macd2) getTitleOpts(top, left int, colorIndex int) []opts.Title {
	c.ci = colorIndex
	return []opts.Title{
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[c.ci],
				FontSize: chartLabelFontSize,
			},
			Title: c.nm + "-Diff",
			Left:  px(left),
			Top:   px(top),
		},
		opts.Title{
			TitleStyle: &opts.TextStyle{
				Color:    colors[c.ci+1],
				FontSize: chartLabelFontSize,
			},
			Title: c.nm + "-Sig",
			Left:  px(left),
			Top:   px(top + chartLabelFontHeight),
		},
	}
}

func (c macd2) genChart(_, _, _, closes, _ []float64, xAxis interface{}, gridIndex int) charts.Overlaper {

	lineItems := []opts.LineData{}
	for _, v := range c.dif {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}
	macdLine := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Diff", lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[c.ci],
				Opacity: opacityMed,
			}),
		)

	lineItems = []opts.LineData{}
	for _, v := range c.eda {
		lineItems = append(lineItems, opts.LineData{Value: v})
	}
	signalLine := charts.NewLine().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Sig", lineItems,
			charts.WithLineChartOpts(opts.LineChart{
				Symbol:     "none",
				XAxisIndex: gridIndex,
				YAxisIndex: gridIndex,
				ZLevel:     100,
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Color:   colors[c.ci+1],
				Opacity: opacityMed,
			}),
		)

	barItems := []opts.BarData{}
	for _, v := range c.macd {
		style := &opts.ItemStyle{
			Color:   colorUpBar,
			Opacity: opacityHeavy,
		}
		if v < 0 {
			style = &opts.ItemStyle{
				Color:   colorDownBar,
				Opacity: opacityHeavy,
			}
		}
		barItems = append(barItems, opts.BarData{Value: v, ItemStyle: style})
	}
	histBar := charts.NewBar().
		SetXAxis(xAxis).
		AddSeries(c.nm+"-Hist", barItems, charts.WithBarChartOpts(opts.BarChart{
			BarWidth:   "60%",
			XAxisIndex: gridIndex,
			YAxisIndex: gridIndex,
			ZLevel:     100,
		}))

	macdLine.Overlap(signalLine, histBar)

	return macdLine
}
