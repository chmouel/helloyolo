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

//Wget get url to a dest locally with wget directly, the builtin one of getting
//via go was for whatever reason too slow
func Wget(url, dest string) {
	wgetExec, err := exec.LookPath("wget")
	CheckError(err)

	_, err = exec.Command(wgetExec, "-c", "-O", dest, url).Output()
	CheckError(err)
}
