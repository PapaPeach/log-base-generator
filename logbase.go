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
var paramLines []int
var customizations []([]string)
var prefix string
var autoComment bool

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

	// Clear input buffer
	var clear string
	fmt.Scanln(&clear)

	// Get parameter(s)
	getParamPassed := false
	for !getParamPassed {
		panelTree, paramLines, getParamPassed = getParam(panelTree)
		if !getParamPassed {
			fmt.Println("Please try again.")
		}
	}
	fmt.Println(panelTree, paramLines)

	// Get value(s)
	valuesDone := false
	for !valuesDone {
		newCustomization := make([]string, len(panelTree))
		copy(newCustomization, panelTree)
		customizations = append(customizations, getValue(newCustomization))
		if !valuesDone {
			var response string
			fmt.Print("Add additional value? [Y] / [N]: ")
			fmt.Scanln(&response)
			if !strings.EqualFold(response, "y") {
				valuesDone = true
			}
		}
	}
	for i := range customizations {
		fmt.Println(i, ": ", customizations[i])
	}

	// Set default value
	
	// Some way of handling "tied" customization options. IE - multiple param and values in a single option

	// Ask if user has more panels they'd like to edit in srcFile. Y = loop, N = break

	// Ask if user has more files they'd like to generate in directory

	// Get Hud name
	hud = filepath.Base(curDir)

	// Search for existing hud/cfg/_.cfg file containing "sixense_clear_bindings;sixense_write_bindings _.txt" || "alias lb_log_open"
	// If exists: prompt user to generate aliases there. Else ask user if they'd like to use hud name as file prefix, or custom prefix
	prefix = ""
	for len(prefix) == 0 {
		// Use autodetected hud name as prefix
		var response string
		fmt.Printf("HUD name \"%v\" found, would you like to use that as your config prefix? [Y] / [N]: ", hud)
		fmt.Scanln(&response)
		if strings.EqualFold(response, "y") {
			prefix = hud
		} else { // Use custom prefix
			fmt.Print("Enter custom prefix: ")
			fmt.Scanln(&response)
			prefix = response
		}

		// Confirm selection
		fmt.Printf("Config files will generate using prefix: %v\nExample: hud/cfg/%v_generate.cfg\n", prefix, prefix)
		fmt.Print("Continue? [Y] / [N]: ")
		fmt.Scanln(&response)
		if !strings.EqualFold(response, "y") {
			prefix = ""
		}
	}

	// Confirm user would like to comment out srcFile panel lines. Comment with //lb
	commentsConfirmed := false
	for !commentsConfirmed {
		var response string
		fmt.Print("Allow program to comment necessary lines from source file? [Y] / [N]: ")
		fmt.Scanln(&response)
		if strings.EqualFold(response, "y") { // Allow automatic comments
			autoComment = true
			commentsConfirmed = true
		} else if strings.EqualFold(response, "n") { // Don't allow automatic comments
			fmt.Println("You will need to manually comment / remove lines with customized parameters or customizations will not work.")
			fmt.Print("Confrim? [Y] / [N]: ")
			fmt.Scanln(&response)
			if strings.EqualFold(response, "y") { // Allow automatic comments
				autoComment = false
				commentsConfirmed = true
			}
		}
	}

	// Populate hud/cfg/_.cfg file
	generateConfig()

	// Comment srcFile panel lines
	if autoComment == true {

	}
}
