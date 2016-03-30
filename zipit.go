package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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
