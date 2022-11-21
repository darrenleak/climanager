package main

/*
TODO:
1. Set yaml files
3. Register commands
	Sometimes exec cannot find a command, eg nvm. The user should be able to register commands and the paths to the executable
		- Auto check on failure if command is registered. Acts more as a try then.
*/

import (
	"CLIManager/commandManager"
	"CLIManager/orchestrator"
	"os"
)

type Runnable struct {
	Name      string   //`yaml:"name"`
	DependsOn []string `yaml:"dependsOn"`
	Command   string   //`yaml:"command"`
	Alias     string   //`yaml:"alias"`
}

type Action struct {
	Name      string     //`yaml:"name"`
	Runnables []Runnable //`yaml:"runnables"`
}

type Actions struct {
	Actions []Action //`yaml:"actions"`
}

type RunningStatus struct {
	Complete []Runnable
	Failed   []Runnable
	Running  []Runnable
	Waiting  []Runnable
}

type RunningStatuses struct {
	Statuses []RunningStatus
}

type RunnableStatus struct {
	Output    []string
	DependsOn []string
}

var config commandManager.Config

func main() {
	args := os.Args
	requireInit, loadedConfig := commandManager.RequireCliSetup()
	hasActioned := commandManager.ParseArgs(args, requireInit)

	// TODO: Test code
	// parsedCommands := commandManager.Parse(args)
	// interpreter.InterpretCommands(parsedCommands, requireInit)

	if hasActioned {
		return
	}

	config = loadedConfig
	// actionToRun := os.Args[1]
	// allRunnables = setupCommands(config.CommandFiles)
	orchestrator.BuildCommandTree(config)

	// runAction(actionToRun)
}
