---
title: "Usage"
classes: wide
permalink: /usage/
---

## CLI mode

> CLI mode is made to be integrated in your command lines workflow.

You can perform a single run:

```shell
clipper list all processes listening on port 8080
```

You can ask for a command line generation, enforcing `🚀 exec` prompt mode usage with `-e`:

```shell
clipper -e show the disk usage of my docker resources
```

You can ask any question, enforcing `💬 chat` prompt mode usage with `-c`:

```shell
clipper -c generate me a go application example using fiber
```

You can also `pipe` input that will be taken into account in your request:

```shell
cat some_script.go | clipper -c generate unit tests
```

Or even:

```shell
cat error.log | clipper -c explain what is wrong here
```

## REPL mode

> REPL mode is made to work in an interactive way.

Just run:

```shell
clipper
```

This will open a [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) interface, with 2 types of prompts (use `tab` to switch)

- `🚀 exec`: will generate a command line to execute for what you're asking
- `💬 chat`: will engage in a discussion to help you the best way possible

You also can use the following **keyboard shortcuts**:
- `↑` `↓`  : Navigate in history                                 
- `tab`    : Switch between `🚀 exec` and `💬 chat` prompt modes 
- `ctrl+h` : Show help                                           
- `ctrl+s` : Edit settings                                       
- `ctrl+r` : Clear terminal and reset discussion history         
- `ctrl+l` : Clear terminal but keep discussion history          
- `ctrl+c` : Exit or interrupt command execution                 
