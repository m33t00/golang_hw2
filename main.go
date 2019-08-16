package main

import (
	"fmt"
	"time"
)

func main() {
	mainJobs := []job{
		job(func(in, out chan interface{}) {
			for _, i := range []int{0, 1, 1, 2, 3, 5, 8} {
				out <- i
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			fmt.Println(<-in)
		}),
	}

	start := time.Now()
	ExecutePipeline(mainJobs...)

	fmt.Println("execution time (seconds):", time.Since(start))
}
