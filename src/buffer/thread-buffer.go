package buffer

type ThreadBuffer chan bool

func NewThreadBuffer(numThreads int) ThreadBuffer {
	buffer := make(ThreadBuffer, numThreads)

	buffer.threadBufferInit(numThreads)

	return buffer
}

func (t ThreadBuffer) threadBufferInit(numThreads int) {
	for i := 1; i <= numThreads; i++ {
		t <- true
	}
}
