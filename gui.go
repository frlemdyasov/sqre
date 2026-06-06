// Fedor Lemdyasov 2026

package main

import (
	"os"

	_ "embed"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed main.ui
var uiXML string

func sqreGui() {
	app := gtk.NewApplication("com.github.frlemdyasov.sqre", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	builder := gtk.NewBuilderFromString(uiXML)

	window := builder.GetObject("MainWindow").Cast().(*gtk.Window)
	renameButton := builder.GetObject("RenameButton").Cast().(*gtk.Button)
	zeroCheck := builder.GetObject("AddZeroesCheck").Cast().(*gtk.CheckButton)
	setOrder := builder.GetObject("OrderByChooser").Cast().(*gtk.DropDown)
	setName := builder.GetObject("FileNamer").Cast().(*gtk.Entry)
	pasteButton := builder.GetObject("pasteButon").Cast().(*gtk.Button)
	filePaths := builder.GetObject("TextInput").Cast().(*gtk.TextView)

	var order string
	location := "."

	// Paste text from clipboard into
	pasteButton.ConnectClicked(func() {

		clipboardObject := (*gdk.Display).Clipboard(gdk.DisplayGetDefault())
		filePaths.Buffer().PasteClipboard(clipboardObject, filePaths.Buffer().StartIter(), true)

	})

	renameButton.ConnectClicked(func() {

		// Read whether to add zeroes
		zeroes := zeroCheck.Active()

		// Read the selected order
		switch int(setOrder.Selected()) {
		case 0:
			order = "date"

		case 1:
			order = "name"

		case 2:
			order = "size"
		}

		// Read the specified file name
		name := setName.Buffer().Text()

		// Read the text box for the file paths
		input := filePaths.Buffer().Text(filePaths.Buffer().StartIter(), filePaths.Buffer().EndIter(), false)

		// Clear the text box
		filePaths.Buffer().Delete(filePaths.Buffer().StartIter(), filePaths.Buffer().EndIter())

		// Run the sorting script
		SequenceRename(input, location, order, name, zeroes)

	})

	app.AddWindow(window)
	window.SetVisible(true)
}
