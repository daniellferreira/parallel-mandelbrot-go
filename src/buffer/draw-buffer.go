package buffer

type Pix struct {
	X  int
	Y  int
	Cr uint8
	Cg uint8
	Cb uint8
}

type DrawBuffer chan Pix

func NewDrawBuffer(pixelTotal int) DrawBuffer {
	return make(chan Pix, pixelTotal)
}
