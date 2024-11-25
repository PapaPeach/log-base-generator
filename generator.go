package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func generateConfig() {
	// Open or create file
	var fileExists bool
	if _, err := os.Stat("cfg/"+prefix+".cfg"); err == nil { // Handle file already exists
		fmt.Println("File exists.")
		fileExists = true
		// TODO
	} else if errors.Is(err, os.ErrNotExist) { // Create fresh file
		fileExists = false
	} else { // Oh shit
		fmt.Println("File is inaccessible!")
		os.Exit(1)
	}

	if fileExists == true {
		os.Exit(1)
	}
	
	// Create main alias file
	file, err := os.Create("cfg/"+prefix+".cfg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}

	defer file.Close()

	// Create log_open alias
	file.WriteString("//Open log for writing alias\n")
	file.WriteString("alias lb_log_open \"sixense_clear_bindings;sixense_write_bindings " + prefix + "_customization_selection.txt;con_timestamp 0;con_logfile cfg/" + prefix + "_customization_selection.txt\"\n")
	file.WriteString("\n")

	// Create save aliases
	file.WriteString("//Declare customization save aliases\n")
	file.WriteString("alias " + prefix + "_" + strings.ToLower(customizations[0][len(customizations[0])-3] + customizations[0][len(customizations[0])-2]) + "_dump \"\"\n")
	file.WriteString("\n")

	// Create default value aliases
	var defaultAlias string
	defaultBraceCount := 0
	for i := 0; i < len(customizations[0])-2; i++ {
		defaultAlias += customizations[0][i] + "{"
		defaultBraceCount++
	}
	// Parameter and value
	defaultAlias += customizations[0][len(customizations[0])-2] + " " + customizations[0][len(customizations[0])-1]
	//Close braces
	for j := 0; j < defaultBraceCount; j++ {
		defaultAlias += "}"
	}
	file.WriteString("//Initialize default values\n")
	file.WriteString("alias " + prefix + "_" + strings.ToLower(customizations[0][len(customizations[0])-3] + customizations[0][len(customizations[0])-2]) + "_write \"echo " + defaultAlias + "\"\n")
	file.WriteString("\n")

	// Create customization definitions
	file.WriteString("//Define customization aliases\n")
	for i := range customizations {
		var alias string
		braceCount := 0
		// Panel path
		for j := 0; j < len(customizations[i])-2; j++ {
			alias += customizations[i][j] + "{"
			braceCount++
		}
		// Parameter and value
		alias += customizations[i][len(customizations[i])-2] + " " + customizations[i][len(customizations[i])-1]
		//Close braces
		for k := 0; k < braceCount; k++ {
			alias += "}"
		}

		file.WriteString("alias " + strings.ToLower(customizations[0][len(customizations[0])-3] + customizations[0][len(customizations[0])-2]) + strconv.Itoa(i+1) + "_write \"echo " + alias + "\"\n")
	}
}