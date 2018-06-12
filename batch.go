package main

import (
	"fmt"
	"github.com/argot42/tagit/dictionary"
	"github.com/argot42/tagit/fileproc"
	"github.com/argot42/tagit/flags"
	"github.com/argot42/tagit/options"
	"github.com/argot42/tagit/utils"
	"os"
	"path/filepath"
	"sync"
)

func batch(defaultDict dictionary.Dictionary) {
	if flags.Args.Verbose {
		fmt.Println("Starting batch tagging")
	}

	if flags.Args.Interactive {
		interactive(defaultDict)
	} else {
		noninteractive(defaultDict)
	}
}

func interactive(defaultDict dictionary.Dictionary) {
	for _, item := range flags.Args.Files {
		dict := loadd(item, defaultDict)

		// parse file
		file, err := fileproc.Parse(item)
		if err != nil {
			warning("Warning while parsing (%s): %s\n", item, err.Error())
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
					o := options.ChooseNumeric(0, matches_plus_tag, "Choose tag ")

					if o < 0 || o >= len(matches_plus_tag) {
						continue
					}

					tagToAdd = matches_plus_tag[o]
					if o == 0 {
						validToDict = true
					}
				}
			}

			// add tag to file (ask user)
			if options.YesNo(options.Yes, "Add tag (%s) to file (%s) ?", tagToAdd, file.Name()) {
				//if err := file.Add(tagToAdd); err != nil { warningFile(tagToAdd, file.Path(), err) }
				file.Add(tagToAdd)

				// add tag to dictionary if flags allows it (ask user)
				if !validToDict {
					continue
				}
				if options.YesNo(options.Yes, "Add tag (%s) to dictionary (%s) ?", tagToAdd, dict.Name()) {
					if err := dict.Add(tagToAdd); err != nil {
						warningDict(tagToAdd, dict.Name(), err)
					}
				}
			}
		}

		// write tags to filename
		if err := file.Write(); err != nil {
			warningFile(file.Path(), err)
		}
	}
}

func noninteractive(defaultDict dictionary.Dictionary) {
	info, err := concur(defaultDict)

out:
	for {
		select {
		case i, more := <-info:
			if !more {
				break out
			}
			if flags.Args.Verbose {
				fmt.Print(i)
			}
		case e := <-err:
			warning(e.Error())
		}
	}
}

func concur(defaultDict dictionary.Dictionary) (info chan string, err chan error) {
	var wg sync.WaitGroup
	info = make(chan string, len(flags.Args.Files))
	err = make(chan error, len(flags.Args.Files))

	// process files
	wg.Add(len(flags.Args.Files))
	for _, f := range flags.Args.Files {
		go processFile(f, loadd(f, defaultDict), info, err, &wg)
	}

	// wait until all goroutines finish and then close channels
	// to signal main to end
	go func() {
		wg.Wait()
		close(info)
	}()

	return
}

func processFile(path string, dict dictionary.Dictionary, info chan string, e chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := fileproc.Parse(path)
	if err != nil {
		e <- fmt.Errorf("Warning while parsing (%s): %s\n", path, err.Error())
		return
	}

	for _, tag := range flags.Args.Tags {
		file.Add(tag)

		// writing to dict
		if err = dict.Add(tag); err != nil {
			e <- fmt.Errorf("Warning adding tag (%s) to dictionary (%s): %s\n", tag, dict.Name(), err.Error())
			continue
		}

		info <- fmt.Sprintf("Tag (%s) added to dictionary (%s)!\n", tag, dict.Name())
	}

	if err = file.Write(); err != nil {
		e <- fmt.Errorf("Warning writing tags to file (%s) failed\n", file.Name())
		return
	}

	info <- fmt.Sprintf("Tags written to file (%s)!\n", file.Name())
}

func loadd(filename string, defaultDict dictionary.Dictionary) (currentDict dictionary.Dictionary) {
	if flags.Args.Verbose {
		fmt.Println("Loading local dictionary")
	}

	itemDir := filepath.Dir(filename)
	currentDict, err := dictionary.LoadDictionary(itemDir)

	switch {
	// if file doesn't exist use defaultDict
	// or ask the user to create a new empty one for that directory
	case os.IsNotExist(err):
		// user interaction
		if flags.Args.Interactive {
			// create a new dictionary (answer: no)
			if !options.YesNo(options.Yes, "Dictionary not found on (%s), use default?", itemDir) {
				currentDict, err = dictionary.CreateDictionary(filepath.Join(itemDir, ".dict.db"))
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
		warning("Warning %s\n", err.Error())
		currentDict = defaultDict
	}

	return
}
