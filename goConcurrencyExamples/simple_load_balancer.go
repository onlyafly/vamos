package goConcurrencyExamples

import (
	"fmt"
	"time"
)

// Helpers

const NumWorkers = 42

func sendLotsOfWork(c chan<- *Work) {
	for i := 0; i < 100; i++ {
		c <- &Work{x: i, y: i * 2, z: 0}
	}
}

func receiveLotsOfResults(c <-chan *Work) {
	for {
		w := <-c
		fmt.Printf("Result: %v\n", w.z)
	}
}

// From: http://talks.golang.org/2012/waza.slide#39

// A unit of work

type Work struct {
	x, y, z int
}

// A worker task

func worker(in <-chan *Work, out chan<- *Work) {
	for w := range in {
		w.z = w.x * w.y
		time.Sleep(time.Duration(w.z))
		out <- w
	}
}

// Task: Must make sure other workers can run when one blocks.

// The runner
func SimpleLoadBalancerMain() {
	in, out := make(chan *Work), make(chan *Work)
	for i := 0; i < NumWorkers; i++ {
		go worker(in, out)
	}
	go sendLotsOfWork(in)
	receiveLotsOfResults(out)
}

// Easy problem but also hard to solve concisely without concurrency.
