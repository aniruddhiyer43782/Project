package web

import (
	"backup-x/entity"
	"backup-x/util"
	"bytes"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"
)

// ViewFunc func
type ViewFunc func(http.ResponseWriter, *http.Request)

type loginDetect struct {
	FailTimes int
}

var ld = &loginDetect{}

// BasicAuth basic auth
func BasicAuth(f ViewFunc) ViewFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conf, _ := entity.GetConfigCache()

	
		if conf.Username == "" && conf.Password == "" {
			
			f(w, r)
			return
		}

		if ld.FailTimes >= 5 {
			log.Printf("%s login failed more than 5 times! Response delayed by 5 minutes\n", r.RemoteAddr)
			time.Sleep(5 * time.Minute)
			if ld.FailTimes >= 5 {
				ld.FailTimes = 0
			}
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		
		basicAuthPrefix := "Basic "

		
		auth := r.Header.Get("Authorization")
		
		if strings.HasPrefix(auth, basicAuthPrefix) {
			
			payload, err := base64.StdEncoding.DecodeString(
				auth[len(basicAuthPrefix):],
			)
			if err == nil {
				pair := bytes.SplitN(payload, []byte(":"), 2)
				pwd, _ := util.DecryptByEncryptKey(conf.EncryptKey, conf.Password)
				if len(pair) == 2 &&
					bytes.Equal(pair[0], []byte(conf.Username)) &&
					bytes.Equal(pair[1], []byte(pwd)) {
					ld.FailTimes = 0
					
					f(w, r)
					return
				}
			}

			ld.FailTimes = ld.FailTimes + 1
			log.Printf("%s login failed!\n", r.RemoteAddr)
		}

		
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("%s requested login!\n", r.RemoteAddr)
	}
}
