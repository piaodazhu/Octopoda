package network

import "sync"

var wg,joinwg sync.WaitGroup

func Run() {
	wg.Add(1)
	joinwg.Add(1)
	KeepAlive()
	ReadAndServe()
	wg.Wait()
}
