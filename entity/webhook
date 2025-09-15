package entity

import (
	"backup-x/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Webhook represents webhook configuration
type Webhook struct {
	WebhookURL         string
	WebhookRequestBody string
}

// ExecWebhook executes the webhook with the given backup result
func (webhook Webhook) ExecWebhook(result BackupResult) {

	if webhook.WebhookURL != "" {
		method := "GET"
		postData := ""
		contentType := "application/x-www-form-urlencoded"

		if webhook.WebhookRequestBody != "" {
			method = "POST"
			postData = webhook.replaceBody(result)
			if json.Valid([]byte(postData)) {
				contentType = "application/json"
			}
		}

		requestURL := webhook.replaceURL(result)
		u, err := url.Parse(requestURL)
		if err != nil {
			log.Println("Invalid URL in webhook configuration")
			return
		}

		req, err := http.NewRequest(method, fmt.Sprintf("%s://%s%s?%s", u.Scheme, u.Host, u.Path, u.Query().Encode()), strings.NewReader(postData))
		if err != nil {
			log.Println("Error creating webhook request, Err:", err)
			return
		}
		req.Header.Add("content-type", contentType)

		client := http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		body, err := util.GetHTTPResponseOrg(resp, requestURL, err)
		if err == nil {
			log.Println(fmt.Sprintf("Webhook call succeeded, response: %s", string(body)))
		} else {
			log.Println(fmt.Sprintf("Webhook call failed, Err: %s", err))
		}
	}
}

// replaceURL replaces placeholders in the webhook URL with actual values
func (webhook Webhook) replaceURL(result BackupResult) (newURL string) {
	newURL = strings.ReplaceAll(webhook.WebhookURL, "#{projectName}", result.ProjectName)
	newURL = strings.ReplaceAll(newURL, "#{fileName}", result.FileName)
	newURL = strings.ReplaceAll(newURL, "#{fileSize}", result.FileSize)
	newURL = strings.ReplaceAll(newURL, "#{result}", result.Result)
	return newURL
}

// replaceBody replaces placeholders in the webhook request body with actual values
func (webhook Webhook) replaceBody(result BackupResult) (newBody string) {
	newBody = strings.ReplaceAll(webhook.WebhookRequestBody, "#{projectName}", result.ProjectName)
	newBody = strings.ReplaceAll(newBody, "#{fileName}", result.FileName)
	newBody = strings.ReplaceAll(newBody, "#{fileSize}", result.FileSize)
	newBody = strings.ReplaceAll(newBody, "#{result}", result.Result)
	return newBody
}
