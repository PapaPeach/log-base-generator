package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getPanel() []string {
	// Get panel
	var panel string
	fmt.Print("Enter name of panel: ")
	fmt.Scanln(&panel)

	// Open source file
	file, err := os.Open(customizations[customizationsCount].srcFile)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Scan file for matches
	count := 0
	level := 0
	var prevWord string
	var panels []string
	var levels []int
	var panelTree []string
	scnr := bufio.NewScanner(file)

	// Create slice of all panelsFind panel
	for scnr.Scan() {
		line := scnr.Text()
		words := strings.Fields(line)

		for _, word := range words {
			// Go to next line if commented
			if strings.HasPrefix(word, "//") {
				break
			} else if word == `{` || strings.HasPrefix(word, "{") { // Count nested level
				panels = append(panels, prevWord)
				levels = append(levels, level)
				level++
			} else if word == `}` || strings.HasSuffix(word, "}") { // Manage lower nested level
				level--
			}

			prevWord = word
		}
	}

	// Search for specified panel in list of all panels
	for i := range panels {
		if strings.EqualFold(panels[i], panel) || strings.EqualFold(panels[i], "\""+panel+"\"") {
			panel = panels[i] // Set panel var to name of actual panel for fanciness
			count++
		}
	}

	// Immediately handle no match found
	if count < 1 {
		fmt.Println("Could not find", panel, "in", customizations[customizationsCount].srcFile)
		return nil
	}

	// Handle open and close brace mismatch
	if len(panels) != len(levels) {
		fmt.Println("Number of panels and number of braces do not correlate correctly.")
		return nil
	}

	// Filer all non-parents of desired panel
	var filteredPanels []string
	for i := range panels {
		if panels[i] == panel { // add matches to slice of to slice of possible panels
			filteredPanels = append(filteredPanels, panels[i])

			panelLevel := levels[i]
			for j := i; j > 0; j-- { // Once panel is found, add to family tree, leave off header (level 0)
				if levels[j] < panelLevel { // Determine direct parent by comparing level to previous
					filteredPanels = append(filteredPanels, panels[j])
					panelLevel = levels[j]
				}
			}
		}
	}

	// Remove double quotes around panel title
	for i := range filteredPanels {
		filteredPanels[i] = strings.ReplaceAll(filteredPanels[i], "\"", "")
	}
	panel = filteredPanels[0]

	// Arange tree in hierarchical order
	prevMatch := 0
	var optionIndexes []int
	for i := 1; i < len(filteredPanels); i++ {
		if filteredPanels[i] == panel { // create tree for instances of panel
			for j := i - 1; j > prevMatch; j-- {
				panelTree = append(panelTree, filteredPanels[j])
			}
			panelTree = append(panelTree, filteredPanels[prevMatch])
			optionIndexes = append(optionIndexes, prevMatch, i-1)
			prevMatch = i
		} else if i == len(filteredPanels)-1 { // Special handling for last instance
			for j := i; j > prevMatch; j-- {
				panelTree = append(panelTree, filteredPanels[j])
			}
			panelTree = append(panelTree, filteredPanels[prevMatch])
			optionIndexes = append(optionIndexes, prevMatch, i)
		}
	}

	// Handle matches
	if count == 1 {
		fmt.Printf("Found %v in %v.\n", panel, customizations[customizationsCount].srcFile)
		// Handle level 1 lone match
		if len(filteredPanels) == 1 {
			return filteredPanels
		}
		return panelTree
	} else if count > 1 { // Duplicates found, user input needed
		fmt.Printf("Found %v instances of %v in %v.\n", count, panel, customizations[customizationsCount].srcFile)
		// Print options
		optionNum := 1
		for i := 0; i < len(optionIndexes)-1; i += 2 {
			fmt.Printf("[%v] ", optionNum)
			for j := optionIndexes[i]; j <= optionIndexes[i+1]; j++ {
				if j != optionIndexes[i+1] {
					fmt.Printf("%v > ", panelTree[j])
				} else {
					fmt.Println(panelTree[j])
				}
			}
			optionNum++
		}

		// Recieve user selection
		fmt.Print("Please select an option: ")
		var option int
		for option < 1 || option > count {
			// Use buffered reader because Scanln sucks
			input, _ := reader.ReadString('\n') // Read to newline
			input = strings.TrimSpace(input)    // Remove newline
			option, err = strconv.Atoi(input)

			// Validate
			if err != nil || option < 1 || option > count {
				fmt.Printf("Please make a selection 1 - %v: ", count)
				option = 0
			}
		}

		// Translate user selection to output slice
		option = (option - 1) * 2
		return panelTree[optionIndexes[option] : optionIndexes[option+1]+1]
	}

	return nil
}

func getParam(tree []string) ([]string, bool) {

	// Get parameter
	var param string
	fmt.Print("Enter parameter to customize: ")
	// Use buffered reader because Scanln sucks
	param, _ = reader.ReadString('\n') // Read to newline
	param = strings.TrimSpace(param)   // Remove newline

	// Check for previous instance of parameter
	if customizations[customizationsCount].numParam > 1 {
		// Increment through just the parameters in the options slice by skipping the panel nav part of the tree
		for i := len(customizations[customizationsCount].options[0]) - (2 * (customizations[customizationsCount].numParam - 1)); i < len(customizations[customizationsCount].options[0]); i += 2 {
			if strings.EqualFold(param, customizations[customizationsCount].options[0][i]) {
				fmt.Printf("Previous instance of: \"%v\" found.\n", param)
				return tree, false
			}
		}
	}

	// Open source file
	file, err := os.Open(customizations[customizationsCount].srcFile)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Scan to panel
	var isParent bool
	lineNum := 1
	level := 0
	scnr := bufio.NewScanner(file)

	for scnr.Scan() {
		line := scnr.Text()

		if strings.Contains(line, "{") { // Count nested level
			level++
		}
		if level > 0 && level <= len(customizations[customizationsCount].panelTree) && strings.Contains(line, customizations[customizationsCount].panelTree[level-1]) { // Navigate through tree
			isParent = true
		} else if level == len(customizations[customizationsCount].panelTree)+1 && isParent && strings.Contains(line, "\""+param+"\"") { // Parameter found in correct panel
			tree = append(tree, param)
			customizations[customizationsCount].paramLines = append(customizations[customizationsCount].paramLines, lineNum)
			return tree, true
		}
		if strings.Contains(line, "}") { // Track tree status and nested level
			isParent = false
			level--
		}

		// Track line number
		lineNum++
	}

	// No match
	fmt.Println("Did not find parameter:", param)
	return tree, false
}

func getValues(tree []string, curNum int, numValues int) []string {
	// Get values
	var value string

	for len(value) < 1 {
		fmt.Printf("Enter value for %v (%v/%v): ", tree[len(tree)-1], curNum+1, numValues)
		// Use buffered reader because Scanln sucks
		value, _ = reader.ReadString('\n') // Read to newline
		value = strings.TrimSpace(value)   // Remove newline

		// Check that value does not contain space
		if len(value) < 1 {
			fmt.Println("Invalid value. Values must contain something.")
		} else if strings.Contains(value, " ") {
			fmt.Printf("\"%v\" is an invalid value. Values cannot contain any spaces.\n", value)
			fmt.Println("If this is a color or font value, create an appropriate definition in your ClientScheme.")
			value = ""
		}
	}

	// Add value to tree
	tree = append(tree, value)

	return tree
}
