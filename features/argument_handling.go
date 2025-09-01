package features

import (
	"flag"
	"fmt"
	"os"
)

func Handle_arguments() {
	argsWithProg := os.Args
	// argsWithoutProg := argsWithProg[1:]
	fmt.Println(argsWithProg)
	heightPtr := flag.Uint("height", 172, "height in centimeters")
	flag.Parse()
	height := float64(*heightPtr)
	lengthFactor := height / 172
	cubeFactor := lengthFactor * lengthFactor * lengthFactor
	weight := 60 * cubeFactor
	fmt.Println("Weight: ", weight, " kg")

}
