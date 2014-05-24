package main

import . "github.com/visionmedia/go-debug"
import "time"

var a = Debug("multiple:a")
var b = Debug("multiple:b")

func work(debug DebugFunction, delay time.Duration) {
	for {
		debug("doing stuff")
		time.Sleep(delay)
	}
}

func main() {
	go work(a, 1000*time.Millisecond)
	go work(b, 100*time.Millisecond)

	time.Sleep(5 * time.Second)
}
