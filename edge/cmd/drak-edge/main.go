// Command drak-edge is the DRAK Edge runtime. S0 skeleton: it only reports
// where the shared contract lives.
package main

import (
	"fmt"
	"os"

	"github.com/Kaiser1308/drak-home/edge/internal/contract"
)

func main() {
	dir, err := contract.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("drak-edge (skeleton); contract:", dir)
}
