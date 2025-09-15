package client

import (
	"backup-x/entity"
	"backup-x/util"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

const minFileSize = 1000

// backupLooper handles backup loops
type backupLooper struct {
	Wg      sync.WaitGroup
	Tickers []*time.Ticker
}

var bl = &backupLooper{Wg: sync.WaitGroup{}}

// RunLoop starts the backup loop
func RunLoop(firstDelay time.Duration) {
	conf, err := entity.GetConfigCache()
	if err != nil {
		return
	}

	time.Sleep(firstDelay)

	// Clear existing tickers
	bl.Tickers = []*time.Ticker{}

	for _, backupConf := range conf.BackupConfig {
		if !backupConf.NotEmptyProject() {
			continue
		}

		if backupConf.Enabled != 0 {
			log.Println(backupConf.ProjectName + " project is disabled")
			continue
		}

		if !backupConf.CheckPeriod() {
			log.Println(backupConf.ProjectName + " project has an invalid period")
			continue
		}

		delay := util.GetDelaySeconds(backupConf.StartTime)
		ticker := time.NewTicker(delay)
		log.Printf("%s project will run after %.1f hours\n", backupConf.ProjectName, delay.Hours())

		bl.Wg.Add(1)
		go func(backupConf entity.BackupConfig) {
			defer bl.Wg.Done()
			for {
				<-ticker.C
				run(conf, backupConf)
				ticker.Reset(time.Minute * time.Duration(backupConf.Period))
				log.Printf("%s project will wait %d minutes before next run\n", backupConf.ProjectName, backupConf.Period)
			}
		}(backupConf)
		bl.Tickers = append(bl.Tickers, ticker)
	}

	bl.Wg.Wait()
}

// StopRunLoop stops all running backup loops
func StopRunLoop() {
	for _, ticker := range bl.Tickers {
		if ticker != nil {
			ticker.Stop()
		}
	}
}

// RunOnce runs all backups once
func RunOnce() {
	conf, err := entity.GetConfigCache()
	if err != nil {
		return
	}

	for _, backupConf := range conf.BackupConfig {
		run(conf, backupConf)
	}
}

// RunByIdx runs a backup for a specific index
func RunByIdx(idx int) {
	conf, err := entity.GetConfigCache()
	if err != nil {
		return
	}

	run(conf, conf.BackupConfig[idx])
}

// run executes a backup task
func run(conf entity.Config, backupConf entity.BackupConfig) {
	if backupConf.NotEmptyProject() && backupConf.Enabled == 0 {
		err := prepare(backupConf)
		if err != nil {
			log.Println(err)
			return
		}

		// Perform backup
		outFileName, err := backup(backupConf, conf.EncryptKey, conf.S3Config)
		result := entity.BackupResult{ProjectName: backupConf.ProjectName, Result: "Failed"}
		if err == nil {
			// Webhook
			if outFileName != nil {
				result.FileName = outFileName.Name()
				result.FileSize = fmt.Sprintf("%d MB", outFileName.Size()/1000/1000)
				// Upload to S3 if configured
				if conf.S3Config.CheckNotEmpty() {
					go conf.S3Config.UploadFile(backupConf.GetProjectPath() + string(os.PathSeparator) + outFileName.Name())
				}
			}
			result.Result = "Success"
		}
		conf.ExecWebhook(result)
	}
}

// prepare creates project folder
func prepare(backupConf entity.BackupConfig) (err error) {
	os.MkdirAll(backupConf.GetProjectPath(), 0750)
	return
}

// backup executes the backup shell command
func backup(backupConf entity.BackupConfig, encryptKey string, s3Conf entity.S3Config) (outFileName os.FileInfo, err error) {
	projectName := backupConf.ProjectName
	log.Printf("Backing up project: %s ...", projectName)

	todayString := time.Now().Format(util.FileNameFormatStr)
	shellString := strings.ReplaceAll(backupConf.Command, "#{DATE}", todayString)

	// Decrypt password
	pwd := ""
	if backupConf.Pwd != "" {
		pwd, err = util.DecryptByEncryptKey(encryptKey, backupConf.Pwd)
		if err != nil {
			err = fmt.Errorf("decryption failed")
			log.Println(err)
			return nil, err
		}
	}

	// Decrypt S3 secret key
	secretKey := ""
	if s3Conf.SecretKey != "" {
		secretKey, err = util.DecryptByEncryptKey(encryptKey, s3Conf.SecretKey)
		if err != nil {
			err = fmt.Errorf("decryption failed")
			log.Println(err)
			return nil, err
		}
	}

	// Replace placeholders
	shellString = strings.ReplaceAll(shellString, "#{PWD}", pwd)
	shellString = strings.ReplaceAll(shellString, "#{AccessKey}", s3Conf.AccessKey)
	shellString = strings.ReplaceAll(shellString, "#{SecretKey}", secretKey)
	shellString = strings.ReplaceAll(shellString, "#{Endpoint}", s3Conf.Endpoint)
	shellString = strings.ReplaceAll(shellString, "#{BucketName}", s3Conf.BucketName)

	// Create shell file
	var shellName string
	if runtime.GOOS == "windows" {
		shellName = time.Now().Format("shell-"+util.FileNameFormatStr+"-") + "backup.bat"
	} else {
		shellString = strings.ReplaceAll(shellString, "\r\n", "\n") // convert windows line endings
		shellName = time.Now().Format("shell-"+util.FileNameFormatStr+"-") + "backup.sh"
	}

	shellFile, err := os.Create(backupConf.GetProjectPath() + string(os.PathSeparator) + shellName)
	shellFile.Chmod(0700)
	if err == nil {
		shellFile.WriteString(shellString)
		shellFile.Close()
	} else {
		log.Println("Error creating shell file: ", err)
	}

	// Execute shell
	var shell *exec.Cmd
	if runtime.GOOS == "windows" {
		shell = exec.Command("cmd", "/c", shellName)
	} else {
		shell = exec.Command("bash", shellName)
	}
	shell.Dir = backupConf.GetProjectPath()
	outputBytes, err := shell.CombinedOutput()
	if len(outputBytes) > 0 {
		if util.IsGBK(outputBytes) {
			outputBytes, _ = util.GbkToUtf8(outputBytes)
		}
		log.Printf("<span style='color: #7983f5;font-weight: bold;'>%s</span> Shell output: <span class='click-layer' onclick='showLayer(this)' tip=\"%s\" style='cursor: pointer; color: #4a3a3a; font-weight: bold; border: 2px dashed;'>Click to view</span>\n", backupConf.ProjectName, util.EscapeShell(string(outputBytes)))
	} else {
		log.Printf("Shell output is empty\n")
	}

	// Check if backup was successful
	if err == nil {
		outFileName, err = findBackupFile(backupConf, todayString)
		if backupConf.BackupType == 0 {
			// Database backup
			if err != nil {
				log.Println(err)
			} else if outFileName.Size() >= minFileSize {
				log.Printf("Successfully backed up project: %s, file: %s\n", projectName, outFileName.Name())
			} else {
				err = fmt.Errorf("%s backup file is smaller than %d bytes, current: %d bytes", projectName, minFileSize, outFileName.Size())
				log.Println(err)
			}
		} else {
			// File sync type
			err = nil
		}
	} else {
		err = fmt.Errorf("Failed to execute backup shell: %s", util.EscapeShell(string(outputBytes)))
		log.Println(err)
	}

	// Remove shell file
	os.Remove(shellFile.Name())

	return
}

// findBackupFile searches for backup file containing today's date
func findBackupFile(backupConf entity.BackupConfig, todayString string) (backupFile os.FileInfo, err error) {
	files, err := ioutil.ReadDir(backupConf.GetProjectPath())
	for _, file := range files {
		if strings.Contains(file.Name(), todayString) && !strings.HasPrefix(file.Name(), "shell-") {
			backupFile = file
			return
		}
	}

	err = fmt.Errorf("Project %s has no output file containing %s", backupConf.ProjectName, todayString)
	return
}
