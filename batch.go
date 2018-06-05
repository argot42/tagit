package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/argot42/tagit/flags"
	"github.com/argot42/tagit/dictionary"
	"github.com/argot42/tagit/options"
	"github.com/argot42/tagit/fileproc"
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
		// if dictionary is not empty ask do a fuzzy search
		if dict.Initialized() {
			// fuzzy search (that fails at exact match)
			matches, err := file.FFind(tag, dict, fileproc.ExactFail)

			// if exact match add to word
			if err == fileproc.ExactMatchErr {
				if options.YesNo(options.Yes, "Add tag (%s) to file (%s) ?", tag, file.Name()) {

					// adding tag to file
					if err := file.Add(tag); err != nil { warningFile(tag, file.Path(), err) }
				}
				continue

			} else if err != nil { // if it is another type of error warn
				if flags.Args.Verbose { 
					fmt.Fprintf(os.Stderr, "Warning fuzzy match (%s): %s\n", tag, err.Error()) 
				}
				continue
			}

			// if no matches add tag to file and dictionary (ask user)
			if len(matches) == 0 {
				if options.YesNo(options.Yes, "Add tag (%s) to file (%s) ?", tag, file.Name()) {
					
					// adding tag to file
					if err := file.Add(tag); err != nil { warningFile(tag, file.Path(), err) }

					if options.YesNo(options.Yes, "Add tag (%s) to dictionary (%s) ?", tag, dict.Name()) {
						if err := dict.Add(tag); err != nil { warningDict(tag, dict.Name(), err) }
					}
				}
				continue
			}

			// if matches
			matches_plus_tag := append(matches, tag)
			o := options.ChooseNumeric(0, matches_plus_tag, "choose tag")	

			if o >= 0 && o < len(matches_plus_tag) {
				if err := file.Add(matches_plus_tag[o]); err != nil { warningFile(tag, file.Path(), err) }

				if options.YesNo(options.Yes, "Add tag (%s) to dictionary (%s) ?", tag, dict.Name()) {
					if err := dict.Add(matches_plus_tag[o]); err != nil {
						warningDict(tag, dict.Name(), err)
					}
				}
			}
		}
	}
}

func noninteractive (filepath string) {
	fmt.Println(filepath)
}
