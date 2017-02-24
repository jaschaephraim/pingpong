// Executable pingpong reads urls for gifs from a text file, and concurrently converts them to
// "pingpong" style gifs.
//
// The executable accepts three flags:
//  --urls  path to a local text file containing a newline separated list of urls
//  --dir   path to a local directory in which the processed gifs should be saved
//  --trans optional boolean value that indicates whether transparency should be corrected
//          on reversal, avoiding a ghosting effect on certain gifs
//
// Example:
//   pingpong --urls ./urls.txt --dir ./gifs --trans
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	urlsFile, dir, trans, err := parseFlags()
	if err != nil {
		logError(err)
		os.Exit(1)
	}

	urls, err := readURLs(urlsFile)
	if err != nil {
		logError(err)
		os.Exit(1)
	}

	results := make(chan *img)
	for _, url := range urls {
		i := newImg(url)
		go i.process(dir, trans, results)
	}

	for range urls {
		i := <-results
		i.logResult()
	}
	fmt.Println("Done")
}

// Parse, validate. and return flags.
func parseFlags() (string, string, bool, error) {
	urlsFlag := flag.String(
		"urls", "",
		"path to a local text file containing a newline separated list of urls")
	dirFlag := flag.String(
		"dir", "",
		"path to a local directory in which the processed gifs should be saved")
	transFlag := flag.Bool(
		"trans", false,
		"optional boolean value that indicates whether transparency should be corrected on reversal, "+
			"avoiding a ghosting effect on certain gifs")
	flag.Parse()

	trans := *transFlag

	urlsFile := *urlsFlag
	if urlsFile == "" {
		return "", "", trans, errors.New("Missing --urls flag")
	}

	dir := *dirFlag
	if dir == "" {
		return urlsFile, "", trans, errors.New("Missing --dir flag")
	}

	return urlsFile, dir, trans, nil
}

// Read and parse urls from file.
func readURLs(filename string) ([]string, error) {
	urlsBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	trimmed := strings.TrimSuffix(string(urlsBytes), "\n")
	urls := strings.Split(trimmed, "\n")
	return urls, nil
}
