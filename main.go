package main

import (
	"github.com/faiface/pixel/pixelgl"
)

// The main function needs to delegate to run()
// to ensure thread consistency.
func main() {
	pixelgl.Run(run)
}

// run is where all the code will be fired from.
func run() {

	// prepare report window
	rwin := NewReportWindow()

	// prepare main window
	mwin := NewMWindow()

	for !mwin.Closed() && !rwin.Closed() {
		// Update should be called in the same order
		// as the windows were declared
		rwin.Update()
		mwin.Update()

	}
}
