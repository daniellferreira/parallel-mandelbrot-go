package buffer

import (
	"math"
)

type WorkItem struct {
	InitialX int
	FinalX   int
	InitialY int
	FinalY   int
}

type WorkBuffer chan WorkItem

func NewWorkBuffer(numBlocks, width, height int) WorkBuffer {
	buffer := make(WorkBuffer, numBlocks)

	buffer.workBufferInit(float64(numBlocks), width, height)

	return buffer
}

func (w WorkBuffer) workBufferInit(numBlocks float64, width, height int) {
	var sqrt = int(math.Sqrt(numBlocks))

	for i := sqrt - 1; i >= 0; i-- {
		for j := 0; j < sqrt; j++ {
			w <- WorkItem{
				InitialX: i * (width / sqrt),
				FinalX:   (i + 1) * (width / sqrt),
				InitialY: j * (height / sqrt),
				FinalY:   (j + 1) * (height / sqrt),
			}
		}
	}
}
