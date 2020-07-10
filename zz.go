package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var ignoreList = []string{
	".",
	"..",
	"node_modules",
	".git",
	"target",
	".idea",
}

func file_line() string {
	_, fileName, fileLine, ok := runtime.Caller(2)
	var s string
	if ok {
		s = fmt.Sprintf("%s:%d", fileName, fileLine)
	} else {
		s = ""
	}
	return s
}

func check(err error) {
	if err != nil {
		log.Fatalf("%s :: %v\n", file_line(), err)
	}
}
func isIgnored(s string) bool {
	if strings.HasPrefix(s, ".") {
		return true
	}

	for _, x := range ignoreList {
		if x == s {
			return true
		}
	}
	return false
}

func getMTimes(root string) (map[string]time.Time, error) {
	dirs := []string{root}
	files := make(map[string]time.Time)

	for len(dirs) > 0 {
		// pop
		dir := dirs[len(dirs) - 1]
		dirs = dirs[:len(dirs) - 1]

		xs, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		Loop:
		for _, x := range xs {
			switch {
			case isIgnored(x.Name()):
				continue Loop
			case x.IsDir():
				dirs = append(dirs, dir + "/" + x.Name())
			default:
				fname := dir + "/" + x.Name()
				files[fname] = x.ModTime()
			}
		}
	}
	return files, nil
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	check(err)
}

func runCmd(xs []string) {
	cmd := exec.Command(xs[0], xs[1:]...)

	// hmm.......  @_@
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	check(err)
}

func mapEqual(m1, m2 map[string]time.Time) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k := range m1 {
		if m1[k] != m2[k] {
			clear()
			//fmt.Printf("%v\n%v\n%v\n", k, m1[k], m2[k])
			return false
		}
	}
	return true
}

func main() {
	watchDir := os.Args[2]
	xs := os.Args[3:]

	oldMTimes := make(map[string]time.Time)
	clear()
	for {
		newMTimes, err := getMTimes(watchDir)
		check(err)

		if !mapEqual(oldMTimes, newMTimes) {
			oldMTimes = newMTimes
			runCmd(xs)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
