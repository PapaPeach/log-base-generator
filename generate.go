package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func generateMainConfig() {
	// Open or create file
	var fileExists bool
	if _, err := os.Stat("cfg/"+prefix+".cfg"); err == nil { // Handle file already exists
		fmt.Println("Main config file exists.")
		fileExists = true
		// TODO
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("Main config file is inaccessible!")
		os.Exit(1)
	}

	if fileExists == true {
		os.Exit(1)
	}
	
	// Create main alias file
	file, err := os.Create("cfg/"+prefix+".cfg")
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

	// Create save aliases
	saveAlias := prefix + "_" + customizations[customizationsCount].customizationName + "_dump"
	file.WriteString("//Declare customization save aliases\n")
	file.WriteString("alias " + saveAlias + " \"\"\n")
	file.WriteString("\n")

	// Create default value aliases
	var defaultAlias string
	defaultBraceCount := 0
	for i := 0; i < len(customizations[customizationsCount].customizations[0])-(2*customizations[customizationsCount].numParam); i++ {
		defaultAlias += customizations[customizationsCount].customizations[0][i] + "{"
		defaultBraceCount++
	}
	// Parameter and value
	defaultAlias += strings.Join(customizations[customizationsCount].customizations[0][len(customizations[customizationsCount].customizations[0])-(2*customizations[customizationsCount].numParam):], " ")
	//Close braces
	for j := 0; j < defaultBraceCount; j++ {
		defaultAlias += "}"
	}
	writeAlias := prefix + "_" + customizations[customizationsCount].customizationName + "_write"
	file.WriteString("//Initialize default values\n")
	file.WriteString("alias " + writeAlias + " \"echo " + defaultAlias + "\"\n")
	file.WriteString("\n")

	// Create customization definitions
	// saveAlias [value];writeAlias [value]
	file.WriteString("//Define customization aliases\n")
	for i := range customizations {
		customizationAlias := customizations[customizationsCount].customizationName + strconv.Itoa(i+1)
		var alias string
		braceCount := 0
		// Panel path
		for j := 0; j < len(customizations[customizationsCount].customizations[i])-(2*customizations[customizationsCount].numParam); j++ {
			alias += customizations[customizationsCount].customizations[i][j] + "{"
			braceCount++
		}
		// Parameter and value
		alias += strings.Join(customizations[customizationsCount].customizations[i][len(customizations[customizationsCount].customizations[i])-(2*customizations[customizationsCount].numParam):], " ")
		//Close braces
		for j := 0; j < braceCount; j++ {
			alias += "}"
		}

		file.WriteString("alias " + customizationAlias + 
						" \"alias " + saveAlias + " echo " + customizationAlias + 
						";alias " + writeAlias + " echo " + alias + "\"\n")
	}
}

func generateSaveConfig() {
	// Open or create file
	var fileExists bool
	if _, err := os.Stat("cfg/"+prefix+"_save.cfg"); err == nil { // Handle file already exists
		fmt.Println("Save file exists.")
		fileExists = true
		// TODO
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("Save file is inaccessible!")
		os.Exit(1)
	}

	if fileExists == true {
		os.Exit(1)
	}
	
	// Create main alias file
	file, err := os.Create("cfg/"+prefix+"_save.cfg")
	if err != nil {
		fmt.Println("Error creating save file:", err)
		os.Exit(1)
	}

	defer file.Close()
	
	// Create aliases to dump save aliases to file
	file.WriteString("//Clear and prep log file\n")
	file.WriteString("lb_log_selection_open\n")
	file.WriteString("\n")

	saveAlias := prefix + "_" + customizations[customizationsCount].customizationName + "_dump"
	file.WriteString("//Dump current aliases to file\n")
	file.WriteString(saveAlias + "\n")
	file.WriteString("\n")

	file.WriteString("//Close log file\n")
	file.WriteString("con_logfile \"\"")
}

func generateGeneratorConfig() {
	// Open or create file
	var fileExists bool
	if _, err := os.Stat("cfg/"+prefix+"_generate.cfg"); err == nil { // Handle file already exists
		fmt.Println("Genertor file exists.")
		fileExists = true
		// TODO
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("Generator file is inaccessible!")
		os.Exit(1)
	}

	if fileExists == true {
		os.Exit(1)
	}
	
	// Create main alias file
	file, err := os.Create("cfg/"+prefix+"_generate.cfg")
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

	writeAlias := prefix + "_" + customizations[customizationsCount].customizationName + "_write"
	file.WriteString("//Write current customizations to file\n")
	file.WriteString(writeAlias + "\n")
	file.WriteString("\n")

	file.WriteString("//Close log file\n")
	file.WriteString("echo \"}\"\n")
	file.WriteString("con_logfile \"\"")
}

func generateValveRc() {
	fmt.Println("TODO generate valve.rc")
}

func generateButtonCommands() {
	fmt.Println("TODO generate button commands")
}

func commentSource() {
	// Open source file
	inputFile, err := os.Open(customizations[customizationsCount].srcFile)
	if err != nil {
		fmt.Println("Error opening source file for reading comments:", err)
		os.Exit(1)
	}
	defer inputFile.Close()
	
	// Create slice containing all the lines and add comments to necessary lines
	var fileContents []string
	lineNum := 1
	paramLinesIndex := 0
	scnr := bufio.NewScanner(inputFile)
	for scnr.Scan() {
		line := scnr.Text()
		if paramLinesIndex < len(customizations[customizationsCount].paramLines) && lineNum == customizations[customizationsCount].paramLines[paramLinesIndex] { // If line is to be commented
			commented := "//lb" + line
			fileContents = append(fileContents, commented)
			paramLinesIndex++
		} else { // Non-commented lines
			fileContents = append(fileContents, line)
		}
		lineNum++
	}
	
	// Rewrite file with comments
	outputFile, err := os.Create(customizations[customizationsCount].srcFile)
	if err != nil {
		fmt.Println("Error opening file for writing comments:", err)
	}
	defer outputFile.Close()

	for i := range fileContents {
		outputFile.WriteString(fileContents[i])
		if i < len(fileContents)-1 {
			outputFile.WriteString("\n")
		}
	}
}