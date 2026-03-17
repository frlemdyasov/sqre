// Fedor Lemdyasov 2026

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	commandName := "seqre"

	// Declare command flags
	editor := flag.String("e", "emacs", "The text editor used to open the path paste file.") // (Must be a GUI program)
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
		fmt.Println("  -o [ORDER] \t \t Order the files will be renamed in. Either: name, date, or size")
		fmt.Println("  -n [NAME] \t \t Specify a custom name for the renamed files. Otherwise use the name of the first file in the order.")
		fmt.Println("  -z \t \t \t Make all the numbers the same length by adding zeroes in front of smaller numbers")
		fmt.Println("  -v \t \t \t Print the version number of the program")
		fmt.Println("  -h \t \t \t Print the usage information for this program")

		os.Exit(0)
	}

	// Check if the version flag is used, if so, exit the program
	if *version == true {
		fmt.Println(commandName + " version: 1.1")
		os.Exit(0)
	}

	// Check if the right order flag is used, if not exit the program
	if *order != "name" && *order != "date" && *order != "size" {
		fmt.Println("Error:", errors.New("Order of files not specified"))
		os.Exit(1)
	}

	// Create the text file that you paste in the image paths
	pasteFile, err := os.CreateTemp("", "sequenceRenamerPasteFile")
	if err != nil {
		fmt.Println("Error:", errors.New("Unable to make pasteFile."))
	}

	// Open the pasteFile with the selected editor
	cmd := exec.Command(*editor, pasteFile.Name())
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", errors.New("Text Editor failed to run."))
		fmt.Println("Make sure the editor does not run in a terminal.")
	}

	// Read the contents of the pasteFile
	input, err := os.ReadFile(pasteFile.Name())
	if err != nil {
		fmt.Println("Error:", errors.New("Unable to read the contents of the pasteFile."))
	}

	// Convert each line of the pasteFile into a list value
	filePaths := strings.Split(string(input), "\n")

	// Create the working directory
	workingFolder, err := os.MkdirTemp("", "sequenceRenamerImages")
	if err != nil {
		fmt.Println("Error:", errors.New("Unable to make the working directory."))
	}

	// Move the desired files into the working directory
	for i := range len(filePaths) {
		err = os.Rename(filePaths[i], workingFolder+"/"+filepath.Base(filePaths[i]))
		if err != nil {
			fmt.Println("Error:", errors.New("Unable to move files to working directory."))
		}
	}

	// Detect the files in the working directory
	entries, err := os.ReadDir(workingFolder)
	if err != nil {
		fmt.Println("Error:", errors.New("Failed to read the working directory."))
	}

	// Copy the original path of the first selected file before entries is sorted
	originalPath := filePaths[0][0 : len(filePaths[0])-len(filepath.Base(filePaths[0]))]
	fmt.Println(originalPath)

	switch *order {
	case "name":
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].Name() < entries[j].Name()
		})
	case "date":
		sort.SliceStable(entries, func(i, j int) bool {
			entryI, _ := entries[i].Info()
			entryJ, _ := entries[j].Info()
			return entryI.ModTime().Before(entryJ.ModTime())
		})
	case "size":
		sort.SliceStable(entries, func(i, j int) bool {
			entryI, _ := entries[i].Info()
			entryJ, _ := entries[j].Info()
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

		newName := originalPath + firstFileExtentionlessName + addZeroes + strconv.Itoa(i+1) + fileExtension
		fmt.Println(workingFolder + entries[i].Name())
		fmt.Println(newName)

		os.Rename(workingFolder+"/"+entries[i].Name(), newName)

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
