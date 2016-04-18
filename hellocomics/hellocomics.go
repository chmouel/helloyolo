package hellocomics

import (
	"bufio"
	"fmt"
	"html"
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

func pack(comicname, episode string) {
	r, err := regexp.Compile("\\d+$")
	utils.CheckError(err)
	match := r.FindString(episode)
	if match != "" {
		episodeNumber, err := strconv.Atoi(match)
		utils.CheckError(err)
		utils.DBupdate(comicname, episodeNumber)
	} else {
		log.Println("Cannot figure out the episode number? not updating DB: " + episode)
	}

	cbzDir := filepath.Join(config["comicDir"], comicname)
	os.MkdirAll(cbzDir, 0755)
	cbzFile := fmt.Sprintf("%s/%s.cbz", cbzDir, episode)
	tmpDir := filepath.Join(os.TempDir(), comicname, episode)
	if _, err := os.Stat(cbzFile); os.IsNotExist(err) {
		utils.Zipit(tmpDir, cbzFile)
		log.Printf("ZIP: %s\n", cbzFile)
	} else {
		log.Printf("ZIP: Skipping %s\n", cbzFile)

	}
	//os.RemoveAll(tmpDir)
}

func helloParse(url string) (nextLink, comicname, episode string) {
	tmpfile, err := ioutil.TempFile("", ".xxxxxxx-download-comics")
	utils.CheckError(err)
	defer os.Remove(tmpfile.Name()) // clean up

	utils.Wget(url, tmpfile.Name())

	rimg, err := regexp.Compile("a href.*img src=\"(.*.jpg)\".*div")
	utils.CheckError(err)
	rnext, err := regexp.Compile("a class=\"nextLink nextBtn\" href=\"([^\"]*)\"")
	utils.CheckError(err)
	inFile, err := os.Open(tmpfile.Name())
	utils.CheckError(err)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		if nextMatch := rnext.FindStringSubmatch(scanner.Text()); len(nextMatch) != 0 {
			nextLink = nextMatch[1]
		}

		matches := rimg.FindStringSubmatch(scanner.Text())
		if len(matches) != 0 {
			img := html.UnescapeString(matches[1])
			comicname = strings.Split(img, "/")[5]
			episode = strings.Split(img, "/")[6]
			dirimg := filepath.Join(os.TempDir(), comicname, episode)
			fullimage := filepath.Join(dirimg, strings.Split(img, "/")[7])
			os.MkdirAll(dirimg, 0755)
			if _, err := os.Stat(fullimage); os.IsNotExist(err) {
				log.Printf("IMG: %s\n", img)
				utils.Wget(img, fullimage)
			}
		}
	}
	return
}

//HelloComics stuff
func HelloComics(cfg map[string]string) {
	var next, comicname, episode string
	config = cfg

	next, comicname, episode = helloParse(config["url"])

	previousEpisode := episode
	for {
		next, comicname, episode = helloParse(next)

		if strings.HasPrefix(next, "http://www.hellocomic.com/comic/view?slug=") || next == "" {
			pack(comicname, episode)
			break
		}
		if episode != previousEpisode {
			pack(comicname, previousEpisode)
		}
		previousEpisode = episode
	}
}
