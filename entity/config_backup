package entity

// BackupConfig represents a backup configuration
type BackupConfig struct {
	ProjectName string // Project name
	Command     string // Command to run
	SaveDays    int    // Number of days to keep local backups
	SaveDaysS3  int    // Number of days to keep backups in object storage (S3)
	StartTime   int    // Start time (0-23)
	Period      int    // Interval period (minutes)
	Pwd         string // Password
	BackupType  int    // Backup type: 0 = Database backup, 1 = File sync
	Enabled     int    // Whether enabled: 0 = Enabled, 1 = Disabled
}

// GetProjectPath returns the path for the project
func (backupConfig *BackupConfig) GetProjectPath() string {
	return parentSavePath + "/" + backupConfig.ProjectName
}

// NotEmptyProject checks if the project is not empty
func (backupConfig *BackupConfig) NotEmptyProject() bool {
	return backupConfig.Command != "" && backupConfig.ProjectName != ""
}

// CheckPeriod validates the start time and interval period
func (backupConfig *BackupConfig) CheckPeriod() bool {
	return backupConfig.StartTime >= 0 && backupConfig.StartTime < 24 && backupConfig.Period > 0
}
