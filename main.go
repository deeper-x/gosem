package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

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
// 1. Read version file. Get a slice []string [maj, min, patch]
// 2. Create an updateMaj func: increment slice[0] and update file
// 3. Create an updateMin func: increment slice[1] and update file
// 4. Create an updatePatch func: increment slice[2] and update file

func main() {
	inFile := semVerFile{
		name:  "./build_version.html",
		regex: `^v\.[0-9]+\.[0-9]+\.[0-9]+`,
	}

	// 1R-GV
	var semVerB []byte = inFile.genGitTag()
	inFile.writeSemVer(semVerB)

	/* 	// 2R-UV
	   	var fileContent []byte = inFile.readCurSemVer()
	   	var stringContent = string(fileContent)

	   	var sliceContent = strings.Split(stringContent, ".")
	   	fmt.Println(sliceContent) */

}

func (inFile semVerFile) isSemVer(inString string) bool {
	matched, err := regexp.Match(inFile.regex, []byte(inString))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
		os.Exit(2)
	}

	return data
}

func (inFile *semVerFile) assignSemVer() {
	value := inFile.readCurSemVer()

	(*inFile).value = value
}

func (inFile semVerFile) genGitTag() []byte {
	cmd := exec.Command("git", "describe --tags")

	out, err := cmd.Output()
	fmt.Println(out)
	if err != nil {
		fmt.Println(err)
	}

	return out
}
