package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"path/filepath"
)


// Parse a string to an integer
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error parsing integer")
		os.Exit(1)
	}
	return i
}

// Format time from int64 to a readable string
func formatTime(timestamp int64) string {
	if timestamp == 0 {
		return "N/A"
	}
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// Calculate the duration between start and end times
func calculateDuration(start, end int64, running bool) string {
	if start == 0 {
		return "Not Started"
	}
	if running {
		end = time.Now().Unix()
	}
	duration := end - start
	return fmt.Sprintf("%d seconds", duration)
}


// Returns the default directory for the application,
// according to the operating system
// e.g. /home/user/.config/edax
func getDefaultDirectory() (string, error) {
	// Get the user config directory
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	// Append your application directory
	dir = filepath.Join(dir, "edax")
	return dir, nil
}
