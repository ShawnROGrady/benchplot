package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	benches, err := parseBenchmarks(os.Stdin)
	if err != nil {
		log.Fatalf("error parsing input: %s", err)
	}

	for _, bench := range benches {
		fmt.Printf("%s\n", bench.name)
		for _, res := range bench.results {
			fmt.Println("\tInputs")
			fmt.Println("\t\tvarValues:")
			for _, varValue := range res.inputs.varValues {
				fmt.Printf("\t\t\t%s\n", varValue)
			}
			fmt.Printf("\t\tsubs: %v\n", res.inputs.subs)
			fmt.Println("\tOutputs")
			fmt.Printf("\t\t%v\n", res.outputs)
		}
	}
}
