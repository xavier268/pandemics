package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// NewConfig creates a configuration object.
func NewConfig() pixelgl.WindowConfig {
	return pixelgl.WindowConfig{
		Title:  "Animation",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
}
