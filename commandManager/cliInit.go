package commandManager

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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

	fmt.Println(inputValues)
}
