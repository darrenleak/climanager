package main

/*
TODO:
1. Set yaml files
3. Register commands
	Sometimes exec cannot find a command, eg nvm. The user should be able to register commands and the paths to the executable
		- Auto check on failure if command is registered. Acts more as a try then.
*/

import (
	"CLIManager/commandManager"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"

	"gopkg.in/yaml.v2"
)

type Runnable struct {
	Name      string   //`yaml:"name"`
	DependsOn []string `yaml:"dependsOn"`
	Command   string   //`yaml:"command"`
	Alias     string   //`yaml:"alias"`
}

type Action struct {
	Name      string     //`yaml:"name"`
	Runnables []Runnable //`yaml:"runnables"`
}

type Actions struct {
	Actions []Action //`yaml:"actions"`
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

var allRunnables = make(map[string]map[string]Runnable)
var channels = make(map[string]chan string)
var shell string
var profilePath string
var config commandManager.Config

func main() {
	args := os.Args
	requireInit, loadedConfig := commandManager.RequireCliSetup()
	hasActioned := commandManager.ParseArgs(args, requireInit)

	// TODO: Test code
	parsedCommands := commandManager.Parser(args)
	commandManager.InterpretCommands(parsedCommands)

	if hasActioned {
		return
	}

	config = loadedConfig
	// actionToRun := os.Args[1]
	allRunnables = setupCommands(config.CommandFiles)

	// runAction(actionToRun)
}

var dependentsMap = make(map[string][]string)

// TODO: Needs to be split up
func setupCommands(filePaths []string) map[string]map[string]Runnable {
	var runnablesStatuses = make(map[string]map[string]Runnable)

	for _, filePath := range filePaths {
		yamlFile, err := ioutil.ReadFile(filePath)

		if err != nil {
			panic(err)
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
	}

	return runnablesStatuses
}

func runAction(actionName string) {
	runnables := allRunnables[actionName]
	feedbackChannel := make(chan string)

	defer close(feedbackChannel)

	var wg sync.WaitGroup

	for _, runnable := range runnables {
		runnablesChannel := make(chan string)
		channels[runnable.Name] = runnablesChannel

		wg.Add(1)
		go runnableRoutine(runnable, feedbackChannel)
	}

	go func() {
		for {
			select {
			case msg := <-feedbackChannel:
				// fmt.Println(msg) // Uncomment to see the actions that have been run
				wg.Done()
				processCompletedRunnable(msg)
			}
		}
	}()

	wg.Wait()
}

func processCompletedRunnable(runnableName string) {
	for _, dependentToNotify := range dependentsMap[runnableName] {
		dependentChannel := channels[dependentToNotify]
		dependentChannel <- runnableName
	}
}

func runnableRoutine(runnable Runnable, feedbackChannel chan string) {
	runnableChannel := channels[runnable.Name]

	if len(runnable.DependsOn) == 0 {
		runCommand(runnable, feedbackChannel)
	} else {
		for {
			select {
			case dependentCompleted := <-runnableChannel:
				selectedIndex := -1

				for index, dependent := range runnable.DependsOn {
					if dependentCompleted == dependent {
						selectedIndex = index
					}
				}

				if selectedIndex > -1 {
					// Remove element from array
					runnable.DependsOn[selectedIndex] = runnable.DependsOn[len(runnable.DependsOn)-1]
					runnable.DependsOn = runnable.DependsOn[:len(runnable.DependsOn)-1]

					if len(runnable.DependsOn) == 0 {
						runCommand(runnable, feedbackChannel)
					}
				}
			}
		}
	}
}

// TODO: Source the users profile. Require them to add it during config
func runCommand(runnable Runnable, feedbackChannel chan string) {
	command := exec.Command(config.Shell, "-c", runnable.Command)
	out, err := command.CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(string(out))
		feedbackChannel <- runnable.Name //Blocking so never ends if you don't listen
		return
	}

	fmt.Println(string(out))
	feedbackChannel <- runnable.Name
}
