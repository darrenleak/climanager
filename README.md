# CLIManager

## What and Why
I started building CLIManager because it seemed like every company I have worked at needed a cli tool that would help them onboard people or just to get people setup quickly, or maybe they needed a bunch of other commands that needed to be orchestrated.

CLIManager tries to solve this(at a very high level at the moment) by allowing anyone to create a yml file with predefined `actions`. These actions allow one to create/orchestrate commands. These commands are called `runnables` and a `runnable` has a `name`, the `command` to execute and `dependencies`. The `dependencies` will run before the runnable's command. 

## Example yml file

```
actions:
  - name: "helloWorld"
    runnables:
      - name: "hello" # Names cannot have spaces
        command: "echo hello"
      - name: "world" # Names cannot have spaces
        dependsOn:
          - "hello"
        command: "echo world"
```