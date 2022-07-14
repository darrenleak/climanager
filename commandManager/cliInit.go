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

func cliInit() {
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
