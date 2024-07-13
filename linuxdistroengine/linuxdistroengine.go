// Written by Mikhail P. Ortiz-Lunyov
//
// Version 0.0.3-beta (July 12th 2024)
//
// This script is licensed under the GNU Public License v3 (GPLv3)
// Intended for use on Linux to check the specific distro running, using native Linux tools.
// This is useful when developing programs to adapt to specific linux environments.
//
// This is the actual engine which defines the Linux Distribution, or at least the family being run.

/*
linuxdistroengine is a Go package that returns the name of the Linux distro being run.

The main method that does this is the DistroResult(*Argument*) method.

Possible *Argument*s:

	"id": Print
	"k": Print the kernel version of the distro
	"pn": Print the 'Pretty Name' of the distro, often including the version number

Exit Codes:

	0: All good
	1: Generic error, see description
	44: Distro nor family not found
	254: Developer has not yet set error
*/
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

// Version constants
const (
	SHORT_VERSION string = "0.0.3"
	VERSION_NAME  string = "July 12th 2024 II"
	DEV_VERSION   string = "-release"
	LONG_VERSION  string = "v" + SHORT_VERSION + DEV_VERSION + " (" + VERSION_NAME + ")"
)

// Functional package constants
const (
	DEFAULT          string = "UNKNOWN_DISTRO"
	OSRELEASE_AMOUNT int    = 2
)

// // Other fields

// An array containing potential paths to the os-release file
var OSRELEASEFILE [OSRELEASE_AMOUNT]string = [OSRELEASE_AMOUNT]string{"/etc/os-release", "/lib/os-release"}

// The string of the final result of the Distro Engine
var found string

// Method to centralize all errors and return the appropriate exit code
func errorAll(errorScore int) int {
	switch errorScore {
	// Incompatible OS
	case 0:
		fmt.Println("This program is intended to run on Linux")
		return 1
	// Reading os-release file failed
	case 1:
		fmt.Println("Reading [*/os-release] file failed")
		return 1
	// Distro unidentified, not even family
	case 2:
		fmt.Println("Unidentifiable distro, not even family identified")
		return 44
	// Unset error messages
	default:
		// Default, if error has not been managed with specific error message
		fmt.Println("ERROR MESSAGE NOT COMPLETE")
		return 254
	}
}

// Checks if the os-release file exists, based on the OSRELEASEFILE field
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

	// Read the os-release file
	content, readErr := os.ReadFile(OSRELEASEFILE[fileNum])
	switch readErr {
	case nil:
		// Convert []byte to string
		fileContent = string(content[:])
		// Slice fileContent file into fileContentSliced
		fileContentSliced = strings.Split(fileContent, "\n")
	default:
		// Failed to read */os-release file
		os.Exit(errorAll(1))
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

// Use alterantive means to find the distro name.
// Such means includes checking package managers installed to at least identify the Distro family.
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
	// var returnVal string
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
			return pkgMan
		}
	}

	// If none of the package managers are detected, return DEFAULT constant
	return DEFAULT
}

// Returns true if the OS being run is Linux
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
		os.Exit(errorAll(0))
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
		// Check if */os-release file exists
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
