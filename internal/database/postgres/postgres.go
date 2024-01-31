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
		Uuid CHAR(42) NOT NULL,
		Machine text NOT NULL REFERENCES Machines (Hostname),
		Name text NOT NULL,
		Brand text NOT NULL,
		DriverVersion text NOT NULL,
		MemoryTotal integer NOT NULL,
		PRIMARY KEY (Uuid)
	);`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Stats (
		Gpu CHAR(42) REFERENCES GPUs (Uuid) NOT NULL,
		Received timestamp NOT NULL,
		MemoryUtilisation real NOT NULL,
		GpuUtilisation real NOT NULL,
		MemoryUsed real NOT NULL,
		FanSpeed real NOT NULL,
		Temp real NOT NULL,
		MemoryTemp real NOT NULL,
		GraphicsVoltage real NOT NULL,
		PowerDraw real NOT NULL,
		GraphicsClock real NOT NULL,
		MaxGraphicsClock real NOT NULL,
		MemoryClock real NOT NULL,
		MaxMemoryClock real NOT NULL,
		PRIMARY KEY (Gpu, Received)
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

	// check if machine exists
	lastSeen, err := getLastSeen(host, tx)

	now := time.Now()
	if err == nil {
		// machine existed, check if time is in future
		if lastSeen.Before(now) {
			// last seen was before now, update it
			err = updateLastSeen(host, now, tx)

			if err != nil {
				return errors.Join(err, tx.Rollback())
			}
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		// this machine isn't in the db, so add it
		err = createMachine(host, now, tx)
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

func (conn postgresConn) AppendDataPoint(sample uplink.GPUStatSample) error {
	now := time.Now()

	_, err := conn.db.Exec(`INSERT INTO Stats
		(Gpu, Received, MemoryUtilisation, GpuUtilisation, MemoryUsed,
		FanSpeed, Temp, MemoryTemp, GraphicsVoltage, PowerDraw,
		GraphicsClock, MaxGraphicsClock, MemoryClock, MaxMemoryClock)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
		$14)`,
		sample.Uuid, now,
		sample.MemoryUtilisation, sample.GPUUtilisation,
		sample.MemoryUsed, sample.FanSpeed, sample.Temp,
		sample.MemoryTemp, sample.GraphicsVoltage, sample.PowerDraw,
		sample.GraphicsClock, sample.MaxGraphicsClock,
		sample.MemoryClock, sample.MaxMemoryClock)

	return err
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

func (conn postgresConn) UpdateGPUContext(host string, packet uplink.GPUInfo) error {
	// Insert the new context we've received into the db, overwriting the
	// existing info
	_, err := conn.db.Exec(`INSERT INTO GPUs
		(Uuid, Machine, Name, Brand, DriverVersion, MemoryTotal)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (Uuid) DO UPDATE
		SET (Uuid, Machine, Name, Brand, DriverVersion, MemoryTotal)
		= (EXCLUDED.Uuid, EXCLUDED.Machine, EXCLUDED.Name,
		EXCLUDED.Brand, EXCLUDED.DriverVersion, EXCLUDED.MemoryTotal)`,
		packet.Uuid, host, packet.Name, packet.Brand,
		packet.DriverVersion, packet.MemoryTotal)

	return err
}

func (conn postgresConn) LatestData() ([]uplink.GpuStatsUpload, error) {
	rows, err := conn.db.Query(`SELECT g.Machine, g.Name, g.Brand,
			g.DriverVersion, g.MemoryTotal, s.MemoryUtilisation,
			s.GpuUtilisation, s.MemoryUsed, s.FanSpeed, s.Temp
		FROM GPUs g INNER JOIN Stats s ON g.Uuid = s.Gpu
		INNER JOIN (
			SELECT Gpu, Max(Received) Received
			FROM Stats
			GROUP BY Gpu
		) latest ON s.Gpu = latest.Gpu AND s.Received = latest.Received
	`)

	if err != nil {
		return nil, err
	}

	// collect rows into a map of hostname to {GPUInfo, Stat} because
	// they'll come out of the db out of hostname order
	type gpus struct {
		infos []uplink.GPUInfo
		stats []uplink.GPUStatSample
	}
	var latest = make(map[string]gpus)

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
		latest[host] = gpus{infos: append(latest[host].infos, info), stats: append(latest[host].stats, stat)}
	}

	// flatten map structure
	var result []uplink.GpuStatsUpload
	for key, value := range latest {
		result = append(result, uplink.GpuStatsUpload{Hostname: key,
			GPUInfos: value.infos, Stats: value.stats})
	}

	return result, rows.Close()
}
