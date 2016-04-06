package main

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
)

func helloPack(comicsDir, comicname, episode string) {
	r, err := regexp.Compile("\\d+$")
	checkError(err)
	match := r.FindString(episode)
	if match == "" {
		log.Fatal("Cannot figure out the episode number?")
	}
	episodeNumber, err := strconv.Atoi(match)
	checkError(err)

	DBupdate(episode, episodeNumber)

	cbzDir := filepath.Join(comicsDir, comicname)
	os.MkdirAll(cbzDir, 0755)
	cbzFile := fmt.Sprintf("%s/%s.cbz", cbzDir, episode)
	tmpDir := filepath.Join(os.TempDir(), comicname, episode)
	if _, err := os.Stat(cbzFile); os.IsNotExist(err) {
		if testMode == false {
			zipit(tmpDir, cbzFile)
		}
		log.Printf("ZIP: %s\n", cbzFile)
	} else {
		log.Printf("ZIP: Skipping %s\n", cbzFile)

	}
	//os.RemoveAll(tmpDir)
}

func helloParse(url string) (nextLink, comicname, episode string) {
	tmpfile, err := ioutil.TempFile("", ".xxxxxxx-download-comics")
	checkError(err)
	defer os.Remove(tmpfile.Name()) // clean up

	wget(url, tmpfile.Name())

	rimg, err := regexp.Compile("a href.*img src=\"(.*.jpg)\".*div")
	checkError(err)
	rnext, err := regexp.Compile("a class=\"nextLink nextBtn\" href=\"([^\"]*)\"")
	checkError(err)
	inFile, err := os.Open(tmpfile.Name())
	checkError(err)
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
				if testMode == false {
					wget(img, fullimage)
				}
			}
		}
	}
	return
}

//HelloComics stuff
func HelloComics(url string) {
	fmt.Println("hello")
	var next, comicname, episode string

	next, comicname, episode = helloParse(url)

	previousEpisode := episode
	for {
		next, comicname, episode = helloParse(next)

		if episode != previousEpisode {
			helloPack(comicsDir, comicname, previousEpisode)
		}
		if strings.HasPrefix(next, "http://www.hellocomic.com/comic/view?slug=") || next == "" {
			helloPack(comicsDir, comicname, episode)
		}
		previousEpisode = episode
	}
}