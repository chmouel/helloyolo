package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chmouel/helloyolo/frcomics"
	"github.com/chmouel/helloyolo/hellocomics"
	"github.com/chmouel/helloyolo/utils"
)

func main() {
	var cfg = make(map[string]string)

	user, err := user.Current()
	utils.CheckError(err)
	parsedComicDir := flag.String("comicdir", filepath.Join(user.HomeDir, "/Documents/Comics"), "Comics download dir.")
	parsedTestMode := flag.Bool("test", false, "Run in a test mode")
	update := flag.Bool("u", false, "Check if needed update")

	flag.Usage = func() {
		fmt.Printf("Usage: helloyolo [options] hello-comics-url\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if _, err := os.Stat(*parsedComicDir); os.IsNotExist(err) {
		err := os.MkdirAll(*parsedComicDir, 755)
		utils.CheckError(err)
	}
	cfg["comicDir"] = *parsedComicDir
	cfg["testMode"] = strconv.FormatBool(*parsedTestMode)

	if *update {
		hellocomics.GetUpdate(cfg)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	url := flag.Args()[0]
	cfg["url"] = url
	if strings.HasPrefix(url, "http://fr.comics-reader.com/") {
		frcomics.Loop(cfg)
	} else if strings.HasPrefix(url, "http://www.hellocomic.com/") {
		hellocomics.HelloComics(cfg)
	}

}
