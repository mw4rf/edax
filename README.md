# edax

`edax` is a CLI time manager written in Golang.

## Use case

Let's create a new timer for the "Design Phase" of my new shiny project. `$ edax create "Design Phase"` Yay!

Now, let's start the "Design Phase" timer, it's time to work! `$ edax start`

I need to list all my timers to check their IDs. `$ edax list` Here's what I have: 1. Design Phase 2. Development 3. Testing

It seems I have other timers as well. I'll start the "Development" timer and stop all others. `$ edax start 2`

Working hard for a while...

I need to check the status of the currently running timer. `$ edax status` Wow... 2 hours already, time for a nap! Let's stop the "Development" timer. `$ edax stop`

_A few moments later_ (I mean, after the nap). Let's see all timers I've started today. `$ edax today` Here's today's list: 1. Design Phase - started at 9:00 AM 2. Development - started at 10:15 AM, stopped 2 PM.

I want to delete the "Testing" timer since I don't need it anymore. `$ edax delete 3` Clean slate!

Finally, I'll reset the "Design Phase" timer to start fresh tomorrow. `$ edax reset 1` Ready for a new day!

## Usage

``` sh
Usage: edax [command] <value>
Commands:
  version               => Show version number
  create <name>         => Create a new timer with a name
  start                 => Start the last timer created
  start <id>            => Start a timer, and stop all others
  stop                  => Stop the running timer
  reset <id>            => Reset a timer
  delete <id>           => Delete a timer
  list                  => List all timers
  today                 => List all timers started today
  search <query>        => Search for timers
  status                => Show the status of the running timer
  status <id>           => Show the status of a timer
```
