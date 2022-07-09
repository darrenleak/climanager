package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v2"
)

type Runnable struct {
	Name      string   //`yaml:"name"`
	DependsOn []string `yaml:"dependsOn"`
	Command   string   //`yaml:"command"`
	Alias     string
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

func main() {
	actionToRun := os.Args[1]
	allRunnables = setupCommands([]string{"/Users/darren/Developer/Projects/CLIManager/test.yml"})

	// TODO: This action will come in from the cli
	runAction(actionToRun)
}

var dependentsMap = make(map[string][]string)

// TODO: Needs to be split up
func setupCommands(filePaths []string) map[string]map[string]Runnable {
	var runnablesStatuses = make(map[string]map[string]Runnable)

	for _, filePath := range filePaths {
		filename, _ := filepath.Abs(filePath)
		yamlFile, err := ioutil.ReadFile(filename)

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

func runCommand(runnable Runnable, feedbackChannel chan string) {
	command := exec.Command("bash", "-c", runnable.Command)
	out, err := command.CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
		feedbackChannel <- runnable.Name //Blocking so never ends if you don't listen
		return
	}

	fmt.Println(string(out))
	feedbackChannel <- runnable.Name
}
