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

	"github.com/PuerkitoBio/goquery"
	"github.com/chmouel/helloyolo/utils"
)

const updatepagetofetch int = 3

var config = make(map[string]string)

// GetUpdatePageNumber get update of page number
func GetUpdatePageNumber(updateMode bool, page int) {
	tmpfile, err := ioutil.TempFile("", ".xxxxxxx-download-comics")
	utils.CheckError(err)
	defer os.Remove(tmpfile.Name()) // clean up
	utils.Wget("http://www.hellocomic.com/site/index?page="+strconv.Itoa(page), tmpfile.Name())

	r, err := os.Open(tmpfile.Name())
	utils.CheckError(err)

	doc, err := goquery.NewDocumentFromReader(r)
	utils.CheckError(err)

	doc.Find("dd").Each(func(i int, s *goquery.Selection) {
		val, exist := s.Children().Attr("href")
		if exist {
			// TODO(chmouel): convert to a regexp
			trimmed := strings.TrimPrefix(val, "https://www.hellocomic.com/")
			trimmed = strings.TrimPrefix(trimmed, "http://www.hellocomic.com/")
			splits := strings.Split(trimmed, "/")
			episodeNumber, err := strconv.Atoi(strings.TrimPrefix(splits[1], "c"))
			utils.CheckError(err)
			comicname := splits[0]

			needupdate := utils.DBCheckLatest(config["comicDir"], comicname, episodeNumber)
			if needupdate && updateMode {
				fmt.Printf("%s-%d -- %s\n", comicname, episodeNumber, val)
				utils.DBupdate(config["comicDir"], comicname, episodeNumber)
			} else if needupdate {
				fmt.Println("Updating", comicname, episodeNumber)
				config["url"] = val
				HelloComics(config)
			}
		}
	})
}

// GetUpdate get all update of the comics
func GetUpdate(cfg map[string]string, updateMode bool) {
	config = cfg
	for i := 1; i < updatepagetofetch+1; i++ {
		GetUpdatePageNumber(updateMode, i)
	}
}

func pack(comicname, episode string) {
	r, err := regexp.Compile("\\d+$")
	utils.CheckError(err)
	match := r.FindString(episode)
	if match != "" {
		episodeNumber, err := strconv.Atoi(match)
		utils.CheckError(err)
		utils.DBupdate(config["comicDir"], comicname, episodeNumber)
	} else {
		log.Println("Cannot figure out the episode number? not updating DB: " + episode)
	}

	cbzDir := filepath.Join(config["comicDir"], comicname)
	err = os.MkdirAll(cbzDir, 0755)
	utils.CheckError(err)

	cbzFile := fmt.Sprintf("%s/%s.cbz", cbzDir, episode)
	tmpDir := filepath.Join(os.TempDir(), comicname, episode)
	if _, err := os.Stat(cbzFile); os.IsNotExist(err) {
		err = utils.Zipit(tmpDir, cbzFile)
		utils.CheckError(err)
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
