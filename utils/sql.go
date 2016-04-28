package utils

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var sqlTable = `
CREATE TABLE IF NOT EXISTS Comics (
	id integer PRIMARY KEY,
	ComicName varchar(255) NOT NULL,
	LastEpisode integer,
    DateTime datetime default current_timestamp,
	CONSTRAINT uc_comicID UNIQUE (ComicName))`

func createTable(comicsDir string) (db *sql.DB) {
	db, err := sql.Open("sqlite3", filepath.Join(comicsDir, ".helloyolo.db"))
	CheckError(err)

	stmt, err := db.Prepare(sqlTable)
	CheckError(err)

	_, err = stmt.Exec()
	CheckError(err)

	return db
}

// DBupdate Update DB
func DBupdate(comicsDir, comicname string, latest int, subscribe bool) {
	// TODO(chmou): Make it using the global flags
	db := createTable(comicsDir)
	defer db.Close()

	var subscribesql string
	if subscribe {
		subscribesql = "1"
	} else {
		// Autodetect the previous subscribtion
		subscribesql = fmt.Sprintf(`(SELECT subscribed FROM Comics WHERE comicname='%s')`, comicname)
	}

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO Comics(ComicName, LastEpisode, Subscribed) VALUES(?,?,%s)`, subscribesql)
	stmt, err := db.Prepare(sql)
	CheckError(err)

	_, err = stmt.Exec(comicname, latest)
	CheckError(err)
}

// DBCheckLatest episode
func DBCheckLatest(comicsDir, comicsname string, latest int) bool {
	var needUpdate int
	db := createTable(comicsDir)
	defer db.Close()

	_ = db.QueryRow("select 1 from comics where comicName = ? and subscribed=1 and lastepisode < ?", comicsname, latest).Scan(&needUpdate)

	if needUpdate == 1 {
		return true
	}

	return false
}

// DBSubscribe episode
func DBSubscribe(comicsDir, comicsname string) {
	db := createTable(comicsDir)
	defer db.Close()

	stmt, err := db.Prepare("UPDATE Comics SET subscribed=1 WHERE comicname=?")
	CheckError(err)

	_, err = stmt.Exec(comicsname)
	CheckError(err)
}
