package homestead

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	// commented until we use pq directly
	_ "github.com/lib/pq"
)

type Log struct {
	SensorName string
	SensorID   int
	Stamp      time.Time
	Metrics    []string
}

func (l *Log) SetID(db *sql.DB) (*int, error) {
	i, err := GetSensor(db, l.SensorName)
	if err != nil {
		return nil, err
	}
	l.SensorID = *i
	return i, nil
}

type User struct {
	ID      int
	Created time.Time
	LName   string
	FName   string
	Email   string
	User    string
	Pass    string
	Hash    string
	Authed  bool
	Admin   bool
}

type DataSet struct {
	JSON string
}

type DataSets []DataSet

type DataBlob struct {
	SensorName string `json:"sensorname"`
	SensorID   int    `json:"sensorid"`
	Metrics    DataSets
}

type TopStat struct {
	Min  float64
	Max  float64
	Avg  float64
	Name string
}

type TopStats []TopStat

var topSQL = `
select
        min(temp::float),
        max(temp::float),
        avg(temp::float),
        name
from (
        select
                metrics -> 'temp' as temp,
                sensorid
        from sensorlogs
          join sensors on (sensors.id = sensorid)
        where
        sensors.name = $1
        sensorlogs.created >= now() - '1 day'::INTERVAL
) as a
join sensors on (sensorid = sensors.id)
group by name
`

var monthData = `
select
      hstore_to_json_loose(metrics)
from
      sensorlogs
where
      created >= now() - '1 month'::interval
order by created desc
`

func GetMonthData(db *sql.DB) (*DataBlob, error) {
	var d = &DataBlob{}
	rows, err := db.Query(monthData)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b = DataSet{}
		rows.Scan(&b.JSON)
		d.Metrics = append(d.Metrics, b)
	}

	return d, nil
}

func GetTopStats(db *sql.DB, s string) (*TopStats, error) {
	var d = &TopStats{}
	rows, err := db.Query(topSQL, s)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var b = TopStat{}
		rows.Scan(&b.Min, &b.Max, &b.Avg, &b.Name)

		*d = append(*d, b)
	}
	return d, nil
}

func GetSensor(db *sql.DB, n string) (*int, error) {
	var i int
	err := db.QueryRow(`
select id from sensors where name = $1
`, n).Scan(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func GetSensors(db *sql.DB) (*string, error) {
	var i string
	err := db.QueryRow(`
select to_json(array_agg(t)) from (select * from sensors) as t
`).Scan(&i)

	if err != nil {
		return nil, err
	}

	return &i, nil
}

func GetCurrent(db *sql.DB, s string) (*string, error) {
	var i string
	err := db.QueryRow(`
select
     hstore_to_json_loose(metrics)
from sensorlogs
join sensors on (sensorid = sensors.id)
where sensors.name = $1
order by sensorlogs.created desc
limit 1`, s).Scan(&i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func Auth(db *sql.DB, u string, p string) (*User, error) {
	var user = &User{}

	err := db.QueryRow(`select id, created, fname, lname, email, username, (hash = crypt
($1, hash)) as authed, admin from users where username = $2`, p, u).Scan(&user.ID, &user.Created, &user.FName, &user.LName, &user.Email, &user.User, &user.Authed, &user.Admin)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func InsertLog(db *sql.DB, log *Log) (*int, error) {
	var id int
	fmt.Printf("%v", log)
	err := db.QueryRow(`
insert into sensorlogs (sensorid, metrics) values ($1, $2) returning id
`, log.SensorID, strings.Join(log.Metrics, ", ")).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
