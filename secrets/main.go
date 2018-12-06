package main

import (
	"fmt"
	"os"

	secretly "github.com/evanphx/go-secretly"
	_ "github.com/evanphx/go-secretly/all"
)

func main() {
	switch len(os.Args) {
	case 2:
		val, err := secretly.Get(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred: %s\n", err)
			os.Exit(1)
		}

		fmt.Println(val)
	case 3:
		err := secretly.Put(os.Args[1], os.Args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "An error occurred: %s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "secrets [path] [value]\n")
		os.Exit(1)
	}
}
