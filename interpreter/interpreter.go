package interpreter

import (
	"CLIManager/commandManager"
	"fmt"
)

func InterpretCommands(
	parsedCommands commandManager.ParsedCommands,
	requireInit bool,
) {
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
	fmt.Println("Process Command With Value: ", command)
}

// If only the command exists but requires input
func interpretCommandWithoutValue(
	command commandManager.ParsedCommand,
) {
	if !commandRequiresInput(command) {
		// Run command
		fmt.Println("COMMAND DOES NOT REQUIRE INPUT: ", command)
		return
	}

	fmt.Println("COMMAND REQUIRES INPUT: ", command)
}

func commandRequiresInput(
	command commandManager.ParsedCommand,
) bool {
	return true
}
