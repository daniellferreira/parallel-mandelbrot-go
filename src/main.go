package main

import (
	"github.com/daniellferreira/parallel-mandelbrot-go/src/window"
	"github.com/faiface/pixel/pixelgl"
)

const (
	defaultMaxIter    = 1000
	defaultSamples    = 200
	defaultNumBlocks  = 64
	defaultNumThreads = 16
)

func main() {
	w := window.NewWindow(defaultMaxIter, defaultSamples, defaultNumBlocks, defaultNumThreads)

	pixelgl.Run(w.Open)
}
