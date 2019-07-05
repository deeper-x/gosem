package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// EXIT STATUS LIST
const errRegex = 10
const errReadFile = 20
const errInArgs = 30
const errFileNotExists = 50

type semVerFile struct {
	name  string
	regex string
	value []byte
}

var fileToWritePath string
var semVerRegex = `^v\.[0-9]+\.[0-9]+\.[0-9]+`
var usageString = fmt.Sprintf("USAGE: %v <file_to_write> <action>", os.Args[0])

// 1R-GV: GEN VERSION
// 1. launch bin: get git tag and write it to file
// 1A. get git describe value
// 1B. overwrite it in file

// 2R-UV: UPDATE Maj-Min-Patch
// 1. Read git tag. Get a slice []string [maj, min, patch]
// 2. Write it to version file
// 3. Create an updateMaj func: increment slice[0] and update file
// 4. Create an updateMin func: increment slice[1] and update file
// 5. Create an updatePatch func: increment slice[2] and update file

func main() {
	// STEP 1 - parameters check
	parseArgs()
}

func parseArgs() {
	// CLI argument is the file to write
	if len(os.Args) != 3 {
		fmt.Println(usageString)
		os.Exit(errInArgs)
	}

	if !fileExists(os.Args[1]) {
		fmt.Println("file path is not correct...")
		os.Exit(errFileNotExists)
	} else {
		fileToWritePath = os.Args[1]
		fmt.Println("writing semver to...", fileToWritePath)
	}

	inFile := semVerFile{
		name:  fileToWritePath,
		regex: semVerRegex,
	}

	switch os.Args[2] {
	case "update":
		fmt.Println("updating...")
		inFile.updateFile()

	case "major":
		fmt.Println("#TODO major...")

	case "minor":
		fmt.Println("#TODO minor...")

	case "patch":
		fmt.Println("#TODO patch...")

	default:
		fmt.Println(usageString)
		os.Exit(errInArgs)
	}
}

func fileExists(pathToFile string) bool {
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func newMajor(inVersion []string) []string {
	actualMaj, err := strconv.Atoi(inVersion[0])

	if err != nil {
		fmt.Println("error in newMajor", err)
	}

	return []string{string(actualMaj + 1), "0", "0"}
}

func newMinor(inVersion []string) []string {
	actualMin, err := strconv.Atoi(inVersion[1])

	if err != nil {
		fmt.Println("error in newMinor", err)
	}

	return []string{inVersion[0], string(actualMin + 1), "0"}
}

func newPatch(inVersion []string) []string {
	actualPatch, err := strconv.Atoi(inVersion[2])

	if err != nil {
		fmt.Println("error in newPatch", err)
	}

	return []string{inVersion[0], inVersion[1], string(actualPatch + 1)}
}

func (inFile semVerFile) updateGitTag() bool {
	var semVerVal string = string(inFile.genGitTag())

	cmd := exec.Command("git", "tag", "-a", semVerVal, "-m", "auto updated")
	err := cmd.Run()

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (inFile semVerFile) updateFile() {
	// STEP 2 - READING GIT DATA
	// 1R-GV
	var semVerB []byte = inFile.genGitTag()

	// git tag must be in the *standard* format
	if !inFile.isSemVer(semVerB) {
		os.Exit(errRegex)
	}

	// STEP 3 - WRITING DATA TO FILE
	// write git tag we read -> to file
	inFile.writeSemVer(semVerB)

	// 2R-UV [WIP]
	// read file content
	var fileContent []byte = inFile.readCurSemVer()

	// check file content
	if !inFile.isSemVer(fileContent) {
		os.Exit(errRegex)
	}

	var result []string = getSemVerSplit(fileContent)

	// Now reading only. Next steps are to increment programmatically
	fmt.Printf("Semver level - Major: %v; Minor: %v; Patch: %v\n", result[0], result[1], result[2])
}

func getSemVerSplit(inVersion []byte) []string {
	// get file content to slice
	var stringContent = string(inVersion)
	var sliceContent = strings.Split(stringContent, ".")

	var major string = strings.TrimSpace(sliceContent[1])
	var minor string = strings.TrimSpace(sliceContent[2])
	var patch string = strings.TrimSpace(sliceContent[3])

	return []string{major, minor, patch}
}

func (inFile semVerFile) isSemVer(inString []byte) bool {
	matched, err := regexp.Match(inFile.regex, inString)

	if err != nil {
		os.Exit(errRegex)
	}

	return matched
}

func (inFile semVerFile) getLabels() []string {
	retSlice := strings.Split(string(inFile.value), ".")

	return retSlice
}

func (inFile semVerFile) writeSemVer(fileContent []byte) {
	ioutil.WriteFile(inFile.name, fileContent, 0644)
}

func (inFile semVerFile) readCurSemVer() []byte {
	data, err := ioutil.ReadFile(inFile.name)

	if err != nil {
		fmt.Println(err)
		os.Exit(errReadFile)
	}

	return data
}

func (inFile *semVerFile) assignSemVer() {
	value := inFile.readCurSemVer()

	(*inFile).value = value
}

func (inFile semVerFile) genGitTag() []byte {
	cmd := exec.Command("git", "describe", "--tags")

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
	}

	return out
}
