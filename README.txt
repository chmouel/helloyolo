Download comics from hello-comics and fr.comics-reader.com

You just provide a hellocomic, fr.comics-reader.com URL and it will download
them and generate zip
files sorted in ~/Documents/Comics/comic-name/comic-episode.cbz

It uses real curl for downloading since the internal one was very slow for
frComic if that's an issue for someone
i can get the internal one as fallback again to be truly os independent.

With the '-u' option you can check for updates (which is only supported with
hellocomic backend) it checks the subscribed comics from your database on
~/Documents/Comics/.helloyolo.db as long it has the subscribed flags. You can
manage the database manually with sqlite3 command line or just use the '-s'
option to subscribe to one.

INSTALL
-------

Install go and configure your environment variable GOPATH to something like ~/go
and add the $GOPATH/bin to your PATH in your shell environement file.

Then just do a :

$ go get -u github.com/chmouel/helloyolo

and the binary should go to $GOPATH/bin, rerun that same command to get it
updated.

you can as well just use the Makefile which would output to _output/ directory
and grab the generated `helloyolo` binary there
