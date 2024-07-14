// Written by Mikhail P. Ortiz-Lunyov
//
// Version 0.0.4-release (July 13th 2024)
//
// This script is licensed under the GNU Public License v3 (GPLv3)
// Intended for use on Linux to check the specific distro running, using native Linux tools.
// This is useful when developing programs to adapt to specific linux environments.
//
// This is the front-end that uses the functions from the Linux Distro Engine

/*
Linux Distro Engine (Go ed.) main package.

This package serves as the main package for the whole Linux Distro Engine.
It contains the flags and default values to run the linuxdistroengine package,
especially as a stand-alone compiled program
*/
package main

// Import packages
import (
	"flag"
	"fmt"
	"os"

	linuxdistorengine "github.com/mportizlunyov/linuxdistroengine-go/linuxdistroengine"
)

// Script-level fields
// // Constants

// Version constants
const (
	SHORT_VERSION string = "0.0.4"
	VERSION_NAME  string = "July 13th 2024"
	DEV_VERSION   string = "-release"
	LONG_VERSION  string = "v" + SHORT_VERSION + DEV_VERSION + " (" + VERSION_NAME + ")"
)

// Verbosity flag field (not implemented as of v0.0.3)
var verboseFlag bool

// Main method of the program.
//
// This method declares the flags and sets the default value to use the
// linuxdistroengine package.
func main() {
	// Declare variables
	var option string

	// Check for flags
	versionShort := flag.Bool("v", false, "Print the version number of this module")
	versionLong := flag.Bool("version", false, "Long form of [-v]")
	// // -vb / --verbose
	verboseShort := flag.Bool("vb", false, "Print verbose output")
	verboseLong := flag.Bool("verbose", false, "Long form of [-vb]")
	// // -pn / --pretty-name
	prettynameShort := flag.Bool("pn", false, "Print the full 'Pretty Name' of the distro, if applicable")
	prettynameLong := flag.Bool("pretty-name", false, "Long form of [-pn]")
	// // -k / --kernel
	kernelShort := flag.Bool("k", false, "Print the kernel version")
	kernelLong := flag.Bool("kernel", false, "Long form of [-v]")
	// // Parse flags
	flag.Parse()
	// // Finalise flags
	prettynameFlag := *prettynameLong || *prettynameShort
	kernelFlag := *kernelLong || *kernelShort
	verboseFlag = *verboseLong || *verboseShort
	versionFlag := *versionShort || *versionLong

	// Print version statement and exit if needed
	switch versionFlag {
	case true:
		fmt.Println("LinuxDistroEngine-Go " + LONG_VERSION)
		os.Exit(0)
	}

	// Filter bad parameters, and prepare options
	if prettynameFlag && kernelFlag {
		fmt.Println("-pn & -k flags are incompatible")
		os.Exit(1)
	} else if prettynameFlag {
		option = "pn"
	} else if kernelFlag {
		option = "k"
	} else {
		option = "id"
	}

	// Print Distro result
	fmt.Println(linuxdistorengine.DistroResult(option, verboseFlag))
}
