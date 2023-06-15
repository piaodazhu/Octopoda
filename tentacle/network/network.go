package network

import "sync"

var wg sync.WaitGroup

func Run() {
	wg.Add(1)
	KeepAlive()
	// ListenAndServe()
	ReadAndServe()
	wg.Wait()
}
