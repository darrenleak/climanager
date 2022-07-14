package commandManager

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	Init         string = "--init"
	Shell        string = "--shell"
	CommandFiles string = "--commandFiles"
)

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

	for argIndex, arg := range args {
		if argIndex == 0 {
			continue
		}

		// If init is the first argument, ignore the rest
		if arg == Init && argIndex == 1 {
			cliInit()
			break
		}

		if arg == Shell {
			if doesNextArgumentExistAndIsNotCommand("") {
				// Get value
			} else {
				processInput("Shell to use: ")
			}
		}

		if arg == CommandFiles {
			if doesNextArgumentExistAndIsNotCommand("") {
				// Get value
			} else {
				processInput("Paths to command files: ")
			}
		}
	}

	// // TODO: This is bad
	// // There should be an argument processor
	// // to handle all of this stuff
	// if args[1] == "init" {

	// 	// } else if args[1] == "shell" {
	// 	// 	shell = args[2]
	// 	// 	actionToRun = args[3]
	// 	// } else {
	// 	// 	shell = "bash"
	// }
}

func doesNextArgumentExistAndIsNotCommand(arg string) bool {
	return false
}

func ensureCliSetup() bool {
	return true
}

func processInput(title string) string {
	fmt.Println(title)

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return line
}
