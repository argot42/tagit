package main

import (
	"fmt"
	"os"
	"github.com/argot42/tagit/flags"
)

func checkError (err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func warning (tag string, obj string, name string, err error) {
	fmt.Fprintf(os.Stderr, "Warning adding tag (%s) to %s (%s): %s\n", tag, obj, name, err.Error())
}

func warningFile (tag string, filepath string, err error) {
	if flags.Args.Verbose { warning(tag, "file", filepath, err) }
}

func warningDict (tag string, dictname string, err error) {
	if flags.Args.Verbose { warning(tag, "dictionary", dictname, err) }
}
