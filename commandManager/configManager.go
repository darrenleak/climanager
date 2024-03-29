package commandManager

import (
	"CLIManager/cliHttp"
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
	Shell             string
	Profile           string
	CommandFiles      []string
	CommandFilesUrls  map[string]string
	DefaultActionFile string
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
			"Specify shell to use(default is [bash]):",
		},
		{
			Profile,
			"zsh/bash profile path(Absolute path):",
		},
		{
			CommandFiles,
			"Action files path/url(Absolute path)",
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

		line = strings.TrimSuffix(line, "\n")

		if len(line) == 0 && question.Arg != Shell {
			fmt.Println("Value is required for: ", question.Arg)
			return
		}

		if len(line) == 0 && question.Arg == Shell {
			line = "bash"
		}

		inputValues = append(inputValues, ConfigPartial{
			ConfigSetting:      question.Arg,
			ConfigSettingValue: line,
		})
	}

	updateConfig := buildConfig(inputValues, config)

	err := WriteConfig(updateConfig)
	if err != nil {
		fmt.Println("Failed to initialise app")
		return
	}

	fmt.Println("Successfully initialised app")
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

		if setting.ConfigSetting == Profile {
			config.Profile = removedNewLineConfigSettingValue
			continue
		}

		if setting.ConfigSetting == CommandFiles {
			config.CommandFiles = strings.Split(removedNewLineConfigSettingValue, ",")
			downloadFileMap := downloadCommandFiles(config.CommandFiles)
			config.CommandFilesUrls = downloadFileMap

			continue
		}
	}

	// Replace config.CommandFiles with file paths of downloaded files
	for commandFileURL := range config.CommandFilesUrls {
		for index, commandFile := range config.CommandFiles {
			if commandFile == commandFileURL {
				config.CommandFiles[index] = config.CommandFilesUrls[commandFileURL]
			}
		}
	}

	return config
}

func downloadCommandFiles(commandFiles []string) map[string]string {
	downloadedFiles := map[string]string{}

	for _, commandFile := range commandFiles {
		if !cliHttp.IsCommandFileURL(commandFile) {
			continue
		}

		fileName, err := cliHttp.DownloadFile(commandFile)
		if err != nil {
			fmt.Println("Could not download command file: ", commandFile)
			return nil
		}

		downloadedFiles[commandFile] = *fileName
	}

	return downloadedFiles
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
	err = WriteConfig(config)

	if err != nil {
		fmt.Println("Failed to append command file: ", commandFilePath)
		return
	}

	err = UseActionFile(config, commandFilePath)
	if err != nil {
		return
	}

	fmt.Println("Appended command file successfully: ", commandFilePath)
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

	err := WriteConfig(config)

	if err != nil {
		fmt.Println("Failed to remove command file: ", commandFilePath)
		return
	}

	fmt.Println("Removed command file successfully: ", commandFilePath)
}

func Updateshell(config Config, shell string) {
	if len(shell) == 0 || shell == "\n" {
		shell = "bash"
	}

	config.Shell = strings.TrimSuffix(shell, "\n")
	err := WriteConfig(config)

	if err != nil {
		fmt.Println("Failed to update shell: ", shell)
		return
	}

	fmt.Println("Updated shell successfully: ", shell)
}

func UpdateProfile(config Config, profilePath string) {
	if len(profilePath) == 0 || profilePath == "\n" {
		return
	}

	config.Profile = strings.TrimSuffix(profilePath, "\n")
	err := WriteConfig(config)

	if err != nil {
		fmt.Println("Failed to update profile: ", profilePath)
		return
	}

	fmt.Println("Updated profile successfully: ", profilePath)
}

func UseActionFile(config Config, actionFilePath string) error {
	if len(actionFilePath) == 0 || actionFilePath == "\n" {
		return nil
	}

	config.DefaultActionFile = strings.TrimSuffix(actionFilePath, "\n")
	err := WriteConfig(config)

	if err != nil {
		fmt.Println("Failed to update default action file: ", actionFilePath)
		return err
	}

	fmt.Println("Updated default action file to use", actionFilePath)
	return nil
}

func WriteConfig(config Config) error {
	file, _ := json.MarshalIndent(config, "", "  ")
	err := ioutil.WriteFile("./config.json", file, 0644)

	if err != nil {
		fmt.Println("Failed to write config: ", err)
		return err
	}

	return nil
}
