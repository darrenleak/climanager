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

func LoadConfig() (Config, error) {
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

func InitConfig() {
	config, _ := LoadConfig()

	questions := []ConfigInitQuestion{
		{
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
	WriteConfig(updateConfig)
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

// TODO: Write command files alphabetically
// TODO: Try to improve how the config file gets written to
// Multiple places are writing to the config with different ways of
// doing it
func AppendCommandFilePath(config Config, commandFilePath string) {
	_, err := ioutil.ReadFile(commandFilePath)
	if err != nil {
		fmt.Println("Error reading command file: ", commandFilePath)
		return
	}

	config.CommandFiles = append(config.CommandFiles, commandFilePath)
	WriteConfig(config)
}

func RemoveCommandFilePath(config Config, commandFilePath string) {
	commandFiles := config.CommandFiles
	newCommandFiles := []string{}

	for index := range commandFiles {
		if commandFiles[index] != commandFilePath {
			newCommandFiles = append(newCommandFiles, commandFiles[index])
		}
	}

	config.CommandFiles = newCommandFiles

	WriteConfig(config)
}

func Updateshell(config Config, shell string) {
	if len(shell) == 0 || shell == "\n" {
		shell = "bash"
	}

	config.Shell = strings.TrimSuffix(shell, "\n")
	WriteConfig(config)
}

func WriteConfig(config Config) {
	file, _ := json.MarshalIndent(config, "", "  ")
	_ = ioutil.WriteFile("./config.json", file, 0644)
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
	isNextArgumentAValue, nextIndex := DoesNextArgumentExistAndIsNotCommand(args, argIndex)

	if isNextArgumentAValue {
		nextValue := args[nextIndex]
		return nextValue
	} else {
		nextValue := ReadUserInput(userPrompt)
		return nextValue
	}
}
