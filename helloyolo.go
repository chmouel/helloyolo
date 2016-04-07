package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/chmouel/helloyolo/frcomics"
	"github.com/chmouel/helloyolo/utils"
)

var (
	comicsDir string
	testMode  bool
)

func main() {
	user, err := user.Current()
	utils.CheckError(err)
	parsedComicDir := flag.String("comicdir", filepath.Join(user.HomeDir, "/Documents/Comics"), "Comics download dir.")
	parsedTestMode := flag.Bool("test", false, "Run in a test mode")

	flag.Usage = func() {
		fmt.Printf("Usage: helloyolo [options] hello-comics-url\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	comicsDir = *parsedComicDir
	if _, err := os.Stat(comicsDir); os.IsNotExist(err) {
		os.MkdirAll(comicsDir, 755)
	}
	testMode = *parsedTestMode
	if testMode {
		log.Println("RUNNING IN TEST MODE")
	}

	url := flag.Args()[0]

	if strings.HasPrefix(url, "http://fr.comics-reader.com/") {
		frcomics.Loop(url)
	} else if strings.HasPrefix(url, "http://www.hellocomic.com/") {
		HelloComics(url)
	}
}
