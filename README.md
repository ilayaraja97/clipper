# Clipper - AI powered terminal assistant

> Forked from https://github.com/ekkinox/yai/

[![build](https://github.com/ilayaraja97/clipper/actions/workflows/build.yml/badge.svg)](https://github.com/ilayaraja97/clipper/actions/workflows/build.yml)
[![release](https://github.com/ilayaraja97/clipper/actions/workflows/release.yml/badge.svg)](https://github.com/ilayaraja97/clipper/actions/workflows/release.yml)

> Unleash the power of artificial intelligence to streamline your command line experience.

![Intro](docs/_assets/intro.gif)

## What is Clipper?

`Clipper` is an assistant for your terminal, using [OpenAI ChatGPT](https://chat.openai.com/) to build and run commands for you. You just need to describe them in your everyday language, it will take care of the rest.

You have any questions on random topics in mind? You can also ask `Clipper`, and get the power of AI without leaving `/home`.

It is already aware of your:
- operating system & distribution
- username, shell & home directory
- preferred editor

And you can also give any supplementary preferences to fine tune your experience.

## Documentation

A complete documentation is available at [https://ilayaraja97.github.io/clipper/](https://ilayaraja97.github.io/clipper/).

## Quick start

To install `Clipper`, simply run:

```shell
curl -sS https://raw.githubusercontent.com/ilayaraja97/clipper/main/install.sh | bash
```

At first run, it will ask you for an [OpenAI API key](https://platform.openai.com/account/api-keys), and use it to create the configuration file in `~/.config/clipper.json`.

See [documentation](https://ilayaraja97.github.io/clipper/getting-started/#configuration) for more information.
