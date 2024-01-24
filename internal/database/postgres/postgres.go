package postgres

import (
	"database/sql"

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
	return nil
}

func (conn postgresConn) AppendDataPoint(packet uplink.GPUStats) error {
	return nil
}

func (conn postgresConn) LatestData() ([]uplink.GPUStats, error) {
	return nil, nil
}
