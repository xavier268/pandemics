package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// ReportWindow is a window that can display a graph of datapoints.
type ReportWindow struct {
	bounds pixel.Rect
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
	rw.bounds = pixel.R(0, 0, 1024, 500)
	rw.imd = imdraw.New(nil)
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

	var v = rw.bounds.Min

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
		v.Y = rw.bounds.Min.Y
	}
}
