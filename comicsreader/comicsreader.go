package comicsreader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

var comicPrefixURL string

// Pages repr
type Pages struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

func parse(nextLink string) (nextURL string) {
	var comicName, fileName, episodeNumber string
	var allImages []Pages

	rComicName, err := regexp.Compile("<h1 class=\"tbtitle dnone\"><a href=\"[^\"]*\" title=\"([^\"]*)\">.* :: ")
	utils.CheckError(err)
	rCurrentEpisode, err := regexp.Compile("<a href=\"([^\"]*)\" onClick=\"return nextPage()")
	utils.CheckError(err)
	rAllPages, err := regexp.Compile("var pages = (.*);$")
	utils.CheckError(err)

	if _, err := os.Stat(nextLink); os.IsNotExist(err) {
		fileName = getURL(nextLink)
		defer os.Remove(fileName) // clean up
	} else {
		//Mostly for testing
		fileName = nextLink
	}

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReaderSize(f, 4*bufio.MaxScanTokenSize)
	line, isPrefix, scanErr := r.ReadLine()
	for scanErr == nil && !isPrefix {
		tmps := strings.TrimSpace(string(line))

		comicNameMatch := rComicName.FindStringSubmatch(tmps)
		if len(comicNameMatch) != 0 {
			comicName = comicNameMatch[1]
		}

		currentEpisode := rCurrentEpisode.FindStringSubmatch(tmps)
		if len(currentEpisode) != 0 {
			episodeNumber = getEpisode(currentEpisode[1])
		}

		jsAllPages := rAllPages.FindStringSubmatch(tmps)
		if len(jsAllPages) != 0 {
			tmps = jsAllPages[1]
			err := json.Unmarshal([]byte(tmps), &allImages)
			utils.CheckError(err)
		}

		if strings.HasPrefix(tmps, "var next_chapter = \"") {
			tmps = strings.TrimPrefix(tmps, "var next_chapter = \"")
			tmps = strings.TrimSuffix(tmps, "/\";")
			if len(strings.Split(strings.TrimPrefix(tmps, comicPrefixURL), "/")) > 1 {
				nextURL = tmps
			}
		}
		line, isPrefix, scanErr = r.ReadLine()
	}

	if isPrefix {
		fmt.Println("buffer size to small")
		return
	}
	if scanErr != io.EOF {
		log.Fatal(err)
	}

	if comicName == "" {
		log.Fatal("I didn't get the comicName which is weird\nMake sure you have the page with the image something like http://fr.comics-reader.com/read/batman__new_52_fr/fr/1/0/")
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
			log.Printf(v.Filename)
			utils.Wget(v.URL, temporaryFileImage)
		}
	}

	targetDir := filepath.Join(config["comicDir"], comicName)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		os.MkdirAll(targetDir, 0755)
	}
	targetCBZ := fmt.Sprintf("%s/%s.cbz", targetDir, episodeNumber)
	if _, err := os.Stat(targetCBZ); os.IsNotExist(err) {
		err := utils.Zipit(dirImg, targetCBZ)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("ZIP: %s ", targetCBZ)
	}

	reg, err := regexp.Compile("(.*)-(\\d+)$")
	utils.CheckError(err)
	match := reg.FindStringSubmatch(episodeNumber)
	if len(match) != 3 {
		log.Fatal("Cannot figure out the episode number?")
	}

	numero, err := strconv.Atoi(match[2])
	utils.CheckError(err)

	utils.DBupdate(config["comicDir"], match[1], numero, false)
	return nextURL
}

// getTMPFile ...
func getURL(url string) string {
	tmpfile, err := ioutil.TempFile("", ".xxxxxxx-download-comics")
	utils.CheckError(err)
	utils.Curl(url, tmpfile.Name(), "-F adult=true")
	return tmpfile.Name()
}

func getEpisode(s string) (ret string) {
	s = strings.TrimSuffix(s, "/")
	sp := strings.Split(strings.TrimPrefix(s, comicPrefixURL), "/")

	// Pretty clumsy but I haven't find a way to this in a better way
	if sp[2] != "0" && sp[3] == "0" {
		ret = fmt.Sprintf("%s-%s", sp[0], sp[2])
	} else {
		ret = fmt.Sprintf("%s-%s", sp[0], sp[len(sp)-3])
	}
	return
}

// Loop over all the links until the is none
func Loop(cfg map[string]string) {
	reg, err := regexp.Compile("^(http://(fr|us)\\.comics-reader\\.com/read/)")
	utils.CheckError(err)
	match := reg.FindStringSubmatch(cfg["url"])
	comicPrefixURL = match[0]

	config = cfg
	nextLink := cfg["url"]
	for {
		nextLink = parse(nextLink)
		if nextLink == "" {
			break
		}
	}
}
