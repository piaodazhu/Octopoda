package ostp

type averager struct {
	buffer     []int64
	next       int
	windowSize int
	currentNum int
	average    int64
	currentSum int64
}

func newAverager(windowSize int) *averager {
	return &averager{
		buffer:     make([]int64, windowSize),
		next:       0,
		windowSize: windowSize,
		currentNum: 0,
		average:    0,
	}
}

func (a *averager) Update(num int64) int64 {
	a.currentSum -= a.buffer[a.next]
	a.buffer[a.next] = num
	a.currentSum += a.buffer[a.next]
	if a.currentNum < a.windowSize {
		a.currentNum++
	}
	a.next++
	if a.next == a.windowSize {
		a.next = 0
	}
	// logger.Exceptions.Printf("[DBG] push %d cal avr %d / %d\n", num, a.currentSum, int64(a.currentNum))
	return a.currentSum / int64(a.currentNum)
}
