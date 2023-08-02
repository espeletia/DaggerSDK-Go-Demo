package main

import (
	"context"
	"coolProject/ci/runner"
	"fmt"
	"os"
)

func main() {
	mode := os.Args[1:]
	if len(mode) == 0 {
		fmt.Println("Usage: ci [build|test]")
		return
	}
	if mode[0] == "build" {
		if err := runner.Build(context.Background()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if mode[0] == "test" {
		if err := runner.DoTests(context.Background()); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage: ci [build|test]")
		return
	}
}
