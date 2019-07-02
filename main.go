package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

func main() {
	flag.Parse()
	err := Run(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Run is the main execution of the application
func Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no URL specified")
	}

	target, err := url.ParseRequestURI(args[0])
	if err != nil {
		return fmt.Errorf("invalid URL")
	}

	fmt.Println(target)

	return nil
}
