package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os/user"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var sqlTable = `
CREATE TABLE IF NOT EXISTS Comics (
	id integer PRIMARY KEY,
	ComicName varchar(255) NOT NULL,
	LastEpisode integer,
    Subscribed BOOLEAN NOT NULL DEFAULT 0,
    DateTime datetime default current_timestamp,
    CHECK (Subscribed IN (0,1)),
	CONSTRAINT uc_comicID UNIQUE (ComicName))`

// DBupdate Update DB
func DBupdate(episode string, latest int) {
	user, err := user.Current()
	// TODO(chmou): Make it using the global flags
	comicsDir := filepath.Join(user.HomeDir, "/Documents/Comics")
	db, err := sql.Open("sqlite3", filepath.Join(comicsDir, ".helloyolo.db"))
	CheckError(err)

	defer db.Close()

	stmt, err := db.Prepare(sqlTable)
	CheckError(err)

	_, err = stmt.Exec()
	CheckError(err)

	stmt, err = db.Prepare("INSERT OR REPLACE INTO Comics(ComicName, LastEpisode) values(?,?)")
	CheckError(err)

	_, err = stmt.Exec(episode, latest)
	CheckError(err)
}

// dbgetLatestEpisode we will use this!
func dbgetLatestEpisode(episode string) {
	db, err := sql.Open("sqlite3", "./foo.db") // TODO(chmouel):
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT LastEpisode FROM Comics where comicname=?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Query(episode)
	if err != nil {
		log.Fatal(err)
	}
	for res.Next() {
		var last int
		err = res.Scan(&last)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(last)
	}

}
