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
)

var head []string

func main() {
	file := flag.String("file", "", "File to process.")
	sid := flag.Int("sid", 0, "Sensor ID.")
	flag.Parse()

	sql := `
insert into sensorlogs (sensorid, metrics) values (%d, '%s');
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
		count++

		var metrics []string
		for i, k := range head {
			// '{{temp=>80}, {baro=>100}}'
			metric := fmt.Sprintf("{%s=>%s}", k, record[i])
			metrics = append(metrics, metric)
		}
		line := fmt.Sprintf(sql, *sid, "{"+strings.Join(metrics, ", ")+"}")
		fmt.Println(line)
	}
}
