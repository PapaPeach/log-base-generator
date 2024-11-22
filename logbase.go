package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const ProgramDir string = "custom"

var curDir string
var srcFile string
var hud string
var panelTree []string

func locationCheck() bool {
	// Check if program is in the right directory
	var err error
	curDir, err = os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return false
	}

	parentDir := filepath.Dir(curDir)

	if filepath.Base(parentDir) == ProgramDir {
		return true
	}

	fmt.Println("Log-Base Generator is in the wrong location.\nPlease put it in tf/custom/yourHud")
	return false
}

func getSourceFile() bool {
	// Get source file
	fmt.Print("Enter filepath of file: ")
	fmt.Scanln(&srcFile)

	// Convert / to \ for program
	srcFile = strings.ReplaceAll(srcFile, "/", "\\")

	// Enforce that file ends in .res
	if !strings.HasSuffix(srcFile, ".res") {
		srcFile += ".res"
	}

	// Validate source file
	_, err := os.Stat(srcFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The filepath", srcFile, "does not exist")
			return false
		} else {
			fmt.Println("Error checking file path:", err)
			return false
		}
	}

	return true
}

func main() {
	// Validate program location
	if !locationCheck() {
		return
	}
	fmt.Println("Program passed location check.")

	// Validate source file
	if !getSourceFile() {
		return
	}

	// Get panel tree
	panelTree = getPanel()
	if panelTree == nil {
		return
	}
	fmt.Println(panelTree)

	// Get parameter(s)

	// Get value

	// Ask if user has more panels they'd like to edit in srcFile. Y = loop, N = break

	// Ask if user has more files they'd like to generate in directory

	// Get Hud name
	hud = filepath.Base(curDir)

	// Search for existing hud/cfg/_.cfg file containing "sixense_clear_bindings;sixense_write_bindings _.txt" || "alias lb_log_open"
	// If exists: prompt user to generate aliases there. Else ask user if they'd like to use hud name as file prefix, or custom prefix

	// Confirm user would like to comment out srcFile panel lines. Comment with //lb

	// Populate hud/cfg/_.cfg file
}
