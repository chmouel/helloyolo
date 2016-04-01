Download comics from hello-comics and fr.comics-reader.com

You just provide a hellocomic, fr.comics-reader.com URL and it will download them and generate zip
files sorted in ~/Documents/Comics/comic-name/comic-episode.cbz

It uses real wget for downloading since the internal one was very slow for frComic if that's an issue for someone 
i can get the internal one as fallback again to be truly os independent.

Future feature would hopefully will include tracking and updating for now it just store in a sqlitedb file in  ~/Documents/Comics/.helloyolo.db  the latest episode number of the serie. 

INSTALL
-------

Just do a :

$ go get -u github.com/chmouel/helloyolo

and the binary should go to GOPATH/bin, just rerun it to get it updated.

you can as well just use the Makefile which would output to _output/ directory
