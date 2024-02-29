package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	//	"reflect"
	"strings"
	"time"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/uplink"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresConn represents an open connection to a control database backed by postgres.
type PostgresConn struct {
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

	return PostgresConn{db}, nil
}

func createTables(db *sql.DB) error {
	// TODO: Find a way to generate this from gpustats.go?

	// We have to make all rows non-null, because we can't scan a null value
	// into a Go variable

	var err error
	_, err = db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS Machines (
		Hostname text NOT NULL,
		GroupName text NOT NULL DEFAULT '%s',
		CPU text,
		Motherboard text,
		Notes text,
		Owner text,
		LastSeen timestamp,
		PRIMARY KEY (Hostname)
	);`, DefaultGroup))

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS GPUs (
		Uuid text NOT NULL,
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Files (
		Hostname text NOT NULL REFERENCES Machines (Hostname),
		Filename text NOT NULL,
		Mime text NOT NULL,
		File text NOT NULL,
		PRIMARY KEY (Hostname, Filename)
	);`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Stats (
		Gpu text REFERENCES GPUs (Uuid) NOT NULL,
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
		InUse boolean NOT NULL,
		UserName text NOT NULL DEFAULT '',
		PRIMARY KEY (Gpu, Received)
	);`)

	return err
}

// implement interface
func (conn PostgresConn) UpdateLastSeen(host string, given_time int64) error {
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

func (conn PostgresConn) AppendDataPoint(sample uplink.GPUStatSample) error {
	now := time.Now()

	inUse, user := sample.RunningProcesses.Summarise()
	_, err := conn.db.Exec(`INSERT INTO Stats
		(Gpu, Received, MemoryUtilisation, GpuUtilisation, MemoryUsed,
		FanSpeed, Temp, MemoryTemp, GraphicsVoltage, PowerDraw,
		GraphicsClock, MaxGraphicsClock, MemoryClock, MaxMemoryClock,
		InUse, UserName)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,
		$14, $15, $16)`,
		sample.Uuid, now,
		sample.MemoryUtilisation, sample.GPUUtilisation,
		sample.MemoryUsed, sample.FanSpeed, sample.Temp,
		sample.MemoryTemp, sample.GraphicsVoltage, sample.PowerDraw,
		sample.GraphicsClock, sample.MaxGraphicsClock,
		sample.MemoryClock, sample.MaxMemoryClock,
		inUse, user)

	// TODO: hacky, untested. Might not work
	if err != nil {
		return ErrGpuNotPresent
	}

	return err
}

func (conn PostgresConn) UpdateGPUContext(host string, packet uplink.GPUInfo) error {
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

func (conn PostgresConn) Downsample(int_now int64) error {
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
    InUse,
    UserName,
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
    (RowNum / 100) AS GroupId,
    bool_or(InUse) AS OrInUse,
    mode() WITHIN GROUP (ORDER BY UserName) as ModeUserName
  FROM OrderedStats
  GROUP BY Gpu, GroupId
)
SELECT * FROM GroupedStats;`

	delete_query := `DELETE FROM Stats
WHERE Received > $1
AND Received <= (SELECT MAX(SampleEndTime) FROM TempDownsampled);
	`

	insert_query := `INSERT INTO Stats (Gpu, Received, MemoryUtilisation, GpuUtilisation, MemoryUsed, FanSpeed, Temp, MemoryTemp, GraphicsVoltage, PowerDraw, GraphicsClock, MaxGraphicsClock, MemoryClock, MaxMemoryClock, InUse, UserName)
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
  AvgMaxMemoryClock,
  OrInUse,
  ModeUserName
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
func (conn PostgresConn) LatestData() (broadcast.Workstations, error) {
	// pull all the machines in, then gpus for those machines
	// has to all be done in a transaction to avoid race conditions
	tx, err := conn.db.Begin()
	if err != nil {
		return nil, err
	}

	groups := make(map[string][]broadcast.Workstation)
	machines, err := tx.Query(`SELECT GroupName, Hostname, CPU, Motherboard,
		Notes, Owner, LastSeen
		FROM Machines`)
	if err != nil {
		return nil, errors.Join(err, tx.Rollback())
	}

	for machines.Next() {
		var groupName string
		var machine broadcast.Workstation
		var lastSeen time.Time

		err = machines.Scan(&groupName, &machine.Name, &machine.CPU,
			&machine.Motherboard, &machine.Notes, &machine.Owner, &lastSeen)
		if err != nil {
			return nil, errors.Join(err, tx.Rollback())
		}

		// coalesce null and empty group names to default
		if strings.TrimSpace(groupName) == "" {
			groupName = DefaultGroup
		}

		machine.LastSeen = time.Since(lastSeen)
		machine.Gpus = nil

		groups[groupName] = append(groups[groupName], machine)
	}

	// check for error whilst iterating, continuing if it's "no results"
	err = machines.Err()
	if err != nil {
		return nil, errors.Join(err, tx.Rollback())
	}

	// attach gpus to all machines
	// can't be done in the previous loop because we can't be iterating
	// through two queries at once
	for group := range groups {
		for i, machine := range groups[group] {
			machine.Gpus, err = getGpus(machine.Name, tx)
			if err != nil {
				return nil, errors.Join(err, tx.Rollback())
			}
			groups[group][i] = machine
		}
	}

	// flatten map
	var result broadcast.Workstations
	for groupName, machines := range groups {
		result = append(result, broadcast.Group{
			Name:         groupName,
			Workstations: machines,
		})
	}

	return result, tx.Commit()
}

// get the latest stat for all the gpus on a machine
func getGpus(host string, tx *sql.Tx) ([]broadcast.GPU, error) {
	result := make([]broadcast.GPU, 0)

	gpus, err := tx.Query(`SELECT g.Uuid, g.Name, g.Brand,
		g.DriverVersion, g.MemoryTotal,
		s.MemoryUtilisation, s.GpuUtilisation,
		s.MemoryUsed, s.FanSpeed, s.Temp, s.MemoryTemp,
		s.GraphicsVoltage, s.PowerDraw, s.GraphicsClock,
		s.MaxGraphicsClock, s.MemoryClock,
		s.MaxMemoryClock, s.InUse, s.UserName
		FROM GPUs g INNER JOIN Stats s ON g.Uuid = s.Gpu
		INNER JOIN (
			SELECT Gpu, Max(Received) Received
			FROM Stats
			GROUP BY Gpu
		) latest ON s.Gpu = latest.Gpu
			AND s.Received = latest.Received
		WHERE g.Machine=$1`,
		host,
	)
	if err != nil {
		return nil, err
	}

	for gpus.Next() {
		var gpu broadcast.GPU
		err = gpus.Scan(&gpu.Uuid, &gpu.Name, &gpu.Brand,
			&gpu.DriverVersion, &gpu.MemoryTotal,
			&gpu.MemoryUtilisation,
			&gpu.GPUUtilisation, &gpu.MemoryUsed,
			&gpu.FanSpeed, &gpu.Temp,
			&gpu.MemoryTemp, &gpu.GraphicsVoltage,
			&gpu.PowerDraw, &gpu.GraphicsClock,
			&gpu.MaxGraphicsClock, &gpu.MemoryClock,
			&gpu.MaxMemoryClock, &gpu.InUse, &gpu.User)
		if err != nil {
			return nil, err
		}

		result = append(result, gpu)
	}

	err = gpus.Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Create new machine
func (conn PostgresConn) NewMachine(machine broadcast.NewMachine) (err error) {
	_, err = conn.db.Exec(`INSERT INTO Machines (Hostname, GroupName)
		VALUES ($1, $2)`,
		machine.Hostname, machine.Group,
	)
	return
}

func (conn PostgresConn) RemoveMachine(machine broadcast.RemoveMachine) error {
	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM Stats
		WHERE Gpu=ANY(SELECT Uuid
			FROM Gpus
			WHERE Machine=$1)`,
		machine.Hostname,
	)

	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	_, err = tx.Exec(`DELETE FROM GPUs
		WHERE Machine=$1`,
		machine.Hostname,
	)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	_, err = tx.Exec(`DELETE FROM Files
		WHERE Hostname=$1`,
		machine.Hostname)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	_, err = tx.Exec(`DELETE FROM Machines
		WHERE Hostname=$1`,
		machine.Hostname,
	)
	if err != nil {
		return errors.Join(err, tx.Rollback())
	}

	return tx.Commit()
}

