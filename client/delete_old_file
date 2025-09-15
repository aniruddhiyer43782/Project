package client

import (
	"backup-x/entity"
	"backup-x/util"
	"log"
	"os"
	"time"
)

// DeleteOldBackup runs a loop to delete old backups
func DeleteOldBackup() {
	for {
		delay := util.GetDelaySeconds(2)
		log.Printf("Deleting expired backup files will run in %.1f hours\n", delay.Hours())
		time.Sleep(delay)

		conf, err := entity.GetConfigCache()
		if err != nil {
			return
		}

		for _, backupConf := range conf.BackupConfig {
			// Skip empty or disabled projects
			if !backupConf.NotEmptyProject() || backupConf.Enabled == 1 {
				continue
			}
			// Delete old local files
			deleteLocalOlderFiles(backupConf)

			// Delete old files from object storage (S3)
			deleteS3OlderFiles(conf.S3Config, backupConf)
		}
	}
}

// deleteLocalOlderFiles deletes expired local backup files
func deleteLocalOlderFiles(backupConf entity.BackupConfig) {
	backupFiles, err := os.ReadDir(backupConf.GetProjectPath())
	if err != nil {
		log.Printf("Failed to read local directory for project %s! ERR: %s\n", backupConf.ProjectName, err)
		return
	}
	if backupConf.SaveDays <= 0 {
		log.Printf("Local retention days setting for project %s is invalid", backupConf.ProjectName)
		return
	}

	backupFileNames := make([]string, 0)
	for _, backupFile := range backupFiles {
		if !backupFile.IsDir() {
			info, err := backupFile.Info()
			if err == nil {
				if info.Size() >= minFileSize {
					backupFileNames = append(backupFileNames, backupFile.Name())
				} else {
					if util.IsFileNameDate(backupFile.Name()) {
						log.Printf("Backup file size %d bytes is less than minimum %d, deleting file: %s", info.Size(), minFileSize, backupConf.GetProjectPath()+string(os.PathSeparator)+backupFile.Name())
						os.Remove(backupConf.GetProjectPath() + string(os.PathSeparator) + backupFile.Name())
					}
				}
			}
		}
	}

	tobeDeleteFiles := util.FileNameBeforeDays(backupConf.SaveDays, backupFileNames, backupConf.ProjectName)

	for i := 0; i < len(tobeDeleteFiles); i++ {
		err := os.Remove(backupConf.GetProjectPath() + string(os.PathSeparator) + tobeDeleteFiles[i])
		if err == nil {
			log.Printf("Successfully deleted expired local file: %s", backupConf.ProjectName+string(os.PathSeparator)+tobeDeleteFiles[i])
		} else {
			log.Printf("Failed to delete expired local file: %s, ERR: %s", backupConf.ProjectName+string(os.PathSeparator)+tobeDeleteFiles[i], err)
		}
	}
}

// deleteS3OlderFiles deletes expired files from object storage
func deleteS3OlderFiles(s3Conf entity.S3Config, backupConf entity.BackupConfig) {
	if !s3Conf.CheckNotEmpty() {
		return
	}
	if backupConf.SaveDaysS3 <= 0 {
		log.Printf("S3 retention days setting for project %s is invalid", backupConf.ProjectName)
		return
	}

	fileNames, err := s3Conf.ListFiles(backupConf.GetProjectPath())
	if err != nil {
		log.Printf("Failed to read S3 directory for project %s! ERR: %s\n", backupConf.ProjectName, err)
		return
	}

	tobeDeleteFiles := util.FileNameBeforeDays(backupConf.SaveDaysS3, fileNames, backupConf.ProjectName)

	for i := 0; i < len(tobeDeleteFiles); i++ {
		err := s3Conf.DeleteFile(tobeDeleteFiles[i])
		if err == nil {
			log.Printf("Successfully deleted expired file from S3: %s", tobeDeleteFiles[i])
		} else {
			log.Printf("Failed to delete expired file from S3: %s, ERR: %s", tobeDeleteFiles[i], err)
		}
	}
}
