package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
)

var command string
var N, G, C int
var host string
var warmup bool
var rps int

func main() {
	flag.StringVar(&ip, "h", "127.0.0.1", "ns host")
	flag.IntVar(&port, "p", 3455, "listening port")
	flag.StringVar(&caCertFile, "ca", "/home/netlab1038/wuyou/netlab1038/ca.pem", "ca certificate")
	flag.StringVar(&cliCertFile, "crt", "/home/netlab1038/wuyou/netlab1038/client.pem", "client certificate")
	flag.StringVar(&cliKeyFile, "key", "/home/netlab1038/wuyou/netlab1038/client.key", "client private key")
	flag.StringVar(&command, "cmd", "basic", "app cmd: basic, queryBench, writeBench")
	flag.IntVar(&N, "N", 1000000, "number of bench requests")
	flag.IntVar(&G, "G", runtime.NumCPU(), "number of bench goroutines")
	flag.IntVar(&C, "C", 10000, "number of bench clients")
	flag.IntVar(&rps, "r", 0, "expected request rate per sec. 0 means batch.")
	flag.BoolVar(&warmup, "w", false, "warm up handshake")
	flag.Parse()

	InitHttpsClient()

	host = fmt.Sprintf("https://%s:%d", ip, port)

	err := PingServer()
	if err != nil {
		log.Panicln("Cannot ping server", err)
	}

	switch command {
	case "basic":
		basicTest()
	case "queryBench":
		queryBench()
	case "writeBench":
		writeBench()
	default:
		fmt.Println("invalid cmd")
	}
}
