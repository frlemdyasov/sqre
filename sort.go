// Fedor Lemdyasov 2026

package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func SequenceRename(input string, location string, order string, name string, zeroes bool) {

	var filePaths []string

	// Convert each line of the pasteFile into a list value
	filePaths = strings.Split(input, "\n")

	// Copy the original path of the first selected file before entries are sorted, and before the first file operation
	originalPath := filePaths[0][0 : len(filePaths[0])-len(filepath.Base(filePaths[0]))]

	// Check if the location flag is used
	var newPath string
	if location == "." {
		newPath = originalPath
	} else {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("Error:", errors.New("Failed to determine executable location"))
			os.Exit(1)
		}
		var slash rune = '/'
		if a := []rune(ex); a[0] == slash {
			newPath = location
		} else {
			newPath = filepath.Dir(ex) + location
		}
	}

	// Read each file's info
	var entries []fs.FileInfo
	for i := range filePaths {
		fileInfo, err := os.Stat(filePaths[i])
		if err != nil {
			fmt.Println("Error:", errors.New("Failed to read the file info."))
		}
		entries = append(entries, fileInfo)
	}

	// Sort the files depending on chosen flag
	switch order {
	case "name":
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].Name() < entries[j].Name()
		})
	case "date":
		sort.SliceStable(entries, func(i, j int) bool {
			entryI := entries[i]
			entryJ := entries[j]
			return entryI.ModTime().Before(entryJ.ModTime())
		})
	case "size":
		sort.SliceStable(entries, func(i, j int) bool {
			entryI := entries[i]
			entryJ := entries[j]
			return entryI.Size() < entryJ.Size()
		})
	}

	var firstFileExtentionlessName string
	if name != "" {
		firstFileExtentionlessName = name
	} else {
		// Rename the files, so that each file carries the same name as the first file in the order.
		firstFileExtentionlessName = entries[0].Name()[0 : len(entries[0].Name())-len(filepath.Ext(entries[0].Name()))]
	}

	// Write a number at the end of the name denoting its order in the sequence.
	// Renaming also changes the path of the file, back to its original location
	var addZeroes string
	newName := make([]string, len(entries))
	for i := range entries {
		fileExtension := filepath.Ext(entries[i].Name())

		if zeroes == true {
			addZeroes = countZeroes(i+1, len(entries))
		} else {
			addZeroes = ""
		}

		newName[i] = newPath + firstFileExtentionlessName + addZeroes + strconv.Itoa(i+1) + fileExtension

		err := os.Rename(originalPath+"/"+entries[i].Name(), newName[i])
		if err != nil {
			fmt.Println("Error:", errors.New("Unable to move files into the directory:"))
			fmt.Println(newPath)
			os.Exit(1)
		}

		fmt.Println(entries[i].Name(), "->", newName[i])

	}

}

// Function to add zeroes in front of the sequence number
func countZeroes(num int, tot int) string {
	numZeros := len(strconv.Itoa(tot)) - len(strconv.Itoa(num))
	zeroes := ""
	for range numZeros {
		zeroes = zeroes + "0"
	}
	return zeroes
}

// Function to manually move a file: copy then delete the original
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	fileInfo, err := inputFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}
	fmt.Println(fileInfo.Mode())

	outputFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close() // for Windows, close before trying to remove: https://stackoverflow.com/a/64943554/246801

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source files: %v", err)
	}

	return nil
}
