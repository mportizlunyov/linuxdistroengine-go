// Written by Mikhail P. Ortiz-Lunyov
//
// Version 0.0.1-release (July 11th 2024)
//
// This script is licensed under the GNU Public License v3 (GPLv3)
// Intended for use on Linux to check the specific distro running, using native Linux tools.
// This is useful when developing programs to adapt to specific linux environments.
//
// This is the actual engine which defines the Linux Distribution, or at least the family being run.

// Package name
package linuxdistorengine

// Import packages
import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Declare fields
// // Set constants
// // // Version
const SHORT_VERSION string = "0.0.1"
const VERSION_NAME string = "July 11th 2024"
const DEV_VERSION string = "-release"
const LONG_VERSION string = "v" + SHORT_VERSION + DEV_VERSION + " (" + VERSION_NAME + ")"

// // // Functional
const DEFAULT string = "UNKNOWN"
const OSRELEASE_AMOUNT = 2

// // Other fields
var OSRELEASEFILE [OSRELEASE_AMOUNT]string = [OSRELEASE_AMOUNT]string{"/etc/os-release", "/lib/os-release"}
var found string

// Function to centralize all errors
func errorAll(compatibility bool, readOsReleaseErr bool, distroNotFound bool) int {
	// Check if error is related to compatibility
	switch compatibility {
	case true:
		fmt.Println("This program is intended to run on Linux")
		return 1
	}
	switch readOsReleaseErr {
	case true:
		fmt.Println("Reading [*/os-release] file failed")
		return 1
	}
	// Check if no distro was found
	switch distroNotFound {
	case true:
		fmt.Println("Unidentifiable distro, not even family identified")
		return 404
	}

	// Default, if error has not been managed with specific error message
	fmt.Println("ERROR MESSAGE NOT COMPLETE")
	return 254
}

// Checks if the /etc/os-release file exists
func OSReleaseFileExist() (bool, int) {
	// Check if os.Stat() method returns an error
	for i := 0; i < OSRELEASE_AMOUNT; i++ {
		_, existsErr := os.Stat(OSRELEASEFILE[i])
		switch existsErr {
		case nil:
			return true, i
		}
	}

	return false, -1
}

// Reads OS Release file
func readOSReleaseFile(option string, fileNum int) string {
	// Declare variables
	var fileContent string
	var fileContentSliced []string
	var selectionCriteria string
	var fileSelection string
	var substringFirst int
	// var returnName string

	// Read the os-release file
	content, readErr := os.ReadFile(OSRELEASEFILE[fileNum])
	switch readErr {
	case nil:
		// Convert []byte to string
		fileContent = string(content[:])
		// Slice fileContent file into fileContentSliced
		fileContentSliced = strings.Split(fileContent, "\n")
	default:
		os.Exit(errorAll(false, true, false))
	}

	// Extract ID (Default), or PRETTY_NAME
	// // Prepare selectionCriteria variable
	switch option {
	// PRETTY_NAME=
	case "pn":
		selectionCriteria = "PRETTY_NAME=\""
		substringFirst = 13
	// ID=
	case "id":
		selectionCriteria = "ID="
		substringFirst = 3
	}

	// // Iterate through fileContentSliced, locating selected criteria
	for i := 0; i < len(fileContentSliced); i++ {
		switch strings.HasPrefix(fileContentSliced[i], selectionCriteria) {
		case true:
			// Select specific slice index
			fileSelection = fileContentSliced[i]
			i = len(fileContentSliced) // End loop early
		}
	}

	// Format fileSelection variable to remove quotation marks, if needed (ID only)
	switch option {
	case "id":
		switch strings.Contains(fileSelection[substringFirst:], "\"") {
		case true:
			// The following works because "\"DISTRO\"" -> ["", "DISTRO", ""]
			//                                                   |-1--|
			return strings.Split(fileSelection[substringFirst:], "\"")[1]
		default:
			return fileSelection[substringFirst:]
		}
	default:
		return fileSelection[substringFirst : len(fileSelection)-1]
	}
}

// Use alterantive means to find the distro name
func oSReleaseAlt() string {
	// Initialise variables
	var potFamily string
	var potReturnVal string
	// Discover family using package manager
	potFamily = pkgManCheck()
	switch potFamily {
	case "dpkg":
		// Check for Ubuntu or derivative
		stdout, cmdErr := exec.Command("lsb_release", "-a").Output()
		switch cmdErr {
		case nil:
			switch strings.Contains(string(stdout[:]), "ubuntu") {
			case true:
				potFamily = "Ubuntu"
			case false:
				potFamily = "Debian"
			}
		}
	case "rpm":
		potFamily = "RedHat"
	}

	return potReturnVal
}

// Checks the Linux family type based on its package manager
func pkgManCheck() string {
	// Initialise variables
	var pkgMan string
	var returnVal string
	// Check package managers
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			pkgMan = "dpkg"
		case 1:
			pkgMan = "rpm"
		default:
			i = 3 // End early
		}

		_, cmdErr := exec.Command(pkgMan, "--help").Output()
		switch cmdErr {
		case nil:
			returnVal = pkgMan
			i = 3 // End early
		}
	}

	return returnVal
}

// Checks if the script is run on a Linux OS
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// Returns result of Linux Distro Engine
func DistroResult(option string, verbose bool) string {
	// Declare variables
	var osReleaseExist bool
	var osReleaseNum int
	// // Set fields
	found = DEFAULT                  // Set default OS label
	option = strings.ToLower(option) // Filter potential bad params

	// Check if the OS is Linux or not
	switch IsLinux() {
	case false:
		os.Exit(errorAll(true, false, false))
	}

	// Check if "k" option is called
	switch option {
	case "k":
		stdout, cmdErr := exec.Command("uname", "-r").Output()
		switch cmdErr {
		case nil:
			found = string(stdout[:len(stdout)-1]) // Remove last newline
		default:
			os.Exit(1)
		}
	default:
		// Check if /etc/os-release exists
		osReleaseExist, osReleaseNum = OSReleaseFileExist()
		switch osReleaseExist {
		case true:
			found = readOSReleaseFile(option, osReleaseNum)
		case false:
			found = oSReleaseAlt()
		}
	}

	// Return final result
	return found
}
