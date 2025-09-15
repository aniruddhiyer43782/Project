package util

import (
	"strconv"
	"testing"
	"time"
)

// TestFileNameUtil
func TestFileNameUtil(t *testing.T) {
	const days = 10
	beforeTomorrow, _ := time.ParseDuration("-" + strconv.Itoa((days+1)*24) + "h")
	fileNames := []string{
		"a2020-10-10-11-12b.sql",
		"test2021-10-10-11-12test.sql",
		time.Now().Add(beforeTomorrow).Format(FileNameFormatStr) + ".sql",
		time.Now().Format(FileNameFormatStr) + ".sql",
	}
	deleteFiles := FileNameBeforeDays(days, fileNames, "test")

	if len(fileNames) != len(deleteFiles)+1 {
		t.Error("TestFileNameUtil Test failed!")
	}
}

// TestFileNameUtilAll
func TestFileNameUtilAll(t *testing.T) {
	const days = 10
	beforeTomorrow, _ := time.ParseDuration("-" + strconv.Itoa((days+1)*24) + "h")
	fileNames := []string{
		"a2020-10-10-11-12b.sql",
		"test2021-10-10-11-12test.sql",
		time.Now().Add(beforeTomorrow).Format(FileNameFormatStr) + ".sql",
	}
	deleteFiles := FileNameBeforeDays(days, fileNames, "test")
	if len(deleteFiles) != 0 {
		t.Error("TestFileNameUtilAll Test failed!")
	}
}