// Update machine info
func (conn PostgresConn) UpdateMachine(machine broadcast.ModifyMachine) error {
	tx, err := conn.db.Begin()
	if err != nil {
		return err
	}

	//	v := reflect.ValueOf(machine)
	//	for _, field := range reflect.VisibleFields(reflect.TypeOf(machine)) {
	//		value := v.FieldByIndex(field.Index)
	//		if v.Kind() == reflect.Pointer && !value.IsNil() {
	//			_, err = tx.Exec(`UPDATE Machines
	//				SET $1=$2
	//				WHERE Hostname=$3`,
	//				field.Name, reflect.Indirect(value), machine.Hostname,
	//			)
	//
	//			if err != nil {
	//				return errors.Join(err, tx.Rollback())
	//			}
	//		}
	//	}

	if machine.CPU != nil {
		slog.Info("Changing CPU", "Hostname", machine.Hostname, "New CPU", *machine.CPU)
		_, err = tx.Exec(`UPDATE Machines
			SET CPU=$1
			WHERE Hostname=$2`,
			*machine.CPU, machine.Hostname,
		)

		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
	}

	if machine.Motherboard != nil {
		slog.Info("Changing Motherboard", "Hostname", machine.Hostname, "New Motherboard", *machine.Motherboard)
		_, err = tx.Exec(`UPDATE Machines
			SET Motherboard=$1
			WHERE Hostname=$2`,
			*machine.Motherboard, machine.Hostname,
		)

		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
	}

	if machine.Notes != nil {
		slog.Info("Changing Notes", "Hostname", machine.Hostname, "New Notes", *machine.Notes)
		_, err = tx.Exec(`UPDATE Machines
			SET Notes=$1
			WHERE Hostname=$2`,
			*machine.Notes, machine.Hostname,
		)

		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
	}

	if machine.Group != nil {
		slog.Info("Changing Group", "Hostname", machine.Hostname, "New Group", *machine.Group)
		_, err = tx.Exec(`UPDATE Machines
			SET GroupName=$1
			WHERE Hostname=$2`,
			*machine.Group, machine.Hostname,
		)

		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
	}

	if machine.Owner != nil {
		slog.Info("Changing Owner", "Hostname", machine.Hostname, "New Owner", *machine.Owner)
		_, err = tx.Exec(`UPDATE Machines
			SET Owner=$1
			WHERE Hostname=$2`,
			*machine.Owner, machine.Hostname,
		)

		if err != nil {
			return errors.Join(err, tx.Rollback())
		}
	}

	return tx.Commit()
}

