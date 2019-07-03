package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

var semverFile = "<PATH>"

func main() {
	semVer := genSemVer()
	writeSemVer(semVer)
}

func writeSemVer(fileContent []byte) {
	ioutil.WriteFile(semverFile, fileContent, 0644)
}

func genSemVer() []byte {
	cmd := exec.Command("git", "describe")

	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
	}

	return out
}
