---
title: "Getting started"
classes: wide
permalink: /getting-started/
---

## What is `Clipper`?

`Clipper` is an assistant for your terminal, unleashing the power of artificial intelligence to streamline your command line experience.

It can work with OpenAI and with OpenAI-compatible providers through LangChain, which means you can route requests through services such as OpenRouter and choose the model that fits your workflow.

It is already aware of your:
- operating system & distribution
- username, shell & home directory
- preferred editor

And you can also give any supplementary preferences to fine tune your experience.

## Installation

### macOS and Linux

```shell
curl -sS https://raw.githubusercontent.com/ilayaraja97/clipper/main/install.sh | bash
```

- this will detect the proper binary to install for your machine
- and upgrade to the latest stable version if already installed

### Windows

In **PowerShell**:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "Invoke-Expression (Invoke-WebRequest -UseBasicParsing https://raw.githubusercontent.com/ilayaraja97/clipper/main/install.ps1).Content"
```

The Windows installer puts `clipper.exe` in `%LOCALAPPDATA%\Programs\clipper` and updates your user `PATH`. If you use **Git Bash**, you can run the same `install.sh` command as on Linux.

You can also install a binary from the [available releases](https://github.com/ilayaraja97/clipper/releases) on GitHub.

### Uninstall

**macOS / Linux / Git Bash on Windows:**

```shell
curl -sS https://raw.githubusercontent.com/ilayaraja97/clipper/main/uninstall.sh | bash
```

**Windows (PowerShell):** 

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "Invoke-Expression (Invoke-WebRequest -UseBasicParsing https://raw.githubusercontent.com/ilayaraja97/clipper/main/uninstall.ps1).Content"
```

Or remove `%LOCALAPPDATA%\Programs\clipper\clipper.exe` and delete `%USERPROFILE%\.config\clipper.json`.

## Configuration

At first run, `Clipper` will ask you to set it up. You need to get the API key and base URL of the provider you want to use.

This makes it easy to:
- use OpenAI directly
- connect to OpenRouter
- switch to another OpenAI-compatible provider without changing how you use `Clipper`

It will then generate your configuration in the file `~/.config/clipper.json`, with the following structure:

```json
{
  "key": "",
  "model": "",
  "base_url": "",
  "proxy": "",
  "temperature": 0.2,
  "max_tokens": 1000,
  "user_default_prompt_mode": "exec",
  "user_preferences": ""
}
```

## Fine tuning

You can fine tune `Clipper` by editing the settings in the `~/.config/clipper.json` configuration file.

Note that in `REPL` mode, you can press anytime `ctrl+s` to edit settings:
- it will open your editor on the configuration file
- and will hot reload settings changes when you're done.

### Model 

You can use the `model` to configure the AI model you want to use.

```json
{
  "model": "gpt-4"
}
```

With OpenRouter or another compatible gateway, this can be any model exposed through the provider's OpenAI-compatible API.

For example:

```json
{
  "model": "openai/gpt-4o-mini"
}
```

### Base URL

You can use `base_url` if you want to send requests to an OpenAI-compatible endpoint instead of the default OpenAI API:

```json
{
  "base_url": "https://openrouter.ai/api/v1"
}
```

This is useful for providers and gateways that expose an OpenAI-compatible API surface.

Combined with `model`, this lets `Clipper` work with a much wider range of LLMs while keeping the same terminal experience.

### Preferences

You can use the `user_preferences` to express any preferences in your natural language:

```json
{
  "user_preferences": "I want you to talk like an humorist, and I want you to always add the -y flag when I use dnf"
}
```

`Clipper` will take them into account.
