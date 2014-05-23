package main

import . "github.com/visionmedia/go-debug"
import "time"

var a DebugFunction = Debug("multiple:a")
var b DebugFunction = Debug("multiple:b")

func work(debug DebugFunction) {
	for {
		debug("doing stuff")
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	go work(a)
	go work(b)

	a("quitting")
	time.Sleep(5 * time.Second)
}
