package amclient

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

func initDB(path string) error {
	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite3", path)
		if err != nil {
			return
		}

		createTables := []string{
			`
			CREATE TABLE IF NOT EXISTS artist_album (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				artist TEXT NOT NULL,
				album TEXT NOT NULL,
				UNIQUE(artist, album)
			);`,
			`
			CREATE TABLE IF NOT EXISTS artwork_urls (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				artist_album_id INTEGER NOT NULL,
				freq INTEGER DEFAULT 0,
				url TEXT,
				FOREIGN KEY (artist_album_id) REFERENCES artist_album(id),
				UNIQUE(url)
			);`,
		}

		for _, table := range createTables {
			if _, err = db.Exec(table); err != nil {
				return
			}
		}
	})

	return err
}

func createDbEntry(artist string, album string, url string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert or ignore artist_album (due to UNIQUE constraint)
	res, err := tx.Exec(`INSERT OR IGNORE INTO artist_album (artist, album) VALUES (?, ?)`, artist, album)
	if err != nil {
		return err
	}

	var artistAlbumID int64
	if lastID, err := res.LastInsertId(); err == nil && lastID != 0 {
		artistAlbumID = lastID
	} else {
		err = tx.QueryRow(`SELECT id FROM artist_album WHERE artist = ? AND album = ?`, artist, album).Scan(&artistAlbumID)
		if err != nil {
			return err
		}
	}

	if url != "" {
		_, err = tx.Exec(`INSERT OR IGNORE INTO artwork_urls (artist_album_id, freq, url) VALUES (?, ?, ?)`,
			artistAlbumID, 1, url)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func getMaxFreqUrl(artist string, album string) (string, error) {
	var artistAlbumID int
	err := db.QueryRow(`SELECT id FROM artist_album WHERE artist = ? AND album = ?`, artist, album).Scan(&artistAlbumID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("artist+album not found")
		}
		return "", err
	}

	var url string
	err = db.QueryRow(`
		SELECT url FROM artwork_urls
		WHERE artist_album_id = ?
		ORDER BY freq DESC
		LIMIT 1
	`, artistAlbumID).Scan(&url)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no artwork URLs found for artist+album")
		}
		return "", err
	}

	return url, nil
}

// CloseDB safely closes the database
func closeDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
		db = nil
	}
}
