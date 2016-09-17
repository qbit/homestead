package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var head []string

func main() {
	file := flag.String("file", "", "File to process.")
	sid := flag.Int("sid", 0, "Sensor ID.")
	flag.Parse()

	sql := `
insert into sensorlogs (sensorid, created, metrics) values (%d, '%s', '%s');
`

	fi, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(fi)
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(strings.NewReader(string(b)))

	count := 0

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			head = record
		}

		if count > 0 {
			var metrics []string
			var t time.Time
			for i, k := range head {
				// '{{temp=>80}, {baro=>100}}'
				if k != "timestamp" {
					metric := fmt.Sprintf(`"%s"=>"%s"`, k, record[i])
					metrics = append(metrics, metric)
				} else {
					// Jan 2, 2006 at 3:04pm (MST)
					// "2006-01-02T15:04:05Z07:00"
					tf := "2006-01-02T15:04:05.999Z"
					t, err = time.Parse(tf, record[i])
					if err != nil {
						fmt.Println(err)
						fmt.Println(record[i], "->", tf)
						os.Exit(1)
					}
				}
			}
			line := fmt.Sprintf(sql, *sid, t.Format(time.RFC3339Nano), strings.Join(metrics, ", "))
			fmt.Println(line)
		}
		count++
	}
}
