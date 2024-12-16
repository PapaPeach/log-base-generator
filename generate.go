package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func copyFile(file string) []string {
	// Open file
	inputFile, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening %v for reading: %v\n", file, err)
		os.Exit(1)
	}

	defer inputFile.Close()

	// Create slice containing all the lines
	var fileContents []string
	scnr := bufio.NewScanner(inputFile)
	for scnr.Scan() {
		line := scnr.Text()
		fileContents = append(fileContents, line)
	}
	return fileContents
}

func pasteFile(file *os.File, fileContents []string) {
	// Repopulate file
	for i := range fileContents {
		file.WriteString(fileContents[i])
		if i < len(fileContents)-1 {
			file.WriteString("\n")
		}
	}
}

func getUniqueFilePrefix(curPrefix string) string {
	// Increment number suffix until fileName is unique
	suffix := 1
	filePrefix := curPrefix + strconv.Itoa(suffix)
	var uniqueFileName bool

	for !uniqueFileName {
		if _, err := os.Stat("cfg/" + filePrefix + ".cfg"); err == nil { // File exists
			filePrefix = curPrefix + strconv.Itoa(suffix+1)
		} else if _, err := os.Stat("cfg/" + filePrefix + ".rc"); err == nil { // valveX.rc exists
			filePrefix = curPrefix + strconv.Itoa(suffix+1)
		} else if _, err := os.Stat(filePrefix + ".txt"); err == nil { // Button commands file exists
			filePrefix = curPrefix + strconv.Itoa(suffix+1)
		} else {
			uniqueFileName = true
		}
	}
	return filePrefix
}

func generateMainConfig() {
	// Open or create file
	fileName := "cfg/" + prefix + ".cfg"
	var fileExists bool
	var appendToTop bool
	var fileContents []string

	if _, err := os.Stat(fileName); err == nil { // Handle file already exists
		fileExists = true

		// Prompt user how to deal with existing file
		var validResponse bool
		fmt.Printf("%v already exists, how would you like to resolve this?\n", fileName)
		for !validResponse {
			fmt.Println("[1] Generate code above the original contents of existing file")
			fmt.Println("[2] Generate code to unique file to be copy + pasted manually")
			fmt.Println("[3] Quit program without generating or modifying the file")
			fmt.Print("Please select and option: ")

			// Use buffered reader because Scanln sucks
			var response string
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if response == "1" {
				appendToTop = true
				validResponse = true
				fileContents = copyFile(fileName)
			} else if response == "2" {
				appendToTop = false
				validResponse = true
				fileName = "cfg/" + getUniqueFilePrefix(prefix) + ".cfg"
			} else if response == "3" {
				fmt.Println("\nProgram exiting, the file was not generated or modified.")
				os.Exit(0)
			} else {
				fmt.Println("Invalid response, how you would like to resolve this?")
				validResponse = false
			}
			fmt.Println()
		}
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("Main config file is inaccessible!")
		os.Exit(1)
	}

	// Create main alias file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating main config file:", err)
		os.Exit(1)
	}

	defer file.Close()

	// Create log_open alias
	file.WriteString("//Open log for writing alias\n")
	file.WriteString("alias lb_log_selection_open \"sixense_clear_bindings;sixense_write_bindings " + prefix + "_customization_selection.txt;con_timestamp 0;con_logfile cfg/" + prefix + "_customization_selection.txt\"\n")
	file.WriteString("alias lb_log_customization_open \"sixense_clear_bindings;sixense_write_bindings " + prefix + "_customizations.txt;con_timestamp 0;con_logfile cfg/" + prefix + "_customizations.txt\"\n")
	file.WriteString("\n")

	// Iterate through customizations and create save aliases
	file.WriteString("//Declare customization save aliases\n")
	for c := range customizations {
		saveAlias := prefix + "_" + customizations[c].customizationName + "_dump"
		file.WriteString("alias " + saveAlias + " \"\"\n")
	}
	file.WriteString("\n")

	// Iterate through customizations and Create default value aliases
	file.WriteString("//Initialize default values\n")
	for c := range customizations {
		var defaultAlias string
		defaultBraceCount := 0
		for i := 0; i < len(customizations[c].options[0])-(2*customizations[c].numParam); i++ {
			defaultAlias += customizations[c].options[0][i] + "{"
			defaultBraceCount++
		}

		// Parameter and value
		defaultAlias += strings.Join(customizations[c].options[0][len(customizations[c].options[0])-(2*customizations[c].numParam):], " ")

		//Close braces
		for j := 0; j < defaultBraceCount; j++ {
			defaultAlias += "}"
		}

		writeAlias := prefix + "_" + customizations[c].customizationName + "_write"
		file.WriteString("alias " + writeAlias + " \"echo " + defaultAlias + "\"\n")
	}
	file.WriteString("\n")

	// Create customization definitions
	file.WriteString("//Define customization aliases\n")
	for c := range customizations {
		for i := range customizations[c].options {
			saveAlias := prefix + "_" + customizations[c].customizationName + "_dump"
			customizationAlias := customizations[c].customizationName + strconv.Itoa(i+1)
			writeAlias := prefix + "_" + customizations[c].customizationName + "_write"
			var panelCode string
			braceCount := 0

			// Panel path
			for j := 0; j < len(customizations[c].options[i])-(2*customizations[c].numParam); j++ {
				panelCode += customizations[c].options[i][j] + "{"
				braceCount++
			}

			// Parameter and value
			panelCode += strings.Join(customizations[c].options[i][len(customizations[c].options[i])-(2*customizations[c].numParam):], " ")
			//Close braces
			for j := 0; j < braceCount; j++ {
				panelCode += "}"
			}

			// alias customizationAlias "alias saveAlias echo customizationAlias;alias writeAlias echo panelCode"
			file.WriteString("alias " + customizationAlias +
				" \"alias " + saveAlias + " echo " + customizationAlias +
				";alias " + writeAlias + " echo " + panelCode + "\"\n")
		}
	}

	// Main config file already exists and user opted to append log-base to the top of it
	if fileExists == true && appendToTop {
		file.WriteString("\n//User contents\n")
		pasteFile(file, fileContents)
	}
}

