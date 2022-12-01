package main

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
