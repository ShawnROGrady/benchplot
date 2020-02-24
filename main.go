package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ShawnROGrady/benchplot/benchmark"
)

func main() {
	benches, err := benchmark.ParseBenchmarks(os.Stdin)
	if err != nil {
		log.Fatalf("error parsing input: %s", err)
	}

	fmt.Printf("len(benches) = %d\n", len(benches))
}
