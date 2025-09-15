package main

import (
	"backup-x/client"
	"backup-x/util"
	"backup-x/web"
	"embed"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/kardianos/service"
)

// Listening address
var listen = flag.String("l", ":9977", "Listening address")

// Service management
var serviceType = flag.String("s", "", "Service management: supports install, uninstall")

// Default backup path is the current working directory
var backupDirDefault, _ = os.Getwd()

// Custom backup directory path
var backupDir = flag.String("d", backupDirDefault, "Custom backup directory path")

//go:embed static
var staticEmbededFiles embed.FS

//go:embed favicon.ico
var faviconEmbededFile embed.FS

// Version
var version = "DEV"

func main() {
	flag.Parse()

	if _, err := net.ResolveTCPAddr("tcp", *listen); err != nil {
		log.Fatalf("Error resolving listening address: %s", err)
	}

	os.Setenv(web.VersionEnv, version)

	switch *serviceType {
	case "install":
		installService()
	case "uninstall":
		uninstallService()
	default:
		if util.IsRunInDocker() {
			run(100 * time.Millisecond)
		} else {
			s := getService()
			status, _ := s.Status()
			if status != service.StatusUnknown {
				// Run as a service
				s.Run()
			} else {
				// Run in non-service mode
				switch s.Platform() {
				case "windows-service":
					log.Println("You can install the service using: .\\backup-x.exe -s install")
				default:
					log.Println("You can install the service using: ./backup-x -s install")
				}
				run(100 * time.Millisecond)
			}
		}
	}
}

// Serve static files
func staticFsFunc(writer http.ResponseWriter, request *http.Request) {
	http.FileServer(http.FS(staticEmbededFiles)).ServeHTTP(writer, request)
}

// Serve favicon
func faviconFsFunc(writer http.ResponseWriter, request *http.Request) {
	http.FileServer(http.FS(faviconEmbededFile)).ServeHTTP(writer, request)
}

// run starts the backup server and loop
func run(firstDelay time.Duration) {
	// Serve static files
	http.HandleFunc("/static/", web.BasicAuth(staticFsFunc))
	http.HandleFunc("/favicon.ico", web.BasicAuth(faviconFsFunc))

	// Web handlers
	http.HandleFunc("/", web.BasicAuth(web.WritingConfig))
	http.HandleFunc("/save", web.BasicAuth(web.Save))
	http.HandleFunc("/logs", web.BasicAuth(web.Logs))
	http.HandleFunc("/clearLog", web.BasicAuth(web.ClearLog))
	http.HandleFunc("/webhookTest", web.BasicAuth(web.WebhookTest))

	// Change working directory
	os.Chdir(*backupDir)

	// Run backup loops
	go client.DeleteOldBackup()
	go client.RunLoop(firstDelay)

	err := http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Println("Error starting server, please check if the port is already in use:", err)
		time.Sleep(time.Minute)
	}
}

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work asynchronously.
	go p.run()
	return nil
}

func (p *program) run() {
	// Service starts after 20-second delay to wait for network readiness
	run(20 * time.Second)
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block
	return nil
}

// getService creates a new service configuration
func getService() service.Service {
	options := make(service.KeyValue)
	var depends []string

	// Ensure service waits for network before starting
	switch service.ChosenSystem().String() {
	case "windows-service":
		// Set Windows service to delayed automatic start
		options["DelayedAutoStart"] = true
	default:
		// Add network dependency for Systemd
		depends = append(depends, "Requires=network.target", "After=network-online.target")
	}

	svcConfig := &service.Config{
		Name:         "backup-x",
		DisplayName:  "backup-x",
		Description:  "Database/File backup tool with web interface",
		Arguments:    []string{"-l", *listen, "-d", *backupDir},
		Dependencies: depends,
		Option:       options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}

// uninstallService uninstalls the service
func uninstallService() {
	s := getService()
	status, _ := s.Status()
	if status != service.StatusUnknown {
		s.Stop()
		if err := s.Uninstall(); err == nil {
			log.Println("backup-x service uninstalled successfully!")
		} else {
			log.Printf("Failed to uninstall backup-x service, ERR: %s\n", err)
		}
	} else {
		log.Println("backup-x service is not installed")
	}
}

// installService installs the service
func installService() {
	s := getService()
	status, err := s.Status()
	if err != nil && status == service.StatusUnknown {
		// Service unknown, create service
		if err = s.Install(); err == nil {
			s.Start()
			log.Println("backup-x service installed successfully! The program will run continuously, including after restart.")
			return
		}

		log.Printf("Failed to install backup-x service, ERR: %s\n", err)
		switch s.Platform() {
		case "windows-service":
			log.Println("Make sure to run: .\\backup-x.exe -s install")
		default:
			log.Println("Make sure to run: ./backup-x -s install")
		}
	}

	if status != service.StatusUnknown {
		log.Println("backup-x service is already installed, no need to install again")
	}
}
