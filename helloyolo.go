package main

import (
	"bufio"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	comicsDir string
	testMode  bool
)

func helloComic(url string) (nextLink, comicname, episode string) {

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

func frComic(url string) (nextLink, comicname, episode string) {
	var img string
	tmpfile, err := ioutil.TempFile("", ".xxxxxxx-download-comics")
	checkError(err)
	defer os.Remove(tmpfile.Name()) // clean up

	wget(url, tmpfile.Name())

	rimg, _ := regexp.Compile("img class=\"open\" src=\"(.*.jpg)\"/>")
	rnext, _ := regexp.Compile("a href=\"([^\"]*)\" onClick=.*nextPage.*")
	repisode, _ := regexp.Compile("^<title>([^:]*)[ ]+::[ ]+Chapitre ([^:]*)[ ]+::.*")

	inFile, err := os.Open(tmpfile.Name())
	checkError(err)
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		episodeMatch := repisode.FindStringSubmatch(scanner.Text())
		if len(episodeMatch) != 0 {
			comicname = episodeMatch[1]
		}

		nextMatch := rnext.FindStringSubmatch(scanner.Text())
		if len(nextMatch) != 0 {
			nextLink = nextMatch[1]
		}

		matches := rimg.FindStringSubmatch(scanner.Text())
		if len(matches) != 0 {
			img = html.UnescapeString(matches[1])
			re := regexp.MustCompile("^[0-9]+")
			number := re.FindString(strings.Split(img, "/")[6])
			re = regexp.MustCompile("_[^_]*$")
			episode = re.ReplaceAllString(strings.Split(img, "/")[5], "${1}") + "-" + number
		}

		if episode != "" && comicname != "" && img != "" && nextLink != "" {
			dirimg := filepath.Join(os.TempDir(), comicname, episode)
			fullimage := filepath.Join(dirimg, strings.Split(img, "/")[7])
			os.MkdirAll(dirimg, 0755)
			if _, err := os.Stat(fullimage); os.IsNotExist(err) {
				log.Printf("IMG: %s\n", img)
				if testMode == false {
					wget(img, fullimage)
				}
			}
			return
		}

	}
	return
}

// loop ...
func loop(url string) {
	var next, comicname, episode string

	if strings.HasPrefix(url, "http://www.hellocomic.com/") {
		next, comicname, episode = helloComic(url)
	} else if strings.HasPrefix(url, "http://fr.comics-reader.com") {
		next, comicname, episode = frComic(url)
	} else {
		log.Fatal("Only hellocomics url for now is supported")
	}

	previousEpisode := episode
	for {
		// retarded
		if strings.HasPrefix(url, "http://www.hellocomic.com/") {
			next, comicname, episode = helloComic(next)
		} else if strings.HasPrefix(url, "http://fr.comics-reader.com") {
			next, comicname, episode = frComic(next)
		}

		if episode != previousEpisode {
			packitShipit(comicsDir, comicname, previousEpisode)
		}
		if strings.HasPrefix(next, "http://www.hellocomic.com/comic/view?slug=") || next == "" {
			packitShipit(comicsDir, comicname, episode)
		}
		previousEpisode = episode
	}
}

func main() {
	user, err := user.Current()
	checkError(err)
	parsed_comicDir := flag.String("comicdir", filepath.Join(user.HomeDir, "/Documents/Comics"), "Comics download dir.")
	parsed_testMode := flag.Bool("test", false, "Run in a test mode")

	flag.Usage = func() {
		fmt.Printf("Usage: helloyolo [options] hello-comics-url\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	comicsDir = *parsed_comicDir
	if _, err := os.Stat(comicsDir); os.IsNotExist(err) {
		os.MkdirAll(comicsDir, 755)
	}
	testMode = *parsed_testMode
	if testMode {
		log.Println("RUNNING IN TEST MODE")
	}

	url := flag.Args()[0]

	loop(url)
}
