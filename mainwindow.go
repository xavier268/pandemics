package main

import (
	"fmt"

	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// MWindow object
type MWindow struct {
	win  *pixelgl.Window
	imd  *imdraw.IMDraw
	pop  *Population
	err  error
	rwin *ReportWindow
}

// NewMWindow constructor
func NewMWindow() *MWindow {

	mw := new(MWindow)

	mw.pop = NewPopulation().Seed(1)

	cfg := pixelgl.WindowConfig{
		Title:  "Animation",
		Bounds: mw.pop.bounds,
		VSync:  true,
	}

	// create the window
	mw.win, mw.err = pixelgl.NewWindow(cfg)
	if mw.err != nil {
		panic(mw.err)
	}

	// Prepare population window
	mw.win.SetSmooth(true)

	mw.imd = imdraw.New(nil)

	mw.rwin = NewReportWindow()

	return mw
}

//Closed wrapper
func (mw *MWindow) Closed() bool {
	return mw.win.Closed()
}

// Update loop
func (mw *MWindow) Update() {
	mw.pop.DeltaTimeCompute()

	// Only move if running  ...
	if mw.pop.running {
		mw.pop.Update()
		mw.pop.CollisionDetect(handlCollision)
		mw.pop.Draw(mw.imd)

		mw.win.Clear(colornames.Skyblue)

		mw.rwin.imd.Draw(mw.win)
		mw.imd.Draw(mw.win)
		mw.fpsDisplay()

	}

	// should we stop running ?
	if mw.pop.running && mw.pop.count[StateTouched] <= 0 {
		mw.pop.running = false
	}

	// Update in anycase, to stay responsive
	mw.win.Update()
}

// fpsDisplay dispaly the FPS in the title bar,
// with the population status
func (mw *MWindow) fpsDisplay() {
	// counting FPS and displaying in title
	select {
	case <-mw.pop.second:
		mw.win.SetTitle(fmt.Sprintf("Time %ds, FPS: %d - Live %d(%.1f%%), Touched %d(%.1f%%), Cured %d(%.1f%%), Dead %d(%.1f%%)",
			mw.pop.elapsed,
			mw.pop.frames,
			mw.pop.count[StateLive], float64(100*mw.pop.count[StateLive])/float64(mw.pop.size),
			mw.pop.count[StateTouched], float64(100*mw.pop.count[StateTouched])/float64(mw.pop.size),
			mw.pop.count[StateCured], float64(100*mw.pop.count[StateCured])/float64(mw.pop.size),
			mw.pop.count[StateDead], float64(100*mw.pop.count[StateDead])/float64(mw.pop.size)))
		mw.pop.frames = 0
		mw.pop.elapsed++

		mw.rwin.Record(stat{
			time:    mw.pop.elapsed,
			live:    mw.pop.count[StateLive],
			dead:    mw.pop.count[StateDead],
			touched: mw.pop.count[StateTouched],
			cured:   mw.pop.count[StateCured],
		})
	default:
		mw.pop.frames++
	}
}
