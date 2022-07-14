package main

import "fmt"

type KeyValueStringPair struct {
	Key   string
	Value string
}

type ConfigProperties struct {
	Shell            string
	ShellProfilePath string
}

type Config struct {
	properties []ConfigProperties
}

func writeConfig(keyValues []KeyValueStringPair) {
	// Convert KeyValueStringPair to Config
	fmt.Println("Write Config")
}

func readConfig() {
	fmt.Println("Read Config")
}
