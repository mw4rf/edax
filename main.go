package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"regexp"
	"strings"
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

// Parse a string to an integer
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Error parsing integer")
		os.Exit(1)
	}
	return i
}

// Timer struct
type Timer struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Start   int64  `json:"start"`
	End     int64  `json:"end"`
	Running bool   `json:"running"`
	Tags    []Tag  `json:"tags"`
}

type Tag struct {
	Name string `json:"name"`
}

// Create a new timer
func createTimer(id int, name string) Timer {
	// Extract tags from the name
	tagRegex := regexp.MustCompile(`#\w+`)
	tags := tagRegex.FindAllString(name, -1)

	// Remove tags from the name
	for _, tag := range tags {
		name = strings.Replace(name, tag, "", -1)
	}

	// Trim any extra spaces from the name
	name = strings.TrimSpace(name)

	// Create Tag objects
	tagObjects := make([]Tag, len(tags))
	for i, tag := range tags {
		tagObjects[i] = Tag{Name: tag}
	}

	return Timer{
		Id:      id,
		Name:    name,
		Start:   0,
		End:     0,
		Running: false,
		Tags:    tagObjects,
	}
}

// Print the list of timers
func printList(timers *[]Timer) {
	const (
		grey      = "\033[90m"
		lightBlue = "\033[94m"
		red       = "\033[91m"
		green     = "\033[92m"
		yellow    = "\033[93m"
		reset     = "\033[0m"
	)

	fmt.Println(grey + "Timers List:" + reset)
	fmt.Println(grey + "------------------------------------------------------" + reset)
	for _, timer := range *timers {
		status, color := getStatus(timer)
		duration := calculateDuration(timer.Start, timer.End, timer.Running)
		// Name
		fmt.Printf("[%d] %s%s%s %s%s%s", timer.Id, color, status, reset, yellow, timer.Name, reset)
		// Tags
		if len(timer.Tags) > 0 {
			fmt.Printf(" ")
			for _, tag := range timer.Tags {
				fmt.Printf("%s%s%s ", lightBlue, tag.Name, reset)
			}
		}
		fmt.Printf("\n")
		// Start, End, Duration
		fmt.Printf("%sStart: %-19s End: %-19s Duration: %s%s\n", grey, formatTime(timer.Start), formatTime(timer.End), duration, reset)
		fmt.Println(grey + "------------------------------------------------------" + reset)
	}
}

func printSearch(timers *[]Timer, query string) {
	fmt.Printf("Searching for '%s'\n", query)
	results := make([]Timer, 0)
	for _, timer := range *timers {
		if strings.Contains(timer.Name, query) {
			results = append(results, timer)
			continue
		}
		for _, tag := range timer.Tags {
			if strings.Contains(tag.Name, query) {
				results = append(results, timer)
				break
			}
		}
	}
	if len(results) == 0 {
		fmt.Println("No results found")
		return
	}
	printList(&results)
}

func printToday(timers *[]Timer) {
	fmt.Println("Today's Timers:")
	results := make([]Timer, 0)
	for _, timer := range *timers {
		if time.Now().Day() == time.Unix(timer.Start, 0).Day() {
			results = append(results, timer)
		}
	}
	if len(results) == 0 {
		fmt.Println("No timers today")
		return
	}
	printList(&results)
}

// Get the status of the timer as an emoji with appropriate color
func getStatus(timer Timer) (string, string) {
	const (
		red    = "\033[91m"
		green  = "\033[92m"
	)

	if timer.Running {
		return "[Running]", green
	} else {
		return "[Stopped]", red
	}
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

// Start a timer identified by its id
func startTimer(timers *[]Timer, id int) {
	if id < 0 || id >= len(*timers) {
		fmt.Println("Invalid timer ID")
		return
	}
	timer := &(*timers)[id]
	if timer.Running {
		fmt.Println("Timer is already running")
		return
	}
	// Start the timer
	// Are we resuming a stopped timer?
	if timer.Start != 0 && timer.End != 0 {
		timer.Start = time.Now().Unix() - (timer.End - timer.Start)
	} else {
		timer.Start = time.Now().Unix()
	}
	timer.End = 0
	timer.Running = true
	// Stop all other timers
	for i := range *timers {
		if i != id {
			(*timers)[i].Running = false
			(*timers)[i].End = time.Now().Unix()
		}
	}
	fmt.Printf("Timer %s started at %s\n", timer.Name, formatTime(timer.Start))
}

// Reset a timer identified by its id
func resetTimer(timers *[]Timer, id int) {
	if id < 0 || id >= len(*timers) {
		fmt.Println("Invalid timer ID")
		return
	}
	timer := &(*timers)[id]
	timer.Start = 0
	timer.End = 0
	timer.Running = false
	fmt.Printf("Timer %s reset\n", timer.Name)
}

// Stop the current timer
func stopTimer(timers *[]Timer) {
	for i := range *timers {
		if (*timers)[i].Running {
			(*timers)[i].End = time.Now().Unix()
			(*timers)[i].Running = false
			fmt.Printf("Timer %s stopped at %s\n", (*timers)[i].Name, formatTime((*timers)[i].End))
			return
		}
	}
}

// Save all timers to a file as a JSON object
func saveTimers(timers *[]Timer) {
	file, err := os.Create("timers.json")
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
	file, err := os.Open("timers.json")
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

// Delete a timer
func deleteTimer(timers *[]Timer, id int) {
	if id < 0 || id >= len(*timers) {
		fmt.Println("Invalid timer ID")
		return
	}
	*timers = append((*timers)[:id], (*timers)[id+1:]...)
}
