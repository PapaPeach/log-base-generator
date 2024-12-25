package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	siblings          []int
}

func locationCheck() bool {
	// Check if program is in the right directory by checking if info.vdf exists
	if _, err := os.Stat("info.vdf"); err == nil { // Handle file already exists
		if curDir, err = os.Getwd(); err != nil {
			fmt.Print("Error getting working directory:", err)
			return false
		}
		return true
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fmt.Println("Log-Base Generator is in the wrong location.\nPlease put it in your HUD's root directory, where info.vdf is located.")
		return false
	} else {
		fmt.Println("Error getting program location:", err)
	}
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

func getResponse(prompt string, failText string, options []map[string]string) string {
	// Get slice containing all options available for selection
	validResponse := false
	hasKeyText := false

	fmt.Print(prompt)

	// Loop until we get a response that matches our expectation
	for !validResponse {
		var response string

		// Check if option has an explaination assiociated with it, if so display it and ask for selection
		for o := range options {
			for key, keyText := range options[o] {
				if len(keyText) > 0 {
					fmt.Printf("\n[%v] %v", key, keyText)
					hasKeyText = true
				}
			}
		}
		if hasKeyText {
			fmt.Print("\nPlease select an option: ")
		}

		// Use buffered reader for getting input because Scanln sucks
		response, _ = reader.ReadString('\n')  // Read to newline
		response = strings.TrimSpace(response) // Remove newline

		// Compare user input to defined keys and only return if it matches
		for o := range options {
			for key := range options[o] {
				if strings.EqualFold(response, key) {
					validResponse = true
					return key
				}
			}
		}
		fmt.Printf("\"%v\" %v", response, failText)
	}
	return ""
}

func getSibling() {
	// Check if panel should be a sibling, only ask if customizationsCount >= 1
	if customizationsCount > 0 {
		prompt := "Should this be used in conjuntion with a previous panel customization? [Y] / [N]: "
		options := []map[string]string{
			{"y": ""},
			{"n": ""},
		}
		failText := "is an invalid response. [Y] / [N]: "
		response := getResponse(prompt, failText, options)
		if response == "y" {
			var siblingIndex int
			if customizationsCount == 1 { // If there is only one possible sibling
				siblingIndex = 0
			} else { // If the user needs to select from several possible siblings
				prompt := fmt.Sprintf("Found %v possible siblings:", customizationsCount)
				options := []map[string]string{}
				// Incrementally create options of previous customization names
				for i := range customizationsCount {
					options = append(options, map[string]string{strconv.Itoa(i + 1): customizations[i].customizationName})
				}
				failText := fmt.Sprintf("is an invalid response. Please make a selection 1 - %v: ", customizationsCount)
				response := getResponse(prompt, failText, options)

				// Assign to sibling accordingly
				var err error
				siblingIndex, err = strconv.Atoi(response)
				if err != nil {
					fmt.Println("Error converting selection to int:", err)
					os.Exit(1)
				}
				siblingIndex--
				// Assign sibling to eldest sibling by detecting if selected customization has an older sibling
				if len(customizations[siblingIndex].siblings) > 0 && customizations[siblingIndex].siblings[0] < siblingIndex {
					siblingIndex = customizations[siblingIndex].siblings[0]
				}
			}
			// Assign sibling to eldest
			customizations[siblingIndex].siblings = append(customizations[siblingIndex].siblings, customizationsCount)
			// Indicate that current customization is a younger sibling
			customizations[customizationsCount].siblings = append(customizations[customizationsCount].siblings, siblingIndex)
		}
		fmt.Println()
	}
}

func main() {
	// Validate program location
	if !locationCheck() {
		os.Exit(1)
	}
	fmt.Println("Program passed location check.")
	fmt.Println()

	// Loop until user is done adding customizations
	for !userDone {
		// Create new customization
		customizations = append(customizations, customization{})

		// Check if panel should be sibling, and match numParam
		getSibling()

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
		validPanel := false
		for !validPanel {
			customizations[customizationsCount].panelTree = getPanel()
			if customizations[customizationsCount].panelTree != nil { // If panel is invalid
				validPanel = true
			}
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
					if getParamPassed {
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
			if customizations[customizationsCount].siblings == nil { // Not a sibling
				for numValues < 2 {
					fmt.Print("Enter quantity of values to set: ")
					fmt.Scanln(&numValues)
					if numValues < 2 {
						fmt.Println("Must set at least 2 values.")
					} else if numValues > 9 { // Check for abnormally high numValues
						prompt := fmt.Sprintf("Detected high quantity of values, this is not necessarily an error.\nConfirm generate %v values? [Y] / [N]: ", numValues)
						options := []map[string]string{
							{"y": ""},
							{"n": ""},
						}
						failText := "is an invalid response. [Y] / [N]: "
						response := getResponse(prompt, failText, options)
						if response == "n" {
							numValues = 0
						}
					}
				}
			} else { // If sibling, get quantity from eldest sibling
				numValues = len(customizations[customizations[customizationsCount].siblings[0]].options)
				fmt.Printf("Matching quantity of values to %v\n", customizations[customizations[customizationsCount].siblings[0]].customizationName)
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
			fmt.Println()
			prompt := "Add more parameters to this customization? [Y] / [N]: "
			options := []map[string]string{
				{"y": ""},
				{"n": ""},
			}
			failText := "is an invalid response. [Y] / [N]: "
			response := getResponse(prompt, failText, options)
			if response == "y" {
				customizations[customizationsCount].numParam++
			} else if response == "n" {
				moreParam = false
			}
		}

		// Get customization name
		for len(customizations[customizationsCount].customizationName) < 3 {
			fmt.Print("Enter name for customization. This will be used generate log-base aliases: ")
			// Use buffered reader because Scanln sucks
			customizations[customizationsCount].customizationName, _ = reader.ReadString('\n')                                               // Read to newline
			customizations[customizationsCount].customizationName = strings.TrimSpace(customizations[customizationsCount].customizationName) // Remove newline
			// Validate
			if len(customizations[customizationsCount].customizationName) < 3 {
				fmt.Printf("\"%v\" is an invalid name. Name must be longer than 3 characters.\n", customizations[customizationsCount].customizationName)
			} else if strings.Contains(customizations[customizationsCount].customizationName, " ") {
				fmt.Printf("\"%v\" is an invalid name. Name cannot contain any spaces.\n", customizations[customizationsCount].customizationName)
				customizations[customizationsCount].customizationName = ""
			} else { // Passed easy checks, now check for duplicate names
				customizations[customizationsCount].customizationName = strings.ToLower(customizations[customizationsCount].customizationName)
				for i := range customizationsCount {
					if customizations[customizationsCount].customizationName == customizations[i].customizationName {
						fmt.Printf("\"%v\" is an invalid name. It matches a previously assigned name.\n", customizations[customizationsCount].customizationName)
						customizations[customizationsCount].customizationName = ""
					}
				}
			}
		}

		// TODO: Set default parameter value

		// TODO: Ask if user has more panels they'd like to edit in srcFile. Y = loop, N = break

		// Ask if user has more files they'd like to generate in directory
		fmt.Println()
		prompt := "Do you have more customizations to generate log-bases for? [Y] / [N]: "
		options := []map[string]string{
			{"y": ""},
			{"n": ""},
		}
		failText := "is an invalid response. [Y] / [N]: "
		response := getResponse(prompt, failText, options)
		if response == "y" { // Generate more log-base customizations
			customizationsCount++
		} else if response == "n" { // Proceed to rest of generation
			userDone = true
		}
	}
	fmt.Println()

	// TODO: Search for existing hud/cfg/_.cfg file containing "sixense_clear_bindings;sixense_write_bindings _.txt" || "alias lb_log_open"
	// If exists: prompt user to generate aliases there. Else ask user if they'd like to use hud name as file prefix, or custom prefix

	// Get config prefix
	prefix = ""
	for len(prefix) == 0 {
		// Use autodetected hud name as prefix
		prefixPrompt := fmt.Sprintf("HUD name \"%v\" found, would you like to use that as your config prefix? [Y] / [N]: ", filepath.Base(curDir))
		prefixOptions := []map[string]string{
			{"y": ""},
			{"n": ""},
		}
		prefixFailText := "is an invalid response. [Y] / [N]: "
		response := getResponse(prefixPrompt, prefixFailText, prefixOptions)
		if response == "y" { // Use HUD name as prefix
			prefix = filepath.Base(curDir)
		} else if response == "n" { // Use custom prefix if valid
			validPrefix := false
			for !validPrefix {
				fmt.Print("Enter custom prefix: ")
				// Use buffered reader because Scanln sucks
				response, _ = reader.ReadString('\n')  // Read to newline
				response = strings.TrimSpace(response) // Remove newline
				if len(response) < 3 {
					fmt.Printf("\"%v\" is an invalid prefix. Prefix must be longer than 3 characters.\n", response)
					prefix = ""
				} else if strings.Contains(response, " ") {
					fmt.Printf("\"%v\" is an invalid prefix. Prefix cannot contain any spaces.\n", response)
					prefix = ""
				} else {
					validPrefix = true
					prefix = response
				}
			}
		}

		// Confirm selection
		confirmPrompt := fmt.Sprintf("Config files will generate using prefix: %v\nExample: hud/cfg/%v_generate.cfg\nContinue? [Y] / [N]: ", prefix, prefix)
		confirmOptions := []map[string]string{
			{"y": ""},
			{"n": ""},
		}
		confirmFailText := "is an invalid response. [Y] / [N]: "
		response = getResponse(confirmPrompt, confirmFailText, confirmOptions)
		if response == "n" { // Reset prefix and restart prefix selection prompts
			prefix = ""
		}
	}
	fmt.Println()

	// Confirm user would like to comment out srcFile panel lines. Comment with //lb
	editsConfirmed := false
	for !editsConfirmed {
		prompt := "Allow program to edit necessary lines from source file? [Y] / [N]: "
		options := []map[string]string{
			{"y": ""},
			{"n": ""},
		}
		failText := "is an invalid response. [Y] / [N]: "
		response := getResponse(prompt, failText, options)
		if response == "y" { // Allow automated editing
			autoEdit = true
			editsConfirmed = true
		} else if response == "n" { // Don't allow automated editing
			prompt := ("You will need to manually remove customized source lines and add #base paths for customizations to work.\nConfirm? [Y] / [N]: ")
			options := []map[string]string{
				{"y": ""},
				{"n": ""},
			}
			failText := "is an invalid response. [Y] / [N]: "
			response := getResponse(prompt, failText, options)
			if response == "y" { // Confirm user will manually handle it
				autoEdit = false
				editsConfirmed = true
			}
		}
	}
	fmt.Println()

	// Create cfg directory if it doesn't already exist
	checkCfg()

	// Populate hud/cfg/_.cfg file
	generateMainConfig()
	generateSaveConfig()
	generateGeneratorConfig()
	generateValveRc()

	// Generate button code pastable
	generateButtonCommands()

	// Comment srcFile panel lines
	if autoEdit {
		editSource()
	}

	// Announce program is done
	fmt.Println("Program has finished :D")
	fmt.Println("Press Enter to exit.")
	anyKey, _ := reader.ReadString('\n')
	if anyKey != "" {
		os.Exit(0)
	}
}
