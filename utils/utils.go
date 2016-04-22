package utils

import (
	"log"
	"os/exec"
)

// CheckError is just a dummy function to check those common pattern in golang
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//Curl get url to a dest locally with curl directly, the builtin one of getting
//via go was for whatever reason too slow, `optionals` would pass optionals args
func Curl(url, dest, optionals string) {
	curlExec, err := exec.LookPath("curl")
	CheckError(err)

	_, err = exec.Command(curlExec, "-o", dest, optionals, url).Output()
	CheckError(err)
}
