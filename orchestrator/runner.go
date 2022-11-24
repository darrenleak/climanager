package orchestrator

import (
	"CLIManager/commandManager"
	"fmt"
	"os/exec"
	"sync"
)

// map[dependent][runnables that depend on dependent]
var dependencyMap = make(map[string][]string)

// map[runnable][list of dependencies]
var runnableDependencies = make(map[string][]string)

var currentConfig commandManager.Config

/*

When a dependency is done, dependencyMap will remove that dependent and get all the dependents and remove the dependent
from the runnableDependencies

*/

func Run(config commandManager.Config, actionToRun string, actions map[string]map[string]Runnable) {
	currentConfig = config
	action := actions[actionToRun]
	actionCompletedChannel := make(chan string)
	makeRunnable := make(chan string)
	defer close(actionCompletedChannel)
	defer close(makeRunnable)

	immediatelyRunnable := runAction(action)
	var waitGroup sync.WaitGroup

	for _, runnableName := range immediatelyRunnable {
		runnable := action[runnableName]

		waitGroup.Add(1)
		go execute(runnable, actionCompletedChannel)
	}

	go func() {
		for {
			select {
			case completedAction := <-actionCompletedChannel:
				newlyRunnableCount := updateAsImmediatelyRunnable(completedAction, makeRunnable)
				waitGroup.Add(newlyRunnableCount)
				waitGroup.Done()
			}
		}
	}()

	go func() {
		for {
			select {
			case newRunnable := <-makeRunnable:
				runnable := action[newRunnable]
				go execute(runnable, actionCompletedChannel)
			}
		}
	}()

	waitGroup.Wait()
}

// Might need a channel
// TODO:Cleanup
func runAction(action map[string]Runnable) []string {
	var immediatelyRunnable = []string{}

	for runnableName := range action {
		runnable := action[runnableName]
		dependencies := runnable.DependsOn

		if len(dependencies) == 0 {
			immediatelyRunnable = append(immediatelyRunnable, runnableName)
			continue
		}

		for _, dependency := range dependencies {
			if dependencyMap[dependency] == nil {
				dependencyMap[dependency] = []string{}
			}

			dependencyMap[dependency] = append(dependencyMap[dependency], runnableName)

			if runnableDependencies[runnableName] == nil {
				runnableDependencies[runnableName] = []string{}
			}

			runnableDependencies[runnableName] = append(runnableDependencies[runnableName], dependency)
		}
	}

	return immediatelyRunnable
}

func updateAsImmediatelyRunnable(runnableName string, makeRunnable chan string) int {
	dependencies := dependencyMap[runnableName]
	count := 0

	// TODO:Cleanup
	for _, dependency := range dependencies {
		currentRunnableDependency := runnableDependencies[dependency]

		for index := range currentRunnableDependency {
			currentRunnableDependency[index] = currentRunnableDependency[len(currentRunnableDependency)-1]
			runnableDependencies[dependency] = currentRunnableDependency[:len(currentRunnableDependency)-1]

			if len(runnableDependencies[dependency]) == 0 {
				count++
				makeRunnable <- dependency
			}
		}
	}

	delete(dependencyMap, runnableName)

	return count
}

func execute(runnable Runnable, actionCompleteChannel chan string) {
	// sourcedCommand := fmt.Sprintf("%s;%s", currentConfig.Profile, runnable.Command)
	command := exec.Command(currentConfig.Shell, "-c", runnable.Command)
	out, err := command.CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(string(out))
		actionCompleteChannel <- runnable.Name
		return
	}

	fmt.Println(string(out))
	actionCompleteChannel <- runnable.Name
}
