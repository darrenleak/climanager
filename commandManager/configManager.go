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

func loadConfig() (Config, error) {
	jsonFile, fileOpenError := os.Open("./config.json")

	if fileOpenError != nil {
		return Config{}, fileOpenError
	}

	configData, configReadError := ioutil.ReadAll(jsonFile)

	if configReadError != nil {
		return Config{}, configReadError
	}

	var config Config
	json.Unmarshal(configData, &config)

	return config, nil
}

func initConfig() {
	config, configLoadError := loadConfig()

	if configLoadError != nil {
		// TODO: Handle this better
		fmt.Println("Config load error. Either no file exists or it failed to read")
	}

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

	updateConfig := buildConfig(inputValues, config)
	writeConfig(updateConfig)
}

/*
TODO:

This needs to be fixed. It needs to understand the the commands that the user
has entered and it needs to update the corresponding value.

Currently this is updating in the order the commands come in, which causes
issues when you do something like:

--config --commandFiles

The above command will update the Shell value in the config, not the CommandFiles
in the config
*/
func buildConfig(initConfigSettings []string, providedConfig Config) Config {
	config := providedConfig

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
	config, configLoadError := loadConfig()

	if configLoadError != nil {
		// TODO: Handle this better
		fmt.Println("Config load error. Either no file exists or it failed to read")
	}

	inputValues := []string{}

	for argIndex, arg := range args {
		if argIndex <= 1 {
			continue
		}

		textHint := getTextHint(arg)

		inputValue := getInputValue(args, argIndex, textHint)
		inputValues = append(inputValues, inputValue)
	}

	return buildConfig(inputValues, config)
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
