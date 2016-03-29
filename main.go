package main

import (
	"bufio"
	"flag"
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
	log.Printf("Getting %s to %s\n", url, dest)
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func loadNextFromFile(url string) (nextLink string) {

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
			comicname := strings.Split(img, "/")[5]
			dirimg := downloadDir + "/" + comicname + "/" + strings.Split(img, "/")[6]
			fullimage := dirimg + "/" + strings.Split(img, "/")[7]
			os.MkdirAll(dirimg, 0755)
			if _, err := os.Stat(fullimage); err == nil {
				log.Printf("Skiping %s\n", fullimage)
			} else {
				wget(img, fullimage)
			}
		}
	}
	return
}

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("Usage: comics-download hello-comics-url")
	}

	url := flag.Args()[0]

	if strings.HasPrefix(url, "http://www.hellocomic.com/") == false {
		log.Fatal("Only hellocomics url for now is supported")
	}

	var next = loadNextFromFile(url)
	for {
		if strings.HasPrefix(next, "http://www.hellocomic.com/comic/view?slug=") || next == "" {
			log.Println("Finished")
			break
		}
		next = loadNextFromFile(next)
	}
}
