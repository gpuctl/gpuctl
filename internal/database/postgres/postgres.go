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

	// We have to make all rows non-null, because we can't scan a null value
	// into a Go variable

	var err error
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Machines (
		Hostname text NOT NULL,
		LastSeen timestamp NOT NULL,
		PRIMARY KEY (Hostname)
	);`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS GPUs (
		Id serial NOT NULL,
		Machine text NOT NULL REFERENCES Machines (Hostname),
		Name text NOT NULL,
		Brand text NOT NULL,
		DriverVersion text NOT NULL,
		MemoryTotal integer NOT NULL,
		PRIMARY KEY (Id)
	);`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Stats (
		Gpu integer REFERENCES GPUs (Id) NOT NULL,
		Recieved timestamp NOT NULL,
		MemoryUtilisation real NOT NULL,
		GpuUtilisation real NOT NULL,
		MemoryUsed real NOT NULL,
		FanSpeed real NOT NULL,
		Temp real NOT NULL,
		PRIMARY KEY (Gpu, Recieved)
	);`)

	return err
}

// implement interface
func (conn postgresConn) UpdateLastSeen(uuid string) error {
	var err error

	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	slog.Debug("checking if machine exists")
	lastSeen, err := getLastSeen(uuid, tx)

	now := time.Now()
	if err == nil {
		// machine existed, check if time is in future
		if lastSeen.Before(now) {
			// last seen was before now, update it
			err = updateLastSeen(uuid, now, tx)

			if err != nil {
				return errors.Join(err, tx.Rollback())
			}
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		// this machine isn't in the db, so add it
		err = createMachine(uuid, now, tx)
		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
	} else {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

func getLastSeen(host string, tx *sql.Tx) (lastSeen time.Time, err error) {
	row := tx.QueryRow(`SELECT LastSeen
		FROM Machines
		WHERE Hostname=$1`,
		host)
	err = row.Scan(&lastSeen)
	return
}

// TODO: in future we may want to consider a list for machines to wait on
// before insertion into the database
func createMachine(host string, now time.Time, tx *sql.Tx) (err error) {
	_, err = tx.Exec(`INSERT INTO Machines (Hostname, LastSeen)
		VALUES ($1, $2)`,
		host, now)
	return
}

func updateLastSeen(host string, now time.Time, tx *sql.Tx) (err error) {
	_, err = tx.Exec(`UPDATE Machines
		SET LastSeen=$1
		WHERE Hostname=$2`,
		now, host)
	return
}

func (conn postgresConn) AppendDataPoint(host string, packet uplink.GPUStatSample) error {
	var err error

	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	// Find matching GPU

	// TODO: replace with Query
	// This silently discards all rows other than the first
	row := tx.QueryRow(`SELECT Id
		FROM GPUs
		WHERE Machine=$1`,
		host)
	var id int64
	err = row.Scan(&id)

	// check error type. If no rows found, add a new GPU
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Didn't find an existing entry, so add a new GPU

			// use TOFU approach, only set the name, brand, etc. on
			// the first packet; future AppendDataPoint calls only
			// add the 'dynamic' data to Stats
			// TODO: reevaluate whether this is a good idea (what
			// happens when a driver update occurs?)

			// TODO: the great refactor made this code impossible. Commented out for now because it wont work
			/*
				id, err = createGPU(host, packet, tx)
				if err != nil {
					return errors.Join(err, tx.Rollback())
				}
			*/
			return errors.Join(err, tx.Rollback())
		} else {
			return errors.Join(err, tx.Rollback())
		}
	}

	// GPU should now exist, so push stats
	now := time.Now()
	err = insertStats(id, packet, now, tx)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

func createGPU(host string, gpuinfo uplink.GPUInfo, tx *sql.Tx) (id int64, err error) {
	newId := tx.QueryRow(`INSERT INTO GPUs
		(Machine, Name, Brand, DriverVersion, MemoryTotal)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING Id`,
		host, gpuinfo.Name, gpuinfo.Brand, gpuinfo.DriverVersion,
		gpuinfo.MemoryTotal)
	err = newId.Scan(&id)
	return
}

func insertStats(id int64, packet uplink.GPUStatSample, now time.Time, tx *sql.Tx) (err error) {
	_, err = tx.Exec(`INSERT INTO Stats
		(Gpu, Recieved, MemoryUtilisation, GpuUtilisation, MemoryUsed,
			FanSpeed, Temp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		id, now, packet.MemoryUtilisation, packet.GPUUtilisation,
		packet.MemoryUsed, packet.FanSpeed, packet.Temp)

	return
}

func (conn postgresConn) LatestData() (map[string][]uplink.GPUStatSample, error) {
	rows, err := conn.db.Query(`SELECT g.Machine, g.Name, g.Brand,
			g.DriverVersion, g.MemoryTotal, s.MemoryUtilisation,
			s.GpuUtilisation, s.MemoryUsed, s.FanSpeed, s.Temp
		FROM GPUs g INNER JOIN Stats s ON g.Id = s.Gpu
		INNER JOIN (
			SELECT Gpu, Max(Recieved) Recieved
			FROM Stats
			GROUP BY Gpu
		) latest ON s.Gpu = latest.Gpu AND s.Recieved = latest.Recieved
	`)

	if err != nil {
		return nil, err
	}

	var latest = make(map[string][]uplink.GPUStatSample)

	for rows.Next() {
		var host string
		var info uplink.GPUInfo
		var stat uplink.GPUStatSample

		err = rows.Scan(&host, &info.Name, &info.Brand,
			&info.DriverVersion, &info.MemoryTotal,
			&stat.MemoryUtilisation, &stat.GPUUtilisation,
			&stat.MemoryUsed, &stat.FanSpeed, &stat.Temp)

		if err != nil {
			return nil, err
		}

		slog.Debug("got stat from table", "host", host, "info", info, "stat", stat)
		latest[host] = append(latest[host], stat)
	}

	return latest, rows.Close()
}
