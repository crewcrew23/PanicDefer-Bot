package grapth

import (
	"bytes"
	"image/color"
	dbmodel "service-healthz-checker/internal/model/dbModel"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func CreateGrapth(data []*dbmodel.History) ([]byte, error) {
	p := plot.New()
	p.Title.Text = "Response Time History"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "ms"
	p.Legend.Top = true

	timePoints := make(plotter.XYs, len(data))
	for i, v := range data {
		timePoints[i].X = float64(v.CreatedAt.Unix())
		timePoints[i].Y = float64(v.ResponseTimeMs)
	}

	successPoints := make(plotter.XYs, 0)
	errorPoints := make(plotter.XYs, 0)

	for _, v := range data {
		point := plotter.XY{
			X: float64(v.CreatedAt.Unix()),
			Y: float64(v.ResponseTimeMs),
		}

		if v.Status >= 200 && v.Status < 300 {
			successPoints = append(successPoints, point)
		} else {
			errorPoints = append(errorPoints, point)
		}
	}

	timeLine, err := plotter.NewLine(timePoints)
	if err != nil {
		return nil, err
	}
	timeLine.Color = color.RGBA{B: 255, A: 255}

	successScater, err := plotter.NewScatter(successPoints)
	if err != nil {
		return nil, err
	}
	successScater.GlyphStyle.Color = color.RGBA{G: 255, A: 255}

	errorScatter, err := plotter.NewScatter(errorPoints)
	if err != nil {
		return nil, err
	}
	errorScatter.GlyphStyle.Color = color.RGBA{R: 255, A: 255}

	p.Add(timeLine, successScater, errorScatter)
	p.Legend.Add("Response Time", timeLine)
	p.Legend.Add("Success (2xx)", successScater)
	p.Legend.Add("Error", errorScatter)

	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04:05"}

	var buf bytes.Buffer
	writer, err := p.WriterTo(15*vg.Inch, 8*vg.Inch, "png")
	if err != nil {
		return nil, err
	}
	_, err = writer.WriteTo(&buf)
	return buf.Bytes(), err

}
