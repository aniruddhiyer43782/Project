package web

import (
	"backup-x/client"
	"backup-x/entity"
	"backup-x/util"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)


var startTime = time.Now()


var saveLimit = time.Duration(30 * time.Minute)


func Save(writer http.ResponseWriter, request *http.Request) {
	oldConf, _ := entity.GetConfigCache()
	conf := &entity.Config{}

	if oldConf.Password == "" {
		if time.Since(startTime) > saveLimit {
			writer.Write([]byte(fmt.Sprintf("Username and password must be set before %s, please restart backup-x", startTime.Add(saveLimit).Format("2006-01-02 15:04:05"))))
			return
		}
	}

	conf.EncryptKey = oldConf.EncryptKey
	if conf.EncryptKey == "" {
		encryptKey, err := util.GenerateEncryptKey()
		if err != nil {
			writer.Write([]byte("Failed to generate key"))
			return
		}
		conf.EncryptKey = encryptKey
	}

	
	conf.Username = strings.TrimSpace(request.FormValue("Username"))
	conf.Password = request.FormValue("Password")

	if conf.Username == "" || conf.Password == "" {
		writer.Write([]byte("Please enter login username/password"))
		return
	}
	if conf.Password != oldConf.Password {
		encryptPasswd, err := util.EncryptByEncryptKey(conf.EncryptKey, conf.Password)
		if err != nil {
			writer.Write([]byte("Encryption failed"))
			return
		}
		conf.Password = encryptPasswd
	}

	forms := request.PostForm
	for index, projectName := range forms["ProjectName"] {
		saveDays, _ := strconv.Atoi(forms["SaveDays"][index])
		saveDaysS3, _ := strconv.Atoi(forms["SaveDaysS3"][index])
		startTime, _ := strconv.Atoi(forms["StartTime"][index])
		period, _ := strconv.Atoi(forms["Period"][index])
		backupType, _ := strconv.Atoi(forms["BackupType"][index])
		enabled, _ := strconv.Atoi(forms["Enabled"][index])
		conf.BackupConfig = append(
			conf.BackupConfig,
			entity.BackupConfig{
				ProjectName: projectName,
				Command:     forms["Command"][index],
				SaveDays:    saveDays,
				SaveDaysS3:  saveDaysS3,
				StartTime:   startTime,
				Period:      period,
				Pwd:         forms["Pwd"][index],
				BackupType:  backupType,
				Enabled:     enabled,
			},
		)
	}

	for i := 0; i < len(conf.BackupConfig); i++ {
		if conf.BackupConfig[i].Pwd != "" &&
			(len(oldConf.BackupConfig) == 0 || conf.BackupConfig[i].Pwd != oldConf.BackupConfig[i].Pwd) {
			encryptPwd, err := util.EncryptByEncryptKey(conf.EncryptKey, conf.BackupConfig[i].Pwd)
			if err != nil {
				writer.Write([]byte("Encryption failed"))
				return
			}
			conf.BackupConfig[i].Pwd = encryptPwd
		}
	}

	// Webhook
	conf.WebhookURL = strings.TrimSpace(request.FormValue("WebhookURL"))
	conf.WebhookRequestBody = strings.TrimSpace(request.FormValue("WebhookRequestBody"))

	// S3
	conf.Endpoint = strings.TrimSpace(request.FormValue("Endpoint"))
	conf.AccessKey = strings.TrimSpace(request.FormValue("AccessKey"))
	conf.SecretKey = strings.TrimSpace(request.FormValue("SecretKey"))
	conf.BucketName = strings.TrimSpace(request.FormValue("BucketName"))
	conf.Region = strings.TrimSpace(request.FormValue("Region"))

	if conf.SecretKey != "" && conf.SecretKey != oldConf.SecretKey {
		secretKey, err := util.EncryptByEncryptKey(conf.EncryptKey, conf.SecretKey)
		if err != nil {
			writer.Write([]byte("Encryption failed"))
			return
		}
		conf.SecretKey = secretKey
	}

	
	err := conf.SaveConfig()

	
	if err == nil {
		conf.CreateBucketIfNotExist()
		if request.URL.Query().Get("backupAll") == "true" {
			go client.RunOnce()
		}
		if request.URL.Query().Get("backupIdx") != "" {
			idx, err := strconv.Atoi(request.URL.Query().Get("backupIdx"))
			if err == nil {
				go client.RunByIdx(idx)
			} else {
				log.Println("Index number is incorrect" + request.URL.Query().Get("backupIdx"))
			}
		}
		
		client.StopRunLoop()
		go client.RunLoop(100 * time.Millisecond)
	}

	
	if err == nil {
		writer.Write([]byte("ok"))
	} else {
		writer.Write([]byte(err.Error()))
	}

}
