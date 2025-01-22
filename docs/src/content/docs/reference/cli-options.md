---
title: CLI
description: A reference page about all the available CLI options for the `depshub` command.
---

## Commands

### `depshub lint`

Runs the linter on the project. The linter is responsible for checking the project for any dependency issues.

### `depshub help`

Shows the help message.

### `depshub version`

Shows the version of DepsHub.

## Flags

### `--config`

Path to the configuration file. If not provided, DepsHub will try to find a configuration file in the current directory.
If not found, it will use the default configuration.

Default value: `.`

Example usage:

```sh
depshub lint . --config ./path/to/config.json
```
