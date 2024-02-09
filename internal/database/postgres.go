package database

import (
	"database/sql"
	"errors"
	"log/slog"
	"reflect"
	"time"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/uplink"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// struct holding database context
// only holds a pointer, so we can pass it around by value
type postgresConn struct {
	db *sql.DB
}

func Postgres(databaseUrl string) (Database, error) {
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
		GroupName text,
		CPU text,
		Motherboard text,
		Notes text,
		LastSeen timestamp,
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
func (conn postgresConn) UpdateLastSeen(host string, given_time int64) error {
	var err error

	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	// check if machine exists
	lastSeen, err := getLastSeen(host, tx)

	now := time.Unix(given_time, 0)

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

func (conn postgresConn) Downsample(int_now int64) error {
	downsample_query := `CREATE TEMPORARY TABLE TempDownsampled AS
WITH OrderedStats AS (
  SELECT
    Gpu,
    Received,
    MemoryUtilisation,
    GpuUtilisation,
    MemoryUsed,
    FanSpeed,
    Temp,
    MemoryTemp,
    GraphicsVoltage,
    PowerDraw,
    GraphicsClock,
    MaxGraphicsClock,
    MemoryClock,
    MaxMemoryClock,
    ROW_NUMBER() OVER (PARTITION BY Gpu ORDER BY Received ASC) - 1 AS RowNum
  FROM Stats
  WHERE Received > $1
),
GroupedStats AS (
  SELECT
    Gpu,
    AVG(MemoryUtilisation) AS AvgMemoryUtilisation,
    AVG(GpuUtilisation) AS AvgGpuUtilisation,
    AVG(MemoryUsed) AS AvgMemoryUsed,
    AVG(FanSpeed) AS AvgFanSpeed,
    AVG(Temp) AS AvgTemp,
    AVG(MemoryTemp) AS AvgMemoryTemp,
    AVG(GraphicsVoltage) AS AvgGraphicsVoltage,
    AVG(PowerDraw) AS AvgPowerDraw,
    AVG(GraphicsClock) AS AvgGraphicsClock,
    AVG(MaxGraphicsClock) AS AvgMaxGraphicsClock,
    AVG(MemoryClock) AS AvgMemoryClock,
    AVG(MaxMemoryClock) AS AvgMaxMemoryClock,
    MIN(Received) AS SampleStartTime,
    MAX(Received) AS SampleEndTime,
    (RowNum / 100) AS GroupId
  FROM OrderedStats
  GROUP BY Gpu, GroupId
)
SELECT * FROM GroupedStats;`

	delete_query := `DELETE FROM Stats
WHERE Received > $1
AND Received <= (SELECT MAX(SampleEndTime) FROM TempDownsampled);
	`

	insert_query := `INSERT INTO Stats (Gpu, Received, MemoryUtilisation, GpuUtilisation, MemoryUsed, FanSpeed, Temp, MemoryTemp, GraphicsVoltage, PowerDraw, GraphicsClock, MaxGraphicsClock, MemoryClock, MaxMemoryClock)
SELECT
  Gpu,
  SampleStartTime, 
  AvgMemoryUtilisation,
  AvgGpuUtilisation,
  AvgMemoryUsed,
  AvgFanSpeed,
  AvgTemp,
  AvgMemoryTemp,
  AvgGraphicsVoltage,
  AvgPowerDraw,
  AvgGraphicsClock,
  AvgMaxGraphicsClock,
  AvgMemoryClock,
  AvgMaxMemoryClock
FROM TempDownsampled;
	`

	cleanup_query := `DROP TABLE TempDownsampled;`

	now := time.Unix(int_now, 0)
	sixMonthsAgo := now.AddDate(0, -6, 0)
	sixMonthsAgoFormatted := sixMonthsAgo.Format("2006-01-02 15:04:05")

	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(downsample_query, sixMonthsAgoFormatted)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(delete_query, sixMonthsAgoFormatted)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(insert_query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = tx.Exec(cleanup_query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// TODO: consider returning workstationGroup
func (conn postgresConn) LatestData() ([]uplink.GpuStatsUpload, error) {
	// we pull Uuid twice so we can put one into Stat sample and the other into Info
	rows, err := conn.db.Query(`SELECT g.Machine, g.Uuid, g.Uuid, g.Name,
			g.Brand, g.DriverVersion, g.MemoryTotal,
			s.MemoryUtilisation, s.GpuUtilisation, s.MemoryUsed,
			s.FanSpeed, s.Temp, s.MemoryTemp, s.GraphicsVoltage,
			s.PowerDraw, s.GraphicsClock, s.MaxGraphicsClock,
			s.MemoryClock, s.MaxMemoryClock
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

		err = rows.Scan(&host, &info.Uuid, &stat.Uuid,
			&info.Name, &info.Brand, &info.DriverVersion,
			&info.MemoryTotal,
			&stat.MemoryUtilisation, &stat.GPUUtilisation,
			&stat.MemoryUsed, &stat.FanSpeed, &stat.Temp,
			&stat.MemoryTemp, &stat.GraphicsVoltage,
			&stat.PowerDraw, &stat.GraphicsClock,
			&stat.MaxGraphicsClock, &stat.MemoryClock,
			&stat.MaxMemoryClock,
		)

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

// Create new machine
func (conn postgresConn) NewMachine(machine broadcast.NewMachine) (err error) {
	_, err = conn.db.Exec(`INSERT INTO Machines (Hostname, GroupName)
		VALUES ($1, $2)`,
		machine.Hostname, machine.Group,
	)
	return
}

// Update machine info
func (conn postgresConn) UpdateMachine(machine broadcast.ModifyMachine) error {
	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	v := reflect.ValueOf(machine)
	for _, field := range reflect.VisibleFields(reflect.TypeOf(machine)) {
		value := v.FieldByIndex(field.Index)
		if v.Kind() == reflect.Pointer && !value.IsNil() {
			_, err = tx.Exec(`UPDATE Machines
				SET $1=$2
				WHERE Hostname=$3`,
				field.Name, reflect.Indirect(value), machine.Hostname,
			)

			if err != nil {
				return errors.Join(err, tx.Rollback())
			}
		}
	}

	return tx.Commit()
}

// drop all tables we create in the database
func (conn postgresConn) Drop() error {
	_, err := conn.db.Exec(`DROP TABLE stats;
		DROP TABLE gpus;
		DROP TABLE machines`)

	return errors.Join(err, conn.db.Close())
}

func (conn postgresConn) LastSeen() ([]uplink.WorkstationSeen, error) {
	rows, err := conn.db.Query(`SELECT * FROM Machines`)

	if err != nil {
		return nil, err
	}

	var seens []uplink.WorkstationSeen

	for rows.Next() {
		var seen_instance uplink.WorkstationSeen
		var t time.Time
		var dud sql.NullString

		err = rows.Scan(&seen_instance.Hostname, &dud, &dud, &dud, &dud, &t)

		seen_instance.LastSeen = t.Unix()

		if err != nil {
			return nil, err
		}

		slog.Debug("Fetched last seen instance from Machine table", "Hostname", seen_instance.Hostname, "LastSeen", seen_instance.LastSeen)
		seens = append(seens, seen_instance)
	}

	return seens, nil
}
