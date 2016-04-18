package utils

import (
	"database/sql"
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
func DBupdate(comicsDir, episode string, latest int) {
	// TODO(chmou): Make it using the global flags
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
