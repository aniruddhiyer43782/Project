package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// GetHTTPResponse processes an HTTP response and unmarshals the JSON body into result.
func GetHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
	body, err := GetHTTPResponseOrg(resp, url, err)

	if err == nil {
		// log.Println(string(body))
		err = json.Unmarshal(body, &result)

		if err != nil {
			log.Printf("Failed to parse JSON response from %s! ERROR: %s\n", url, err)
		}
	}

	return err
}

// GetHTTPResponseOrg processes an HTTP response and returns the raw body as bytes.
func GetHTTPResponseOrg(resp *http.Response, url string, err error) ([]byte, error) {
	if err != nil {
		log.Printf("Request to %s failed! ERROR: %s\n", url, err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Failed to read response body from %s! ERROR: %s\n", url, err)
	}

	// Any status code 300 or above is considered an error
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("Request to %s failed! Response body: %s, Status code: %d\n",
			url, string(body), resp.StatusCode)
		log.Println(errMsg)
		err = fmt.Errorf(errMsg)
	}

	return body, err
}
