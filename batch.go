package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/argot42/tagit/flags"
	"github.com/argot42/tagit/dictionary"
	"github.com/argot42/tagit/options"
	"github.com/argot42/tagit/fileproc"
	"github.com/argot42/tagit/utils"
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
			if flags.Args.Interactive {
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
			if flags.Args.Verbose { fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error()) }
			currentDict = defaultDict
		}

		// start processing filenames
		if flags.Args.Interactive {
			interactive(item, currentDict)
		} else {
			noninteractive(item)
		}
	}
}

func interactive (filepath string, dict dictionary.Dictionary) {
	file, err := fileproc.Parse(filepath)
	if err != nil {
		if flags.Args.Verbose { fmt.Fprintf(os.Stderr, "Warning while parsing: %s\n", err.Error()) }
		return
	}

	for _, tag := range flags.Args.Tags {
		tagToAdd := tag
		validToDict := false // flag used to ask to add or not to dictionary

		// if dictionary is initialized do a fuzzy search
		if dict.Initialized() {
			// fuzzy sarch (that returns an error on exact match)
			matches, err := dict.FFind(tag, dictionary.ExactFail)

			switch {
			case err != nil && err != dictionary.ExactMatchErr:
				warning("Warning fuzzy match (%s): %s\n", tag, err.Error())
				continue
			case err == dictionary.ExactMatchErr:
				// if exact match do nothing
			case len(matches) == 0:
				// if there're no matches
				validToDict = true
			default:
				// if there're matches
				matches_plus_tag := utils.Prepend(tag, matches)
				o := options.ChooseNumeric(0, matches_plus_tag, "choose tag")

				if o < 0 || o >= len(matches_plus_tag) { continue }

				tagToAdd = matches_plus_tag[o]
				validToDict = true
			}
		}

		// add tag to file (ask user)
		if options.YesNo(options.Yes, "Add tag (%s) to file (%s) ?", tagToAdd, file.Name()) {
			if err := file.Add(tagToAdd); err != nil { warningFile(tagToAdd, file.Path(), err) }
			
			// add tag to dictionary if flags allows it (ask user)
			if !validToDict { continue }
			if options.YesNo(options.Yes, "Add tag (%s) to dictionary (%s) ?", tagToAdd, dict.Name()) {
				if err := dict.Add(tagToAdd); err != nil { warningDict(tagToAdd, dict.Name(), err) }	
			}
		}
	}
}

func noninteractive (filepath string) {
	fmt.Println(filepath)
}
