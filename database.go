package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

var DB *sql.DB

func connect() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		panic("DB_PATH is not set")
	}
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open the database at %s: %v", dbPath, err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database at %s: %v", dbPath, err)
	}

}

func createRecord(name string) (Record, error) {
	result, err := DB.Exec(`INSERT INTO records (name) VALUES (?)`, name)
	if err != nil {
		return Record{}, fmt.Errorf("Error creating record: %v", err)
	}
  id, err := result.LastInsertId()
  if err != nil {
    return Record{}, fmt.Errorf("Error getting last insert id: %v", err)
  }
	var record Record
	err = DB.QueryRow(`SELECT id, name, created_at FROM records WHERE id = ?`, id).Scan(&record.ID, &record.Name, &record.CreatedAt)
	if err != nil {
		return Record{}, fmt.Errorf("Error retrieving the record: %v", err)
	}
	return record, nil
}

func getRecords() ([]Record, error) {
	var records []Record
	rows, err := DB.Query(`SELECT * FROM records`)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.ID, &r.Name, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		records = append(records, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return records, nil
}

func editRecord(id int64, name string) error {
	result, err := DB.Exec(`UPDATE records SET name = ? WHERE id = ?`, name, id)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No record found with id: %d", id)
	}
	return nil
}

func deleteRecord(id string) error {
	result, err := DB.Exec(`DELETE FROM records WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No record found with id: %s", id)
	}
	return nil
}
