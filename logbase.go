package main

import (
	"bufio"
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
var numParam int
var paramLines []int
var customizations []([]string)
var customizationName string
var prefix string
var autoComment bool
var reader = bufio.NewReader(os.Stdin)

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
	// Use buffered reader because Scanln sucks
	srcFile, _ = reader.ReadString('\n') // Read to newline
	srcFile = strings.TrimSpace(srcFile) // Remove newline

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
		fmt.Println("Program is in the wrong location. Please place it in your rood Hud folder, where info.vdf is located.")
		os.Exit(1)
	}
	fmt.Println("Program passed location check.")
	fmt.Println()

	// Validate source file
	getSourceFilePassed := false
	for !getSourceFilePassed {
		getSourceFilePassed = getSourceFile()
		if !getSourceFilePassed {
			fmt.Println("Please try again.")
			fmt.Println()
		}
	}

	// Get panel tree
	panelTree = getPanel()
	if panelTree == nil {
		return
	}
	fmt.Println()

	// Handle tied parameters
	numParam = 1
	moreParam := true
	var numValues int
	var customizationTree []string
	for moreParam {
		// Get parameter(s)
		getParamPassed := false
		for !getParamPassed {
			if numParam == 1 { // Only one parameter
				customizationTree, getParamPassed = getParam(panelTree)
			} else { // Additional parameters
				var tempTree []string
				tempTree, getParamPassed = getParam(customizationTree)
				if getParamPassed == true {
					for i := 0; i < len(customizations); i++ {
						customizations[i] = append(customizations[i], tempTree[len(tempTree)-1])
					}
				}
			}
			if !getParamPassed { // Parameter not found
				fmt.Println("Please try again.")
				fmt.Println()
			}
		}

		// Get total number of values
		for numValues < 2 {
			fmt.Print("Enter quantity of values to set: ")
			fmt.Scanln(&numValues)
		}
		// Get values
		for i := 0; i < numValues; i++ {
			if numParam == 1 { // Only one parameter
				newCustomization := make([]string, len(customizationTree))
				copy(newCustomization, customizationTree)
				customizations = append(customizations, getValues(newCustomization, i, numValues))
				/*if !valuesDone {
					var response string
					fmt.Print("Add additional value? [Y] / [N]: ")
					fmt.Scanln(&response)
					if !strings.EqualFold(response, "y") {
						valuesDone = true
					}
				}*/
			} else { // Additional parameters
				customizations[i] = getValues(customizations[i], i, numValues)
			}
		}
		// For debug output
		/*for i := range customizations {
			fmt.Println(i+1, ":", customizations[i])
		}*/

		// Ask for more parameters
		var response string
		for len(response) == 0 {
			fmt.Print("Add more parameters to this customization? [Y] / [N]: ")
			// Use buffered reader because Scanln sucks
			response, _ = reader.ReadString('\n') // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if strings.EqualFold(response, "y") {
				numParam++
			} else if strings.EqualFold(response, "n") {
				moreParam = false
			} else {
				response = ""
			}
		}
		fmt.Println()
	}
	
	// Get customization name
	for len(customizationName) < 3 {
		fmt.Print("Enter name for customization. This will be used generate log-base aliases: ")
		// Use buffered reader because Scanln sucks
		customizationName, _ = reader.ReadString('\n') // Read to newline
		customizationName = strings.TrimSpace(customizationName) // Remove newline
		// Validate
		if strings.Contains(customizationName, " ") {
			fmt.Println("Name cannot contain spaces. Please try again.")
			fmt.Println()
			customizationName = ""
		} else if len(customizationName) < 3 {
			fmt.Println("Name must be more than 3 characters. Please try again.")
			fmt.Println()
		} else { // Passed
			customizationName = strings.ToLower(customizationName)
		}
	}
	fmt.Println()

	// Set default parameter value

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
		// Use buffered reader because Scanln sucks
		response, _ = reader.ReadString('\n') // Read to newline
		response = strings.TrimSpace(response) // Remove newline
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
		// Use buffered reader because Scanln sucks
		response, _ = reader.ReadString('\n') // Read to newline
		response = strings.TrimSpace(response) // Remove newline
		if !strings.EqualFold(response, "y") {
			prefix = ""
		}
	}
	fmt.Println()

	// Confirm user would like to comment out srcFile panel lines. Comment with //lb
	commentsConfirmed := false
	for !commentsConfirmed {
		var response string
		fmt.Print("Allow program to comment necessary lines from source file? [Y] / [N]: ")
		// Use buffered reader because Scanln sucks
		response, _ = reader.ReadString('\n') // Read to newline
		response = strings.TrimSpace(response) // Remove newline
		if strings.EqualFold(response, "y") { // Allow automatic comments
			autoComment = true
			commentsConfirmed = true
		} else if strings.EqualFold(response, "n") { // Don't allow automatic comments
			fmt.Println("You will need to manually comment / remove lines with customized parameters or customizations will not work.")
			fmt.Print("Confirm? [Y] / [N]: ")
			// Use buffered reader because Scanln sucks
			response, _ = reader.ReadString('\n') // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if strings.EqualFold(response, "y") { // Allow automatic comments
				autoComment = false
				commentsConfirmed = true
			}
		}
	}
	fmt.Println()

	// Populate hud/cfg/_.cfg file
	generateMainConfig()
	generateSaveConfig()
	generateGeneratorConfig()
	generateValveRc()
	
	// Generate button code pastable
	generateButtonCommands()

	// Comment srcFile panel lines
	if autoComment == true {
		commentSource()
	}
}