func generateSaveConfig() {
	// Open or create file
	fileName := "cfg/" + prefix + "_save.cfg"
	var fileExists bool
	var appendToTop bool
	var fileContents []string

	if _, err := os.Stat(fileName); err == nil { // Handle file already exists
		fileExists = true

		// Prompt user how to deal with existing file
		var validResponse bool
		fmt.Printf("%v already exists, how would you like to resolve this?\n", fileName)
		for !validResponse {
			fmt.Println("[1] Generate code above the original contents of existing file")
			fmt.Println("[2] Generate code to unique file to be copy + pasted manually")
			fmt.Println("[3] Quit program without generating or modifying the file")
			fmt.Print("Please select and option: ")

			// Use buffered reader because Scanln sucks
			var response string
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if response == "1" {
				appendToTop = true
				validResponse = true
				fileContents = copyFile(fileName)
			} else if response == "2" {
				appendToTop = false
				validResponse = true
				fileName = "cfg/" + getUniqueFilePrefix(prefix+"_save") + ".cfg"
			} else if response == "3" {
				fmt.Println("\nProgram exiting, the file was not generated or modified.")
				os.Exit(0)
			} else {
				fmt.Println("Invalid response, how you would like to resolve this?")
				validResponse = false
			}
			fmt.Println()
		}
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("Save file is inaccessible!")
		os.Exit(1)
	}

	// Create save file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating save file:", err)
		os.Exit(1)
	}

	defer file.Close()

	// Create aliases to dump save aliases to file
	file.WriteString("//Clear and prep log file\n")
	file.WriteString("lb_log_selection_open\n")
	file.WriteString("\n")

	file.WriteString("//Dump current aliases to file\n")
	for c := range customizations {
		saveAlias := prefix + "_" + customizations[c].customizationName + "_dump"
		file.WriteString(saveAlias + "\n")
	}
	file.WriteString("\n")

	file.WriteString("//Close log file\n")
	file.WriteString("con_logfile \"\"")

	// Save config file already exists and user opted to append log-base to the top of it
	if fileExists == true && appendToTop {
		file.WriteString("\n\n//User contents\n")
		pasteFile(file, fileContents)
	}
}

