// Written by Mikhail P. Ortiz-Lunyov (mportizlunyov)
//
// Version 1.0.0-release (August 19th 2024)
//
// This script is licensed under the GNU Public License v3 (GPLv3)
// Intended for use on Linux to check the specific distro running, using native Linux tools.
// This is useful when developing programs to adapt to specific linux environments.
//
// This is the front-end that uses the functions from the Linux Distro Engine.
// This is only compiled/run if being used as a stand-alone program.

/*
Linux Distro Engine (Go edition) main package.

This package serves as a stand-alone application of the linuxdistroengine go package.
It exists to demonstrate the capabilities of the linuxdistroengine package.
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

// Main method of the program.
//
// This method declares the flags and sets the default value to use the
// linuxdistroengine package.
func main() {
	// Declare variables and their defaults
	var option string = "id" // Default argument for the engine
	var argsCount int = 0

	// Check for flags
	// // -v / --version
	versionShort := flag.Bool("v", false, "Print the version number of this module")
	versionLong := flag.Bool("version", false, "Long form of [-v]")
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
	versionFlag := *versionShort || *versionLong

	// Prepare options based on parameter
	switch true {
	case versionFlag:
		option = "v"
		argsCount++
	}
	switch true {
	case prettynameFlag:
		option = "pn"
		argsCount++
	}
	switch true {
	case kernelFlag:
		option = "k"
		argsCount++
	}

	// Define final action based on parameters
	switch argsCount {
	case 0:
		fallthrough
	case 1:
		pkgResult, pkgErrCode := linuxdistorengine.DistroResult(option)
		switch pkgErrCode {
		case 0:
			fmt.Println(pkgResult)
		case 1:
			fallthrough
		case 2:
			fallthrough
		case 3:
			fallthrough
		case 4:
			fallthrough
		case 5:
			switch pkgErrCode {
			case 1:
				fmt.Println("This program is intended to run on Linux")
			case 2:
				fmt.Println("[uname -r] command failed!")
			case 3:
				fmt.Println("Reading [os-release] file failed")
			case 4:
				fmt.Println("Neither ID nor PRETTY_NAME found in */os-release file")
			case 5:
				fmt.Println("Bad argument for DistroResult() method [ ", option, " ]")
			}
			os.Exit(1)
		case 44:
			fmt.Println(pkgResult)
			os.Exit(44)
		default:
			fmt.Println("INTERNAL ERROR, MISSING ERROR CODE [", pkgErrCode, "]")
			os.Exit(254)
		}
	default:
		fmt.Println("Incompatible arguments")
		os.Exit(1)
	}
}
