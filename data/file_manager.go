package data

import (
	"bufio"
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

var backupFileName = "backup.txt"

// BuildDictionaryFromFile reads the contents of dictionary.txt and imports it into a Dictionary object
func BuildDictionaryFromFile() (*Dictionary, bool) {
	file, err := os.Open(fileName)
	if err != nil {
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
	if !copyFile(backupFileName, fileName) {
		return false
	}

	// Calling create on an existing file truncates it, and we want a fresh file here
	file, err := os.Create(fileName)
	if err != nil {
		return false
	}
	defer file.Close()

	dict.Sort()
	for _, entry := range dict.Entries {
		_, err := file.WriteString(fmt.Sprintf("%v\n", entry))
		if err != nil {
			copyFile(fileName, backupFileName)
			err = os.Remove(backupFileName)
			return false
		}
	}

	err = os.Remove(backupFileName)

	return true
}

// copyFile copies the contents of the files in the original filepath to the destination filepath
func copyFile(destination string, original string) bool {
	originalFile, err := os.Open(original)
	if err != nil {
		return false
	}
	defer originalFile.Close()

	destinationFile, err := os.Create(destination)
	if err != nil {
		return false
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, originalFile)
	if err != nil {
		return false
	}

	return true
}
