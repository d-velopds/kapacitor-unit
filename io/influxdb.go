package io

import (
	"bytes"
	"net/http"

	"github.com/golang/glog"
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
		glog.Info("DEBUG:: Influxdb added [" + d + "] to " + url)
	}
	return nil
}

// Creates db and rp where tests will run
func (influxdb Influxdb) Setup(db string, rp string) error {
	glog.Info("DEGUB:: Influxdb setup ", db+":"+rp)
	// If no retention policy is defined, use "autogen"
	if rp == "" {
		rp = "autogen"
	}
	q := "q=CREATE DATABASE \"" + db + "\" WITH DURATION 1h REPLICATION 1 NAME \"" + rp + "\""
	baseUrl := influxdb.Host + "/query"
	_, err := influxdb.Client.Post(baseUrl, "application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(q)))
	if err != nil {
		return err
	}
	return nil
}

func (influxdb Influxdb) CleanUp(db string) error {
	q := "q=DROP DATABASE \"" + db + "\""
	baseUrl := influxdb.Host + "/query"
	_, err := influxdb.Client.Post(baseUrl, "application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(q)))
	if err != nil {
		return err
	}
	glog.Info("DEBUG:: Influxdb cleanup database ", q)
	return nil
}

func (influxdb Influxdb) CleanUpSeries(db string) error {
	q := "q=DROP SERIES FROM /.*/"
	baseUrl := influxdb.Host + "/query"
	_, err := influxdb.Client.Post(baseUrl, "application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(q)))
	if err != nil {
		return err
	}
	glog.Info("DEBUG:: Influxdb cleanup series ", q)
	return nil
}
