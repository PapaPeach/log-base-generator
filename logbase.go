package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const ProgramDir string = "custom"

var userDone bool = false
var curDir string
var prefix string
var autoEdit bool
var reader = bufio.NewReader(os.Stdin)
var customizations []customization
var customizationsCount int = 0

type customization struct {
	srcFile           string
	panelTree         []string
	numParam          int
	paramLines        []int
	options           []([]string)
	customizationName string
}

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
	customizations[customizationsCount].srcFile, _ = reader.ReadString('\n')                                     // Read to newline
	customizations[customizationsCount].srcFile = strings.TrimSpace(customizations[customizationsCount].srcFile) // Remove newline

	// Convert / to \ for program
	customizations[customizationsCount].srcFile = strings.ReplaceAll(customizations[customizationsCount].srcFile, "/", "\\")

	// Enforce that file ends in .res
	if !strings.HasSuffix(customizations[customizationsCount].srcFile, ".res") {
		customizations[customizationsCount].srcFile += ".res"
	}

	// Validate source file
	_, err := os.Stat(customizations[customizationsCount].srcFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The filepath", customizations[customizationsCount].srcFile, "does not exist")
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

	// Loop until user is done adding customizations
	for !userDone {
		// Create new customization
		customizations = append(customizations, customization{})
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
		customizations[customizationsCount].panelTree = getPanel()
		if customizations[customizationsCount].panelTree == nil {
			fmt.Println("panelTree returned as nil")
			os.Exit(1)
		}
		fmt.Println()

		// Handle tied parameters
		customizations[customizationsCount].numParam = 1
		moreParam := true
		var numValues int
		var customizationTree []string
		for moreParam {
			// Get parameter(s)
			getParamPassed := false
			for !getParamPassed {
				if customizations[customizationsCount].numParam == 1 { // Only one parameter
					customizationTree, getParamPassed = getParam(customizations[customizationsCount].panelTree)
				} else { // Additional parameters
					var tempTree []string
					tempTree, getParamPassed = getParam(customizationTree)
					if getParamPassed == true {
						for i := 0; i < len(customizations[customizationsCount].options); i++ {
							customizations[customizationsCount].options[i] = append(customizations[customizationsCount].options[i], tempTree[len(tempTree)-1])
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
				if customizations[customizationsCount].numParam == 1 { // Only one parameter
					newCustomization := make([]string, len(customizationTree))
					copy(newCustomization, customizationTree)
					customizations[customizationsCount].options = append(customizations[customizationsCount].options, getValues(newCustomization, i, numValues))
				} else { // Additional parameters
					customizations[customizationsCount].options[i] = getValues(customizations[customizationsCount].options[i], i, numValues)
				}
			}

			// Ask for more parameters
			var response string
			for len(response) == 0 {
				fmt.Print("Add more parameters to this customization? [Y] / [N]: ")
				// Use buffered reader because Scanln sucks
				response, _ = reader.ReadString('\n')  // Read to newline
				response = strings.TrimSpace(response) // Remove newline
				if strings.EqualFold(response, "y") {
					customizations[customizationsCount].numParam++
					fmt.Println()
				} else if strings.EqualFold(response, "n") {
					moreParam = false
				} else {
					response = ""
				}
			}
		}

		// Get customization name
		for len(customizations[customizationsCount].customizationName) < 3 {
			fmt.Print("Enter name for customization. This will be used generate log-base aliases: ")
			// Use buffered reader because Scanln sucks
			customizations[customizationsCount].customizationName, _ = reader.ReadString('\n')                                               // Read to newline
			customizations[customizationsCount].customizationName = strings.TrimSpace(customizations[customizationsCount].customizationName) // Remove newline
			// Validate
			if strings.Contains(customizations[customizationsCount].customizationName, " ") {
				fmt.Println("Name cannot contain spaces. Please try again.")
				fmt.Println()
				customizations[customizationsCount].customizationName = ""
			} else if len(customizations[customizationsCount].customizationName) < 3 {
				fmt.Println("Name must be more than 3 characters. Please try again.")
				fmt.Println()
			} else { // Passed
				customizations[customizationsCount].customizationName = strings.ToLower(customizations[customizationsCount].customizationName)
			}
		}
		fmt.Println()

		// Set default parameter value

		// Ask if user has more panels they'd like to edit in srcFile. Y = loop, N = break

		// Ask if user has more files they'd like to generate in directory
		validResponse := false
		for !validResponse {
			var response string
			fmt.Print("Do you have more customizations to generate log-bases for? [Y] / [N]: ")
			// Use buffered reader because Scanln sucks
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if strings.EqualFold(response, "y") {  // Generate more log-base customizations
				validResponse = true
				customizationsCount++
				fmt.Println()
			} else if strings.EqualFold(response, "n") { // Proceed to rest of generation
				userDone = true
				validResponse = true
			} else {
				fmt.Printf("\"%v\" is not a valid response.\n", response)
			}
		}
	}
	fmt.Println()

	// Search for existing hud/cfg/_.cfg file containing "sixense_clear_bindings;sixense_write_bindings _.txt" || "alias lb_log_open"
	// If exists: prompt user to generate aliases there. Else ask user if they'd like to use hud name as file prefix, or custom prefix
	prefix = ""
	for len(prefix) == 0 {
		// Use autodetected hud name as prefix
		var response string
		fmt.Printf("HUD name \"%v\" found, would you like to use that as your config prefix? [Y] / [N]: ", filepath.Base(curDir))
		// Use buffered reader because Scanln sucks
		response, _ = reader.ReadString('\n')  // Read to newline
		response = strings.TrimSpace(response) // Remove newline
		if strings.EqualFold(response, "y") {
			prefix = filepath.Base(curDir)
		} else { // Use custom prefix
			fmt.Print("Enter custom prefix: ")
			fmt.Scanln(&response)
			prefix = response
		}

		// Confirm selection
		fmt.Printf("Config files will generate using prefix: %v\nExample: hud/cfg/%v_generate.cfg\n", prefix, prefix)
		fmt.Print("Continue? [Y] / [N]: ")
		// Use buffered reader because Scanln sucks
		response, _ = reader.ReadString('\n')  // Read to newline
		response = strings.TrimSpace(response) // Remove newline
		if !strings.EqualFold(response, "y") {
			prefix = ""
		}
	}
	fmt.Println()

	// Confirm user would like to comment out srcFile panel lines. Comment with //lb
	editsConfirmed := false
	for !editsConfirmed {
		var response string
		fmt.Print("Allow program to edit necessary lines from source file? [Y] / [N]: ")
		// Use buffered reader because Scanln sucks
		response, _ = reader.ReadString('\n')  // Read to newline
		response = strings.TrimSpace(response) // Remove newline
		if strings.EqualFold(response, "y") {  // Allow automatic comments
			autoEdit = true
			editsConfirmed = true
		} else if strings.EqualFold(response, "n") { // Don't allow automatic comments
			fmt.Println("You will need to manually handle removing customized source lines and adding #base paths or customizations will not work.")
			fmt.Print("Confirm? [Y] / [N]: ")
			// Use buffered reader because Scanln sucks
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if strings.EqualFold(response, "y") {  // Allow automatic comments
				autoEdit = false
				editsConfirmed = true
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
	if autoEdit == true {
		editSource()
	}
}
