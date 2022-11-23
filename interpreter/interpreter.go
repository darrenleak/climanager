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
	parsedActions = allCommands
	currentConfig = config

	for _, command := range parsedCommands {
		commandExists := len(command.Command) > 0
		valueExists := len(command.Value) > 0

		if commandExists && !valueExists {
			interpretCommandWithoutValue(command)
		} else if commandExists && valueExists {
			interpretCommandWithValue(command)
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
			commandFile := commandManager.ReadUserInput("Command file to append")
			commandFilesAppend(commandFile)
		case commandManager.CommandFilesRemove:
			commandFile := commandManager.ReadUserInput("Command file to remove")
			commandFilesRemove(commandFile)
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
