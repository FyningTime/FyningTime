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
        date DATETIME NOT NULL UNIQUE,
		breaktime DATETIME DEFAULT '00:30:00'
    );

	CREATE TABLE IF NOT EXISTS worktime(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		time DATETIME DEFAULT (datetime('now', 'localtime')),
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
	query := `INSERT INTO worktime(type, workday, time) VALUES(?, ?, ?)`
	res, err := r.db.Exec(query, worktime.Type, worktime.Workday.ID, worktime.Time)
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
	query := `SELECT id, date from workday WHERE date = ? ORDER BY date DESC LIMIT 1`
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
	query := `SELECT id, date FROM workday ORDER BY date DESC`

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

	loc, _ := time.LoadLocation("Europe/Berlin")

	rows, err := r.db.Query(query, workday.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var worktimes []*db.Worktime
	for rows.Next() {
		var w db.Worktime
		tmpTime := w.Time.In(loc)
		err := rows.Scan(&w.ID, &w.Type, &tmpTime, &w.Workday.ID)
		w.Time = tmpTime
		if err != nil {
			log.Error(err)
			return nil, err
		}
		worktimes = append(worktimes, &w)
	}
	log.Debug("worktimes", "size", len(worktimes))
	return worktimes, nil
}

func (r *SQLiteRepository) DeleteWorktime(worktime *db.Worktime) (int64, error) {
	log.Info("Deleting worktime", "worktime-id", worktime.ID)
	query := `DELETE FROM worktime WHERE id = ?`

	res, err := r.db.Exec(query, worktime.ID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

func (r *SQLiteRepository) DeleteWorkday(workday *db.Workday) (int64, error) {
	log.Info("Deleting workday", "workday-id", workday.ID)
	queryDeleteWorktimes := `DELETE FROM worktime WHERE workday = ?`
	resWt, errWt := r.db.Exec(queryDeleteWorktimes, workday.ID)
	if errWt == nil {
		query := `DELETE FROM workday WHERE id = ?`
		resWd, _ := r.db.Exec(query, workday.ID)
		return resWd.RowsAffected()
	} else {
		return resWt.RowsAffected()
	}

}

func (r *SQLiteRepository) UpdateWorktime(worktime *db.Worktime) (int64, error) {
	log.Info("Updating worktime", "worktime-id", worktime.ID)
	query := `UPDATE worktime SET type = ?, time = ?, breaktime = ? WHERE id = ?`

	loc, _ := time.LoadLocation("Europe/Berlin")
	tmpTime := worktime.Time.In(loc)
	res, err := r.db.Exec(query, worktime.Type, tmpTime, worktime.Breaktime, worktime.ID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}
