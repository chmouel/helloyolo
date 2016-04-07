package frcomics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/chmouel/helloyolo/utils"
)

var config = make(map[string]string)

const frComicPrefixURL string = "http://fr.comics-reader.com/read/"

// Pages repr
type Pages struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

func parse(nextLink string) (nextURL string) {
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
			utils.CheckError(err)
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
			utils.Wget(v.URL, temporaryFileImage)
		}
	}

	targetDir := filepath.Join(config["comicdir"], "/Documents/Comics", comicName)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		os.MkdirAll(targetDir, 0755)
	}
	targetCBZ := fmt.Sprintf("%s/%s.cbz", targetDir, episodeNumber)
	if _, err := os.Stat(targetCBZ); os.IsNotExist(err) {
		utils.Zipit(dirImg, targetCBZ)
		log.Printf("ZIP: %s ", targetCBZ)
	}

	reg, err := regexp.Compile("\\d+$")
	utils.CheckError(err)
	match := reg.FindString(episodeNumber)
	if match == "" {
		log.Fatal("Cannot figure out the episode number?")
	}
	numero, err := strconv.Atoi(match)
	utils.CheckError(err)

	utils.DBupdate(episodeNumber, numero)
	return nextURL
}

// getTMPFile ...
func getURL(url string) string {
	tmpfile, err := ioutil.TempFile("", ".xxxxxxx-download-comics")
	utils.CheckError(err)
	utils.Wget(url, tmpfile.Name())
	return tmpfile.Name()
}

func getEpisode(s string) string {
	s = strings.TrimSuffix(s, "/")
	sp := strings.Split(strings.TrimPrefix(s, frComicPrefixURL), "/")
	return fmt.Sprintf("%s-%s", sp[0], sp[len(sp)-1])
}

// Loop over all the links until the is none
func Loop(cfg map[string]string) {
	config = cfg
	for {
		nextLink := parse(cfg["url"])
		if nextLink == "" {
			break
		}
	}
}
