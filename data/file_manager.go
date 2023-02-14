package data

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"os"
)

/**
 * File Format:
 * n [such that n is the number of items in this file, excluding the first line]
 * entry 1
 * ...
 * entry n
 */
var fileName = "dictionary.txt"

//go:embed dictionary.txt
var embeddedDictionary []byte

// BuildDictionaryFromFile reads the contents of dictionary.txt and imports it into a Dictionary object
func BuildDictionaryFromFile() (*Dictionary, bool) {
	file, ok := getDictionaryFile()
	if !ok {
		return nil, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	dict := new(Dictionary)
	for scanner.Scan() {
		dict.Add(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, false
	}

	return dict, true
}

// SaveDictionaryToFile accepts a dictionary reference and then writes its data to the configured output file
func SaveDictionaryToFile(dict *Dictionary) bool {
	// Calling create on an existing file truncates it, and we want a fresh file here
	file, ok := getDictionaryFile()
	if !ok {
		return false
	}
	defer file.Close()

	dict.Sort()
	for _, entry := range dict.Entries {
		_, err := file.WriteString(fmt.Sprintf("%v\n", entry))
		if err != nil {
			return false
		}
	}

	return true
}

// getDictionaryFile will pull embedded dictionary file into build for a starting point, then clone it into running
// service. This would need to be replaced with a database or something persistent in production
func getDictionaryFile() (*os.File, bool) {
	if _, err := os.Stat(fileName); err == nil {
		file, err := os.Open(fileName)
		return file, err == nil
	}

	file, createErr := os.Create(fileName)
	if createErr != nil {
		return nil, false
	}

	_, writeErr := file.Write(embeddedDictionary)
	if writeErr != nil {
		return nil, false
	}

	_, seekErr := file.Seek(0, io.SeekStart)
	if seekErr != nil {
		return nil, false
	}

	return file, true
}
