package main

import (
	"fmt"
	"os"
	"log"
	"path/filepath"
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

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments.")
	}

	var S SiteList
	
	blacklist := []string{"all", "default"}

	S.Path = filepath.Clean(os.Args[1])

	var visit = func(p string, f os.FileInfo, e error) error {
		if f.IsDir() && p != S.Path {
			_, file := filepath.Split(p)
			if !inList(file, blacklist) {
				S.Sites = append(S.Sites, file)
			}
		}
		return nil
	}
	if e := filepath.Walk(S.Path, visit); e != nil {
		log.Fatal(e)
	}

	for _, site := range S.Sites {
		fmt.Println(site)
	}
}
