package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	pass = color.FgGreen
	link = color.FgHiBlack
	skip = color.FgYellow
	fail = color.FgHiRed

	skipnotest bool
)

var rootDir = ""
var rootDirReg *regexp.Regexp

var testLineReg = regexp.MustCompile("_test")
var urlReg = regexp.MustCompile("(http:|https:)")

var sourceFileReg = regexp.MustCompile("(.go|.ts|.js|.tsx|.jsx|.dart)")

func main() {
	loadFileDir()
	os.Exit(runCmd(os.Args[1:]))
}

func loadFileDir() {
	file, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	rootDir = file
	rootDirReg = regexp.MustCompile(rootDir + "/")
}

func runCmd(args []string) int {
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		splitLine(string(out))
		return 1
	}
	splitLine(string(out))
	return 0
}

func splitLine(str string) {
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		parse(line)
	}
}

func parse(line string) {
	trimmed := strings.TrimSpace(line)

	var c color.Attribute

	if urlReg.MatchString(line) {
		c = pass
	}

	if sourceFileReg.MatchString(line) {
		if testLineReg.MatchString(line) {
			c = link
		} else {
			c = fail
		}
	}

	switch {
	case strings.Contains(trimmed, rootDir):
		if line == rootDir {
			return
		}
		if rootDirReg.Match([]byte(line)) {
			line = rootDirReg.ReplaceAllString(line, "")
		} else {
			line = strings.ReplaceAll(line, rootDir, ".")
		}

		line = strings.Trim(line, " ")
		line = strings.Trim(line, "\t")

		if testLineReg.MatchString(line) {
			c = link
		} else {
			c = fail
		}
	case strings.Contains(trimmed, "so.go:"):
		return
	case strings.Contains(trimmed, "[no test files]"):
		if skipnotest {
			return
		}
	case strings.HasPrefix(trimmed, "--- PASS"): // passed
		fallthrough
	case strings.HasPrefix(trimmed, "ok"):
		fallthrough
	case strings.HasPrefix(trimmed, "PASS"):
		c = pass

	// skipped
	case strings.HasPrefix(trimmed, "--- SKIP"):
		c = skip

	// failed
	case strings.HasPrefix(trimmed, "--- FAIL"):
		fallthrough
	case strings.HasPrefix(trimmed, "FAIL"):
		c = fail
		color.Set(color.Bold)
	}

	if c != 0 {
		color.Set(c)
		fmt.Printf("%s\n", line)
		color.Unset()
	} else {
		fmt.Printf("%s\n", line)
	}

}
