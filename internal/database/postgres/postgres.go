package postgres

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/gpuctl/gpuctl/internal/database"
	"github.com/gpuctl/gpuctl/internal/uplink"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// struct holding database context
// only holds a pointer, so we can pass it around by value
type postgresConn struct {
	db *sql.DB
}

func New(databaseUrl string) (database.Database, error) {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return nil, err
	}

	// sql.Open won't make a connection til use
	// so try pinging database to verify connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return postgresConn{db}, nil
}

func createTables(db *sql.DB) error {
	// TODO: Find a way to generate this from gpustats.go?
	// TODO: Allow passing in a parameter to create temporary tables for use
	// with the unit tests

	var err error
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Machines (
		Hostname text,
		LastSeen timestamp,
		PRIMARY KEY (Hostname)
	);`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS GPUs (
		Id integer,
		Machine text REFERENCES Machines (Hostname),
		Name text,
		Brand text,
		DriverVersion text,
		MemoryTotal integer,
		PRIMARY KEY (Id)
	);`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Stats (
		Gpu integer REFERENCES GPUs (Id),
		Recieved timestamp,
		MemoryUtilisation real,
		GpuUtilisation real,
		MemoryUsed real,
		FanSpeed real,
		Temp real,
		PRIMARY KEY (Gpu, Recieved)
	);`)

	return err
}

// implement interface
func (conn postgresConn) UpdateLastSeen(host string) error {
	var err error

	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	slog.Debug("checking if machine exists")

	rows, err := tx.Query(`SELECT LastSeen
		FROM Machines
		WHERE Hostname=$1`,
		host)

	if err != nil {
		// rolling back a transaction can fail, so join it with the
		// error that caused us to rollback
		return errors.Join(err, tx.Rollback())
	}

	if rows.Next() {
		slog.Debug("machine existed, check if time is in future")

		var lastSeen time.Time
		err = rows.Scan(&lastSeen)
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}

		// if not shut before the next transaction, causes
		// "driver: bad connection" errors
		err = rows.Close()
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}

		now := time.Now()
		if lastSeen.Before(now) {
			slog.Debug("last seen was before now, update it")
			_, err = tx.Exec(`UPDATE Machines
				SET LastSeen=$1
				WHERE Hostname=$2`,
				now, host)

			if err != nil {
				return errors.Join(err, tx.Rollback())
			}
		}
	} else {
		slog.Debug("row for this hostname doesn't exist, make it",
			"hostname", host)

		// Next() failing might be because of an error
		err = rows.Err()
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}

		// a machine with this hostname wasn't found, so make a new row
		// for it
		_, err = tx.Exec(`INSERT INTO Machines (Hostname, LastSeen)
			VALUES ($1, $2)`,
			host, time.Now())

		if err != nil {
			return errors.Join(err, tx.Rollback())
		}

		// TODO: in future, it should be added to a list to wait for
		// approval
	}

	return tx.Commit()
}

func (conn postgresConn) AppendDataPoint(packet uplink.GPUStats) error {
	return nil
}

func (conn postgresConn) LatestData() ([]uplink.GPUStats, error) {
	return nil, nil
}
