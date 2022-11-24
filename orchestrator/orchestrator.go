package orchestrator

import (
	"CLIManager/commandManager"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Actions struct {
	Actions []Action //`yaml:"actions"`
}

type Action struct {
	Name      string     //`yaml:"name"`
	Runnables []Runnable //`yaml:"runnables"`
}

type Runnable struct {
	Name      string   //`yaml:"name"`
	DependsOn []string `yaml:"dependsOn"`
	Command   string   //`yaml:"command"`
	Alias     string   //`yaml:"alias"`
}

type RunningStatus struct {
	Complete []Runnable
	Failed   []Runnable
	Running  []Runnable
	Waiting  []Runnable
}

type RunningStatuses struct {
	Statuses []RunningStatus
}

type RunnableStatus struct {
	Output    []string
	DependsOn []string
}

func BuildCommandTree(
	config commandManager.Config,
) []map[string]map[string]Runnable {
	var commands []map[string]map[string]Runnable

	actionFilePaths := config.CommandFiles
	for _, path := range actionFilePaths {
		commandTree, err := setupCommands(path)
		if err != nil {
			continue
		}

		commands = append(commands, commandTree)
	}

	return commands
}

var dependentsMap = make(map[string][]string)

func setupCommands(filePath string) (map[string]map[string]Runnable, error) {
	var runnablesStatuses = make(map[string]map[string]Runnable)

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	actions := Actions{}
	err = yaml.Unmarshal(yamlFile, &actions)

	if err != nil {
		panic("Unmarshal failed")
	}

	for _, currentAction := range actions.Actions {
		runnables := currentAction.Runnables
		runnablesStatuses[currentAction.Name] = make(map[string]Runnable)

		for _, runnable := range runnables {
			runnablesStatuses[currentAction.Name][runnable.Name] = runnable
		}

		for _, runnable := range runnables {
			for _, dependent := range runnable.DependsOn {
				dependentsMap[dependent] = append(dependentsMap[dependent], runnable.Name)
			}
		}
	}

	return runnablesStatuses, nil
}
