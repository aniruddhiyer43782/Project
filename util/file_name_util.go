package util

import (
	"log"
	"regexp"
	"strconv"
	"time"
)

const FileNameFormatStr = "2006-01-02-15-04"
const fileNameRegStr = `([\d]{4})-([\d]{2})-([\d]{2})-([\d]{2})-([\d]{2})`

// FileNameBeforeDays 
func FileNameBeforeDays(days int, fileNames []string, projectName string) []string {
	oldFiles := make([]string, 0)
	// 2006-01-02-15-04
	fileRegxp := regexp.MustCompile(fileNameRegStr)
	subDuration, _ := time.ParseDuration("-" + strconv.Itoa(days*24) + "h")
	before := time.Now().Add(subDuration)
	for i := 0; i < len(fileNames); i++ {
		dateString := fileRegxp.FindString(fileNames[i])
		if dateString != "" {
			if fileTime, err := time.Parse(FileNameFormatStr, dateString); err == nil && fileTime.Before(before) {
				oldFiles = append(oldFiles, fileNames[i])
			}
		}

	}
	
	if len(oldFiles) > 0 && len(oldFiles)-len(fileNames) >= 0 {
		log.Printf("Project %s expired files include all files, no deletion will be performed!\n", projectName)
		return []string{}
	}
	return oldFiles
}

// FileNameDate 
func IsFileNameDate(fileName string) bool {
	// 2006-01-02-15-04
	fileRegxp := regexp.MustCompile(fileNameRegStr)
	return fileRegxp.FindString(fileName) != ""
}
