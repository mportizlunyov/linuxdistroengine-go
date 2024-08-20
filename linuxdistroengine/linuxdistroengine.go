// Written by Mikhail P. Ortiz-Lunyov (mportizlunyov)
//
// Version 1.0.0-release (August 19th 2024)
//
// This script is licensed under the GNU Public License v3 (GPLv3)
// Intended for use on Linux to check the specific distro running, using native Linux tools.
// This is useful when developing programs to adapt to specific linux environments.
//
// This is the actual engine which identified the Linux Distribution,
//  or at least the family being run.

/*
linuxdistroengine is a Go package that returns the name of the Linux distro being run.
This is the main file that actually serves as the importable package.

The primary method that returns the name (or at least the family) of the current Linux distro is the DistroResult() method.

Possible DistroResult() string arguments:

	"id": Print the basic, minimal ID
	"k":  Print the kernel version of the distro
	"pn": Print the 'Pretty Name' of the distro, often including the version number
	"v":  Print the version of this engine

The DistroResult() method also provides the following error Codes:

	0:   Complete, correct completion
	1:   OS run is NOT Linux
	2:   'uname -r' command failed
	3:   /etc/os-release and /lib/os-release not found
	4:   Required section (ID & PRETTY_NAME) not found in os-release file
	5:   Invalid argument for DistroResult() method
	44:  Distro nor family not found (reference to HTTP 404 error)
	254: Developer has not yet set error
*/
package linuxdistroengine

// Import packages
import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Declare fields
// // Set constants

// Version constants
const (
	SHORT_VERSION string = "1.0.0"
	VERSION_NAME  string = "August 19th 2024"
	DEV_VERSION   string = "-release"
	LONG_VERSION  string = "v" + SHORT_VERSION + DEV_VERSION + " (" + VERSION_NAME + ")"
)

// Functional package constants
const (
	DEFAULT          string = "UNKNOWN_DISTRO"
	INVALID          string = "INVALID_OPERATION"
	OSRELEASE_AMOUNT int    = 2
)

// // Other fields

// An array containing potential paths to the os-release file
var OSRELEASEFILE [OSRELEASE_AMOUNT]string = [OSRELEASE_AMOUNT]string{"/etc/os-release", "/lib/os-release"}

// The string of the final result of the Distro Engine
var found string

// Checks if the os-release file exists, based on the OSRELEASEFILE[] field
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
func readOSReleaseFile(option string, fileNum int) (string, int) {
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
		// Failed to read */os-release file, return INVALID constant
		return INVALID, 3
	}

	// Prepare filtering variables for the respective method param.
	switch option {
	case "pn":
		selectionCriteria = "PRETTY_NAME=\""
		substringFirst = 13
	case "id":
		selectionCriteria = "ID="
		substringFirst = 3
	}

	// // Iterate through fileContentSliced, locating and isolating the selected criteria
	for i := 0; i < len(fileContentSliced); i++ {
		switch strings.HasPrefix(fileContentSliced[i], selectionCriteria) {
		case true:
			// Select specific slice index
			fileSelection = fileContentSliced[i]
			i = len(fileContentSliced) // End loop early
		}
		// If iteration is at the last line, but fileSelection variable is still empty, return INVALID constant.
		if i == len(fileContentSliced)-1 && fileSelection == "" {
			return INVALID, 4
		}
	}

	// Format fileSelection variable to remove quotation marks, if needed (ID only)
	switch option {
	case "id":
		switch strings.Contains(fileSelection[substringFirst:], "\"") {
		case true:
			// The following works because "\"DISTRO\"" -> ["", "DISTRO", ""]
			//                                              0|   |-1--|   |2
			return strings.Split(fileSelection[substringFirst:], "\"")[1], 0
		default:
			return fileSelection[substringFirst:], 0
		}
	default:
		return fileSelection[substringFirst : len(fileSelection)-1], 0
	}
}

// Checks the Linux family type based on its package manager
func pkgManCheck() string {
	// Initialise variables
	var pkgMan string
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

/*
Use alternative means to find the distro name.
Such means includes checking package managers installed to at least identify the Distro family.
*/
func oSReleaseAlt() string {
	// Initialise variables
	var potReturnVal string
	// Discover family using package manager
	switch pkgManCheck() {
	case "dpkg":
		// Check for Ubuntu or derivative
		stdout, cmdErr := exec.Command("lsb_release", "-a").Output()
		switch cmdErr {
		case nil:
			switch strings.Contains(string(stdout[:]), "ubuntu") {
			case true:
				potReturnVal = "Ubuntu"
			case false:
				potReturnVal = "Debian"
			}
		}
	case "rpm":
		potReturnVal = "RedHat"
	default:
		potReturnVal = DEFAULT
	}

	return potReturnVal
}

// Returns result of the Linux Distro Engine (see top of page for details)
func DistroResult(option string) (string, int) {
	// Declare variables
	var invalid_OpScore int

	// Set fields with values
	found = DEFAULT // Set default OS label
	option = strings.ToLower(option)

	// Check if the kernel is Linux or not
	switch runtime.GOOS == "linux" {
	case false:
		return INVALID, 1
	}

	// Take action based on method parameter
	switch option {
	case "id": // Identify distro
		fallthrough
	case "pn": // Same as above
		var osReleaseExist bool
		var osReleaseNum int
		// Check if */os-release file exists
		osReleaseExist, osReleaseNum = OSReleaseFileExist()
		switch osReleaseExist {
		case true:
			found, invalid_OpScore = readOSReleaseFile(option, osReleaseNum)
			switch invalid_OpScore {
			case 0: // Nothing invalid, continue
			default:
				return INVALID, invalid_OpScore
			}
		case false:
			found = oSReleaseAlt()
		}
	case "k": // Get Kernel version
		stdout, cmdErr := exec.Command("uname", "-r").Output()
		switch cmdErr {
		case nil:
			found = string(stdout[:len(stdout)-1]) // Remove last newline
		default:
			return INVALID, 2
		}
	case "v": // Get Linux Distro Engine version
		found = "LinuxDistroEngine-Go " + LONG_VERSION
	default: // Return bad argument error
		return INVALID, 5
	}

	// Return final result
	switch found {
	case DEFAULT:
		return found, 44
	default:
		return found, 0
	}
}
