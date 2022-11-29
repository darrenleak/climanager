# bluff

## What and Why
I started building bluff because it seemed like every company I have worked at needed a cli tool that would help them onboard people or just to get people setup quickly, or maybe they needed a bunch of other commands that needed to be orchestrated.

bluff tries to solve this(at a very high level at the moment) by allowing anyone to create a yml file with predefined `actions`. These actions allow one to create/orchestrate commands. These commands are called `runnables` and a `runnable` has a `name`, the `command` to execute and `dependencies`. The `dependencies` will run before the runnable's command. 

## Example yml file

```
actions:
  - name: "helloWorld"
    runnables:
      - name: "hello" # Names cannot have spaces
        command: "echo hello"
      - name: "helloOne" # Names cannot have spaces
        command: "echo helloOne"
      - name: "world" # Names cannot have spaces
        dependsOn:
          - "helloOne"
        command: "echo world"
```

# How to setup bluff
## Building bluff
1. First step, clone this repo. 
2. Build the project. This should be as simple as:
```
go build -o bluff
```

## After building/downloading bluff
1. Create an actions yml file. You can use the example provided above. Or you can test with this file `https://raw.githubusercontent.com/darrenleak/CLIManager/main/actions.yml`
2. Once the project is built, run the following:
```
./bluff --init
```
3.1 For shell, use `zsh`

3.2 The action files, you need to specify the absolute path to your action files.

# Using bluff
Once you have built bluff you can do the following:
```
./CLIMananger helloWorld
```

The `helloWorld` argument is the `action` name from the yml file example above. When you add new actions, you would use that action's name instead.

# bluff Commands
```
--init                Setup the config file by asking a few questions
--shell               Allow you to update the shell setting in the config
--profile             Allow you to update the profile setting in the config
--commandFiles        Allow you to update the command files in the config
--commandFilesAppend  Allow you to append to the command files in the config
--commandFilesRemove  Allow you to remove from the command files in the config
--listCommands        List all the actions
--viewConfig          Print out the current config file
--help                Shows help, what you are seeing now :)
```