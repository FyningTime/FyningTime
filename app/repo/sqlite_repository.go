package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/FyningTime/FyningTime/app/model/db"

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

type SORTING string

const (
	ASC  SORTING = "ASC"
	DESC SORTING = "DESC"
)

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
		time TEXT default '01.01.1970',
		breaktime INTEGER DEFAULT 0
    );

	CREATE TABLE IF NOT EXISTS worktime(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		time DATETIME,
		workday INTEGER,
		FOREIGN KEY(workday) REFERENCES workday(id)
	);

	CREATE TABLE IF NOT EXISTS vacations(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		startdate DATETIME NOT NULL UNIQUE,
		enddate DATETIME NOT NULL UNIQUE
	);
    `

	_, err := r.db.Exec(query)
	if err != nil {
		log.Error(err)
		return err
	}

	query = `
		ALTER TABLE workday ADD COLUMN overtime TEXT DEFAULT "";

		CREATE INDEX IF NOT EXISTS idx_worktime_workday ON worktime(workday);
		CREATE INDEX IF NOT EXISTS idx_workday_date ON workday(date);
	`
	_, err = r.db.Exec(query)
	if err != nil {
		log.Warn(err)
		// Ignoring error as columns might already exist
	}

	return nil
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
	query := `SELECT id, date, time, breaktime, overtime
	from workday WHERE date = ? ORDER BY date DESC LIMIT 1`

	var w db.Workday
	loc, _ := time.LoadLocation("Europe/Berlin")
	tmpDate := w.Date.In(loc)
	qd := date.In(loc).Format(time.DateOnly)
	log.Debug("Query date", "date", qd)
	err := r.db.QueryRow(query, qd).
		Scan(&w.ID, &tmpDate, &w.Time, &w.Breaktime, &w.Overtime)
	w.Date = tmpDate

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &w, nil
}

func (r *SQLiteRepository) GetAllWorkday(sorting SORTING) ([]*db.Workday, error) {
	log.Debug("Getting all workdays")
	query := fmt.Sprintf(`SELECT id, date, time, breaktime, overtime FROM workday ORDER BY date %s`, sorting)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workdays []*db.Workday
	for rows.Next() {
		var w db.Workday
		err := rows.Scan(&w.ID, &w.Date, &w.Time, &w.Breaktime, &w.Overtime)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		log.Debug("Workday", "id", w.ID, "date", w.Date)
		workdays = append(workdays, &w)
	}

	return workdays, nil
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
	query := `UPDATE worktime SET type = ?, time = ? WHERE id = ?`

	loc, _ := time.LoadLocation("Europe/Berlin")
	tmpTime := worktime.Time.In(loc)
	res, err := r.db.Exec(query, worktime.Type, tmpTime, worktime.ID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

func (r *SQLiteRepository) UpdateWorkday(workday *db.Workday) (int64, error) {
	log.Info("Updating workday", "workday", workday)
	query := `UPDATE workday
		SET breaktime = ?, time = ?
		WHERE id = ?`

	res, err := r.db.Exec(query,
		workday.Breaktime,
		workday.Time,
		workday.ID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

func (r *SQLiteRepository) AddVacation(vacation *db.Vacation) (*db.Vacation, error) {
	log.Info("Adding vacation", "start", vacation.StartDate, "end", vacation.EndDate)
	query := `INSERT INTO vacations(startdate, enddate) VALUES(?, ?)`

	loc, _ := time.LoadLocation("Europe/Berlin")
	vacation.StartDate = vacation.StartDate.In(loc)
	vacation.EndDate = vacation.EndDate.In(loc)

	res, err := r.db.Exec(query,
		vacation.StartDate,
		vacation.EndDate,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	vacation.ID = id
	return vacation, nil
}

func (r *SQLiteRepository) GetAllVacation() ([]*db.Vacation, error) {
	log.Info("Getting all vacations")
	query := `SELECT ID, startdate, enddate FROM vacations ORDER BY startdate DESC`

	var v []*db.Vacation
	rows, err := r.db.Query(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var vacation db.Vacation
		loc, _ := time.LoadLocation("Europe/Berlin")

		sd := vacation.StartDate.In(loc)
		ed := vacation.EndDate.In(loc)

		err := rows.Scan(&vacation.ID, &sd, &ed)

		if err != nil {
			log.Error(err)
			return nil, err
		}
		vacation.StartDate = sd
		vacation.EndDate = ed
		v = append(v, &vacation)
	}
	log.Debug("Vacations", "size", len(v))

	return v, nil
}

func (r *SQLiteRepository) DeleteVacation(vacation *db.Vacation) (int64, error) {
	log.Info("Deleting vacation", "vacation-id", vacation.ID)
	query := `DELETE FROM vacations WHERE id = ?`

	res, err := r.db.Exec(query, vacation.ID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

func (r *SQLiteRepository) UpdateVacation(vacation *db.Vacation) (int64, error) {
	log.Info("Updating vacation", "vacation-id", vacation.ID)
	query := `UPDATE vacations SET startdate = ?, enddate = ? WHERE id = ?`

	loc, _ := time.LoadLocation("Europe/Berlin")
	startDate := vacation.StartDate.In(loc)
	endDate := vacation.EndDate.In(loc)

	res, err := r.db.Exec(query, startDate, endDate, vacation.ID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

func (r *SQLiteRepository) UpdateOvertimes(workday *db.Workday) (int64, error) {
	log.Info("Updating overtimes", "wd", workday)
	query := `UPDATE Workday SET overtime = ? WHERE id = ?`

	res, err := r.db.Exec(query, workday.Overtime, workday.ID)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}
