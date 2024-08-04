package repo

import (
	"database/sql"
	"errors"
	"log"
	"main/app/model/db"
	"time"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS workday(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date DATETIME NOT NULL UNIQUE
    );

	CREATE TABLE IF NOT EXISTS worktime(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		time DATETIME NOT NULL,
		workday INTEGER,
		FOREIGN KEY(workday) REFERENCES workday(id)
	);
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) AddWorkday(workday db.Workday) (*db.Workday, error) {
	query := `INSERT INTO workday(null, date) VALUES(?)`
	res, err := r.db.Exec(query, workday.Date)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	workday.ID = id

	return &workday, nil
}

func (r *SQLiteRepository) AddWorktime(worktime db.Worktime) (*db.Worktime, error) {
	query := `INSERT INTO worktime(null, type, time, workday) VALUES(?, ?, ?)`
	res, err := r.db.Exec(query, worktime.Type, worktime.Time)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	worktime.ID = id
	return &worktime, nil
}

func (r *SQLiteRepository) GetWorkday(date time.Time) (*db.Workday, error) {
	query := `SELECT id, date FROM workday WHERE date = ?`
	var workday db.Workday

	err := r.db.QueryRow(query, date).Scan(&workday.ID, &workday.Date)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &workday, nil
}

func (r *SQLiteRepository) GetAllWorktime(workday db.Workday) ([]db.Worktime, error) {
	query := `SELECT id, type, time, workday FROM worktime WHERE workday = ?`

	rows, err := r.db.Query(query, workday.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var worktime []db.Worktime
	for rows.Next() {
		var w db.Worktime
		err := rows.Scan(&w)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		worktime = append(worktime, w)
	}

	return worktime, nil
}
