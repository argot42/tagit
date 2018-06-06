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

func warning (format string, a ...interface{}) {
	if flags.Args.Verbose {
		fmt.Fprintf(os.Stderr, format, a)
	}
}

func warningTag (tag string, obj string, name string, err error) {
	warning("Warning adding tag (%s) to %s (%s): %s\n", tag, obj, name, err.Error())
}

func warningFile (tag string, filepath string, err error) {
	warningTag(tag, "file", filepath, err)
}

func warningDict (tag string, dictname string, err error) {
	warningTag(tag, "dictionary", dictname, err)
}
