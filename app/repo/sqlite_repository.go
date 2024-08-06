package repo

import (
	"database/sql"
	"errors"
	"main/app/model/db"
	"time"

	"github.com/charmbracelet/log"
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
	// TODO might be a good idea to move this to a separate file with history table
	query := `
    CREATE TABLE IF NOT EXISTS workday(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        date DATETIME NOT NULL UNIQUE
    );

	CREATE TABLE IF NOT EXISTS worktime(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		time DATETIME DEFAULT CURRENT_TIMESTAMP,
		breaktime DATETIME,
		workday INTEGER,
		FOREIGN KEY(workday) REFERENCES workday(id)
	);
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) AddWorkday(workday *db.Workday) (*db.Workday, error) {
	log.Info("Adding workday", "date", workday.Date)
	query := `INSERT INTO workday(date) VALUES(?)`
	res, err := r.db.Exec(query, workday.Date.Format(time.DateOnly))
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	workday.ID = id

	return workday, nil
}

func (r *SQLiteRepository) AddWorktime(worktime *db.Worktime) (*db.Worktime, error) {
	log.Info("Adding worktime", "type", worktime.Type, "time", worktime.Time)
	query := `INSERT INTO worktime(type, workday) VALUES(?, ?)`
	res, err := r.db.Exec(query, worktime.Type, worktime.Workday.ID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	worktime.ID = id
	return worktime, nil
}

func (r *SQLiteRepository) GetWorkday(date time.Time) (*db.Workday, error) {
	log.Info("Getting workday", "date", date)
	query := `SELECT id, date from workday WHERE date = ?`
	var workday db.Workday

	err := r.db.QueryRow(query, date.Format(time.DateOnly)).Scan(&workday.ID, &workday.Date)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &workday, nil
}

func (r *SQLiteRepository) GetAllWorkday() ([]*db.Workday, error) {
	log.Info("Getting all workdays")
	query := `SELECT id, date FROM workday`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workday []*db.Workday
	for rows.Next() {
		var w db.Workday
		err := rows.Scan(&w.ID, &w.Date)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		log.Debug("Workday", "id", w.ID, "date", w.Date)
		workday = append(workday, &w)
	}

	return workday, nil
}

func (r *SQLiteRepository) GetAllWorktime(workday *db.Workday) ([]*db.Worktime, error) {
	log.Info("Getting all worktimes", "workday-id", workday.ID)
	query := `SELECT id, type, time, workday FROM worktime WHERE workday = ?`

	rows, err := r.db.Query(query, workday.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var worktimes []*db.Worktime
	for rows.Next() {
		var w db.Worktime
		err := rows.Scan(&w.ID, &w.Type, &w.Time, &w.Workday.ID)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		worktimes = append(worktimes, &w)
	}
	log.Debug("worktimes", "size", len(worktimes))
	return worktimes, nil
}
