package main

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

func multiClient(method, url string) {
	var fail int64 = 0
	var reqsent, handshack int64 = 0, 0
	limiter := rate.NewLimiter(rate.Limit(rps), G)
	latency := []int{}
	wg := sync.WaitGroup{}
	wg.Add(C)
	clients := make([]*http.Client, C)
	for i := range clients {
		clients[i] = NewHttpsClient()
	}

	if warmup {
		fmt.Println("warming up ----")
		for i := range clients {
			clients[i].Get(host + "/ping") // warm up
		}
		fmt.Println("warming up DONE")
	}

	start := time.Now()
	var reqRate float64
	onceDone := sync.Once{}
	for c := 0; c < C; c++ {
		go func(idx int) {
			for i := 0; i < N/C; i++ {
				req, _ := http.NewRequest(method, url, nil)
				reqStart := time.Now()
				if time.Since(start) > time.Second*3 {
					onceDone.Do(func() {
						reqRate = float64(reqsent) / time.Since(start).Seconds()
					})
					break
				}
				atomic.AddInt64(&reqsent, 1)
				if i == 0 {
					atomic.AddInt64(&handshack, 1)
				}

				if rps > 0 {
					limiter.Wait(context.Background())
				}
				rsp, err := clients[idx].Do(req)
				if err != nil {
					atomic.AddInt64(&fail, 1)
					continue
				}
				reqDuration := time.Since(reqStart)
				rsp.Body.Close()

				latency = append(latency, int(reqDuration.Microseconds()))
			}
			wg.Done()
		}(c)
	}
	wg.Wait()

	duration := time.Since(start)
	fmt.Printf("%d client send %d requests with %d handshack with %f rps. avr= %d us. succ rate=%f\n", C, reqsent, handshack, reqRate, duration.Microseconds()/int64(reqsent), float64(reqsent-fail)/float64(reqsent))

	sort.Ints(latency)
	fmt.Printf("10%% requests finish in %d us.\n", latency[len(latency)/100*10])
	fmt.Printf("30%% requests finish in %d us.\n", latency[len(latency)/100*30])
	fmt.Printf("50%% requests finish in %d us.\n", latency[len(latency)/100*50])
	fmt.Printf("80%% requests finish in %d us.\n", latency[len(latency)/100*80])
	fmt.Printf("100%% requests finish in %d us.\n", latency[len(latency)-1])
}

func queryBench() {
	url := host + "/query?name=test"
	multiClient("GET", url)
}

func writeBench() {
	url := host + "/register?name=test&ip=1.1.1.1&port=11&type=brain&description=abc&ttl=1000000"
	multiClient("POST", url)
}
