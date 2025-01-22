---
title: Examples
---

## Basic usage

After you have [installed](/installation) DepsHub you can run it with the following command:

```bash
depshub lint .
```

This will run DepsHub on the current directory and all its subdirectories. It respects the `.gitignore` file and will not lint files that are ignored by Git.

## Ignoring files

To further configure DepsHub, you can create a configuration file. This file should be named `depshub.yml`. DepsHub automatically reads this file when it is present in the current directory.

Now, let's say you want to ignore some specific directory. You can do this by adding the following to your `depshub.yml` file:

```yaml
version: 1
ignore:
  - "**/testdata/**"
```

## Configuring specific rules

If you want to configure a specific rule for a specific manifest file, you can do this by adding the following to your `depshub.yml` file:

```yaml
version: 1
ignore:
  - "**/testdata/**"
manifest_files:
  - filter: "**" # This rule applies to all manifest files
    rules:
      - name: "allowed-licenses"
        value: ["", "MIT", "Apache-2.0", "BSD-3-Clause"]
  - filter: "**/docs/package.json"
    packages: ["@astrojs/check", "@astrojs/starlight", "sharp"]
    rules:
      - name: "no-unstable"
        disabled: true
```

See all the available configuration options in the [configuration reference](/reference/configuration-file/).
