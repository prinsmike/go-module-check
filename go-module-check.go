package main

import (
	"bytes"
	"fmt"
	"os"
	"log"
	"path/filepath"
	"os/exec"
	"github.com/mgutz/ansi"
	"strings"
)

type SiteList struct {
	Path  string
	Sites []string
}

func inList(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func main() {

	var e error

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments.")
	}

	var S SiteList
	
	blacklist := []string{"all", "default"}

	S.Path, e = os.Getwd()
	if e != nil {
		log.Fatal(e)
	}

	module := os.Args[1]

	// Get all file paths in the current working directory.
	files, e := filepath.Glob(fmt.Sprintf("%s/*", S.Path))
	if e != nil {
		log.Fatal(e)
	}

	// Loop over files and append only directories, that are not blacklisted, to S.Sites.
	for _, file := range files {
		_, f := filepath.Split(file)

		finfo, e := os.Lstat(file)
		if e != nil {
			log.Fatal(e)
		}

		if finfo.IsDir() && !inList(f, blacklist) {
			S.Sites = append(S.Sites, f)
		}
	}

	var tcount, encount, dcount int

	green := ansi.ColorCode("green+h:black")
	red := ansi.ColorCode("red+h:black")
	reset := ansi.ColorCode("reset")

	fmt.Println(S.Path)

	for _, site := range S.Sites {
		tcount++
		drushcommand := exec.Command("drush", "-l", site, "pmi", module)
		grepcommand := exec.Command("grep", "Status")
		grepcommand.Stdin, _ = drushcommand.StdoutPipe()

		// Create a buffer of bytes.
		var b bytes.Buffer

		// Assign the address of our buffer to grepcommand.Stdout.
		grepcommand.Stdout = &b

		// Start grepcommand.
		_ = grepcommand.Start()

		// Run syscommand
		_ = drushcommand.Run()

		// Wait for grepcommand to exit.
		_ = grepcommand.Wait()

		s := fmt.Sprintf("%s", &b)

		if strings.Contains(s, "enabled") {
			encount++
			fmt.Printf("%sModule %s is enabled on %s.%s\n", green, module, site, reset)
		} else {
			dcount++
			fmt.Printf("%sModule %s is not enabled on %s.%s\n", red, module, site, reset)
		}
	}

	fmt.Printf("\n\nModule: %s", module)
	fmt.Printf("Total number of sites: %d\n", tcount)
	fmt.Printf("%s%s is enabled on %d sites.%s\n", green, module, encount, reset)
	fmt.Printf("%s%s is disabled on %d sites.%s\n", red, module, dcount, reset)
}
