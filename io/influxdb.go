package io

import (
	"bytes"
	"github.com/golang/glog"
	"net/http"
)

// Influxdb service configurations
type Influxdb struct {
	Host   string
	Client http.Client
}

func NewInfluxdb(host string) Influxdb {
	return Influxdb{
		host,
		http.Client{},
	}
}

// Adds test data to influxdb
func (influxdb Influxdb) Data(data []string, db string, rp string) error {
	url := influxdb.Host + influxdb_write + "db=" + db + "&rp=" + rp
	for _, d := range data {
		_, err := influxdb.Client.Post(url, "application/x-www-form-urlencoded",
			bytes.NewBuffer([]byte(d)))
		if err != nil {
			return err
		}
		glog.Info("Influxdb added data:: ", d)
		glog.Info("Influxdb added TO URL:: ", url)

	}
	return nil
}

// Creates db and rp where tests will run
func (influxdb Influxdb) Setup(db string, rp string) error {
	glog.Info("Influxdb Setup:: ", db+":"+rp)
	err := influxdb.post("q=CREATE DATABASE "+db)
	if err != nil {
		return err
	}
	// rp autogen is created by default
	if rp != "autogen" {
		err = influxdb.post("q=CREATE RETENTION POLICY "+rp+" ON "+db)
		if err != nil {
			return err
		}		
	}
	return nil
}

func (influxdb Influxdb) post(q string) error {
	url := influxdb.Host + "/query"
	_, err := influxdb.Client.Post(url, "application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(q)))
	if err != nil {
		return err
	}
	glog.Info("Influxdb POST:: ", q)
	return nil
}

