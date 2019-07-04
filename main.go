package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

type semVerFile struct {
	name  string
	regex string
}

func main() {
	inFile := semVerFile{
		name:  "<PATH>",
		regex: `^v[0-9]+\.[0-9]+\.[0-9]+`,
	}

	//semVerTag := genSemVerTag()
	curSemVer := inFile.readCurSemVer()
	//inFile.writeSemVer(semVerTag)
	itIs := inFile.isSemVer(string(curSemVer))

	if !itIs {
		fmt.Println(string(curSemVer), "is not")
	} else {
		fmt.Println("it is")
	}
}

func (inFile semVerFile) isSemVer(inString string) bool {
	matched, err := regexp.Match(inFile.regex, []byte(inString))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return matched
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

func genSemVerTag() []byte {
	cmd := exec.Command("git", "describe")

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
	}

	return out
}
