package window

import (
	"fmt"
	"github.com/daniellferreira/parallel-mandelbrot-go/src/worker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	"image/color"
	"log"
	"time"
)

const (
	imgWidth   = 1024
	imgHeight  = 1024
	pixelTotal = imgWidth * imgHeight

	showProgress = true
	closeOnEnd   = false
)

type Window interface {
	Open()
}

type window struct {
	img             *image.RGBA
	pixelCount      int
	worker          worker.Worker
	finishedProcess bool
}

func NewWindow(maxIter, samples, numBlocks, numThreads int) Window {
	return &window{
		img:    image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight)),
		worker: worker.NewWorker(maxIter, samples, numBlocks, numThreads, imgWidth, imgHeight),
	}
}

func (w *window) Open() {
	log.Println("Initial processing...")

	cfg := pixelgl.WindowConfig{
		Title:  "Parallel Mandelbrot - Goroutines",
		Bounds: pixel.R(0, 0, imgWidth, imgHeight),
		VSync:  true,
		// Invisible: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	log.Println("Rendering...")

	start := time.Now()

	go w.worker.Run()
	go w.drawThread(win)

	for !win.Closed() {
		pic := pixel.PictureDataFromImage(w.img)
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()

		if w.finishedProcess {
			if closeOnEnd {
				break
			}
			continue
		}

		if showProgress {
			fmt.Printf("\rDrawed pixels: %d/%d (%d%%)", w.pixelCount, pixelTotal, int(100*(float64(w.pixelCount)/float64(pixelTotal))))
		}

		if w.pixelCount == pixelTotal {
			end := time.Now()
			fmt.Println("\nFinalizou com tempo = ", end.Sub(start))
			w.finishedProcess = true
		}
	}
}

func (w *window) drawThread(win *pixelgl.Window) {
	for i := range w.worker.GetDrawBuffer() {
		w.img.SetRGBA(i.X, i.Y, color.RGBA{R: i.Cr, G: i.Cg, B: i.Cb, A: 255})
		w.pixelCount++
	}
}
