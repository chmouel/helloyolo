package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func packitShipit(comicsDir, comicname, episode string) {
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

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	checkError(err)

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		checkError(err)

		header, err := zip.FileInfoHeader(info)
		checkError(err)

		if baseDir != "" {
			header.Name = strings.TrimPrefix(path, source)
		}

		if info.IsDir() {
			return nil
		}
		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		checkError(err)

		file, err := os.Open(path)
		checkError(err)

		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}
