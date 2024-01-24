package service

import (
	"bytes"
	"os"
	"time"

	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func ExportChart(xValues []time.Time, yValues []float64, title string, filename string) {
	// lineChartStyle := chart.Style{
	// 	Padding: chart.Box{
	// 		Top:    30,
	// 		Left:   30,
	// 		Right:  30,
	// 		Bottom: 30,
	// 	},
	// }

	graph := chart.Chart{
		Title: title,
		// Background: lineChartStyle,
		Width:  1400,
		Height: 600,
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xValues,
				YValues: yValues,
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   drawing.ColorFromHex("9ADFEA"),
				},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	f, _ := os.Create(filename + ".png")
	_, _ = f.Write(buffer.Bytes())
}
