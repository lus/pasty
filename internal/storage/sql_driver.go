package storage

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/lus/pasty/internal/env"
	"github.com/lus/pasty/internal/pastes"
	_ "github.com/mattn/go-sqlite3"
)

// SQLDriver represents the SQL storage driver
type SQLDriver struct {
	database *sql.DB
	table    string
}

// Initialize initializes the SQL storage driver
func (driver *SQLDriver) Initialize() error {
	// Parse the DSN and create a database object
	db, err := sql.Open(env.Get("STORAGE_SQL_DRIVER", "sqlite3"), env.Get("STORAGE_SQL_DSN", "./db"))
	if err != nil {
		return err
	}

	// Ping the database
	err = db.Ping()
	if err != nil {
		return err
	}

	// Migrate the database
	table := env.Get("STORAGE_SQL_TABLE", "pasty")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS ? (
			id varchar NOT NULL PRIMARY KEY,
			content varchar NOT NULL,
			suggestedSyntaxType varchar NOT NULL,
			deletionToken varchar NOT NULL,
			created bigint NOT NULL,
			autoDelete bool NOT NULL
		);
    `, table)
	if err != nil {
		return err
	}

	// Set the database object and table name of the SQL driver
	driver.database = db
	driver.table = table
	return nil
}

// Terminate terminates the SQL storage driver
func (driver *SQLDriver) Terminate() error {
	return driver.database.Close()
}

// ListIDs returns a list of all existing paste IDs
func (driver *SQLDriver) ListIDs() ([]string, error) {
	// Execute a SELECT query to retrieve all the paste IDs
	rows, err := driver.database.Query("SELECT id FROM ?", driver.table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the rows into a slice of IDs and return it
	var ids []string
	err = rows.Scan(&ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// Get loads a paste
func (driver *SQLDriver) Get(id string) (*pastes.Paste, error) {
	// Execute a SELECT query to retrieve the paste
	row := driver.database.QueryRow("SELECT * FROM ? WHERE id = ?", driver.table, id)
	err := row.Err()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Scan the row into a paste and return it
	paste := new(pastes.Paste)
	err = row.Scan(&paste)
	if err != nil {
		return nil, err
	}
	return paste, nil
}

// Save saves a paste
func (driver *SQLDriver) Save(paste *pastes.Paste) error {
	// Execute an INSERT statement to create the paste
	_, err := driver.database.Exec("INSERT INTO ? (?, ?, ?, ?, ?, ?)", driver.table, paste.ID, paste.Content, paste.SuggestedSyntaxType, paste.DeletionToken, paste.Created, paste.AutoDelete)
	return err
}

// Delete deletes a paste
func (driver *SQLDriver) Delete(id string) error {
	// Execute a DELETE statement to delete the paste
	_, err := driver.database.Exec("DELETE FROM ? WHERE id = ?", driver.table, id)
	return err
}

// Cleanup cleans up the expired pastes
func (driver *SQLDriver) Cleanup() (int, error) {
	// Retrieve all paste IDs
	ids, err := driver.ListIDs()
	if err != nil {
		return 0, err
	}

	// Define the amount of deleted items
	deleted := 0

	// Loop through all pastes
	for _, id := range ids {
		// Retrieve the paste object
		paste, err := driver.Get(id)
		if err != nil {
			return 0, err
		}

		// Delete the paste if it is expired
		lifetime := env.Duration("AUTODELETE_LIFETIME", 30*24*time.Hour)
		if paste.AutoDelete && paste.Created+int64(lifetime.Seconds()) < time.Now().Unix() {
			err = driver.Delete(id)
			if err != nil {
				return 0, err
			}
			deleted++
		}
	}
	return deleted, nil
}
