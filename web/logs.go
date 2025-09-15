package web

import (
	"io"
	"log"
	"net/http"
	"os"
)


type MemoryLogs struct {
	MaxNum int      
	Logs   []string 
}

func (mlogs *MemoryLogs) Write(p []byte) (n int, err error) {
	mlogs.Logs = append(mlogs.Logs, string(p))
	
	if len(mlogs.Logs) > mlogs.MaxNum {
		mlogs.Logs = mlogs.Logs[len(mlogs.Logs)-mlogs.MaxNum:]
	}
	return len(p), nil
}

var mlogs = &MemoryLogs{MaxNum: 50}


func init() {
	log.SetOutput(io.MultiWriter(mlogs, os.Stdout))
	// log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// Logs web
func Logs(writer http.ResponseWriter, request *http.Request) {
	for _, log := range mlogs.Logs {
		writer.Write([]byte(log))
		writer.Write([]byte("<br/>"))
	}
}

// ClearLog
func ClearLog(writer http.ResponseWriter, request *http.Request) {
	mlogs.Logs = mlogs.Logs[:0]
}
