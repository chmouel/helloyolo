package main

import (
	"log"
	"os/exec"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func wget(url, dest string) {
	wgetExec, err := exec.LookPath("wget")
	checkError(err)

	_, err = exec.Command(wgetExec, "-c", "-O", dest, url).Output()
	checkError(err)
}
