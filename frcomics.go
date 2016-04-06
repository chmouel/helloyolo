package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const frComicPrefixURL string = "http://fr.comics-reader.com/read/"

// Pages repr
type Pages struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

func frcomicsParse(nextLink string) (nextURL string) {
	var comicName, fileName, episodeNumber string
	var allImages []Pages

	rComicName, _ := regexp.Compile("<h1 class=\"tbtitle dnone\"><a href=\"[^\"]*\" title=\"([^\"]*)\">.* :: ")

	if _, err := os.Stat(nextLink); os.IsNotExist(err) {
		fileName = getURL(nextLink)
		defer os.Remove(fileName) // clean up
	} else {
		//Mostly for testing
		fileName = nextLink
	}

	r, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		tmps := strings.TrimSpace(scanner.Text())

		comicNameMatch := rComicName.FindStringSubmatch(tmps)
		if len(comicNameMatch) != 0 {
			comicName = comicNameMatch[1]
		}

		if strings.HasPrefix(tmps, "var pages = [{") {
			tmps = strings.TrimPrefix(tmps, "var pages = ")
			tmps = strings.TrimSuffix(tmps, ";")
			err := json.Unmarshal([]byte(tmps), &allImages)
			checkError(err)
		}

		if strings.HasPrefix(tmps, "var next_chapter = \"") {
			tmps = strings.TrimPrefix(tmps, "var next_chapter = \"")
			tmps = strings.TrimSuffix(tmps, "/\";")
			if len(strings.Split(strings.TrimPrefix(tmps, frComicPrefixURL), "/")) > 1 {
				nextURL = tmps
			}
		}

		if strings.HasPrefix(tmps, "var base_url =") {
			tmps = strings.TrimPrefix(tmps, "var base_url = '")
			tmps = strings.TrimSuffix(tmps, "';")
			episodeNumber = getEpisode(tmps)
		}
	}

	if comicName == "" {
		log.Fatal("I didn't get the comicName which is weird")
	}
	if episodeNumber == "" {
		log.Fatal("I didn't get the episodeNumber which is weird")
	}

	if len(allImages) == 0 {
		log.Fatal("I didn't get the any images to download which is weird")
	}

	log.Printf("Downloading: %s ", episodeNumber)
	dirImg := filepath.Join(os.TempDir(), comicName, episodeNumber)
	if _, err := os.Stat(dirImg); os.IsNotExist(err) {
		os.MkdirAll(dirImg, 0755)
	}
	for _, v := range allImages {
		temporaryFileImage := filepath.Join(dirImg, v.Filename)
		if _, err := os.Stat(temporaryFileImage); os.IsNotExist(err) {
			log.Printf("%s ", v.Filename)
			wget(v.URL, temporaryFileImage)
		}
	}

	user, err := user.Current()
	targetDir := filepath.Join(user.HomeDir, "/Documents/Comics", comicName)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		os.MkdirAll(targetDir, 0755)
	}
	targetCBZ := fmt.Sprintf("%s/%s.cbz", targetDir, episodeNumber)
	if _, err := os.Stat(targetCBZ); os.IsNotExist(err) {
		zipit(dirImg, targetCBZ)
		log.Printf("ZIP: %s ", targetCBZ)
	}

	reg, err := regexp.Compile("\\d+$")
	checkError(err)
	match := reg.FindString(episodeNumber)
	if match == "" {
		log.Fatal("Cannot figure out the episode number?")
	}
	numero, err := strconv.Atoi(match)
	checkError(err)

	DBupdate(episodeNumber, numero)
	return nextURL
}

// getTMPFile ...
func getURL(url string) string {
	tmpfile, err := ioutil.TempFile("", ".xxxxxxx-download-comics")
	checkError(err)
	wget(url, tmpfile.Name())
	return tmpfile.Name()
}

func getEpisode(s string) string {
	s = strings.TrimSuffix(s, "/")
	sp := strings.Split(strings.TrimPrefix(s, frComicPrefixURL), "/")
	return fmt.Sprintf("%s-%s", sp[0], sp[len(sp)-1])
}

// FRComics loop over all the links
func FRComics(nextLink string) {
	for {
		nextLink = frcomicsParse(nextLink)
		if nextLink == "" {
			break
		}
	}
}
