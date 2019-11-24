package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/common-nighthawk/go-figure"
	b "github.com/graniet/bfc/core"
	"log"
	"os"
)

const (
	VERSION = "1.0"
)

func main() {
	parser := argparse.NewParser(
		"bff",
		"Best Friend For Curious (BFC) - This tools execute a routine created by users.",
		)

	verbose := parser.Flag("v", "verbose", &argparse.Options{
		Required: false,
		Help: "Don't print a BestFriend banner",
	})

	routine := parser.String("s", "source", &argparse.Options{
		Required: true,
		Help:     "Execute selected routine",
	})

	parameters := parser.String("p", "parameters", &argparse.Options{
		Required: false,
		Help: "Set dynamically routine parameters",
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if !*verbose {
		banner := figure.NewFigure("BFC", "", true)
		banner.Print()
		fmt.Println("A best friend for curious")
		fmt.Println()
	}

	bff, err := b.NewBffExecution(*routine)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	bff.Execute(*parameters)
	return
}