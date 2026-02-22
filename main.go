package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
)

// Hardcoded configuration
const (
	Port = "8080"
)

var (
	// Defined commands
	Commands = map[string]string{
		// 直接调用 shutdown.exe，不走 cmd
		"shutdown": "shutdown /s /t 0",
		"restart":  "shutdown /r /t 0",
		"abort":    "shutdown /a",
		// 睡眠和锁屏在服务模式下可能无效（受 Session 0 隔离限制）
		"sleep":    "rundll32.exe powrprof.dll,SetSuspendState 0,1,0",
		"lock":     "rundll32.exe user32.dll,LockWorkStation",
	}
	serviceLog service.Logger
	logFile    *os.File
)

// Program structures
type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	setupLogging()
	logToFile("Service is starting...")
	go p.run()
	return nil
}

func (p *program) run() {
	http.HandleFunc("/execute/", executeHandler)
	
	// Add a root handler to check if server is running
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Shutdown Tool is running! Use /execute/{command} to trigger actions."))
	})

	logMsg := fmt.Sprintf("Server starting on port %s... Available commands: %v", Port, Commands)
	logToFile(logMsg)
	
	err := http.ListenAndServe(":"+Port, nil)
	if err != nil {
		logToFile(fmt.Sprintf("Server failed to start: %v", err))
	}
}

func (p *program) Stop(s service.Service) error {
	logToFile("Service is stopping...")
	if logFile != nil {
		logFile.Close()
	}
	return nil
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	// Path should be /execute/{command}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 || pathParts[1] != "execute" {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	cmdName := pathParts[2]

	cmdStr, ok := Commands[cmdName]
	if !ok {
		http.Error(w, fmt.Sprintf("Command '%s' not found", cmdName), http.StatusNotFound)
		return
	}

	logToFile(fmt.Sprintf("Received request: %s -> %s", cmdName, cmdStr))

	// Execute the command
	var cmd *exec.Cmd
	
	// Split command and arguments
	parts := strings.Fields(cmdStr)
	head := parts[0]
	args := parts[1:]

	cmd = exec.Command(head, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Error executing command: %s\nOutput: %s", err, output)
		logToFile(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	successMsg := fmt.Sprintf("Command '%s' executed successfully.\nOutput: %s", cmdName, output)
	logToFile(successMsg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(successMsg))
}

func setupLogging() {
	// Get executable path
	ex, err := os.Executable()
	if err != nil {
		return
	}
	exPath := filepath.Dir(ex)
	logPath := filepath.Join(exPath, "shutdown-tool.log")

	f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	logFile = f
	
	// MultiWriter to print to console (if running interactively) and file
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
}

func logToFile(msg string) {
	log.Printf("%s\n", msg)
	if serviceLog != nil {
		// Also send to Windows Event Log if available
		serviceLog.Info(msg)
	}
}

func main() {
	svcConfig := &service.Config{
		Name:        "RemoteShutdown",
		DisplayName: "Remote Shutdown Service",
		Description: "Allows remote control of PC via HTTP requests.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	serviceLog, _ = s.Logger(nil)

	// Check for command line arguments (install, start, stop, uninstall)
	if len(os.Args) > 1 {
		// Setup logging for CLI operations too
		setupLogging()
		
		action := os.Args[1]
		err = service.Control(s, action)
		if err != nil {
			fmt.Printf("Action '%s' failed: %v\n", action, err)
		} else {
			fmt.Printf("Action '%s' succeeded.\n", action)
		}
		return
	}

	// Run the service
	// If run interactively, it executes here
	// If run as service, it executes here too
	err = s.Run()
	if err != nil {
		if serviceLog != nil {
			serviceLog.Error(err.Error())
		}
	}
}
