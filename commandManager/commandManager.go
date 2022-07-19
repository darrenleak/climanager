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
func ParseArgs(args []string) {
	if !ensureCliSetup() {
		fmt.Println("Please run: CLIManager --init, to setup CLIManager")
		return
	}

	argsCount := len(args)

	// If init is the first argument, ignore the rest
	if argsCount > 1 && args[1] == Init {
		initConfig()
	}

	if argsCount > 1 && args[1] == SetConfig {
		config := processConfigInput(args)

		writeConfig(config)
	}
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

func ensureCliSetup() bool {
	_, configLoadError := loadConfig()

	if configLoadError != nil {
		fmt.Println("Error ensuring cli is setup")
		return false
	}

	return true
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
