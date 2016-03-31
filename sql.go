package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var sqlTable = `
CREATE TABLE IF NOT EXISTS Comics (
	id integer PRIMARY KEY,
	ComicName varchar(255) NOT NULL,
	Last integer,
	CONSTRAINT uc_comicID UNIQUE (ComicName, Last))`

// update
func DBupdate(episode string, latest int) {
	db, err := sql.Open("sqlite3", filepath.Join(comicsDir, ".helloyolo.db")) // TODO(chmouel):
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(sqlTable)
	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec()

	stmt, err = db.Prepare("INSERT OR REPLACE INTO Comics(comicname, last) values(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(episode, latest)

	if err != nil {
		log.Fatal(err)
	}
}

func DBgetLatestEpisode(episode string) {
	db, err := sql.Open("sqlite3", "./foo.db") // TODO(chmouel):
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT last FROM Comics where comicname=?")
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
