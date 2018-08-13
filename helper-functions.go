package main

import (
	"fmt"
	"os"
	"time"
)

func dieOnError(err error) {
	if err != nil {
		printError(err)
		os.Exit(1)
	}
}

func outputError(err error) {
	if err != nil {
		printError(err)
	}
}

func printError(err error) {
	errorMessage := fmt.Sprintf("%s Error: %v", time.Now().Format("2006/01/02 15:04:05"), err)
	fmt.Fprintf(os.Stderr, "%s\n", errorMessage)
}
