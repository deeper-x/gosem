package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// EXIT STATUS LIST
const errRegex = 10
const errReadFile = 20

type semVerFile struct {
	name  string
	regex string
	value []byte
}

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
	inFile := semVerFile{
		name:  "./build_version.html",
		regex: `^v\.[0-9]+\.[0-9]+\.[0-9]+`,
	}

	// 1R-GV
	var semVerB []byte = inFile.genGitTag()
	inFile.writeSemVer(semVerB)

	// 2R-UV [WIP]
	var fileContent []byte = inFile.readCurSemVer()
	var stringContent = string(fileContent)

	var sliceContent = strings.Split(stringContent, ".")

	var major string = strings.TrimSpace(sliceContent[1])
	var minor string = strings.TrimSpace(sliceContent[2])
	var patch string = strings.TrimSpace(sliceContent[3])

	// Now reading only. Next steps are to increment programmatically
	fmt.Printf("Semver level - Major: %v; Minor: %v; Patch: %v\n", major, minor, patch)
}

func (inFile semVerFile) isSemVer(inString string) bool {
	matched, err := regexp.Match(inFile.regex, []byte(inString))

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
