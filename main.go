// Fedor Lemdyasov 2026

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	commandName := "sqre"

	// Declare command flags
	editor := flag.String("e", "emacs", "The text editor used to open the path paste file.") // (Must be a GUI program)
	location := flag.String("l", ".", "The directory the renamed files will be placed in.")
	order := flag.String("o", "date", "Order the files will be renamed in. Either: name, date, or size")
	name := flag.String("n", "", "Specify a custom name for the renamed files. Leaving this flag empty defaults to using the name of the first file in the order.")
	zeroes := flag.Bool("z", false, "Make all the numbers the same length by adding zeroes in front of smaller numbers")
	version := flag.Bool("v", false, "Print the version number of the program")
	help := flag.Bool("h", false, "Print the usage information for this program.")
	flag.Parse()

	// Check if the help flag is being used, if so, exit the program
	if *help == true {
		fmt.Println("Usage: " + commandName + " [OPTIONS]...")
		fmt.Println("")
		fmt.Println("List of options:")
		fmt.Println("  -e [EDITOR] \t \t The text editor used to open the path paste file.")
		fmt.Println("  -l [LOCATION] \t The directory the renamed files will be placed in.")
		fmt.Println("  -o [ORDER] \t \t Order the files will be renamed in. Either: name, date, or size")
		fmt.Println("  -n [NAME] \t \t Specify a custom name for the renamed files. Otherwise use the name of the first file in the order.")
		fmt.Println("  -z \t \t \t Make all the numbers the same length by adding zeroes in front of smaller numbers")
		fmt.Println("  -v \t \t \t Print the version number of the program")
		fmt.Println("  -h \t \t \t Print the usage information for this program")
		fmt.Println("\n Example: \n sqre -e gedit -l ~Downloads")

		os.Exit(0)
	}

	// Check if the order flag is used, if not exit the program
	if *order != "name" && *order != "date" && *order != "size" {
		fmt.Println("Error:", errors.New("Order of files not specified"))
		os.Exit(1)
	}

	// Check if the version flag is used, if so, exit the program
	if *version == true {
		fmt.Println(commandName + " version: 1.2")
		os.Exit(0)
	}

	// Create the text file that you paste in the image paths
	pasteFile, err := os.CreateTemp("", "sequenceRenamerPasteFile")
	if err != nil {
		fmt.Println("Error:", errors.New("Unable to make pasteFile."))
		os.Exit(1)
	}

	// Open the pasteFile with the selected editor
	cmd := exec.Command(*editor, pasteFile.Name())
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", errors.New("Text Editor failed to run."))
		fmt.Println("Make sure the editor does not run in a terminal.")
		os.Exit(1)
	}

	// Read the contents of the pasteFile
	input, err := os.ReadFile(pasteFile.Name())
	if err != nil {
		fmt.Println("Error:", errors.New("Unable to read the contents of the pasteFile."))
	}

	// Remove temporary pasteFile
	err = os.Remove(pasteFile.Name())
	if err != nil {
		fmt.Println("Error:", errors.New("Couldn't remove temporary text file"))
		os.Exit(1)
	}

	// Convert each line of the pasteFile into a list value
	filePaths := strings.Split(string(input), "\n")

	// Copy the original path of the first selected file before entries is sorted, and before the first file operation
	originalPath := filePaths[0][0 : len(filePaths[0])-len(filepath.Base(filePaths[0]))]

	// Check if the location flag is used
	var newPath string
	if *location == "." {
		newPath = originalPath
	} else {
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("Error:", errors.New("Failed to determine executable location"))
			os.Exit(1)
		}
		var slash rune = '/'
		if a := []rune(ex); a[0] == slash {
			newPath = *location
		} else {
			newPath = filepath.Dir(ex) + *location
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

	// sort the files depending on chosen flag
	switch *order {
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
	if *name != "" {
		firstFileExtentionlessName = *name
	} else {
		// Rename the files, so that each file carries the same name as the first file in the order.
		firstFileExtentionlessName = entries[0].Name()[0 : len(entries[0].Name())-len(filepath.Ext(entries[0].Name()))]
	}

	// Write a number at the end of the name denoting its order in the sequence.
	// Renaming also changes the path of the file, back to its original location
	var addZeroes string
	for i := range entries {
		fileExtension := filepath.Ext(entries[i].Name())

		if *zeroes == true {
			addZeroes = countZeroes(i+1, len(entries))
		} else {
			addZeroes = ""
		}

		newName := newPath + firstFileExtentionlessName + addZeroes + strconv.Itoa(i+1) + fileExtension

		err = os.Rename(originalPath+"/"+entries[i].Name(), newName)
		if err != nil {
			fmt.Println("Error:", errors.New("Unable to move files into the directory:"))
			fmt.Println(newPath)
			os.Exit(1)
		}

	}

}

func countZeroes(num int, tot int) string {
	numZeros := len(strconv.Itoa(tot)) - len(strconv.Itoa(num))
	zeroes := ""
	for range numZeros {
		zeroes = zeroes + "0"
	}
	return zeroes
}

// https://stackoverflow.com/questions/50740902/move-a-file-to-a-different-drive-with-go
func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	// https://www.programmershelp.net/go/tutorial-on-copying-files-in-go.php
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
