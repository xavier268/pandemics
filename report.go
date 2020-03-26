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
	rw.bounds = pixel.R(0, 0, 1000, 150)
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
}

// Update the window
func (rw *ReportWindow) Update() {
	rw.win.Clear(colornames.Beige)
	rw.imd.Draw(rw.win)
	// Always keep updateting to stay responsive
	rw.win.Update()
}

// Closed test for close request
func (rw *ReportWindow) Closed() bool {
	return rw.win.Closed()
}
