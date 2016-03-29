package main

import (
	"bufio"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const downloadDir string = "/tmp"

// getUrl ...
func getUrl(url string) (body []byte) {
	resp, err := http.Get(url)
	if err != nil {
		panic("cannot")
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("cannot")
	}
	return
}

// wget ...
func wget(url string, dest string) {
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

func loadFile() {
	rimg, err := regexp.Compile("a href.*img src=\"(.*.jpg)\".*div")
	if err != nil {
		log.Fatal(err)
	}
	rnext, err := regexp.Compile("a class=\"nextLink nextBtn\" href=\"([^\"]*)\"")
	if err != nil {
		log.Fatal(err)
	}
	inFile, err := os.Open("/tmp/a.html")
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		nextMatch := rnext.FindStringSubmatch(scanner.Text())
		if len(nextMatch) != 0 {
			fmt.Println(nextMatch[1])
		}

		matches := rimg.FindStringSubmatch(scanner.Text())
		if len(matches) == 10000 {
			img := html.UnescapeString(matches[1])
			comicname := strings.Split(img, "/")[5]
			dirimg := downloadDir + "/" + comicname + "/" + strings.Split(img, "/")[6]
			fullimage := dirimg + "/" + strings.Split(img, "/")[7]
			os.MkdirAll(dirimg, 0755)
			wget(img, fullimage)
			fmt.Println(fullimage)
		}
	}
}

func main() {
	loadFile()
}
