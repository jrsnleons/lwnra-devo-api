package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lwnra-devo-api/models"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the database connection and provides methods for database operations
type DB struct {
	conn *sql.DB
}

// New creates a new database connection and initializes the schema
func New(dbPath string) (*DB, error) {
	// Ensure the directory exists for the database file
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %v", err)
		}
	}

	// Add file: prefix for sqlite to handle file creation properly
	if !strings.HasPrefix(dbPath, "file:") {
		dbPath = "file:" + dbPath + "?cache=shared&mode=rwc"
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// createTables creates the necessary database tables
func (db *DB) createTables() error {
	query := `CREATE TABLE IF NOT EXISTS devotionals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		reading TEXT,
		version TEXT,
		passage TEXT,
		refqs TEXT,
		title TEXT,
		author TEXT,
		body TEXT,
		prayer TEXT,
		UNIQUE(date, title)
	)`

	_, err := db.conn.Exec(query)
	return err
}

// SaveDevotional saves a devotional to the database
func (db *DB) SaveDevotional(devo models.Devotional) error {
	query := `INSERT OR IGNORE INTO devotionals
		(date, reading, version, passage, refqs, title, author, body, prayer)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := db.conn.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		devo.Date,
		devo.Reading,
		devo.Version,
		devo.Passage,
		strings.Join(devo.ReflectionQs, "\n"),
		devo.Title,
		devo.Author,
		devo.Body,
		devo.Prayer,
	)

	return err
}

// GetDevotionals retrieves a limited number of devotionals from the database
func (db *DB) GetDevotionals(limit int) ([]models.Devotional, error) {
	query := `SELECT date, reading, version, passage, refqs, title, author, body, prayer 
			  FROM devotionals 
			  ORDER BY date DESC 
			  LIMIT ?`
	
	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devotionals []models.Devotional
	for rows.Next() {
		var devo models.Devotional
		var refqs string
		
		err := rows.Scan(
			&devo.Date,
			&devo.Reading,
			&devo.Version,
			&devo.Passage,
			&refqs,
			&devo.Title,
			&devo.Author,
			&devo.Body,
			&devo.Prayer,
		)
		if err != nil {
			return nil, err
		}
		
		// Convert refqs back to slice
		if refqs != "" {
			devo.ReflectionQs = strings.Split(refqs, "\n")
		}
		
		devotionals = append(devotionals, devo)
	}
	
	return devotionals, nil
}

// GetDevotionalByDate retrieves a devotional by date
func (db *DB) GetDevotionalByDate(date string) (*models.Devotional, error) {
	query := `SELECT date, reading, version, passage, refqs, title, author, body, prayer 
			  FROM devotionals 
			  WHERE date = ? 
			  LIMIT 1`
	
	var devo models.Devotional
	var refqs string
	
	err := db.conn.QueryRow(query, date).Scan(
		&devo.Date,
		&devo.Reading,
		&devo.Version,
		&devo.Passage,
		&refqs,
		&devo.Title,
		&devo.Author,
		&devo.Body,
		&devo.Prayer,
	)
	
	if err != nil {
		return nil, err
	}
	
	// Convert refqs back to slice
	if refqs != "" {
		devo.ReflectionQs = strings.Split(refqs, "\n")
	}
	
	return &devo, nil
}
