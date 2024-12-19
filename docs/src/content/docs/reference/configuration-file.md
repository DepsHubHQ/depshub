---
title: depshub.yaml
description: A referee for the depshub.yaml configuration file
---

All the configuration for DepsHub is stored in a file called `depshub.yaml`.
This file is used to configure the behavior of the tool and to specify the dependencies to be checked.

`depshub.yaml` **isn't required** by default but it's recommended to have one in your repository to configure the behavior of the tool.

## Basic example

Here is a basic example of `depshub.yaml` file that disables `no-unstable` rule for all the manifest files:

```yaml
version: 1                   # The version of the configuration file
manifest_files:              # The list of manifest files to check
  - filter: "**"             # The glob pattern to specify the files to check
    rules:                   # The list of rules to apply to the manifest file
      - name: "no-unstable"  # The name of the rule to apply
        disabled: true       # Use this option to disable the rule. See more options below.
```

## Configuration options

### `version`
Specifies the version of the configuration file. The current version is `1`.

### `ignore`
Specifies the list of files to ignore. The files specified in this list will be ignored by the tool.
Use the glob pattern to specify the files to ignore. The default value is `[]`.

```yaml
version: 1
ignore:
  - "**/test-requirements.txt"
```

### `manifest_files`
Specifies the list of manifest files to check.

#### `filter`
Use a glob pattern to specify the files to check. For example, `**/requirements.txt` will match all the `requirements.txt` files in the repository. `**` matches all the manifest files.

Example:

#### `rules`
The list of rules to apply to the manifest file.

##### `name`
The name of the rule to apply. The list of available rules can be found in the [rules](/reference/rules) section.

##### `disabled`
Use this option to disable the rule. The default value is `false`.

##### `value`
Use this option to specify the value for the rule. The value is optional and depends on the rule. See the [rules](/reference/rules) section for more information.

Example:
```yaml
version: 1
manifest_files:
  - filter: "**"
    rules:
      - name: "max-libyear"
        value: 100
```

#####  `level`
Use this option to specify the level of the rule. The level can be `error` or `warning`. The default value depends on the rule.

Example:

```yaml
version: 1
manifest_files:
  - filter: "**"
    rules:
      - name: "max-libyear"
        level: "warning"
```

#### `packages`
An array of package names to check. If this option is specified, only the specified packages will be checked.

Example: 

```yaml
version: 1
manifest_files:
  - filter: "**/requirements.txt"
    packages: ["requests", "flask"]
    rules:
      - name: "no-unstable"
        disabled: true
```

## Further reading

- Read about all the supported [rules](/reference/rules) that can be applied to the manifest files.
