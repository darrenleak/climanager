package interpreter

import (
	"CLIManager/commandManager"
	"CLIManager/orchestrator"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"golang.org/x/exp/slices"
)

var parsedActions []map[string]map[string]orchestrator.Runnable
var currentConfig commandManager.Config

func InterpretCommands(
	parsedCommands commandManager.ParsedCommands,
	requireInit bool,
	allCommands []map[string]map[string]orchestrator.Runnable,
	config commandManager.Config,
) {
	if parsedCommands[0].Command == commandManager.Help {
		printHelp()
		return
	}

	parsedActions = allCommands
	currentConfig = config

	for _, command := range parsedCommands {
		commandExists := len(command.Command) > 0
		valueExists := len(command.Value) > 0

		if commandExists && !valueExists {
			interpretCommandWithoutValue(command)
		} else if commandExists && valueExists {
			interpretCommandWithValue(command)
		} else if !commandExists {
			// TODO: This is not great, as this is acting like a value, but it is a command
			orchestrator.Run(currentConfig, command.Value, allCommands[0])
		}
	}
}

func interpretCommandWithValue(
	command commandManager.ParsedCommand,
) {
	if commandRequiresInput(command) {
		switch command.Command {
		case commandManager.CommandFilesAppend:
			commandFilesAppend(command.Value)
		case commandManager.CommandFilesRemove:
			commandFilesRemove(command.Value)
		case commandManager.Shell:
			updateShell(command.Value)
		}

		return
	}
}

// If only the command exists but requires input
func interpretCommandWithoutValue(
	command commandManager.ParsedCommand,
) {
	if commandRequiresInput(command) {
		switch command.Command {
		case commandManager.CommandFilesAppend:
			commandFile := commandManager.ReadUserInput("Command file path(absolute path) to append")
			commandFilesAppend(commandFile)
		case commandManager.CommandFilesRemove:
			commandFile := commandManager.ReadUserInput("Command file path(absolute path) to remove")
			commandFilesRemove(commandFile)
		case commandManager.Shell:
			shell := commandManager.ReadUserInput("What shell are you using?")
			updateShell(shell)
		case commandManager.Profile:
			profilePath := commandManager.ReadUserInput("Absolute path to profile file(zsh, bash, etc)")
			updateProfile(profilePath)
		}

		return
	}

	switch command.Command {
	case commandManager.ViewConfig:
		viewConfig()
	case commandManager.ListCommands:
		listCommands()
	}
}

func updateShell(shell string) {
	commandManager.Updateshell(currentConfig, shell)
}

func updateProfile(profilePath string) {
	commandManager.UpdateProfile(currentConfig, profilePath)
}

func commandFilesAppend(commandFilePath string) {
	commandManager.AppendCommandFilePath(currentConfig, commandFilePath)
}

func commandFilesRemove(commandFilePath string) {
	commandManager.RemoveCommandFilePath(currentConfig, commandFilePath)
}

func viewConfig() {
	marshalledJson, err := json.Marshal(currentConfig)
	if err != nil {
		fmt.Println("Error unmarshalling config")
		return
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, marshalledJson, "", "\t")
	if err != nil {
		log.Println("JSON parse error")
		return
	}

	fmt.Println(prettyJSON.String())
}

func listCommands() {
	var outputActions []string
	actions := parsedActions[0]

	for action := range actions {
		outputActions = append(outputActions, action)
	}

	sort.Strings(outputActions)
	for _, action := range outputActions {
		fmt.Println(action)
	}
}

func commandRequiresInput(
	command commandManager.ParsedCommand,
) bool {
	commandsRequiringInput := []string{
		commandManager.CommandFiles,
		commandManager.ListCommands,
		commandManager.ViewConfig,
	}

	return !slices.Contains(commandsRequiringInput, command.Command)
}

func printHelp() {
	fmt.Println(commandManager.Init, "               Setup the config file by asking a few questions")
	fmt.Println(commandManager.Shell, "              Allow you to update the shell setting in the config")
	fmt.Println(commandManager.Profile, "            Allow you to update the profile setting in the config")
	fmt.Println(commandManager.CommandFiles, "       Allow you to update the command files in the config")
	fmt.Println(commandManager.CommandFilesAppend, " Allow you to append to the command files in the config")
	fmt.Println(commandManager.CommandFilesRemove, " Allow you to remove from the command files in the config")
	fmt.Println(commandManager.ListCommands, "       List all the actions")
	fmt.Println(commandManager.ViewConfig, "         Print out the current config file")
	fmt.Println(commandManager.Help, "               Shows help, what you are seeing now :)")
}
