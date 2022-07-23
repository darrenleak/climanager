package commandManager

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
TODO: Implement
- Add commandFileAppend
- Add commandFileRemove
- Add viewConfig

Without these two, when you update your config, it will overwrite all your
commandFiles with the new value.
*/
const (
	Init               string = "--init"
	Shell              string = "--shell"
	CommandFiles       string = "--commandFiles"
	CommandFilesAppend string = "--commandFilesAppend"
	CommandFilesRemove string = "--commandFilesRemove"
	SetConfig          string = "--config"
	ViewConfig         string = "--viewConfig"
)

type ParsedCommand struct {
	Command string
	Value   string
}

type ParsedCommands struct {
	actions []ParsedCommand
}

/*
What about a middleware architecture for this?

Basically like a lifecycle. This will allow custmizations
in the future if they are needed. Maybe something will need to be
injected or processed during the parsing?
*/

/*
This needs to parse commands, not run them. It should return the commands
that the user entered in a way that the program can understand and something
else needs to run the commands
*/
func ParseArgs(args []string, requireInit bool) bool {
	argsCount := len(args)
	isInitCommand := argsCount > 1 && args[1] == Init

	if requireInit && !isInitCommand {
		fmt.Println("Please run: CLIManager --init, to setup CLIManager")
		return true
	}

	// If init is the first argument, ignore the rest
	if isInitCommand {
		initConfig()

		return true
	}

	if argsCount > 1 && args[1] == SetConfig {
		config := processConfigInput(args)

		writeConfig(config)

		return true
	}

	return false
}

func doesNextArgumentExistAndIsNotCommand(args []string, index int) (bool, int) {
	argsCount := len(args) - 1

	if argsCount <= index {
		return false, index
	}

	nextIndex := index + 1

	if strings.HasPrefix(args[nextIndex], "--") {
		return false, index
	}

	return true, nextIndex
}

func RequireCliSetup() (bool, Config) {
	config, configLoadError := LoadConfig()

	if configLoadError != nil {
		return true, Config{}
	}

	return false, config
}

func readUserInput(userPrompt string) string {
	// Displays the prompt to the user
	fmt.Println(userPrompt)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return line
}
