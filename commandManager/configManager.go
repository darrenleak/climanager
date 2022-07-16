package commandManager

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	Shell        string
	CommandFiles []string
}

func initConfig() {
	initValues := []string{
		"Shell(default is Bash):",
		"Action files(file path separated by comma)",
	}
	inputValues := []string{}

	for _, question := range initValues {
		fmt.Println(question)

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		inputValues = append(inputValues, line)
	}

	config := buildConfig(inputValues)
	writeConfig(config)
}

func buildConfig(initConfigSettings []string) Config {
	config := Config{}

	for settingIndex, setting := range initConfigSettings {
		if settingIndex == 0 {
			config.Shell = setting
			continue
		}

		if settingIndex == 1 {
			config.CommandFiles = strings.Split(setting, ",")
			continue
		}
	}

	return config
}

func writeConfig(config Config) {
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile("./config.json", file, 0644)
}

func processConfigInput(args []string) Config {
	inputValues := []string{}

	for argIndex, arg := range args {
		if argIndex <= 1 {
			continue
		}

		textHint := getTextHint(arg)

		inputValue := getInputValue(args, argIndex, textHint)
		inputValues = append(inputValues, inputValue)
	}

	return buildConfig(inputValues)
}

func getTextHint(arg string) string {
	if arg == Shell {
		return "Shell to use: "
	}

	if arg == CommandFiles {
		return "Paths to command files: "
	}

	return ""
}

func getInputValue(
	args []string,
	argIndex int,
	textHint string,
) string {
	isNextArgumentAValue, nextIndex := doesNextArgumentExistAndIsNotCommand(args, argIndex)

	if isNextArgumentAValue {
		nextValue := args[nextIndex]
		return nextValue
	} else {
		nextValue := processInput(textHint)
		return nextValue
	}
}
