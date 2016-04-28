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
	comicDir := flag.String("comicdir", filepath.Join(user.HomeDir, "/Documents/Comics"), "Comics download dir.")
	testMode := flag.Bool("test", false, "Run in a test mode")
	subscribe := flag.String("s", "", "Subscribe to a comic already in DB")
	update := flag.Bool("u", false, "Check if needed update")
	updatePrint := flag.Bool("up", false, "Check update but only print and update DB (cron notification mode)")
	addcomic := flag.Bool("a", false, "Add or Update a comic in database (and subscribe): EPISODE_NAME EPISODE_NUMBER")

	flag.Usage = func() {
		fmt.Printf("Usage: helloyolo [options] hello-comics-url\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if _, err := os.Stat(*comicDir); os.IsNotExist(err) {
		err := os.MkdirAll(*comicDir, 755)
		utils.CheckError(err)
	}
	cfg["comicDir"] = *comicDir
	cfg["testMode"] = strconv.FormatBool(*testMode)

	if *subscribe != "" {
		utils.DBSubscribe(*comicDir, *subscribe)
		fmt.Println(*subscribe, "has been subscribed.")
		os.Exit(0)
	}

	if *addcomic {
		if len(flag.Args()) != 2 {
			fmt.Println("You need to specify two arguments: EPISODE_NAME EPISODE_NUMBER")
			os.Exit(1)
		}
		var episodeNumber int
		var err error
		if episodeNumber, err = strconv.Atoi(flag.Args()[1]); err != nil {
			fmt.Printf("%s does not seem a number\n", flag.Args()[1])
		}
		utils.DBupdate(*comicDir, flag.Args()[0], episodeNumber, true)
		fmt.Printf("%s episode %d has been subscribed/added to the database\n", flag.Args()[0], episodeNumber)
		os.Exit(0)
	}

	if *update {
		hellocomics.GetUpdate(cfg, false)
		os.Exit(0)
	}

	if *updatePrint {
		hellocomics.GetUpdate(cfg, true)
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
