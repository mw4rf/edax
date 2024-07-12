package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Main
func main() {
	// Get command line arguments
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage: edax [command] <value>")
		fmt.Println("Commands:")
		fmt.Println("  version		=> Show version number")
		fmt.Println("  create <name>		=> Create a new timer with a name")
		fmt.Println("  start			=> Start the last timer created")
		fmt.Println("  start <id>		=> Start a timer, and stop all others")
		fmt.Println("  stop			=> Stop the running timer")
		fmt.Println("  reset <id>		=> Reset a timer")
		fmt.Println("  delete <id>		=> Delete a timer")
		fmt.Println("  list			=> List all timers")
		fmt.Println("  today			=> List all timers started today")
		fmt.Println("  search <query>	=> Search for timers")
		os.Exit(1)
	}

	// Load the timers from JSON file
	timers := loadTimers()
	defer saveTimers(&timers)

	// Switch the command
	switch args[0] {
	case "version":
		fmt.Println("Edax v0.2")
	case "create":
		t := createTimer(len(timers), args[1])
		timers = append(timers, t)
	case "start":
		// No id provided: start the last timer
		id := len(timers) - 1
		// Id provided: start the timer with that id
		if len(args) >= 1 {
			id = parseInt(args[1])
		}
		startTimer(&timers, id)
	case "reset":
		if len(args) < 2 {
			fmt.Println("Usage: edax reset <id>")
			os.Exit(1)
		}
		id := args[1]
		resetTimer(&timers, parseInt(id))
	case "stop":
		stopTimer(&timers)
	case "delete":
		if len(args) < 2 {
			fmt.Println("Usage: edax delete <id>")
			os.Exit(1)
		}
		id := args[1]
		deleteTimer(&timers, parseInt(id))
	case "list":
		printList(&timers)
	case "today":
		printToday(&timers)
	case "search":
		if len(args) < 2 {
			fmt.Println("Usage: edax search <query>")
			os.Exit(1)
		}
		printSearch(&timers, args[1])
	}

}

// Save all timers to a file as a JSON object
func saveTimers(timers *[]Timer) {
	// Get the default directory, create it if it doesn't exist
	dir, err := getDefaultDirectory()
	if err != nil {
		fmt.Println("Error getting default directory")
		return
	}
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory")
		return
	}

	// Create the full file path
	filePath := filepath.Join(dir, "timers.json")

	// Create the file
	// If the file already exists, it will be overwritten
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(timers)
	if err != nil {
		fmt.Println("Error writing JSON object")
		return
	}
}

// Load all timers from a file as a JSON object
func loadTimers() []Timer {
	// Get the default directory,
	// return an empty Timer array if it doesn't exist
	dir, err := getDefaultDirectory()
	if err != nil {
		return make([]Timer, 0)
	}

	// Create the full file path
	filePath := filepath.Join(dir, "timers.json")

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return make([]Timer, 0)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	timers := make([]Timer, 0)
	err = decoder.Decode(&timers)
	if err != nil {
		return make([]Timer, 0)
	}
	return timers
}