// Drop drops all tables on the connected database, then closes the connection.
//
// This should only be used for testing purposes
func (conn PostgresConn) Drop() error {
	_, err := conn.db.Exec(`DROP TABLE stats;
		DROP TABLE gpus;
		DROP TABLE files;
		DROP TABLE machines`)
	if err != nil {
		return err
	}

	return conn.db.Close()
}

func (conn PostgresConn) LastSeen() ([]broadcast.WorkstationSeen, error) {
	rows, err := conn.db.Query(`SELECT Hostname, LastSeen FROM Machines`)

	if err != nil {
		return nil, err
	}

	var seens []broadcast.WorkstationSeen

	for rows.Next() {
		var seen_instance broadcast.WorkstationSeen
		var t time.Time

		err = rows.Scan(&seen_instance.Hostname, &t)

		seen_instance.LastSeen = t.Unix()

		if err != nil {
			return nil, err
		}

		slog.Debug("Fetched last seen instance from Machine table", "Hostname", seen_instance.Hostname, "LastSeen", seen_instance.LastSeen)
		seens = append(seens, seen_instance)
	}

	return seens, nil
}

func (conn PostgresConn) AttachFile(attach broadcast.AttachFile) error {
	// Insert into db
	_, err := conn.db.Exec(`INSERT INTO Files (Hostname, Mime, Filename, File)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (Hostname, Filename) DO UPDATE
		SET (Mime, File) = (EXCLUDED.Mime, Excluded.File);`,
		attach.Hostname, attach.Mime, attach.Filename, attach.EncodedFile,
	)
	if err != nil {
		return ErrNoSuchMachine
	}
	return nil
}

func (conn PostgresConn) GetFile(hostname string, filename string) (broadcast.AttachFile, error) {
	file := broadcast.AttachFile{Hostname: hostname, Filename: filename}

	row := conn.db.QueryRow(`SELECT Mime, File
		FROM Files
		WHERE Hostname=$1 AND Filename=$2`,
		hostname, filename)
	err := row.Scan(&file.Mime, &file.EncodedFile)
	if errors.Is(err, sql.ErrNoRows) {
		return file, ErrFileNotPresent
	}

	return file, err
}

func (conn PostgresConn) ListFiles(hostname string) ([]string, error) {
	rows, err := conn.db.Query(`SELECT Filename
		FROM Files
		WHERE Hostname=$1`,
		hostname)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	for rows.Next() {
		val := ""
		err = rows.Scan(&val)
		if err != nil {
			return res, err
		}
		res = append(res, val)
	}
	return res, nil
}

func (conn PostgresConn) RemoveFile(remove broadcast.RemoveFile) error {
	rows, err := conn.db.Query(`DELETE FROM Files
		WHERE Hostname=$1 AND Filename=$2
		RETURNING Filename;`,
		remove.Hostname, remove.Filename)
	if err != nil {
		return err
	}

	found := false
	for rows.Next() {
		found = true
		rows.Scan()
	}

	if !found {
		return ErrFileNotPresent
	}
	return err
}
