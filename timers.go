package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

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

	// Stop all other running timers
	for i := range *timers {
		if i != id && (*timers)[i].Running {
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

// Delete a timer
func deleteTimer(timers *[]Timer, id int) {
	if id < 0 || id >= len(*timers) {
		fmt.Println("Invalid timer ID")
		return
	}
	*timers = append((*timers)[:id], (*timers)[id+1:]...)
}
