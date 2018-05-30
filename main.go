package main

import (
	"fmt"
	"log"
	"os"

	"github.com/carldanley/scribe/src/commands"
)

func main() {
	// setup the logger
	log.SetPrefix("==== [scribe] ==== ")
	log.SetFlags(0)

	// run the CLI and print any errors that occur
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
