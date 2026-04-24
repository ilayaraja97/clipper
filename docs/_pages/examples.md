---
title: "Examples"
classes: wide
permalink: /examples/
---

> The examples below use an older command name, but the current binary is `clipper`.

`Clipper` now works with more than a single hosted model. If your provider exposes an OpenAI-compatible API, you can keep the same workflows while choosing the model through your configuration, including setups routed through OpenRouter.

## CLI mode

`Clipper` can be used in `CLI` mode to streamline your terminal workflow.

That means the same command-generation and chat flows can now be powered by the model you prefer for speed, price, or reasoning quality.

For example, this is how it can assist you while developing an application:

![](https://raw.githubusercontent.com/ilayaraja97/clipper/main/docs/_assets/dev.gif)

You can also pipe any input and make it work with it:

![](https://raw.githubusercontent.com/ilayaraja97/clipper/main/docs/_assets/pipe.gif)

## REPL mode

`Clipper` can be used in `REPL` mode to offer interactive prompts and chain instructions.

In practice, this makes it easy to keep one familiar terminal assistant experience while swapping the backend model behind it.

### Command lines

An example on how it can help you to manage your system:

![](https://raw.githubusercontent.com/ilayaraja97/clipper/main/docs/_assets/system.gif)

Or help you to manage your packages:

![](https://raw.githubusercontent.com/ilayaraja97/clipper/main/docs/_assets/pkg.gif)

It can also go way further, for example help you to manage docker resources:

![](https://raw.githubusercontent.com/ilayaraja97/clipper/main/docs/_assets/docker.gif)

Or even help you while using a k8s cluster:

![](https://raw.githubusercontent.com/ilayaraja97/clipper/main/docs/_assets/k8s.gif)

### Chat

`Clipper` is not made to just generate command lines. You can also engage in a discussion with it about any topics.

This is especially helpful when you want to use one provider for shell-heavy tasks and another for longer debugging or explanation sessions.

![](https://raw.githubusercontent.com/ilayaraja97/clipper/main/docs/_assets/chat.gif)
