package main

import (
	"flag"
	"fmt"

	"github.com/third-net/euler"
)

var (
	addr = flag.String("b", ":10010", "euler dns service address")
)

func main() {
	flag.Parse()
	fmt.Print(euler.ListenAndServe(*addr))
}
