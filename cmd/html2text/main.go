package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/JayJamieson/striphtml"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	opts := striphtml.Options{}
	out, err := striphtml.FromReader(reader, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(out)
}
