package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func wget(url, dest string) {
	response, err := http.Get(url)
	checkError(err)

	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(dest)
	checkError(err)

	_, err = io.Copy(file, response.Body)
	checkError(err)
	file.Close()
}
