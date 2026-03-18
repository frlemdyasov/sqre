# Sequence Renamer
Version 1.1 

Inspired by [renameutils](https://www.nongnu.org/renameutils/), sqre is a bulk renaming tool that prioritizes speed and convenience.

## The Workflow
In your file manager, highlight and copy a series of files -> Run sqre -> Paste the file paths into the text editor, save and exit -> Renamed files appear in their original directory.

This distills down to:  
Highlight + Copy -> ./sqre -> Paste + Save -> Done

## Building
1. Download the repository as a .zip file
2. Extract the file using, `7z x sqre-main`, for example
3. Move the working directory: `cd sqre-main`
4. run `go build`

## Requirements
Go 1.25.7 worked for me
