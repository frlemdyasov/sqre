// Fedor Lemdyasov 2026

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	commandName := "sqre"
	versionNumber := "2.0"

	// Declare command flags
	editor := flag.String("e", "emacs", "The text editor used to open the path paste file.") // (Must be a GUI program)
	location := flag.String("l", ".", "The directory the renamed files will be placed in.")
	order := flag.String("o", "date", "Order the files will be renamed in. Either: name, date, or size")
	name := flag.String("n", "", "Specify a custom name for the renamed files. Leaving this flag empty defaults to using the name of the first file in the order.")
	cli := flag.Bool("cli", false, "Start the Command Line Interface. This enables the other options.")
	zeroes := flag.Bool("z", false, "Make all the numbers the same length by adding zeroes in front of smaller numbers")
	version := flag.Bool("v", false, "Print the version number of the program")
	help := flag.Bool("h", false, "Print the usage information for this program.")
	flag.Parse()

	// Check if the help flag is being used, if so, exit the program
	if *help == true {
		fmt.Println("Usage: " + commandName + " [OPTIONS]...\n")
		fmt.Println("List of options:")
		fmt.Println("  -cli \t \t \t Start the Command Line Interface. This enables the other options.")
		fmt.Println("  -e [EDITOR] \t \t The text editor used to open the path paste file.")
		fmt.Println("  -l [LOCATION] \t The directory the renamed files will be placed in.")
		fmt.Println("  -o [ORDER] \t \t Order the files will be renamed in. Either: name, date, or size")
		fmt.Println("  -n [NAME] \t \t Specify a custom name for the renamed files. Otherwise use the name of the first file in the order.")
		fmt.Println("  -z \t \t \t Make all the numbers the same length by adding zeroes in front of smaller numbers.")
		fmt.Println("  -v \t \t \t Print the version number of the program.")
		fmt.Println("  -h \t \t \t Print the usage information for this program.")
		fmt.Println("\n Example: \n sqre -e gedit -l ~Downloads")

		os.Exit(0)
	}

	// Check if the GUI needs to be launched, if so disregard launching a text editor
	if *cli == false {
		sqreGui()
	} else {
		// Check if the order flag is used, if not exit the program
		if *order != "name" && *order != "date" && *order != "size" {
			fmt.Println("Error:", errors.New("Order of files not specified"))
			os.Exit(1)
		}

		// Check if the version flag is used, if so, exit the program
		if *version == true {
			fmt.Println(commandName+" version: ", versionNumber)
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

		// Run the sorting script
		SequenceRename(string(input), *location, *order, *name, *zeroes)
	}
}
