package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// ReportWindow is a window that can display a graph of datapoints.
type ReportWindow struct {
	bounds pixel.Rect
	win    *pixelgl.Window
	imd    *imdraw.IMDraw
	err    error
	stats  []stat
}

// stat record at a given time
type stat struct {
	time                       int
	live, dead, touched, cured int
}

// NewReportWindow constructor.
func NewReportWindow() *ReportWindow {
	rw := new(ReportWindow)
	rw.bounds = pixel.R(0, 0, 1000, 600)
	rw.win, rw.err = pixelgl.NewWindow(pixelgl.WindowConfig{Bounds: rw.bounds})
	if rw.err != nil {
		panic(rw.err)
	}
	rw.imd = imdraw.New(nil)
	rw.win.Clear(colornames.Lightgray)
	rw.stats = make([]stat, 0, 100)
	return rw
}

// Record a new data point
func (rw *ReportWindow) Record(s stat) {
	rw.stats = append(rw.stats, s)
	rw.prepareReport()
}

// Prepare the report by drawing on the imd object.
func (rw *ReportWindow) prepareReport() {

	var v pixel.Vec
	for _, stat := range rw.stats {

		rw.imd.Color = StateColor(StateLive)
		rw.imd.Push(v)
		v.Y += float64(stat.live)
		rw.imd.Push(v)
		rw.imd.Line(5)

		rw.imd.Color = StateColor(StateTouched)
		rw.imd.Push(v)
		v.Y += float64(stat.touched)
		rw.imd.Push(v)
		rw.imd.Line(5)

		rw.imd.Color = StateColor(StateCured)
		rw.imd.Push(v)
		v.Y += float64(stat.cured)
		rw.imd.Push(v)
		rw.imd.Line(5)

		rw.imd.Color = StateColor(StateDead)
		rw.imd.Push(v)
		v.Y += float64(stat.dead)
		rw.imd.Push(v)
		rw.imd.Line(5)

		v.X += 5
		v.Y = 0
	}
}

// Update the window
func (rw *ReportWindow) Update() {
	rw.win.Clear(colornames.Beige)
	rw.imd.Draw(rw.win)
	// Always keep updating to stay responsive
	rw.win.Update()
}

// Closed test for close request
func (rw *ReportWindow) Closed() bool {
	return rw.win.Closed()
}
