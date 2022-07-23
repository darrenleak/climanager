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

type ConfigPartial struct {
	ConfigSetting      string
	ConfigSettingValue string
}

type Config struct {
	Shell        string
	CommandFiles []string
}

type ConfigInitQuestion struct {
	Arg        string
	UserPrompt string
}

type ConfigInitQuestions struct {
	InitQuestions []ConfigInitQuestion
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

	questions := []ConfigInitQuestion{{
		Shell,
		"Shell(default is Bash):",
	},
		{
			CommandFiles,
			"Action files(file path separated by comma)",
		},
	}
	initValues := ConfigInitQuestions{
		InitQuestions: questions,
	}
	inputValues := []ConfigPartial{}

	for _, question := range initValues.InitQuestions {
		fmt.Println(question.UserPrompt)

		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		inputValues = append(inputValues, ConfigPartial{
			ConfigSetting:      question.Arg,
			ConfigSettingValue: line,
		})
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
func buildConfig(initConfigSettings []ConfigPartial, providedConfig Config) Config {
	config := providedConfig

	fmt.Println(config)
	fmt.Println(initConfigSettings)

	for _, setting := range initConfigSettings {
		removedNewLineConfigSettingValue := strings.TrimSuffix(setting.ConfigSettingValue, "\n")

		if setting.ConfigSetting == Shell {
			config.Shell = removedNewLineConfigSettingValue
			continue
		}

		if setting.ConfigSetting == CommandFiles {
			config.CommandFiles = strings.Split(removedNewLineConfigSettingValue, ",")
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

	inputValues := []ConfigPartial{}

	for argIndex, arg := range args {
		if argIndex <= 1 {
			continue
		}

		userPrompt := getTextHint(arg)
		inputValue := getInputValue(args, argIndex, userPrompt)

		configPartial := ConfigPartial{
			ConfigSetting:      arg,
			ConfigSettingValue: inputValue,
		}
		inputValues = append(inputValues, configPartial)
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
	userPrompt string,
) string {
	isNextArgumentAValue, nextIndex := doesNextArgumentExistAndIsNotCommand(args, argIndex)

	if isNextArgumentAValue {
		nextValue := args[nextIndex]
		return nextValue
	} else {
		nextValue := readUserInput(userPrompt)
		return nextValue
	}
}
