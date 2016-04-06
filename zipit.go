package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
