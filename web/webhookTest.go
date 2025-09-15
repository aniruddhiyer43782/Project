package web

import (
	"backup-x/entity"
	"log"
	"net/http"
	"strings"
)


func WebhookTest(writer http.ResponseWriter, request *http.Request) {
	url := strings.TrimSpace(request.FormValue("URL"))
	requestBody := strings.TrimSpace(request.FormValue("RequestBody"))
	if url != "" {
		wb := entity.Webhook{WebhookURL: url, WebhookRequestBody: requestBody}
		wb.ExecWebhook(entity.BackupResult{ProjectName: "Simulation test", FileName: "2021-11-11_01_01.sql", FileSize: "100 MB", Result: "Success"})
	} else {
		log.Println("Please enter the Webhook URL")
	}

}
