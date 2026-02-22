package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the config.yaml file
type Config struct {
	Commands map[string]string `yaml:"commands"`
	Port     string            `yaml:"port"` // Optional: Allow custom port
}

var config Config

func loadConfig() error {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &config)
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

	log.Printf("Executing command: %s (%s)", cmdName, cmdStr)

	// Execute the command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", cmdStr)
	} else {
		// Fallback for Linux/Mac (mostly for testing purposes)
		cmd = exec.Command("sh", "-c", cmdStr)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("Error executing command: %s\nOutput: %s", err, output)
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Command '%s' executed successfully.\nOutput: %s", cmdName, output)))
}

func main() {
	// Load configuration
	if err := loadConfig(); err != nil {
		log.Printf("Error loading config.yaml: %v", err)
		log.Println("Using default config...")
		config = Config{
			Commands: map[string]string{
				"shutdown": "shutdown /s /t 0",
			},
		}
	}

	// Set default port if not specified
	port := config.Port
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/execute/", executeHandler)
	
	// Add a root handler to check if server is running
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Shutdown Tool is running! Use /execute/{command} to trigger actions."))
	})

	fmt.Printf("Server starting on port %s...\n", port)
	fmt.Printf("Available commands: %v\n", config.Commands)
	
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
