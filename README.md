# Sequence Renamer
Version 2.0

Inspired by [renameutils](https://www.nongnu.org/renameutils/), sqre is a bulk renaming tool that prioritizes speed and convenience.

![](https://github.com/frlemdyasov/sqre/blob/main/sqreScreenshot.png)

> [!NOTE]
> This program was designed for my own personal use case, so alternate functionality may not be added. 

## Features
- Sequentially rename files in order of: Name, Modification date, Size
- Set a custom name for the ordered files
- Change the sequence labels to have similar digit lengths by adding zeros in front of small numbers
- Insert clipboard with convenient button
- The CLI can additionally:
  - Specify output location (within the same partition/drive)



## The Workflow
GUI: Highlight + Copy -> Paste + Rename -> Done

CLI: Highlight + Copy -> Run sqre -> Paste + Save + Exit -> Done

## Building
If running Nix with nix-command and flake options enabled: 
	1. Run: `nix develop`
	2. Run: `go build .`
	
If not using Nix, the just run `go build` and hope you have the dependencies installed.

## Usage
GUI: Run `./sqre` and a persistant window will open. Highlight several files from your GUI file manager, and paste them into the text box. Then press the rename button.

TUI: Running `./sqre -cli` will open an emacs window by default. Highlight several files from your GUI file manager, and paste them into emacs. Then save and exit emacs.

Change the default text editor using the -e flag. (Ex: `./sqre -cli -e gedit`)

Find all of the other flags by running: `./sqre -h`

