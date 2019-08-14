package test

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestValidateRecAndData(t *testing.T) {
	r := Result{}
	d := []string{"data1", "data2"}
	tst := Test{}

	tst.Data = d
	tst.RecId = "e24db07d-1646-4bb3-a445-828f5049bea0"
	tst.Result = r

	tst.Validate()

	if tst.Result.Error != true {
		t.Error("Test initialized with recording_id and test must be invalid")
	}
}

func TestValidateRecOk(t *testing.T) {
	r := Result{}
	tst := Test{}

	tst.RecId = "e24db07d-1646-4bb3-a445-828f5049bea0"
	tst.Result = r

	tst.Validate()

	if tst.Result.Error != false {
		t.Error("Test initialized only with recording_id must be valid")
	}
}

func TestValidateDataOk(t *testing.T) {
	r := Result{}
	tst := Test{}
	d := []string{"data1", "data2"}

	tst.Data = d
	tst.Result = r

	tst.Validate()

	if tst.Result.Error != false {
		t.Error("Test initialized only with data must be valid")
	}
}

func TestValidateRecNotOk(t *testing.T) {
	tst := NewTest()

	tst.Data = []string{"data1"}
	tst.Result = Result{}
	tst.RecId = "some_id"

	tst.Validate()

	if tst.Result.Error != true {
		t.Error("Test configuration with recording id and protocol line data is invalid")
	}
}

func extractTimestamp(s string) (int64, error) {
	parts := strings.Split(s, " ")
	stampString := parts[len(parts)-1]
	stamp, err := strconv.Atoi(stampString)
	if err != nil {
		return 0, err
	}
	return int64(stamp), nil
}

func TestCreatesNegativeMinuteDiffTimestamps(t *testing.T) {
	tst := NewTest()

	tst.Data = []string{
		"cpu,type=cpu-all,host=a01 usage_idle=82 -1m",
		"cpu,type=cpu-all,host=a01 test=\"hey\" -5m",
	}

	now := time.Now().Unix()
	tst.CreateTimestamps()

	stamp1, err := extractTimestamp(tst.Data[0])
	if err != nil {
		t.Error("Could not extract last part of sample 1")
	}
	diff := now - stamp1
	if diff < 58 || diff > 62 {
		t.Error("Converted timestamp of sample 1 not in expected range")
	}

	stamp2, err := extractTimestamp(tst.Data[1])
	if err != nil {
		t.Error("Could not extract last part of sample 1")
	}
	diff = now - stamp2
	if diff < 298 || diff > 302 {
		t.Error("Converted timestamp of sample 2 not in expected range")
	}
}

func TestCreatesPositiveMinuteDiffTimestamps(t *testing.T) {
	tst := NewTest()

	tst.Data = []string{
		"cpu,type=cpu-all,host=a01 usage_idle=82 +1m",
		"cpu,type=cpu-all,host=a01 test=\"hey\" +3m",
	}

	now := time.Now().Unix()
	tst.CreateTimestamps()

	stamp1, err := extractTimestamp(tst.Data[0])
	if err != nil {
		t.Error("Could not extract last part of sample 1")
	}
	diff := now - stamp1
	if diff > -58 || diff < -62 {
		t.Error("Converted timestamp of sample 1 not in expected range")
	}

	stamp2, err := extractTimestamp(tst.Data[1])
	if err != nil {
		t.Error("Could not extract last part of sample 1")
	}
	diff = now - stamp2
	if diff > -178 || diff < -182 {
		t.Error("Converted timestamp of sample 2 not in expected range")
	}
}
