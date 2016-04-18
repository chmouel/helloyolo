package utils

import (
	"database/sql"
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
