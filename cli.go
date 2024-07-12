package main

import (
	"fmt"
	"strings"
	"time"
)


const (
	grey      = "\033[37m"
	lightBlue = "\033[94m"
	red       = "\033[91m"
	green     = "\033[92m"
	yellow    = "\033[93m"
	reset     = "\033[0m"
)

// Print the list of timers
func printList(timers *[]Timer) {
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

func printStatus(timers *[]Timer, id *int) {
	// If no id is provided, print the status of the running timer
	if id == nil {
		for _, timer := range *timers {
			if timer.Running {
				id = &timer.Id
				break
			}
		}
	}

	// If no running timer is found
	if id == nil {
		fmt.Println("No running timer")
		return
	}

	// Print the status of the timer
	timer := (*timers)[*id]
	status, color := getStatus(timer)
	duration := calculateDuration(timer.Start, timer.End, timer.Running)
	fmt.Printf("%s%s%s %s%s%s %s%s%s\n", color, status, reset, yellow, timer.Name, reset, grey, duration, reset)
}
