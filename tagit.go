package main

import (
	"fmt"
	"github.com/argot42/tagit/dictionary"
	"github.com/argot42/tagit/flags"
	"github.com/argot42/tagit/options"
	"os"
)

func main() {
	err := flags.Init_flags()
	checkError(err)

	if flags.Args.Verbose {
		fmt.Println("Loading default dictionary")
	}

	defaultDictionary, err := dictionary.LoadDictionary(flags.Args.DB)

	switch {
	// if file doesn't exist create a new one by default
	// or let the user select an empty dictionary
	case os.IsNotExist(err):
		// user interaction
		if flags.Args.Interactive || flags.Args.Client {
			// create new dictionary (answer: yes)
			if options.YesNo(options.Yes, "Default dictionary not found, create a new one on (%s)", flags.Args.DB) {
				defaultDictionary, err = dictionary.CreateDictionary(flags.Args.DB)
				// if there's an error exit
				checkError(err)
			}
		} else { // no user interaction
			defaultDictionary, err = dictionary.CreateDictionary(flags.Args.DB)
			// if there's an error and has the verbose flag print a warn
			// and continue with an empty dictionary
			if err != nil {
				if flags.Args.Verbose {
					fmt.Fprintf(os.Stderr, "warning: %s\n", err.Error())
				}
			}
		}

	// if there's another kind of error trying to load the default dictionary
	// print a warning if verbose and continue with an empty dictionary
	case err != nil:
		if flags.Args.Verbose {
			fmt.Fprintf(os.Stderr, "warning: %s\n", err.Error())
		}
	}

	if flags.Args.Client {
		client(defaultDictionary)
	} else {
		batch(defaultDictionary)
	}
}
