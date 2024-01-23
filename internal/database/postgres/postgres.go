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

	// try pinging database to verify connection
	err = db.Ping()
	if err != nil {
		return nil, err
	} else {
		return postgresConn{db}, nil
	}
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
