# timex

A simple, terminal-based time tracking utility written in Go.

## Description

`timex` helps you track the amount of time you spend on different projects directly from your command line. No GUIs, no distractions.

## Usage

Here are the commands currently planned for `timex`:

*   `timex start <project>`: Start a timer for a project.
*   `timex stop`: Stop the currently running timer.
*   `timex switch <project>`: Stop the current timer and start a new one for a different project.
*   `timex status`: Show the currently running timer and elapsed time.
*   `timex log <project> --from "HH:MM" --to "HH:MM"`: Manually log a time entry.
*   `timex report [today|week|all]`: Generate a report of time spent.
*   `timex list`: List all projects.
*   `timex add <project>`: Add a new project.
*   `timex remove <project>`: Remove a project.
*   `timex export --format [csv|json]`: Export time tracking data.

## Building

To build the project, you can run:

```sh
go build
```

To install it to your `$GOPATH/bin`:

```sh
go install
```