func generateGeneratorConfig() {
	// Open or create file
	fileName := "cfg/" + prefix + "_generate.cfg"
	var fileExists bool
	var appendToTop bool
	var fileContents []string

	if _, err := os.Stat(fileName); err == nil { // Handle file already exists
		fileExists = true

		// Prompt user how to deal with existing file
		var validResponse bool
		fmt.Printf("%v already exists, how would you like to resolve this?\n", fileName)
		for !validResponse {
			fmt.Println("[1] Generate code above the original contents of existing file")
			fmt.Println("[2] Generate code to unique file to be copy + pasted manually")
			fmt.Println("[3] Quit program without generating or modifying the file")
			fmt.Print("Please select and option: ")

			// Use buffered reader because Scanln sucks
			var response string
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if response == "1" {
				appendToTop = true
				validResponse = true
				fileContents = copyFile(fileName)
			} else if response == "2" {
				appendToTop = false
				validResponse = true
				fileName = "cfg/" + getUniqueFilePrefix(prefix+"_generate") + ".cfg"
			} else if response == "3" {
				fmt.Println("\nProgram exiting, the file was not generated or modified.")
				os.Exit(0)
			} else {
				fmt.Println("Invalid response, how you would like to resolve this?")
				validResponse = false
			}
			fmt.Println()
		}
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("Generator file is inaccessible!")
		os.Exit(1)
	}

	// Create generate file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating generator file:", err)
		os.Exit(1)
	}

	defer file.Close()

	// Create aliases to dump save aliases to file
	file.WriteString("//Clear and prep log file\n")
	file.WriteString("lb_log_customization_open\n")
	file.WriteString("\n")

	file.WriteString("//Setup file\n")
	file.WriteString("echo \"x{\"\n")
	file.WriteString("\n")

	file.WriteString("//Write current customizations to file\n")
	for c := range customizations {
		writeAlias := prefix + "_" + customizations[c].customizationName + "_write"
		file.WriteString(writeAlias + "\n")
	}
	file.WriteString("\n")

	file.WriteString("//Close log file\n")
	file.WriteString("echo \"}\"\n")
	file.WriteString("con_logfile \"\"")

	// Generator config file already exists and user opted to append log-base to the top of it
	if fileExists == true && appendToTop {
		file.WriteString("\n\n//User contents\n")
		pasteFile(file, fileContents)
	}
}

func generateValveRc() {
	// Open or create file
	fileName := "cfg/valve.rc"
	var fileExists bool
	var appendToTop bool
	var fileContents []string

	if _, err := os.Stat(fileName); err == nil { // Handle file already exists
		fileExists = true

		// Prompt user how to deal with existing file
		var validResponse bool
		fmt.Printf("%v already exists, how would you like to resolve this?\n", fileName)
		for !validResponse {
			fmt.Println("[1] Generate code above the original contents of existing file")
			fmt.Println("[2] Generate code to unique file to be copy + pasted manually")
			fmt.Println("[3] Quit program without generating or modifying the file")
			fmt.Print("Please select and option: ")

			// Use buffered reader because Scanln sucks
			var response string
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if response == "1" {
				appendToTop = true
				validResponse = true
				fileContents = copyFile(fileName)
			} else if response == "2" {
				appendToTop = false
				validResponse = true
				fileName = "cfg/" + getUniqueFilePrefix("valve") + ".rc"
			} else if response == "3" {
				fmt.Println("\nProgram exiting, the file was not generated or modified.")
				os.Exit(0)
			} else {
				fmt.Println("Invalid response, how you would like to resolve this?")
				validResponse = false
			}
			fmt.Println()
		}
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("valve.rc is inaccessible!")
		os.Exit(1)
	}

	// Create valve.rc
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating valve.rc:", err)
		os.Exit(1)
	}

	defer file.Close()

	// Create valve.rc with default values and main config
	file.WriteString("r_decal_cullsize 1\n")
	file.WriteString("exec joystick.cfg\n")
	file.WriteString("exec autoexec.cfg\n")
	file.WriteString("exec " + prefix + ".cfg\n")
	file.WriteString("stuffcmds\n")
	file.WriteString("startupmenu\n")
	file.WriteString("sv_unlockedchapters 99")

	// Valve.rc already exists and user opted to append log-base to the top of it
	if fileExists == true && appendToTop {
		file.WriteString("\n\n//User contents\n")
		pasteFile(file, fileContents)
	}
}

