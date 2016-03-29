package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

const downloadDir string = "/tmp"

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		log.Fatal(err)
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Fatal(err)
		}

		if baseDir != "" {
			header.Name = strings.TrimPrefix(path, source)
		}

		if info.IsDir() {
			return nil
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func wget(url, dest string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}

	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func loadNextFromFile(url string) (nextLink, comicname, episode string) {

	tmpfile, err := ioutil.TempFile("", ".xxxxxxx-download-comics")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	wget(url, tmpfile.Name())

	rimg, err := regexp.Compile("a href.*img src=\"(.*.jpg)\".*div")
	if err != nil {
		log.Fatal(err)
	}
	rnext, err := regexp.Compile("a class=\"nextLink nextBtn\" href=\"([^\"]*)\"")
	if err != nil {
		log.Fatal(err)
	}
	inFile, err := os.Open(tmpfile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		nextMatch := rnext.FindStringSubmatch(scanner.Text())
		if len(nextMatch) != 0 {
			nextLink = nextMatch[1]
		}

		matches := rimg.FindStringSubmatch(scanner.Text())
		if len(matches) != 0 {
			img := html.UnescapeString(matches[1])
			comicname = strings.Split(img, "/")[5]
			episode = strings.Split(img, "/")[6]
			dirimg := filepath.Join(downloadDir, comicname, episode)
			fullimage := filepath.Join(dirimg, strings.Split(img, "/")[7])
			os.MkdirAll(dirimg, 0755)
			if _, err := os.Stat(fullimage); err == nil {
				//log.Printf("Skiping %s\n", fullimage)
			} else {
				log.Printf("IMG: %s\n", img)
				wget(img, fullimage)
			}
		}
	}
	return
}

func packitShipit(comicsDir, comicname, episode string) {
	cbzDir := filepath.Join(comicsDir, comicname)
	os.MkdirAll(cbzDir, 0755)
	cbzFile := fmt.Sprintf("%s/%s.cbz", cbzDir, episode)
	tmpDir := filepath.Join(downloadDir, comicname, episode)
	if _, err := os.Stat(cbzFile); err != nil {
		if os.IsNotExist(err) {
			zipit(tmpDir, cbzFile)
		}
	}
	os.RemoveAll(tmpDir)
	log.Printf("ZIP: %s\n", cbzFile)
}

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	comicsDir := filepath.Join(user.HomeDir, "/Documents/Comics")

	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal("Usage: comics-download hello-comics-url")
	}

	url := flag.Args()[0]
	if strings.HasPrefix(url, "http://www.hellocomic.com/") == false {
		log.Fatal("Only hellocomics url for now is supported")
	}

	var next, comicname, episode = loadNextFromFile(url)
	var previousEpisode = episode
	for {
		next, comicname, episode = loadNextFromFile(next)
		if episode != previousEpisode {
			packitShipit(comicsDir, comicname, previousEpisode)
		}
		if strings.HasPrefix(next, "http://www.hellocomic.com/comic/view?slug=") || next == "" {
			packitShipit(comicsDir, comicname, episode)
			log.Println("Finished")
			break
		}
		previousEpisode = episode
	}
}
