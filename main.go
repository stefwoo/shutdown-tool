package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kardianos/service"
	"gopkg.in/yaml.v3"
)

// Config represents the structure of the config.yaml file
type Config struct {
	Commands map[string]string `yaml:"commands"`
	Port     string            `yaml:"port"`
}

var (
	config      Config
	serviceLog  service.Logger
)

// Program structures
type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	// Determine executable directory for config file
	// This is crucial because when running as a service, the CWD is usually System32
	ex, err := os.Executable()
	if err != nil {
		if serviceLog != nil {
			serviceLog.Error(fmt.Sprintf("Failed to get executable path: %v", err))
		}
		ex = "."
	}
	exPath := filepath.Dir(ex)
	configFile := filepath.Join(exPath, "config.yaml")

	// Load configuration
	data, err := os.ReadFile(configFile)
	if err != nil {
		if serviceLog != nil {
			serviceLog.Warning(fmt.Sprintf("Could not read config at %s: %v. Using defaults.", configFile, err))
		}
		// Default config
		config = Config{
			Commands: map[string]string{
				"shutdown": "shutdown /s /t 0",
			},
			Port: "8080",
		}
	} else {
		if err := yaml.Unmarshal(data, &config); err != nil {
			if serviceLog != nil {
				serviceLog.Error(fmt.Sprintf("Failed to parse config: %v", err))
			}
		}
	}

	// Set default port if not specified
	if config.Port == "" {
		config.Port = "8080"
	}

	http.HandleFunc("/execute/", executeHandler)
	
	// Add a root handler to check if server is running
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Shutdown Tool is running! Use /execute/{command} to trigger actions."))
	})

	logMsg := fmt.Sprintf("Server starting on port %s... Available commands: %v", config.Port, config.Commands)
	if serviceLog != nil {
		serviceLog.Info(logMsg)
	} else {
		fmt.Println(logMsg)
	}
	
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		if serviceLog != nil {
			serviceLog.Error(fmt.Sprintf("Server failed to start: %v", err))
		} else {
			log.Fatal("Server failed to start:", err)
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Any cleanup work goes here
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

	cmdStr, ok := config.Commands[cmdName]
	if !ok {
		http.Error(w, fmt.Sprintf("Command '%s' not found", cmdName), http.StatusNotFound)
		return
	}

	msg := fmt.Sprintf("Executing command: %s (%s)", cmdName, cmdStr)
	if serviceLog != nil {
		serviceLog.Info(msg)
	} else {
		log.Println(msg)
	}

	// Execute the command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		// Fallback for Linux/Mac
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Error executing command: %s\nOutput: %s", err, output)
		if serviceLog != nil {
			serviceLog.Error(errMsg)
		} else {
			log.Println(errMsg)
		}
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Command '%s' executed successfully.\nOutput: %s", cmdName, output)))
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
		action := os.Args[1]
		err = service.Control(s, action)
		if err != nil {
			fmt.Printf("Action '%s' failed: %v\n", action, err)
			// Don't exit with error code strictly, just inform user.
			// Unless it's critical.
		} else {
			fmt.Printf("Action '%s' succeeded.\n", action)
		}
		return
	}

	// Run the service
	// If run interactively (not as service), this will run the program normally
	err = s.Run()
	if err != nil {
		if serviceLog != nil {
			serviceLog.Error(err.Error())
		}
	}
}
