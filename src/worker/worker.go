package worker

import (
	"github.com/daniellferreira/parallel-mandelbrot-go/src/buffer"
	"image/color"
)

const (
	posX   = -2
	posY   = -1.2
	height = 2.5
)

type Worker interface {
	Run()
	GetDrawBuffer() buffer.DrawBuffer
}

type worker struct {
	width  int
	height int
	ratio  float64

	maxIter int
	samples int

	numBlocks  int
	numThreads int

	workBuffer   buffer.WorkBuffer
	threadBuffer buffer.ThreadBuffer
	drawBuffer   buffer.DrawBuffer
}

func NewWorker(maxIter, samples, numBlocks, numThreads, width, height int) Worker {
	return &worker{
		maxIter:      maxIter,
		samples:      samples,
		numBlocks:    numBlocks,
		numThreads:   numThreads,
		width:        width,
		height:       height,
		ratio:        float64(width) / float64(height),
		workBuffer:   buffer.NewWorkBuffer(numBlocks, width, height),
		threadBuffer: buffer.NewThreadBuffer(numThreads),
		drawBuffer:   buffer.NewDrawBuffer(width * height),
	}
}

func (w *worker) Run() {
	w.runWork()
}

func (w *worker) GetDrawBuffer() buffer.DrawBuffer {
	return w.drawBuffer
}

func (w *worker) runWork() {
	for range w.threadBuffer {
		workItem := <-w.workBuffer

		go w.workerThread(workItem)
	}
}

func (w *worker) workerThread(workItem buffer.WorkItem) {
	for x := workItem.InitialX; x < workItem.FinalX; x++ {
		for y := workItem.InitialY; y < workItem.FinalY; y++ {
			var colorR, colorG, colorB int
			for k := 0; k < w.samples; k++ {
				a := height*w.ratio*((float64(x)+RandFloat64())/float64(w.width)) + posX
				b := height*((float64(y)+RandFloat64())/float64(w.height)) + posY
				c := pixelColor(w.mandelbrotIteraction(a, b))
				colorR += int(c.R)
				colorG += int(c.G)
				colorB += int(c.B)
			}
			var cr, cg, cb uint8
			cr = uint8(float64(colorR) / float64(w.samples))
			cg = uint8(float64(colorG) / float64(w.samples))
			cb = uint8(float64(colorB) / float64(w.samples))

			w.drawBuffer <- buffer.Pix{
				X: x, Y: y, Cr: cr, Cg: cg, Cb: cb,
			}
		}
	}
	w.threadBuffer <- true
}

func (w *worker) mandelbrotIteraction(a, b float64) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < w.maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		// xn+1 = x^2 - y^2 + a
		x = xx - yy + a
		// yn+1 = 2xy + b
		y = 2*xy + b
	}

	return xx + yy, w.maxIter
}

func pixelColor(r float64, iter int) color.RGBA {
	insideSet := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// validar se está dentro do conjunto
	// https://pt.wikipedia.org/wiki/Conjunto_de_Mandelbrot
	if r > 4 {
		return hslToRGB(float64(0.70)-float64(iter)/3500*r, 1, 0.5)
		// return hslToRGB(float64(iter)/100*r, 1, 0.5)
	}

	return insideSet
}
