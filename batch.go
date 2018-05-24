package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/argot42/tagit/flags"
	"github.com/argot42/tagit/dictionary"
	"github.com/argot42/tagit/options"
)

func batch (defaultDict dictionary.Dictionary) {
	if flags.Args.Verbose { fmt.Println("Starting batch tagging") }

	for _, item := range flags.Args.Files {
		if flags.Args.Verbose { fmt.Println("Loading local dictionary") }
		
		var currentDict dictionary.Dictionary
		itemDir := filepath.Dir(item)
		currentDict, err := dictionary.LoadDictionary(itemDir)

		switch {
		// if file doesn't exist use defaultDict by default
		// or ask the user to create a new empty one for that directory
		case os.IsNotExist(err):
			// user interaction
			if flags.Args.Interactive || flags.Args.Client {
				// create new dictionary (answer: no)
				if !options.YesNo(options.Yes, "Dictionary not found on (%s), use default?", itemDir) {
					currentDict, err = dictionary.CreateDictionary(itemDir)
					// if there's an error exit
					checkError(err)
				} else { // use default (answer: yes)
					currentDict = defaultDict
				}
			} else { // no user interaction
				currentDict = defaultDict
			}
		// if there's another kind of error reading the dictionary file
		// print a warning if verbose and continue with the default dictionary
		case err != nil:
			if flags.Args.Verbose { fmt.Fprintf(os.Stderr, "warning: %s\n", err.Error()) }
			currentDict = defaultDict
		}

		fmt.Println("DIR:", currentDict)
	}
}
