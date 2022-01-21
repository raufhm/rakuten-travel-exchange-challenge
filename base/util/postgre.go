package util

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/lib/pq"

	src "github.com/raufhm/rtxTest/base"
	"github.com/sirupsen/logrus"
)

const (
	XMLRate = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
)

func GetPostgreClient() (db *sql.DB, err error) {
	db, err = sql.Open("postgres", "postgresql://root:secret@localhost:5432/rtx_test?sslmode=disable")
	if err != nil {
		logrus.Error(err)
		return
	}
	err = db.Ping()
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func BulkImport(db *sql.DB) (err error) {

	// load xml data
	results, err := LoadXmLData(XMLRate)
	if err != nil {
		logrus.Error(err)
		return
	}

	if len(results) == 0 {
		err = fmt.Errorf("an error occured")
		logrus.Error(err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		logrus.Error(err)
		return
	}

	// bulk import to database
	s, err := tx.Prepare(pq.CopyIn("envelopes", "time", "currency", "rate"))
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, record := range results {
		_, err = s.Exec(record.Time, record.Currency, record.Rate)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
	_, err = s.Exec()
	if err != nil {
		logrus.Error(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		return
	}

	return
}

func LoadXmLData(url string) (output []*src.Output, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %v", err)
	}

	var result src.Envelope
	err = xml.Unmarshal(data, &result)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	output = []*src.Output{}
	for _, cc := range result.Cube.Cube {
		for _, c := range cc.Cube {
			times, _ := time.Parse("2006-01-02", cc.Time)
			rate, _ := strconv.ParseFloat(c.Rate, 64)
			output = append(output, &src.Output{
				Time:     times,
				Currency: c.Currency,
				Rate:     rate,
			})
		}
	}

	return
}
