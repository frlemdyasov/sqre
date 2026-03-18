# Sequence Renamer
Version 1.1 

Inspired by [renameutils](https://www.nongnu.org/renameutils/), sqre is a bulk renaming tool that prioritizes speed and convenience.

> [!NOTE]
> This program was designed for my own personal use case, so alternate functionality may not be added. 

## The Workflow
In your file manager, highlight and copy a series of files -> Run sqre -> Paste the file paths into the text editor, save and exit -> Renamed files appear in their original directory.

This distills down to:  
Highlight + Copy -> Run sqre -> Paste + Save + Exit -> Done

## Building
1. Download the repository as a .zip file
2. Extract the file using, `7z x sqre-main`, for example
3. Move into the working directory: `cd sqre-main`
4. Run `go build`

## Usage
Running `./sqre` will open an emacs window. Highlight several files from your GUI file manager, and paste them into emacs. Then save and exit emacs.

Change the default text editor using the -e flag. (Ex: `./sqre -n gedit`)

Find all of the other flags by running: `./sqre -h`

## Requirements
- Go
- Emacs