func generateButtonCommands() {
	// Open or create file
	fileName := "logbase_button_copypasta.txt"
	var fileExists bool
	var appendToTop bool
	var fileContents []string

	if _, err := os.Stat(fileName); err == nil { // Handle file already exists
		fileExists = true

		// Prompt user how to deal with existing file
		var validResponse bool
		fmt.Printf("%v already exists, how would you like to resolve this?\n", fileName)
		for !validResponse {
			fmt.Println("[1] Generate code above the original contents of existing file")
			fmt.Println("[2] Generate code to unique file to be copy + pasted manually")
			fmt.Println("[3] Quit program without generating or modifying the file")
			fmt.Print("Please select and option: ")

			// Use buffered reader because Scanln sucks
			var response string
			response, _ = reader.ReadString('\n')  // Read to newline
			response = strings.TrimSpace(response) // Remove newline
			if response == "1" {
				appendToTop = true
				validResponse = true
				fileContents = copyFile(fileName)
			} else if response == "2" {
				appendToTop = false
				validResponse = true
				fileName = getUniqueFilePrefix("logbase_button_copypasta") + ".txt"
			} else if response == "3" {
				fmt.Println("\nProgram exiting, the file was not generated or modified.")
				os.Exit(0)
			} else {
				fmt.Println("Invalid response, how you would like to resolve this?")
				validResponse = false
			}
			fmt.Println()
		}
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("Button commands file is inaccessible!")
		os.Exit(1)
	}

	// Create button commands file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating button commands file:", err)
		os.Exit(1)
	}

	defer file.Close()

	// Create file containing copy + paste template for button code
	file.WriteString("This file contains the command parameter and value for each unique customization option.\n")
	file.WriteString("Create your button as normal, then copy + paste the button code in the appropriate location.\n")
	file.WriteString("You will have to handle the aesthetics and ActionSignalLevel on your own.\n")
	for c := range customizations {
		file.WriteString("\nCustomization " + strconv.Itoa(c+1) + " buttons\n")
		for i := range customizations[c].options {
			customizationAlias := customizations[c].customizationName + strconv.Itoa(i+1)
			file.WriteString("\"command\"\t\t\"engine " + customizationAlias + "\"\n")
		}
	}

	// Button commands file already exists and user opted to append log-base to the top of it
	if fileExists == true && appendToTop {
		file.WriteString("\n//User contents\n")
		pasteFile(file, fileContents)
	}
}

func getBasePath(c int) string {
	srcPath := strings.Split(customizations[c].srcFile, "\\")
	var basePath string

	// Iterate through srcPath and replace directories from root dir with back nav
	for i := 0; i < len(srcPath)-1; i++ {
		basePath += "../"
	}
	// Append with path to customizations file
	basePath += "cfg/" + prefix + "_customizations.txt"

	return basePath
}

func editSource() {
	for c := range customizations {
		// Open source file
		inputFile, err := os.Open(customizations[c].srcFile)
		if err != nil {
			fmt.Println("Error opening source file for reading comments:", err)
			os.Exit(1)
		}

		// Close at end of loop!

		// Create slice containing all the lines and add comments to necessary lines
		var fileContents []string
		lineNum := 1
		paramLinesIndex := 0
		scnr := bufio.NewScanner(inputFile)
		for scnr.Scan() {
			line := scnr.Text()
			if paramLinesIndex < len(customizations[c].paramLines) && lineNum == customizations[c].paramLines[paramLinesIndex] { // If line is to be commented
				commented := "//lb" + line
				fileContents = append(fileContents, commented)
				paramLinesIndex++
			} else { // Non-commented lines
				fileContents = append(fileContents, line)
			}
			lineNum++
		}

		// Rewrite file with comments
		outputFile, err := os.Create(customizations[c].srcFile)
		if err != nil {
			fmt.Println("Error opening file for writing comments:", err)
		}

		// Close at end of loop!

		// Added #base path to top of file if not already there
		basePath := getBasePath(c)
		if !strings.Contains(fileContents[0], basePath) {
			outputFile.WriteString("#base \"" + basePath + "\"\n\n")
		}

		// Repopulate file
		for i := range fileContents {
			outputFile.WriteString(fileContents[i])
			if i < len(fileContents)-1 {
				outputFile.WriteString("\n")
			}
		}
		inputFile.Close()
		outputFile.Close()
	}
}
